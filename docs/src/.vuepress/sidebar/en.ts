// @ts-ignore
import { sidebar } from "vuepress-theme-hope";

export const enSidebarConfig = sidebar({
    "/": [
        "",
        {
            text: "Quick Start",
            icon: "streamline-sharp-color:startup",
            prefix: "getting-started/",
            children: "structure",
            collapsible: true,
            expanded: true,
        },
        {
            text: "Guide",
            icon: "icon-park-twotone:guide-board",
            prefix: "guide/",
            children: "structure",
            collapsible: true,
            expanded: true,
        },
        {
            text: "Components",
            icon: "iconamoon:component-bold",
            prefix: "component/",
            children: "structure",
            collapsible: true,
            expanded: true,
        },
        {
            text: "Ecosystem",
            icon: "material-icon-theme:pm2-ecosystem",
            prefix: "ecosystem/",
            children: "structure",
            collapsible: true,
            expanded: true,
        },
        {
            text: "Community",
            icon: "iconoir:community",
            prefix: "community/",
            children: "structure",
            collapsible: true,
            expanded: true,
        },
        {
            text: "Release Notes",
            icon: "catppuccin:release",
            prefix: "release/",
            children: "structure",
            collapsible: true,
            expanded: true,
        },
        {
            text: "Blog",
            icon: "fluent-mdl2:blog",
            prefix: "blog/",
            children: "structure",
            collapsible: true,
            expanded: false,
        },
    ]
});
