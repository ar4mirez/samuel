# Autonomous Iteration Prompt

You are running in autonomous mode as part of the Ralph Wiggum methodology.
Each iteration is independent -- you start with a fresh context window.

## Your Task

1. **Read project context**:
   - Read `CLAUDE.md` or `AGENTS.md` for project guardrails
   - Read `.claude/auto/progress.txt` for learnings from prior iterations
   - Read `.claude/auto/prd.json` to find the task list and current state

2. **Select the next task**:
   - Find the highest-priority task with status "pending"
   - Respect dependencies: skip tasks whose `depends_on` tasks are not yet "completed" or "skipped"
   - Prefer tasks with priority "critical" > "high" > "medium" > "low"
   - If priorities are equal, prefer lower-numbered task IDs

3. **Implement the task**:
   - Update the task's status to "in_progress" in prd.json
   - Follow project guardrails from CLAUDE.md
   - Write tests alongside code
   - Keep changes atomic -- one task per iteration

4. **Run quality checks**:
   - Execute the commands listed in `prd.json` under `config.quality_checks`
   - All checks must pass before committing
   - If a check fails, fix the issue and retry

5. **Commit changes**:
   - Use conventional commit format: `type(scope): description`
   - Include task ID in commit message
   - Example: `feat(auth): task 1.1 - create user schema`

6. **Update state**:
   - Set the task's status to "completed" in prd.json
   - Record the commit SHA in the task's `commit_sha` field
   - Update `progress.total_tasks` and `progress.completed_tasks`

7. **Document learnings**:
   - Append any insights, gotchas, or decisions to `.claude/auto/progress.txt`
   - Format: `[timestamp] [iteration:N] [task:ID] LEARNING: description`

## Rules

- Complete exactly ONE task per iteration
- Never skip quality checks
- If stuck for too long, mark the task as "blocked" and document why
- Keep functions ≤50 lines, files ≤300 lines (project guardrails)
- All exported functions need documentation
- Write tests for all new code

## Error Recovery

If you encounter errors:
1. Try to fix them within this iteration
2. If unfixable, mark the task as "blocked" with a description
3. Append the error details to progress.txt as a LEARNING entry
4. The next iteration will have fresh context and can try a different approach
