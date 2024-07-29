import{_ as r}from"./plugin-vue_export-helper-DlAUqK2U.js";import{c as t,a as c,w as e,b as i,r as o,o as u,d as s,e as n}from"./app-BvxCvYPe.js";const d={},m=i('<div class="hint-container tip"><p class="hint-container-title">jzero 支持多 proto 进行管理 proto(goctl 原生工具不支持).</p><p>jzero 在自动生成代码的时候会自动识别 daemon/desc/proto 下的文件并自动注册到 zrpc 上. jzero 默认支持对 proto 的字段校验, 使用 protoc-gen-validate.</p></div><p>jzero 框架的理念是:</p><ul><li>不同模块分在不同的 proto 文件下. 如一个系统, 凭证模块即 credential.proto, 主机模块即 machine.proto.</li><li>每个 proto 文件可以有多个 service. 对于复杂模块可以使用多个 service.</li><li>应该考虑接口版本控制, 如 v1/hello.proto, v2/hello_v2.proto</li></ul><p>jzero 中 proto 规范:</p><ul><li>proto 文件引用规范: 依据于 go-zero 的 proto 规范， 即 service 的 rpc 方法中入参和出参的 proto 不能是 import 的 proto 文件中的 message, 只能在当前文件</li></ul><h2 id="proto-文件示例" tabindex="-1"><a class="header-anchor" href="#proto-文件示例"><span>proto 文件示例</span></a></h2>',6),v=n("div",{class:"language-protobuf line-numbers-mode","data-highlighter":"shiki","data-ext":"protobuf","data-title":"protobuf",style:{"--shiki-light":"#24292e","--shiki-dark":"#abb2bf","--shiki-light-bg":"#fff","--shiki-dark-bg":"#282c34"}},[n("pre",{class:"shiki shiki-themes github-light one-dark-pro vp-code"},[n("code",null,[n("span",{class:"line"},[n("span",null,'syntax = "proto3";')]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"package credentialpb;")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,'import "google/api/annotations.proto";')]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,'option go_package = "./pb/credentialpb";')]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"message Empty {}")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"message CredentialVersionResponse {")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string version = 1;")]),s(`
`),n("span",{class:"line"},[n("span",null,"}")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"message CreateCredentialRequest {")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string name = 1;")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string type = 2;")]),s(`
`),n("span",{class:"line"},[n("span",null,"}")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"message CreateCredentialResponse {")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string name = 1;")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string type = 2;")]),s(`
`),n("span",{class:"line"},[n("span",null,"}")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"service credential {")]),s(`
`),n("span",{class:"line"},[n("span",null,"  rpc CredentialVersion(Empty) returns(CredentialVersionResponse) {")]),s(`
`),n("span",{class:"line"},[n("span",null,"    option (google.api.http) = {")]),s(`
`),n("span",{class:"line"},[n("span",null,'      get: "/api/v1.0/credential/version"')]),s(`
`),n("span",{class:"line"},[n("span",null,"    };")]),s(`
`),n("span",{class:"line"},[n("span",null,"  };")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"  rpc CreateCredential(CreateCredentialRequest) returns(CreateCredentialResponse) {")]),s(`
`),n("span",{class:"line"},[n("span",null,"    option (google.api.http) = {")]),s(`
`),n("span",{class:"line"},[n("span",null,'      post: "/api/v1.0/credential/create"')]),s(`
`),n("span",{class:"line"},[n("span",null,'      body: "*"')]),s(`
`),n("span",{class:"line"},[n("span",null,"    };")]),s(`
`),n("span",{class:"line"},[n("span",null,"  }")]),s(`
`),n("span",{class:"line"},[n("span",null,"}")])])]),n("div",{class:"line-numbers","aria-hidden":"true",style:{"counter-reset":"line-number 0"}},[n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"})])],-1),b=n("div",{class:"language-protobuf line-numbers-mode","data-highlighter":"shiki","data-ext":"protobuf","data-title":"protobuf",style:{"--shiki-light":"#24292e","--shiki-dark":"#abb2bf","--shiki-light-bg":"#fff","--shiki-dark-bg":"#282c34"}},[n("pre",{class:"shiki shiki-themes github-light one-dark-pro vp-code"},[n("code",null,[n("span",{class:"line"},[n("span",null,'syntax = "proto3";')]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"package machinepb;")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,'import "google/api/annotations.proto";')]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,'option go_package = "./pb/machinepb";')]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"message Empty {}")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"message MachineVersionResponse {")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string version = 1;")]),s(`
`),n("span",{class:"line"},[n("span",null,"}")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"message CreateMachineRequest {")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string name = 1;")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string type = 2;")]),s(`
`),n("span",{class:"line"},[n("span",null,"}")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"message CreateMachineResponse {")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string name = 1;")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string type = 2;")]),s(`
`),n("span",{class:"line"},[n("span",null,"}")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"service credential {")]),s(`
`),n("span",{class:"line"},[n("span",null,"  rpc MachineVersion(Empty) returns(MachineVersionResponse) {")]),s(`
`),n("span",{class:"line"},[n("span",null,"    option (google.api.http) = {")]),s(`
`),n("span",{class:"line"},[n("span",null,'      get: "/api/v1.0/machine/version"')]),s(`
`),n("span",{class:"line"},[n("span",null,"    };")]),s(`
`),n("span",{class:"line"},[n("span",null,"  };")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"  rpc CreateMachine(CreateMachineRequest) returns(CreateMachineResponse) {")]),s(`
`),n("span",{class:"line"},[n("span",null,"    option (google.api.http) = {")]),s(`
`),n("span",{class:"line"},[n("span",null,'      post: "/api/v1.0/machine/create"')]),s(`
`),n("span",{class:"line"},[n("span",null,'      body: "*"')]),s(`
`),n("span",{class:"line"},[n("span",null,"    };")]),s(`
`),n("span",{class:"line"},[n("span",null,"  }")]),s(`
`),n("span",{class:"line"},[n("span",null,"}")])])]),n("div",{class:"line-numbers","aria-hidden":"true",style:{"counter-reset":"line-number 0"}},[n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"})])],-1),h=n("div",{class:"language-protobuf line-numbers-mode","data-highlighter":"shiki","data-ext":"protobuf","data-title":"protobuf",style:{"--shiki-light":"#24292e","--shiki-dark":"#abb2bf","--shiki-light-bg":"#fff","--shiki-dark-bg":"#282c34"}},[n("pre",{class:"shiki shiki-themes github-light one-dark-pro vp-code"},[n("code",null,[n("span",{class:"line"},[n("span",null,'syntax = "proto3";')]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"package chainpb;")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,'import "google/api/annotations.proto";')]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,'option go_package = "./pb/chainpb";')]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"message Empty {}")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"message CreateNodeRequest {")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string name = 1;")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string type = 2;")]),s(`
`),n("span",{class:"line"},[n("span",null,"}")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"message CreateNodeResponse {")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string name = 1;")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string type = 2;")]),s(`
`),n("span",{class:"line"},[n("span",null,"}")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"message CreateNamespaceRequest {")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string name = 1;")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string type = 2;")]),s(`
`),n("span",{class:"line"},[n("span",null,"}")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"message CreateNamespaceResponse {")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string name = 1;")]),s(`
`),n("span",{class:"line"},[n("span",null,"  string type = 2;")]),s(`
`),n("span",{class:"line"},[n("span",null,"}")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"service node {")]),s(`
`),n("span",{class:"line"},[n("span",null,"  rpc CreateNode(CreateNodeRequest) returns(CreateNodeResponse) {")]),s(`
`),n("span",{class:"line"},[n("span",null,"    option (google.api.http) = {")]),s(`
`),n("span",{class:"line"},[n("span",null,'      post: "/api/v1.0/chain/node/create"')]),s(`
`),n("span",{class:"line"},[n("span",null,'      body: "*"')]),s(`
`),n("span",{class:"line"},[n("span",null,"    };")]),s(`
`),n("span",{class:"line"},[n("span",null,"  }")]),s(`
`),n("span",{class:"line"},[n("span",null,"}")]),s(`
`),n("span",{class:"line"},[n("span")]),s(`
`),n("span",{class:"line"},[n("span",null,"service namespace {")]),s(`
`),n("span",{class:"line"},[n("span",null,"  rpc CreateNamespace(CreateNamespaceRequest) returns(CreateNamespaceResponse) {")]),s(`
`),n("span",{class:"line"},[n("span",null,"    option (google.api.http) = {")]),s(`
`),n("span",{class:"line"},[n("span",null,'      post: "/api/v1.0/chain/namespace/create"')]),s(`
`),n("span",{class:"line"},[n("span",null,'      body: "*"')]),s(`
`),n("span",{class:"line"},[n("span",null,"    };")]),s(`
`),n("span",{class:"line"},[n("span",null,"  }")]),s(`
`),n("span",{class:"line"},[n("span",null,"}")])])]),n("div",{class:"line-numbers","aria-hidden":"true",style:{"counter-reset":"line-number 0"}},[n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"})])],-1),g=i(`<h2 id="proto-字段校验" tabindex="-1"><a class="header-anchor" href="#proto-字段校验"><span>proto 字段校验</span></a></h2><p>基于 <a href="https://github.com/bufbuild/protoc-gen-validate" target="_blank" rel="noopener noreferrer">protoc-gen-validate</a> 进行二次开发, 支持了自定义错误信息</p><p>插件地址: <a href="https://github.com/jzero-io/protoc-gen-validate" target="_blank" rel="noopener noreferrer">protoc-gen-validate</a></p><div class="language-shell line-numbers-mode" data-highlighter="shiki" data-ext="shell" data-title="shell" style="--shiki-light:#24292e;--shiki-dark:#abb2bf;--shiki-light-bg:#fff;--shiki-dark-bg:#282c34;"><pre class="shiki shiki-themes github-light one-dark-pro vp-code"><code><span class="line"><span style="--shiki-light:#6F42C1;--shiki-dark:#61AFEF;">go</span><span style="--shiki-light:#032F62;--shiki-dark:#98C379;"> install</span><span style="--shiki-light:#032F62;--shiki-dark:#98C379;"> github.com/jzero-io/protoc-gen-validate@latest</span></span></code></pre><div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0;"><div class="line-number"></div></div></div><p><a href="https://github.com/jzero-io/protoc-gen-validate/blob/main/validate/validate.proto" target="_blank" rel="noopener noreferrer">确保 validate.proto 文件内容</a></p><div class="hint-container tip"><p class="hint-container-title">提示</p><p>需要自定义错误信息时, 在原始校验规则后加上 _message 即可</p></div><div class="language-protobuf line-numbers-mode" data-highlighter="shiki" data-ext="protobuf" data-title="protobuf" style="--shiki-light:#24292e;--shiki-dark:#abb2bf;--shiki-light-bg:#fff;--shiki-dark-bg:#282c34;"><pre class="shiki shiki-themes github-light one-dark-pro vp-code"><code><span class="line"><span>syntax = &quot;proto3&quot;;</span></span>
<span class="line"><span></span></span>
<span class="line"><span>package hellopb;</span></span>
<span class="line"><span></span></span>
<span class="line"><span>import &quot;validate/validate.proto&quot;;</span></span>
<span class="line"><span>option go_package = &quot;./pb/hellopb&quot;;</span></span>
<span class="line"><span></span></span>
<span class="line"><span>message SayHelloRequest {</span></span>
<span class="line"><span>  string message = 1 [(validate.rules).string = {</span></span>
<span class="line"><span>    max_len: 10,</span></span>
<span class="line"><span>    max_len_message: &quot;最大长度为 10&quot;  // 自定义错误信息</span></span>
<span class="line"><span>  }];</span></span>
<span class="line"><span>}</span></span></code></pre><div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0;"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div>`,7);function k(f,y){const p=o("CodeTabs");return u(),t("div",null,[m,c(p,{id:"38",data:[{id:"credential.proto"},{id:"machine.proto"},{id:"chain.proto(最复杂场景 proto 多 service)"}]},{title0:e(({value:l,isActive:a})=>[s("credential.proto")]),title1:e(({value:l,isActive:a})=>[s("machine.proto")]),title2:e(({value:l,isActive:a})=>[s("chain.proto(最复杂场景 proto 多 service)")]),tab0:e(({value:l,isActive:a})=>[v]),tab1:e(({value:l,isActive:a})=>[b]),tab2:e(({value:l,isActive:a})=>[h]),_:1}),g])}const z=r(d,[["render",k],["__file","proto.html.vue"]]),j=JSON.parse('{"path":"/guide/develop/proto.html","title":"proto 规范","lang":"zh-CN","frontmatter":{"title":"proto 规范","icon":"vscode-icons:file-type-protobuf","star":true,"order":1,"category":"开发","tag":["Guide"],"description":"jzero 支持多 proto 进行管理 proto(goctl 原生工具不支持). jzero 在自动生成代码的时候会自动识别 daemon/desc/proto 下的文件并自动注册到 zrpc 上. jzero 默认支持对 proto 的字段校验, 使用 protoc-gen-validate. jzero 框架的理念是: 不同模块分在不同的 pr...","head":[["meta",{"property":"og:url","content":"https://jzero.jaronnie.com/guide/develop/proto.html"}],["meta",{"property":"og:site_name","content":"Jzero Framework"}],["meta",{"property":"og:title","content":"proto 规范"}],["meta",{"property":"og:description","content":"jzero 支持多 proto 进行管理 proto(goctl 原生工具不支持). jzero 在自动生成代码的时候会自动识别 daemon/desc/proto 下的文件并自动注册到 zrpc 上. jzero 默认支持对 proto 的字段校验, 使用 protoc-gen-validate. jzero 框架的理念是: 不同模块分在不同的 pr..."}],["meta",{"property":"og:type","content":"article"}],["meta",{"property":"og:locale","content":"zh-CN"}],["meta",{"property":"og:updated_time","content":"2024-07-15T06:04:30.000Z"}],["meta",{"property":"article:author","content":"jaronnie"}],["meta",{"property":"article:tag","content":"Guide"}],["meta",{"property":"article:modified_time","content":"2024-07-15T06:04:30.000Z"}],["script",{"type":"application/ld+json"},"{\\"@context\\":\\"https://schema.org\\",\\"@type\\":\\"Article\\",\\"headline\\":\\"proto 规范\\",\\"image\\":[\\"\\"],\\"dateModified\\":\\"2024-07-15T06:04:30.000Z\\",\\"author\\":[{\\"@type\\":\\"Person\\",\\"name\\":\\"jaronnie\\",\\"url\\":\\"https://github.com/jaronnie\\"}]}"]]},"headers":[{"level":2,"title":"proto 文件示例","slug":"proto-文件示例","link":"#proto-文件示例","children":[]},{"level":2,"title":"proto 字段校验","slug":"proto-字段校验","link":"#proto-字段校验","children":[]}],"git":{"createdTime":1713332628000,"updatedTime":1721023470000,"contributors":[{"name":"jaronnie","email":"jaron@jaronnie.com","commits":8}]},"readingTime":{"minutes":1.77,"words":530},"filePathRelative":"guide/develop/proto.md","localizedDate":"2024年4月17日","excerpt":"<div class=\\"hint-container tip\\">\\n<p class=\\"hint-container-title\\">jzero 支持多 proto 进行管理 proto(goctl 原生工具不支持).</p>\\n<p>jzero 在自动生成代码的时候会自动识别 daemon/desc/proto 下的文件并自动注册到 zrpc 上.\\njzero 默认支持对 proto 的字段校验, 使用 protoc-gen-validate.</p>\\n</div>\\n<p>jzero 框架的理念是:</p>\\n<ul>\\n<li>不同模块分在不同的 proto 文件下. 如一个系统, 凭证模块即 credential.proto, 主机模块即 machine.proto.</li>\\n<li>每个 proto 文件可以有多个 service. 对于复杂模块可以使用多个 service.</li>\\n<li>应该考虑接口版本控制, 如 v1/hello.proto, v2/hello_v2.proto</li>\\n</ul>","autoDesc":true}');export{z as comp,j as data};
