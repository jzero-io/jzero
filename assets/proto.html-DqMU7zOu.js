import{_ as c}from"./plugin-vue_export-helper-DlAUqK2U.js";import{r as l,o as p,c as i,d as u,w as a,e as r,b as s,a as n}from"./app-G3wKh349.js";const k={},d=r('<div class="hint-container tip"><p class="hint-container-title">jzero 支持多 proto 进行管理 proto(goctl 原生工具不支持).</p><p>jzero 在自动生成代码的时候会自动识别 daemon/desc/proto 下的文件并自动注册到 zrpc 上.</p></div><p>jzero 框架的理念是:</p><ul><li>不同模块分在不同的 proto 文件下. 如一个系统, 凭证模块即 credential.proto, 主机模块即 machine.proto.</li><li>每个 proto 文件可以有多个 service. 对于复杂模块可以使用多个 service.</li><li>如需对模块进行版本管理, 应该是 credential.proto, credential_v2.proto 规范.</li></ul><p>proto 规范:</p><ul><li>依据于 go-zero 的 proto 规范. 即 service 的 rpc 方法中 入参和出参的 proto 不能是 import 的 proto 文件</li></ul><p>规范文件实例:</p>',6),m=n("div",{class:"language-protobuf line-numbers-mode","data-ext":"protobuf","data-title":"protobuf"},[n("pre",{class:"language-protobuf"},[n("code",null,[n("span",{class:"token keyword"},"syntax"),s(),n("span",{class:"token operator"},"="),s(),n("span",{class:"token string"},'"proto3"'),n("span",{class:"token punctuation"},";"),s(`

`),n("span",{class:"token keyword"},"package"),s(" credentialpb"),n("span",{class:"token punctuation"},";"),s(`

`),n("span",{class:"token keyword"},"import"),s(),n("span",{class:"token string"},'"google/api/annotations.proto"'),n("span",{class:"token punctuation"},";"),s(`

`),n("span",{class:"token keyword"},"option"),s(" go_package "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token string"},'"./pb/credentialpb"'),n("span",{class:"token punctuation"},";"),s(`

`),n("span",{class:"token keyword"},"message"),s(),n("span",{class:"token class-name"},"Empty"),s(),n("span",{class:"token punctuation"},"{"),n("span",{class:"token punctuation"},"}"),s(`

`),n("span",{class:"token keyword"},"message"),s(),n("span",{class:"token class-name"},"CredentialVersionResponse"),s(),n("span",{class:"token punctuation"},"{"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" version "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"1"),n("span",{class:"token punctuation"},";"),s(`
`),n("span",{class:"token punctuation"},"}"),s(`

`),n("span",{class:"token keyword"},"message"),s(),n("span",{class:"token class-name"},"CreateCredentialRequest"),s(),n("span",{class:"token punctuation"},"{"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" name "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"1"),n("span",{class:"token punctuation"},";"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" type "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"2"),n("span",{class:"token punctuation"},";"),s(`
`),n("span",{class:"token punctuation"},"}"),s(`

`),n("span",{class:"token keyword"},"message"),s(),n("span",{class:"token class-name"},"CreateCredentialResponse"),s(),n("span",{class:"token punctuation"},"{"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" name "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"1"),n("span",{class:"token punctuation"},";"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" type "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"2"),n("span",{class:"token punctuation"},";"),s(`
`),n("span",{class:"token punctuation"},"}"),s(`

`),n("span",{class:"token keyword"},"service"),s(),n("span",{class:"token class-name"},"credential"),s(),n("span",{class:"token punctuation"},"{"),s(`
  `),n("span",{class:"token keyword"},"rpc"),s(),n("span",{class:"token function"},"CredentialVersion"),n("span",{class:"token punctuation"},"("),n("span",{class:"token class-name"},"Empty"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token keyword"},"returns"),n("span",{class:"token punctuation"},"("),n("span",{class:"token class-name"},"CredentialVersionResponse"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token punctuation"},"{"),s(`
    `),n("span",{class:"token keyword"},"option"),s(),n("span",{class:"token punctuation"},"("),s("google"),n("span",{class:"token punctuation"},"."),s("api"),n("span",{class:"token punctuation"},"."),s("http"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token operator"},"="),s(),n("span",{class:"token punctuation"},"{"),s(`
      get`),n("span",{class:"token punctuation"},":"),s(),n("span",{class:"token string"},'"/api/v1.0/credential/version"'),s(`
    `),n("span",{class:"token punctuation"},"}"),n("span",{class:"token punctuation"},";"),s(`
  `),n("span",{class:"token punctuation"},"}"),n("span",{class:"token punctuation"},";"),s(`

  `),n("span",{class:"token keyword"},"rpc"),s(),n("span",{class:"token function"},"CreateCredential"),n("span",{class:"token punctuation"},"("),n("span",{class:"token class-name"},"CreateCredentialRequest"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token keyword"},"returns"),n("span",{class:"token punctuation"},"("),n("span",{class:"token class-name"},"CreateCredentialResponse"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token punctuation"},"{"),s(`
    `),n("span",{class:"token keyword"},"option"),s(),n("span",{class:"token punctuation"},"("),s("google"),n("span",{class:"token punctuation"},"."),s("api"),n("span",{class:"token punctuation"},"."),s("http"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token operator"},"="),s(),n("span",{class:"token punctuation"},"{"),s(`
      post`),n("span",{class:"token punctuation"},":"),s(),n("span",{class:"token string"},'"/api/v1.0/credential/create"'),s(`
      body`),n("span",{class:"token punctuation"},":"),s(),n("span",{class:"token string"},'"*"'),s(`
    `),n("span",{class:"token punctuation"},"}"),n("span",{class:"token punctuation"},";"),s(`
  `),n("span",{class:"token punctuation"},"}"),s(`
`),n("span",{class:"token punctuation"},"}"),s(`
`)])]),n("div",{class:"line-numbers","aria-hidden":"true"},[n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"})])],-1),b=n("div",{class:"language-protobuf line-numbers-mode","data-ext":"protobuf","data-title":"protobuf"},[n("pre",{class:"language-protobuf"},[n("code",null,[n("span",{class:"token keyword"},"syntax"),s(),n("span",{class:"token operator"},"="),s(),n("span",{class:"token string"},'"proto3"'),n("span",{class:"token punctuation"},";"),s(`

`),n("span",{class:"token keyword"},"package"),s(" machinepb"),n("span",{class:"token punctuation"},";"),s(`

`),n("span",{class:"token keyword"},"import"),s(),n("span",{class:"token string"},'"google/api/annotations.proto"'),n("span",{class:"token punctuation"},";"),s(`

`),n("span",{class:"token keyword"},"option"),s(" go_package "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token string"},'"./pb/machinepb"'),n("span",{class:"token punctuation"},";"),s(`

`),n("span",{class:"token keyword"},"message"),s(),n("span",{class:"token class-name"},"Empty"),s(),n("span",{class:"token punctuation"},"{"),n("span",{class:"token punctuation"},"}"),s(`

`),n("span",{class:"token keyword"},"message"),s(),n("span",{class:"token class-name"},"MachineVersionResponse"),s(),n("span",{class:"token punctuation"},"{"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" version "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"1"),n("span",{class:"token punctuation"},";"),s(`
`),n("span",{class:"token punctuation"},"}"),s(`

`),n("span",{class:"token keyword"},"message"),s(),n("span",{class:"token class-name"},"CreateMachineRequest"),s(),n("span",{class:"token punctuation"},"{"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" name "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"1"),n("span",{class:"token punctuation"},";"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" type "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"2"),n("span",{class:"token punctuation"},";"),s(`
`),n("span",{class:"token punctuation"},"}"),s(`

`),n("span",{class:"token keyword"},"message"),s(),n("span",{class:"token class-name"},"CreateMachineResponse"),s(),n("span",{class:"token punctuation"},"{"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" name "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"1"),n("span",{class:"token punctuation"},";"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" type "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"2"),n("span",{class:"token punctuation"},";"),s(`
`),n("span",{class:"token punctuation"},"}"),s(`

`),n("span",{class:"token keyword"},"service"),s(),n("span",{class:"token class-name"},"credential"),s(),n("span",{class:"token punctuation"},"{"),s(`
  `),n("span",{class:"token keyword"},"rpc"),s(),n("span",{class:"token function"},"MachineVersion"),n("span",{class:"token punctuation"},"("),n("span",{class:"token class-name"},"Empty"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token keyword"},"returns"),n("span",{class:"token punctuation"},"("),n("span",{class:"token class-name"},"MachineVersionResponse"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token punctuation"},"{"),s(`
    `),n("span",{class:"token keyword"},"option"),s(),n("span",{class:"token punctuation"},"("),s("google"),n("span",{class:"token punctuation"},"."),s("api"),n("span",{class:"token punctuation"},"."),s("http"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token operator"},"="),s(),n("span",{class:"token punctuation"},"{"),s(`
      get`),n("span",{class:"token punctuation"},":"),s(),n("span",{class:"token string"},'"/api/v1.0/machine/version"'),s(`
    `),n("span",{class:"token punctuation"},"}"),n("span",{class:"token punctuation"},";"),s(`
  `),n("span",{class:"token punctuation"},"}"),n("span",{class:"token punctuation"},";"),s(`

  `),n("span",{class:"token keyword"},"rpc"),s(),n("span",{class:"token function"},"CreateMachine"),n("span",{class:"token punctuation"},"("),n("span",{class:"token class-name"},"CreateMachineRequest"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token keyword"},"returns"),n("span",{class:"token punctuation"},"("),n("span",{class:"token class-name"},"CreateMachineResponse"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token punctuation"},"{"),s(`
    `),n("span",{class:"token keyword"},"option"),s(),n("span",{class:"token punctuation"},"("),s("google"),n("span",{class:"token punctuation"},"."),s("api"),n("span",{class:"token punctuation"},"."),s("http"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token operator"},"="),s(),n("span",{class:"token punctuation"},"{"),s(`
      post`),n("span",{class:"token punctuation"},":"),s(),n("span",{class:"token string"},'"/api/v1.0/machine/create"'),s(`
      body`),n("span",{class:"token punctuation"},":"),s(),n("span",{class:"token string"},'"*"'),s(`
    `),n("span",{class:"token punctuation"},"}"),n("span",{class:"token punctuation"},";"),s(`
  `),n("span",{class:"token punctuation"},"}"),s(`
`),n("span",{class:"token punctuation"},"}"),s(`
`)])]),n("div",{class:"line-numbers","aria-hidden":"true"},[n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"})])],-1),v=n("div",{class:"language-protobuf line-numbers-mode","data-ext":"protobuf","data-title":"protobuf"},[n("pre",{class:"language-protobuf"},[n("code",null,[n("span",{class:"token keyword"},"syntax"),s(),n("span",{class:"token operator"},"="),s(),n("span",{class:"token string"},'"proto3"'),n("span",{class:"token punctuation"},";"),s(`

`),n("span",{class:"token keyword"},"package"),s(" chainpb"),n("span",{class:"token punctuation"},";"),s(`

`),n("span",{class:"token keyword"},"import"),s(),n("span",{class:"token string"},'"google/api/annotations.proto"'),n("span",{class:"token punctuation"},";"),s(`

`),n("span",{class:"token keyword"},"option"),s(" go_package "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token string"},'"./pb/chainpb"'),n("span",{class:"token punctuation"},";"),s(`

`),n("span",{class:"token keyword"},"message"),s(),n("span",{class:"token class-name"},"Empty"),s(),n("span",{class:"token punctuation"},"{"),n("span",{class:"token punctuation"},"}"),s(`

`),n("span",{class:"token keyword"},"message"),s(),n("span",{class:"token class-name"},"CreateNodeRequest"),s(),n("span",{class:"token punctuation"},"{"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" name "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"1"),n("span",{class:"token punctuation"},";"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" type "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"2"),n("span",{class:"token punctuation"},";"),s(`
`),n("span",{class:"token punctuation"},"}"),s(`

`),n("span",{class:"token keyword"},"message"),s(),n("span",{class:"token class-name"},"CreateNodeResponse"),s(),n("span",{class:"token punctuation"},"{"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" name "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"1"),n("span",{class:"token punctuation"},";"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" type "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"2"),n("span",{class:"token punctuation"},";"),s(`
`),n("span",{class:"token punctuation"},"}"),s(`

`),n("span",{class:"token keyword"},"message"),s(),n("span",{class:"token class-name"},"CreateNamespaceRequest"),s(),n("span",{class:"token punctuation"},"{"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" name "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"1"),n("span",{class:"token punctuation"},";"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" type "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"2"),n("span",{class:"token punctuation"},";"),s(`
`),n("span",{class:"token punctuation"},"}"),s(`

`),n("span",{class:"token keyword"},"message"),s(),n("span",{class:"token class-name"},"CreateNamespaceResponse"),s(),n("span",{class:"token punctuation"},"{"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" name "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"1"),n("span",{class:"token punctuation"},";"),s(`
  `),n("span",{class:"token builtin"},"string"),s(" type "),n("span",{class:"token operator"},"="),s(),n("span",{class:"token number"},"2"),n("span",{class:"token punctuation"},";"),s(`
`),n("span",{class:"token punctuation"},"}"),s(`

`),n("span",{class:"token keyword"},"service"),s(),n("span",{class:"token class-name"},"node"),s(),n("span",{class:"token punctuation"},"{"),s(`
  `),n("span",{class:"token keyword"},"rpc"),s(),n("span",{class:"token function"},"CreateNode"),n("span",{class:"token punctuation"},"("),n("span",{class:"token class-name"},"CreateNodeRequest"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token keyword"},"returns"),n("span",{class:"token punctuation"},"("),n("span",{class:"token class-name"},"CreateNodeResponse"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token punctuation"},"{"),s(`
    `),n("span",{class:"token keyword"},"option"),s(),n("span",{class:"token punctuation"},"("),s("google"),n("span",{class:"token punctuation"},"."),s("api"),n("span",{class:"token punctuation"},"."),s("http"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token operator"},"="),s(),n("span",{class:"token punctuation"},"{"),s(`
      post`),n("span",{class:"token punctuation"},":"),s(),n("span",{class:"token string"},'"/api/v1.0/chain/node/create"'),s(`
      body`),n("span",{class:"token punctuation"},":"),s(),n("span",{class:"token string"},'"*"'),s(`
    `),n("span",{class:"token punctuation"},"}"),n("span",{class:"token punctuation"},";"),s(`
  `),n("span",{class:"token punctuation"},"}"),s(`
`),n("span",{class:"token punctuation"},"}"),s(`

`),n("span",{class:"token keyword"},"service"),s(),n("span",{class:"token class-name"},"namespace"),s(),n("span",{class:"token punctuation"},"{"),s(`
  `),n("span",{class:"token keyword"},"rpc"),s(),n("span",{class:"token function"},"CreateNamespace"),n("span",{class:"token punctuation"},"("),n("span",{class:"token class-name"},"CreateNamespaceRequest"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token keyword"},"returns"),n("span",{class:"token punctuation"},"("),n("span",{class:"token class-name"},"CreateNamespaceResponse"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token punctuation"},"{"),s(`
    `),n("span",{class:"token keyword"},"option"),s(),n("span",{class:"token punctuation"},"("),s("google"),n("span",{class:"token punctuation"},"."),s("api"),n("span",{class:"token punctuation"},"."),s("http"),n("span",{class:"token punctuation"},")"),s(),n("span",{class:"token operator"},"="),s(),n("span",{class:"token punctuation"},"{"),s(`
      post`),n("span",{class:"token punctuation"},":"),s(),n("span",{class:"token string"},'"/api/v1.0/chain/namespace/create"'),s(`
      body`),n("span",{class:"token punctuation"},":"),s(),n("span",{class:"token string"},'"*"'),s(`
    `),n("span",{class:"token punctuation"},"}"),n("span",{class:"token punctuation"},";"),s(`
  `),n("span",{class:"token punctuation"},"}"),s(`
`),n("span",{class:"token punctuation"},"}"),s(`
`)])]),n("div",{class:"line-numbers","aria-hidden":"true"},[n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"}),n("div",{class:"line-number"})])],-1);function g(y,h){const o=l("CodeTabs");return p(),i("div",null,[d,u(o,{id:"38",data:[{id:"credential.proto"},{id:"machine.proto"},{id:"chain.proto(最复杂场景 proto 多 service)"}]},{title0:a(({value:e,isActive:t})=>[s("credential.proto")]),title1:a(({value:e,isActive:t})=>[s("machine.proto")]),title2:a(({value:e,isActive:t})=>[s("chain.proto(最复杂场景 proto 多 service)")]),tab0:a(({value:e,isActive:t})=>[m]),tab1:a(({value:e,isActive:t})=>[b]),tab2:a(({value:e,isActive:t})=>[v]),_:1})])}const _=c(k,[["render",g],["__file","proto.html.vue"]]),f=JSON.parse('{"path":"/guide/develop/proto.html","title":"proto 规范","lang":"zh-CN","frontmatter":{"title":"proto 规范","icon":"puzzle-piece","star":true,"order":1,"category":"开发","tag":["Guide"],"description":"jzero 支持多 proto 进行管理 proto(goctl 原生工具不支持). jzero 在自动生成代码的时候会自动识别 daemon/desc/proto 下的文件并自动注册到 zrpc 上. jzero 框架的理念是: 不同模块分在不同的 proto 文件下. 如一个系统, 凭证模块即 credential.proto, 主机模块即 mac...","head":[["meta",{"property":"og:url","content":"https://jzero.jaronnie.com/guide/develop/proto.html"}],["meta",{"property":"og:site_name","content":"Jzero Framework"}],["meta",{"property":"og:title","content":"proto 规范"}],["meta",{"property":"og:description","content":"jzero 支持多 proto 进行管理 proto(goctl 原生工具不支持). jzero 在自动生成代码的时候会自动识别 daemon/desc/proto 下的文件并自动注册到 zrpc 上. jzero 框架的理念是: 不同模块分在不同的 proto 文件下. 如一个系统, 凭证模块即 credential.proto, 主机模块即 mac..."}],["meta",{"property":"og:type","content":"article"}],["meta",{"property":"og:locale","content":"zh-CN"}],["meta",{"property":"og:updated_time","content":"2024-04-23T05:45:21.000Z"}],["meta",{"property":"article:author","content":"jaronnie"}],["meta",{"property":"article:tag","content":"Guide"}],["meta",{"property":"article:modified_time","content":"2024-04-23T05:45:21.000Z"}],["script",{"type":"application/ld+json"},"{\\"@context\\":\\"https://schema.org\\",\\"@type\\":\\"Article\\",\\"headline\\":\\"proto 规范\\",\\"image\\":[\\"\\"],\\"dateModified\\":\\"2024-04-23T05:45:21.000Z\\",\\"author\\":[{\\"@type\\":\\"Person\\",\\"name\\":\\"jaronnie\\",\\"url\\":\\"https://github.com/jaronnie\\"}]}"]]},"headers":[],"git":{"createdTime":1713332628000,"updatedTime":1713851121000,"contributors":[{"name":"jaronnie","email":"jaron@jaronnie.com","commits":4}]},"readingTime":{"minutes":1.27,"words":380},"filePathRelative":"guide/develop/proto.md","localizedDate":"2024年4月17日","autoDesc":true}');export{_ as comp,f as data};
