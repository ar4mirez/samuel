---
title: TypeScript Guide
description: TypeScript and JavaScript development guardrails and best practices
---

# TypeScript Guide

> **Applies to**: TypeScript, JavaScript (ES6+), React, Node.js, Next.js, Vue, Angular

---

## Core Principles

1. **Type Safety First**: Use TypeScript strict mode, avoid `any`
2. **Immutability**: Prefer `const`, use readonly, avoid mutations
3. **Async/Await**: Modern async patterns over callbacks/raw promises
4. **Functional Patterns**: Pure functions, map/filter/reduce over loops
5. **Composition**: Small, composable functions over large classes

---

## Language-Specific Guardrails

### TypeScript Configuration

```
✓ Use strict mode: "strict": true in tsconfig.json
✓ Enable noUncheckedIndexedAccess, noImplicitReturns
✓ No any types without explicit // @ts-expect-error comment with justification
✓ All function parameters and return types explicitly typed
✓ Use unknown instead of any when type truly unknown
```

### Code Style

```
✓ Prefer const over let, never use var
✓ Use async/await over raw promises (better error handling)
✓ Arrow functions for callbacks: array.map(x => x * 2)
✓ Template literals over string concatenation: `Hello ${name}`
✓ Destructuring for object/array access: const { id, name } = user
✓ Optional chaining: user?.address?.city instead of nested checks
✓ Nullish coalescing: value ?? defaultValue instead of ||
```

### React-Specific

```
✓ Functional components over class components
✓ Hooks over HOCs or render props
✓ Component files ≤200 lines (split into smaller components)
✓ Props validated with TypeScript interfaces or Zod
✓ No inline functions in JSX (causes unnecessary re-renders)
✓ Use React.memo() for expensive components
✓ Custom hooks for reusable logic (prefix with use)
```

### Node.js-Specific

```
✓ Use ES modules (import/export) over CommonJS (require)
✓ Environment variables via process.env with validation (Zod)
✓ Async error handling with try/catch or error middleware
✓ Graceful shutdown (handle SIGTERM, SIGINT)
✓ Rate limiting on public API endpoints
```

---

## Validation & Input Handling

### Recommended Libraries

| Library | Description |
|---------|-------------|
| **Zod** | Runtime type validation + TypeScript inference (recommended) |
| **Yup** | Schema validation (older but still popular) |
| **io-ts** | Functional runtime type checking |

### Pattern

```typescript
import { z } from 'zod';

const UserSchema = z.object({
  email: z.string().email(),
  age: z.number().int().positive(),
  role: z.enum(['admin', 'user']),
});

type User = z.infer<typeof UserSchema>;

// Validate API input
app.post('/users', (req, res) => {
  const result = UserSchema.safeParse(req.body);
  if (!result.success) {
    return res.status(400).json({ errors: result.error });
  }
  const user: User = result.data; // Type-safe!
});
```

---

## Testing

### Frameworks

| Framework | Use Case |
|-----------|----------|
| **Vitest** | Fast, modern (recommended for new projects) |
| **Jest** | Industry standard (mature ecosystem) |
| **Playwright/Cypress** | E2E testing |

### Guardrails

```
✓ Test files alongside code: user.service.ts + user.service.test.ts
✓ Test naming: describe('UserService') → it('should throw error when email invalid')
✓ Use beforeEach for setup, avoid test interdependencies
✓ Mock external dependencies (APIs, database, file system)
✓ Test both success and error paths
✓ Coverage target: >80% for business logic
```

### Example

```typescript
import { describe, it, expect, beforeEach } from 'vitest';
import { UserService } from './user.service';

describe('UserService', () => {
  let userService: UserService;

  beforeEach(() => {
    userService = new UserService();
  });

  it('should create user with valid data', async () => {
    const user = await userService.create({
      email: 'test@example.com',
      age: 25,
    });
    expect(user.id).toBeDefined();
    expect(user.email).toBe('test@example.com');
  });

  it('should throw error when email invalid', async () => {
    await expect(
      userService.create({ email: 'invalid', age: 25 })
    ).rejects.toThrow('Invalid email');
  });
});
```

---

## Tooling

### Essential Tools

| Tool | Purpose |
|------|---------|
| **ESLint** | Linting (detect bugs, enforce style) |
| **Prettier** | Formatting (consistent code style) |
| **TypeScript** | Type checking |
| **ts-node/tsx** | Run TypeScript directly |

### Configuration Files

=== "tsconfig.json"

    ```json
    {
      "compilerOptions": {
        "strict": true,
        "target": "ES2022",
        "module": "ESNext",
        "moduleResolution": "bundler",
        "esModuleInterop": true,
        "skipLibCheck": true,
        "noUncheckedIndexedAccess": true,
        "noImplicitReturns": true
      }
    }
    ```

=== ".eslintrc.json"

    ```json
    {
      "extends": [
        "eslint:recommended",
        "plugin:@typescript-eslint/recommended",
        "prettier"
      ],
      "rules": {
        "@typescript-eslint/no-explicit-any": "error",
        "@typescript-eslint/explicit-function-return-type": "warn"
      }
    }
    ```

### Pre-Commit Commands

```bash
# Type check
tsc --noEmit

# Lint
eslint . --ext .ts,.tsx

# Format
prettier --write .

# Test
npm test
```

---

## Common Pitfalls

### Don't Do This

```typescript
// ❌ Using any
function process(data: any) { ... }

// ❌ Mutation
const user = { name: 'John' };
user.name = 'Jane'; // Mutating object

// ❌ Ignoring errors
try {
  await riskyOperation();
} catch (e) {
  // Empty catch - error swallowed
}

// ❌ Inline functions in React
<button onClick={() => handleClick()}>Click</button>

// ❌ No return type
function calculate(a, b) {
  return a + b;
}
```

### Do This Instead

```typescript
// ✅ Proper typing
function process(data: UserData): ProcessedData { ... }

// ✅ Immutability
const user = { name: 'John' } as const;
const updatedUser = { ...user, name: 'Jane' };

// ✅ Proper error handling
try {
  await riskyOperation();
} catch (error) {
  logger.error('Operation failed:', error);
  throw new AppError('Failed to process', { cause: error });
}

// ✅ Memoized callback in React
const handleClick = useCallback(() => {
  // handler logic
}, [dependencies]);
<button onClick={handleClick}>Click</button>

// ✅ Explicit return type
function calculate(a: number, b: number): number {
  return a + b;
}
```

---

## Framework-Specific Patterns

### React + TypeScript

```typescript
// Component with props
interface UserCardProps {
  user: User;
  onEdit?: (user: User) => void;
}

export function UserCard({ user, onEdit }: UserCardProps) {
  return (
    <div>
      <h3>{user.name}</h3>
      {onEdit && <button onClick={() => onEdit(user)}>Edit</button>}
    </div>
  );
}

// Custom hook
function useUser(userId: string) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchUser(userId).then(setUser).finally(() => setLoading(false));
  }, [userId]);

  return { user, loading };
}
```

### Express + TypeScript

```typescript
import express, { Request, Response, NextFunction } from 'express';

// Typed request handler
app.get('/users/:id', async (req: Request<{ id: string }>, res: Response) => {
  const user = await userService.findById(req.params.id);
  if (!user) {
    return res.status(404).json({ error: 'User not found' });
  }
  res.json(user);
});

// Error middleware
app.use((err: Error, req: Request, res: Response, next: NextFunction) => {
  logger.error(err);
  res.status(500).json({ error: 'Internal server error' });
});
```

---

## Performance Considerations

### Optimization Guardrails

```
✓ Bundle size < 200KB initial load (use code splitting)
✓ Lazy load routes: const Page = lazy(() => import('./Page'))
✓ Memoize expensive computations: useMemo, React.memo
✓ Debounce user input handlers (search, autocomplete)
✓ Use pagination for large datasets (not load all)
✓ Avoid unnecessary re-renders (React DevTools Profiler)
```

### Example

```typescript
// Code splitting
const Dashboard = lazy(() => import('./pages/Dashboard'));

// Memoization
const expensiveValue = useMemo(() => {
  return heavyComputation(data);
}, [data]);

// Debouncing
const debouncedSearch = useDeferredValue(searchTerm);
```

---

## References

- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html)
- [React TypeScript Cheatsheet](https://react-typescript-cheatsheet.netlify.app/)
- [Effective TypeScript (Book)](https://effectivetypescript.com/)
- [TypeScript ESLint](https://typescript-eslint.io/)
- [Zod Documentation](https://zod.dev/)
