import{_ as n}from"./plugin-vue_export-helper-DlAUqK2U.js";import{c as e,a,o as i}from"./app-M4fn2MCd.js";const p={};function l(t,s){return i(),e("div",null,s[0]||(s[0]=[a(`<div class="hint-container tip"><p class="hint-container-title">提示</p><p><a href="https://go-zero.dev/docs/tutorials" target="_blank" rel="noopener noreferrer">go-zero api 文档</a></p></div><h2 id="api-字段校验" tabindex="-1"><a class="header-anchor" href="#api-字段校验"><span>api 字段校验</span></a></h2><blockquote><p>jzero 集成 <a href="https://github.com/go-playground/validator" target="_blank" rel="noopener noreferrer">https://github.com/go-playground/validator</a> 进行字段校验</p></blockquote><div class="language-api line-numbers-mode" data-highlighter="shiki" data-ext="api" data-title="api" style="--shiki-light:#24292e;--shiki-dark:#abb2bf;--shiki-light-bg:#fff;--shiki-dark-bg:#282c34;"><pre class="shiki shiki-themes github-light one-dark-pro vp-code"><code><span class="line"><span>syntax = &quot;v1&quot;</span></span>
<span class="line"><span></span></span>
<span class="line"><span>type CreateRequest {</span></span>
<span class="line"><span>    name string \`json:&quot;name&quot; validate:&quot;gte=2,lte=30&quot;\` // 名称</span></span>
<span class="line"><span>}</span></span></code></pre><div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0;"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><h2 id="api-types-文件分组分文件夹" tabindex="-1"><a class="header-anchor" href="#api-types-文件分组分文件夹"><span>api types 文件分组分文件夹</span></a></h2><div class="hint-container tip"><p class="hint-container-title">提示</p><p>保证 .jzero.yaml 文件中的 gen.split-api-types-dir 配置为 true, 否则不生效</p></div><div class="language-api line-numbers-mode" data-highlighter="shiki" data-ext="api" data-title="api" style="--shiki-light:#24292e;--shiki-dark:#abb2bf;--shiki-light-bg:#fff;--shiki-dark-bg:#282c34;"><pre class="shiki shiki-themes github-light one-dark-pro vp-code"><code><span class="line"><span>syntax = &quot;v1&quot;</span></span>
<span class="line"><span></span></span>
<span class="line"><span>info (</span></span>
<span class="line"><span>	go_package: &quot;version&quot;</span></span>
<span class="line"><span>)</span></span></code></pre><div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0;"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>jzero 脚手架推荐的 api 文件内容如下:</p><p>可以通过如下命令生成改文件:</p><div class="language-shell line-numbers-mode" data-highlighter="shiki" data-ext="shell" data-title="shell" style="--shiki-light:#24292e;--shiki-dark:#abb2bf;--shiki-light-bg:#fff;--shiki-dark-bg:#282c34;"><pre class="shiki shiki-themes github-light one-dark-pro vp-code"><code><span class="line"><span style="--shiki-light:#6F42C1;--shiki-dark:#61AFEF;">jzero</span><span style="--shiki-light:#032F62;--shiki-dark:#98C379;"> ivm</span><span style="--shiki-light:#032F62;--shiki-dark:#98C379;"> add</span><span style="--shiki-light:#032F62;--shiki-dark:#98C379;"> api</span><span style="--shiki-light:#005CC5;--shiki-dark:#D19A66;"> --name</span><span style="--shiki-light:#032F62;--shiki-dark:#98C379;"> user</span></span></code></pre><div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0;"><div class="line-number"></div></div></div><div class="language-api line-numbers-mode" data-highlighter="shiki" data-ext="api" data-title="api" style="--shiki-light:#24292e;--shiki-dark:#abb2bf;--shiki-light-bg:#fff;--shiki-dark-bg:#282c34;"><pre class="shiki shiki-themes github-light one-dark-pro vp-code"><code><span class="line"><span>syntax = &quot;v1&quot;</span></span>
<span class="line"><span></span></span>
<span class="line"><span>info (</span></span>
<span class="line"><span>	go_package: &quot;user&quot;</span></span>
<span class="line"><span>)</span></span>
<span class="line"><span></span></span>
<span class="line"><span>type CreateRequest {}</span></span>
<span class="line"><span></span></span>
<span class="line"><span>type CreateResponse {}</span></span>
<span class="line"><span></span></span>
<span class="line"><span>type ListRequest {}</span></span>
<span class="line"><span></span></span>
<span class="line"><span>type ListResponse {}</span></span>
<span class="line"><span></span></span>
<span class="line"><span>type GetRequest {}</span></span>
<span class="line"><span></span></span>
<span class="line"><span>type GetResponse {}</span></span>
<span class="line"><span></span></span>
<span class="line"><span>type EditRequest {}</span></span>
<span class="line"><span></span></span>
<span class="line"><span>type EditResponse {}</span></span>
<span class="line"><span></span></span>
<span class="line"><span>type DeleteRequest {}</span></span>
<span class="line"><span></span></span>
<span class="line"><span>type DeleteResponse {}</span></span>
<span class="line"><span></span></span>
<span class="line"><span>@server (</span></span>
<span class="line"><span>	prefix: /api/v1</span></span>
<span class="line"><span>	group:  user</span></span>
<span class="line"><span>)</span></span>
<span class="line"><span>service ntls {</span></span>
<span class="line"><span>	@handler CreateHandler</span></span>
<span class="line"><span>	post /user/create (CreateRequest) returns (CreateResponse)</span></span>
<span class="line"><span></span></span>
<span class="line"><span>	@handler ListHandler</span></span>
<span class="line"><span>	get /user/list (ListRequest) returns (ListResponse)</span></span>
<span class="line"><span></span></span>
<span class="line"><span>	@handler GetHandler</span></span>
<span class="line"><span>	get /user/get (GetRequest) returns (GetResponse)</span></span>
<span class="line"><span></span></span>
<span class="line"><span>	@handler EditHandler</span></span>
<span class="line"><span>	post /user/edit (EditRequest) returns (EditResponse)</span></span>
<span class="line"><span></span></span>
<span class="line"><span>	@handler DeleteHandler</span></span>
<span class="line"><span>	get /user/delete (DeleteRequest) returns (DeleteResponse)</span></span>
<span class="line"><span>}</span></span></code></pre><div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0;"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div>`,11)]))}const c=n(p,[["render",l],["__file","api.html.vue"]]),o=JSON.parse('{"path":"/guide/develop/api.html","title":"api 教程","lang":"zh-CN","frontmatter":{"title":"api 教程","icon":"eos-icons:api","star":true,"order":0.2,"category":"开发","tag":["Guide"],"description":"提示 go-zero api 文档 api 字段校验 jzero 集成 https://github.com/go-playground/validator 进行字段校验 api types 文件分组分文件夹 提示 保证 .jzero.yaml 文件中的 gen.split-api-types-dir 配置为 true, 否则不生效 jzero 脚手架...","head":[["meta",{"property":"og:url","content":"https://jzero.jaronnie.com/guide/develop/api.html"}],["meta",{"property":"og:site_name","content":"Jzero Framework"}],["meta",{"property":"og:title","content":"api 教程"}],["meta",{"property":"og:description","content":"提示 go-zero api 文档 api 字段校验 jzero 集成 https://github.com/go-playground/validator 进行字段校验 api types 文件分组分文件夹 提示 保证 .jzero.yaml 文件中的 gen.split-api-types-dir 配置为 true, 否则不生效 jzero 脚手架..."}],["meta",{"property":"og:type","content":"article"}],["meta",{"property":"og:locale","content":"zh-CN"}],["meta",{"property":"og:updated_time","content":"2024-10-17T06:14:19.000Z"}],["meta",{"property":"article:author","content":"jaronnie"}],["meta",{"property":"article:tag","content":"Guide"}],["meta",{"property":"article:modified_time","content":"2024-10-17T06:14:19.000Z"}],["script",{"type":"application/ld+json"},"{\\"@context\\":\\"https://schema.org\\",\\"@type\\":\\"Article\\",\\"headline\\":\\"api 教程\\",\\"image\\":[\\"\\"],\\"dateModified\\":\\"2024-10-17T06:14:19.000Z\\",\\"author\\":[{\\"@type\\":\\"Person\\",\\"name\\":\\"jaronnie\\",\\"url\\":\\"https://github.com/jaronnie\\"}]}"]]},"headers":[{"level":2,"title":"api 字段校验","slug":"api-字段校验","link":"#api-字段校验","children":[]},{"level":2,"title":"api types 文件分组分文件夹","slug":"api-types-文件分组分文件夹","link":"#api-types-文件分组分文件夹","children":[]}],"git":{"createdTime":1724234753000,"updatedTime":1729145659000,"contributors":[{"name":"jaron","email":"jaron@jaronnie.com","commits":1},{"name":"jaronnie","email":"jaron@jaronnie.com","commits":1}]},"readingTime":{"minutes":0.67,"words":202},"filePathRelative":"guide/develop/api.md","localizedDate":"2024年8月21日","excerpt":"<div class=\\"hint-container tip\\">\\n<p class=\\"hint-container-title\\">提示</p>\\n<p><a href=\\"https://go-zero.dev/docs/tutorials\\" target=\\"_blank\\" rel=\\"noopener noreferrer\\">go-zero api 文档</a></p>\\n</div>\\n<h2>api 字段校验</h2>\\n<blockquote>\\n<p>jzero 集成 <a href=\\"https://github.com/go-playground/validator\\" target=\\"_blank\\" rel=\\"noopener noreferrer\\">https://github.com/go-playground/validator</a> 进行字段校验</p>\\n</blockquote>","autoDesc":true}');export{c as comp,o as data};
