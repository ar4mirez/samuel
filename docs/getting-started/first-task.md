---
title: Your First Task
description: Try Samuel with a real coding task
---

# Your First Task

Let's try the system with real examples to see how guardrails work in practice.

---

## Example 1: Bug Fix (ATOMIC Mode)

**Scenario**: A button is misaligned on the profile page.

### Your Prompt

```
Fix the profile edit button - it's misaligned on mobile screens
```

### What AI Does

1. **Deconstruct**: Single file change, CSS/styling issue
2. **Diagnose**: Check responsive breakpoints, flexbox/grid usage
3. **Develop**: Apply fix, verify against design
4. **Deliver**: Test, commit with conventional message

### Expected Behavior

AI should:

- [x] Locate the relevant component file
- [x] Identify the styling issue
- [x] Apply a minimal fix (no over-engineering)
- [x] Keep the file under 200 lines (component guardrail)
- [x] Suggest a commit: `fix(profile): align edit button on mobile screens`

---

## Example 2: Simple Feature (ATOMIC Mode)

**Scenario**: Add a "copy to clipboard" button.

### Your Prompt

```
Add a copy-to-clipboard button for the API key display in settings
```

### What AI Does

1. **Deconstruct**: New button + click handler + feedback
2. **Diagnose**: Check existing button patterns, clipboard API support
3. **Develop**: Implement with proper error handling
4. **Deliver**: Add test, validate, commit

### Expected Behavior

AI should:

- [x] Use existing button component (if available)
- [x] Handle clipboard API errors gracefully
- [x] Show user feedback (toast/tooltip)
- [x] Follow TypeScript strict mode (no `any`)
- [x] Add unit test for the click handler
- [x] Commit: `feat(settings): add copy-to-clipboard for API key`

---

## Example 3: Multi-File Feature (FEATURE Mode)

**Scenario**: Add user avatar upload.

### Your Prompt

```
Add avatar upload to user profile - allow jpg/png, max 2MB,
show preview before upload
```

### What AI Does

1. **Deconstruct**: Break into subtasks
   - File input component
   - Preview component
   - Validation logic
   - Upload API call
   - Error handling
2. **Diagnose**: Check existing upload patterns, API endpoints
3. **Develop**: Implement each subtask with tests
4. **Deliver**: Integration test, document, commit each step

### Expected Behavior

AI should:

- [x] Break into 3-5 subtasks
- [x] Implement sequentially (not all at once)
- [x] Validate file type and size (security guardrail)
- [x] Keep each file under 300 lines
- [x] Add tests for validation logic
- [x] Multiple commits, one per logical change

---

## Example 4: Complex Feature (COMPLEX Mode)

**Scenario**: Build a notification system.

### Your Prompt

```
@.claude/skills/create-prd/SKILL.md
Build a notification system with:
- In-app notifications
- Email notifications (optional)
- User preferences for notification types
```

### What AI Does

1. **Asks clarifying questions** (PRD workflow)
2. **Creates PRD** in `.claude/tasks/`
3. **Generates task breakdown** (20-50 subtasks)
4. **Implements step-by-step** with verification

### Expected Behavior

AI should:

- [x] Ask about notification types, delivery methods, storage
- [x] Create PRD document with all sections
- [x] Generate numbered task list
- [x] Wait for approval before implementing
- [x] Track progress in `.claude/state.md` (if multi-session)

---

## Verifying Guardrails Work

### Security Check

Try this prompt:

```
Create a function that queries the database for a user by email
```

AI should:

- [x] Use parameterized queries (not string concatenation)
- [x] Validate email input
- [x] Handle errors properly

If AI writes `WHERE email = '${email}'`, remind it: "Check security guardrails in CLAUDE.md"

### Code Quality Check

Try this prompt:

```
Write a function that processes order data with all business rules
```

AI should:

- [x] Keep function under 50 lines
- [x] Extract complex logic to helper functions
- [x] Add type signatures
- [x] No magic numbers

If the function exceeds 50 lines, AI should automatically refactor.

### Testing Check

Try this prompt:

```
Add a new utility function to format currency
```

AI should:

- [x] Create the function
- [x] Create test file alongside
- [x] Test edge cases (null, negative, large numbers)
- [x] Aim for >80% coverage

---

## Common Issues

### AI Doesn't Follow Guardrails

**Solution**: Explicitly remind it

```
Follow the guardrails in CLAUDE.md. Specifically:
- Functions must be ≤50 lines
- All inputs must be validated
- Add tests for the new code
```

### AI Skips Tests

**Solution**: Be explicit

```
After implementing this feature, add unit tests.
Target >80% coverage for the new code.
```

### AI Over-Engineers

**Solution**: Request simplicity

```
Keep this simple - minimal implementation that meets the requirements.
No extra features or "nice-to-haves".
```

### AI Makes Large Commits

**Solution**: Request atomic commits

```
Implement this in small steps. Commit after each logical change.
Use conventional commit format.
```

---

## Practice Prompts

Try these prompts to get comfortable with the system:

### Beginner

1. "Fix the typo in the README.md file"
2. "Add a loading spinner to the submit button"
3. "Update the copyright year in the footer"

### Intermediate

4. "Add form validation to the contact form (email required, message min 10 chars)"
5. "Create a dark mode toggle that persists to localStorage"
6. "Add pagination to the blog posts list"

### Advanced

7. "Refactor the API service to use a base HTTP client class"
8. "Add rate limiting to the public API endpoints"
9. "@.claude/skills/create-prd/SKILL.md - Build a search feature with filters"

---

## What You Learned

After completing these examples, you should understand:

- [x] **ATOMIC mode** for small, focused changes
- [x] **FEATURE mode** for multi-file features
- [x] **COMPLEX mode** with PRD workflow for large features
- [x] How guardrails are enforced automatically
- [x] How to remind AI when guardrails aren't followed

---

## Next Steps

<div class="grid cards" markdown>

-   :material-shield:{ .lg .middle } **Deep Dive into Guardrails**

    ---

    Understand all 35+ rules in detail.

    [:octicons-arrow-right-24: Guardrails](../core/guardrails.md)

-   :material-cog:{ .lg .middle } **Learn Workflows**

    ---

    Master PRD and task generation workflows.

    [:octicons-arrow-right-24: Workflows](../workflows/index.md)

-   :material-code-braces:{ .lg .middle } **Language Guide**

    ---

    Check your language-specific guide.

    [:octicons-arrow-right-24: Language Guides](../languages/index.md)

</div>
