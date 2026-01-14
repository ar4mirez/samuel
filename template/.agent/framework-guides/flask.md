# Flask Framework Guide

> **Applies to**: Flask 3.0+, REST APIs, Microservices, Web Applications
> **Language Guide**: @.agent/language-guides/python.md

---

## Overview

Flask is a lightweight WSGI web framework that provides the basics for building web applications while allowing flexibility in choosing components. It's ideal for microservices, APIs, and applications where you want full control over the architecture.

**Use Flask when:**
- Building microservices or small-to-medium APIs
- You want flexibility to choose your own ORM, auth, etc.
- Rapid prototyping is needed
- You prefer explicit over implicit behavior

**Consider alternatives when:**
- You need async support (use FastAPI or Quart)
- You want batteries-included (use Django)
- High performance async is critical (use FastAPI)

---

## Project Structure

```
myproject/
├── app/
│   ├── __init__.py           # Application factory
│   ├── config.py             # Configuration classes
│   ├── extensions.py         # Flask extensions
│   ├── models/
│   │   ├── __init__.py
│   │   ├── base.py           # Base model class
│   │   └── user.py
│   ├── api/
│   │   ├── __init__.py       # API blueprint
│   │   ├── users.py          # User endpoints
│   │   └── auth.py           # Auth endpoints
│   ├── services/
│   │   ├── __init__.py
│   │   └── user_service.py   # Business logic
│   ├── schemas/
│   │   ├── __init__.py
│   │   └── user.py           # Marshmallow schemas
│   └── utils/
│       ├── __init__.py
│       └── security.py       # Auth utilities
├── migrations/               # Alembic migrations
├── tests/
│   ├── conftest.py
│   ├── test_api/
│   │   └── test_users.py
│   └── test_services/
│       └── test_user_service.py
├── .env.example
├── requirements.txt
├── requirements-dev.txt
├── pyproject.toml
└── run.py                    # Entry point
```

---

## Application Factory

### app/__init__.py
```python
"""Flask application factory."""
from flask import Flask

from app.config import config
from app.extensions import db, migrate, ma, jwt, cors


def create_app(config_name: str = "development") -> Flask:
    """Create and configure the Flask application.

    Args:
        config_name: Configuration environment name.

    Returns:
        Configured Flask application instance.
    """
    app = Flask(__name__)
    app.config.from_object(config[config_name])

    # Initialize extensions
    register_extensions(app)

    # Register blueprints
    register_blueprints(app)

    # Register error handlers
    register_error_handlers(app)

    # Register CLI commands
    register_commands(app)

    return app


def register_extensions(app: Flask) -> None:
    """Initialize Flask extensions."""
    db.init_app(app)
    migrate.init_app(app, db)
    ma.init_app(app)
    jwt.init_app(app)
    cors.init_app(app)


def register_blueprints(app: Flask) -> None:
    """Register Flask blueprints."""
    from app.api import api_bp

    app.register_blueprint(api_bp, url_prefix="/api/v1")


def register_error_handlers(app: Flask) -> None:
    """Register error handlers."""
    from app.utils.errors import (
        handle_validation_error,
        handle_not_found,
        handle_internal_error,
    )
    from marshmallow import ValidationError

    app.register_error_handler(ValidationError, handle_validation_error)
    app.register_error_handler(404, handle_not_found)
    app.register_error_handler(500, handle_internal_error)


def register_commands(app: Flask) -> None:
    """Register CLI commands."""
    from app.commands import seed_db

    app.cli.add_command(seed_db)
```

---

## Configuration

### app/config.py
```python
"""Application configuration."""
import os
from datetime import timedelta
from typing import Type


class Config:
    """Base configuration."""

    # Flask
    SECRET_KEY = os.getenv("SECRET_KEY", "dev-secret-key-change-in-production")

    # SQLAlchemy
    SQLALCHEMY_TRACK_MODIFICATIONS = False
    SQLALCHEMY_ENGINE_OPTIONS = {
        "pool_pre_ping": True,
        "pool_recycle": 300,
    }

    # JWT
    JWT_SECRET_KEY = os.getenv("JWT_SECRET_KEY", SECRET_KEY)
    JWT_ACCESS_TOKEN_EXPIRES = timedelta(hours=1)
    JWT_REFRESH_TOKEN_EXPIRES = timedelta(days=30)
    JWT_TOKEN_LOCATION = ["headers"]

    # CORS
    CORS_ORIGINS = os.getenv("CORS_ORIGINS", "*").split(",")


class DevelopmentConfig(Config):
    """Development configuration."""

    DEBUG = True
    SQLALCHEMY_DATABASE_URI = os.getenv(
        "DATABASE_URL",
        "postgresql://postgres:postgres@localhost:5432/myapp_dev"
    )
    SQLALCHEMY_ECHO = True


class TestingConfig(Config):
    """Testing configuration."""

    TESTING = True
    SQLALCHEMY_DATABASE_URI = os.getenv(
        "TEST_DATABASE_URL",
        "postgresql://postgres:postgres@localhost:5432/myapp_test"
    )
    JWT_ACCESS_TOKEN_EXPIRES = timedelta(seconds=5)


class ProductionConfig(Config):
    """Production configuration."""

    DEBUG = False
    SQLALCHEMY_DATABASE_URI = os.environ["DATABASE_URL"]

    # Stricter settings
    SESSION_COOKIE_SECURE = True
    SESSION_COOKIE_HTTPONLY = True
    SESSION_COOKIE_SAMESITE = "Lax"


config: dict[str, Type[Config]] = {
    "development": DevelopmentConfig,
    "testing": TestingConfig,
    "production": ProductionConfig,
}
```

---

## Extensions

### app/extensions.py
```python
"""Flask extensions initialization."""
from flask_sqlalchemy import SQLAlchemy
from flask_migrate import Migrate
from flask_marshmallow import Marshmallow
from flask_jwt_extended import JWTManager
from flask_cors import CORS

# SQLAlchemy for ORM
db = SQLAlchemy()

# Alembic migrations
migrate = Migrate()

# Marshmallow for serialization
ma = Marshmallow()

# JWT authentication
jwt = JWTManager()

# CORS handling
cors = CORS()
```

---

## Models

### app/models/base.py
```python
"""Base model classes."""
from datetime import datetime
from typing import Optional
import uuid

from sqlalchemy import DateTime, String
from sqlalchemy.orm import Mapped, mapped_column, declared_attr

from app.extensions import db


class TimestampMixin:
    """Mixin for created_at and updated_at timestamps."""

    created_at: Mapped[datetime] = mapped_column(
        DateTime,
        default=datetime.utcnow,
        nullable=False,
    )
    updated_at: Mapped[datetime] = mapped_column(
        DateTime,
        default=datetime.utcnow,
        onupdate=datetime.utcnow,
        nullable=False,
    )


class UUIDMixin:
    """Mixin for UUID primary key."""

    id: Mapped[str] = mapped_column(
        String(36),
        primary_key=True,
        default=lambda: str(uuid.uuid4()),
    )


class SoftDeleteMixin:
    """Mixin for soft delete functionality."""

    deleted_at: Mapped[Optional[datetime]] = mapped_column(
        DateTime,
        nullable=True,
        default=None,
    )

    @property
    def is_deleted(self) -> bool:
        return self.deleted_at is not None

    def soft_delete(self) -> None:
        self.deleted_at = datetime.utcnow()

    def restore(self) -> None:
        self.deleted_at = None


class BaseModel(db.Model, TimestampMixin):
    """Base model with timestamps."""

    __abstract__ = True

    def save(self) -> "BaseModel":
        """Save model to database."""
        db.session.add(self)
        db.session.commit()
        return self

    def delete(self) -> None:
        """Delete model from database."""
        db.session.delete(self)
        db.session.commit()

    @classmethod
    def get_by_id(cls, id: int) -> Optional["BaseModel"]:
        """Get model by ID."""
        return cls.query.get(id)

    @classmethod
    def get_all(cls) -> list["BaseModel"]:
        """Get all models."""
        return cls.query.all()
```

### app/models/user.py
```python
"""User model."""
from typing import Optional

from sqlalchemy import String, Boolean, Integer
from sqlalchemy.orm import Mapped, mapped_column, relationship
from werkzeug.security import generate_password_hash, check_password_hash

from app.models.base import BaseModel, SoftDeleteMixin


class User(BaseModel, SoftDeleteMixin):
    """User model."""

    __tablename__ = "users"

    id: Mapped[int] = mapped_column(Integer, primary_key=True)
    email: Mapped[str] = mapped_column(
        String(255),
        unique=True,
        nullable=False,
        index=True,
    )
    password_hash: Mapped[str] = mapped_column(String(255), nullable=False)
    first_name: Mapped[str] = mapped_column(String(100), nullable=False)
    last_name: Mapped[str] = mapped_column(String(100), nullable=False)
    is_active: Mapped[bool] = mapped_column(Boolean, default=True)
    is_admin: Mapped[bool] = mapped_column(Boolean, default=False)

    # Relationships
    posts: Mapped[list["Post"]] = relationship(
        "Post",
        back_populates="author",
        lazy="dynamic",
    )

    def __repr__(self) -> str:
        return f"<User {self.email}>"

    @property
    def full_name(self) -> str:
        return f"{self.first_name} {self.last_name}"

    def set_password(self, password: str) -> None:
        """Hash and set password."""
        self.password_hash = generate_password_hash(password)

    def check_password(self, password: str) -> bool:
        """Verify password."""
        return check_password_hash(self.password_hash, password)

    @classmethod
    def get_by_email(cls, email: str) -> Optional["User"]:
        """Get user by email."""
        return cls.query.filter_by(email=email, deleted_at=None).first()

    @classmethod
    def get_active_users(cls) -> list["User"]:
        """Get all active users."""
        return cls.query.filter_by(is_active=True, deleted_at=None).all()


class Post(BaseModel):
    """Post model example."""

    __tablename__ = "posts"

    id: Mapped[int] = mapped_column(Integer, primary_key=True)
    title: Mapped[str] = mapped_column(String(200), nullable=False)
    content: Mapped[str] = mapped_column(String, nullable=False)
    author_id: Mapped[int] = mapped_column(
        Integer,
        db.ForeignKey("users.id"),
        nullable=False,
    )

    # Relationships
    author: Mapped["User"] = relationship("User", back_populates="posts")
```

---

## Schemas (Marshmallow)

### app/schemas/user.py
```python
"""User schemas for serialization/validation."""
from marshmallow import fields, validate, validates, ValidationError, post_load

from app.extensions import ma
from app.models.user import User


class UserSchema(ma.SQLAlchemyAutoSchema):
    """Schema for User serialization."""

    class Meta:
        model = User
        load_instance = True
        exclude = ("password_hash", "deleted_at")

    id = fields.Int(dump_only=True)
    email = fields.Email(required=True)
    first_name = fields.Str(required=True, validate=validate.Length(min=1, max=100))
    last_name = fields.Str(required=True, validate=validate.Length(min=1, max=100))
    full_name = fields.Str(dump_only=True)
    is_active = fields.Bool(dump_only=True)
    is_admin = fields.Bool(dump_only=True)
    created_at = fields.DateTime(dump_only=True)
    updated_at = fields.DateTime(dump_only=True)


class UserCreateSchema(ma.Schema):
    """Schema for creating a user."""

    email = fields.Email(required=True)
    password = fields.Str(
        required=True,
        load_only=True,
        validate=validate.Length(min=8, max=128),
    )
    first_name = fields.Str(required=True, validate=validate.Length(min=1, max=100))
    last_name = fields.Str(required=True, validate=validate.Length(min=1, max=100))

    @validates("email")
    def validate_email_unique(self, value: str) -> None:
        """Validate email is unique."""
        if User.get_by_email(value):
            raise ValidationError("Email already registered.")


class UserUpdateSchema(ma.Schema):
    """Schema for updating a user."""

    first_name = fields.Str(validate=validate.Length(min=1, max=100))
    last_name = fields.Str(validate=validate.Length(min=1, max=100))
    is_active = fields.Bool()


class LoginSchema(ma.Schema):
    """Schema for login."""

    email = fields.Email(required=True)
    password = fields.Str(required=True, load_only=True)


class TokenSchema(ma.Schema):
    """Schema for JWT tokens."""

    access_token = fields.Str(required=True)
    refresh_token = fields.Str(required=True)
    token_type = fields.Str(default="bearer")


# Schema instances
user_schema = UserSchema()
users_schema = UserSchema(many=True)
user_create_schema = UserCreateSchema()
user_update_schema = UserUpdateSchema()
login_schema = LoginSchema()
token_schema = TokenSchema()
```

---

## Services

### app/services/user_service.py
```python
"""User business logic."""
from typing import Optional

from flask_jwt_extended import create_access_token, create_refresh_token

from app.extensions import db
from app.models.user import User
from app.utils.errors import NotFoundError, UnauthorizedError, ConflictError


class UserService:
    """Service for user operations."""

    @staticmethod
    def get_user(user_id: int) -> User:
        """Get user by ID.

        Raises:
            NotFoundError: If user not found.
        """
        user = User.query.filter_by(id=user_id, deleted_at=None).first()
        if not user:
            raise NotFoundError(f"User with id {user_id} not found")
        return user

    @staticmethod
    def get_all_users(
        page: int = 1,
        per_page: int = 20,
        active_only: bool = True,
    ) -> tuple[list[User], int]:
        """Get paginated users.

        Returns:
            Tuple of (users, total_count).
        """
        query = User.query.filter_by(deleted_at=None)

        if active_only:
            query = query.filter_by(is_active=True)

        pagination = query.order_by(User.created_at.desc()).paginate(
            page=page,
            per_page=per_page,
            error_out=False,
        )

        return pagination.items, pagination.total

    @staticmethod
    def create_user(
        email: str,
        password: str,
        first_name: str,
        last_name: str,
    ) -> User:
        """Create a new user.

        Raises:
            ConflictError: If email already exists.
        """
        if User.get_by_email(email):
            raise ConflictError(f"User with email {email} already exists")

        user = User(
            email=email,
            first_name=first_name,
            last_name=last_name,
        )
        user.set_password(password)

        db.session.add(user)
        db.session.commit()

        return user

    @staticmethod
    def update_user(user_id: int, **kwargs) -> User:
        """Update user.

        Raises:
            NotFoundError: If user not found.
        """
        user = UserService.get_user(user_id)

        for key, value in kwargs.items():
            if hasattr(user, key) and value is not None:
                setattr(user, key, value)

        db.session.commit()
        return user

    @staticmethod
    def delete_user(user_id: int) -> None:
        """Soft delete user.

        Raises:
            NotFoundError: If user not found.
        """
        user = UserService.get_user(user_id)
        user.soft_delete()
        db.session.commit()

    @staticmethod
    def authenticate(email: str, password: str) -> dict:
        """Authenticate user and return tokens.

        Raises:
            UnauthorizedError: If credentials invalid.
        """
        user = User.get_by_email(email)

        if not user or not user.check_password(password):
            raise UnauthorizedError("Invalid email or password")

        if not user.is_active:
            raise UnauthorizedError("Account is deactivated")

        access_token = create_access_token(
            identity=user.id,
            additional_claims={"is_admin": user.is_admin},
        )
        refresh_token = create_refresh_token(identity=user.id)

        return {
            "access_token": access_token,
            "refresh_token": refresh_token,
            "token_type": "bearer",
        }

    @staticmethod
    def refresh_tokens(user_id: int) -> dict:
        """Refresh access token.

        Raises:
            NotFoundError: If user not found.
            UnauthorizedError: If user inactive.
        """
        user = UserService.get_user(user_id)

        if not user.is_active:
            raise UnauthorizedError("Account is deactivated")

        access_token = create_access_token(
            identity=user.id,
            additional_claims={"is_admin": user.is_admin},
        )

        return {
            "access_token": access_token,
            "token_type": "bearer",
        }
```

---

## API Blueprints

### app/api/__init__.py
```python
"""API blueprint registration."""
from flask import Blueprint

api_bp = Blueprint("api", __name__)

# Import routes to register them
from app.api import users, auth  # noqa: F401, E402
```

### app/api/users.py
```python
"""User API endpoints."""
from flask import request, jsonify
from flask_jwt_extended import jwt_required, get_jwt_identity
from marshmallow import ValidationError

from app.api import api_bp
from app.services.user_service import UserService
from app.schemas.user import (
    user_schema,
    users_schema,
    user_create_schema,
    user_update_schema,
)
from app.utils.decorators import admin_required


@api_bp.route("/users", methods=["GET"])
@jwt_required()
def get_users():
    """Get paginated list of users.

    Query params:
        page: Page number (default: 1)
        per_page: Items per page (default: 20, max: 100)
        active_only: Filter active users only (default: true)

    Returns:
        200: List of users with pagination info.
    """
    page = request.args.get("page", 1, type=int)
    per_page = min(request.args.get("per_page", 20, type=int), 100)
    active_only = request.args.get("active_only", "true").lower() == "true"

    users, total = UserService.get_all_users(
        page=page,
        per_page=per_page,
        active_only=active_only,
    )

    return jsonify({
        "users": users_schema.dump(users),
        "page": page,
        "per_page": per_page,
        "total": total,
        "pages": (total + per_page - 1) // per_page,
    })


@api_bp.route("/users/<int:user_id>", methods=["GET"])
@jwt_required()
def get_user(user_id: int):
    """Get user by ID.

    Returns:
        200: User data.
        404: User not found.
    """
    user = UserService.get_user(user_id)
    return jsonify(user_schema.dump(user))


@api_bp.route("/users", methods=["POST"])
def create_user():
    """Create a new user.

    Returns:
        201: Created user.
        400: Validation error.
        409: Email already exists.
    """
    data = user_create_schema.load(request.get_json())

    user = UserService.create_user(
        email=data["email"],
        password=data["password"],
        first_name=data["first_name"],
        last_name=data["last_name"],
    )

    return jsonify(user_schema.dump(user)), 201


@api_bp.route("/users/<int:user_id>", methods=["PATCH"])
@jwt_required()
def update_user(user_id: int):
    """Update user.

    Returns:
        200: Updated user.
        400: Validation error.
        403: Not authorized.
        404: User not found.
    """
    current_user_id = get_jwt_identity()

    # Users can only update themselves unless admin
    if current_user_id != user_id:
        from flask_jwt_extended import get_jwt
        claims = get_jwt()
        if not claims.get("is_admin"):
            return jsonify({"error": "Not authorized"}), 403

    data = user_update_schema.load(request.get_json())
    user = UserService.update_user(user_id, **data)

    return jsonify(user_schema.dump(user))


@api_bp.route("/users/<int:user_id>", methods=["DELETE"])
@jwt_required()
@admin_required
def delete_user(user_id: int):
    """Delete user (admin only).

    Returns:
        204: User deleted.
        403: Not authorized.
        404: User not found.
    """
    UserService.delete_user(user_id)
    return "", 204


@api_bp.route("/users/me", methods=["GET"])
@jwt_required()
def get_current_user():
    """Get current authenticated user.

    Returns:
        200: Current user data.
    """
    user_id = get_jwt_identity()
    user = UserService.get_user(user_id)
    return jsonify(user_schema.dump(user))
```

### app/api/auth.py
```python
"""Authentication endpoints."""
from flask import request, jsonify
from flask_jwt_extended import jwt_required, get_jwt_identity

from app.api import api_bp
from app.services.user_service import UserService
from app.schemas.user import login_schema, token_schema


@api_bp.route("/auth/login", methods=["POST"])
def login():
    """Authenticate user and return tokens.

    Returns:
        200: Access and refresh tokens.
        400: Validation error.
        401: Invalid credentials.
    """
    data = login_schema.load(request.get_json())
    tokens = UserService.authenticate(
        email=data["email"],
        password=data["password"],
    )
    return jsonify(token_schema.dump(tokens))


@api_bp.route("/auth/refresh", methods=["POST"])
@jwt_required(refresh=True)
def refresh():
    """Refresh access token.

    Returns:
        200: New access token.
        401: Invalid refresh token.
    """
    user_id = get_jwt_identity()
    tokens = UserService.refresh_tokens(user_id)
    return jsonify(tokens)


@api_bp.route("/auth/logout", methods=["POST"])
@jwt_required()
def logout():
    """Logout user (client should discard tokens).

    For token blacklisting, implement a token blocklist.

    Returns:
        200: Logout successful.
    """
    # In production, add token to blocklist
    return jsonify({"message": "Successfully logged out"})
```

---

## Error Handling

### app/utils/errors.py
```python
"""Custom exceptions and error handlers."""
from flask import jsonify
from marshmallow import ValidationError as MarshmallowValidationError


class AppError(Exception):
    """Base application error."""

    def __init__(self, message: str, status_code: int = 400):
        super().__init__(message)
        self.message = message
        self.status_code = status_code


class NotFoundError(AppError):
    """Resource not found error."""

    def __init__(self, message: str = "Resource not found"):
        super().__init__(message, status_code=404)


class UnauthorizedError(AppError):
    """Unauthorized access error."""

    def __init__(self, message: str = "Unauthorized"):
        super().__init__(message, status_code=401)


class ForbiddenError(AppError):
    """Forbidden access error."""

    def __init__(self, message: str = "Forbidden"):
        super().__init__(message, status_code=403)


class ConflictError(AppError):
    """Conflict error (e.g., duplicate resource)."""

    def __init__(self, message: str = "Conflict"):
        super().__init__(message, status_code=409)


def handle_validation_error(error: MarshmallowValidationError):
    """Handle Marshmallow validation errors."""
    return jsonify({
        "error": "Validation error",
        "details": error.messages,
    }), 400


def handle_app_error(error: AppError):
    """Handle application errors."""
    return jsonify({
        "error": error.message,
    }), error.status_code


def handle_not_found(error):
    """Handle 404 errors."""
    return jsonify({
        "error": "Not found",
    }), 404


def handle_internal_error(error):
    """Handle 500 errors."""
    return jsonify({
        "error": "Internal server error",
    }), 500
```

---

## Decorators

### app/utils/decorators.py
```python
"""Custom decorators."""
from functools import wraps
from typing import Callable

from flask import jsonify
from flask_jwt_extended import get_jwt, verify_jwt_in_request


def admin_required(fn: Callable) -> Callable:
    """Decorator to require admin privileges."""
    @wraps(fn)
    def wrapper(*args, **kwargs):
        verify_jwt_in_request()
        claims = get_jwt()
        if not claims.get("is_admin"):
            return jsonify({"error": "Admin access required"}), 403
        return fn(*args, **kwargs)
    return wrapper


def roles_required(*roles: str) -> Callable:
    """Decorator to require specific roles."""
    def decorator(fn: Callable) -> Callable:
        @wraps(fn)
        def wrapper(*args, **kwargs):
            verify_jwt_in_request()
            claims = get_jwt()
            user_role = claims.get("role", "user")
            if user_role not in roles:
                return jsonify({"error": "Insufficient permissions"}), 403
            return fn(*args, **kwargs)
        return wrapper
    return decorator
```

---

## CLI Commands

### app/commands.py
```python
"""Custom CLI commands."""
import click
from flask.cli import with_appcontext

from app.extensions import db
from app.models.user import User


@click.command("seed-db")
@with_appcontext
def seed_db():
    """Seed database with initial data."""
    click.echo("Seeding database...")

    # Create admin user
    admin = User.query.filter_by(email="admin@example.com").first()
    if not admin:
        admin = User(
            email="admin@example.com",
            first_name="Admin",
            last_name="User",
            is_admin=True,
        )
        admin.set_password("admin123")
        db.session.add(admin)
        click.echo("Created admin user")

    # Create test users
    for i in range(1, 6):
        email = f"user{i}@example.com"
        if not User.query.filter_by(email=email).first():
            user = User(
                email=email,
                first_name=f"User",
                last_name=f"{i}",
            )
            user.set_password("password123")
            db.session.add(user)
            click.echo(f"Created user: {email}")

    db.session.commit()
    click.echo("Database seeded successfully!")


@click.command("create-admin")
@click.argument("email")
@click.argument("password")
@with_appcontext
def create_admin(email: str, password: str):
    """Create an admin user."""
    if User.query.filter_by(email=email).first():
        click.echo(f"User {email} already exists")
        return

    user = User(
        email=email,
        first_name="Admin",
        last_name="User",
        is_admin=True,
    )
    user.set_password(password)
    db.session.add(user)
    db.session.commit()

    click.echo(f"Admin user {email} created successfully!")
```

---

## Testing

### tests/conftest.py
```python
"""Test configuration and fixtures."""
import pytest
from flask import Flask
from flask.testing import FlaskClient

from app import create_app
from app.extensions import db
from app.models.user import User


@pytest.fixture(scope="session")
def app() -> Flask:
    """Create test application."""
    app = create_app("testing")

    with app.app_context():
        db.create_all()
        yield app
        db.drop_all()


@pytest.fixture
def client(app: Flask) -> FlaskClient:
    """Create test client."""
    return app.test_client()


@pytest.fixture
def db_session(app: Flask):
    """Create database session for tests."""
    with app.app_context():
        yield db.session
        db.session.rollback()


@pytest.fixture
def user(db_session) -> User:
    """Create test user."""
    user = User(
        email="test@example.com",
        first_name="Test",
        last_name="User",
    )
    user.set_password("password123")
    db_session.add(user)
    db_session.commit()
    return user


@pytest.fixture
def admin_user(db_session) -> User:
    """Create test admin user."""
    user = User(
        email="admin@example.com",
        first_name="Admin",
        last_name="User",
        is_admin=True,
    )
    user.set_password("admin123")
    db_session.add(user)
    db_session.commit()
    return user


@pytest.fixture
def auth_headers(client: FlaskClient, user: User) -> dict:
    """Get auth headers for test user."""
    response = client.post("/api/v1/auth/login", json={
        "email": user.email,
        "password": "password123",
    })
    token = response.get_json()["access_token"]
    return {"Authorization": f"Bearer {token}"}


@pytest.fixture
def admin_headers(client: FlaskClient, admin_user: User) -> dict:
    """Get auth headers for admin user."""
    response = client.post("/api/v1/auth/login", json={
        "email": admin_user.email,
        "password": "admin123",
    })
    token = response.get_json()["access_token"]
    return {"Authorization": f"Bearer {token}"}
```

### tests/test_api/test_users.py
```python
"""Tests for user API endpoints."""
import pytest
from flask.testing import FlaskClient

from app.models.user import User


class TestGetUsers:
    """Tests for GET /users endpoint."""

    def test_get_users_requires_auth(self, client: FlaskClient):
        """Test that authentication is required."""
        response = client.get("/api/v1/users")
        assert response.status_code == 401

    def test_get_users_returns_list(
        self,
        client: FlaskClient,
        auth_headers: dict,
        user: User,
    ):
        """Test successful user list retrieval."""
        response = client.get("/api/v1/users", headers=auth_headers)

        assert response.status_code == 200
        data = response.get_json()
        assert "users" in data
        assert "total" in data
        assert "page" in data

    def test_get_users_pagination(
        self,
        client: FlaskClient,
        auth_headers: dict,
    ):
        """Test pagination parameters."""
        response = client.get(
            "/api/v1/users?page=1&per_page=5",
            headers=auth_headers,
        )

        assert response.status_code == 200
        data = response.get_json()
        assert data["per_page"] == 5


class TestCreateUser:
    """Tests for POST /users endpoint."""

    def test_create_user_success(self, client: FlaskClient):
        """Test successful user creation."""
        response = client.post("/api/v1/users", json={
            "email": "new@example.com",
            "password": "password123",
            "first_name": "New",
            "last_name": "User",
        })

        assert response.status_code == 201
        data = response.get_json()
        assert data["email"] == "new@example.com"
        assert "password" not in data

    def test_create_user_duplicate_email(
        self,
        client: FlaskClient,
        user: User,
    ):
        """Test duplicate email validation."""
        response = client.post("/api/v1/users", json={
            "email": user.email,
            "password": "password123",
            "first_name": "New",
            "last_name": "User",
        })

        assert response.status_code == 400

    def test_create_user_invalid_email(self, client: FlaskClient):
        """Test email validation."""
        response = client.post("/api/v1/users", json={
            "email": "invalid-email",
            "password": "password123",
            "first_name": "New",
            "last_name": "User",
        })

        assert response.status_code == 400


class TestDeleteUser:
    """Tests for DELETE /users endpoint."""

    def test_delete_user_requires_admin(
        self,
        client: FlaskClient,
        auth_headers: dict,
        user: User,
    ):
        """Test that admin is required for deletion."""
        response = client.delete(
            f"/api/v1/users/{user.id}",
            headers=auth_headers,
        )
        assert response.status_code == 403

    def test_delete_user_as_admin(
        self,
        client: FlaskClient,
        admin_headers: dict,
        user: User,
    ):
        """Test admin can delete users."""
        response = client.delete(
            f"/api/v1/users/{user.id}",
            headers=admin_headers,
        )
        assert response.status_code == 204
```

### tests/test_services/test_user_service.py
```python
"""Tests for user service."""
import pytest

from app.services.user_service import UserService
from app.utils.errors import NotFoundError, ConflictError, UnauthorizedError
from app.models.user import User


class TestUserService:
    """Tests for UserService."""

    def test_create_user(self, db_session):
        """Test user creation."""
        user = UserService.create_user(
            email="service_test@example.com",
            password="password123",
            first_name="Service",
            last_name="Test",
        )

        assert user.id is not None
        assert user.email == "service_test@example.com"
        assert user.check_password("password123")

    def test_create_user_duplicate_email(self, user: User):
        """Test duplicate email raises error."""
        with pytest.raises(ConflictError):
            UserService.create_user(
                email=user.email,
                password="password123",
                first_name="Duplicate",
                last_name="User",
            )

    def test_get_user_not_found(self):
        """Test getting non-existent user."""
        with pytest.raises(NotFoundError):
            UserService.get_user(99999)

    def test_authenticate_success(self, user: User):
        """Test successful authentication."""
        tokens = UserService.authenticate(
            email=user.email,
            password="password123",
        )

        assert "access_token" in tokens
        assert "refresh_token" in tokens

    def test_authenticate_wrong_password(self, user: User):
        """Test authentication with wrong password."""
        with pytest.raises(UnauthorizedError):
            UserService.authenticate(
                email=user.email,
                password="wrongpassword",
            )

    def test_soft_delete_user(self, user: User):
        """Test soft delete."""
        UserService.delete_user(user.id)

        assert user.is_deleted

        with pytest.raises(NotFoundError):
            UserService.get_user(user.id)
```

---

## Entry Point

### run.py
```python
"""Application entry point."""
import os

from app import create_app

config_name = os.getenv("FLASK_ENV", "development")
app = create_app(config_name)

if __name__ == "__main__":
    app.run(
        host=os.getenv("HOST", "0.0.0.0"),
        port=int(os.getenv("PORT", 5000)),
        debug=config_name == "development",
    )
```

---

## Dependencies

### requirements.txt
```
Flask>=3.0.0
Flask-SQLAlchemy>=3.1.0
Flask-Migrate>=4.0.0
Flask-Marshmallow>=0.15.0
Flask-JWT-Extended>=4.6.0
Flask-CORS>=4.0.0
marshmallow-sqlalchemy>=0.29.0
psycopg2-binary>=2.9.9
python-dotenv>=1.0.0
gunicorn>=21.0.0
```

### requirements-dev.txt
```
-r requirements.txt
pytest>=7.4.0
pytest-flask>=1.2.0
pytest-cov>=4.1.0
black>=23.0.0
isort>=5.12.0
mypy>=1.5.0
ruff>=0.1.0
```

---

## Commands Reference

```bash
# Install dependencies
pip install -r requirements.txt
pip install -r requirements-dev.txt

# Set environment variables
export FLASK_APP=run.py
export FLASK_ENV=development

# Initialize database
flask db init
flask db migrate -m "Initial migration"
flask db upgrade

# Seed database
flask seed-db

# Create admin user
flask create-admin admin@example.com secretpassword

# Run development server
flask run
# Or
python run.py

# Run with Gunicorn (production)
gunicorn -w 4 -b 0.0.0.0:5000 "app:create_app('production')"

# Run tests
pytest
pytest -v --cov=app --cov-report=html

# Lint and format
black .
isort .
ruff check .
mypy app/

# Database shell
flask shell
```

---

## Best Practices

### Flask-Specific Guidelines
- ✓ Use application factory pattern
- ✓ Use blueprints for modular routing
- ✓ Use Flask extensions appropriately
- ✓ Configure via environment variables
- ✓ Use Marshmallow for validation/serialization
- ✓ Implement proper error handlers
- ✓ Use Flask-Migrate for database migrations
- ✓ Use Flask-JWT-Extended for authentication

### Testing Guidelines
- ✓ Use pytest with Flask test client
- ✓ Use fixtures for test data
- ✓ Test both success and error cases
- ✓ Mock external services
- ✓ Use separate test database

### Security Guidelines
- ✓ Never store plain passwords
- ✓ Use environment variables for secrets
- ✓ Implement proper authentication
- ✓ Validate all user input
- ✓ Use HTTPS in production
- ✓ Set secure cookie options

---

## References

- [Flask Documentation](https://flask.palletsprojects.com/)
- [Flask-SQLAlchemy](https://flask-sqlalchemy.palletsprojects.com/)
- [Flask-Migrate](https://flask-migrate.readthedocs.io/)
- [Flask-JWT-Extended](https://flask-jwt-extended.readthedocs.io/)
- [Marshmallow](https://marshmallow.readthedocs.io/)
- [SQLAlchemy](https://www.sqlalchemy.org/)
