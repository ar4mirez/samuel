---
title: Quick Start
description: Get started with AI Code Template in 60 seconds
---

# Quick Start Guide

Get up and running with AI Code Template in under a minute.

---

## 60-Second Setup

```bash
# 1. Copy to your project
cp -r /path/to/ai-code-template/{CLAUDE.md,.agent} ./

# 2. (Optional) For cross-tool compatibility
ln -s CLAUDE.md AGENTS.md

# 3. Start coding - AI automatically follows guardrails
# That's it! No configuration needed.
```

**The system works immediately:**

- [x] AI loads CLAUDE.md automatically (500 lines of guardrails + operations)
- [x] Language guides auto-load based on file extensions
- [x] Workflows available when you need them
- [x] Progressive - starts minimal, grows with your project
- [x] Cross-tool: Works with any AI assistant via AGENTS.md

---

## Choose Your Path

### :material-plus-circle: Path 1: New Project

```
@.agent/workflows/initialize-project.md

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

### :material-folder-open: Path 2: Existing Project

```
@.agent/workflows/initialize-project.md

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

### :material-lightning-bolt: Path 3: Jump Right In

Just start coding. AI follows guardrails automatically.

=== "Simple Tasks"

    ```
    "Fix the login button alignment"
    ```

    AI uses **ATOMIC mode** - single file, quick fix, tests, commit.

=== "Features"

    ```
    "Add user profile editing"
    ```

    AI uses **FEATURE mode** - breaks into subtasks, implements systematically.

=== "Complex Work"

    ```
    "Build real-time chat with WebSockets"
    ```

    AI suggests **COMPLEX mode** - offers PRD workflow for structured approach.

---

## What Happens Next?

Once you've copied the files, the system is ready. Here's what to expect:

### Automatic Loading

When you start working with an AI assistant:

1. **CLAUDE.md loads automatically** - 500 lines of guardrails and operations
2. **Language guides auto-load** - Based on file extensions you're working with
3. **Guardrails are enforced** - Code quality, security, testing standards

### Progressive Growth

The `.agent/` directory grows with your project:

| Timeline | What Gets Created |
|----------|-------------------|
| **Day 1** | Only CLAUDE.md + templates |
| **Week 1** | `.agent/project.md` (tech stack) |
| **Month 1** | `.agent/patterns.md` (conventions) |
| **Ongoing** | `.agent/memory/` (decisions) |

!!! tip "Don't Over-Document"

    Let the documentation grow naturally. Don't create `project.md` on day one - wait until you make architecture decisions.

---

## Verify It's Working

After setup, try a simple task:

```
"Add a console.log statement to debug the user login function"
```

AI should:

1. Locate the function
2. Add the debug statement
3. Mention this is temporary (guardrails discourage leaving debug code)
4. Follow conventional commit format if committing

If AI references guardrails or mentions code quality standards, **the system is working!**

---

## Next Steps

<div class="grid cards" markdown>

-   :material-book:{ .lg .middle } **Learn the Methodology**

    ---

    Understand the 4D approach: Deconstruct, Diagnose, Develop, Deliver.

    [:octicons-arrow-right-24: 4D Methodology](../core/methodology.md)

-   :material-shield:{ .lg .middle } **Review Guardrails**

    ---

    See all 35+ rules that AI will follow.

    [:octicons-arrow-right-24: All Guardrails](../core/guardrails.md)

-   :material-code-braces:{ .lg .middle } **Language Guides**

    ---

    Check the guide for your programming language.

    [:octicons-arrow-right-24: Language Guides](../languages/index.md)

-   :material-cog:{ .lg .middle } **Try a Workflow**

    ---

    Use PRD workflow for a complex feature.

    [:octicons-arrow-right-24: Workflows](../workflows/index.md)

</div>
