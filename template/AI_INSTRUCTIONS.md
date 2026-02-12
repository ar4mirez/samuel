# AICoF - Quick Start Guide

> **Purpose**: Get started with the AICoF (Artificial Intelligence Coding Framework) in 5 minutes
>
> **Status**: Production ready, use immediately
> **Version**: 1.8.0
> **Cross-Tool**: AGENTS.md compatible (works with Cursor, Codex, Copilot, etc.)

---

## 🚀 60-Second Quick Start

```bash
# 1. Copy to your project
cp -r /path/to/aicof/{CLAUDE.md,.agent} ./

# 2. (Optional) For cross-tool compatibility
ln -s CLAUDE.md AGENTS.md

# 3. Start coding - AI automatically follows guardrails
# That's it! No configuration needed.
```

**The system works immediately:**
- ✅ AI loads CLAUDE.md automatically (500 lines of guardrails + operations)
- ✅ Language guides auto-load based on file extensions
- ✅ Workflows available when you need them
- ✅ Progressive - starts minimal, grows with your project
- ✅ Cross-tool: Works with any AI assistant via AGENTS.md

---

## 📖 What Is AICoF?

AICoF (Artificial Intelligence Coding Framework) is an **opinionated AI development framework** with:
- **35+ specific guardrails** (not vague suggestions)
- **21 language guide skills** (TypeScript, Python, Go, Rust, Kotlin, Java, C#, PHP, Swift, C/C++, Ruby, SQL, Shell, R, Dart, HTML/CSS, Lua, Assembly, CUDA, Solidity, Zig)
- **15 workflow skills** (PRD, tasks, init, troubleshooting, AGENTS.md generator, code-review, security-audit, testing-strategy, cleanup-project, refactoring, dependency-update, document-work, update-framework, create-rfd, create-skill)
- **33 framework skills** (React, Next.js, Express, Django, FastAPI, Flask, Gin, and more)
- **3 modes** (ATOMIC/FEATURE/COMPLEX - scales from bugs to architecture)
- **Cross-tool compatible** (AGENTS.md standard)

**Philosophy**: Small, validated changes. Quality enforced. Documentation grows organically.

---

## 🎯 Choose Your Path

### Path 1: New Project (10 minutes)

```
@.agent/skills/initialize-project/SKILL.md

"Initialize a new [TypeScript/Python/Go/Rust] project with [describe your project]"
```

**AI will ask:**
1. Tech stack? (language, framework, versions)
2. Architecture? (monolith, microservices, serverless)
3. Testing approach? (unit, integration, e2e)
4. Deployment target? (AWS, Vercel, etc.)

**AI will create:**
- `.agent/project.md` (your tech stack documented)
- Directory structure (best practices for your stack)
- Config files (tsconfig, package.json, etc.)
- `.gitignore`, `.env.example`, `README.md`

**Then**: Start coding! AI follows all guardrails automatically.

---

### Path 2: Existing Project (5 minutes)

```
@.agent/skills/initialize-project/SKILL.md

"This is an existing project - analyze the codebase"
```

**AI will:**
1. Scan for tech stack (package.json, requirements.txt, etc.)
2. Examine directory structure
3. Analyze code patterns
4. Review recent commits
5. Create `.agent/project.md` with findings
6. Extract patterns to `.agent/patterns.md`
7. Identify gaps and suggest improvements

**You confirm or correct**, then start coding with guardrails.

---

### Path 3: Jump Right In (0 setup)

Just start coding. AI follows guardrails automatically.

**For simple tasks:**
```
"Fix the login button alignment"
```
AI uses ATOMIC mode - single file, quick fix, tests, commit.

**For features:**
```
"Add user profile editing"
```
AI uses FEATURE mode - breaks into subtasks, implements systematically.

**For complex work:**
```
"Build real-time chat with WebSockets"
```
AI suggests COMPLEX mode - offers PRD workflow for structured approach.

---

## 🎨 How It Works

### The 3 Modes

**ATOMIC** (<5 files, clear scope)
- Direct implementation
- Quick validation
- One commit
- Example: Bug fixes, styling, simple features

**FEATURE** (5-10 files)
- Break into 3-5 subtasks
- Implement sequentially
- Integration testing
- Example: New component, API endpoint, refactoring

**COMPLEX** (>10 files, new subsystem)
- Optional: Create PRD (Product Requirements Document)
- Generate task breakdown
- Step-by-step implementation
- Example: Authentication, payments, analytics

**AI auto-detects which mode to use** based on your request.

---

### The Guardrails (35+ rules)

**Code Quality:**
- ✓ Functions ≤50 lines
- ✓ Files ≤300 lines
- ✓ Complexity ≤10 per function
- ✓ All exports have types/docs

**Security (CRITICAL):**
- ✓ All inputs validated
- ✓ Parameterized queries only
- ✓ No secrets in code
- ✓ Dependencies checked for vulnerabilities + licenses

**Testing (CRITICAL):**
- ✓ >80% coverage for business logic
- ✓ >60% overall coverage
- ✓ Tests for all public APIs
- ✓ Regression tests for bugs

**Git:**
- ✓ Conventional commits (`feat:`, `fix:`, etc.)
- ✓ One logical change per commit
- ✓ All tests pass before push
- ✓ PRs required (no direct commits to main)

**Performance:**
- ✓ No N+1 queries
- ✓ Pagination for large datasets
- ✓ API responses <200ms
- ✓ Frontend bundles <200KB

**Full list**: See [CLAUDE.md](CLAUDE.md) lines 33-80

---

### Language Guides (Auto-Load)

AI automatically loads the right guide based on file extensions:

**TypeScript/JavaScript** (`.ts`, `.tsx`, `.js`, `.jsx`)
- Strict mode, Zod validation, React patterns, async/await
- [@.agent/skills/typescript-guide/SKILL.md](.agent/skills/typescript-guide/SKILL.md)

**Python** (`.py`)
- Type hints, Pydantic, Django/FastAPI patterns, pytest
- [@.agent/skills/python-guide/SKILL.md](.agent/skills/python-guide/SKILL.md)

**Go** (`.go`)
- Error handling, concurrency, interfaces, table tests
- [@.agent/skills/go-guide/SKILL.md](.agent/skills/go-guide/SKILL.md)

**Rust** (`.rs`)
- Ownership, Result<T,E>, async, zero-cost abstractions
- [@.agent/skills/rust-guide/SKILL.md](.agent/skills/rust-guide/SKILL.md)

**Kotlin** (`.kt`, `.kts`)
- Null safety, coroutines, data classes, extension functions
- [@.agent/skills/kotlin-guide/SKILL.md](.agent/skills/kotlin-guide/SKILL.md)

**No manual selection needed** - it just works!

---

## 🛠️ Common Workflows

### Build a Complex Feature

```
Step 1: Create PRD
@.agent/skills/create-prd/SKILL.md
"Build user authentication with OAuth"

Step 2: Generate Tasks
@.agent/skills/generate-tasks/SKILL.md
(AI uses the PRD you just created)

Step 3: Implement
"Start on task 1.1"
AI implements → Tests → Validates guardrails → Commits
"Continue with task 1.2"
Repeat until done ✅
```

### Fix a Bug

```
"The user profile page crashes when email is null"

AI will:
1. Find the issue
2. Fix it (with null check)
3. Add regression test
4. Validate guardrails
5. Commit: "fix(profile): handle null email gracefully"
```

### Refactor Code

```
"Refactor the API service to use dependency injection"

AI will:
1. Analyze current structure
2. Break into subtasks
3. Implement incrementally (keeps tests passing)
4. Update documentation
5. Conventional commits for each step
```

### I'm Stuck!

```
@.agent/skills/troubleshooting/SKILL.md

AI will:
1. Help you debug systematically
2. Check fundamentals (deps, config, versions)
3. Provide language-specific troubleshooting
4. Document solution in .agent/memory/
```

---

## 📚 File Structure Overview

```
your-project/
├── CLAUDE.md                    # Core guardrails (400 lines)
├── AI_INSTRUCTIONS.md           # This file (quick start)
└── .agent/                      # Project context (grows over time)
    ├── README.md                # How to use .agent/
    ├── project.md               # Your tech stack (created when chosen)
    ├── patterns.md              # Your code patterns (created when emerge)
    ├── state.md                 # Current work (optional, multi-session)
    │
    ├── skills/                  # Agent Skills (auto-load based on file type)
    │   ├── go-guide/            # Go language guide skill
    │   │   ├── SKILL.md         # Core guardrails and patterns
    │   │   └── references/      # Detailed patterns, pitfalls
    │   ├── typescript-guide/    # TypeScript language guide skill
    │   ├── python-guide/        # Python language guide skill
    │   ├── ...                  # 18 more language guide skills
    │   ├── create-prd/          # Workflow skills (on-demand)
    │   │   └── SKILL.md         # Product Requirements Doc
    │   ├── generate-tasks/
    │   │   └── SKILL.md         # Task breakdown
    │   ├── initialize-project/
    │   │   └── SKILL.md         # New/existing project setup
    │   ├── troubleshooting/
    │   │   └── SKILL.md         # Debug guide
    │   ├── generate-agents-md/
    │   │   └── SKILL.md         # Cross-tool compatibility
    │   └── ... (8 more)         # code-review, security-audit, testing-strategy, cleanup-project, refactoring, dependency-update, document-work, create-skill
    │
    ├── tasks/                   # PRDs and task lists (complex features)
    │   └── (created on demand)
    │
    ├── memory/                  # Decision logs (key learnings)
    │   └── (created on demand)
    │
    ├── project.md.template      # Template for project.md
    └── state.md.template        # Template for state.md
```

**Day 1**: Only CLAUDE.md + .agent/ templates
**Week 1**: project.md created
**Month 1**: patterns.md populated
**Ongoing**: memory/ captures decisions

---

## 💡 Pro Tips

### Tip 1: Trust the Auto-Loading
You don't need to manually load language guides. Just start coding:
```
"Create a new Express API endpoint for users"
```
AI detects TypeScript/Express → Loads typescript.md → Follows best practices automatically.

### Tip 2: Use Quick Reference
CLAUDE.md has a Quick Reference section (lines 7-30) with:
- Task classification
- Common guardrails
- Emergency links
- Line numbers for quick navigation

### Tip 3: Workflows Are Optional
Don't use PRD workflow for simple features. It's for:
- ❌ Bug fixes (use ATOMIC)
- ❌ Simple features (use FEATURE)
- ✅ New subsystems (use COMPLEX + PRD)
- ✅ Unclear requirements (use PRD to clarify)

### Tip 4: Let .agent/ Grow Naturally
Don't create project.md on Day 1. Let AI create it when:
- You make your first architecture decision
- You choose a tech stack
- You establish a pattern

**Don't over-document upfront** - it's anti-pattern.

### Tip 5: Commit Often
The system enforces atomic commits:
```
feat(auth): add login endpoint
fix(profile): handle null email
refactor(api): extract validation to middleware
test(auth): add OAuth flow tests
```

One logical change per commit. Tests must pass. Guardrails must validate.

---

## 🎓 Learning Path

### Week 1: Learn the Basics
- [ ] Initialize your first project
- [ ] Write 5 features using ATOMIC mode
- [ ] Review guardrails in CLAUDE.md
- [ ] Check which language guide applies to you

### Week 2: Try Complex Features
- [ ] Use PRD workflow for a medium feature
- [ ] Generate task breakdown
- [ ] Implement step-by-step
- [ ] Notice how .agent/project.md grows

### Week 3: Customize
- [ ] Add project-specific patterns to .agent/patterns.md
- [ ] Create first decision log in .agent/memory/
- [ ] Review and refine .agent/project.md
- [ ] Experiment with different modes

### Month 1: Master It
- [ ] Used on 10+ features
- [ ] System adapted to your workflow
- [ ] .agent/ directory reflects your project
- [ ] Would recommend to others

---

## ❓ FAQ

### Q: Do I need to read all 400 lines of CLAUDE.md?
**A:** No. The Quick Reference (lines 7-30) covers 80% of daily use. Read the rest when you need specific guidance.

### Q: Which language guide should I use?
**A:** AI auto-loads based on file extensions. You don't choose - it just works.

### Q: What if my language isn't covered?
**A:** The core guardrails (CLAUDE.md) are language-agnostic. 90% still applies. Consider contributing a new language guide!

### Q: Is PRD workflow required?
**A:** No! Only for complex features (>10 files, new subsystems). Most work is ATOMIC or FEATURE mode.

### Q: Can I customize the guardrails?
**A:** Yes. Edit CLAUDE.md for your team. Common customizations:
- File length limits (change 300 → 500 if needed)
- Coverage targets (change 80% → 90% if stricter)
- Commit format (add your own types)

### Q: How do I add a new language?
**A:** Create a new skill directory at `.agent/skills/<lang>-guide/` with a `SKILL.md` file. See `.agent/skills/README.md` for the format and existing guides for examples.

### Q: What if AI doesn't follow guardrails?
**A:** Explicitly remind: "Validate all guardrails from CLAUDE.md before committing." AI will check each ✓.

### Q: Can I use this with [Cursor/GitHub Copilot/other AI]?
**A:** Yes! Create a symlink: `ln -s CLAUDE.md AGENTS.md`. This follows the AGENTS.md standard that 20+ AI tools support. Or use the workflow: `@.agent/skills/generate-agents-md/SKILL.md`

---

## 🚨 Emergency Procedures

### Tests Failing After Change?
```
@.agent/skills/troubleshooting/SKILL.md
"Tests broke after my last commit"
```

### Stuck on Same Issue >30 Minutes?
```
@.agent/skills/troubleshooting/SKILL.md
"Been stuck on [issue] for 45 minutes, tried [what you tried]"
```

### Security Issue Found?
```
"SECURITY: Found SQL injection in user.service.ts line 42"
```
AI will:
1. Stop everything (security is CRITICAL)
2. Fix immediately
3. Add regression test
4. Review similar code for same issue
5. Document in .agent/memory/

### Build Broken?
```
@.agent/skills/troubleshooting/SKILL.md
"Build failing with [error message]"
```

AI will check: dependencies, versions, config, environment variables, common issues.

---

## 📞 Getting Help

### Built-in Help
- **Quick Reference**: CLAUDE.md lines 7-30
- **Troubleshooting**: @.agent/skills/troubleshooting/SKILL.md
- **Initialization**: @.agent/skills/initialize-project/SKILL.md
- **Language-specific**: @.agent/skills/[your-language]-guide/SKILL.md

### Documentation
- **Core guardrails**: [CLAUDE.md](CLAUDE.md)
- **System overview**: [.agent/README.md](.agent/README.md)
- **Skills**: [.agent/skills/README.md](.agent/skills/README.md)

### Ask AI
```
"How do I use the PRD workflow?"
"What guardrails apply to API endpoints?"
"Show me the TypeScript testing patterns"
```

AI has full context of the system and can explain anything.

---

## 🎯 Success Criteria

**You're successful when:**
- ✅ AI automatically follows guardrails (you don't remind it)
- ✅ Language guides load automatically (you don't think about it)
- ✅ Code quality is consistently high (fewer bugs)
- ✅ .agent/ directory reflects your project naturally
- ✅ You spend more time building, less time debugging

**Measure:**
- Week 1: Used on 1+ tasks
- Month 1: .agent/ has 3+ files (project.md, patterns.md, etc.)
- Quarter 1: Would recommend to team/friends

---

## 🎉 You're Ready!

**The system is complete and production-ready.**

**Next steps:**
1. Copy CLAUDE.md + .agent/ to your project
2. Start with initialization (new project) or jump right in (existing)
3. Let AI guide you with guardrails
4. Watch .agent/ grow with your project

**Remember:**
- Small, validated changes
- Trust the auto-loading
- Don't over-document upfront
- Let patterns emerge naturally

**Happy coding! 🚀**

---

## 📋 Quick Command Reference

```bash
# Initialize new project
@.agent/skills/initialize-project/SKILL.md

# Create PRD for complex feature
@.agent/skills/create-prd/SKILL.md

# Generate tasks from PRD
@.agent/skills/generate-tasks/SKILL.md

# Get unstuck / debug
@.agent/skills/troubleshooting/SKILL.md

# Generate AGENTS.md for cross-tool compatibility
@.agent/skills/generate-agents-md/SKILL.md

# Check guardrails
@CLAUDE.md lines 33-80

# Load language guide manually (usually auto-loads)
@.agent/skills/typescript-guide/SKILL.md
@.agent/skills/python-guide/SKILL.md
@.agent/skills/go-guide/SKILL.md
@.agent/skills/rust-guide/SKILL.md
@.agent/skills/kotlin-guide/SKILL.md

# Check current state
cat .agent/project.md
cat .agent/patterns.md
cat .agent/state.md

# Cross-tool setup
ln -s CLAUDE.md AGENTS.md
```

---

*Version: 1.8.0*
*Last Updated: 2025-01-14*
*Status: Production Ready*
*Cross-Tool: AGENTS.md Compatible*

**For full details, see [CLAUDE.md](CLAUDE.md)**
