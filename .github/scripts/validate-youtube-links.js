const fs = require('node:fs');
const path = require('node:path');

// Helper to extract video ID from YouTube URL
function getVideoId(url) {
  // Regex to extract video ID from youtube embed url
  // https://www.youtube-nocookie.com/embed/xtuRdL8PTSA?si=...
  const match = url.match(/\/embed\/([a-zA-Z0-9_-]+)/);
  return match ? match[1] : null;
}

// Helper to check if video exists using oEmbed
async function checkVideo(videoId) {
  const oembedUrl = `https://www.youtube.com/oembed?url=https://www.youtube.com/watch?v=${videoId}&format=json`;
  try {
    const response = await fetch(oembedUrl);
    return response.status === 200;
  } catch (error) {
    console.warn(`  [WARN] Request failed for ${videoId}: ${error.message}`);
    return false;
  }
}

async function validateFile(filePath) {
  console.log(`Validating ${filePath}...`);

  let data;
  try {
    const content = fs.readFileSync(filePath, 'utf-8');
    data = JSON.parse(content);
  } catch (error) {
    console.error(`  [ERROR] Failed to load JSON: ${error.message}`);
    return false;
  }

  let failed = false;

  for (const drill of data) {
    const slug = drill.slug || 'Unknown';
    const videoUrls = drill.video_url || [];

    for (const url of videoUrls) {
      if (!url || url.trim() === '') {
        continue;
      }
      const videoId = getVideoId(url);

      if (!videoId) {
        console.error(`  [ERROR] Could not extract ID from URL: ${url} (Drill: ${slug})`);
        failed = true;
        continue;
      }

      const exists = await checkVideo(videoId);
      if (!exists) {
        console.error(`  [ERROR] Video unavailable: ${videoId} (URL: ${url}, Drill: ${slug})`);
        failed = true;
      }
    }
  }

  return !failed;
}

async function main() {
  // Resolve paths relative to the project root
  // This script is in .github/scripts/, so root is ../../
  const rootDir = path.resolve(__dirname, '../../');
  const filesToCheck = [
    path.join(rootDir, 'data', 'drills', 'de.json'),
    path.join(rootDir, 'data', 'drills', 'en.json')
  ];

  let allPassed = true;

  for (const fp of filesToCheck) {
    if (!fs.existsSync(fp)) {
      console.error(`File not found: ${fp}`);
      allPassed = false;
      continue;
    }

    const passed = await validateFile(fp);
    if (!passed) {
      allPassed = false;
    }
  }

  if (!allPassed) {
    process.exit(1);
  }

  console.log("All videos validated successfully.");
}

main();
