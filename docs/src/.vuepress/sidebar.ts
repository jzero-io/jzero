// @ts-ignore
import { sidebar } from "vuepress-theme-hope";

export default sidebar({
  "/": [
    "",
    {
        text: "快速开始",
        icon: "streamline-sharp-color:startup",
        prefix: "getting-started/",
        children: "structure",
        collapsible: true,
        expanded: true,
    },
    {
        text: "指南",
        icon: "icon-park-twotone:guide-board",
        prefix: "guide/",
        children: "structure",
        collapsible: true,
        expanded: true,
    },
    {
        text: "组件",
        icon: "iconamoon:component-bold",
        prefix: "component/",
        children: "structure",
        collapsible: true,
        expanded: true,
    },
    {
        text: "生态",
        icon: "material-icon-theme:pm2-ecosystem",
        prefix: "ecosystem/",
        children: "structure",
        collapsible: true,
        expanded: true,
    },
    {
        text: "社区",
        icon: "iconoir:community",
        prefix: "community/",
        children: "structure",
        collapsible: true,
        expanded: true,
    },
    {
          text: "版本记录",
          icon: "catppuccin:release",
          prefix: "release/",
          children: "structure",
          collapsible: true,
          expanded: true,
    },
    {
        text: "博客",
        icon: "logos:bloglovin",
        prefix: "blog/",
        children: "structure",
        collapsible: true,
        expanded: false,
    },
  ],
});
