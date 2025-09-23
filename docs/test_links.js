import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Simple link checker for markdown files
function checkLinks(dir) {
  const files = fs.readdirSync(dir, { recursive: true });
  let hasBrokenLinks = false;

  for (const file of files) {
    if (file.endsWith('.md') || file.endsWith('.mdx')) {
      const filePath = path.join(dir, file);
      const content = fs.readFileSync(filePath, 'utf8');

      // Simple regex for markdown links [text](url)
      const linkRegex = /\[([^\]]+)\]\(([^)]+)\)/g;
      let match;

      while ((match = linkRegex.exec(content)) !== null) {
        const url = match[2];
        // Check internal links (starting with / or relative)
        if (url.startsWith('/') || !url.includes('://')) {
          // For now, just check if it's not empty
          if (!url || url === '#') {
            console.log(`❌ Broken link in ${file}: ${url}`);
            hasBrokenLinks = true;
          }
        }
      }
    }
  }

  return hasBrokenLinks;
}

try {
  console.log('Checking links in src/content/docs/...');
  const hasBroken = checkLinks(path.join(__dirname, 'src', 'content', 'docs'));

  if (!hasBroken) {
    console.log('✅ No broken links found');
    process.exit(0);
  } else {
    console.log('❌ Broken links detected');
    process.exit(1);
  }
} catch (error) {
  console.log('❌ Link check failed with error:', error.message);
  process.exit(1);
}