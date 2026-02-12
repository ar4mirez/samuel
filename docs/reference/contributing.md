---
title: Contributing
description: How to contribute to AICoF
---

# Contributing

Thank you for your interest in contributing to AICoF!

---

## Ways to Contribute

### 1. New Language Guides

Add support for languages not yet covered:

**Needed**:

- Java
- C#
- PHP
- Swift
- Ruby
- Scala

**How to contribute**:

1. Copy an existing guide from `.agent/skills/<language>-guide/`
2. Create a new skill directory for your language (e.g., `java-guide/SKILL.md`)
3. Adapt all sections for your language
4. Test with real projects
5. Submit PR

**Guide structure**:

```markdown
# [Language] Guide

## Core Principles
## Language-Specific Guardrails
## Validation & Input Handling
## Testing
## Tooling
## Common Pitfalls
## Framework-Specific Patterns
## Performance Considerations
## Security Best Practices
## References
```

---

### 2. Framework Templates

Create templates for specific frameworks:

**Needed**:

- Next.js
- Django REST
- Ruby on Rails
- Spring Boot
- Laravel

**How to contribute**:

1. Create `.agent/templates/[framework]/`
2. Include framework-specific patterns
3. Include common configurations
4. Document usage

---

### 3. Documentation Improvements

- Fix typos and errors
- Clarify confusing sections
- Add more examples
- Improve organization

---

### 4. Bug Reports

Found something wrong?

1. Check existing issues first
2. Open a new issue with:
   - Clear title
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details

---

### 5. Feature Suggestions

Have an idea?

1. Open a discussion first
2. Describe the use case
3. Explain the benefit
4. Discuss implementation approach

---

## Development Setup

### Clone the Repository

```bash
git clone https://github.com/ar4mirez/aicof.git
cd aicof
```

### Install Documentation Dependencies

```bash
pip install -r requirements-docs.txt
```

### Run Documentation Locally

```bash
mkdocs serve
```

Open http://localhost:8000 to preview.

---

## Contribution Guidelines

### Code Quality

Follow the guardrails in CLAUDE.md:

- Functions ≤50 lines
- Files ≤300 lines
- Clear, descriptive names
- No dead code

### Commit Messages

Use conventional commits:

```
feat(language-guide): add Java guide
fix(docs): correct TypeScript example
docs(readme): update installation instructions
```

### Pull Requests

1. Fork the repository
2. Create a feature branch: `feat/java-guide`
3. Make your changes
4. Test thoroughly
5. Submit PR with clear description

### PR Description Template

```markdown
## Summary
Brief description of changes

## Type of Change
- [ ] New language guide
- [ ] Framework template
- [ ] Documentation improvement
- [ ] Bug fix
- [ ] Other

## Testing
How did you test this?

## Checklist
- [ ] Follows CLAUDE.md guardrails
- [ ] Documentation updated
- [ ] Tested with real project
```

---

## Language Guide Standards

When creating a language guide:

### Required Sections

1. **Core Principles** - 5 key principles for the language
2. **Language-Specific Guardrails** - Rules beyond universal ones
3. **Validation** - Input validation patterns
4. **Testing** - Frameworks and patterns
5. **Tooling** - Linters, formatters, configs
6. **Common Pitfalls** - Do/Don't examples
7. **Framework Patterns** - Major framework examples
8. **Performance** - Optimization guidelines
9. **Security** - Security best practices
10. **References** - Official docs, books, resources

### Quality Standards

- [ ] All code examples tested and working
- [ ] Consistent with existing guides' style
- [ ] Covers major frameworks for the language
- [ ] Includes real-world patterns
- [ ] References official documentation

### Example Structure

See existing guides for reference:

- [TypeScript Guide](../languages/typescript.md)
- [Python Guide](../languages/python.md)
- [Go Guide](../languages/go.md)

---

## Documentation Standards

When improving documentation:

### Style

- Clear, concise language
- Active voice
- Present tense
- Short paragraphs

### Formatting

- Use headers for organization
- Include code examples
- Use admonitions for notes/warnings
- Add links to related content

### Code Examples

- Test all examples
- Include both good and bad examples
- Show complete, runnable code
- Add comments for clarity

---

## Review Process

1. **Automated checks** - Linting, formatting
2. **Maintainer review** - Content and quality
3. **Testing** - Verified with real use cases
4. **Merge** - Squash and merge to main

Typical review time: 1-5 days

---

## Community

### Discussions

Use GitHub Discussions for:

- Questions
- Ideas
- Showcasing usage
- General discussion

### Issues

Use GitHub Issues for:

- Bug reports
- Feature requests
- Documentation errors

---

## Recognition

Contributors are recognized in:

- README acknowledgments
- Release notes
- This documentation

---

## Code of Conduct

Be respectful and constructive:

- Welcome newcomers
- Provide helpful feedback
- Focus on ideas, not people
- Assume good intentions

---

## Questions?

- Open a [Discussion](https://github.com/ar4mirez/aicof/discussions)
- Check existing [Issues](https://github.com/ar4mirez/aicof/issues)
- Review this documentation

Thank you for contributing!
