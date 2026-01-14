# Project Initialization Workflow

> **Purpose**: Set up new projects or adopt CLAUDE.md system in existing projects
>
> **When to use**: First time setup, onboarding existing codebase

---

## For New Projects

### Discovery Questions

**AI will ask these questions to set up the project:**

1. **What tech stack?**
   - Language: TypeScript, Python, Go, Rust, other?
   - Framework: React, Django, Express, Axum, other?
   - Versions: Which versions should we target?

2. **What architecture?**
   - Monolith (single deployable)
   - Microservices (multiple services)
   - Serverless (functions/lambdas)
   - Jamstack (static + APIs)
   - Other?

3. **What testing approach?**
   - Unit tests (which framework?)
   - Integration tests needed?
   - E2E tests (Playwright, Cypress)?
   - Coverage targets?

4. **What deployment target?**
   - Cloud: AWS, GCP, Azure, Vercel, Railway?
   - On-premise servers?
   - Hybrid?
   - Container platform: Docker, Kubernetes?

5. **What database?**
   - PostgreSQL, MySQL, MongoDB, SQLite?
   - ORM/Query builder preferences?
   - Migrations approach?

6. **Additional requirements?**
   - Authentication needed?
   - API type: REST, GraphQL, gRPC?
   - Real-time features (WebSockets)?
   - Background jobs?

### AI Will Create

**1. `.agent/project.md`**
Document all answers from discovery questions.

**2. Initial Directory Structure**

**TypeScript/Node.js:**
```
project/
├── src/
│   ├── index.ts
│   ├── config/
│   ├── routes/
│   ├── services/
│   └── types/
├── tests/
├── package.json
├── tsconfig.json
├── .eslintrc.json
├── .prettierrc
├── .gitignore
└── README.md
```

**Python:**
```
project/
├── src/
│   └── __init__.py
├── tests/
├── requirements.txt
├── pyproject.toml
├── .gitignore
└── README.md
```

**Go:**
```
project/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
├── pkg/
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

**Rust:**
```
project/
├── src/
│   └── main.rs
├── tests/
├── Cargo.toml
├── Cargo.lock
├── .gitignore
└── README.md
```

**3. Configuration Files**

**TypeScript:**
- `tsconfig.json` (strict mode enabled)
- `package.json` (with dev dependencies: ESLint, Prettier, testing)
- `.eslintrc.json` (recommended rules)
- `.prettierrc` (consistent formatting)

**Python:**
- `pyproject.toml` (Black, isort, mypy config)
- `requirements.txt` or `poetry.toml`
- `.python-version` (pin Python version)

**Go:**
- `go.mod` (module definition)
- `.golangci.yml` (linter config)

**Rust:**
- `Cargo.toml` (dependencies, profile settings)
- `.clippy.toml` (linter config)

**4. Essential Files**

**`.gitignore`:**
```
# Dependencies
node_modules/
venv/
vendor/
target/

# Environment
.env
.env.local

# Build artifacts
dist/
build/
*.pyc
__pycache__/

# IDE
.vscode/
.idea/
*.swp

# OS
.DS_Store
Thumbs.db

# Logs
*.log
logs/
```

**`.env.example`:**
```bash
# Database
DATABASE_URL=postgresql://localhost:5432/mydb

# API Keys (never commit real keys)
API_KEY=your_api_key_here

# Environment
NODE_ENV=development
PORT=3000
```

**`README.md`:**
```markdown
# Project Name

## Setup

1. Install dependencies:
   \`\`\`bash
   npm install  # or pip install -r requirements.txt, cargo build
   \`\`\`

2. Copy environment variables:
   \`\`\`bash
   cp .env.example .env
   # Edit .env with your values
   \`\`\`

3. Run migrations (if applicable):
   \`\`\`bash
   npm run migrate  # or python manage.py migrate
   \`\`\`

4. Start development server:
   \`\`\`bash
   npm run dev  # or python manage.py runserver, cargo run
   \`\`\`

## Testing

\`\`\`bash
npm test  # or pytest, go test, cargo test
\`\`\`

## Deployment

[Deployment instructions]
```

---

## For Existing Projects

### Analysis Process

**AI will automatically:**

1. **Scan for Tech Stack**
   ```bash
   # Check for package managers
   ls package.json requirements.txt go.mod Cargo.toml

   # Identify language
   find . -name "*.ts" -o -name "*.py" -o -name "*.go" -o -name "*.rs"

   # Check framework (look for imports/dependencies)
   ```

2. **Examine Directory Structure**
   - Frontend? (src/, components/, pages/)
   - Backend? (api/, routes/, controllers/)
   - Monorepo? (packages/, apps/)
   - Testing setup? (tests/, __tests__, *_test.go)

3. **Analyze Code Patterns**
   ```bash
   # Naming conventions
   grep -r "function\|class\|interface" src/

   # Testing patterns
   head -20 tests/*.test.ts

   # Database usage
   grep -r "SELECT\|INSERT\|UPDATE" src/
   ```

4. **Review Recent Commits**
   ```bash
   # Commit message style
   git log --pretty=format:"%s" -10

   # Code review practices
   git log --merges -10

   # Branching strategy
   git branch -a
   ```

5. **Check Existing Tooling**
   - Linters: ESLint, Ruff, golangci-lint, Clippy?
   - Formatters: Prettier, Black, gofmt, rustfmt?
   - CI/CD: GitHub Actions, GitLab CI, CircleCI?
   - Testing: Jest, pytest, Go testing, cargo test?

### AI Will Create

**1. `.agent/project.md`**
Document discovered tech stack and patterns.

**Example:**
```markdown
# Project: ExistingApp

## Discovered Tech Stack
- **Language**: TypeScript 5.1
- **Framework**: Next.js 14 (App Router)
- **Database**: PostgreSQL with Prisma ORM
- **Testing**: Jest + React Testing Library
- **Styling**: Tailwind CSS
- **Deployment**: Vercel

## Observed Patterns
- API routes in `app/api/`
- Components use `use client` directive
- Server actions for mutations
- Zod for API validation
- Conventional commits (mostly followed)

## Gaps Identified
- Test coverage: 45% (target: >60%)
- No pre-commit hooks
- ESLint config basic (missing recommended rules)
- Some files exceed 300 lines (needs refactoring)
- Missing .env.example
```

**2. `.agent/patterns.md`**
Extract coding patterns from existing codebase.

**Example:**
```markdown
# Patterns

## API Error Handling
All API routes use this pattern:
\`\`\`typescript
try {
  const data = await validateInput(req.body);
  const result = await service.process(data);
  return NextResponse.json(result);
} catch (error) {
  return handleAPIError(error);
}
\`\`\`

## Database Queries
Always use Prisma with error handling:
\`\`\`typescript
const user = await prisma.user.findUnique({
  where: { id }
}).catch(handlePrismaError);
\`\`\`
```

**3. Gap Analysis & Recommendations**

**AI will propose:**
- Add missing configuration files (`.env.example`, linter configs)
- Improve test coverage (identify untested modules)
- Refactor oversized files (list files >300 lines)
- Add pre-commit hooks (husky + lint-staged)
- Update dependencies (security patches)
- Improve documentation (missing README sections)

### Adoption Checklist

**User confirms or corrects AI's understanding:**

- [ ] Tech stack correctly identified?
- [ ] Architecture pattern accurate?
- [ ] Code conventions match observations?
- [ ] Testing approach correct?
- [ ] Any custom patterns AI should know?
- [ ] Any legacy code to avoid modifying?
- [ ] Any protected files/directories?

**Then AI proceeds to:**
- [ ] Apply guardrails to new code (don't refactor existing immediately)
- [ ] Suggest incremental improvements
- [ ] Document conventions in `.agent/`
- [ ] Set up recommended tooling (optional)

---

## Post-Initialization

### Verify Setup

**New Projects:**
```bash
# Install dependencies
npm install  # or pip install -r requirements.txt

# Run tests (should pass even if empty)
npm test

# Start dev server
npm run dev

# Check linting
npm run lint
```

**Existing Projects:**
```bash
# Verify .agent/ structure created
ls .agent/

# Read project documentation
cat .agent/project.md

# Review proposed improvements
cat .agent/patterns.md

# Check gaps to address
# (listed in project.md or as TODOs)
```

### Next Steps

1. **Review `.agent/project.md`**
   - Confirm tech stack accurate
   - Add any missing context

2. **Review `.agent/patterns.md`** (if existing project)
   - Confirm patterns are correct
   - Add any unobserved patterns

3. **Prioritize Improvements** (if existing project)
   - Security issues: Fix immediately
   - Test coverage: Incremental improvement
   - Refactoring: Gradual (as you touch files)
   - Tooling: Set up recommended tools

4. **Start Development**
   - AI will now follow guardrails
   - AI will load language guide automatically
   - AI will reference patterns from `.agent/`

---

## Common Initialization Issues

### Issue: Dependencies Won't Install

**Node.js:**
```bash
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
```

**Python:**
```bash
python -m venv venv
source venv/bin/activate  # or venv\Scripts\activate on Windows
pip install --upgrade pip
pip install -r requirements.txt
```

**Go:**
```bash
go clean -modcache
go mod download
```

**Rust:**
```bash
cargo clean
cargo build
```

### Issue: Wrong Node/Python/Go/Rust Version

**Use version managers:**
```bash
# Node.js
nvm install 20
nvm use 20

# Python
pyenv install 3.11
pyenv local 3.11

# Go
gvm install go1.21
gvm use go1.21

# Rust
rustup update stable
```

### Issue: AI Not Detecting Tech Stack

**Manually specify in `.agent/project.md`:**
```markdown
# Project: MyApp

**Tech Stack**: TypeScript + React + Express
**Database**: PostgreSQL
**Testing**: Jest

[Add to top of project.md]
```

Then tell AI: `"I've updated project.md with our tech stack. Please load @.agent/language-guides/typescript.md"`

---

## Customization

### Modify Initial Structure

Edit this file (`.agent/workflows/initialize-project.md`) to customize:
- Default directory structure
- Configuration file templates
- Questions AI asks
- Files AI creates

### Add Company Standards

Create `.agent/company-standards.md`:
```markdown
# Company Standards

## Required Tools
- ESLint with company config
- Prettier with company config
- Husky for pre-commit hooks

## Required Files
- SECURITY.md
- CONTRIBUTING.md
- LICENSE (MIT)

## Deployment
- All projects deploy to AWS
- Use Terraform for infrastructure
- CI/CD via GitHub Actions
```

AI will reference this during initialization.

---

**Remember**: Initialization is a one-time setup. The .agent/ directory will grow organically as the project evolves. Don't over-document upfront - let patterns emerge naturally.
