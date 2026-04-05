import fs from 'fs';
import path from 'path';
import https from 'https';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// 从 iconify API 下载 SVG
function downloadIcon(prefix, name) {
  const url = `https://api.iconify.design/${prefix}/${name}.svg`;

  return new Promise((resolve, reject) => {
    https.get(url, (res) => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => {
        if (res.statusCode === 200) {
          resolve(data);
        } else {
          reject(new Error(`HTTP ${res.statusCode}`));
        }
      });
    }).on('error', reject);
  });
}

// 从 markdown 文件中提取所有图标
function extractIconsFromMarkdown(dir) {
  const icons = new Set();

  function traverse(currentDir) {
    const files = fs.readdirSync(currentDir);

    files.forEach(file => {
      const filePath = path.join(currentDir, file);
      const stat = fs.statSync(filePath);

      if (stat.isDirectory() && !filePath.includes('node_modules')) {
        traverse(filePath);
      } else if (file.endsWith('.md')) {
        const content = fs.readFileSync(filePath, 'utf-8');
        const matches = content.match(/^icon:\s+(\S+)$/gm);
        if (matches) {
          matches.forEach(match => {
            const iconMatch = match.match(/icon:\s+(\S+)/);
            if (iconMatch && !iconMatch[1].startsWith('/icons/')) {
              icons.add(iconMatch[1]);
            }
          });
        }
      }
    });
  }

  traverse(dir);
  return icons;
}

// 主函数
async function main() {
  const srcDir = path.join(__dirname, '../src');
  const outputDir = path.join(__dirname, '../src/.vuepress/public/icons');

  if (!fs.existsSync(outputDir)) {
    fs.mkdirSync(outputDir, { recursive: true });
  }

  console.log('扫描 markdown 文件...');
  const icons = extractIconsFromMarkdown(srcDir);
  console.log(`找到 ${icons.size} 个唯一图标`);

  const iconMap = {};
  let successCount = 0;
  let failCount = 0;

  for (const icon of icons) {
    const [prefix, name] = icon.split(':');
    const filename = `${prefix}-${name}.svg`;
    const filepath = path.join(outputDir, filename);

    try {
      console.log(`下载 ${icon}...`);
      const svg = await downloadIcon(prefix, name);
      fs.writeFileSync(filepath, svg);
      iconMap[icon] = `/icons/${filename}`;
      successCount++;
      console.log(`  ✓ 保存到 ${filename}`);
    } catch (error) {
      failCount++;
      console.error(`  ✗ 失败: ${error.message}`);
    }
  }

  // 保存图标映射
  const mapPath = path.join(__dirname, 'icon-map.json');
  fs.writeFileSync(mapPath, JSON.stringify(iconMap, null, 2));

  console.log(`\n图标映射已保存到 ${path.relative(__dirname, mapPath)}`);
  console.log(`成功: ${successCount}, 失败: ${failCount}`);
}

main().catch(console.error);
