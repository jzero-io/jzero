import{_ as e}from"./plugin-vue_export-helper-DlAUqK2U.js";import{o as n,c as a,a as t}from"./app-BEDc3YmB.js";const i={},s=t(`<h2 id="安装" tabindex="-1"><a class="header-anchor" href="#安装"><span>安装</span></a></h2><div class="language-bash line-numbers-mode" data-ext="sh" data-title="sh"><pre class="language-bash"><code>go <span class="token function">install</span> github.com/jaronnie/jzero@latest
<span class="token comment"># 初始化</span>
jzero init
<span class="token comment"># 启动服务</span>
jzero jzerod
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><h2 id="开发" tabindex="-1"><a class="header-anchor" href="#开发"><span>开发</span></a></h2><div class="language-bash line-numbers-mode" data-ext="sh" data-title="sh"><pre class="language-bash"><code><span class="token comment"># gencode</span>
task gencode

<span class="token comment"># run</span>
task run

<span class="token comment"># test</span>
<span class="token comment"># unix</span>
<span class="token function">curl</span> <span class="token parameter variable">-s</span> --unix-socket ./jzero.sock http://localhost:8001/api/v1.0/credential/version
<span class="token comment"># gateway</span>
http://localhost:8001/api/v1.0/credential/version
<span class="token comment"># grpc</span>
grpcurl <span class="token parameter variable">-plaintext</span> localhost:8000 credentialpb.credential/CredentialVersion
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div>`,4),r=[s];function o(l,c){return n(),a("div",null,r)}const m=e(i,[["render",o],["__file","index.html.vue"]]),u=JSON.parse('{"path":"/guide/","title":"指南","lang":"zh-CN","frontmatter":{"title":"指南","icon":"lightbulb","description":"安装 开发","head":[["meta",{"property":"og:url","content":"https://jaronnie.github.io/jzero/guide/"}],["meta",{"property":"og:site_name","content":"Jzero Framework"}],["meta",{"property":"og:title","content":"指南"}],["meta",{"property":"og:description","content":"安装 开发"}],["meta",{"property":"og:type","content":"article"}],["meta",{"property":"og:locale","content":"zh-CN"}],["meta",{"property":"og:updated_time","content":"2024-04-11T08:57:13.000Z"}],["meta",{"property":"article:author","content":"jaronnie"}],["meta",{"property":"article:modified_time","content":"2024-04-11T08:57:13.000Z"}],["script",{"type":"application/ld+json"},"{\\"@context\\":\\"https://schema.org\\",\\"@type\\":\\"Article\\",\\"headline\\":\\"指南\\",\\"image\\":[\\"\\"],\\"dateModified\\":\\"2024-04-11T08:57:13.000Z\\",\\"author\\":[{\\"@type\\":\\"Person\\",\\"name\\":\\"jaronnie\\",\\"url\\":\\"https://github.com/jaronnie\\"}]}"]]},"headers":[{"level":2,"title":"安装","slug":"安装","link":"#安装","children":[]},{"level":2,"title":"开发","slug":"开发","link":"#开发","children":[]}],"git":{"createdTime":1712825833000,"updatedTime":1712825833000,"contributors":[{"name":"jaronnie","email":"jaron@jaronnie.com","commits":1}]},"readingTime":{"minutes":0.17,"words":51},"filePathRelative":"guide/README.md","localizedDate":"2024年4月11日","autoDesc":true}');export{m as comp,u as data};
