import{_ as p}from"./plugin-vue_export-helper-DlAUqK2U.js";import{r as l,o as m,c as d,a as e,b as a,d as i,w as n,e as o}from"./app-C5jEYKWE.js";const u={},v={href:"https://go-zero.dev",target:"_blank",rel:"noopener noreferrer"},b=o('<div class="hint-container tip"><p class="hint-container-title">目前还在定制规范中, 不能保证兼容性. 待 v1.0 后保证兼容性</p></div><div style="text-align:center;"><img src="https://oss.jaronnie.com/jzero.jpg" style="width:33%;" alt=""></div><h2 id="特性" tabindex="-1"><a class="header-anchor" href="#特性"><span>特性</span></a></h2><ul><li>企业级代码规范</li><li>grpc, grpc-gateway, api 三合一, 满足绝大部分场景业务需要</li><li>支持监听 unix socket 本地通信</li><li>支持多 proto 多 service, 减少开发耦合性</li><li>一键创建项目, 快速拓展新业务, 减少心理负担</li><li>一键生成各种代码, 大大提高开发效率</li><li>支持流量治理, 减少线上风险</li><li>支持链路追踪, 监控等, 快速定位问题</li><li>所有工具链跨平台支持</li></ul><h2 id="快速开始" tabindex="-1"><a class="header-anchor" href="#快速开始"><span>快速开始</span></a></h2>',5),k=e("div",{class:"language-bash line-numbers-mode","data-ext":"sh","data-title":"sh"},[e("pre",{class:"language-bash"},[e("code",null,[e("span",{class:"token comment"},"# 一键创建项目"),a(`
`),e("span",{class:"token function"},"docker"),a(" run "),e("span",{class:"token parameter variable"},"--rm"),a(),e("span",{class:"token punctuation"},"\\"),a(`
  `),e("span",{class:"token parameter variable"},"-v"),a(" ./app1:/app/app1 jaronnie/jzero:latest "),e("span",{class:"token punctuation"},"\\"),a(`
  new `),e("span",{class:"token parameter variable"},"--module"),e("span",{class:"token operator"},"="),a("github.com/jaronnie/app1 "),e("span",{class:"token punctuation"},"\\"),a(`
  `),e("span",{class:"token parameter variable"},"--dir"),e("span",{class:"token operator"},"="),a("./app1 "),e("span",{class:"token parameter variable"},"--app"),e("span",{class:"token operator"},"="),a(`app1
  
`),e("span",{class:"token comment"},"# 一键生成代码"),a(`
`),e("span",{class:"token function"},"docker"),a(" run "),e("span",{class:"token parameter variable"},"--rm"),a(),e("span",{class:"token punctuation"},"\\"),a(`
  `),e("span",{class:"token parameter variable"},"-v"),a(" ./app1:/app/app1 jaronnie/jzero:latest "),e("span",{class:"token punctuation"},"\\"),a(`
  gen `),e("span",{class:"token parameter variable"},"-w"),a(` app1

`),e("span",{class:"token builtin class-name"},"cd"),a(` app1
`),e("span",{class:"token comment"},"# 下载依赖"),a(`
go mod tidy
`),e("span",{class:"token comment"},"# 启动项目"),a(`
go run main.go daemon `),e("span",{class:"token parameter variable"},"--config"),a(` config.toml
`)])]),e("div",{class:"line-numbers","aria-hidden":"true"},[e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"})])],-1),h=e("div",{class:"language-bash line-numbers-mode","data-ext":"sh","data-title":"sh"},[e("pre",{class:"language-bash"},[e("code",null,[e("span",{class:"token comment"},"# 一键创建项目"),a(`
`),e("span",{class:"token function"},"docker"),a(" run "),e("span",{class:"token parameter variable"},"--rm"),a(),e("span",{class:"token punctuation"},"\\"),a(`
  `),e("span",{class:"token parameter variable"},"-v"),a(" ./app1:/app/app1 jaronnie/jzero:latest-arm64 "),e("span",{class:"token punctuation"},"\\"),a(`
  new `),e("span",{class:"token parameter variable"},"--module"),e("span",{class:"token operator"},"="),a("github.com/jaronnie/app1 "),e("span",{class:"token punctuation"},"\\"),a(`
  `),e("span",{class:"token parameter variable"},"--dir"),e("span",{class:"token operator"},"="),a("./app1 "),e("span",{class:"token parameter variable"},"--app"),e("span",{class:"token operator"},"="),a(`app1
  
`),e("span",{class:"token comment"},"# 一键生成代码"),a(`
`),e("span",{class:"token function"},"docker"),a(" run "),e("span",{class:"token parameter variable"},"--rm"),a(),e("span",{class:"token punctuation"},"\\"),a(`
  `),e("span",{class:"token parameter variable"},"-v"),a(" ./app1:/app/app1 jaronnie/jzero:latest-arm64 "),e("span",{class:"token punctuation"},"\\"),a(`
  gen `),e("span",{class:"token parameter variable"},"-w"),a(` app1

`),e("span",{class:"token builtin class-name"},"cd"),a(` app1
`),e("span",{class:"token comment"},"# 下载依赖"),a(`
go mod tidy
`),e("span",{class:"token comment"},"# 启动项目"),a(`
go run main.go daemon `),e("span",{class:"token parameter variable"},"--config"),a(` config.toml
`)])]),e("div",{class:"line-numbers","aria-hidden":"true"},[e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"})])],-1),g=e("div",{class:"language-bash line-numbers-mode","data-ext":"sh","data-title":"sh"},[e("pre",{class:"language-bash"},[e("code",null,[e("span",{class:"token comment"},"# 安装 goctl"),a(`
go `),e("span",{class:"token function"},"install"),a(` github.com/zeromicro/go-zero/tools/goctl@latest
`),e("span",{class:"token comment"},"# 一键安装相关工具"),a(`
goctl `),e("span",{class:"token function"},"env"),a(" check "),e("span",{class:"token parameter variable"},"--install"),a(),e("span",{class:"token parameter variable"},"--verbose"),a(),e("span",{class:"token parameter variable"},"--force"),a(`
`),e("span",{class:"token comment"},"# 安装 jzero"),a(`
go `),e("span",{class:"token function"},"install"),a(` github.com/jaronnie/jzero@latest
`),e("span",{class:"token comment"},"# 一键创建项目"),a(`
jzero new `),e("span",{class:"token parameter variable"},"--module"),e("span",{class:"token operator"},"="),a("github.com/jaronnie/app1 "),e("span",{class:"token parameter variable"},"--dir"),e("span",{class:"token operator"},"="),a("./app1 "),e("span",{class:"token parameter variable"},"--app"),e("span",{class:"token operator"},"="),a(`app1
`),e("span",{class:"token builtin class-name"},"cd"),a(` app1
`),e("span",{class:"token comment"},"# 一键生成代码"),a(`
jzero gen
`),e("span",{class:"token comment"},"# 下载依赖"),a(`
go mod tidy
`),e("span",{class:"token comment"},"# 启动项目"),a(`
go run main.go daemon `),e("span",{class:"token parameter variable"},"--config"),a(` config.toml
`)])]),e("div",{class:"line-numbers","aria-hidden":"true"},[e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"}),e("div",{class:"line-number"})])],-1),f=o(`<h2 id="验证" tabindex="-1"><a class="header-anchor" href="#验证"><span>验证</span></a></h2><div class="language-bash line-numbers-mode" data-ext="sh" data-title="sh"><pre class="language-bash"><code><span class="token comment"># test</span>
<span class="token comment"># gateway</span>
<span class="token function">curl</span> http://localhost:8001/api/v1.0/credential/version
<span class="token comment"># grpc</span>
grpcurl <span class="token parameter variable">-plaintext</span> localhost:8000 credentialpb.credential/CredentialVersion
<span class="token comment"># api</span>
<span class="token function">curl</span> http://localhost:8001/api/v1/hello/me
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div>`,2);function _(j,x){const r=l("ExternalLinkIcon"),c=l("CodeTabs");return m(),d("div",null,[e("p",null,[a("基于 "),e("a",v,[a("go-zero"),i(r)]),a(" 框架定制的企业级后端代码框架")]),b,i(c,{id:"59",data:[{id:"Docker(amd64)"},{id:"Docker(arm64)"},{id:"jzero"}],"tab-id":"shell"},{title0:n(({value:s,isActive:t})=>[a("Docker(amd64)")]),title1:n(({value:s,isActive:t})=>[a("Docker(arm64)")]),title2:n(({value:s,isActive:t})=>[a("jzero")]),tab0:n(({value:s,isActive:t})=>[k]),tab1:n(({value:s,isActive:t})=>[h]),tab2:n(({value:s,isActive:t})=>[g]),_:1},8,["data"]),f])}const w=p(u,[["render",_],["__file","index.html.vue"]]),T=JSON.parse('{"path":"/","title":"首页","lang":"zh-CN","frontmatter":{"home":false,"icon":"home","title":"首页","description":"基于 go-zero 框架定制的企业级后端代码框架 目前还在定制规范中, 不能保证兼容性. 待 v1.0 后保证兼容性 特性 企业级代码规范 grpc, grpc-gateway, api 三合一, 满足绝大部分场景业务需要 支持监听 unix socket 本地通信 支持多 proto 多 service, 减少开发耦合性 一键创建项目, 快速拓展新...","head":[["meta",{"property":"og:url","content":"https://jaronnie.github.io/jzero/"}],["meta",{"property":"og:site_name","content":"Jzero Framework"}],["meta",{"property":"og:title","content":"首页"}],["meta",{"property":"og:description","content":"基于 go-zero 框架定制的企业级后端代码框架 目前还在定制规范中, 不能保证兼容性. 待 v1.0 后保证兼容性 特性 企业级代码规范 grpc, grpc-gateway, api 三合一, 满足绝大部分场景业务需要 支持监听 unix socket 本地通信 支持多 proto 多 service, 减少开发耦合性 一键创建项目, 快速拓展新..."}],["meta",{"property":"og:type","content":"article"}],["meta",{"property":"og:locale","content":"zh-CN"}],["meta",{"property":"og:updated_time","content":"2024-04-24T06:31:29.000Z"}],["meta",{"property":"article:author","content":"jaronnie"}],["meta",{"property":"article:modified_time","content":"2024-04-24T06:31:29.000Z"}],["script",{"type":"application/ld+json"},"{\\"@context\\":\\"https://schema.org\\",\\"@type\\":\\"Article\\",\\"headline\\":\\"首页\\",\\"image\\":[\\"\\"],\\"dateModified\\":\\"2024-04-24T06:31:29.000Z\\",\\"author\\":[{\\"@type\\":\\"Person\\",\\"name\\":\\"jaronnie\\",\\"url\\":\\"https://github.com/jaronnie\\"}]}"]]},"headers":[{"level":2,"title":"特性","slug":"特性","link":"#特性","children":[]},{"level":2,"title":"快速开始","slug":"快速开始","link":"#快速开始","children":[]},{"level":2,"title":"验证","slug":"验证","link":"#验证","children":[]}],"git":{"createdTime":1712825833000,"updatedTime":1713940289000,"contributors":[{"name":"jaronnie","email":"jaron@jaronnie.com","commits":11}]},"readingTime":{"minutes":1.38,"words":414},"filePathRelative":"README.md","localizedDate":"2024年4月11日","autoDesc":true}');export{w as comp,T as data};
