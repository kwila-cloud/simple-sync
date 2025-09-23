# Research: Build a docs site to replace our docs directory

## Astro Framework Research
**Decision**: Use Astro as the static site generator for the documentation site.

**Rationale**: Astro is designed for content-focused websites like documentation, provides excellent performance with static generation, and has great integration with Starlight for documentation sites. It's lightweight and focused on content delivery.

**Alternatives considered**:
- Next.js: More complex for a docs site, overkill for static content
- Hugo: Good for docs but less flexible for customization
- Docusaurus: Similar to Starlight but Astro provides better performance

## Starlight Theme Research
**Decision**: Use Starlight as the documentation theme on top of Astro.

**Rationale**: Starlight is specifically built for documentation sites, provides excellent navigation, search, and theming out of the box. It integrates seamlessly with Astro and supports markdown content structure.

**Alternatives considered**:
- Docusaurus: Similar features but heavier
- MkDocs: Python-based, not JavaScript
- Custom theme: Would require more development time

## GitHub Pages Hosting Research
**Decision**: Host the site on GitHub Pages with automatic deployment.

**Rationale**: GitHub Pages is free, integrates directly with the repository, and supports custom domains. GitHub Actions can automate the build and deployment process.

**Alternatives considered**:
- Netlify: More features but adds external dependency
- Vercel: Similar to Netlify
- Self-hosted: More complex infrastructure

## PDF Generation Research
**Decision**: Use starlight-to-pdf tool for generating PDF versions of the documentation.

**Rationale**: The tool is specifically designed for Starlight sites, integrates into the build process, and provides a clean PDF output for offline reading.

**Alternatives considered**:
- Puppeteer custom script: More complex to implement
- Other PDF tools: May not integrate as well with Starlight

## Build Process Research
**Decision**: Use GitHub Actions for CI/CD with build and deployment automation.

**Rationale**: Integrates with GitHub Pages, can run on every push to main branch, and can include PDF generation in the workflow.

**Alternatives considered**:
- Manual deployment: Error-prone and time-consuming
- Other CI services: Adds external dependencies