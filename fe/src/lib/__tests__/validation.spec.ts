import { describe, it, expect } from 'vitest'
import { loginSchema, registerSchema, requestSchema, validate } from '../validation'

describe('validation schemas', () => {
  it('accepts valid login input', async () => {
    expect(await validate(loginSchema, { email: 'a@b.com', password: 'secret' })).toBeNull()
  })

  it('rejects invalid email and short password', async () => {
    const errors = await validate(loginSchema, { email: 'not-an-email', password: '123' })
    expect(errors).not.toBeNull()
    expect(errors?.email).toBeTruthy()
    expect(errors?.password).toBeTruthy()
  })

  it('requires matching passwords on register', async () => {
    const errors = await validate(registerSchema, {
      fullName: 'Jane Doe',
      email: 'jane@miva.edu',
      password: 'secret1',
      confirmPassword: 'secret2',
    })
    expect(errors?.confirmPassword).toBe('Passwords do not match')
  })

  it('validates a service request', async () => {
    const ok = await validate(requestSchema, {
      title: 'Broken socket',
      categoryId: 'c-elec',
      location: 'Room 12',
      priority: 'high',
      description: 'The socket sparks when used.',
    })
    expect(ok).toBeNull()
  })
})
