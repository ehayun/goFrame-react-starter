const fs = require('fs')
const path = require('path')

// Helper to copy directory recursively
function copyDir(src, dest) {
  if (!fs.existsSync(dest)) {
    fs.mkdirSync(dest, { recursive: true })
  }

  const entries = fs.readdirSync(src, { withFileTypes: true })

  for (let entry of entries) {
    const srcPath = path.join(src, entry.name)
    const destPath = path.join(dest, entry.name)

    if (entry.isDirectory()) {
      copyDir(srcPath, destPath)
    } else {
      fs.copyFileSync(srcPath, destPath)
    }
  }
}

// Concatenate CSS files into main.css
function concatenateCSS() {
  console.log('Concatenating CSS files...')

  const cssFiles = [
    'assets/css/remixicon.css',
    'assets/css/lib/bootstrap.min.css',
    'assets/css/lib/apexcharts.css',
    'assets/css/lib/dataTables.min.css',
    'assets/css/style.css'
  ]

  let combinedCSS = ''

  for (const file of cssFiles) {
    if (fs.existsSync(file)) {
      const content = fs.readFileSync(file, 'utf8')
      combinedCSS += `/* ${path.basename(file)} */\n${content}\n\n`
    }
  }

  // Ensure public/css directory exists
  if (!fs.existsSync('public/css')) {
    fs.mkdirSync('public/css', { recursive: true })
  }

  fs.writeFileSync('public/css/main.css', combinedCSS)
  console.log('✓ CSS concatenated to main.css')
}

// Concatenate JS files into main.js
function concatenateJS() {
  console.log('Concatenating JS files...')

  const jsFiles = [
    'assets/js/lib/jquery-3.7.1.min.js',
    'assets/js/lib/bootstrap.bundle.min.js',
    'assets/js/lib/apexcharts.min.js',
    'assets/js/lib/dataTables.min.js'
  ]

  let combinedJS = ''

  for (const file of jsFiles) {
    if (fs.existsSync(file)) {
      const content = fs.readFileSync(file, 'utf8')
      combinedJS += `/* ${path.basename(file)} */\n${content}\n\n`
    }
  }

  // Ensure public/js directory exists
  if (!fs.existsSync('public/js')) {
    fs.mkdirSync('public/js', { recursive: true })
  }

  fs.writeFileSync('public/js/main.js', combinedJS)
  console.log('✓ JS concatenated to main.js')
}

async function buildAssets() {
  console.log('Building assets...')

  // Concatenate CSS into main.css
  concatenateCSS()

  // Concatenate JS into main.js
  concatenateJS()

  // Copy remaining CSS files (for backup/reference)
  console.log('Copying CSS files...')
  if (fs.existsSync('assets/css')) {
    copyDir('assets/css', 'public/css')
    console.log('✓ CSS copied')
  }

  // Copy JavaScript files
  console.log('Copying JavaScript files...')
  if (fs.existsSync('assets/js')) {
    copyDir('assets/js', 'public/js')
    console.log('✓ JavaScript copied')
  }

  // Copy images
  console.log('Copying images...')
  if (fs.existsSync('assets/images')) {
    copyDir('assets/images', 'public/images')
    console.log('✓ Images copied')
  }

  // Copy fonts
  console.log('Copying fonts...')
  if (fs.existsSync('assets/fonts')) {
    copyDir('assets/fonts', 'public/fonts')
    console.log('✓ Fonts copied')
  }

  // Copy webfonts
  if (fs.existsSync('assets/webfonts')) {
    copyDir('assets/webfonts', 'public/webfonts')
    console.log('✓ Webfonts copied')
  }

  console.log('Asset build complete!')
}

buildAssets().catch(err => {
  console.error('Build failed:', err)
  process.exit(1)
})
