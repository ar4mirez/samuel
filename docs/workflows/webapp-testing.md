# Web App Testing

> **Purpose**: Test local web applications using Playwright. Supports verifying frontend functionality, debugging UI behavior, capturing browser screenshots, and viewing browser logs.

---

## When to Use

| Trigger | Description |
|---------|-------------|
| **Frontend Testing** | Verifying web application functionality |
| **UI Debugging** | Debugging visual or behavioral issues in the browser |
| **Browser Automation** | Automating browser interactions for testing |
| **Screenshot Capture** | Capturing browser state for inspection |

---

## Process

### Decision Tree

```text
User task --> Is it static HTML?
    |-- Yes --> Read HTML file to identify selectors
    |           |-- Success --> Write Playwright script
    |           +-- Fails --> Treat as dynamic (below)
    |
    +-- No (dynamic webapp) --> Is the server already running?
        |-- No --> Use scripts/with_server.py helper
        +-- Yes --> Reconnaissance-then-action:
            1. Navigate and wait for networkidle
            2. Take screenshot or inspect DOM
            3. Identify selectors from rendered state
            4. Execute actions with discovered selectors
```

### Step 1: Start the Server (if needed)

Use the `with_server.py` helper script:

```bash
# Single server
python scripts/with_server.py --server "npm run dev" --port 5173 -- python your_test.py

# Multiple servers (backend + frontend)
python scripts/with_server.py \
  --server "cd backend && python server.py" --port 3000 \
  --server "cd frontend && npm run dev" --port 5173 \
  -- python your_test.py
```

### Step 2: Write Playwright Script

```python
from playwright.sync_api import sync_playwright

with sync_playwright() as p:
    browser = p.chromium.launch(headless=True)
    page = browser.new_page()
    page.goto('http://localhost:5173')
    page.wait_for_load_state('networkidle')  # CRITICAL: wait for JS
    # ... your automation logic
    browser.close()
```

### Step 3: Reconnaissance-Then-Action

1. **Inspect rendered DOM** — take screenshots, read content, list elements
2. **Identify selectors** from inspection results
3. **Execute actions** using discovered selectors

---

## Best Practices

- **Use bundled scripts as black boxes** — run `--help` first, don't read source
- Use `sync_playwright()` for synchronous scripts
- Always close the browser when done
- Use descriptive selectors: `text=`, `role=`, CSS selectors, or IDs
- Add appropriate waits: `page.wait_for_selector()` or `page.wait_for_timeout()`

---

## Common Pitfall

Don't inspect the DOM before waiting for `networkidle` on dynamic apps. Always call `page.wait_for_load_state('networkidle')` before inspection.

---

## Resources

| Resource | Purpose |
|----------|---------|
| `scripts/with_server.py` | Server lifecycle management |
| `examples/element_discovery.py` | Discovering buttons, links, inputs |
| `examples/static_html_automation.py` | Local HTML with file:// URLs |
| `examples/console_logging.py` | Capturing console logs |

---

## See Also

- [CLI Reference: skill](../reference/cli.md#skill) - Full command reference
- [The .claude Directory](../core/agent-directory.md) - Where skills live
- [Web Artifacts Builder](web-artifacts-builder.md) - Build React applications
- [Testing Strategy](testing-strategy.md) - Test planning and coverage strategy
