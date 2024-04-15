import{_ as e}from"./plugin-vue_export-helper-DlAUqK2U.js";import{o as n,c as a,a as s}from"./app-BndQm_1n.js";const i={},t=s(`<p>基于 go-zero 框架项目的代码设计</p><h2 id="技术栈" tabindex="-1"><a class="header-anchor" href="#技术栈"><span>技术栈</span></a></h2><ul><li>cobra 实现命令行管理</li><li>go-zero 提供 grpc 和 http 请求等</li></ul><h2 id="特性" tabindex="-1"><a class="header-anchor" href="#特性"><span>特性</span></a></h2><ul><li>支持将 grpc 通过 gateway 转化为 http 请求, 并支持自定义 http 请求</li><li>同时支持在项目中使用 grpc, grpc-gateway, api</li><li>支持监听 unix socket</li><li>支持多 proto 多 service(多人开发友好)</li><li>一键创建项目(jzero new)</li><li>一键生成各种代码(jzero gen)</li></ul><h2 id="安装" tabindex="-1"><a class="header-anchor" href="#安装"><span>安装</span></a></h2><div class="language-bash line-numbers-mode" data-ext="sh" data-title="sh"><pre class="language-bash"><code>go <span class="token function">install</span> github.com/jaronnie/jzero@latest
<span class="token comment"># 初始化</span>
jzero init
<span class="token comment"># 启动样例服务</span>
jzero daemon
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><h2 id="一键创建项目" tabindex="-1"><a class="header-anchor" href="#一键创建项目"><span>一键创建项目</span></a></h2><div class="language-bash line-numbers-mode" data-ext="sh" data-title="sh"><pre class="language-bash"><code><span class="token comment"># 安装 goctl</span>
go <span class="token function">install</span> github.com/zeromicro/go-zero/tools/goctl@latest
<span class="token comment"># 一键安装相关工具</span>
goctl <span class="token function">env</span> check <span class="token parameter variable">--install</span> <span class="token parameter variable">--verbose</span> <span class="token parameter variable">--force</span>
<span class="token comment"># 安装 jzero</span>
go <span class="token function">install</span> github.com/jaronnie/jzero@latest
<span class="token comment"># 一键创建项目</span>
jzero new <span class="token parameter variable">--module</span><span class="token operator">=</span>github.com/jaronnie/app1 <span class="token parameter variable">--dir</span><span class="token operator">=</span>./app1 <span class="token parameter variable">--app</span><span class="token operator">=</span>app1
<span class="token builtin class-name">cd</span> app1
<span class="token comment"># 一键生成代码</span>
jzero gen
<span class="token comment"># 下载依赖</span>
go mod tidy
<span class="token comment"># 启动项目</span>
go run main.go daemon <span class="token parameter variable">--config</span> config.toml
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><h2 id="开发" tabindex="-1"><a class="header-anchor" href="#开发"><span>开发</span></a></h2><div class="language-bash line-numbers-mode" data-ext="sh" data-title="sh"><pre class="language-bash"><code><span class="token comment"># gencode</span>
jzero gen

<span class="token comment"># run</span>
go run main.go daemon <span class="token parameter variable">--config</span> config.toml

<span class="token comment"># test</span>
<span class="token comment"># gateway</span>
<span class="token function">curl</span> http://localhost:8001/api/v1.0/credential/version
<span class="token comment"># grpc</span>
grpcurl <span class="token parameter variable">-plaintext</span> localhost:8000 credentialpb.credential/CredentialVersion
<span class="token comment"># api</span>
<span class="token function">curl</span> http://localhost:8001/api/v1/hello/me
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div>`,11),l=[t];function o(r,c){return n(),a("div",null,l)}const m=e(i,[["render",o],["__file","index.html.vue"]]),v=JSON.parse('{"path":"/","title":"首页","lang":"zh-CN","frontmatter":{"home":false,"icon":"home","title":"首页","description":"基于 go-zero 框架项目的代码设计 技术栈 cobra 实现命令行管理 go-zero 提供 grpc 和 http 请求等 特性 支持将 grpc 通过 gateway 转化为 http 请求, 并支持自定义 http 请求 同时支持在项目中使用 grpc, grpc-gateway, api 支持监听 unix socket 支持多 prot...","head":[["meta",{"property":"og:url","content":"https://jaronnie.github.io/jzero/"}],["meta",{"property":"og:site_name","content":"Jzero Framework"}],["meta",{"property":"og:title","content":"首页"}],["meta",{"property":"og:description","content":"基于 go-zero 框架项目的代码设计 技术栈 cobra 实现命令行管理 go-zero 提供 grpc 和 http 请求等 特性 支持将 grpc 通过 gateway 转化为 http 请求, 并支持自定义 http 请求 同时支持在项目中使用 grpc, grpc-gateway, api 支持监听 unix socket 支持多 prot..."}],["meta",{"property":"og:type","content":"article"}],["meta",{"property":"og:locale","content":"zh-CN"}],["meta",{"property":"og:updated_time","content":"2024-04-15T05:30:28.000Z"}],["meta",{"property":"article:author","content":"jaronnie"}],["meta",{"property":"article:modified_time","content":"2024-04-15T05:30:28.000Z"}],["script",{"type":"application/ld+json"},"{\\"@context\\":\\"https://schema.org\\",\\"@type\\":\\"Article\\",\\"headline\\":\\"首页\\",\\"image\\":[\\"\\"],\\"dateModified\\":\\"2024-04-15T05:30:28.000Z\\",\\"author\\":[{\\"@type\\":\\"Person\\",\\"name\\":\\"jaronnie\\",\\"url\\":\\"https://github.com/jaronnie\\"}]}"]]},"headers":[{"level":2,"title":"技术栈","slug":"技术栈","link":"#技术栈","children":[]},{"level":2,"title":"特性","slug":"特性","link":"#特性","children":[]},{"level":2,"title":"安装","slug":"安装","link":"#安装","children":[]},{"level":2,"title":"一键创建项目","slug":"一键创建项目","link":"#一键创建项目","children":[]},{"level":2,"title":"开发","slug":"开发","link":"#开发","children":[]}],"git":{"createdTime":1712825833000,"updatedTime":1713159028000,"contributors":[{"name":"jaronnie","email":"jaron@jaronnie.com","commits":3}]},"readingTime":{"minutes":0.8,"words":239},"filePathRelative":"README.md","localizedDate":"2024年4月11日","autoDesc":true}');export{m as comp,v as data};
