import * as yup from 'yup'

export const loginSchema = yup.object({
  email: yup.string().required('Email is required').email('Enter a valid email'),
  password: yup.string().required('Password is required').min(6, 'At least 6 characters'),
})

export const registerSchema = yup.object({
  fullName: yup.string().required('Full name is required').min(2, 'Too short'),
  email: yup.string().required('Email is required').email('Enter a valid email'),
  roleId: yup.string().required('Select an account type'),
  phone: yup
    .string()
    .optional()
    .matches(/^$|^[+\d][\d\s-]{6,}$/, 'Enter a valid phone number'),
  password: yup.string().required('Password is required').min(6, 'At least 6 characters'),
  confirmPassword: yup
    .string()
    .required('Please confirm your password')
    .oneOf([yup.ref('password')], 'Passwords do not match'),
})

export const requestSchema = yup.object({
  title: yup.string().required('Title is required').min(5, 'Give a clearer title'),
  categoryId: yup.string().required('Select a category'),
  location: yup.string().required('Location is required'),
  priority: yup
    .string()
    .required('Select a priority')
    .oneOf(['low', 'medium', 'high'], 'Invalid priority'),
  description: yup.string().required('Description is required').min(10, 'Add more detail'),
})

/**
 * Validate a yup schema and return a flat `{ field: message }` map of errors,
 * or `null` when the data is valid.
 */
export async function validate<T extends yup.AnyObject>(
  schema: yup.ObjectSchema<T>,
  data: unknown,
): Promise<Record<string, string> | null> {
  try {
    await schema.validate(data, { abortEarly: false })
    return null
  } catch (err) {
    const errors: Record<string, string> = {}
    if (err instanceof yup.ValidationError) {
      for (const e of err.inner) {
        if (e.path && !errors[e.path]) errors[e.path] = e.message
      }
    }
    return errors
  }
}
