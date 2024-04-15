import{_ as e}from"./plugin-vue_export-helper-DlAUqK2U.js";import{o as n,c as i,a}from"./app-3tSUa4gl.js";const l={},s=a(`<p>jzero gen 根据 daemon/api 和 daemon/proto 文件生成代码. 生成代码的逻辑是调用 goctl 工具完成.</p><p>jzero 会自动检测对应文件夹下的内容, 然后进行自动生成, 使用者无需关心.</p><h2 id="生成代码" tabindex="-1"><a class="header-anchor" href="#生成代码"><span>生成代码</span></a></h2><div class="language-bash line-numbers-mode" data-ext="sh" data-title="sh"><pre class="language-bash"><code><span class="token builtin class-name">cd</span> app1
jzero gen
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div></div></div><h2 id="下载依赖" tabindex="-1"><a class="header-anchor" href="#下载依赖"><span>下载依赖</span></a></h2><div class="language-bash line-numbers-mode" data-ext="sh" data-title="sh"><pre class="language-bash"><code>go mod tidy
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div></div></div><h2 id="运行项目" tabindex="-1"><a class="header-anchor" href="#运行项目"><span>运行项目</span></a></h2><div class="language-bash line-numbers-mode" data-ext="sh" data-title="sh"><pre class="language-bash"><code>go run main.go daemon <span class="token parameter variable">--config</span> config.toml
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div></div></div><h2 id="测试接口" tabindex="-1"><a class="header-anchor" href="#测试接口"><span>测试接口</span></a></h2><div class="language-bash line-numbers-mode" data-ext="sh" data-title="sh"><pre class="language-bash"><code><span class="token comment"># gateway</span>
<span class="token function">curl</span> http://localhost:8001/api/v1.0/credential/version
<span class="token comment"># grpc</span>
grpcurl <span class="token parameter variable">-plaintext</span> localhost:8000 credentialpb.credential/CredentialVersion
<span class="token comment"># api</span>
<span class="token function">curl</span> http://localhost:8001/api/v1/hello/me
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>执行命令后的代码结构为:</p><div class="language-bash line-numbers-mode" data-ext="sh" data-title="sh"><pre class="language-bash"><code>$ tree
<span class="token builtin class-name">.</span>
├── cmd
│   ├── daemon.go
│   └── root.go
├── config.toml
├── daemon
│   ├── api
│   │   ├── app1.api
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
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div>`,12),d=[s];function r(c,o){return n(),i("div",null,d)}const m=e(l,[["render",r],["__file","gen.html.vue"]]),u=JSON.parse('{"path":"/guide/gen.html","title":"生成代码","lang":"zh-CN","frontmatter":{"title":"生成代码","icon":"code","order":4,"description":"jzero gen 根据 daemon/api 和 daemon/proto 文件生成代码. 生成代码的逻辑是调用 goctl 工具完成. jzero 会自动检测对应文件夹下的内容, 然后进行自动生成, 使用者无需关心. 生成代码 下载依赖 运行项目 测试接口 执行命令后的代码结构为:","head":[["meta",{"property":"og:url","content":"https://jaronnie.github.io/jzero/guide/gen.html"}],["meta",{"property":"og:site_name","content":"Jzero Framework"}],["meta",{"property":"og:title","content":"生成代码"}],["meta",{"property":"og:description","content":"jzero gen 根据 daemon/api 和 daemon/proto 文件生成代码. 生成代码的逻辑是调用 goctl 工具完成. jzero 会自动检测对应文件夹下的内容, 然后进行自动生成, 使用者无需关心. 生成代码 下载依赖 运行项目 测试接口 执行命令后的代码结构为:"}],["meta",{"property":"og:type","content":"article"}],["meta",{"property":"og:locale","content":"zh-CN"}],["meta",{"property":"og:updated_time","content":"2024-04-15T06:59:31.000Z"}],["meta",{"property":"article:author","content":"jaronnie"}],["meta",{"property":"article:modified_time","content":"2024-04-15T06:59:31.000Z"}],["script",{"type":"application/ld+json"},"{\\"@context\\":\\"https://schema.org\\",\\"@type\\":\\"Article\\",\\"headline\\":\\"生成代码\\",\\"image\\":[\\"\\"],\\"dateModified\\":\\"2024-04-15T06:59:31.000Z\\",\\"author\\":[{\\"@type\\":\\"Person\\",\\"name\\":\\"jaronnie\\",\\"url\\":\\"https://github.com/jaronnie\\"}]}"]]},"headers":[{"level":2,"title":"生成代码","slug":"生成代码","link":"#生成代码","children":[]},{"level":2,"title":"下载依赖","slug":"下载依赖","link":"#下载依赖","children":[]},{"level":2,"title":"运行项目","slug":"运行项目","link":"#运行项目","children":[]},{"level":2,"title":"测试接口","slug":"测试接口","link":"#测试接口","children":[]}],"git":{"createdTime":1713163888000,"updatedTime":1713164371000,"contributors":[{"name":"jaronnie","email":"jaron@jaronnie.com","commits":2}]},"readingTime":{"minutes":0.67,"words":202},"filePathRelative":"guide/gen.md","localizedDate":"2024年4月15日","autoDesc":true}');export{m as comp,u as data};
