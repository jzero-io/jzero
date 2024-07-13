// @ts-ignore
import { hopeTheme } from "vuepress-theme-hope";
import navbar from "./navbar.js";
import sidebar from "./sidebar.js";

export default hopeTheme({
  hostname: "https://jzero.jaronnie.com",

  author: {
    name: "jaronnie",
    url: "https://github.com/jaronnie",
  },

  iconAssets: "iconify",

  // made by https://gopherize.me
  // favicon.ico made by https://www.bitbug.net
  logo: "https://oss.jaronnie.com/jzero.jpg",

  repo: "jzero-io/jzero",

  docsDir: "docs/src",

  // 导航栏
  navbar,

  // 侧边栏
  sidebar,

  // 页脚
  footer: "",
  displayFooter: true,

  // 多语言配置
  metaLocales: {
    editLink: "在 GitHub 上编辑此页",
  },

  // 如果想要实时查看任何改变，启用它。注: 这对更新性能有很大负面影响
  // hotReload: true,

  // 在这里配置主题提供的插件
  plugins: {
    blog: {
      category: "category",
      tag: "tag",
      star: "star",
    },
    comment: {
      provider: "Giscus",
      repo: "jzero-io/jzero",
      repoId: "R_kgDOLq1_9Q",
      category: "Announcements",
      categoryId: "DIC_kwDOLq1_9c4Cf5lp",
    },

    components: {
      components: ["Badge", "VPCard"],
    },

    // 此处开启了很多功能用于演示，你应仅保留用到的功能。
    mdEnhance: {
      align: true,
      attrs: true,
      codetabs: true,
      component: true,
      demo: true,
      figure: true,
      imgLazyload: true,
      imgSize: true,
      include: true,
      mark: true,
      stylize: [
        {
          matcher: "Recommended",
          replacer: ({ tag }) => {
            if (tag === "em")
              return {
                tag: "Badge",
                attrs: { type: "tip" },
                content: "Recommended",
              };
          },
        },
      ],
      sub: true,
      sup: true,
      tabs: true,
      tasklist: true,
      vPre: true,
    },
  },
});
