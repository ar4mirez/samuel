# React Framework Guide

> **Framework**: React 18+
> **Language**: TypeScript/JavaScript
> **Use Cases**: Single Page Applications, Component Libraries, Web Apps

---

## Overview

React is a declarative, component-based JavaScript library for building user interfaces. It uses a virtual DOM for efficient updates and supports server-side rendering.

---

## Project Setup

### Create New Project
```bash
# Vite (recommended)
npm create vite@latest my-app -- --template react-ts
cd my-app
npm install

# Create React App (legacy)
npx create-react-app my-app --template typescript

# Next.js (for SSR/SSG)
npx create-next-app@latest my-app --typescript
```

### Project Structure
```
my-app/
├── public/
│   └── index.html
├── src/
│   ├── assets/              # Static assets
│   ├── components/          # Reusable components
│   │   ├── ui/             # Base UI components
│   │   └── features/       # Feature-specific components
│   ├── hooks/              # Custom hooks
│   ├── contexts/           # React contexts
│   ├── services/           # API services
│   ├── utils/              # Utility functions
│   ├── types/              # TypeScript types
│   ├── pages/              # Page components
│   ├── App.tsx
│   ├── main.tsx
│   └── index.css
├── tests/
├── package.json
├── tsconfig.json
├── vite.config.ts
└── README.md
```

---

## Component Patterns

### Functional Component
```tsx
import { useState, useEffect } from 'react';

interface UserCardProps {
  userId: string;
  onSelect?: (user: User) => void;
  className?: string;
}

export function UserCard({ userId, onSelect, className }: UserCardProps) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const data = await userService.getById(userId);
        setUser(data);
      } finally {
        setLoading(false);
      }
    };
    fetchUser();
  }, [userId]);

  if (loading) return <Skeleton />;
  if (!user) return null;

  return (
    <div className={cn('user-card', className)} onClick={() => onSelect?.(user)}>
      <Avatar src={user.avatar} alt={user.name} />
      <h3>{user.name}</h3>
      <p>{user.email}</p>
    </div>
  );
}
```

### Component with Children
```tsx
interface CardProps {
  title: string;
  children: React.ReactNode;
  footer?: React.ReactNode;
}

export function Card({ title, children, footer }: CardProps) {
  return (
    <div className="card">
      <div className="card-header">
        <h2>{title}</h2>
      </div>
      <div className="card-body">{children}</div>
      {footer && <div className="card-footer">{footer}</div>}
    </div>
  );
}
```

### Compound Components
```tsx
interface TabsContextValue {
  activeTab: string;
  setActiveTab: (tab: string) => void;
}

const TabsContext = createContext<TabsContextValue | null>(null);

function Tabs({ children, defaultTab }: { children: React.ReactNode; defaultTab: string }) {
  const [activeTab, setActiveTab] = useState(defaultTab);

  return (
    <TabsContext.Provider value={{ activeTab, setActiveTab }}>
      <div className="tabs">{children}</div>
    </TabsContext.Provider>
  );
}

function TabList({ children }: { children: React.ReactNode }) {
  return <div className="tab-list" role="tablist">{children}</div>;
}

function Tab({ value, children }: { value: string; children: React.ReactNode }) {
  const context = useContext(TabsContext);
  if (!context) throw new Error('Tab must be used within Tabs');

  return (
    <button
      role="tab"
      aria-selected={context.activeTab === value}
      onClick={() => context.setActiveTab(value)}
    >
      {children}
    </button>
  );
}

function TabPanel({ value, children }: { value: string; children: React.ReactNode }) {
  const context = useContext(TabsContext);
  if (!context) throw new Error('TabPanel must be used within Tabs');

  if (context.activeTab !== value) return null;
  return <div role="tabpanel">{children}</div>;
}

Tabs.List = TabList;
Tabs.Tab = Tab;
Tabs.Panel = TabPanel;

export { Tabs };

// Usage
<Tabs defaultTab="profile">
  <Tabs.List>
    <Tabs.Tab value="profile">Profile</Tabs.Tab>
    <Tabs.Tab value="settings">Settings</Tabs.Tab>
  </Tabs.List>
  <Tabs.Panel value="profile">Profile content</Tabs.Panel>
  <Tabs.Panel value="settings">Settings content</Tabs.Panel>
</Tabs>
```

---

## Hooks

### Custom Data Fetching Hook
```tsx
interface UseQueryResult<T> {
  data: T | null;
  loading: boolean;
  error: Error | null;
  refetch: () => Promise<void>;
}

function useQuery<T>(fetcher: () => Promise<T>, deps: unknown[] = []): UseQueryResult<T> {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const fetch = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await fetcher();
      setData(result);
    } catch (e) {
      setError(e instanceof Error ? e : new Error('Unknown error'));
    } finally {
      setLoading(false);
    }
  }, deps);

  useEffect(() => {
    fetch();
  }, [fetch]);

  return { data, loading, error, refetch: fetch };
}

// Usage
const { data: users, loading, error, refetch } = useQuery(
  () => userService.getAll(),
  []
);
```

### useLocalStorage Hook
```tsx
function useLocalStorage<T>(key: string, initialValue: T) {
  const [storedValue, setStoredValue] = useState<T>(() => {
    try {
      const item = window.localStorage.getItem(key);
      return item ? JSON.parse(item) : initialValue;
    } catch {
      return initialValue;
    }
  });

  const setValue = useCallback((value: T | ((val: T) => T)) => {
    try {
      const valueToStore = value instanceof Function ? value(storedValue) : value;
      setStoredValue(valueToStore);
      window.localStorage.setItem(key, JSON.stringify(valueToStore));
    } catch (error) {
      console.error('Error saving to localStorage:', error);
    }
  }, [key, storedValue]);

  return [storedValue, setValue] as const;
}
```

### useDebounce Hook
```tsx
function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
}

// Usage
const [searchTerm, setSearchTerm] = useState('');
const debouncedSearch = useDebounce(searchTerm, 300);

useEffect(() => {
  if (debouncedSearch) {
    searchUsers(debouncedSearch);
  }
}, [debouncedSearch]);
```

---

## State Management

### Context API Pattern
```tsx
// contexts/AuthContext.tsx
interface AuthContextValue {
  user: User | null;
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);

  const login = async (credentials: LoginCredentials) => {
    const user = await authService.login(credentials);
    setUser(user);
  };

  const logout = () => {
    authService.logout();
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, logout, isAuthenticated: !!user }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
}
```

### Zustand (Recommended for Complex State)
```tsx
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface CartStore {
  items: CartItem[];
  addItem: (item: CartItem) => void;
  removeItem: (id: string) => void;
  clearCart: () => void;
  total: () => number;
}

const useCartStore = create<CartStore>()(
  persist(
    (set, get) => ({
      items: [],
      addItem: (item) => set((state) => ({
        items: [...state.items, item]
      })),
      removeItem: (id) => set((state) => ({
        items: state.items.filter((item) => item.id !== id)
      })),
      clearCart: () => set({ items: [] }),
      total: () => get().items.reduce((sum, item) => sum + item.price * item.quantity, 0),
    }),
    { name: 'cart-storage' }
  )
);

// Usage
function Cart() {
  const { items, removeItem, total } = useCartStore();

  return (
    <div>
      {items.map((item) => (
        <CartItem key={item.id} item={item} onRemove={() => removeItem(item.id)} />
      ))}
      <p>Total: ${total()}</p>
    </div>
  );
}
```

---

## Data Fetching

### React Query / TanStack Query
```tsx
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

// Queries
function useUsers() {
  return useQuery({
    queryKey: ['users'],
    queryFn: () => userService.getAll(),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

function useUser(id: string) {
  return useQuery({
    queryKey: ['users', id],
    queryFn: () => userService.getById(id),
    enabled: !!id,
  });
}

// Mutations
function useCreateUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateUserDTO) => userService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });
}

// Usage in component
function UserList() {
  const { data: users, isLoading, error } = useUsers();
  const createUser = useCreateUser();

  if (isLoading) return <Spinner />;
  if (error) return <ErrorMessage error={error} />;

  return (
    <div>
      {users?.map((user) => <UserCard key={user.id} user={user} />)}
      <button onClick={() => createUser.mutate({ name: 'New User' })}>
        Add User
      </button>
    </div>
  );
}
```

---

## Forms

### React Hook Form
```tsx
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

const userSchema = z.object({
  email: z.string().email('Invalid email'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
  name: z.string().min(2, 'Name is required'),
});

type UserFormData = z.infer<typeof userSchema>;

function UserForm({ onSubmit }: { onSubmit: (data: UserFormData) => void }) {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<UserFormData>({
    resolver: zodResolver(userSchema),
  });

  const handleFormSubmit = async (data: UserFormData) => {
    await onSubmit(data);
    reset();
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)}>
      <div>
        <label htmlFor="email">Email</label>
        <input id="email" type="email" {...register('email')} />
        {errors.email && <span className="error">{errors.email.message}</span>}
      </div>

      <div>
        <label htmlFor="password">Password</label>
        <input id="password" type="password" {...register('password')} />
        {errors.password && <span className="error">{errors.password.message}</span>}
      </div>

      <div>
        <label htmlFor="name">Name</label>
        <input id="name" {...register('name')} />
        {errors.name && <span className="error">{errors.name.message}</span>}
      </div>

      <button type="submit" disabled={isSubmitting}>
        {isSubmitting ? 'Submitting...' : 'Submit'}
      </button>
    </form>
  );
}
```

---

## Testing

### Component Testing with React Testing Library
```tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { UserCard } from './UserCard';

describe('UserCard', () => {
  const mockUser = {
    id: '1',
    name: 'John Doe',
    email: 'john@example.com',
  };

  it('renders user information', () => {
    render(<UserCard user={mockUser} />);

    expect(screen.getByText('John Doe')).toBeInTheDocument();
    expect(screen.getByText('john@example.com')).toBeInTheDocument();
  });

  it('calls onSelect when clicked', async () => {
    const onSelect = vi.fn();
    render(<UserCard user={mockUser} onSelect={onSelect} />);

    await userEvent.click(screen.getByRole('button'));

    expect(onSelect).toHaveBeenCalledWith(mockUser);
  });

  it('shows loading state', () => {
    render(<UserCard userId="1" />);

    expect(screen.getByTestId('skeleton')).toBeInTheDocument();
  });

  it('fetches and displays user data', async () => {
    render(<UserCard userId="1" />);

    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
    });
  });
});
```

### Hook Testing
```tsx
import { renderHook, act } from '@testing-library/react';
import { useCounter } from './useCounter';

describe('useCounter', () => {
  it('initializes with default value', () => {
    const { result } = renderHook(() => useCounter());

    expect(result.current.count).toBe(0);
  });

  it('increments counter', () => {
    const { result } = renderHook(() => useCounter());

    act(() => {
      result.current.increment();
    });

    expect(result.current.count).toBe(1);
  });

  it('decrements counter', () => {
    const { result } = renderHook(() => useCounter(10));

    act(() => {
      result.current.decrement();
    });

    expect(result.current.count).toBe(9);
  });
});
```

---

## Performance Optimization

### Memoization
```tsx
import { memo, useMemo, useCallback } from 'react';

// Memoize component
const ExpensiveList = memo(function ExpensiveList({ items, onItemClick }: Props) {
  return (
    <ul>
      {items.map((item) => (
        <li key={item.id} onClick={() => onItemClick(item)}>
          {item.name}
        </li>
      ))}
    </ul>
  );
});

// Memoize computed values
function Dashboard({ data }: { data: DataPoint[] }) {
  const statistics = useMemo(() => {
    return computeExpensiveStatistics(data);
  }, [data]);

  return <StatsDisplay stats={statistics} />;
}

// Memoize callbacks
function ParentComponent() {
  const [items, setItems] = useState<Item[]>([]);

  const handleItemClick = useCallback((item: Item) => {
    console.log('Clicked:', item);
  }, []);

  return <ExpensiveList items={items} onItemClick={handleItemClick} />;
}
```

### Code Splitting
```tsx
import { lazy, Suspense } from 'react';

// Lazy load components
const Dashboard = lazy(() => import('./pages/Dashboard'));
const Settings = lazy(() => import('./pages/Settings'));

function App() {
  return (
    <Suspense fallback={<LoadingSpinner />}>
      <Routes>
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/settings" element={<Settings />} />
      </Routes>
    </Suspense>
  );
}
```

---

## Best Practices

### Guardrails
- ✓ Use TypeScript for all components
- ✓ Keep components small and focused (<200 lines)
- ✓ Extract reusable logic into custom hooks
- ✓ Use proper prop types and interfaces
- ✓ Implement error boundaries for graceful error handling
- ✓ Use React.memo() for expensive pure components
- ✓ Prefer composition over prop drilling
- ✓ Use proper accessibility attributes (ARIA)
- ✓ Test components with React Testing Library
- ✓ Use ESLint with react-hooks plugin

### File Naming
- Components: `PascalCase.tsx` (e.g., `UserCard.tsx`)
- Hooks: `camelCase.ts` with `use` prefix (e.g., `useAuth.ts`)
- Utils: `camelCase.ts` (e.g., `formatDate.ts`)
- Tests: `*.test.tsx` or `*.spec.tsx`

---

## Dependencies

### Essential
```json
{
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.x",
    "@tanstack/react-query": "^5.x",
    "zustand": "^4.x"
  },
  "devDependencies": {
    "@types/react": "^18.x",
    "@types/react-dom": "^18.x",
    "@testing-library/react": "^14.x",
    "vitest": "^1.x"
  }
}
```

---

## References

- [React Documentation](https://react.dev/)
- [React TypeScript Cheatsheet](https://react-typescript-cheatsheet.netlify.app/)
- [TanStack Query](https://tanstack.com/query/latest)
- [Zustand](https://zustand-demo.pmnd.rs/)
- [React Testing Library](https://testing-library.com/docs/react-testing-library/intro/)
