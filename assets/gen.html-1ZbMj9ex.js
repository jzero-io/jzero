import{_ as e}from"./plugin-vue_export-helper-DlAUqK2U.js";import{o as n,c as i,a as l}from"./app-BMdoMXZF.js";const a={},s=l(`<p>jzero gen 根据 daemon/api 和 daemon/proto 文件生成代码. 生成代码的逻辑是调用 goctl 工具完成.</p><p>jzero 会自动检测对应文件夹下的内容, 然后进行自动生成, 使用者无需关心.</p><div class="language-bash line-numbers-mode" data-ext="sh" data-title="sh"><pre class="language-bash"><code>jzero gen
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div></div></div><p>执行命令后的代码结构为:</p><div class="language-bash line-numbers-mode" data-ext="sh" data-title="sh"><pre class="language-bash"><code>$ tree
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
│   │   ├── handler
│   │   │   ├── <span class="token function">file</span>
│   │   │   │   ├── downloadhandler.go
│   │   │   │   └── uploadhandler.go
│   │   │   ├── hello
│   │   │   │   ├── helloparamhandler.go
│   │   │   │   ├── hellopathhandler.go
│   │   │   │   └── helloposthandler.go
│   │   │   ├── myhandler.go
│   │   │   ├── myroutes.go
│   │   │   └── routes.go
│   │   ├── logic
│   │   │   ├── credential
│   │   │   │   └── credentialversionlogic.go
│   │   │   ├── credentialv2
│   │   │   │   └── credentialversionlogic.go
│   │   │   ├── <span class="token function">file</span>
│   │   │   │   ├── downloadlogic.go
│   │   │   │   └── uploadlogic.go
│   │   │   ├── hello
│   │   │   │   ├── helloparamlogic.go
│   │   │   │   ├── hellopathlogic.go
│   │   │   │   └── hellopostlogic.go
│   │   │   ├── machine
│   │   │   │   └── machineversionlogic.go
│   │   │   └── machinev2
│   │   │       └── machineversionlogic.go
│   │   ├── server
│   │   │   ├── credential
│   │   │   │   └── credentialserver.go
│   │   │   ├── credentialv2
│   │   │   │   └── credentialv2server.go
│   │   │   ├── machine
│   │   │   │   └── machineserver.go
│   │   │   └── machinev2
│   │   │       └── machinev2server.go
│   │   ├── svc
│   │   │   └── servicecontext.go
│   │   └── types
│   │       └── types.go
│   ├── pb
│   │   ├── credentialpb
│   │   │   ├── credential.pb.go
│   │   │   └── credential_grpc.pb.go
│   │   └── machinepb
│   │       ├── machine.pb.go
│   │       └── machine_grpc.pb.go
│   └── proto
│       ├── credential.proto
│       └── machine.proto
├── go.mod
└── main.go

<span class="token number">27</span> directories, <span class="token number">39</span> files
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div>`,5),d=[s];function r(o,c){return n(),i("div",null,d)}const m=e(a,[["render",r],["__file","gen.html.vue"]]),u=JSON.parse('{"path":"/guide/gen.html","title":"生成代码","lang":"zh-CN","frontmatter":{"title":"生成代码","icon":"code","order":4,"description":"jzero gen 根据 daemon/api 和 daemon/proto 文件生成代码. 生成代码的逻辑是调用 goctl 工具完成. jzero 会自动检测对应文件夹下的内容, 然后进行自动生成, 使用者无需关心. 执行命令后的代码结构为:","head":[["meta",{"property":"og:url","content":"https://jaronnie.github.io/jzero/guide/gen.html"}],["meta",{"property":"og:site_name","content":"Jzero Framework"}],["meta",{"property":"og:title","content":"生成代码"}],["meta",{"property":"og:description","content":"jzero gen 根据 daemon/api 和 daemon/proto 文件生成代码. 生成代码的逻辑是调用 goctl 工具完成. jzero 会自动检测对应文件夹下的内容, 然后进行自动生成, 使用者无需关心. 执行命令后的代码结构为:"}],["meta",{"property":"og:type","content":"article"}],["meta",{"property":"og:locale","content":"zh-CN"}],["meta",{"property":"og:updated_time","content":"2024-04-15T06:51:28.000Z"}],["meta",{"property":"article:author","content":"jaronnie"}],["meta",{"property":"article:modified_time","content":"2024-04-15T06:51:28.000Z"}],["script",{"type":"application/ld+json"},"{\\"@context\\":\\"https://schema.org\\",\\"@type\\":\\"Article\\",\\"headline\\":\\"生成代码\\",\\"image\\":[\\"\\"],\\"dateModified\\":\\"2024-04-15T06:51:28.000Z\\",\\"author\\":[{\\"@type\\":\\"Person\\",\\"name\\":\\"jaronnie\\",\\"url\\":\\"https://github.com/jaronnie\\"}]}"]]},"headers":[],"git":{"createdTime":1713163888000,"updatedTime":1713163888000,"contributors":[{"name":"jaronnie","email":"jaron@jaronnie.com","commits":1}]},"readingTime":{"minutes":0.52,"words":156},"filePathRelative":"guide/gen.md","localizedDate":"2024年4月15日","autoDesc":true}');export{m as comp,u as data};
