package repository

import squirrel "github.com/Masterminds/squirrel"

// SQL builds PostgreSQL queries using Squirrel with $1, $2 placeholders.
var SQL = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
