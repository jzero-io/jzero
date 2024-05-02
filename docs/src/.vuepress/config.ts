// @ts-ignore
import { defineUserConfig } from "vuepress";
import theme from "./theme.js";

export default defineUserConfig({
  base: "/",

  lang: "zh-CN",
  title: "Jzero Framework",
  description: "Jzero docs",

  theme,

  // 和 PWA 一起启用
  // shouldPrefetch: false,
});
