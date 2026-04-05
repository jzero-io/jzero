import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// 读取图标映射
const iconMapPath = path.join(__dirname, 'icon-map.json');
const iconMap = JSON.parse(fs.readFileSync(iconMapPath, 'utf-8'));

console.log(`找到 ${Object.keys(iconMap).length} 个图标映射`);

// 获取所有 markdown 文件
function getAllMarkdownFiles(dir, fileList = []) {
  const files = fs.readdirSync(dir);

  files.forEach(file => {
    const filePath = path.join(dir, file);
    const stat = fs.statSync(filePath);

    if (stat.isDirectory() && !filePath.includes('node_modules')) {
      getAllMarkdownFiles(filePath, fileList);
    } else if (file.endsWith('.md')) {
      fileList.push(filePath);
    }
  });

  return fileList;
}

// 替换图标引用
function replaceIconInFile(filePath) {
  let content = fs.readFileSync(filePath, 'utf-8');
  let modified = false;

  // 替换 icon: prefix:name 为 icon: /icons/prefix-name.svg
  content = content.replace(/^icon:\s+(\S+)$/gm, (match, icon) => {
    if (iconMap[icon] && !icon.startsWith('/icons/')) {
      modified = true;
      return `icon: ${iconMap[icon]}`;
    }
    return match;
  });

  if (modified) {
    fs.writeFileSync(filePath, content);
    return true;
  }
  return false;
}

// 主函数
function main() {
  const srcDir = path.join(__dirname, '../src');
  const markdownFiles = getAllMarkdownFiles(srcDir);

  console.log(`找到 ${markdownFiles.length} 个 markdown 文件`);

  let modifiedCount = 0;
  markdownFiles.forEach(file => {
    const relativePath = path.relative(__dirname, file);
    if (replaceIconInFile(file)) {
      modifiedCount++;
      console.log(`✓ 已更新: ${relativePath}`);
    }
  });

  console.log(`\n总共更新了 ${modifiedCount} 个文件`);
}

main();
