import { sidebar } from "vuepress-theme-hope";

export default sidebar({
  "/": [
    "",
    {
      text: "文档",
      icon: "vscode-icons:folder-type-docs",
      prefix: "guide/",
      children: "structure",
      collapsible: true,
      expanded: true,
    },
    {
      text: "项目实战",
      icon: "mdi:arrow-projectile-multiple",
      prefix: "project/",
      children: "structure",
      collapsible: true,
    },
    {
      text: "问题与解决",
      icon: "mdi:faq",
      prefix: "faq/",
      children: "structure",
      collapsible: true,
    },
    {
      text: "RoadMap",
      icon: "mdi:roadmap",
      prefix: "roadmap/",
      children: "structure",
      collapsible: true,
    },
  ],
});
