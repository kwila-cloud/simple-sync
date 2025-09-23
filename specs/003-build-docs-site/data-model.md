# Data Model: Documentation Site

## Project Structure (Astro + Starlight)

### Source Code Layout
```
docs/
├── astro.config.mjs         # Astro configuration with Starlight integration
├── package.json             # Dependencies and scripts
├── tsconfig.json            # TypeScript configuration
├── src/
│   ├── content/
│   │   └── docs/
│   │       ├── index.md     # Homepage
│   │       ├── api/
│   │       │   ├── v1.md   # API v1 documentation
│   │       │   └── ...
│   │       ├── acl.md      # ACL documentation
│   │       ├── tech-stack.md # Technology stack
│   │       └── ...
│   └── env.d.ts            # TypeScript environment
├── public/                  # Static assets (images, etc.)
└── dist/                    # Build output (generated)
```

### Content Organization
- **Root Level**: Main sections in src/content/docs/
- **Subdirectories**: Organized by feature or component
- **Files**: Markdown files with frontmatter for metadata

### Frontmatter Schema
Each markdown file includes:
```yaml
---
title: "Page Title"
description: "Brief description"
sidebar:
  label: "Display Label"
  order: 1
---
```

### Navigation Model
- **Sidebar**: Hierarchical navigation based on directory structure
- **Breadcrumbs**: Path-based navigation
- **Search**: Full-text search across all content
- **Table of Contents**: Auto-generated from headings

### Content Types
- **Reference Docs**: API endpoints, configuration
- **Guides**: Tutorials, setup instructions
- **Examples**: Code samples, use cases

## Build Model

### Static Generation
- Markdown files in src/content/docs/ processed at build time
- HTML output in dist/ with navigation and styling
- Assets (images, CSS, JS) optimized and bundled

### Deployment Model
- Built site in dist/ pushed to GitHub Pages
- PDF version generated and included
- Automatic updates on repository changes
