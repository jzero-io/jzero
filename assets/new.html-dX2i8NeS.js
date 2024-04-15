import{_ as e}from"./plugin-vue_export-helper-DlAUqK2U.js";import{o as n,c as a,a as i}from"./app-BMdoMXZF.js";const t={},s=i(`<div class="language-bash line-numbers-mode" data-ext="sh" data-title="sh"><pre class="language-bash"><code>jzero new <span class="token parameter variable">--module</span><span class="token operator">=</span>github.com/jaronnie/app1 <span class="token parameter variable">--dir</span><span class="token operator">=</span>./app1 <span class="token parameter variable">--app</span><span class="token operator">=</span>app1
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div></div></div><p>flag 解释:</p><ul><li>module 表示新建项目的 go module</li><li>dir 表示创建的项目目录路径</li><li>app 表示项目名</li></ul><p>生成的代码结构:</p><div class="language-bash line-numbers-mode" data-ext="sh" data-title="sh"><pre class="language-bash"><code>$ tree                           
<span class="token builtin class-name">.</span>
├── cmd
│   ├── daemon.go
│   └── root.go
├── config.toml
├── daemon
│   ├── api
│   │   ├── app5.api
│   │   ├── file.api
│   │   └── hello.api
│   ├── daemon.go
│   ├── internal
│   │   ├── config
│   │   │   └── config.go
│   │   └── handler
│   │       ├── myhandler.go
│   │       └── myroutes.go
│   └── proto
│       ├── credential.proto
│       └── machine.proto
├── go.mod
└── main.go

<span class="token number">8</span> directories, <span class="token number">14</span> files
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div>`,5),o=[s];function r(l,d){return n(),a("div",null,o)}const p=e(t,[["render",r],["__file","new.html.vue"]]),u=JSON.parse('{"path":"/guide/new.html","title":"新建项目","lang":"zh-CN","frontmatter":{"title":"新建项目","icon":"clone","order":3,"description":"flag 解释: module 表示新建项目的 go module dir 表示创建的项目目录路径 app 表示项目名 生成的代码结构:","head":[["meta",{"property":"og:url","content":"https://jaronnie.github.io/jzero/guide/new.html"}],["meta",{"property":"og:site_name","content":"Jzero Framework"}],["meta",{"property":"og:title","content":"新建项目"}],["meta",{"property":"og:description","content":"flag 解释: module 表示新建项目的 go module dir 表示创建的项目目录路径 app 表示项目名 生成的代码结构:"}],["meta",{"property":"og:type","content":"article"}],["meta",{"property":"og:locale","content":"zh-CN"}],["meta",{"property":"og:updated_time","content":"2024-04-15T06:51:28.000Z"}],["meta",{"property":"article:author","content":"jaronnie"}],["meta",{"property":"article:modified_time","content":"2024-04-15T06:51:28.000Z"}],["script",{"type":"application/ld+json"},"{\\"@context\\":\\"https://schema.org\\",\\"@type\\":\\"Article\\",\\"headline\\":\\"新建项目\\",\\"image\\":[\\"\\"],\\"dateModified\\":\\"2024-04-15T06:51:28.000Z\\",\\"author\\":[{\\"@type\\":\\"Person\\",\\"name\\":\\"jaronnie\\",\\"url\\":\\"https://github.com/jaronnie\\"}]}"]]},"headers":[],"git":{"createdTime":1713163888000,"updatedTime":1713163888000,"contributors":[{"name":"jaronnie","email":"jaron@jaronnie.com","commits":1}]},"readingTime":{"minutes":0.28,"words":84},"filePathRelative":"guide/new.md","localizedDate":"2024年4月15日","autoDesc":true}');export{p as comp,u as data};
