# FastAPI Framework Guide

> **Applies to**: FastAPI 0.100+, Pydantic v2, SQLAlchemy 2.0
> **Language**: Python 3.10+
> **Type**: Async API Framework

---

## Overview

FastAPI is a modern, fast (high-performance) web framework for building APIs with Python based on standard Python type hints. It's built on Starlette for the web parts and Pydantic for data validation.

**Use FastAPI when:**
- Building REST APIs or GraphQL backends
- Need async/await support natively
- Want automatic OpenAPI documentation
- Need high performance (comparable to Node.js/Go)
- Want strong type validation with Pydantic

**Consider alternatives when:**
- Need full-stack with templates (consider Django)
- Building simple scripts (consider Flask)
- Need extensive admin interface (consider Django)

---

## Project Structure

### Standard FastAPI Project
```
myproject/
├── pyproject.toml
├── alembic.ini
├── alembic/
│   ├── versions/
│   └── env.py
├── app/
│   ├── __init__.py
│   ├── main.py                 # Application entry point
│   ├── config.py               # Settings configuration
│   ├── database.py             # Database connection
│   ├── dependencies.py         # Shared dependencies
│   ├── models/                 # SQLAlchemy models
│   │   ├── __init__.py
│   │   ├── base.py
│   │   └── user.py
│   ├── schemas/                # Pydantic schemas
│   │   ├── __init__.py
│   │   └── user.py
│   ├── api/                    # API routes
│   │   ├── __init__.py
│   │   ├── deps.py             # API dependencies
│   │   └── v1/
│   │       ├── __init__.py
│   │       ├── router.py
│   │       └── endpoints/
│   │           ├── users.py
│   │           └── products.py
│   ├── services/               # Business logic
│   │   ├── __init__.py
│   │   └── user.py
│   ├── repositories/           # Data access
│   │   ├── __init__.py
│   │   └── user.py
│   └── core/                   # Core utilities
│       ├── __init__.py
│       ├── security.py
│       └── exceptions.py
├── tests/
│   ├── conftest.py
│   ├── test_users.py
│   └── factories.py
└── docker-compose.yml
```

---

## Application Setup

### Main Application
```python
# app/main.py
from contextlib import asynccontextmanager
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.api.v1.router import api_router
from app.config import settings
from app.database import engine
from app.models.base import Base


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan events."""
    # Startup
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    yield
    # Shutdown
    await engine.dispose()


app = FastAPI(
    title=settings.PROJECT_NAME,
    version=settings.VERSION,
    openapi_url=f"{settings.API_V1_STR}/openapi.json",
    lifespan=lifespan,
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include routers
app.include_router(api_router, prefix=settings.API_V1_STR)


@app.get("/health")
async def health_check():
    return {"status": "healthy"}
```

### Configuration
```python
# app/config.py
from functools import lru_cache
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        case_sensitive=True,
    )

    # Application
    PROJECT_NAME: str = "FastAPI App"
    VERSION: str = "1.0.0"
    DEBUG: bool = False
    API_V1_STR: str = "/api/v1"

    # Database
    DATABASE_URL: str = "postgresql+asyncpg://user:pass@localhost/db"

    # Security
    SECRET_KEY: str
    ACCESS_TOKEN_EXPIRE_MINUTES: int = 30
    REFRESH_TOKEN_EXPIRE_DAYS: int = 7

    # CORS
    CORS_ORIGINS: list[str] = ["http://localhost:3000"]

    # Redis
    REDIS_URL: str = "redis://localhost:6379"


@lru_cache
def get_settings() -> Settings:
    return Settings()


settings = get_settings()
```

---

## Database Setup

### Async SQLAlchemy
```python
# app/database.py
from sqlalchemy.ext.asyncio import (
    AsyncSession,
    async_sessionmaker,
    create_async_engine,
)

from app.config import settings

engine = create_async_engine(
    settings.DATABASE_URL,
    echo=settings.DEBUG,
    pool_pre_ping=True,
    pool_size=5,
    max_overflow=10,
)

AsyncSessionLocal = async_sessionmaker(
    engine,
    class_=AsyncSession,
    expire_on_commit=False,
    autocommit=False,
    autoflush=False,
)


async def get_db() -> AsyncSession:
    """Dependency for getting async database session."""
    async with AsyncSessionLocal() as session:
        try:
            yield session
            await session.commit()
        except Exception:
            await session.rollback()
            raise
        finally:
            await session.close()
```

### Base Model
```python
# app/models/base.py
from datetime import datetime
from uuid import uuid4

from sqlalchemy import DateTime, func
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column


class Base(DeclarativeBase):
    pass


class TimestampMixin:
    created_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True),
        server_default=func.now(),
        nullable=False,
    )
    updated_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True),
        server_default=func.now(),
        onupdate=func.now(),
        nullable=False,
    )


class UUIDMixin:
    id: Mapped[str] = mapped_column(
        UUID(as_uuid=True),
        primary_key=True,
        default=uuid4,
    )
```

---

## Models

### SQLAlchemy Models
```python
# app/models/user.py
from sqlalchemy import Boolean, String, Text, ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.models.base import Base, TimestampMixin, UUIDMixin


class User(Base, UUIDMixin, TimestampMixin):
    __tablename__ = "users"

    email: Mapped[str] = mapped_column(String(255), unique=True, index=True)
    hashed_password: Mapped[str] = mapped_column(String(255))
    first_name: Mapped[str | None] = mapped_column(String(100))
    last_name: Mapped[str | None] = mapped_column(String(100))
    is_active: Mapped[bool] = mapped_column(Boolean, default=True)
    is_superuser: Mapped[bool] = mapped_column(Boolean, default=False)

    # Relationships
    profile: Mapped["Profile"] = relationship(back_populates="user", uselist=False)
    posts: Mapped[list["Post"]] = relationship(back_populates="author")

    @property
    def full_name(self) -> str:
        if self.first_name and self.last_name:
            return f"{self.first_name} {self.last_name}"
        return self.email


class Profile(Base, UUIDMixin, TimestampMixin):
    __tablename__ = "profiles"

    user_id: Mapped[str] = mapped_column(
        ForeignKey("users.id", ondelete="CASCADE"),
        unique=True,
    )
    bio: Mapped[str | None] = mapped_column(Text)
    avatar_url: Mapped[str | None] = mapped_column(String(500))

    user: Mapped["User"] = relationship(back_populates="profile")


class Post(Base, UUIDMixin, TimestampMixin):
    __tablename__ = "posts"

    title: Mapped[str] = mapped_column(String(255))
    content: Mapped[str] = mapped_column(Text)
    is_published: Mapped[bool] = mapped_column(Boolean, default=False)
    author_id: Mapped[str] = mapped_column(
        ForeignKey("users.id", ondelete="CASCADE"),
    )

    author: Mapped["User"] = relationship(back_populates="posts")
```

---

## Schemas (Pydantic v2)

```python
# app/schemas/user.py
from datetime import datetime
from uuid import UUID

from pydantic import BaseModel, ConfigDict, EmailStr, Field


# Base schemas
class UserBase(BaseModel):
    email: EmailStr
    first_name: str | None = None
    last_name: str | None = None


# Create schema
class UserCreate(UserBase):
    password: str = Field(..., min_length=8, max_length=100)


# Update schema
class UserUpdate(BaseModel):
    email: EmailStr | None = None
    first_name: str | None = None
    last_name: str | None = None
    password: str | None = Field(None, min_length=8, max_length=100)


# Response schema
class UserResponse(UserBase):
    model_config = ConfigDict(from_attributes=True)

    id: UUID
    is_active: bool
    created_at: datetime
    full_name: str


# With relationships
class UserWithProfile(UserResponse):
    profile: "ProfileResponse | None" = None


class ProfileResponse(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    id: UUID
    bio: str | None
    avatar_url: str | None


# Pagination
class PaginatedResponse[T](BaseModel):
    items: list[T]
    total: int
    page: int
    size: int
    pages: int


# Auth schemas
class Token(BaseModel):
    access_token: str
    refresh_token: str
    token_type: str = "bearer"


class TokenPayload(BaseModel):
    sub: str
    exp: datetime
    type: str
```

---

## Repositories (Data Access Layer)

```python
# app/repositories/base.py
from typing import Any, Generic, TypeVar
from uuid import UUID

from sqlalchemy import select, func
from sqlalchemy.ext.asyncio import AsyncSession

from app.models.base import Base

ModelType = TypeVar("ModelType", bound=Base)


class BaseRepository(Generic[ModelType]):
    def __init__(self, model: type[ModelType], session: AsyncSession):
        self.model = model
        self.session = session

    async def get(self, id: UUID) -> ModelType | None:
        result = await self.session.execute(
            select(self.model).where(self.model.id == id)
        )
        return result.scalar_one_or_none()

    async def get_multi(
        self,
        *,
        skip: int = 0,
        limit: int = 100,
    ) -> list[ModelType]:
        result = await self.session.execute(
            select(self.model).offset(skip).limit(limit)
        )
        return list(result.scalars().all())

    async def create(self, obj_in: dict[str, Any]) -> ModelType:
        db_obj = self.model(**obj_in)
        self.session.add(db_obj)
        await self.session.flush()
        await self.session.refresh(db_obj)
        return db_obj

    async def update(
        self,
        db_obj: ModelType,
        obj_in: dict[str, Any],
    ) -> ModelType:
        for field, value in obj_in.items():
            if value is not None:
                setattr(db_obj, field, value)
        await self.session.flush()
        await self.session.refresh(db_obj)
        return db_obj

    async def delete(self, id: UUID) -> bool:
        obj = await self.get(id)
        if obj:
            await self.session.delete(obj)
            await self.session.flush()
            return True
        return False

    async def count(self) -> int:
        result = await self.session.execute(
            select(func.count()).select_from(self.model)
        )
        return result.scalar_one()


# app/repositories/user.py
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import selectinload

from app.models.user import User
from app.repositories.base import BaseRepository


class UserRepository(BaseRepository[User]):
    def __init__(self, session: AsyncSession):
        super().__init__(User, session)

    async def get_by_email(self, email: str) -> User | None:
        result = await self.session.execute(
            select(User).where(User.email == email)
        )
        return result.scalar_one_or_none()

    async def get_with_profile(self, id: str) -> User | None:
        result = await self.session.execute(
            select(User)
            .options(selectinload(User.profile))
            .where(User.id == id)
        )
        return result.scalar_one_or_none()

    async def get_active_users(
        self,
        *,
        skip: int = 0,
        limit: int = 100,
    ) -> list[User]:
        result = await self.session.execute(
            select(User)
            .where(User.is_active == True)
            .offset(skip)
            .limit(limit)
        )
        return list(result.scalars().all())
```

---

## Services (Business Logic)

```python
# app/services/user.py
from uuid import UUID

from sqlalchemy.ext.asyncio import AsyncSession

from app.core.security import get_password_hash, verify_password
from app.core.exceptions import (
    NotFoundException,
    ConflictException,
    UnauthorizedException,
)
from app.models.user import User
from app.repositories.user import UserRepository
from app.schemas.user import UserCreate, UserUpdate


class UserService:
    def __init__(self, session: AsyncSession):
        self.session = session
        self.repository = UserRepository(session)

    async def get(self, user_id: UUID) -> User:
        user = await self.repository.get(user_id)
        if not user:
            raise NotFoundException(f"User {user_id} not found")
        return user

    async def get_by_email(self, email: str) -> User | None:
        return await self.repository.get_by_email(email)

    async def get_multi(
        self,
        *,
        skip: int = 0,
        limit: int = 100,
    ) -> tuple[list[User], int]:
        users = await self.repository.get_multi(skip=skip, limit=limit)
        total = await self.repository.count()
        return users, total

    async def create(self, user_in: UserCreate) -> User:
        # Check if email already exists
        existing = await self.repository.get_by_email(user_in.email)
        if existing:
            raise ConflictException("Email already registered")

        # Create user with hashed password
        user_data = user_in.model_dump(exclude={"password"})
        user_data["hashed_password"] = get_password_hash(user_in.password)

        return await self.repository.create(user_data)

    async def update(self, user_id: UUID, user_in: UserUpdate) -> User:
        user = await self.get(user_id)

        update_data = user_in.model_dump(exclude_unset=True, exclude={"password"})

        if user_in.password:
            update_data["hashed_password"] = get_password_hash(user_in.password)

        return await self.repository.update(user, update_data)

    async def delete(self, user_id: UUID) -> bool:
        await self.get(user_id)  # Verify exists
        return await self.repository.delete(user_id)

    async def authenticate(self, email: str, password: str) -> User:
        user = await self.repository.get_by_email(email)
        if not user:
            raise UnauthorizedException("Invalid credentials")
        if not verify_password(password, user.hashed_password):
            raise UnauthorizedException("Invalid credentials")
        if not user.is_active:
            raise UnauthorizedException("User is inactive")
        return user
```

---

## API Endpoints

### Router Setup
```python
# app/api/v1/router.py
from fastapi import APIRouter

from app.api.v1.endpoints import auth, users, products

api_router = APIRouter()
api_router.include_router(auth.router, prefix="/auth", tags=["auth"])
api_router.include_router(users.router, prefix="/users", tags=["users"])
api_router.include_router(products.router, prefix="/products", tags=["products"])
```

### User Endpoints
```python
# app/api/v1/endpoints/users.py
from uuid import UUID

from fastapi import APIRouter, Depends, Query, status
from sqlalchemy.ext.asyncio import AsyncSession

from app.api.deps import get_current_user, get_current_superuser
from app.database import get_db
from app.models.user import User
from app.schemas.user import (
    UserCreate,
    UserUpdate,
    UserResponse,
    PaginatedResponse,
)
from app.services.user import UserService

router = APIRouter()


@router.get("", response_model=PaginatedResponse[UserResponse])
async def list_users(
    page: int = Query(1, ge=1),
    size: int = Query(20, ge=1, le=100),
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user),
):
    """Get list of users with pagination."""
    service = UserService(db)
    skip = (page - 1) * size
    users, total = await service.get_multi(skip=skip, limit=size)

    return PaginatedResponse(
        items=[UserResponse.model_validate(u) for u in users],
        total=total,
        page=page,
        size=size,
        pages=(total + size - 1) // size,
    )


@router.post("", response_model=UserResponse, status_code=status.HTTP_201_CREATED)
async def create_user(
    user_in: UserCreate,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_superuser),
):
    """Create new user. Only superusers."""
    service = UserService(db)
    user = await service.create(user_in)
    return UserResponse.model_validate(user)


@router.get("/me", response_model=UserResponse)
async def get_current_user_info(
    current_user: User = Depends(get_current_user),
):
    """Get current user info."""
    return UserResponse.model_validate(current_user)


@router.get("/{user_id}", response_model=UserResponse)
async def get_user(
    user_id: UUID,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user),
):
    """Get user by ID."""
    service = UserService(db)
    user = await service.get(user_id)
    return UserResponse.model_validate(user)


@router.patch("/{user_id}", response_model=UserResponse)
async def update_user(
    user_id: UUID,
    user_in: UserUpdate,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_user),
):
    """Update user."""
    # Only allow users to update themselves or superusers
    if str(current_user.id) != str(user_id) and not current_user.is_superuser:
        from app.core.exceptions import ForbiddenException
        raise ForbiddenException("Not allowed to update this user")

    service = UserService(db)
    user = await service.update(user_id, user_in)
    return UserResponse.model_validate(user)


@router.delete("/{user_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_user(
    user_id: UUID,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_superuser),
):
    """Delete user. Only superusers."""
    service = UserService(db)
    await service.delete(user_id)
```

### Auth Endpoints
```python
# app/api/v1/endpoints/auth.py
from fastapi import APIRouter, Depends
from fastapi.security import OAuth2PasswordRequestForm
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.security import create_access_token, create_refresh_token
from app.database import get_db
from app.schemas.user import Token, UserCreate, UserResponse
from app.services.user import UserService

router = APIRouter()


@router.post("/login", response_model=Token)
async def login(
    form_data: OAuth2PasswordRequestForm = Depends(),
    db: AsyncSession = Depends(get_db),
):
    """OAuth2 compatible token login."""
    service = UserService(db)
    user = await service.authenticate(form_data.username, form_data.password)

    return Token(
        access_token=create_access_token(str(user.id)),
        refresh_token=create_refresh_token(str(user.id)),
    )


@router.post("/register", response_model=UserResponse)
async def register(
    user_in: UserCreate,
    db: AsyncSession = Depends(get_db),
):
    """Register new user."""
    service = UserService(db)
    user = await service.create(user_in)
    return UserResponse.model_validate(user)


@router.post("/refresh", response_model=Token)
async def refresh_token(
    refresh_token: str,
    db: AsyncSession = Depends(get_db),
):
    """Refresh access token."""
    from app.core.security import verify_token

    payload = verify_token(refresh_token, token_type="refresh")
    user_id = payload.sub

    return Token(
        access_token=create_access_token(user_id),
        refresh_token=create_refresh_token(user_id),
    )
```

---

## Dependencies

```python
# app/api/deps.py
from fastapi import Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.security import verify_token
from app.database import get_db
from app.models.user import User
from app.repositories.user import UserRepository

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="/api/v1/auth/login")


async def get_current_user(
    token: str = Depends(oauth2_scheme),
    db: AsyncSession = Depends(get_db),
) -> User:
    """Get current authenticated user."""
    payload = verify_token(token)

    repository = UserRepository(db)
    user = await repository.get(payload.sub)

    if not user:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User not found",
        )

    if not user.is_active:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Inactive user",
        )

    return user


async def get_current_superuser(
    current_user: User = Depends(get_current_user),
) -> User:
    """Get current superuser."""
    if not current_user.is_superuser:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not enough permissions",
        )
    return current_user


async def get_optional_user(
    token: str | None = Depends(oauth2_scheme),
    db: AsyncSession = Depends(get_db),
) -> User | None:
    """Get current user if authenticated, None otherwise."""
    if not token:
        return None

    try:
        return await get_current_user(token, db)
    except HTTPException:
        return None
```

---

## Security

```python
# app/core/security.py
from datetime import datetime, timedelta

from jose import JWTError, jwt
from passlib.context import CryptContext

from app.config import settings
from app.core.exceptions import UnauthorizedException
from app.schemas.user import TokenPayload

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

ALGORITHM = "HS256"


def verify_password(plain_password: str, hashed_password: str) -> bool:
    return pwd_context.verify(plain_password, hashed_password)


def get_password_hash(password: str) -> str:
    return pwd_context.hash(password)


def create_access_token(subject: str) -> str:
    expire = datetime.utcnow() + timedelta(
        minutes=settings.ACCESS_TOKEN_EXPIRE_MINUTES
    )
    to_encode = {
        "sub": subject,
        "exp": expire,
        "type": "access",
    }
    return jwt.encode(to_encode, settings.SECRET_KEY, algorithm=ALGORITHM)


def create_refresh_token(subject: str) -> str:
    expire = datetime.utcnow() + timedelta(days=settings.REFRESH_TOKEN_EXPIRE_DAYS)
    to_encode = {
        "sub": subject,
        "exp": expire,
        "type": "refresh",
    }
    return jwt.encode(to_encode, settings.SECRET_KEY, algorithm=ALGORITHM)


def verify_token(token: str, token_type: str = "access") -> TokenPayload:
    try:
        payload = jwt.decode(token, settings.SECRET_KEY, algorithms=[ALGORITHM])
        token_data = TokenPayload(**payload)

        if token_data.type != token_type:
            raise UnauthorizedException("Invalid token type")

        return token_data

    except JWTError:
        raise UnauthorizedException("Invalid token")
```

---

## Exception Handling

```python
# app/core/exceptions.py
from fastapi import HTTPException, status


class AppException(HTTPException):
    """Base application exception."""

    def __init__(
        self,
        detail: str,
        status_code: int = status.HTTP_500_INTERNAL_SERVER_ERROR,
    ):
        super().__init__(status_code=status_code, detail=detail)


class NotFoundException(AppException):
    def __init__(self, detail: str = "Not found"):
        super().__init__(detail=detail, status_code=status.HTTP_404_NOT_FOUND)


class ConflictException(AppException):
    def __init__(self, detail: str = "Conflict"):
        super().__init__(detail=detail, status_code=status.HTTP_409_CONFLICT)


class UnauthorizedException(AppException):
    def __init__(self, detail: str = "Unauthorized"):
        super().__init__(detail=detail, status_code=status.HTTP_401_UNAUTHORIZED)


class ForbiddenException(AppException):
    def __init__(self, detail: str = "Forbidden"):
        super().__init__(detail=detail, status_code=status.HTTP_403_FORBIDDEN)


class ValidationException(AppException):
    def __init__(self, detail: str = "Validation error"):
        super().__init__(
            detail=detail, status_code=status.HTTP_422_UNPROCESSABLE_ENTITY
        )


# Exception handlers
from fastapi import Request
from fastapi.responses import JSONResponse


async def app_exception_handler(request: Request, exc: AppException):
    return JSONResponse(
        status_code=exc.status_code,
        content={"detail": exc.detail},
    )


# Register in main.py
# app.add_exception_handler(AppException, app_exception_handler)
```

---

## Testing

```python
# tests/conftest.py
import asyncio
from typing import AsyncGenerator

import pytest
import pytest_asyncio
from httpx import AsyncClient
from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine
from sqlalchemy.orm import sessionmaker

from app.config import settings
from app.database import get_db
from app.main import app
from app.models.base import Base

# Test database
TEST_DATABASE_URL = settings.DATABASE_URL.replace("mydb", "mydb_test")

engine = create_async_engine(TEST_DATABASE_URL, echo=False)
TestSessionLocal = sessionmaker(engine, class_=AsyncSession, expire_on_commit=False)


@pytest.fixture(scope="session")
def event_loop():
    loop = asyncio.get_event_loop_policy().new_event_loop()
    yield loop
    loop.close()


@pytest_asyncio.fixture(scope="function")
async def db_session() -> AsyncGenerator[AsyncSession, None]:
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)

    async with TestSessionLocal() as session:
        yield session

    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.drop_all)


@pytest_asyncio.fixture(scope="function")
async def client(db_session: AsyncSession) -> AsyncGenerator[AsyncClient, None]:
    async def override_get_db():
        yield db_session

    app.dependency_overrides[get_db] = override_get_db

    async with AsyncClient(app=app, base_url="http://test") as ac:
        yield ac

    app.dependency_overrides.clear()


@pytest_asyncio.fixture
async def authenticated_client(
    client: AsyncClient,
    db_session: AsyncSession,
) -> AsyncClient:
    from app.services.user import UserService
    from app.schemas.user import UserCreate
    from app.core.security import create_access_token

    service = UserService(db_session)
    user = await service.create(
        UserCreate(
            email="test@example.com",
            password="testpassword123",
        )
    )
    await db_session.commit()

    token = create_access_token(str(user.id))
    client.headers["Authorization"] = f"Bearer {token}"

    return client


# tests/test_users.py
import pytest
from httpx import AsyncClient


@pytest.mark.asyncio
async def test_register_user(client: AsyncClient):
    response = await client.post(
        "/api/v1/auth/register",
        json={
            "email": "newuser@example.com",
            "password": "password123",
        },
    )

    assert response.status_code == 200
    data = response.json()
    assert data["email"] == "newuser@example.com"
    assert "id" in data


@pytest.mark.asyncio
async def test_login(client: AsyncClient, db_session):
    # Create user first
    from app.services.user import UserService
    from app.schemas.user import UserCreate

    service = UserService(db_session)
    await service.create(
        UserCreate(email="login@example.com", password="password123")
    )
    await db_session.commit()

    response = await client.post(
        "/api/v1/auth/login",
        data={
            "username": "login@example.com",
            "password": "password123",
        },
    )

    assert response.status_code == 200
    data = response.json()
    assert "access_token" in data
    assert "refresh_token" in data


@pytest.mark.asyncio
async def test_get_current_user(authenticated_client: AsyncClient):
    response = await authenticated_client.get("/api/v1/users/me")

    assert response.status_code == 200
    data = response.json()
    assert data["email"] == "test@example.com"


@pytest.mark.asyncio
async def test_unauthorized_access(client: AsyncClient):
    response = await client.get("/api/v1/users/me")
    assert response.status_code == 401
```

---

## Background Tasks

```python
# app/tasks/email.py
from fastapi import BackgroundTasks

from app.core.email import send_email


def send_welcome_email(
    background_tasks: BackgroundTasks,
    email: str,
    name: str,
):
    """Add welcome email to background tasks."""
    background_tasks.add_task(
        send_email,
        to=email,
        subject="Welcome!",
        template="welcome.html",
        context={"name": name},
    )


# Using in endpoint
@router.post("/register")
async def register(
    user_in: UserCreate,
    background_tasks: BackgroundTasks,
    db: AsyncSession = Depends(get_db),
):
    service = UserService(db)
    user = await service.create(user_in)

    # Send email in background
    send_welcome_email(
        background_tasks,
        email=user.email,
        name=user.full_name,
    )

    return UserResponse.model_validate(user)
```

---

## Dependencies (pyproject.toml)

```toml
[project]
name = "myproject"
version = "1.0.0"
requires-python = ">=3.10"
dependencies = [
    "fastapi>=0.109.0",
    "uvicorn[standard]>=0.27.0",
    "pydantic>=2.5.0",
    "pydantic-settings>=2.1.0",
    "sqlalchemy[asyncio]>=2.0.0",
    "asyncpg>=0.29.0",
    "alembic>=1.13.0",
    "python-jose[cryptography]>=3.3.0",
    "passlib[bcrypt]>=1.7.4",
    "python-multipart>=0.0.6",
    "httpx>=0.26.0",
    "redis>=5.0.0",
]

[project.optional-dependencies]
dev = [
    "pytest>=7.4.0",
    "pytest-asyncio>=0.23.0",
    "pytest-cov>=4.1.0",
    "ruff>=0.1.0",
    "mypy>=1.8.0",
    "pre-commit>=3.6.0",
]
```

---

## Commands Reference

```bash
# Development
uvicorn app.main:app --reload --port 8000

# Production
uvicorn app.main:app --host 0.0.0.0 --port 8000 --workers 4

# Database migrations
alembic init alembic
alembic revision --autogenerate -m "Initial migration"
alembic upgrade head
alembic downgrade -1

# Testing
pytest
pytest -v --cov=app --cov-report=html
pytest tests/test_users.py -k "test_register"

# Linting
ruff check app/
ruff check app/ --fix
mypy app/
```

---

## Best Practices

### Do
- ✓ Use Pydantic for all request/response validation
- ✓ Use dependency injection for services
- ✓ Keep endpoints thin, logic in services
- ✓ Use async for I/O operations
- ✓ Handle errors with custom exceptions
- ✓ Use type hints everywhere
- ✓ Write comprehensive tests

### Don't
- ❌ Put business logic in endpoints
- ❌ Use sync database calls in async context
- ❌ Expose internal models directly
- ❌ Skip input validation
- ❌ Ignore error handling
- ❌ Use global state

---

## References

- [FastAPI Documentation](https://fastapi.tiangolo.com/)
- [Pydantic v2 Documentation](https://docs.pydantic.dev/)
- [SQLAlchemy 2.0 Documentation](https://docs.sqlalchemy.org/en/20/)
- [Starlette Documentation](https://www.starlette.io/)
