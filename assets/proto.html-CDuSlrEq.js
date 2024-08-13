import{_ as t}from"./plugin-vue_export-helper-DlAUqK2U.js";import{c as r,a as o,w as e,b as l,r as c,o as d,d as n,e as s}from"./app-B7_UXdqy.js";const u={},h=l('<div class="hint-container tip"><p class="hint-container-title">jzero 支持多 proto 进行管理 proto(goctl 原生工具不支持).</p><p>jzero 在自动生成代码的时候会自动识别 daemon/desc/proto 下的文件并自动注册到 zrpc 上. jzero 默认支持对 proto 的字段校验, 使用 protoc-gen-validate.</p></div><p>jzero 框架的理念是:</p><ul><li>不同模块分在不同的 proto 文件下. 如一个系统, 凭证模块即 credential.proto, 主机模块即 machine.proto.</li><li>每个 proto 文件可以有多个 service. 对于复杂模块可以使用多个 service.</li><li>应该考虑接口版本控制, 如 v1/hello.proto, v2/hello_v2.proto</li></ul><p>jzero 中 proto 规范:</p><ul><li>proto 文件引用规范: 依据于 go-zero 的 proto 规范， 即 service 的 rpc 方法中入参和出参的 proto 不能是 import 的 proto 文件中的 message, 只能在当前文件</li></ul><h2 id="proto-文件示例" tabindex="-1"><a class="header-anchor" href="#proto-文件示例"><span>proto 文件示例</span></a></h2>',6),m=s("div",{class:"language-protobuf line-numbers-mode","data-highlighter":"shiki","data-ext":"protobuf","data-title":"protobuf",style:{"--shiki-light":"#24292e","--shiki-dark":"#abb2bf","--shiki-light-bg":"#fff","--shiki-dark-bg":"#282c34"}},[s("pre",{class:"shiki shiki-themes github-light one-dark-pro vp-code"},[s("code",null,[s("span",{class:"line"},[s("span",null,'syntax = "proto3";')]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"package credentialpb;")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,'import "google/api/annotations.proto";')]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,'option go_package = "./pb/credentialpb";')]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"message Empty {}")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"message CredentialVersionResponse {")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string version = 1;")]),n(`
`),s("span",{class:"line"},[s("span",null,"}")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"message CreateCredentialRequest {")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string name = 1;")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string type = 2;")]),n(`
`),s("span",{class:"line"},[s("span",null,"}")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"message CreateCredentialResponse {")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string name = 1;")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string type = 2;")]),n(`
`),s("span",{class:"line"},[s("span",null,"}")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"service credential {")]),n(`
`),s("span",{class:"line"},[s("span",null,"  rpc CredentialVersion(Empty) returns(CredentialVersionResponse) {")]),n(`
`),s("span",{class:"line"},[s("span",null,"    option (google.api.http) = {")]),n(`
`),s("span",{class:"line"},[s("span",null,'      get: "/api/v1.0/credential/version"')]),n(`
`),s("span",{class:"line"},[s("span",null,"    };")]),n(`
`),s("span",{class:"line"},[s("span",null,"  };")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"  rpc CreateCredential(CreateCredentialRequest) returns(CreateCredentialResponse) {")]),n(`
`),s("span",{class:"line"},[s("span",null,"    option (google.api.http) = {")]),n(`
`),s("span",{class:"line"},[s("span",null,'      post: "/api/v1.0/credential/create"')]),n(`
`),s("span",{class:"line"},[s("span",null,'      body: "*"')]),n(`
`),s("span",{class:"line"},[s("span",null,"    };")]),n(`
`),s("span",{class:"line"},[s("span",null,"  }")]),n(`
`),s("span",{class:"line"},[s("span",null,"}")])])]),s("div",{class:"line-numbers","aria-hidden":"true",style:{"counter-reset":"line-number 0"}},[s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"})])],-1),v=s("div",{class:"language-protobuf line-numbers-mode","data-highlighter":"shiki","data-ext":"protobuf","data-title":"protobuf",style:{"--shiki-light":"#24292e","--shiki-dark":"#abb2bf","--shiki-light-bg":"#fff","--shiki-dark-bg":"#282c34"}},[s("pre",{class:"shiki shiki-themes github-light one-dark-pro vp-code"},[s("code",null,[s("span",{class:"line"},[s("span",null,'syntax = "proto3";')]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"package machinepb;")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,'import "google/api/annotations.proto";')]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,'option go_package = "./pb/machinepb";')]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"message Empty {}")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"message MachineVersionResponse {")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string version = 1;")]),n(`
`),s("span",{class:"line"},[s("span",null,"}")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"message CreateMachineRequest {")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string name = 1;")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string type = 2;")]),n(`
`),s("span",{class:"line"},[s("span",null,"}")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"message CreateMachineResponse {")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string name = 1;")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string type = 2;")]),n(`
`),s("span",{class:"line"},[s("span",null,"}")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"service credential {")]),n(`
`),s("span",{class:"line"},[s("span",null,"  rpc MachineVersion(Empty) returns(MachineVersionResponse) {")]),n(`
`),s("span",{class:"line"},[s("span",null,"    option (google.api.http) = {")]),n(`
`),s("span",{class:"line"},[s("span",null,'      get: "/api/v1.0/machine/version"')]),n(`
`),s("span",{class:"line"},[s("span",null,"    };")]),n(`
`),s("span",{class:"line"},[s("span",null,"  };")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"  rpc CreateMachine(CreateMachineRequest) returns(CreateMachineResponse) {")]),n(`
`),s("span",{class:"line"},[s("span",null,"    option (google.api.http) = {")]),n(`
`),s("span",{class:"line"},[s("span",null,'      post: "/api/v1.0/machine/create"')]),n(`
`),s("span",{class:"line"},[s("span",null,'      body: "*"')]),n(`
`),s("span",{class:"line"},[s("span",null,"    };")]),n(`
`),s("span",{class:"line"},[s("span",null,"  }")]),n(`
`),s("span",{class:"line"},[s("span",null,"}")])])]),s("div",{class:"line-numbers","aria-hidden":"true",style:{"counter-reset":"line-number 0"}},[s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"})])],-1),b=s("div",{class:"language-protobuf line-numbers-mode","data-highlighter":"shiki","data-ext":"protobuf","data-title":"protobuf",style:{"--shiki-light":"#24292e","--shiki-dark":"#abb2bf","--shiki-light-bg":"#fff","--shiki-dark-bg":"#282c34"}},[s("pre",{class:"shiki shiki-themes github-light one-dark-pro vp-code"},[s("code",null,[s("span",{class:"line"},[s("span",null,'syntax = "proto3";')]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"package chainpb;")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,'import "google/api/annotations.proto";')]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,'option go_package = "./pb/chainpb";')]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"message Empty {}")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"message CreateNodeRequest {")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string name = 1;")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string type = 2;")]),n(`
`),s("span",{class:"line"},[s("span",null,"}")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"message CreateNodeResponse {")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string name = 1;")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string type = 2;")]),n(`
`),s("span",{class:"line"},[s("span",null,"}")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"message CreateNamespaceRequest {")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string name = 1;")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string type = 2;")]),n(`
`),s("span",{class:"line"},[s("span",null,"}")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"message CreateNamespaceResponse {")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string name = 1;")]),n(`
`),s("span",{class:"line"},[s("span",null,"  string type = 2;")]),n(`
`),s("span",{class:"line"},[s("span",null,"}")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"service node {")]),n(`
`),s("span",{class:"line"},[s("span",null,"  rpc CreateNode(CreateNodeRequest) returns(CreateNodeResponse) {")]),n(`
`),s("span",{class:"line"},[s("span",null,"    option (google.api.http) = {")]),n(`
`),s("span",{class:"line"},[s("span",null,'      post: "/api/v1.0/chain/node/create"')]),n(`
`),s("span",{class:"line"},[s("span",null,'      body: "*"')]),n(`
`),s("span",{class:"line"},[s("span",null,"    };")]),n(`
`),s("span",{class:"line"},[s("span",null,"  }")]),n(`
`),s("span",{class:"line"},[s("span",null,"}")]),n(`
`),s("span",{class:"line"},[s("span")]),n(`
`),s("span",{class:"line"},[s("span",null,"service namespace {")]),n(`
`),s("span",{class:"line"},[s("span",null,"  rpc CreateNamespace(CreateNamespaceRequest) returns(CreateNamespaceResponse) {")]),n(`
`),s("span",{class:"line"},[s("span",null,"    option (google.api.http) = {")]),n(`
`),s("span",{class:"line"},[s("span",null,'      post: "/api/v1.0/chain/namespace/create"')]),n(`
`),s("span",{class:"line"},[s("span",null,'      body: "*"')]),n(`
`),s("span",{class:"line"},[s("span",null,"    };")]),n(`
`),s("span",{class:"line"},[s("span",null,"  }")]),n(`
`),s("span",{class:"line"},[s("span",null,"}")])])]),s("div",{class:"line-numbers","aria-hidden":"true",style:{"counter-reset":"line-number 0"}},[s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"}),s("div",{class:"line-number"})])],-1),g=l(`<h2 id="proto-字段校验" tabindex="-1"><a class="header-anchor" href="#proto-字段校验"><span>proto 字段校验</span></a></h2><p>基于 <a href="https://github.com/bufbuild/protoc-gen-validate" target="_blank" rel="noopener noreferrer">protoc-gen-validate</a> 进行二次开发, 支持了自定义错误信息</p><p>插件地址: <a href="https://github.com/jzero-io/protoc-gen-validate" target="_blank" rel="noopener noreferrer">protoc-gen-validate</a></p><div class="language-shell line-numbers-mode" data-highlighter="shiki" data-ext="shell" data-title="shell" style="--shiki-light:#24292e;--shiki-dark:#abb2bf;--shiki-light-bg:#fff;--shiki-dark-bg:#282c34;"><pre class="shiki shiki-themes github-light one-dark-pro vp-code"><code><span class="line"><span style="--shiki-light:#6F42C1;--shiki-dark:#61AFEF;">go</span><span style="--shiki-light:#032F62;--shiki-dark:#98C379;"> install</span><span style="--shiki-light:#032F62;--shiki-dark:#98C379;"> github.com/jzero-io/protoc-gen-validate@latest</span></span></code></pre><div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0;"><div class="line-number"></div></div></div><p><a href="https://github.com/jzero-io/protoc-gen-validate/blob/main/validate/validate.proto" target="_blank" rel="noopener noreferrer">确保 validate.proto 文件内容</a></p><div class="hint-container tip"><p class="hint-container-title">提示</p><p>需要自定义错误信息时, 在原始校验规则后加上 _message 即可</p></div><div class="language-protobuf line-numbers-mode" data-highlighter="shiki" data-ext="protobuf" data-title="protobuf" style="--shiki-light:#24292e;--shiki-dark:#abb2bf;--shiki-light-bg:#fff;--shiki-dark-bg:#282c34;"><pre class="shiki shiki-themes github-light one-dark-pro vp-code"><code><span class="line"><span>syntax = &quot;proto3&quot;;</span></span>
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
<span class="line"><span>}</span></span></code></pre><div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0;"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><h2 id="proto-扩展" tabindex="-1"><a class="header-anchor" href="#proto-扩展"><span>proto 扩展</span></a></h2><h3 id="middleware-的分组管理" tabindex="-1"><a class="header-anchor" href="#middleware-的分组管理"><span>middleware 的分组管理</span></a></h3><div class="hint-container tip"><p class="hint-container-title">确保存在 desc/proto/jzero/api 文件夹</p><p>如果不存在, 请下载到本地 https://github.com/jzero-io/desc/tree/main/proto/jzero/api</p></div><p>添加 middleware</p><div class="language-protobuf line-numbers-mode" data-highlighter="shiki" data-ext="protobuf" data-title="protobuf" style="--shiki-light:#24292e;--shiki-dark:#abb2bf;--shiki-light-bg:#fff;--shiki-dark-bg:#282c34;"><pre class="shiki shiki-themes github-light one-dark-pro vp-code"><code><span class="line"><span>import &quot;jzero/api/http.proto&quot;;</span></span>
<span class="line"><span>import &quot;jzero/api/zrpc.proto&quot;;</span></span>
<span class="line"><span></span></span>
<span class="line"><span>service User {</span></span>
<span class="line"><span>    option (jzero.api.http_group) = {</span></span>
<span class="line"><span>        middleware: &quot;auth&quot;,</span></span>
<span class="line"><span>    };</span></span>
<span class="line"><span></span></span>
<span class="line"><span>    rpc CreateUser(CreateUserRequest) returns(CreateUserResponse) {</span></span>
<span class="line"><span>        option (google.api.http) = {</span></span>
<span class="line"><span>            post: &quot;/api/v1/user/create&quot;,</span></span>
<span class="line"><span>            body: &quot;*&quot;</span></span>
<span class="line"><span>        };</span></span>
<span class="line"><span>        option (jzero.api.zrpc) = {</span></span>
<span class="line"><span>            middleware: &quot;withValue1&quot;,</span></span>
<span class="line"><span>        };</span></span>
<span class="line"><span>    };</span></span>
<span class="line"><span></span></span>
<span class="line"><span>    rpc ListUser(ListUserRequest) returns(ListUserResponse) {</span></span>
<span class="line"><span>        option (google.api.http) = {</span></span>
<span class="line"><span>            get: &quot;/api/v1/user/{username}/list&quot;,</span></span>
<span class="line"><span>        };</span></span>
<span class="line"><span>    };</span></span>
<span class="line"><span>}</span></span></code></pre><div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0;"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>详细解释:</p><ul><li>option (jzero.api.http_group) 即将该 service 下的所有 method 都新增 http 中间件</li><li>option (google.api.http) 只针对某个 method 新增 http 中间件</li><li>option (jzero.api.zrpc_group) 即将该 service 下的所有 method 都新增 zrpc 中间件</li><li>option (google.api.zrpc) 只针对某个 method 新增 zrpc 中间件</li></ul><p>执行 <code>jzero gen</code> 后将会生成一下文件, 以 auth 为例:</p><ul><li>internal/middleware/authmiddleware.go</li><li>internal/middleware/middleware_gen.go</li></ul><p>修改一下 cmd/server 代码, 将生成的文件注册到 server 中:</p><div class="language-go line-numbers-mode" data-highlighter="shiki" data-ext="go" data-title="go" style="--shiki-light:#24292e;--shiki-dark:#abb2bf;--shiki-light-bg:#fff;--shiki-dark-bg:#282c34;"><pre class="shiki shiki-themes github-light one-dark-pro vp-code"><code><span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E06C75;">zrpc</span><span style="--shiki-light:#D73A49;--shiki-dark:#E5C07B;"> :=</span><span style="--shiki-light:#24292E;--shiki-dark:#E06C75;"> server</span><span style="--shiki-light:#24292E;--shiki-dark:#ABB2BF;">.</span><span style="--shiki-light:#6F42C1;--shiki-dark:#61AFEF;">RegisterZrpc</span><span style="--shiki-light:#24292E;--shiki-dark:#ABB2BF;">(</span><span style="--shiki-light:#24292E;--shiki-dark:#E06C75;">svcCtx</span><span style="--shiki-light:#24292E;--shiki-dark:#ABB2BF;">.</span><span style="--shiki-light:#24292E;--shiki-dark:#E06C75;">Config</span><span style="--shiki-light:#24292E;--shiki-dark:#ABB2BF;">, </span><span style="--shiki-light:#24292E;--shiki-dark:#E06C75;">svcCtx</span><span style="--shiki-light:#24292E;--shiki-dark:#ABB2BF;">)</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E06C75;">gw</span><span style="--shiki-light:#D73A49;--shiki-dark:#E5C07B;"> :=</span><span style="--shiki-light:#24292E;--shiki-dark:#E06C75;"> gateway</span><span style="--shiki-light:#24292E;--shiki-dark:#ABB2BF;">.</span><span style="--shiki-light:#6F42C1;--shiki-dark:#61AFEF;">MustNewServer</span><span style="--shiki-light:#24292E;--shiki-dark:#ABB2BF;">(</span><span style="--shiki-light:#24292E;--shiki-dark:#E06C75;">svcCtx</span><span style="--shiki-light:#24292E;--shiki-dark:#ABB2BF;">.</span><span style="--shiki-light:#24292E;--shiki-dark:#E06C75;">Config</span><span style="--shiki-light:#24292E;--shiki-dark:#ABB2BF;">.</span><span style="--shiki-light:#24292E;--shiki-dark:#E06C75;">Gateway</span><span style="--shiki-light:#24292E;--shiki-dark:#ABB2BF;">.</span><span style="--shiki-light:#24292E;--shiki-dark:#E06C75;">GatewayConf</span><span style="--shiki-light:#24292E;--shiki-dark:#ABB2BF;">, </span><span style="--shiki-light:#24292E;--shiki-dark:#E06C75;">middleware</span><span style="--shiki-light:#24292E;--shiki-dark:#ABB2BF;">.</span><span style="--shiki-light:#6F42C1;--shiki-dark:#61AFEF;">WithHeaderProcessor</span><span style="--shiki-light:#24292E;--shiki-dark:#ABB2BF;">())</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E06C75;">middleware</span><span style="--shiki-light:#24292E;--shiki-dark:#ABB2BF;">.</span><span style="--shiki-light:#6F42C1;--shiki-dark:#61AFEF;">RegisterGen</span><span style="--shiki-light:#24292E;--shiki-dark:#ABB2BF;">(</span><span style="--shiki-light:#24292E;--shiki-dark:#E06C75;">zrpc</span><span style="--shiki-light:#24292E;--shiki-dark:#ABB2BF;">, </span><span style="--shiki-light:#24292E;--shiki-dark:#E06C75;">gw</span><span style="--shiki-light:#24292E;--shiki-dark:#ABB2BF;">)</span></span></code></pre><div class="line-numbers" aria-hidden="true" style="counter-reset:line-number 0;"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div>`,18);function k(y,f){const p=c("CodeTabs");return d(),r("div",null,[h,o(p,{id:"38",data:[{id:"credential.proto"},{id:"machine.proto"},{id:"chain.proto(最复杂场景 proto 多 service)"}]},{title0:e(({value:a,isActive:i})=>[n("credential.proto")]),title1:e(({value:a,isActive:i})=>[n("machine.proto")]),title2:e(({value:a,isActive:i})=>[n("chain.proto(最复杂场景 proto 多 service)")]),tab0:e(({value:a,isActive:i})=>[m]),tab1:e(({value:a,isActive:i})=>[v]),tab2:e(({value:a,isActive:i})=>[b]),_:1}),g])}const B=t(u,[["render",k],["__file","proto.html.vue"]]),z=JSON.parse('{"path":"/guide/develop/proto.html","title":"proto 规范","lang":"zh-CN","frontmatter":{"title":"proto 规范","icon":"vscode-icons:file-type-protobuf","star":true,"order":1,"category":"开发","tag":["Guide"],"description":"jzero 支持多 proto 进行管理 proto(goctl 原生工具不支持). jzero 在自动生成代码的时候会自动识别 daemon/desc/proto 下的文件并自动注册到 zrpc 上. jzero 默认支持对 proto 的字段校验, 使用 protoc-gen-validate. jzero 框架的理念是: 不同模块分在不同的 pr...","head":[["meta",{"property":"og:url","content":"https://jzero.jaronnie.com/guide/develop/proto.html"}],["meta",{"property":"og:site_name","content":"Jzero Framework"}],["meta",{"property":"og:title","content":"proto 规范"}],["meta",{"property":"og:description","content":"jzero 支持多 proto 进行管理 proto(goctl 原生工具不支持). jzero 在自动生成代码的时候会自动识别 daemon/desc/proto 下的文件并自动注册到 zrpc 上. jzero 默认支持对 proto 的字段校验, 使用 protoc-gen-validate. jzero 框架的理念是: 不同模块分在不同的 pr..."}],["meta",{"property":"og:type","content":"article"}],["meta",{"property":"og:locale","content":"zh-CN"}],["meta",{"property":"og:updated_time","content":"2024-08-13T10:43:33.000Z"}],["meta",{"property":"article:author","content":"jaronnie"}],["meta",{"property":"article:tag","content":"Guide"}],["meta",{"property":"article:modified_time","content":"2024-08-13T10:43:33.000Z"}],["script",{"type":"application/ld+json"},"{\\"@context\\":\\"https://schema.org\\",\\"@type\\":\\"Article\\",\\"headline\\":\\"proto 规范\\",\\"image\\":[\\"\\"],\\"dateModified\\":\\"2024-08-13T10:43:33.000Z\\",\\"author\\":[{\\"@type\\":\\"Person\\",\\"name\\":\\"jaronnie\\",\\"url\\":\\"https://github.com/jaronnie\\"}]}"]]},"headers":[{"level":2,"title":"proto 文件示例","slug":"proto-文件示例","link":"#proto-文件示例","children":[]},{"level":2,"title":"proto 字段校验","slug":"proto-字段校验","link":"#proto-字段校验","children":[]},{"level":2,"title":"proto 扩展","slug":"proto-扩展","link":"#proto-扩展","children":[{"level":3,"title":"middleware 的分组管理","slug":"middleware-的分组管理","link":"#middleware-的分组管理","children":[]}]}],"git":{"createdTime":1713332628000,"updatedTime":1723545813000,"contributors":[{"name":"jaronnie","email":"jaron@jaronnie.com","commits":9},{"name":"jaron","email":"jaron@jaronnie.com","commits":1}]},"readingTime":{"minutes":2.39,"words":718},"filePathRelative":"guide/develop/proto.md","localizedDate":"2024年4月17日","excerpt":"<div class=\\"hint-container tip\\">\\n<p class=\\"hint-container-title\\">jzero 支持多 proto 进行管理 proto(goctl 原生工具不支持).</p>\\n<p>jzero 在自动生成代码的时候会自动识别 daemon/desc/proto 下的文件并自动注册到 zrpc 上.\\njzero 默认支持对 proto 的字段校验, 使用 protoc-gen-validate.</p>\\n</div>\\n<p>jzero 框架的理念是:</p>\\n<ul>\\n<li>不同模块分在不同的 proto 文件下. 如一个系统, 凭证模块即 credential.proto, 主机模块即 machine.proto.</li>\\n<li>每个 proto 文件可以有多个 service. 对于复杂模块可以使用多个 service.</li>\\n<li>应该考虑接口版本控制, 如 v1/hello.proto, v2/hello_v2.proto</li>\\n</ul>","autoDesc":true}');export{B as comp,z as data};
