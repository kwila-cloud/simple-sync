import { execSync } from 'child_process';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

async function testPerformance() {
  try {
    console.log('Testing build performance...');
    const startTime = Date.now();

    execSync('npm run build', { stdio: 'pipe' });

    const buildTime = Date.now() - startTime;
    console.log(`✅ Build completed in ${buildTime}ms`);

    // Check bundle size
    const distPath = path.join(__dirname, 'dist');
    const stats = fs.statSync(distPath);
    console.log(`📦 Dist size: ${stats.size} bytes`);

    // Thresholds
    if (buildTime > 30000) { // 30 seconds
      console.log('⚠️  Build time exceeds 30s threshold');
      process.exit(1);
    }

    if (stats.size > 50 * 1024 * 1024) { // 50MB
      console.log('⚠️  Bundle size exceeds 50MB threshold');
      process.exit(1);
    }

    console.log('✅ Performance test passed');
    process.exit(0);

  } catch (error) {
    console.log('❌ Performance test failed:', error.message);
    process.exit(1);
  }
}

testPerformance();