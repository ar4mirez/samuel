# Django Framework Guide

> **Applies to**: Django 4.2+, Django REST Framework, Django Channels
> **Language**: Python 3.10+
> **Type**: Full-stack Web Framework

---

## Overview

Django is a high-level Python web framework that encourages rapid development and clean, pragmatic design. It follows the "batteries included" philosophy, providing everything needed to build web applications out of the box.

**Use Django when:**
- Building full-stack web applications with templates
- Need an admin interface out of the box
- Want ORM with migrations included
- Building content management systems
- Require session-based authentication

**Consider alternatives when:**
- Building pure APIs (consider FastAPI)
- Need async-first architecture
- Want minimal framework overhead

---

## Project Structure

### Standard Django Project
```
myproject/
├── manage.py
├── pyproject.toml
├── requirements/
│   ├── base.txt
│   ├── dev.txt
│   └── prod.txt
├── config/                    # Project configuration
│   ├── __init__.py
│   ├── settings/
│   │   ├── __init__.py
│   │   ├── base.py
│   │   ├── dev.py
│   │   └── prod.py
│   ├── urls.py
│   ├── wsgi.py
│   └── asgi.py
├── apps/                      # Django applications
│   ├── users/
│   │   ├── __init__.py
│   │   ├── admin.py
│   │   ├── apps.py
│   │   ├── models.py
│   │   ├── views.py
│   │   ├── urls.py
│   │   ├── forms.py
│   │   ├── serializers.py     # If using DRF
│   │   ├── services.py        # Business logic
│   │   └── tests/
│   │       ├── __init__.py
│   │       ├── test_models.py
│   │       ├── test_views.py
│   │       └── test_services.py
│   └── core/                  # Shared utilities
│       ├── __init__.py
│       ├── models.py          # Abstract base models
│       └── mixins.py
├── templates/
│   ├── base.html
│   └── components/
├── static/
│   ├── css/
│   └── js/
└── tests/
    └── conftest.py
```

---

## Settings Configuration

### Base Settings
```python
# config/settings/base.py
from pathlib import Path
from decouple import config, Csv

BASE_DIR = Path(__file__).resolve().parent.parent.parent

SECRET_KEY = config("SECRET_KEY")

INSTALLED_APPS = [
    "django.contrib.admin",
    "django.contrib.auth",
    "django.contrib.contenttypes",
    "django.contrib.sessions",
    "django.contrib.messages",
    "django.contrib.staticfiles",
    # Third-party
    "rest_framework",
    "corsheaders",
    "django_extensions",
    # Local apps
    "apps.users",
    "apps.core",
]

MIDDLEWARE = [
    "django.middleware.security.SecurityMiddleware",
    "whitenoise.middleware.WhiteNoiseMiddleware",
    "corsheaders.middleware.CorsMiddleware",
    "django.contrib.sessions.middleware.SessionMiddleware",
    "django.middleware.common.CommonMiddleware",
    "django.middleware.csrf.CsrfViewMiddleware",
    "django.contrib.auth.middleware.AuthenticationMiddleware",
    "django.contrib.messages.middleware.MessageMiddleware",
    "django.middleware.clickjacking.XFrameOptionsMiddleware",
]

ROOT_URLCONF = "config.urls"

TEMPLATES = [
    {
        "BACKEND": "django.template.backends.django.DjangoTemplates",
        "DIRS": [BASE_DIR / "templates"],
        "APP_DIRS": True,
        "OPTIONS": {
            "context_processors": [
                "django.template.context_processors.debug",
                "django.template.context_processors.request",
                "django.contrib.auth.context_processors.auth",
                "django.contrib.messages.context_processors.messages",
            ],
        },
    },
]

# Database
DATABASES = {
    "default": {
        "ENGINE": "django.db.backends.postgresql",
        "NAME": config("DB_NAME"),
        "USER": config("DB_USER"),
        "PASSWORD": config("DB_PASSWORD"),
        "HOST": config("DB_HOST", default="localhost"),
        "PORT": config("DB_PORT", default="5432"),
    }
}

# Custom user model
AUTH_USER_MODEL = "users.User"

# Static files
STATIC_URL = "static/"
STATIC_ROOT = BASE_DIR / "staticfiles"
STATICFILES_DIRS = [BASE_DIR / "static"]

# Media files
MEDIA_URL = "media/"
MEDIA_ROOT = BASE_DIR / "media"

# REST Framework
REST_FRAMEWORK = {
    "DEFAULT_AUTHENTICATION_CLASSES": [
        "rest_framework.authentication.SessionAuthentication",
        "rest_framework_simplejwt.authentication.JWTAuthentication",
    ],
    "DEFAULT_PERMISSION_CLASSES": [
        "rest_framework.permissions.IsAuthenticated",
    ],
    "DEFAULT_PAGINATION_CLASS": "rest_framework.pagination.PageNumberPagination",
    "PAGE_SIZE": 20,
    "DEFAULT_THROTTLE_RATES": {
        "anon": "100/hour",
        "user": "1000/hour",
    },
}
```

### Development Settings
```python
# config/settings/dev.py
from .base import *

DEBUG = True
ALLOWED_HOSTS = ["localhost", "127.0.0.1"]

# Debug toolbar
INSTALLED_APPS += ["debug_toolbar"]
MIDDLEWARE.insert(0, "debug_toolbar.middleware.DebugToolbarMiddleware")
INTERNAL_IPS = ["127.0.0.1"]

# Email backend for development
EMAIL_BACKEND = "django.core.mail.backends.console.EmailBackend"

# CORS for local development
CORS_ALLOW_ALL_ORIGINS = True
```

---

## Models

### Abstract Base Model
```python
# apps/core/models.py
from django.db import models
import uuid


class TimeStampedModel(models.Model):
    """Abstract base model with created/updated timestamps."""

    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        abstract = True


class UUIDModel(models.Model):
    """Abstract base model with UUID primary key."""

    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)

    class Meta:
        abstract = True


class SoftDeleteModel(models.Model):
    """Abstract base model with soft delete."""

    is_deleted = models.BooleanField(default=False)
    deleted_at = models.DateTimeField(null=True, blank=True)

    class Meta:
        abstract = True

    def delete(self, using=None, keep_parents=False):
        from django.utils import timezone

        self.is_deleted = True
        self.deleted_at = timezone.now()
        self.save(update_fields=["is_deleted", "deleted_at"])

    def hard_delete(self):
        super().delete()
```

### Custom User Model
```python
# apps/users/models.py
from django.contrib.auth.models import AbstractUser
from django.db import models

from apps.core.models import TimeStampedModel, UUIDModel


class User(AbstractUser, UUIDModel, TimeStampedModel):
    """Custom user model with UUID and timestamps."""

    email = models.EmailField(unique=True)
    avatar = models.ImageField(upload_to="avatars/", null=True, blank=True)
    bio = models.TextField(max_length=500, blank=True)

    # Use email for authentication
    USERNAME_FIELD = "email"
    REQUIRED_FIELDS = ["username"]

    class Meta:
        db_table = "users"
        ordering = ["-created_at"]
        indexes = [
            models.Index(fields=["email"]),
            models.Index(fields=["username"]),
        ]

    def __str__(self) -> str:
        return self.email

    @property
    def full_name(self) -> str:
        return f"{self.first_name} {self.last_name}".strip() or self.username


class Profile(TimeStampedModel):
    """Extended user profile."""

    user = models.OneToOneField(User, on_delete=models.CASCADE, related_name="profile")
    phone = models.CharField(max_length=20, blank=True)
    address = models.TextField(blank=True)
    birth_date = models.DateField(null=True, blank=True)

    class Meta:
        db_table = "profiles"

    def __str__(self) -> str:
        return f"Profile of {self.user.email}"
```

### Model with Relationships
```python
# apps/products/models.py
from django.db import models
from django.core.validators import MinValueValidator
from decimal import Decimal

from apps.core.models import TimeStampedModel, UUIDModel


class Category(TimeStampedModel):
    name = models.CharField(max_length=100, unique=True)
    slug = models.SlugField(unique=True)
    parent = models.ForeignKey(
        "self",
        on_delete=models.CASCADE,
        null=True,
        blank=True,
        related_name="children",
    )

    class Meta:
        db_table = "categories"
        verbose_name_plural = "categories"

    def __str__(self) -> str:
        return self.name


class Product(UUIDModel, TimeStampedModel):
    class Status(models.TextChoices):
        DRAFT = "draft", "Draft"
        PUBLISHED = "published", "Published"
        ARCHIVED = "archived", "Archived"

    name = models.CharField(max_length=255)
    slug = models.SlugField(unique=True)
    description = models.TextField()
    price = models.DecimalField(
        max_digits=10,
        decimal_places=2,
        validators=[MinValueValidator(Decimal("0.01"))],
    )
    category = models.ForeignKey(
        Category,
        on_delete=models.PROTECT,
        related_name="products",
    )
    status = models.CharField(
        max_length=20,
        choices=Status.choices,
        default=Status.DRAFT,
    )
    stock = models.PositiveIntegerField(default=0)

    class Meta:
        db_table = "products"
        ordering = ["-created_at"]
        indexes = [
            models.Index(fields=["slug"]),
            models.Index(fields=["status"]),
            models.Index(fields=["category", "status"]),
        ]

    def __str__(self) -> str:
        return self.name

    @property
    def is_available(self) -> bool:
        return self.status == self.Status.PUBLISHED and self.stock > 0
```

---

## Views

### Class-Based Views
```python
# apps/products/views.py
from django.views.generic import ListView, DetailView, CreateView, UpdateView
from django.contrib.auth.mixins import LoginRequiredMixin
from django.urls import reverse_lazy
from django.db.models import Q

from .models import Product
from .forms import ProductForm


class ProductListView(ListView):
    model = Product
    template_name = "products/list.html"
    context_object_name = "products"
    paginate_by = 20

    def get_queryset(self):
        queryset = Product.objects.filter(status=Product.Status.PUBLISHED)

        # Search
        search = self.request.GET.get("q")
        if search:
            queryset = queryset.filter(
                Q(name__icontains=search) | Q(description__icontains=search)
            )

        # Category filter
        category = self.request.GET.get("category")
        if category:
            queryset = queryset.filter(category__slug=category)

        return queryset.select_related("category")


class ProductDetailView(DetailView):
    model = Product
    template_name = "products/detail.html"
    context_object_name = "product"
    slug_url_kwarg = "slug"

    def get_queryset(self):
        return Product.objects.select_related("category").prefetch_related(
            "images", "reviews"
        )


class ProductCreateView(LoginRequiredMixin, CreateView):
    model = Product
    form_class = ProductForm
    template_name = "products/form.html"
    success_url = reverse_lazy("products:list")

    def form_valid(self, form):
        form.instance.created_by = self.request.user
        return super().form_valid(form)
```

### Django REST Framework Views
```python
# apps/products/views_api.py
from rest_framework import viewsets, status, filters
from rest_framework.decorators import action
from rest_framework.response import Response
from rest_framework.permissions import IsAuthenticated, IsAuthenticatedOrReadOnly
from django_filters.rest_framework import DjangoFilterBackend

from .models import Product
from .serializers import ProductSerializer, ProductCreateSerializer
from .filters import ProductFilter


class ProductViewSet(viewsets.ModelViewSet):
    """
    ViewSet for Product CRUD operations.

    list: Get all published products
    retrieve: Get single product by ID
    create: Create new product (authenticated)
    update: Update product (owner only)
    destroy: Delete product (owner only)
    """

    queryset = Product.objects.select_related("category").prefetch_related("images")
    permission_classes = [IsAuthenticatedOrReadOnly]
    filter_backends = [DjangoFilterBackend, filters.SearchFilter, filters.OrderingFilter]
    filterset_class = ProductFilter
    search_fields = ["name", "description"]
    ordering_fields = ["price", "created_at", "name"]
    ordering = ["-created_at"]

    def get_serializer_class(self):
        if self.action in ["create", "update", "partial_update"]:
            return ProductCreateSerializer
        return ProductSerializer

    def get_queryset(self):
        queryset = super().get_queryset()
        if self.action == "list":
            queryset = queryset.filter(status=Product.Status.PUBLISHED)
        return queryset

    def perform_create(self, serializer):
        serializer.save(created_by=self.request.user)

    @action(detail=True, methods=["post"])
    def publish(self, request, pk=None):
        product = self.get_object()
        product.status = Product.Status.PUBLISHED
        product.save(update_fields=["status"])
        return Response({"status": "published"})

    @action(detail=False, methods=["get"])
    def featured(self, request):
        featured = self.get_queryset().filter(is_featured=True)[:10]
        serializer = self.get_serializer(featured, many=True)
        return Response(serializer.data)
```

---

## Serializers

```python
# apps/products/serializers.py
from rest_framework import serializers
from django.utils.text import slugify

from .models import Product, Category


class CategorySerializer(serializers.ModelSerializer):
    class Meta:
        model = Category
        fields = ["id", "name", "slug"]


class ProductSerializer(serializers.ModelSerializer):
    category = CategorySerializer(read_only=True)
    is_available = serializers.BooleanField(read_only=True)

    class Meta:
        model = Product
        fields = [
            "id",
            "name",
            "slug",
            "description",
            "price",
            "category",
            "status",
            "stock",
            "is_available",
            "created_at",
            "updated_at",
        ]
        read_only_fields = ["id", "slug", "created_at", "updated_at"]


class ProductCreateSerializer(serializers.ModelSerializer):
    category_id = serializers.UUIDField(write_only=True)

    class Meta:
        model = Product
        fields = [
            "name",
            "description",
            "price",
            "category_id",
            "stock",
        ]

    def validate_category_id(self, value):
        if not Category.objects.filter(id=value).exists():
            raise serializers.ValidationError("Category not found")
        return value

    def validate_price(self, value):
        if value <= 0:
            raise serializers.ValidationError("Price must be positive")
        return value

    def create(self, validated_data):
        category_id = validated_data.pop("category_id")
        validated_data["category_id"] = category_id
        validated_data["slug"] = slugify(validated_data["name"])
        return super().create(validated_data)
```

---

## Services (Business Logic)

```python
# apps/products/services.py
from django.db import transaction
from django.core.exceptions import ValidationError
from typing import Optional

from .models import Product, Category


class ProductService:
    """Business logic for products."""

    @staticmethod
    @transaction.atomic
    def create_product(
        *,
        name: str,
        description: str,
        price: float,
        category_id: int,
        stock: int = 0,
        created_by,
    ) -> Product:
        """Create a new product with validation."""
        # Validate category exists
        try:
            category = Category.objects.get(id=category_id)
        except Category.DoesNotExist:
            raise ValidationError("Category not found")

        # Check for duplicate name in category
        if Product.objects.filter(name=name, category=category).exists():
            raise ValidationError("Product with this name already exists in category")

        product = Product.objects.create(
            name=name,
            description=description,
            price=price,
            category=category,
            stock=stock,
            created_by=created_by,
        )

        return product

    @staticmethod
    def update_stock(product_id: int, quantity: int) -> Product:
        """Update product stock with optimistic locking."""
        with transaction.atomic():
            product = Product.objects.select_for_update().get(id=product_id)
            new_stock = product.stock + quantity

            if new_stock < 0:
                raise ValidationError("Insufficient stock")

            product.stock = new_stock
            product.save(update_fields=["stock", "updated_at"])

            return product

    @staticmethod
    def get_products_by_category(
        category_slug: str,
        limit: Optional[int] = None,
    ) -> list[Product]:
        """Get published products by category."""
        queryset = Product.objects.filter(
            category__slug=category_slug,
            status=Product.Status.PUBLISHED,
        ).select_related("category")

        if limit:
            queryset = queryset[:limit]

        return list(queryset)
```

---

## Forms

```python
# apps/products/forms.py
from django import forms
from django.core.exceptions import ValidationError
from django.utils.text import slugify

from .models import Product


class ProductForm(forms.ModelForm):
    class Meta:
        model = Product
        fields = ["name", "description", "price", "category", "stock"]
        widgets = {
            "description": forms.Textarea(attrs={"rows": 4}),
            "price": forms.NumberInput(attrs={"step": "0.01", "min": "0"}),
        }

    def clean_name(self):
        name = self.cleaned_data.get("name")
        slug = slugify(name)

        # Check for existing slug (excluding current instance)
        qs = Product.objects.filter(slug=slug)
        if self.instance.pk:
            qs = qs.exclude(pk=self.instance.pk)

        if qs.exists():
            raise ValidationError("A product with this name already exists")

        return name

    def save(self, commit=True):
        instance = super().save(commit=False)
        instance.slug = slugify(instance.name)

        if commit:
            instance.save()
        return instance
```

---

## URL Configuration

```python
# config/urls.py
from django.contrib import admin
from django.urls import path, include
from django.conf import settings
from django.conf.urls.static import static

urlpatterns = [
    path("admin/", admin.site.urls),
    path("api/", include("apps.api.urls")),
    path("", include("apps.products.urls")),
    path("users/", include("apps.users.urls")),
]

if settings.DEBUG:
    urlpatterns += static(settings.MEDIA_URL, document_root=settings.MEDIA_ROOT)
    urlpatterns += [path("__debug__/", include("debug_toolbar.urls"))]


# apps/api/urls.py
from django.urls import path, include
from rest_framework.routers import DefaultRouter
from rest_framework_simplejwt.views import TokenObtainPairView, TokenRefreshView

from apps.products.views_api import ProductViewSet
from apps.users.views_api import UserViewSet

router = DefaultRouter()
router.register(r"products", ProductViewSet)
router.register(r"users", UserViewSet)

urlpatterns = [
    path("", include(router.urls)),
    path("token/", TokenObtainPairView.as_view(), name="token_obtain_pair"),
    path("token/refresh/", TokenRefreshView.as_view(), name="token_refresh"),
]
```

---

## Testing

```python
# apps/products/tests/test_models.py
import pytest
from decimal import Decimal
from django.core.exceptions import ValidationError

from apps.products.models import Product, Category


@pytest.mark.django_db
class TestProductModel:
    def test_create_product(self, category):
        product = Product.objects.create(
            name="Test Product",
            slug="test-product",
            description="Test description",
            price=Decimal("19.99"),
            category=category,
            stock=10,
        )

        assert product.name == "Test Product"
        assert product.price == Decimal("19.99")
        assert product.is_available is False  # Draft status

    def test_product_is_available_when_published_with_stock(self, category):
        product = Product.objects.create(
            name="Available Product",
            slug="available-product",
            description="Test",
            price=Decimal("10.00"),
            category=category,
            status=Product.Status.PUBLISHED,
            stock=5,
        )

        assert product.is_available is True

    def test_product_not_available_when_out_of_stock(self, category):
        product = Product.objects.create(
            name="Out of Stock",
            slug="out-of-stock",
            description="Test",
            price=Decimal("10.00"),
            category=category,
            status=Product.Status.PUBLISHED,
            stock=0,
        )

        assert product.is_available is False


# apps/products/tests/test_views.py
import pytest
from django.urls import reverse
from rest_framework import status


@pytest.mark.django_db
class TestProductAPI:
    def test_list_products(self, api_client, published_products):
        url = reverse("product-list")
        response = api_client.get(url)

        assert response.status_code == status.HTTP_200_OK
        assert len(response.data["results"]) == len(published_products)

    def test_create_product_authenticated(self, authenticated_client, category):
        url = reverse("product-list")
        data = {
            "name": "New Product",
            "description": "New description",
            "price": "29.99",
            "category_id": str(category.id),
            "stock": 10,
        }

        response = authenticated_client.post(url, data)

        assert response.status_code == status.HTTP_201_CREATED
        assert Product.objects.filter(name="New Product").exists()

    def test_create_product_unauthenticated(self, api_client, category):
        url = reverse("product-list")
        data = {"name": "New Product"}

        response = api_client.post(url, data)

        assert response.status_code == status.HTTP_401_UNAUTHORIZED


# conftest.py
import pytest
from rest_framework.test import APIClient
from decimal import Decimal

from apps.users.models import User
from apps.products.models import Product, Category


@pytest.fixture
def api_client():
    return APIClient()


@pytest.fixture
def user(db):
    return User.objects.create_user(
        username="testuser",
        email="test@example.com",
        password="testpass123",
    )


@pytest.fixture
def authenticated_client(api_client, user):
    api_client.force_authenticate(user=user)
    return api_client


@pytest.fixture
def category(db):
    return Category.objects.create(name="Electronics", slug="electronics")


@pytest.fixture
def published_products(db, category):
    products = []
    for i in range(5):
        products.append(
            Product.objects.create(
                name=f"Product {i}",
                slug=f"product-{i}",
                description="Test",
                price=Decimal("10.00"),
                category=category,
                status=Product.Status.PUBLISHED,
                stock=10,
            )
        )
    return products
```

---

## Admin Configuration

```python
# apps/products/admin.py
from django.contrib import admin
from django.utils.html import format_html

from .models import Product, Category


@admin.register(Category)
class CategoryAdmin(admin.ModelAdmin):
    list_display = ["name", "slug", "parent", "product_count"]
    prepopulated_fields = {"slug": ("name",)}
    search_fields = ["name"]

    def product_count(self, obj):
        return obj.products.count()

    product_count.short_description = "Products"


@admin.register(Product)
class ProductAdmin(admin.ModelAdmin):
    list_display = ["name", "category", "price", "stock", "status", "is_available"]
    list_filter = ["status", "category", "created_at"]
    search_fields = ["name", "description"]
    prepopulated_fields = {"slug": ("name",)}
    readonly_fields = ["created_at", "updated_at"]
    ordering = ["-created_at"]

    fieldsets = (
        (None, {"fields": ("name", "slug", "description")}),
        ("Pricing", {"fields": ("price", "stock")}),
        ("Classification", {"fields": ("category", "status")}),
        ("Timestamps", {"fields": ("created_at", "updated_at"), "classes": ("collapse",)}),
    )

    def is_available(self, obj):
        if obj.is_available:
            return format_html('<span style="color: green;">✓</span>')
        return format_html('<span style="color: red;">✗</span>')

    is_available.short_description = "Available"
```

---

## Management Commands

```python
# apps/products/management/commands/import_products.py
from django.core.management.base import BaseCommand, CommandError
import csv
from decimal import Decimal

from apps.products.models import Product, Category


class Command(BaseCommand):
    help = "Import products from CSV file"

    def add_arguments(self, parser):
        parser.add_argument("csv_file", type=str, help="Path to CSV file")
        parser.add_argument(
            "--dry-run",
            action="store_true",
            help="Preview import without saving",
        )

    def handle(self, *args, **options):
        csv_file = options["csv_file"]
        dry_run = options["dry_run"]

        try:
            with open(csv_file, "r") as f:
                reader = csv.DictReader(f)
                products = []

                for row in reader:
                    category = Category.objects.get(slug=row["category_slug"])
                    product = Product(
                        name=row["name"],
                        slug=row["slug"],
                        description=row["description"],
                        price=Decimal(row["price"]),
                        category=category,
                        stock=int(row["stock"]),
                    )
                    products.append(product)

                if dry_run:
                    self.stdout.write(f"Would import {len(products)} products")
                else:
                    Product.objects.bulk_create(products)
                    self.stdout.write(
                        self.style.SUCCESS(f"Imported {len(products)} products")
                    )

        except FileNotFoundError:
            raise CommandError(f"File not found: {csv_file}")
        except Category.DoesNotExist:
            raise CommandError(f"Category not found: {row['category_slug']}")
```

---

## Signals

```python
# apps/products/signals.py
from django.db.models.signals import post_save, pre_delete
from django.dispatch import receiver
from django.core.cache import cache

from .models import Product


@receiver(post_save, sender=Product)
def invalidate_product_cache(sender, instance, **kwargs):
    """Invalidate cache when product is saved."""
    cache.delete(f"product_{instance.id}")
    cache.delete("featured_products")


@receiver(pre_delete, sender=Product)
def cleanup_product_files(sender, instance, **kwargs):
    """Clean up associated files before deletion."""
    for image in instance.images.all():
        image.file.delete(save=False)


# apps/products/apps.py
from django.apps import AppConfig


class ProductsConfig(AppConfig):
    default_auto_field = "django.db.models.BigAutoField"
    name = "apps.products"

    def ready(self):
        import apps.products.signals  # noqa
```

---

## Middleware

```python
# apps/core/middleware.py
import time
import logging

logger = logging.getLogger(__name__)


class RequestTimingMiddleware:
    """Log request timing."""

    def __init__(self, get_response):
        self.get_response = get_response

    def __call__(self, request):
        start_time = time.time()
        response = self.get_response(request)
        duration = time.time() - start_time

        logger.info(
            f"{request.method} {request.path} - {response.status_code} - {duration:.3f}s"
        )

        return response


class CurrentUserMiddleware:
    """Make current user available in thread local."""

    def __init__(self, get_response):
        self.get_response = get_response

    def __call__(self, request):
        from threading import local

        _thread_locals = local()
        _thread_locals.user = getattr(request, "user", None)

        return self.get_response(request)
```

---

## Celery Tasks

```python
# apps/products/tasks.py
from celery import shared_task
from django.core.mail import send_mail
from django.conf import settings

from .models import Product


@shared_task
def send_low_stock_alert(product_id: int):
    """Send alert when product stock is low."""
    product = Product.objects.get(id=product_id)

    if product.stock <= 5:
        send_mail(
            subject=f"Low stock alert: {product.name}",
            message=f"Product {product.name} has only {product.stock} items left.",
            from_email=settings.DEFAULT_FROM_EMAIL,
            recipient_list=[settings.ADMIN_EMAIL],
        )


@shared_task
def update_product_statistics():
    """Daily task to update product statistics."""
    from django.db.models import Count, Avg

    products = Product.objects.annotate(
        review_count=Count("reviews"),
        avg_rating=Avg("reviews__rating"),
    )

    for product in products:
        product.review_count_cached = product.review_count
        product.avg_rating_cached = product.avg_rating or 0

    Product.objects.bulk_update(
        products, ["review_count_cached", "avg_rating_cached"]
    )
```

---

## Dependencies

```txt
# requirements/base.txt
Django>=4.2,<5.0
djangorestframework>=3.14
django-cors-headers>=4.0
django-filter>=23.0
djangorestframework-simplejwt>=5.3
python-decouple>=3.8
psycopg2-binary>=2.9
Pillow>=10.0
celery>=5.3
redis>=5.0
whitenoise>=6.6

# requirements/dev.txt
-r base.txt
pytest>=7.4
pytest-django>=4.5
pytest-cov>=4.1
factory-boy>=3.3
django-debug-toolbar>=4.2
django-extensions>=3.2
ipython>=8.0
black>=23.0
ruff>=0.1.0
mypy>=1.5
django-stubs>=4.2
```

---

## Commands Reference

```bash
# Development
python manage.py runserver
python manage.py shell_plus  # django-extensions

# Migrations
python manage.py makemigrations
python manage.py migrate
python manage.py showmigrations

# Testing
pytest
pytest -v --cov=apps --cov-report=html
pytest apps/products/ -k "test_create"

# Database
python manage.py dbshell
python manage.py dumpdata products > fixtures/products.json
python manage.py loaddata fixtures/products.json

# Static files
python manage.py collectstatic

# Custom commands
python manage.py import_products data/products.csv --dry-run

# Celery
celery -A config worker -l info
celery -A config beat -l info
```

---

## Best Practices

### Do
- ✓ Use `select_related` and `prefetch_related` for queries
- ✓ Create indexes for frequently queried fields
- ✓ Use transactions for multi-step operations
- ✓ Validate at model level and serializer level
- ✓ Write services for business logic
- ✓ Use signals sparingly (prefer explicit calls)
- ✓ Cache expensive queries
- ✓ Use environment variables for configuration

### Don't
- ❌ Put business logic in views
- ❌ Use raw SQL without parameterization
- ❌ Ignore N+1 query problems
- ❌ Store sensitive data in settings files
- ❌ Use `Model.objects.all()` without limits
- ❌ Skip migrations in production

---

## References

- [Django Documentation](https://docs.djangoproject.com/)
- [Django REST Framework](https://www.django-rest-framework.org/)
- [Two Scoops of Django](https://www.feldroy.com/books/two-scoops-of-django-3-x)
- [Django Best Practices](https://django-best-practices.readthedocs.io/)
