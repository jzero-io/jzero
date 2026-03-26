// @ts-ignore
import { hopeTheme } from "vuepress-theme-hope";
import { enNavbar, zhNavbar } from "./navbar/index.js";
import { enSidebarConfig, zhSidebarConfig } from "./sidebar/index.js";

// @ts-ignore
export default hopeTheme({
  hostname: "https://docs.jzero.io",

  author: {
    name: "jaronnie",
    url: "https://github.com/jaronnie",
  },

  copyright: 'Copyright © 2024-2026 jzero-io',

  // made by https://gopherize.me
  // favicon.ico made by https://www.bitbug.net
  logo: "https://oss.jaronnie.com/jzero.svg",

  repo: "jzero-io/jzero",

  docsDir: "docs/src",

  locales: {
    "/": {
      // 导航栏
      navbar: enNavbar,

      // 侧边栏
      sidebar: enSidebarConfig,

      // 页脚
      footer: "",
      displayFooter: true,

      // Page meta
      metaLocales: {
        editLink: "Edit this page on GitHub",
      },
    },
    "/zh-CN/": {
      // 导航栏
      navbar: zhNavbar,

      sidebar: zhSidebarConfig,

      // Page meta
      metaLocales: {
        editLink: "在 GitHub 上编辑此页",
      },
    },
  },

  // These features are enabled for demo, only preserve features you need here
  markdown: {
    align: true,
    attrs: true,
    codeTabs: true,
    component: true,
    demo: true,
    figure: true,
    gfm: true,
    imgLazyload: true,
    imgSize: true,
    include: true,
    mark: true,
    plantuml: true,
    spoiler: true,
    stylize: [
      {
        matcher: "Recommended",
        replacer: ({ tag }) => {
          if (tag === "em") {
            return {
              tag: "Badge",
              attrs: { type: "tip" },
              content: "Recommended",
            };
          }
        },
      },
    ],
    sub: true,
    sup: true,
    tabs: true,
    tasklist: true,
    vPre: true,

    // uncomment these if you need TeX support
    // math: {
    //   // install katex before enabling it
    //   type: "katex",
    //   // or install @mathjax/src before enabling it
    //   type: "mathjax",
    // },

    // install chart.js before enabling it
    // chartjs: true,

    // install echarts before enabling it
    // echarts: true,

    // install flowchart.ts before enabling it
    // flowchart: true,

    // install mermaid before enabling it
    // mermaid: true,

    // playground: {
    //   presets: ["ts", "vue"],
    // },

    // install @vue/repl before enabling it
    // vuePlayground: true,

    // install sandpack-vue3 before enabling it
    // sandpack: true,

    // install @vuepress/plugin-revealjs and uncomment these if you need slides
    // revealjs: {
    //   plugins: ["highlight", "math", "search", "notes", "zoom"],
    // },
  },

  // 在这里配置主题提供的插件
  plugins: {
    icon: {
      assets: "iconify"
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
  },
});
