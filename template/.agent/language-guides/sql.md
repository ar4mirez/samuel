# SQL Guide

> **Applies to**: PostgreSQL, MySQL, SQLite, SQL Server, Oracle, MariaDB

---

## Core Principles

1. **Data Integrity**: Constraints, foreign keys, transactions
2. **Query Optimization**: Indexes, query plans, efficient joins
3. **Security First**: Parameterized queries, least privilege
4. **Normalization**: Proper schema design, avoid redundancy
5. **Maintainability**: Clear naming, documentation, migrations

---

## Language-Specific Guardrails

### General SQL Standards
- ✓ Use UPPERCASE for SQL keywords (`SELECT`, `FROM`, `WHERE`)
- ✓ Use snake_case for table and column names
- ✓ Use singular names for tables (`user` not `users`)
- ✓ Always specify column names (no `SELECT *` in production)
- ✓ Use meaningful aliases for complex queries
- ✓ Include comments for complex logic
- ✓ One statement per line for readability

### Schema Design
- ✓ Every table has a primary key
- ✓ Use appropriate data types (don't store numbers as strings)
- ✓ Add foreign key constraints for referential integrity
- ✓ Use `NOT NULL` unless null has semantic meaning
- ✓ Add indexes for frequently queried columns
- ✓ Use `UNIQUE` constraints where appropriate
- ✓ Include `created_at` and `updated_at` timestamps

### Query Best Practices
- ✓ Use parameterized queries (never string concatenation)
- ✓ Use `EXISTS` instead of `IN` for subqueries when possible
- ✓ Use `JOIN` instead of subqueries when appropriate
- ✓ Limit result sets with `LIMIT`/`TOP`
- ✓ Use `EXPLAIN`/`EXPLAIN ANALYZE` to check query plans
- ✓ Avoid `SELECT *` (specify needed columns)
- ✓ Use `COALESCE` for null handling

### Transactions
- ✓ Use transactions for multi-statement operations
- ✓ Keep transactions short (avoid long locks)
- ✓ Handle rollback scenarios
- ✓ Use appropriate isolation levels
- ✓ Avoid deadlocks with consistent ordering

### Security
- ✓ NEVER concatenate user input into queries
- ✓ Use prepared statements/parameterized queries
- ✓ Apply principle of least privilege
- ✓ Audit sensitive data access
- ✓ Encrypt sensitive data at rest
- ✓ Use row-level security when appropriate

---

## Schema Design Patterns

### Table Creation
```sql
-- PostgreSQL example
CREATE TABLE user (
    id              BIGSERIAL PRIMARY KEY,
    email           VARCHAR(255) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    role            VARCHAR(50) NOT NULL DEFAULT 'user',
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT user_role_check CHECK (role IN ('admin', 'user', 'guest'))
);

-- Index for common queries
CREATE INDEX idx_user_email ON user(email);
CREATE INDEX idx_user_role ON user(role) WHERE is_active = true;

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_updated_at
    BEFORE UPDATE ON user
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

### Foreign Keys and Relationships
```sql
-- One-to-many relationship
CREATE TABLE order (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES user(id) ON DELETE CASCADE,
    total_amount    DECIMAL(10, 2) NOT NULL,
    status          VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT order_status_check CHECK (status IN ('pending', 'paid', 'shipped', 'delivered', 'cancelled'))
);

CREATE INDEX idx_order_user_id ON order(user_id);
CREATE INDEX idx_order_status ON order(status);

-- Many-to-many relationship
CREATE TABLE product (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    price           DECIMAL(10, 2) NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_item (
    id              BIGSERIAL PRIMARY KEY,
    order_id        BIGINT NOT NULL REFERENCES order(id) ON DELETE CASCADE,
    product_id      BIGINT NOT NULL REFERENCES product(id) ON DELETE RESTRICT,
    quantity        INTEGER NOT NULL CHECK (quantity > 0),
    unit_price      DECIMAL(10, 2) NOT NULL,

    UNIQUE(order_id, product_id)
);
```

---

## Query Patterns

### Basic CRUD Operations
```sql
-- Create
INSERT INTO user (email, password_hash, role)
VALUES ('user@example.com', '$2b$12$...', 'user')
RETURNING id, email, created_at;

-- Read (single)
SELECT id, email, role, created_at
FROM user
WHERE id = $1 AND is_active = true;

-- Read (list with pagination)
SELECT id, email, role, created_at
FROM user
WHERE is_active = true
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- Update
UPDATE user
SET email = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, email, updated_at;

-- Soft delete
UPDATE user
SET is_active = false, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- Hard delete
DELETE FROM user WHERE id = $1;
```

### Joins
```sql
-- INNER JOIN
SELECT
    o.id AS order_id,
    o.total_amount,
    o.status,
    u.email AS user_email
FROM order o
INNER JOIN user u ON o.user_id = u.id
WHERE o.status = 'pending';

-- LEFT JOIN with aggregation
SELECT
    u.id,
    u.email,
    COUNT(o.id) AS order_count,
    COALESCE(SUM(o.total_amount), 0) AS total_spent
FROM user u
LEFT JOIN order o ON u.id = o.user_id
WHERE u.is_active = true
GROUP BY u.id, u.email
ORDER BY total_spent DESC;

-- Multiple joins
SELECT
    o.id AS order_id,
    u.email,
    p.name AS product_name,
    oi.quantity,
    oi.unit_price,
    (oi.quantity * oi.unit_price) AS line_total
FROM order o
INNER JOIN user u ON o.user_id = u.id
INNER JOIN order_item oi ON o.id = oi.order_id
INNER JOIN product p ON oi.product_id = p.id
WHERE o.id = $1;
```

### Subqueries and CTEs
```sql
-- Subquery
SELECT id, email
FROM user
WHERE id IN (
    SELECT DISTINCT user_id
    FROM order
    WHERE created_at > CURRENT_DATE - INTERVAL '30 days'
);

-- Common Table Expression (CTE) - preferred
WITH recent_orders AS (
    SELECT DISTINCT user_id
    FROM order
    WHERE created_at > CURRENT_DATE - INTERVAL '30 days'
)
SELECT u.id, u.email
FROM user u
INNER JOIN recent_orders ro ON u.id = ro.user_id;

-- Recursive CTE (for hierarchical data)
WITH RECURSIVE category_tree AS (
    -- Base case
    SELECT id, name, parent_id, 1 AS depth
    FROM category
    WHERE parent_id IS NULL

    UNION ALL

    -- Recursive case
    SELECT c.id, c.name, c.parent_id, ct.depth + 1
    FROM category c
    INNER JOIN category_tree ct ON c.parent_id = ct.id
    WHERE ct.depth < 10  -- Prevent infinite recursion
)
SELECT * FROM category_tree ORDER BY depth, name;
```

### Window Functions
```sql
-- Row number for pagination
SELECT
    id,
    email,
    ROW_NUMBER() OVER (ORDER BY created_at DESC) AS row_num
FROM user
WHERE is_active = true;

-- Ranking
SELECT
    u.id,
    u.email,
    SUM(o.total_amount) AS total_spent,
    RANK() OVER (ORDER BY SUM(o.total_amount) DESC) AS spending_rank
FROM user u
INNER JOIN order o ON u.id = o.user_id
GROUP BY u.id, u.email;

-- Running total
SELECT
    id,
    created_at,
    total_amount,
    SUM(total_amount) OVER (ORDER BY created_at) AS running_total
FROM order
WHERE user_id = $1;

-- Partitioned aggregation
SELECT
    user_id,
    created_at,
    total_amount,
    SUM(total_amount) OVER (PARTITION BY user_id ORDER BY created_at) AS user_running_total,
    AVG(total_amount) OVER (PARTITION BY user_id) AS user_avg_order
FROM order;
```

---

## Performance Optimization

### Index Guidelines
```sql
-- Single column index
CREATE INDEX idx_user_email ON user(email);

-- Composite index (order matters!)
CREATE INDEX idx_order_user_status ON order(user_id, status);

-- Partial index (for filtered queries)
CREATE INDEX idx_active_users ON user(email) WHERE is_active = true;

-- Expression index
CREATE INDEX idx_user_email_lower ON user(LOWER(email));

-- Check if index is used
EXPLAIN ANALYZE
SELECT * FROM user WHERE email = 'test@example.com';
```

### Query Optimization Tips
```sql
-- Use EXISTS instead of IN for large subqueries
-- Bad
SELECT * FROM user WHERE id IN (SELECT user_id FROM order);

-- Good
SELECT * FROM user u WHERE EXISTS (
    SELECT 1 FROM order o WHERE o.user_id = u.id
);

-- Use LIMIT for top-N queries
SELECT * FROM order ORDER BY created_at DESC LIMIT 10;

-- Avoid functions on indexed columns in WHERE
-- Bad (can't use index)
SELECT * FROM user WHERE LOWER(email) = 'test@example.com';

-- Good (if you have expression index, or)
SELECT * FROM user WHERE email = 'test@example.com';

-- Use UNION ALL instead of UNION when duplicates are OK
SELECT id FROM active_user
UNION ALL
SELECT id FROM archived_user;
```

---

## Transactions

### Basic Transaction
```sql
BEGIN;

-- Check balance
SELECT balance FROM account WHERE id = $1 FOR UPDATE;

-- Deduct from source
UPDATE account SET balance = balance - $3 WHERE id = $1;

-- Add to destination
UPDATE account SET balance = balance + $3 WHERE id = $2;

-- Record transfer
INSERT INTO transfer (from_account, to_account, amount)
VALUES ($1, $2, $3);

COMMIT;
-- Or ROLLBACK on error
```

### Savepoints
```sql
BEGIN;

INSERT INTO order (user_id, total_amount) VALUES ($1, $2);

SAVEPOINT before_items;

INSERT INTO order_item (order_id, product_id, quantity, unit_price)
VALUES ($3, $4, $5, $6);

-- If item insert fails
ROLLBACK TO SAVEPOINT before_items;
-- Continue with other operations

COMMIT;
```

---

## Migrations

### Migration Best Practices
```sql
-- Always include both up and down migrations

-- Up migration
-- 001_create_user_table.up.sql
CREATE TABLE user (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Down migration
-- 001_create_user_table.down.sql
DROP TABLE IF EXISTS user;

-- Adding columns (non-blocking)
-- 002_add_user_role.up.sql
ALTER TABLE user ADD COLUMN role VARCHAR(50) DEFAULT 'user';

-- 002_add_user_role.down.sql
ALTER TABLE user DROP COLUMN role;

-- Adding NOT NULL to existing column (requires default or backfill)
-- 003_make_role_not_null.up.sql
UPDATE user SET role = 'user' WHERE role IS NULL;
ALTER TABLE user ALTER COLUMN role SET NOT NULL;

-- 003_make_role_not_null.down.sql
ALTER TABLE user ALTER COLUMN role DROP NOT NULL;
```

---

## Database-Specific Notes

### PostgreSQL
```sql
-- UPSERT (INSERT ... ON CONFLICT)
INSERT INTO user (email, role)
VALUES ('test@example.com', 'user')
ON CONFLICT (email) DO UPDATE SET role = EXCLUDED.role
RETURNING *;

-- JSON support
CREATE TABLE event (
    id BIGSERIAL PRIMARY KEY,
    data JSONB NOT NULL
);

SELECT data->>'name' AS name
FROM event
WHERE data @> '{"type": "click"}';

-- Array support
SELECT * FROM product WHERE tags @> ARRAY['featured'];
```

### MySQL
```sql
-- UPSERT (INSERT ... ON DUPLICATE KEY)
INSERT INTO user (email, role)
VALUES ('test@example.com', 'user')
ON DUPLICATE KEY UPDATE role = VALUES(role);

-- Use InnoDB for transactions
CREATE TABLE user (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE
) ENGINE=InnoDB;
```

### SQLite
```sql
-- UPSERT
INSERT INTO user (email, role)
VALUES ('test@example.com', 'user')
ON CONFLICT(email) DO UPDATE SET role = excluded.role;

-- Enable foreign keys (off by default)
PRAGMA foreign_keys = ON;
```

---

## Security

### Parameterized Queries (Application Code)
```python
# Python (psycopg2) - CORRECT
cursor.execute(
    "SELECT * FROM user WHERE email = %s",
    (email,)
)

# Python - WRONG (SQL injection vulnerable)
cursor.execute(f"SELECT * FROM user WHERE email = '{email}'")
```

```javascript
// Node.js (pg) - CORRECT
const result = await client.query(
    'SELECT * FROM user WHERE email = $1',
    [email]
);

// Node.js - WRONG
const result = await client.query(
    `SELECT * FROM user WHERE email = '${email}'`
);
```

### Least Privilege
```sql
-- Create read-only role
CREATE ROLE readonly_user;
GRANT CONNECT ON DATABASE myapp TO readonly_user;
GRANT USAGE ON SCHEMA public TO readonly_user;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO readonly_user;

-- Create application role
CREATE ROLE app_user;
GRANT CONNECT ON DATABASE myapp TO app_user;
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO app_user;
```

---

## Common Anti-Patterns

### Avoid These
```sql
-- SELECT * in production
SELECT * FROM user;  -- Bad: fetches unnecessary data

-- String concatenation for queries
"SELECT * FROM user WHERE email = '" + email + "'"  -- SQL injection!

-- N+1 queries (in application code)
-- Fetching users then looping to fetch orders

-- Missing indexes on foreign keys
-- Missing WHERE clause on UPDATE/DELETE

-- Using LIKE with leading wildcard
SELECT * FROM user WHERE email LIKE '%@gmail.com';  -- Can't use index
```

### Do This Instead
```sql
-- Specify columns
SELECT id, email, role FROM user;

-- Use parameterized queries (see Security section)

-- Use JOINs to avoid N+1
SELECT u.*, o.* FROM user u LEFT JOIN order o ON u.id = o.user_id;

-- Add indexes on foreign keys
CREATE INDEX idx_order_user_id ON order(user_id);

-- Always have WHERE on UPDATE/DELETE
UPDATE user SET role = 'admin' WHERE id = $1;

-- Use suffix wildcards or full-text search
SELECT * FROM user WHERE email LIKE 'john%';  -- Can use index
```

---

## References

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [MySQL Documentation](https://dev.mysql.com/doc/)
- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [Use The Index, Luke](https://use-the-index-luke.com/)
- [SQL Style Guide](https://www.sqlstyle.guide/)
