package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"UMSRMS/internal/config"
	"UMSRMS/internal/handler"
	"UMSRMS/internal/middleware"
	"UMSRMS/internal/models"
	"UMSRMS/internal/repository"
	"UMSRMS/internal/service"
	"UMSRMS/internal/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "UMSRMS/docs"
)

// @title University Maintenance Service API
// @version 1.0
// @description API documentation for the University Maintenance & Service Request Management System.
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.

// ping godoc
// @Summary Health check
// @Description Returns pong if service is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func main() {
	seedOnStartup := seedFlag()
	flag.Parse()

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	gin.SetMode(cfg.GinMode)
	utils.RegisterValidationTagNames()

	db, err := config.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer db.Close()

	roleRepo := repository.NewRoleRepository(db)
	if *seedOnStartup {
		insertedCount, err := roleRepo.SeedDefaultRoles(context.Background())
		if err != nil {
			log.Fatalf("failed to seed default roles: %v", err)
		}
		if insertedCount > 0 {
			log.Printf("role seed completed: inserted %d default role(s)", insertedCount)
		} else {
			log.Printf("role seed skipped: default roles already exist")
		}
	} else {
		log.Printf("role seed disabled: start with --seed=true or SEED_ON_STARTUP=true")
	}

	jwtManager, err := utils.NewJWTManager(cfg)
	if err != nil {
		log.Fatalf("failed to initialize jwt manager: %v", err)
	}
	tokenBanList := utils.NewTokenBanList()

	userRepo := repository.NewUserRepository(db)
	requestRepo := repository.NewServiceRequestRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	reportRepo := repository.NewReportRepository(db)

	handlers := appHandlers{
		auth:         handler.NewAuthHandler(service.NewAuthService(userRepo, roleRepo, jwtManager, tokenBanList)),
		request:      handler.NewServiceRequestHandler(service.NewServiceRequestService(requestRepo, notificationRepo, cfg)),
		category:     handler.NewCategoryHandler(service.NewCategoryService(categoryRepo)),
		user:         handler.NewUserHandler(service.NewUserService(userRepo, roleRepo)),
		notification: handler.NewNotificationHandler(service.NewNotificationService(notificationRepo)),
		audit:        handler.NewAuditHandler(service.NewAuditService(auditRepo)),
		report:       handler.NewReportHandler(service.NewReportService(reportRepo)),
	}

	requireAuth := middleware.RequireAuth(jwtManager, tokenBanList)
	auditLog := middleware.AuditLog(auditRepo)

	router := routes(cfg, handlers, requireAuth, auditLog)

	if err := router.Run(cfg.Address()); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func seedFlag() *bool {
	defaultSeed := false
	if raw, ok := os.LookupEnv("SEED_ON_STARTUP"); ok {
		if parsed, err := strconv.ParseBool(strings.TrimSpace(raw)); err == nil {
			defaultSeed = parsed
		}
	}
	return flag.Bool("seed", defaultSeed, "run startup seed data")
}

// appHandlers bundles the HTTP handlers wired into the router.
type appHandlers struct {
	auth         *handler.AuthHandler
	request      *handler.ServiceRequestHandler
	category     *handler.CategoryHandler
	user         *handler.UserHandler
	notification *handler.NotificationHandler
	audit        *handler.AuditHandler
	report       *handler.ReportHandler
}

func routes(cfg *config.EnvConfig, h appHandlers, requireAuth, auditLog gin.HandlerFunc) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORS(cfg))
	router.Use(middleware.RateLimit(cfg))

	api := router.Group("api/v1")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	authRoutes := api.Group("/auth")
	authRoutes.Use(middleware.RateLimitPerMinute(5))
	{
		authRoutes.POST("/register", h.auth.Register)
		authRoutes.POST("/login", h.auth.Login)
		authRoutes.GET("/registration-data", h.auth.RegistrationData)
		authRoutes.POST("/logout", requireAuth, h.auth.Logout)
		authRoutes.GET("/me", requireAuth, h.auth.Me)
	}

	requestRoutes := api.Group("/requests")
	requestRoutes.Use(requireAuth, auditLog)
	{
		requestRoutes.POST("", middleware.RequireRoles(models.RoleStudentStaff), h.request.Create)
		requestRoutes.GET("", h.request.List)
		requestRoutes.GET("/:id", h.request.GetByID)
		requestRoutes.PUT("/:id/status", middleware.RequireRoles(models.RoleMaintenanceOfficer, models.RoleAdmin), h.request.UpdateStatus)
		requestRoutes.POST("/:id/assign", middleware.RequireRoles(models.RoleAdmin), h.request.Assign)
		requestRoutes.DELETE("/:id", middleware.RequireRoles(models.RoleAdmin), h.request.Delete)
		requestRoutes.POST("/:id/evidence", middleware.RequireRoles(models.RoleStudentStaff), h.request.UploadEvidence)
	}

	userRoutes := api.Group("/users")
	userRoutes.Use(requireAuth, middleware.RequireRoles(models.RoleAdmin), auditLog)
	{
		userRoutes.GET("", h.user.List)
		userRoutes.GET("/officers", h.user.ListOfficers)
		userRoutes.PUT("/:id/role", h.user.UpdateRole)
		userRoutes.DELETE("/:id", h.user.Delete)
	}

	notificationRoutes := api.Group("/notifications")
	notificationRoutes.Use(requireAuth)
	{
		notificationRoutes.GET("", h.notification.List)
		notificationRoutes.PUT("/read-all", h.notification.MarkAllRead)
		notificationRoutes.PUT("/:id/read", h.notification.MarkRead)
	}

	api.GET("/categories", requireAuth, h.category.List)
	api.GET("/reports/summary", requireAuth, middleware.RequireRoles(models.RoleAdmin), h.report.Summary)
	api.GET("/audit-logs", requireAuth, middleware.RequireRoles(models.RoleAdmin), h.audit.List)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
