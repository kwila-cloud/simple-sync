# Quickstart: Build and Deploy Docs Site

## Prerequisites
- Node.js 18+
- GitHub repository with docs/ directory

## Local Development

### 1. Install Dependencies
```bash
cd docs
npm install
```

### 2. Start Development Server
```bash
npm run dev
```
Open http://localhost:4321 to view the site.

### 3. Build for Production
```bash
npm run build
```

## Deployment

### GitHub Pages Setup
1. Go to repository Settings > Pages
2. Set source to "GitHub Actions"
3. The workflow will automatically deploy on pushes to main

### Manual Deployment
```bash
npm run build
# Copy dist/ contents to GitHub Pages branch
```

## PDF Generation
The PDF is automatically generated during the build process using starlight-to-pdf.

## Testing
- Check all links work
- Verify search functionality
- Test on different browsers
- Validate PDF generation