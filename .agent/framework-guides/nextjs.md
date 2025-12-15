# Next.js Framework Guide

> **Framework**: Next.js 14+
> **Language**: TypeScript/JavaScript
> **Use Cases**: Full-Stack Web Apps, SSR/SSG, E-commerce, Blogs, Dashboards

---

## Overview

Next.js is a React framework that enables server-side rendering, static site generation, API routes, and full-stack development with a single codebase. Version 14+ uses the App Router as the default.

---

## Project Setup

### Create New Project
```bash
# Create new Next.js app
npx create-next-app@latest my-app --typescript --tailwind --eslint --app

# Options explained:
# --typescript: Use TypeScript
# --tailwind: Include Tailwind CSS
# --eslint: Include ESLint
# --app: Use App Router (default in 14+)

cd my-app
npm run dev
```

### Project Structure (App Router)
```
my-app/
├── app/
│   ├── (auth)/                 # Route group (no URL segment)
│   │   ├── login/
│   │   │   └── page.tsx
│   │   └── register/
│   │       └── page.tsx
│   ├── dashboard/
│   │   ├── page.tsx           # /dashboard
│   │   ├── loading.tsx        # Loading UI
│   │   ├── error.tsx          # Error UI
│   │   └── layout.tsx         # Dashboard layout
│   ├── api/
│   │   └── users/
│   │       └── route.ts       # API route
│   ├── globals.css
│   ├── layout.tsx             # Root layout
│   └── page.tsx               # Home page
├── components/
│   ├── ui/                    # Reusable UI components
│   └── features/              # Feature components
├── lib/
│   ├── db.ts                  # Database client
│   └── utils.ts               # Utility functions
├── hooks/                     # Custom hooks
├── types/                     # TypeScript types
├── public/                    # Static assets
├── next.config.js
├── tailwind.config.ts
└── package.json
```

---

## Routing

### File-Based Routing
```
app/
├── page.tsx                   # / (home)
├── about/
│   └── page.tsx              # /about
├── blog/
│   ├── page.tsx              # /blog
│   └── [slug]/
│       └── page.tsx          # /blog/my-post (dynamic)
├── shop/
│   └── [...categories]/
│       └── page.tsx          # /shop/a/b/c (catch-all)
└── (marketing)/
    ├── pricing/
    │   └── page.tsx          # /pricing (grouped)
    └── features/
        └── page.tsx          # /features (grouped)
```

### Page Component
```tsx
// app/blog/[slug]/page.tsx
interface PageProps {
  params: { slug: string };
  searchParams: { [key: string]: string | string[] | undefined };
}

export default function BlogPost({ params, searchParams }: PageProps) {
  return (
    <article>
      <h1>Post: {params.slug}</h1>
    </article>
  );
}

// Generate static params for SSG
export async function generateStaticParams() {
  const posts = await getPosts();
  return posts.map((post) => ({
    slug: post.slug,
  }));
}

// Metadata
export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const post = await getPost(params.slug);
  return {
    title: post.title,
    description: post.excerpt,
    openGraph: {
      title: post.title,
      images: [post.coverImage],
    },
  };
}
```

### Layouts
```tsx
// app/layout.tsx (Root Layout)
import { Inter } from 'next/font/google';
import './globals.css';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: {
    default: 'My App',
    template: '%s | My App',
  },
  description: 'My awesome application',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <Header />
        <main>{children}</main>
        <Footer />
      </body>
    </html>
  );
}

// app/dashboard/layout.tsx (Nested Layout)
export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex">
      <Sidebar />
      <div className="flex-1">{children}</div>
    </div>
  );
}
```

---

## Server Components vs Client Components

### Server Components (Default)
```tsx
// app/users/page.tsx
// Server Component - runs on server only
import { db } from '@/lib/db';

export default async function UsersPage() {
  // Direct database access (no API needed)
  const users = await db.user.findMany();

  return (
    <div>
      <h1>Users</h1>
      <ul>
        {users.map((user) => (
          <li key={user.id}>{user.name}</li>
        ))}
      </ul>
    </div>
  );
}
```

### Client Components
```tsx
// components/Counter.tsx
'use client';

import { useState } from 'react';

export function Counter() {
  const [count, setCount] = useState(0);

  return (
    <button onClick={() => setCount(count + 1)}>
      Count: {count}
    </button>
  );
}
```

### Composition Pattern
```tsx
// app/dashboard/page.tsx (Server Component)
import { ClientSidebar } from '@/components/ClientSidebar';
import { db } from '@/lib/db';

export default async function Dashboard() {
  // Fetch data on server
  const stats = await db.stats.get();

  return (
    <div>
      {/* Pass server data to client component */}
      <ClientSidebar initialStats={stats} />

      {/* Server-rendered content */}
      <DashboardContent stats={stats} />
    </div>
  );
}
```

---

## Data Fetching

### Server Component Fetching
```tsx
// app/products/page.tsx
async function getProducts() {
  const res = await fetch('https://api.example.com/products', {
    next: { revalidate: 3600 }, // Revalidate every hour
  });

  if (!res.ok) {
    throw new Error('Failed to fetch products');
  }

  return res.json();
}

export default async function ProductsPage() {
  const products = await getProducts();

  return (
    <div>
      {products.map((product) => (
        <ProductCard key={product.id} product={product} />
      ))}
    </div>
  );
}
```

### Parallel Data Fetching
```tsx
// app/dashboard/page.tsx
async function getUser(id: string) {
  const res = await fetch(`/api/users/${id}`);
  return res.json();
}

async function getOrders(userId: string) {
  const res = await fetch(`/api/users/${userId}/orders`);
  return res.json();
}

export default async function Dashboard({ params }: { params: { id: string } }) {
  // Fetch in parallel
  const [user, orders] = await Promise.all([
    getUser(params.id),
    getOrders(params.id),
  ]);

  return (
    <div>
      <UserProfile user={user} />
      <OrderList orders={orders} />
    </div>
  );
}
```

### Streaming with Suspense
```tsx
// app/dashboard/page.tsx
import { Suspense } from 'react';

export default function Dashboard() {
  return (
    <div>
      <h1>Dashboard</h1>

      {/* Immediate render */}
      <WelcomeMessage />

      {/* Stream in when ready */}
      <Suspense fallback={<StatsSkeleton />}>
        <Stats />
      </Suspense>

      <Suspense fallback={<RecentOrdersSkeleton />}>
        <RecentOrders />
      </Suspense>
    </div>
  );
}

async function Stats() {
  const stats = await fetchStats(); // Slow query
  return <StatsDisplay stats={stats} />;
}
```

---

## API Routes

### Route Handlers
```tsx
// app/api/users/route.ts
import { NextRequest, NextResponse } from 'next/server';
import { db } from '@/lib/db';
import { z } from 'zod';

const userSchema = z.object({
  name: z.string().min(2),
  email: z.string().email(),
});

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const page = parseInt(searchParams.get('page') || '1');
  const limit = parseInt(searchParams.get('limit') || '10');

  const users = await db.user.findMany({
    skip: (page - 1) * limit,
    take: limit,
  });

  return NextResponse.json(users);
}

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    const validated = userSchema.parse(body);

    const user = await db.user.create({
      data: validated,
    });

    return NextResponse.json(user, { status: 201 });
  } catch (error) {
    if (error instanceof z.ZodError) {
      return NextResponse.json(
        { errors: error.errors },
        { status: 400 }
      );
    }
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    );
  }
}

// app/api/users/[id]/route.ts
export async function GET(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  const user = await db.user.findUnique({
    where: { id: params.id },
  });

  if (!user) {
    return NextResponse.json(
      { error: 'User not found' },
      { status: 404 }
    );
  }

  return NextResponse.json(user);
}

export async function DELETE(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  await db.user.delete({
    where: { id: params.id },
  });

  return new NextResponse(null, { status: 204 });
}
```

---

## Server Actions

### Form Actions
```tsx
// app/actions.ts
'use server';

import { revalidatePath } from 'next/cache';
import { redirect } from 'next/navigation';
import { z } from 'zod';
import { db } from '@/lib/db';

const createPostSchema = z.object({
  title: z.string().min(1),
  content: z.string().min(10),
});

export async function createPost(formData: FormData) {
  const validated = createPostSchema.parse({
    title: formData.get('title'),
    content: formData.get('content'),
  });

  const post = await db.post.create({
    data: validated,
  });

  revalidatePath('/posts');
  redirect(`/posts/${post.id}`);
}

export async function deletePost(id: string) {
  await db.post.delete({
    where: { id },
  });

  revalidatePath('/posts');
}

// app/posts/new/page.tsx
import { createPost } from '../actions';

export default function NewPost() {
  return (
    <form action={createPost}>
      <input name="title" placeholder="Title" required />
      <textarea name="content" placeholder="Content" required />
      <button type="submit">Create Post</button>
    </form>
  );
}
```

### With useFormState
```tsx
// app/actions.ts
'use server';

type ActionState = {
  errors?: {
    title?: string[];
    content?: string[];
  };
  message?: string;
};

export async function createPost(
  prevState: ActionState,
  formData: FormData
): Promise<ActionState> {
  const validatedFields = createPostSchema.safeParse({
    title: formData.get('title'),
    content: formData.get('content'),
  });

  if (!validatedFields.success) {
    return {
      errors: validatedFields.error.flatten().fieldErrors,
    };
  }

  try {
    await db.post.create({ data: validatedFields.data });
    revalidatePath('/posts');
    return { message: 'Post created successfully' };
  } catch {
    return { message: 'Failed to create post' };
  }
}

// components/PostForm.tsx
'use client';

import { useFormState, useFormStatus } from 'react-dom';
import { createPost } from '@/app/actions';

function SubmitButton() {
  const { pending } = useFormStatus();
  return (
    <button type="submit" disabled={pending}>
      {pending ? 'Creating...' : 'Create Post'}
    </button>
  );
}

export function PostForm() {
  const [state, formAction] = useFormState(createPost, {});

  return (
    <form action={formAction}>
      <div>
        <input name="title" placeholder="Title" />
        {state.errors?.title && <p className="error">{state.errors.title}</p>}
      </div>
      <div>
        <textarea name="content" placeholder="Content" />
        {state.errors?.content && <p className="error">{state.errors.content}</p>}
      </div>
      <SubmitButton />
      {state.message && <p>{state.message}</p>}
    </form>
  );
}
```

---

## Authentication

### NextAuth.js Setup
```tsx
// app/api/auth/[...nextauth]/route.ts
import NextAuth from 'next-auth';
import { authOptions } from '@/lib/auth';

const handler = NextAuth(authOptions);
export { handler as GET, handler as POST };

// lib/auth.ts
import { NextAuthOptions } from 'next-auth';
import { PrismaAdapter } from '@auth/prisma-adapter';
import GitHubProvider from 'next-auth/providers/github';
import CredentialsProvider from 'next-auth/providers/credentials';
import { db } from './db';

export const authOptions: NextAuthOptions = {
  adapter: PrismaAdapter(db),
  providers: [
    GitHubProvider({
      clientId: process.env.GITHUB_ID!,
      clientSecret: process.env.GITHUB_SECRET!,
    }),
    CredentialsProvider({
      name: 'credentials',
      credentials: {
        email: { label: 'Email', type: 'email' },
        password: { label: 'Password', type: 'password' },
      },
      async authorize(credentials) {
        // Validate credentials
        const user = await validateUser(credentials);
        return user;
      },
    }),
  ],
  session: {
    strategy: 'jwt',
  },
  pages: {
    signIn: '/login',
  },
};

// lib/auth-utils.ts
import { getServerSession } from 'next-auth';
import { authOptions } from './auth';

export async function getSession() {
  return await getServerSession(authOptions);
}

export async function getCurrentUser() {
  const session = await getSession();
  return session?.user;
}

// Protect server component
export default async function ProtectedPage() {
  const user = await getCurrentUser();

  if (!user) {
    redirect('/login');
  }

  return <div>Welcome, {user.name}</div>;
}
```

---

## Middleware

```tsx
// middleware.ts
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import { getToken } from 'next-auth/jwt';

export async function middleware(request: NextRequest) {
  const token = await getToken({ req: request });
  const isAuthPage = request.nextUrl.pathname.startsWith('/login');
  const isProtectedRoute = request.nextUrl.pathname.startsWith('/dashboard');

  if (isAuthPage) {
    if (token) {
      return NextResponse.redirect(new URL('/dashboard', request.url));
    }
    return NextResponse.next();
  }

  if (isProtectedRoute && !token) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/dashboard/:path*', '/login'],
};
```

---

## Error Handling

### Error Boundaries
```tsx
// app/dashboard/error.tsx
'use client';

import { useEffect } from 'react';

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <div className="error-container">
      <h2>Something went wrong!</h2>
      <button onClick={() => reset()}>Try again</button>
    </div>
  );
}

// app/global-error.tsx (for root layout errors)
'use client';

export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <html>
      <body>
        <h2>Something went wrong!</h2>
        <button onClick={() => reset()}>Try again</button>
      </body>
    </html>
  );
}
```

### Not Found
```tsx
// app/not-found.tsx
import Link from 'next/link';

export default function NotFound() {
  return (
    <div>
      <h2>Not Found</h2>
      <p>Could not find requested resource</p>
      <Link href="/">Return Home</Link>
    </div>
  );
}

// Trigger programmatically
import { notFound } from 'next/navigation';

export default async function Page({ params }: { params: { id: string } }) {
  const item = await getItem(params.id);

  if (!item) {
    notFound();
  }

  return <div>{item.name}</div>;
}
```

---

## Best Practices

### Guardrails
- ✓ Use Server Components by default
- ✓ Add 'use client' only when needed (interactivity, hooks)
- ✓ Colocate data fetching with components
- ✓ Use parallel data fetching with Promise.all
- ✓ Implement loading.tsx and error.tsx for each route
- ✓ Use Server Actions for mutations
- ✓ Cache and revalidate appropriately
- ✓ Use proper metadata for SEO
- ✓ Implement proper error boundaries
- ✓ Use TypeScript for type safety

### Performance Tips
- Use `next/image` for optimized images
- Use `next/font` for optimized fonts
- Implement proper caching strategies
- Use streaming with Suspense
- Minimize client-side JavaScript

---

## Configuration

### next.config.js
```js
/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '**.example.com',
      },
    ],
  },
  experimental: {
    serverActions: true,
  },
  async redirects() {
    return [
      {
        source: '/old-path',
        destination: '/new-path',
        permanent: true,
      },
    ];
  },
  async headers() {
    return [
      {
        source: '/api/:path*',
        headers: [
          { key: 'Access-Control-Allow-Origin', value: '*' },
        ],
      },
    ];
  },
};

module.exports = nextConfig;
```

---

## References

- [Next.js Documentation](https://nextjs.org/docs)
- [Next.js App Router](https://nextjs.org/docs/app)
- [NextAuth.js](https://next-auth.js.org/)
- [Vercel Deployment](https://vercel.com/docs)
