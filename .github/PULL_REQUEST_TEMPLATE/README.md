# Pull Request Templates

This directory contains multiple PR templates to help you create comprehensive and consistent pull requests.

## 📋 Available Templates

### 1. Default Template
**File:** `../pull_request_template.md`
- **When to use:** Comprehensive template for most PRs
- **Best for:** Major features, significant changes, or when you want detailed documentation

### 2. Simple Template
**File:** `simple.md`
- **When to use:** Quick fixes, minor changes, or simple updates
- **Best for:** Small bug fixes, typos, documentation updates

### 3. Bug Fix Template
**File:** `bug_fix.md`
- **When to use:** When fixing bugs or issues
- **Best for:** Bug reports, error corrections, issue resolutions

### 4. Feature Template
**File:** `feature.md`
- **When to use:** When adding new features or functionality
- **Best for:** New features, enhancements, major additions

## 🚀 How to Use Templates

### Method 1: URL Parameters (Recommended)
Add a template parameter to your PR URL:

```
https://github.com/yourorg/yourrepo/compare/main...your-branch?template=simple.md
```

### Method 2: Manual Copy-Paste
1. Choose the appropriate template from the list above
2. Copy the template content
3. Paste it into your PR description
4. Fill in the relevant sections

### Method 3: GitHub CLI
```bash
# Create PR with specific template
gh pr create --template .github/PULL_REQUEST_TEMPLATE/simple.md
```

## 📝 Template Selection Guide

| Change Type | Recommended Template | Why? |
|-------------|---------------------|------|
| Bug Fix | `bug_fix.md` | Includes root cause analysis and testing verification |
| New Feature | `feature.md` | Comprehensive feature documentation and testing |
| Documentation | `simple.md` | Lightweight template for docs changes |
| Refactoring | Default | Detailed analysis of changes and impact |
| Quick Fix | `simple.md` | Minimal overhead for small changes |
| Breaking Change | `feature.md` or Default | Detailed impact analysis required |

## ✅ Best Practices

### Before Creating a PR
- [ ] Choose the right template for your change type
- [ ] Read through the entire template
- [ ] Remove unused sections (marked as optional)
- [ ] Fill in all required information

### Writing Descriptions
- **Be specific:** Avoid vague descriptions like "fixed stuff"
- **Include context:** Explain why the change was needed
- **Add testing details:** How did you verify the change works?
- **Link related issues:** Use `Fixes #123` or `Relates to #456`

### For Reviewers
- Check that the appropriate template was used
- Ensure all required sections are completed
- Verify testing information is provided
- Look for clear problem/solution explanation

## 🔧 Customizing Templates

### Adding New Templates
1. Create a new `.md` file in this directory
2. Follow the existing naming convention
3. Include clear sections and checkboxes
4. Update this README with the new template

### Modifying Existing Templates
1. Edit the template file directly
2. Test with a sample PR
3. Update documentation if needed
4. Consider backward compatibility

## 📚 Template Sections Explained

### Common Sections
- **Description:** What the PR does and why
- **Type of Change:** Categorizes the change
- **Testing:** How the change was verified
- **Checklist:** Pre-merge verification items

### Advanced Sections
- **Performance Impact:** For changes that affect performance
- **Security Review:** For security-sensitive changes
- **Breaking Changes:** For changes that break existing functionality
- **Deployment Notes:** For changes requiring special deployment steps

## 🎯 Examples

### Good PR Title Examples
- `feat: add PostgreSQL data scraping functionality`
- `fix: resolve memory leak in article processing`
- `docs: update installation instructions`
- `refactor: optimize database connection pooling`

### Bad PR Title Examples
- `updates`
- `fix stuff`
- `changes`
- `wip`

## 📞 Need Help?

- **Questions about templates:** Create an issue with the `documentation` label
- **Template improvements:** Submit a PR with your suggested changes
- **Usage guidance:** Check with team leads or maintainers

---

*Remember: Good PR descriptions help reviewers understand your changes quickly and lead to faster, more effective code reviews!*