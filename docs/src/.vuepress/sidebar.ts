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
        text: "社区",
        icon: "iconoir:community",
        prefix: "community/",
        children: "structure",
        collapsible: true,
        expanded: true,
    },
  ],
});
