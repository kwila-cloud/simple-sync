import { execSync } from 'child_process';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

try {
  console.log('Running build...');
  execSync('npm run build', { stdio: 'inherit' });

  const distPath = path.join(__dirname, 'dist');
  if (fs.existsSync(distPath)) {
    console.log('✅ Build successful - dist/ directory created');
    process.exit(0);
  } else {
    console.log('❌ Build failed - dist/ directory not found');
    process.exit(1);
  }
} catch (error) {
  console.log('❌ Build failed with error:', error.message);
  process.exit(1);
}