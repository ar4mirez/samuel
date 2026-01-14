# Python Guide

> **Applies to**: Python 3.9+, Django, FastAPI, Flask, Data Science

---

## Core Principles

1. **Readability Counts**: "Code is read more often than written" (PEP 20)
2. **Type Hints**: Use type annotations for better IDE support and validation
3. **Explicit is Better**: Avoid magic, prefer clarity over cleverness
4. **Batteries Included**: Use standard library when possible
5. **Virtual Environments**: Always isolate dependencies

---

## Language-Specific Guardrails

### Python Version & Setup
- ✓ Use Python 3.9+ (3.11+ recommended for performance)
- ✓ Use virtual environments (venv, poetry, or conda)
- ✓ Pin dependencies in requirements.txt or pyproject.toml
- ✓ Include python version in README or .python-version file

### Type Hints (PEP 484)
- ✓ All function signatures have type hints
- ✓ Use `mypy` for static type checking
- ✓ Use `Optional[T]` or `T | None` for nullable types
- ✓ Complex types via `typing` module: `List`, `Dict`, `Union`, `Callable`
- ✓ Type aliases for complex types: `UserId = int`

### Code Style (PEP 8)
- ✓ Follow PEP 8 style guide (enforced by Black)
- ✓ Line length: 88 characters (Black default) or 100 max
- ✓ Use `snake_case` for functions and variables
- ✓ Use `PascalCase` for classes
- ✓ Use `SCREAMING_SNAKE_CASE` for constants
- ✓ 2 blank lines between top-level definitions
- ✓ Import order: stdlib → third-party → local (isort handles this)

### Data Validation
- ✓ Use Pydantic for data models (FastAPI, APIs)
- ✓ Use dataclasses for simple data structures
- ✓ Validate user inputs before processing
- ✓ Use `attrs` or `pydantic` over manual `__init__`

### Django-Specific (if applicable)
- ✓ Follow Django's MTV (Model-Template-View) pattern
- ✓ Use Django ORM (avoid raw SQL without justification)
- ✓ Settings via `django-environ` or environment variables
- ✓ Model files ≤300 lines (split into multiple apps)
- ✓ Use `select_related` and `prefetch_related` to avoid N+1
- ✓ All views have permission checks

### FastAPI-Specific (if applicable)
- ✓ Use Pydantic models for request/response validation
- ✓ Dependency injection for shared logic (database, auth)
- ✓ Async endpoints for I/O-bound operations
- ✓ OpenAPI docs automatically generated and reviewed

---

## Validation & Input Handling

### Recommended Libraries
- **Pydantic**: Data validation with type hints
- **Marshmallow**: Serialization/deserialization (older projects)
- **attrs**: Better classes with validation

### Pattern
```python
from pydantic import BaseModel, EmailStr, validator

class UserCreate(BaseModel):
    email: EmailStr
    age: int
    role: str

    @validator('age')
    def age_must_be_positive(cls, v):
        if v <= 0:
            raise ValueError('Age must be positive')
        return v

    @validator('role')
    def role_must_be_valid(cls, v):
        if v not in ['admin', 'user', 'guest']:
            raise ValueError('Invalid role')
        return v

# Usage
def create_user(data: dict) -> User:
    validated = UserCreate(**data)  # Raises ValidationError if invalid
    return User.objects.create(**validated.dict())
```

---

## Testing

### Frameworks
- **pytest**: Industry standard (recommended)
- **unittest**: Built-in, verbose but works
- **hypothesis**: Property-based testing

### Guardrails
- ✓ Test files: `test_*.py` or `*_test.py` (pytest convention)
- ✓ Test functions: `test_function_name_should_do_something`
- ✓ Use fixtures for setup: `@pytest.fixture`
- ✓ Mock external dependencies: `unittest.mock` or `pytest-mock`
- ✓ Test both success and error paths
- ✓ Coverage target: >80% for business logic
- ✓ Use `pytest-cov` for coverage reports

### Example
```python
import pytest
from myapp.services import UserService
from myapp.exceptions import InvalidEmailError

@pytest.fixture
def user_service():
    return UserService()

def test_create_user_with_valid_data(user_service):
    user = user_service.create(
        email='test@example.com',
        age=25
    )
    assert user.id is not None
    assert user.email == 'test@example.com'

def test_create_user_with_invalid_email_raises_error(user_service):
    with pytest.raises(InvalidEmailError):
        user_service.create(
            email='invalid',
            age=25
        )

@pytest.mark.parametrize('age', [0, -1, -100])
def test_create_user_with_invalid_age(user_service, age):
    with pytest.raises(ValueError, match='Age must be positive'):
        user_service.create(email='test@example.com', age=age)
```

---

## Tooling

### Essential Tools
- **Black**: Code formatter (opinionated, consistent)
- **isort**: Import sorting
- **mypy**: Static type checker
- **ruff**: Fast linter (replaces flake8, pylint)
- **poetry** or **pip-tools**: Dependency management

### Configuration Files
```toml
# pyproject.toml
[tool.black]
line-length = 88
target-version = ['py311']

[tool.isort]
profile = "black"
line_length = 88

[tool.mypy]
python_version = "3.11"
strict = true
warn_return_any = true
warn_unused_configs = true

[tool.pytest.ini_options]
testpaths = ["tests"]
python_files = ["test_*.py", "*_test.py"]
addopts = "--cov=myapp --cov-report=html --cov-report=term"
```

### Pre-Commit Commands
```bash
# Format
black .

# Sort imports
isort .

# Type check
mypy .

# Lint
ruff check .

# Test
pytest
```

---

## Common Pitfalls

### ❌ Don't Do This
```python
# No type hints
def calculate(a, b):
    return a + b

# Mutable default arguments
def add_item(item, items=[]):
    items.append(item)
    return items

# Bare except
try:
    risky_operation()
except:
    pass

# Using `==` for None
if value == None:
    pass

# Not using context managers
f = open('file.txt')
data = f.read()
f.close()
```

### ✅ Do This Instead
```python
# Proper type hints
def calculate(a: int, b: int) -> int:
    return a + b

# Immutable defaults
def add_item(item: str, items: list[str] | None = None) -> list[str]:
    if items is None:
        items = []
    items.append(item)
    return items

# Specific exception handling
try:
    risky_operation()
except ValueError as e:
    logger.error(f'Operation failed: {e}')
    raise

# Use `is` for None
if value is None:
    pass

# Context managers
with open('file.txt') as f:
    data = f.read()
```

---

## Framework-Specific Patterns

### Django
```python
# models.py
from django.db import models

class User(models.Model):
    email = models.EmailField(unique=True)
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        ordering = ['-created_at']
        indexes = [models.Index(fields=['email'])]

# views.py
from django.contrib.auth.decorators import login_required
from django.http import JsonResponse

@login_required
def user_detail(request, user_id: int):
    user = User.objects.select_related('profile').get(id=user_id)
    return JsonResponse({
        'id': user.id,
        'email': user.email,
    })
```

### FastAPI
```python
from fastapi import FastAPI, Depends, HTTPException
from pydantic import BaseModel

app = FastAPI()

class UserCreate(BaseModel):
    email: str
    age: int

class UserResponse(BaseModel):
    id: int
    email: str

    class Config:
        orm_mode = True

@app.post('/users', response_model=UserResponse)
async def create_user(user: UserCreate):
    # user is automatically validated
    db_user = await create_user_in_db(user)
    return db_user

@app.get('/users/{user_id}', response_model=UserResponse)
async def get_user(user_id: int):
    user = await get_user_from_db(user_id)
    if not user:
        raise HTTPException(status_code=404, detail='User not found')
    return user
```

---

## Performance Considerations

### Optimization Guardrails
- ✓ Use async/await for I/O-bound operations
- ✓ Database: Use `select_related`, `prefetch_related` (Django)
- ✓ Pagination for large datasets (not `all()`)
- ✓ Cache expensive computations: `functools.lru_cache`
- ✓ Profile before optimizing: `cProfile`, `line_profiler`
- ✓ Use generators for large data processing

### Example
```python
from functools import lru_cache
from typing import Iterator

# Caching
@lru_cache(maxsize=128)
def expensive_computation(n: int) -> int:
    # Complex calculation
    return result

# Generators for memory efficiency
def process_large_file(filename: str) -> Iterator[dict]:
    with open(filename) as f:
        for line in f:
            yield process_line(line)

# Async I/O
async def fetch_users() -> list[User]:
    async with aiohttp.ClientSession() as session:
        async with session.get('/api/users') as response:
            return await response.json()
```

---

## Security Best Practices

### Guardrails
- ✓ Never use `eval()` or `exec()` on user input
- ✓ Use parameterized queries (ORM prevents SQL injection)
- ✓ Hash passwords: `bcrypt` or `argon2`
- ✓ Validate all user inputs with Pydantic/Marshmallow
- ✓ Use environment variables for secrets (python-decouple, django-environ)
- ✓ CSRF protection enabled (Django middleware)
- ✓ SQL injection prevention (use ORM or parameterized queries)

### Example
```python
# Password hashing
from passlib.context import CryptContext

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

def hash_password(password: str) -> str:
    return pwd_context.hash(password)

def verify_password(plain: str, hashed: str) -> bool:
    return pwd_context.verify(plain, hashed)

# Environment variables
from decouple import config

DATABASE_URL = config('DATABASE_URL')
SECRET_KEY = config('SECRET_KEY')
DEBUG = config('DEBUG', default=False, cast=bool)
```

---

## Dependency Management

### Using Poetry (Recommended)
```toml
# pyproject.toml
[tool.poetry]
name = "myproject"
version = "0.1.0"
description = ""
authors = ["Your Name <you@example.com>"]

[tool.poetry.dependencies]
python = "^3.11"
fastapi = "^0.104.0"
pydantic = "^2.0.0"

[tool.poetry.group.dev.dependencies]
pytest = "^7.4.0"
black = "^23.0.0"
mypy = "^1.5.0"
```

Commands:
```bash
poetry install              # Install dependencies
poetry add requests         # Add new dependency
poetry add --group dev pytest  # Add dev dependency
poetry update              # Update dependencies
poetry lock                # Update lock file
```

---

## Data Science Specific

### Libraries
- **NumPy**: Numerical computing
- **Pandas**: Data manipulation
- **Scikit-learn**: Machine learning
- **Matplotlib/Seaborn**: Visualization

### Guardrails
- ✓ Use type hints even in notebooks (helps with refactoring)
- ✓ Vectorize operations (avoid Python loops on arrays)
- ✓ Use `pandas` methods over iteration
- ✓ Set random seeds for reproducibility
- ✓ Version data alongside code (DVC, Git LFS)

```python
import numpy as np
import pandas as pd

# Vectorization (fast)
df['total'] = df['price'] * df['quantity']

# NOT this (slow)
for index, row in df.iterrows():
    df.at[index, 'total'] = row['price'] * row['quantity']

# Reproducibility
np.random.seed(42)
```

---

## References

- [PEP 8 Style Guide](https://peps.python.org/pep-0008/)
- [Python Type Hints](https://docs.python.org/3/library/typing.html)
- [Pydantic Documentation](https://docs.pydantic.dev/)
- [pytest Documentation](https://docs.pytest.org/)
- [Django Best Practices](https://django-best-practices.readthedocs.io/)
- [FastAPI Documentation](https://fastapi.tiangolo.com/)
