import { sidebar } from "vuepress-theme-hope";

export default sidebar({
  "/": [
    "",
    {
      text: "文档",
      icon: "vscode-icons:folder-type-docs",
      prefix: "guide/",
      children: "structure",
    },
  ],
});
