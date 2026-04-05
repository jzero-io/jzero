# Icon Scripts

用于从 iconify.design 下载图标并更新 VuePress 文档中的图标引用。

## 使用方法

### 1. 下载图标

```bash
node scripts/download-icons.js
```

此脚本会：
- 扫描 `src/` 目录下所有 markdown 文件
- 提取所有 iconify 图标引用（格式：`icon: prefix:name`）
- 从 iconify API 下载对应的 SVG 文件
- 保存到 `src/.vuepress/public/icons/` 目录
- 生成图标映射到 `scripts/icon-map.json`

### 2. 替换图标引用

```bash
node scripts/replace-icons.js
```

此脚本会：
- 读取 `scripts/icon-map.json` 图标映射
- 遍历所有 markdown 文件
- 将 `icon: prefix:name` 替换为 `icon: /icons/prefix-name.svg`

## 完整流程

```bash
# 1. 下载所有图标
node scripts/download-icons.js

# 2. 更新文件引用
node scripts/replace-icons.js
```

## 文件说明

- `download-icons.js` - 下载图标脚本
- `replace-icons.js` - 替换图标引用脚本
- `icon-map.json` - 图标映射文件（自动生成）
