import{_ as o}from"./plugin-vue_export-helper-DlAUqK2U.js";import{r,o as l,c,a as m,w as a,b as p,d as n,e}from"./app--ZFwTmN-.js";const d={},u=p('<p>可支持任意框架的脚手架 jzero, 默认支持 go-zero</p><div class="hint-container tip"><p class="hint-container-title">目前还在定制规范中, 不能保证兼容性. 待 v1.0 后保证兼容性</p></div><div style="text-align:center;"><img src="https://oss.jaronnie.com/jzero.jpg" style="width:33%;" alt=""></div><h2 id="特性" tabindex="-1"><a class="header-anchor" href="#特性"><span>特性</span></a></h2><ul><li>企业级代码规范</li><li>grpc, grpc-gateway, api 三合一, 满足绝大部分场景业务需要</li><li>集成命令行框架 cobra, 轻松编写具备生产可用的命令行工具</li><li>支持多 proto 多 service, 减少开发耦合性</li><li>不修改源码, 完全同步 go-zero 新特性</li><li>一键创建项目, 快速拓展新业务, 减少心理负担</li><li>一键生成服务端代码, 数据库代码, 客户端 sdk, 大大提高开发测试效率</li><li>支持自定义模板, 基于模板新建项目和生成代码</li></ul><h2 id="快速开始" tabindex="-1"><a class="header-anchor" href="#快速开始"><span>快速开始</span></a></h2><figure><img src="https://oss.jaronnie.com/2024-04-30_10-10-52.gif" alt="2024-04-30_10-10-52" tabindex="0" loading="lazy"><figcaption>2024-04-30_10-10-52</figcaption></figure><div class="hint-container tip"><p class="hint-container-title">Windows 用户请在 powershell 下执行所有指令</p></div>',8),v=e("div",{class:"language-bash line-numbers-mode","data-ext":"sh","data-title":"sh"},[e("pre",{class:"language-bash"},[e("code",null,[e("span",{class:"token comment"},"# 一键创建项目"),n(`
`),e("span",{class:"token function"},"docker"),n(" run "),e("span",{class:"token parameter variable"},"--rm"),n(),e("span",{class:"token parameter variable"},"-v"),n(),e("span",{class:"token variable"},[n("${"),e("span",{class:"token environment constant"},"PWD"),n("}")]),n(`/quickstart:/app/quickstart jaronnie/jzero:latest new quickstart
`),e("span",{class:"token builtin class-name"},"cd"),n(` quickstart 
`),e("span",{class:"token comment"},"# 一键生成代码"),n(`
`),e("span",{class:"token function"},"docker"),n(" run "),e("span",{class:"token parameter variable"},"--rm"),n(),e("span",{class:"token parameter variable"},"-v"),n(),e("span",{class:"token variable"},[n("${"),e("span",{class:"token environment constant"},"PWD"),n("}")]),n(":/app/quickstart jaronnie/jzero:latest gen "),e("span",{class:"token parameter variable"},"-w"),n(` quickstart
`),e("span",{class:"token comment"},"# 下载依赖"),n(`
go mod tidy
`),e("span",{class:"token comment"},"# 启动项目"),n(`
go run main.go server
`)])]),e("div",{class:"line-numbers","aria-hidden":"true"},[e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"})])],-1),g=e("div",{class:"language-bash line-numbers-mode","data-ext":"sh","data-title":"sh"},[e("pre",{class:"language-bash"},[e("code",null,[e("span",{class:"token comment"},"# 安装 jzero"),n(`
go `),e("span",{class:"token function"},"install"),n(` github.com/jzero-io/jzero@latest
`),e("span",{class:"token comment"},"# 一键安装所需的工具"),n(`
jzero check
`),e("span",{class:"token comment"},"# 一键创建项目"),n(`
jzero new quickstart
`),e("span",{class:"token builtin class-name"},"cd"),n(` quickstart
`),e("span",{class:"token comment"},"# 一键生成代码"),n(`
jzero gen
`),e("span",{class:"token comment"},"# 下载依赖"),n(`
go mod tidy
`),e("span",{class:"token comment"},"# 启动服务端程序"),n(`
go run main.go server
`)])]),e("div",{class:"line-numbers","aria-hidden":"true"},[e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"})])],-1);function h(b,k){const i=r("CodeTabs");return l(),c("div",null,[u,m(i,{id:"59",data:[{id:"Docker"},{id:"jzero"}],"tab-id":"shell"},{title0:a(({value:t,isActive:s})=>[n("Docker")]),title1:a(({value:t,isActive:s})=>[n("jzero")]),tab0:a(({value:t,isActive:s})=>[v]),tab1:a(({value:t,isActive:s})=>[g]),_:1})])}const f=o(d,[["render",h],["__file","index.html.vue"]]),z=JSON.parse('{"path":"/","title":"首页","lang":"zh-CN","frontmatter":{"home":false,"icon":"home","title":"首页","description":"可支持任意框架的脚手架 jzero, 默认支持 go-zero 目前还在定制规范中, 不能保证兼容性. 待 v1.0 后保证兼容性 特性 企业级代码规范 grpc, grpc-gateway, api 三合一, 满足绝大部分场景业务需要 集成命令行框架 cobra, 轻松编写具备生产可用的命令行工具 支持多 proto 多 service, 减少开发耦...","head":[["meta",{"property":"og:url","content":"https://jzero.jaronnie.com/"}],["meta",{"property":"og:site_name","content":"Jzero Framework"}],["meta",{"property":"og:title","content":"首页"}],["meta",{"property":"og:description","content":"可支持任意框架的脚手架 jzero, 默认支持 go-zero 目前还在定制规范中, 不能保证兼容性. 待 v1.0 后保证兼容性 特性 企业级代码规范 grpc, grpc-gateway, api 三合一, 满足绝大部分场景业务需要 集成命令行框架 cobra, 轻松编写具备生产可用的命令行工具 支持多 proto 多 service, 减少开发耦..."}],["meta",{"property":"og:type","content":"article"}],["meta",{"property":"og:image","content":"https://oss.jaronnie.com/2024-04-30_10-10-52.gif"}],["meta",{"property":"og:locale","content":"zh-CN"}],["meta",{"property":"og:updated_time","content":"2024-06-03T14:28:45.000Z"}],["meta",{"property":"article:author","content":"jaronnie"}],["meta",{"property":"article:modified_time","content":"2024-06-03T14:28:45.000Z"}],["script",{"type":"application/ld+json"},"{\\"@context\\":\\"https://schema.org\\",\\"@type\\":\\"Article\\",\\"headline\\":\\"首页\\",\\"image\\":[\\"https://oss.jaronnie.com/2024-04-30_10-10-52.gif\\"],\\"dateModified\\":\\"2024-06-03T14:28:45.000Z\\",\\"author\\":[{\\"@type\\":\\"Person\\",\\"name\\":\\"jaronnie\\",\\"url\\":\\"https://github.com/jaronnie\\"}]}"]]},"headers":[{"level":2,"title":"特性","slug":"特性","link":"#特性","children":[]},{"level":2,"title":"快速开始","slug":"快速开始","link":"#快速开始","children":[]}],"git":{"createdTime":1712825833000,"updatedTime":1717424925000,"contributors":[{"name":"jaronnie","email":"jaron@jaronnie.com","commits":26}]},"readingTime":{"minutes":1.19,"words":356},"filePathRelative":"README.md","localizedDate":"2024年4月11日","autoDesc":true}');export{f as comp,z as data};