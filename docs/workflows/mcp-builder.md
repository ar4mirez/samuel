# MCP Builder

> **Purpose**: Guide for creating high-quality MCP (Model Context Protocol) servers that enable LLMs to interact with external services through well-designed tools. Supports both TypeScript (recommended) and Python implementations.

---

## When to Use

| Trigger | Description |
|---------|-------------|
| **MCP Server Development** | Building servers that expose tools to LLMs |
| **API Integration** | Integrating external APIs or services for AI agent access |
| **Tool Design** | Designing well-structured tools for LLM consumption |
| **Service Bridging** | Enabling AI agents to interact with external systems |

---

## Process

### Phase 1: Deep Research and Planning

1. **Understand modern MCP design** — balance API coverage with workflow tools; use clear naming and concise descriptions; return actionable error messages
2. **Study the MCP specification** — start with `https://modelcontextprotocol.io/sitemap.xml`
3. **Choose your stack** — TypeScript (recommended) or Python; Streamable HTTP for remote, stdio for local
4. **Plan your implementation** — review the service API, list endpoints to implement

### Phase 2: Implementation

1. **Set up project structure** using language-specific patterns
2. **Build core infrastructure** — API client, error handling, response formatting, pagination
3. **Implement tools** with:
    - Input schemas (Zod for TypeScript, Pydantic for Python)
    - Output schemas with `structuredContent`
    - Clear tool descriptions
    - Async/await for I/O operations
    - Annotations (`readOnlyHint`, `destructiveHint`, `idempotentHint`, `openWorldHint`)

### Phase 3: Review and Test

1. **Code quality review** — no duplication, consistent errors, full type coverage
2. **Build and test** — compile/verify, test with MCP Inspector (`npx @modelcontextprotocol/inspector`)

### Phase 4: Create Evaluations

1. **Create 10 evaluation questions** — complex, realistic, read-only, verifiable
2. **Verify answers** — solve each question yourself
3. **Output** as XML evaluation file

---

## Resources

The skill includes reference documentation:

| Resource | Purpose |
|----------|---------|
| `references/mcp_best_practices.md` | Core MCP guidelines |
| `references/node_mcp_server.md` | TypeScript patterns and examples |
| `references/python_mcp_server.md` | Python patterns and examples |
| `references/evaluation.md` | Evaluation creation guide |

---

## See Also

- [CLI Reference: skill](../reference/cli.md#skill) - Full command reference
- [The .claude Directory](../core/agent-directory.md) - Where skills live
- [MCP Specification](https://modelcontextprotocol.io/) - Official MCP documentation
