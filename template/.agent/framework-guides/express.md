# Express.js Framework Guide

> **Framework**: Express.js 4.x/5.x
> **Language**: TypeScript/JavaScript
> **Use Cases**: REST APIs, Web Servers, Middleware, Microservices

---

## Overview

Express.js is a minimal and flexible Node.js web application framework that provides a robust set of features for web and mobile applications. It's the most popular Node.js framework for building APIs and web servers.

---

## Project Setup

### Create New Project
```bash
mkdir my-api
cd my-api
npm init -y

# Install dependencies
npm install express cors helmet morgan dotenv

# TypeScript setup
npm install -D typescript @types/express @types/node @types/cors @types/morgan
npm install -D ts-node-dev

# Initialize TypeScript
npx tsc --init
```

### Project Structure
```
my-api/
├── src/
│   ├── controllers/         # Route handlers
│   │   └── user.controller.ts
│   ├── middlewares/         # Custom middleware
│   │   ├── auth.middleware.ts
│   │   ├── error.middleware.ts
│   │   └── validate.middleware.ts
│   ├── routes/              # Route definitions
│   │   ├── index.ts
│   │   └── user.routes.ts
│   ├── services/            # Business logic
│   │   └── user.service.ts
│   ├── repositories/        # Data access
│   │   └── user.repository.ts
│   ├── models/              # Data models
│   │   └── user.model.ts
│   ├── types/               # TypeScript types
│   │   └── index.ts
│   ├── utils/               # Utility functions
│   │   └── logger.ts
│   ├── config/              # Configuration
│   │   └── index.ts
│   ├── app.ts               # Express app setup
│   └── server.ts            # Server entry point
├── tests/
│   └── user.test.ts
├── .env
├── .env.example
├── package.json
├── tsconfig.json
└── README.md
```

### Configuration Files

```json
// tsconfig.json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "lib": ["ES2020"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true,
    "baseUrl": "./src",
    "paths": {
      "@/*": ["./*"]
    }
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist"]
}
```

```json
// package.json scripts
{
  "scripts": {
    "dev": "ts-node-dev --respawn --transpile-only src/server.ts",
    "build": "tsc",
    "start": "node dist/server.js",
    "test": "jest",
    "lint": "eslint src/**/*.ts"
  }
}
```

---

## Application Setup

### Basic App Configuration
```typescript
// src/app.ts
import express, { Application, Request, Response, NextFunction } from 'express';
import cors from 'cors';
import helmet from 'helmet';
import morgan from 'morgan';
import { errorHandler } from './middlewares/error.middleware';
import { notFoundHandler } from './middlewares/notFound.middleware';
import routes from './routes';

const app: Application = express();

// Security middleware
app.use(helmet());
app.use(cors({
  origin: process.env.CORS_ORIGIN || '*',
  credentials: true,
}));

// Request parsing
app.use(express.json({ limit: '10mb' }));
app.use(express.urlencoded({ extended: true }));

// Logging
if (process.env.NODE_ENV !== 'test') {
  app.use(morgan('combined'));
}

// Health check
app.get('/health', (req: Request, res: Response) => {
  res.json({ status: 'ok', timestamp: new Date().toISOString() });
});

// API routes
app.use('/api/v1', routes);

// Error handling
app.use(notFoundHandler);
app.use(errorHandler);

export default app;
```

### Server Entry Point
```typescript
// src/server.ts
import 'dotenv/config';
import app from './app';
import { logger } from './utils/logger';

const PORT = process.env.PORT || 3000;

const server = app.listen(PORT, () => {
  logger.info(`Server running on port ${PORT}`);
});

// Graceful shutdown
const shutdown = () => {
  logger.info('Shutting down gracefully...');
  server.close(() => {
    logger.info('Server closed');
    process.exit(0);
  });

  // Force close after 10s
  setTimeout(() => {
    logger.error('Forced shutdown');
    process.exit(1);
  }, 10000);
};

process.on('SIGTERM', shutdown);
process.on('SIGINT', shutdown);

// Handle unhandled rejections
process.on('unhandledRejection', (reason: Error) => {
  logger.error('Unhandled Rejection:', reason);
  throw reason;
});

process.on('uncaughtException', (error: Error) => {
  logger.error('Uncaught Exception:', error);
  process.exit(1);
});
```

---

## Routing

### Route Organization
```typescript
// src/routes/index.ts
import { Router } from 'express';
import userRoutes from './user.routes';
import productRoutes from './product.routes';
import authRoutes from './auth.routes';

const router = Router();

router.use('/auth', authRoutes);
router.use('/users', userRoutes);
router.use('/products', productRoutes);

export default router;

// src/routes/user.routes.ts
import { Router } from 'express';
import { UserController } from '../controllers/user.controller';
import { authMiddleware } from '../middlewares/auth.middleware';
import { validate } from '../middlewares/validate.middleware';
import { createUserSchema, updateUserSchema } from '../schemas/user.schema';

const router = Router();
const controller = new UserController();

router.get('/', controller.getAll);
router.get('/:id', controller.getById);
router.post('/', validate(createUserSchema), controller.create);
router.put('/:id', authMiddleware, validate(updateUserSchema), controller.update);
router.delete('/:id', authMiddleware, controller.delete);

export default router;
```

---

## Controllers

### Controller Pattern
```typescript
// src/controllers/user.controller.ts
import { Request, Response, NextFunction } from 'express';
import { UserService } from '../services/user.service';
import { CreateUserDTO, UpdateUserDTO } from '../types';
import { AppError } from '../utils/errors';

export class UserController {
  private userService: UserService;

  constructor() {
    this.userService = new UserService();
  }

  getAll = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { page = 1, limit = 10, sort = 'createdAt' } = req.query;

      const result = await this.userService.findAll({
        page: Number(page),
        limit: Number(limit),
        sort: String(sort),
      });

      res.json({
        data: result.users,
        meta: {
          page: result.page,
          limit: result.limit,
          total: result.total,
          totalPages: result.totalPages,
        },
      });
    } catch (error) {
      next(error);
    }
  };

  getById = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { id } = req.params;
      const user = await this.userService.findById(id);

      if (!user) {
        throw new AppError('User not found', 404);
      }

      res.json({ data: user });
    } catch (error) {
      next(error);
    }
  };

  create = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const data: CreateUserDTO = req.body;
      const user = await this.userService.create(data);

      res.status(201).json({ data: user });
    } catch (error) {
      next(error);
    }
  };

  update = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { id } = req.params;
      const data: UpdateUserDTO = req.body;

      const user = await this.userService.update(id, data);

      if (!user) {
        throw new AppError('User not found', 404);
      }

      res.json({ data: user });
    } catch (error) {
      next(error);
    }
  };

  delete = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { id } = req.params;
      await this.userService.delete(id);

      res.status(204).send();
    } catch (error) {
      next(error);
    }
  };
}
```

---

## Middleware

### Authentication Middleware
```typescript
// src/middlewares/auth.middleware.ts
import { Request, Response, NextFunction } from 'express';
import jwt from 'jsonwebtoken';
import { AppError } from '../utils/errors';

interface JwtPayload {
  userId: string;
  email: string;
}

declare global {
  namespace Express {
    interface Request {
      user?: JwtPayload;
    }
  }
}

export const authMiddleware = async (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  try {
    const authHeader = req.headers.authorization;

    if (!authHeader?.startsWith('Bearer ')) {
      throw new AppError('No token provided', 401);
    }

    const token = authHeader.split(' ')[1];
    const decoded = jwt.verify(
      token,
      process.env.JWT_SECRET!
    ) as JwtPayload;

    req.user = decoded;
    next();
  } catch (error) {
    if (error instanceof jwt.JsonWebTokenError) {
      next(new AppError('Invalid token', 401));
    } else {
      next(error);
    }
  }
};

// Role-based access control
export const requireRole = (...roles: string[]) => {
  return (req: Request, res: Response, next: NextFunction) => {
    if (!req.user) {
      return next(new AppError('Authentication required', 401));
    }

    // Assuming user has a role property
    const userRole = (req.user as any).role;
    if (!roles.includes(userRole)) {
      return next(new AppError('Insufficient permissions', 403));
    }

    next();
  };
};
```

### Validation Middleware
```typescript
// src/middlewares/validate.middleware.ts
import { Request, Response, NextFunction } from 'express';
import { ZodSchema, ZodError } from 'zod';
import { AppError } from '../utils/errors';

export const validate = (schema: ZodSchema) => {
  return (req: Request, res: Response, next: NextFunction) => {
    try {
      schema.parse({
        body: req.body,
        query: req.query,
        params: req.params,
      });
      next();
    } catch (error) {
      if (error instanceof ZodError) {
        const errors = error.errors.map((e) => ({
          field: e.path.join('.'),
          message: e.message,
        }));
        next(new AppError('Validation failed', 400, errors));
      } else {
        next(error);
      }
    }
  };
};

// src/schemas/user.schema.ts
import { z } from 'zod';

export const createUserSchema = z.object({
  body: z.object({
    email: z.string().email('Invalid email format'),
    password: z.string().min(8, 'Password must be at least 8 characters'),
    name: z.string().min(2, 'Name must be at least 2 characters'),
  }),
});

export const updateUserSchema = z.object({
  body: z.object({
    email: z.string().email().optional(),
    name: z.string().min(2).optional(),
  }),
  params: z.object({
    id: z.string().uuid('Invalid user ID'),
  }),
});
```

### Error Handling Middleware
```typescript
// src/middlewares/error.middleware.ts
import { Request, Response, NextFunction } from 'express';
import { AppError } from '../utils/errors';
import { logger } from '../utils/logger';

export const errorHandler = (
  error: Error,
  req: Request,
  res: Response,
  next: NextFunction
) => {
  logger.error('Error:', {
    message: error.message,
    stack: error.stack,
    path: req.path,
    method: req.method,
  });

  if (error instanceof AppError) {
    return res.status(error.statusCode).json({
      status: 'error',
      message: error.message,
      ...(error.errors && { errors: error.errors }),
    });
  }

  // Mongoose validation error
  if (error.name === 'ValidationError') {
    return res.status(400).json({
      status: 'error',
      message: 'Validation error',
      errors: error.message,
    });
  }

  // Default error
  res.status(500).json({
    status: 'error',
    message: process.env.NODE_ENV === 'production'
      ? 'Internal server error'
      : error.message,
  });
};

// src/middlewares/notFound.middleware.ts
export const notFoundHandler = (req: Request, res: Response) => {
  res.status(404).json({
    status: 'error',
    message: `Route ${req.originalUrl} not found`,
  });
};

// src/utils/errors.ts
export class AppError extends Error {
  constructor(
    message: string,
    public statusCode: number = 500,
    public errors?: any[]
  ) {
    super(message);
    this.name = 'AppError';
    Error.captureStackTrace(this, this.constructor);
  }
}
```

### Rate Limiting
```typescript
// src/middlewares/rateLimit.middleware.ts
import rateLimit from 'express-rate-limit';
import RedisStore from 'rate-limit-redis';
import { createClient } from 'redis';

// Basic rate limiter
export const apiLimiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 100, // Limit each IP to 100 requests per windowMs
  message: {
    status: 'error',
    message: 'Too many requests, please try again later',
  },
  standardHeaders: true,
  legacyHeaders: false,
});

// Stricter limiter for auth routes
export const authLimiter = rateLimit({
  windowMs: 60 * 60 * 1000, // 1 hour
  max: 5, // 5 failed attempts per hour
  message: {
    status: 'error',
    message: 'Too many login attempts, please try again later',
  },
});

// Redis-based limiter for distributed systems
export const createRedisLimiter = async () => {
  const client = createClient({ url: process.env.REDIS_URL });
  await client.connect();

  return rateLimit({
    windowMs: 15 * 60 * 1000,
    max: 100,
    store: new RedisStore({
      sendCommand: (...args: string[]) => client.sendCommand(args),
    }),
  });
};
```

---

## Services

### Service Layer Pattern
```typescript
// src/services/user.service.ts
import { UserRepository } from '../repositories/user.repository';
import { CreateUserDTO, UpdateUserDTO, User, PaginationResult } from '../types';
import { hashPassword } from '../utils/crypto';
import { AppError } from '../utils/errors';

interface FindAllOptions {
  page: number;
  limit: number;
  sort: string;
}

export class UserService {
  private repository: UserRepository;

  constructor() {
    this.repository = new UserRepository();
  }

  async findAll(options: FindAllOptions): Promise<PaginationResult<User>> {
    const { page, limit, sort } = options;
    const skip = (page - 1) * limit;

    const [users, total] = await Promise.all([
      this.repository.findAll({ skip, limit, sort }),
      this.repository.count(),
    ]);

    return {
      users,
      page,
      limit,
      total,
      totalPages: Math.ceil(total / limit),
    };
  }

  async findById(id: string): Promise<User | null> {
    return this.repository.findById(id);
  }

  async findByEmail(email: string): Promise<User | null> {
    return this.repository.findByEmail(email);
  }

  async create(data: CreateUserDTO): Promise<User> {
    // Check if email exists
    const existing = await this.findByEmail(data.email);
    if (existing) {
      throw new AppError('Email already in use', 409);
    }

    // Hash password
    const hashedPassword = await hashPassword(data.password);

    return this.repository.create({
      ...data,
      password: hashedPassword,
    });
  }

  async update(id: string, data: UpdateUserDTO): Promise<User | null> {
    const user = await this.findById(id);
    if (!user) {
      return null;
    }

    // Check email uniqueness if changing
    if (data.email && data.email !== user.email) {
      const existing = await this.findByEmail(data.email);
      if (existing) {
        throw new AppError('Email already in use', 409);
      }
    }

    return this.repository.update(id, data);
  }

  async delete(id: string): Promise<void> {
    const user = await this.findById(id);
    if (!user) {
      throw new AppError('User not found', 404);
    }

    await this.repository.delete(id);
  }
}
```

---

## Testing

### Integration Tests with Supertest
```typescript
// tests/user.test.ts
import request from 'supertest';
import app from '../src/app';
import { db } from '../src/config/database';

describe('User API', () => {
  beforeAll(async () => {
    await db.connect();
  });

  afterAll(async () => {
    await db.disconnect();
  });

  beforeEach(async () => {
    await db.clear();
  });

  describe('GET /api/v1/users', () => {
    it('should return empty array when no users', async () => {
      const response = await request(app)
        .get('/api/v1/users')
        .expect('Content-Type', /json/)
        .expect(200);

      expect(response.body.data).toEqual([]);
    });

    it('should return users with pagination', async () => {
      // Create test users
      await createTestUsers(15);

      const response = await request(app)
        .get('/api/v1/users?page=1&limit=10')
        .expect(200);

      expect(response.body.data).toHaveLength(10);
      expect(response.body.meta.total).toBe(15);
      expect(response.body.meta.totalPages).toBe(2);
    });
  });

  describe('POST /api/v1/users', () => {
    it('should create a new user', async () => {
      const userData = {
        email: 'test@example.com',
        password: 'password123',
        name: 'Test User',
      };

      const response = await request(app)
        .post('/api/v1/users')
        .send(userData)
        .expect(201);

      expect(response.body.data.email).toBe(userData.email);
      expect(response.body.data.name).toBe(userData.name);
      expect(response.body.data.password).toBeUndefined();
    });

    it('should return 400 for invalid email', async () => {
      const response = await request(app)
        .post('/api/v1/users')
        .send({
          email: 'invalid-email',
          password: 'password123',
          name: 'Test',
        })
        .expect(400);

      expect(response.body.status).toBe('error');
    });

    it('should return 409 for duplicate email', async () => {
      await request(app)
        .post('/api/v1/users')
        .send({
          email: 'test@example.com',
          password: 'password123',
          name: 'Test',
        });

      const response = await request(app)
        .post('/api/v1/users')
        .send({
          email: 'test@example.com',
          password: 'password456',
          name: 'Another',
        })
        .expect(409);

      expect(response.body.message).toContain('already in use');
    });
  });

  describe('GET /api/v1/users/:id', () => {
    it('should return user by id', async () => {
      const user = await createTestUser();

      const response = await request(app)
        .get(`/api/v1/users/${user.id}`)
        .expect(200);

      expect(response.body.data.id).toBe(user.id);
    });

    it('should return 404 for non-existent user', async () => {
      await request(app)
        .get('/api/v1/users/non-existent-id')
        .expect(404);
    });
  });
});
```

---

## Database Integration

### Prisma Setup
```typescript
// src/config/database.ts
import { PrismaClient } from '@prisma/client';

const prisma = new PrismaClient({
  log: process.env.NODE_ENV === 'development' ? ['query', 'error', 'warn'] : ['error'],
});

export { prisma as db };

// src/repositories/user.repository.ts
import { db } from '../config/database';
import { User, CreateUserDTO, UpdateUserDTO } from '../types';

export class UserRepository {
  async findAll(options: { skip: number; limit: number; sort: string }) {
    return db.user.findMany({
      skip: options.skip,
      take: options.limit,
      orderBy: { [options.sort]: 'desc' },
      select: {
        id: true,
        email: true,
        name: true,
        createdAt: true,
        updatedAt: true,
      },
    });
  }

  async findById(id: string) {
    return db.user.findUnique({
      where: { id },
      select: {
        id: true,
        email: true,
        name: true,
        createdAt: true,
        updatedAt: true,
      },
    });
  }

  async findByEmail(email: string) {
    return db.user.findUnique({
      where: { email },
    });
  }

  async create(data: CreateUserDTO) {
    return db.user.create({
      data,
      select: {
        id: true,
        email: true,
        name: true,
        createdAt: true,
      },
    });
  }

  async update(id: string, data: UpdateUserDTO) {
    return db.user.update({
      where: { id },
      data,
      select: {
        id: true,
        email: true,
        name: true,
        updatedAt: true,
      },
    });
  }

  async delete(id: string) {
    return db.user.delete({
      where: { id },
    });
  }

  async count() {
    return db.user.count();
  }
}
```

---

## Best Practices

### Guardrails
- ✓ Use TypeScript for type safety
- ✓ Implement proper error handling middleware
- ✓ Validate all inputs with Zod or Joi
- ✓ Use helmet for security headers
- ✓ Implement rate limiting
- ✓ Use environment variables for configuration
- ✓ Log all requests and errors
- ✓ Implement graceful shutdown
- ✓ Use async/await with proper error handling
- ✓ Follow REST conventions

### Security Checklist
- ✓ Enable CORS properly
- ✓ Use HTTPS in production
- ✓ Sanitize user inputs
- ✓ Implement authentication/authorization
- ✓ Use parameterized queries
- ✓ Set secure cookie options
- ✓ Implement rate limiting
- ✓ Keep dependencies updated

---

## Dependencies

```json
{
  "dependencies": {
    "express": "^4.18.2",
    "cors": "^2.8.5",
    "helmet": "^7.1.0",
    "morgan": "^1.10.0",
    "dotenv": "^16.3.1",
    "zod": "^3.22.4",
    "jsonwebtoken": "^9.0.2",
    "bcryptjs": "^2.4.3",
    "@prisma/client": "^5.x"
  },
  "devDependencies": {
    "typescript": "^5.x",
    "@types/express": "^4.17.21",
    "@types/node": "^20.x",
    "ts-node-dev": "^2.0.0",
    "jest": "^29.x",
    "supertest": "^6.x",
    "@types/supertest": "^6.x"
  }
}
```

---

## References

- [Express.js Documentation](https://expressjs.com/)
- [Express Best Practices](https://expressjs.com/en/advanced/best-practice-security.html)
- [Prisma with Express](https://www.prisma.io/express)
- [Zod Documentation](https://zod.dev/)
