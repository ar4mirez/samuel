# HTML & CSS Guide

> **Applies to**: HTML5, CSS3, Sass/SCSS, Tailwind CSS, Responsive Design, Accessibility

---

## Core Principles

1. **Semantic HTML**: Use elements for their meaning, not appearance
2. **Accessibility First**: ARIA, keyboard navigation, screen readers
3. **Progressive Enhancement**: Base experience works without JS/CSS
4. **Mobile-First**: Design for mobile, enhance for larger screens
5. **Performance**: Minimize CSS, optimize loading

---

## HTML Guardrails

### Document Structure
- ✓ Use HTML5 doctype: `<!DOCTYPE html>`
- ✓ Include `lang` attribute: `<html lang="en">`
- ✓ Include meta viewport for responsive design
- ✓ Include meta charset: `<meta charset="UTF-8">`
- ✓ Use semantic elements (`header`, `nav`, `main`, `article`, `section`, `aside`, `footer`)

### Semantic HTML
- ✓ One `<main>` element per page
- ✓ Use `<article>` for self-contained content
- ✓ Use `<section>` for thematic groupings (with heading)
- ✓ Use `<nav>` for navigation links
- ✓ Use `<aside>` for tangentially related content
- ✓ Headings in order (`h1` → `h2` → `h3`, no skipping)

### Accessibility (a11y)
- ✓ All images have `alt` text (empty for decorative: `alt=""`)
- ✓ Form inputs have associated `<label>` elements
- ✓ Interactive elements are focusable and keyboard accessible
- ✓ Color is not the only means of conveying information
- ✓ Sufficient color contrast (WCAG AA: 4.5:1 for text)
- ✓ Use ARIA only when native HTML is insufficient

### Forms
- ✓ Every input has a label (visible or `aria-label`)
- ✓ Use appropriate input types (`email`, `tel`, `number`, etc.)
- ✓ Group related inputs with `<fieldset>` and `<legend>`
- ✓ Provide clear error messages
- ✓ Mark required fields

---

## CSS Guardrails

### Code Style
- ✓ Use lowercase for selectors and properties
- ✓ Use kebab-case for class names
- ✓ 2-space indentation
- ✓ One property per line
- ✓ Space after colon, before value
- ✓ Semicolon after every declaration
- ✓ Blank line between rule sets

### Naming Conventions (BEM)
- ✓ Block: `.card`
- ✓ Element: `.card__title`, `.card__body`
- ✓ Modifier: `.card--featured`, `.card__title--large`

### Selectors
- ✓ Avoid ID selectors for styling
- ✓ Keep specificity low
- ✓ Avoid deep nesting (max 3 levels)
- ✓ Don't use `!important` except for utilities
- ✓ Avoid tag selectors (except reset/base)

### Modern CSS
- ✓ Use CSS custom properties (variables)
- ✓ Use Flexbox and Grid for layout
- ✓ Use `clamp()` for fluid typography
- ✓ Use logical properties (`margin-inline`, `padding-block`)
- ✓ Use `@container` queries where supported

---

## HTML Templates

### Basic Document
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="description" content="Page description for SEO">
  <title>Page Title</title>
  <link rel="stylesheet" href="styles.css">
</head>
<body>
  <header>
    <nav aria-label="Main navigation">
      <a href="/" aria-current="page">Home</a>
      <a href="/about">About</a>
      <a href="/contact">Contact</a>
    </nav>
  </header>

  <main id="main-content">
    <h1>Page Heading</h1>
    <!-- Main content -->
  </main>

  <footer>
    <p>&copy; 2024 Company Name</p>
  </footer>

  <script src="script.js" defer></script>
</body>
</html>
```

### Semantic Article
```html
<article>
  <header>
    <h2>Article Title</h2>
    <p>
      <time datetime="2024-01-15">January 15, 2024</time>
      by <a href="/authors/john">John Doe</a>
    </p>
  </header>

  <p>Article introduction...</p>

  <section>
    <h3>Section Heading</h3>
    <p>Section content...</p>
  </section>

  <footer>
    <p>Tags:
      <a href="/tags/html" rel="tag">HTML</a>,
      <a href="/tags/css" rel="tag">CSS</a>
    </p>
  </footer>
</article>
```

### Accessible Form
```html
<form action="/submit" method="POST" novalidate>
  <fieldset>
    <legend>Contact Information</legend>

    <div class="form-group">
      <label for="name">Full Name <span aria-hidden="true">*</span></label>
      <input
        type="text"
        id="name"
        name="name"
        required
        autocomplete="name"
        aria-describedby="name-hint"
      >
      <small id="name-hint">Enter your first and last name</small>
    </div>

    <div class="form-group">
      <label for="email">Email <span aria-hidden="true">*</span></label>
      <input
        type="email"
        id="email"
        name="email"
        required
        autocomplete="email"
        aria-invalid="false"
      >
    </div>

    <div class="form-group">
      <label for="message">Message</label>
      <textarea
        id="message"
        name="message"
        rows="5"
        aria-describedby="message-hint"
      ></textarea>
      <small id="message-hint">Maximum 500 characters</small>
    </div>
  </fieldset>

  <button type="submit">Send Message</button>
</form>
```

### Navigation with Accessibility
```html
<nav aria-label="Main">
  <ul role="list">
    <li><a href="/" aria-current="page">Home</a></li>
    <li><a href="/products">Products</a></li>
    <li>
      <button
        aria-expanded="false"
        aria-controls="services-menu"
      >
        Services
      </button>
      <ul id="services-menu" hidden>
        <li><a href="/services/web">Web Development</a></li>
        <li><a href="/services/mobile">Mobile Apps</a></li>
      </ul>
    </li>
    <li><a href="/contact">Contact</a></li>
  </ul>
</nav>
```

---

## CSS Patterns

### CSS Reset/Normalize
```css
/* Modern CSS Reset */
*,
*::before,
*::after {
  box-sizing: border-box;
}

* {
  margin: 0;
  padding: 0;
}

html {
  -webkit-text-size-adjust: none;
  text-size-adjust: none;
}

body {
  min-height: 100vh;
  line-height: 1.5;
}

img,
picture,
video,
canvas,
svg {
  display: block;
  max-width: 100%;
}

input,
button,
textarea,
select {
  font: inherit;
}

p,
h1,
h2,
h3,
h4,
h5,
h6 {
  overflow-wrap: break-word;
}
```

### Custom Properties (Variables)
```css
:root {
  /* Colors */
  --color-primary: #3b82f6;
  --color-primary-dark: #2563eb;
  --color-secondary: #10b981;
  --color-text: #1f2937;
  --color-text-muted: #6b7280;
  --color-background: #ffffff;
  --color-surface: #f3f4f6;
  --color-border: #e5e7eb;
  --color-error: #ef4444;
  --color-success: #22c55e;

  /* Typography */
  --font-sans: system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
  --font-mono: 'Fira Code', Consolas, monospace;

  --text-xs: clamp(0.75rem, 0.7rem + 0.25vw, 0.875rem);
  --text-sm: clamp(0.875rem, 0.8rem + 0.35vw, 1rem);
  --text-base: clamp(1rem, 0.9rem + 0.5vw, 1.125rem);
  --text-lg: clamp(1.125rem, 1rem + 0.6vw, 1.25rem);
  --text-xl: clamp(1.25rem, 1rem + 1.25vw, 1.5rem);
  --text-2xl: clamp(1.5rem, 1rem + 2.5vw, 2rem);

  /* Spacing */
  --space-xs: 0.25rem;
  --space-sm: 0.5rem;
  --space-md: 1rem;
  --space-lg: 1.5rem;
  --space-xl: 2rem;
  --space-2xl: 3rem;

  /* Border Radius */
  --radius-sm: 0.25rem;
  --radius-md: 0.5rem;
  --radius-lg: 1rem;
  --radius-full: 9999px;

  /* Shadows */
  --shadow-sm: 0 1px 2px rgb(0 0 0 / 0.05);
  --shadow-md: 0 4px 6px rgb(0 0 0 / 0.1);
  --shadow-lg: 0 10px 15px rgb(0 0 0 / 0.1);

  /* Transitions */
  --transition-fast: 150ms ease;
  --transition-normal: 250ms ease;
}

/* Dark mode */
@media (prefers-color-scheme: dark) {
  :root {
    --color-text: #f9fafb;
    --color-text-muted: #9ca3af;
    --color-background: #111827;
    --color-surface: #1f2937;
    --color-border: #374151;
  }
}
```

### Layout with CSS Grid
```css
/* Holy Grail Layout */
.layout {
  display: grid;
  grid-template-areas:
    "header header header"
    "nav    main   aside"
    "footer footer footer";
  grid-template-columns: 200px 1fr 200px;
  grid-template-rows: auto 1fr auto;
  min-height: 100vh;
}

.layout > header { grid-area: header; }
.layout > nav { grid-area: nav; }
.layout > main { grid-area: main; }
.layout > aside { grid-area: aside; }
.layout > footer { grid-area: footer; }

@media (max-width: 768px) {
  .layout {
    grid-template-areas:
      "header"
      "nav"
      "main"
      "aside"
      "footer";
    grid-template-columns: 1fr;
  }
}

/* Card Grid */
.card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(min(300px, 100%), 1fr));
  gap: var(--space-lg);
}
```

### Flexbox Patterns
```css
/* Centered content */
.center {
  display: flex;
  justify-content: center;
  align-items: center;
}

/* Space between with wrap */
.flex-wrap {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: var(--space-md);
}

/* Sticky footer */
.page {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.page > main {
  flex: 1;
}

/* Navigation */
.nav {
  display: flex;
  align-items: center;
  gap: var(--space-md);
}

.nav__links {
  display: flex;
  gap: var(--space-sm);
  margin-inline-start: auto;
}
```

### BEM Component
```css
/* Card component using BEM */
.card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.card__image {
  aspect-ratio: 16 / 9;
  object-fit: cover;
  width: 100%;
}

.card__body {
  padding: var(--space-md);
}

.card__title {
  font-size: var(--text-lg);
  font-weight: 600;
  margin-block-end: var(--space-xs);
}

.card__description {
  color: var(--color-text-muted);
  font-size: var(--text-sm);
}

.card__footer {
  border-top: 1px solid var(--color-border);
  display: flex;
  gap: var(--space-sm);
  justify-content: flex-end;
  padding: var(--space-sm) var(--space-md);
}

/* Modifiers */
.card--featured {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-md);
}

.card--horizontal {
  display: grid;
  grid-template-columns: 200px 1fr;
}

.card--horizontal .card__image {
  aspect-ratio: 1;
  height: 100%;
}
```

### Responsive Typography
```css
/* Fluid typography using clamp() */
h1 {
  font-size: clamp(2rem, 1.5rem + 2.5vw, 3.5rem);
  line-height: 1.1;
}

h2 {
  font-size: clamp(1.5rem, 1.25rem + 1.25vw, 2.25rem);
  line-height: 1.2;
}

h3 {
  font-size: clamp(1.25rem, 1rem + 1.25vw, 1.75rem);
  line-height: 1.3;
}

p {
  font-size: var(--text-base);
  max-width: 65ch; /* Optimal line length */
}
```

---

## Accessibility Patterns

### Skip Link
```html
<a href="#main-content" class="skip-link">Skip to main content</a>
```

```css
.skip-link {
  background: var(--color-primary);
  color: white;
  left: 50%;
  padding: var(--space-sm) var(--space-md);
  position: absolute;
  transform: translate(-50%, -100%);
  transition: transform var(--transition-fast);
  z-index: 100;
}

.skip-link:focus {
  transform: translate(-50%, 0);
}
```

### Focus Styles
```css
/* Remove default outline, add custom focus */
:focus {
  outline: none;
}

:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 2px;
}

/* High contrast focus for buttons */
button:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 2px;
  box-shadow: 0 0 0 4px rgb(59 130 246 / 0.3);
}
```

### Screen Reader Only
```css
.sr-only {
  border: 0;
  clip: rect(0, 0, 0, 0);
  height: 1px;
  margin: -1px;
  overflow: hidden;
  padding: 0;
  position: absolute;
  white-space: nowrap;
  width: 1px;
}

/* Make visible when focused */
.sr-only-focusable:focus {
  clip: auto;
  height: auto;
  margin: 0;
  overflow: visible;
  position: static;
  white-space: normal;
  width: auto;
}
```

### Reduced Motion
```css
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
    scroll-behavior: auto !important;
  }
}
```

---

## SCSS/Sass Patterns

### Variables and Mixins
```scss
// _variables.scss
$breakpoints: (
  sm: 640px,
  md: 768px,
  lg: 1024px,
  xl: 1280px,
);

$colors: (
  primary: #3b82f6,
  secondary: #10b981,
  danger: #ef4444,
);

// _mixins.scss
@mixin respond-to($breakpoint) {
  $value: map-get($breakpoints, $breakpoint);
  @if $value {
    @media (min-width: $value) {
      @content;
    }
  }
}

@mixin flex-center {
  display: flex;
  justify-content: center;
  align-items: center;
}

@mixin truncate($lines: 1) {
  @if $lines == 1 {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  } @else {
    display: -webkit-box;
    -webkit-line-clamp: $lines;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
}

// Usage
.card {
  @include respond-to(md) {
    display: grid;
    grid-template-columns: 200px 1fr;
  }
}

.title {
  @include truncate(2);
}
```

---

## Common Pitfalls

### HTML
```html
<!-- Bad: Div soup -->
<div class="header">
  <div class="nav">...</div>
</div>

<!-- Good: Semantic -->
<header>
  <nav>...</nav>
</header>

<!-- Bad: Missing alt -->
<img src="photo.jpg">

<!-- Good: Descriptive alt -->
<img src="photo.jpg" alt="Team meeting in conference room">

<!-- Bad: Using outline: none -->
<button style="outline: none;">Click</button>

<!-- Good: Visible focus state -->
<button class="btn">Click</button>
```

### CSS
```css
/* Bad: ID selectors */
#header { }

/* Good: Class selectors */
.header { }

/* Bad: Deep nesting */
.nav ul li a span { }

/* Good: Flat selectors */
.nav__link-text { }

/* Bad: Hardcoded colors */
color: #3b82f6;

/* Good: CSS variables */
color: var(--color-primary);

/* Bad: Fixed widths */
width: 500px;

/* Good: Flexible widths */
max-width: 500px;
width: 100%;
```

---

## Tools

### Linting & Formatting
- **HTMLHint**: HTML linting
- **Stylelint**: CSS linting
- **Prettier**: Code formatting

### Testing
- **Lighthouse**: Performance & accessibility
- **axe DevTools**: Accessibility testing
- **Wave**: Accessibility evaluation

### Configuration (stylelint)
```json
{
  "extends": ["stylelint-config-standard"],
  "rules": {
    "selector-class-pattern": "^[a-z][a-z0-9]*(-[a-z0-9]+)*(__[a-z0-9]+(-[a-z0-9]+)*)?(--[a-z0-9]+(-[a-z0-9]+)*)?$",
    "declaration-block-no-redundant-longhand-properties": true,
    "no-descending-specificity": null
  }
}
```

---

## References

- [MDN Web Docs](https://developer.mozilla.org/)
- [HTML Living Standard](https://html.spec.whatwg.org/)
- [CSS Specifications](https://www.w3.org/Style/CSS/)
- [WCAG Guidelines](https://www.w3.org/WAI/WCAG21/quickref/)
- [A11y Project](https://www.a11yproject.com/)
- [CSS-Tricks](https://css-tricks.com/)
- [Every Layout](https://every-layout.dev/)
- [BEM Methodology](https://getbem.com/)
