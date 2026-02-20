# Doc Co-Authoring

> **Purpose**: Guide users through a structured workflow for co-authoring documentation. Helps efficiently transfer context, refine content through iteration, and verify the document works for readers.

---

## When to Use

| Trigger | Description |
|---------|-------------|
| **Writing Documentation** | Drafting proposals, technical specs, or decision docs |
| **Structured Content** | Creating PRDs, design docs, RFCs, or similar structured documents |
| **Collaborative Drafting** | Co-authoring content with AI through iterative refinement |
| **Document Quality** | Ensuring docs work for readers through structured testing |

---

## Process

### Stage 1: Context Gathering

Close the gap between what the user knows and what the AI knows:

1. **Initial questions** — document type, audience, desired impact, template, constraints
2. **Info dumping** — encourage the user to dump all context (background, discussions, architecture, constraints)
3. **Clarifying questions** — generate 5-10 targeted questions based on gaps
4. **Exit condition** — sufficient context gathered when edge cases and trade-offs can be discussed

### Stage 2: Refinement & Structure

Build the document section by section:

1. **Section ordering** — start with whichever section has the most unknowns
2. **For each section**:
    - Ask 5-10 clarifying questions
    - Brainstorm 5-20 options for content
    - User curates (keep/remove/combine)
    - Gap check for missing items
    - Draft the section
    - Iterative refinement through surgical edits
3. **Near completion** — re-read entire document checking for flow, redundancy, and filler
4. **Quality checking** — after 3 iterations with no changes, ask if anything can be removed

### Stage 3: Reader Testing

Test the document with a fresh perspective:

1. **Predict reader questions** — generate 5-10 realistic questions readers would ask
2. **Test with sub-agent** — use a fresh AI instance (no context) to answer the questions
3. **Run additional checks** — check for ambiguity, false assumptions, contradictions
4. **Report and fix** — loop back to refinement for any problematic sections
5. **Exit condition** — reader consistently answers correctly with no new gaps

---

## Key Principles

- **Quality over speed** — each iteration should make meaningful improvements
- **User agency** — always give the user control to adjust the process
- **Context management** — don't let gaps accumulate, address them as they come up
- **Direct and procedural** — explain rationale briefly, don't oversell the approach

---

## See Also

- [CLI Reference: skill](../reference/cli.md#skill) - Full command reference
- [The .claude Directory](../core/agent-directory.md) - Where skills live
- [Document Work](document-work.md) - Capture patterns and decisions after work
