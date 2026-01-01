// @ts-ignore
import { hopeTheme } from "vuepress-theme-hope";
import { enNavbar, zhNavbar } from "./navbar/index.js";
import sidebar from "./sidebar.js";

// @ts-ignore
export default hopeTheme({
  hostname: "https://docs.jzero.io",

  author: {
    name: "jaronnie",
    url: "https://github.com/jaronnie",
  },

  iconAssets: "iconify",

  copyright: 'Copyright © 2024-2026 jzero-io',

  // made by https://gopherize.me
  // favicon.ico made by https://www.bitbug.net
  logo: "https://oss.jaronnie.com/jzero.svg",

  repo: "jzero-io/jzero",

  docsDir: "docs/src",

  locales: {
    "/": {
      // 导航栏
      navbar: zhNavbar,

      // 侧边栏
      sidebar,

      // 页脚
      footer: "",
      displayFooter: true,

      // Page meta
      metaLocales: {
        editLink: "在 GitHub 上编辑此页",
      },
    },
    "/en/": {
      // 导航栏
      navbar: enNavbar,

      // Page meta
      metaLocales: {
        editLink: "Edit this page on GitHub",
      },
    },
  },

  // 在这里配置主题提供的插件
  plugins: {
    blog: {
      // category: "category",
      // tag: "tag",
      // star: "star",
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
      // demo: true,
      figure: true,
      imgLazyload: true,
      imgSize: true,
      include: true,
      mark: true,
      // stylize: [
      //   {
      //     matcher: "Recommended",
      //     replacer: ({ tag }) => {
      //       if (tag === "em")
      //         return {
      //           tag: "Badge",
      //           attrs: { type: "tip" },
      //           content: "Recommended",
      //         };
      //     },
      //   },
      // ],
      // sub: true,
      // sup: true,
      // tabs: true,
      // tasklist: true,
      // vPre: true,
    },
  },
});
