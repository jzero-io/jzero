---
title: "jzero × Vercel 生态打通：Go 语言无缝接入 Vercel 无服务器平台"
icon: /icons/emojione-v1-rocket.svg
---

## jzero × Vercel：生态打通的强大组合

### 为什么选择 jzero + Vercel？

**什么是 jzero？**

**jzero** 是基于 go-zero 框架开发的增强型开发工具：

🏗️ **通过模板生成基础框架代码**：基于描述文件自动生成框架代码（api → api 框架代码、proto → proto 框架代码、sql/远程数据库地址 → model 代码）

🤖 **通过 Agent Skills 生成业务代码**：内置 jzero-skills，让 AI 生成符合最佳实践的业务逻辑代码

**核心价值与设计理念**：

- ✅ **开发体验优先**：提供简单好用的一站式生产可用解决方案，一键初始化 api/rpc/gateway 项目，极简指令生成基础框架代码
- ✅ **AI 赋能**：内置 jzero-skills，让 AI 生成符合最佳实践的业务逻辑代码
- ✅ **模板驱动**：默认生成即最佳实践，支持自定义模板，可基于远程模板仓库打造企业专属底座
- ✅ **插件化架构**：模块分层、插件设计，团队协作更顺畅
- ✅ **内置开发组件**：包含缓存(cache)、数据库迁移(migrate)、配置中心(configcenter)、数据库查询(condition)等常用工具
- ✅ **生态兼容**：不修改 go-zero，保持生态兼容，解决已有痛点问题并扩展新功能
- ✅ **接口灵活**：不依赖特定数据库/缓存/配置中心，可根据实际需求自由选择

---

**为什么选择 jzero + Vercel？**

- ✅ **零配置部署**：jzero 生成的代码完全兼容 Vercel 平台，从 Git 仓库直接部署，无需额外配置
- ✅ **全球边缘网络**：借助 Vercel 的全球基础设施，将 Go 函数部署到离用户最近的边缘节点，实现毫秒级响应
- ✅ **免费域名与 HTTPS**：自动获得 `.vercel.app` 生产级域名，内置 CDN 加速和 SSL 证书
- ✅ **预览环境**：每次提交 PR 自动生成独立预览 URL

> 💡 **核心价值**：jzero 深度打通 Vercel 生态，让 Go 开发者享受前端级别的部署体验！通过 `.api` 定义自动生成符合 Vercel 规范的无服务器函数，真正实现"一次定义，处处运行"。

## 快速开始：一键接入 Vercel 生态

### 创建 Vercel 兼容项目

jzero 提供了专为 Vercel 生态优化的项目模板，生成的代码结构完全符合 Vercel 平台规范：

```bash
# 创建新的 Vercel 无服务器项目
jzero new jzero-api-vercel-example --branch api-vercel
```

**项目结构**：
```
jzero-api-vercel-example/
├── vercel/
│   └── client.go        # vercel go 运行时入口
├── desc/                # API 定义
│   └── api/             # .api 文件
├── server/              # 服务端
│   ├── handler/         # HTTP 处理器
│   ├── logic/           # 业务逻辑
│   └── types/           # 类型定义
├── vercel.json          # Vercel 平台配置
├── main.go              # 本地运行入口
└── go.mod               # Go 模块文件
```

## 一键部署：深度集成 Vercel 工作流

### Git 推送即部署

jzero 生成的项目完全兼容 Vercel 的 Git 工作流，实现代码推送自动部署：

```bash
cd jzero-api-vercel-example

# 初始化 git
git init
git add .
git commit -m "Initial commit"

# 在 GitHub 上创建仓库并推送
git remote add origin https://github.com/your-username/jzero-api-vercel-example.git
git push -u origin main
```

### Vercel 平台自动识别

1. 访问 [Vercel 控制台](https://vercel.com/dashboard)
2. 点击"Add New Project"

3. 导入您的 GitHub 仓库
> ![image-20260404112504507](https://oss.jaronnie.com/image-20260404112504507.png)
4. **vercel.json 配置让 Vercel 自动识别为 Go 项目**

5. 点击"Deploy"

> ![image-20260404112527983](https://oss.jaronnie.com/image-20260404112527983.png)

**🎉 部署完成！您的 Go API 已接入 Vercel 全球网络**

![image-20260404112555862](https://oss.jaronnie.com/image-20260404112555862.png)

6. 访问 api 

> ![image-20260404112631987](https://oss.jaronnie.com/image-20260404112631987.png)

**觉得有用？请给 jzero 一个 ⭐ Star 支持我们持续改进！**

GitHub: [https://github.com/jzero-io/jzero](https://github.com/jzero-io/jzero)
Jzero 官网: [https://jzero.io](https://jzero.io)
Vercel 模板: [https://github.com/jzero-io/templates/tree/api-vercel](https://github.com/jzero-io/templates/tree/api-vercel)
