# Web Artifacts Builder

> **Purpose**: Build elaborate, multi-component web applications using React, TypeScript, Tailwind CSS, and shadcn/ui. Bundles everything into a single self-contained HTML file. Use for complex projects requiring state management, routing, or shadcn/ui components.

---

## When to Use

| Trigger | Description |
|---------|-------------|
| **Complex Web Apps** | Multi-component applications with state management or routing |
| **React + shadcn/ui** | Projects requiring shadcn/ui components and Tailwind CSS |
| **Self-Contained Output** | Need a single HTML file that works in any browser |
| **Rapid Prototyping** | Quickly scaffold rich frontend applications |

---

## Process

### Step 1: Initialize Project

```bash
bash scripts/init-artifact.sh <project-name>
cd <project-name>
```

Creates a fully configured project with:

- React + TypeScript (via Vite)
- Tailwind CSS 3.4.1 with shadcn/ui theming
- Path aliases (`@/`) configured
- 40+ shadcn/ui components pre-installed
- Parcel configured for bundling

### Step 2: Develop Your Application

Edit the generated files to build your application. See the skill's reference documentation for common development tasks.

### Step 3: Bundle to Single HTML

```bash
bash scripts/bundle-artifact.sh
```

Creates `bundle.html` — a self-contained file with all JavaScript, CSS, and dependencies inlined. Opens in any browser with no setup required.

### Step 4: Share Output

Share the bundled HTML file with the user for viewing in any browser.

### Step 5: Test (Optional)

Test the output using available tools (Playwright, Puppeteer, or other skills). Avoid testing upfront — test after presenting the result if requested.

---

## Tech Stack

| Component | Technology |
|-----------|-----------|
| **Framework** | React 18 |
| **Language** | TypeScript |
| **Build Tool** | Vite + Parcel (bundling) |
| **Styling** | Tailwind CSS |
| **Components** | shadcn/ui (40+ pre-installed) |
| **Output** | Single HTML file |

---

## Design Guidelines

Avoid generic AI aesthetics:

- No excessive centered layouts
- No purple gradients on white backgrounds
- No uniform rounded corners everywhere
- No Inter font as default

---

## See Also

- [CLI Reference: skill](../reference/cli.md#skill) - Full command reference
- [The .claude Directory](../core/agent-directory.md) - Where skills live
- [Frontend Design](frontend-design.md) - Design-driven frontend development
- [Web App Testing](webapp-testing.md) - Test web applications with Playwright
- [shadcn/ui Components](https://ui.shadcn.com/docs/components) - Component reference
