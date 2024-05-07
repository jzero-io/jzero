import{_ as n}from"./plugin-vue_export-helper-DlAUqK2U.js";import{o as s,c as a,e}from"./app-mSjmIJ7g.js";const p={},t=e(`<p>修改 config.yaml, 增加以下配置, 设置最大 qps 100</p><div class="language-yaml line-numbers-mode" data-ext="yml" data-title="yml"><pre class="language-yaml"><code><span class="token key atrule">Gateway</span><span class="token punctuation">:</span>
  <span class="token key atrule">MaxConns</span><span class="token punctuation">:</span> <span class="token number">100</span>

<span class="token comment"># 替换成自己的 App 名称</span>
<span class="token key atrule">App1</span><span class="token punctuation">:</span>
  <span class="token key atrule">GrpcMaxConns</span><span class="token punctuation">:</span> <span class="token number">100</span>
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div><p>由于 jzero 集成了 go-zero 三个特性</p><ul><li>rpc</li><li>api</li><li>gateway</li></ul><p>我们依次测试这三种类型的接口</p><div class="hint-container tip"><p class="hint-container-title">提示</p><p>https://github.com/zeromicro/go-zero/issues/4097</p><p>两种路由的限流都没生效</p><ul><li>api 生成的 handler 注册到 gateway server 后</li><li>gateway server AddRoute 的路由</li></ul></div><div class="language-bash line-numbers-mode" data-ext="sh" data-title="sh"><pre class="language-bash"><code><span class="token comment"># test rpc</span>
ghz <span class="token parameter variable">--insecure</span> <span class="token parameter variable">-c</span> <span class="token number">110</span> <span class="token parameter variable">-n</span> <span class="token number">110</span>  <span class="token punctuation">\\</span>
  <span class="token parameter variable">--call</span> credentialpb.credential.CredentialVersion <span class="token punctuation">\\</span>
  <span class="token number">0.0</span>.0.0:8000
  
$ ghz <span class="token parameter variable">--insecure</span> <span class="token parameter variable">-c</span> <span class="token number">110</span> <span class="token parameter variable">-n</span> <span class="token number">110</span>  <span class="token punctuation">\\</span>
  <span class="token parameter variable">--call</span> credentialpb.credential.CredentialVersion <span class="token punctuation">\\</span>
  <span class="token number">0.0</span>.0.0:8000

Summary:
  Count:	<span class="token number">110</span>
  Total:	<span class="token number">117.02</span> ms
  Slowest:	<span class="token number">106.65</span> ms
  Fastest:	<span class="token number">51.92</span> ms
  Average:	<span class="token number">77.79</span> ms
  Requests/sec:	<span class="token number">940.03</span>

Response <span class="token function">time</span> histogram:
  <span class="token number">51.918</span>  <span class="token punctuation">[</span><span class="token number">1</span><span class="token punctuation">]</span>  <span class="token operator">|</span>∎∎
  <span class="token number">57.391</span>  <span class="token punctuation">[</span><span class="token number">4</span><span class="token punctuation">]</span>  <span class="token operator">|</span>∎∎∎∎∎∎∎∎∎
  <span class="token number">62.865</span>  <span class="token punctuation">[</span><span class="token number">3</span><span class="token punctuation">]</span>  <span class="token operator">|</span>∎∎∎∎∎∎∎
  <span class="token number">68.339</span>  <span class="token punctuation">[</span><span class="token number">8</span><span class="token punctuation">]</span>  <span class="token operator">|</span>∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  <span class="token number">73.813</span>  <span class="token punctuation">[</span><span class="token number">9</span><span class="token punctuation">]</span>  <span class="token operator">|</span>∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  <span class="token number">79.286</span>  <span class="token punctuation">[</span><span class="token number">6</span><span class="token punctuation">]</span>  <span class="token operator">|</span>∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  <span class="token number">84.760</span>  <span class="token punctuation">[</span><span class="token number">11</span><span class="token punctuation">]</span> <span class="token operator">|</span>∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  <span class="token number">90.234</span>  <span class="token punctuation">[</span><span class="token number">12</span><span class="token punctuation">]</span> <span class="token operator">|</span>∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  <span class="token number">95.707</span>  <span class="token punctuation">[</span><span class="token number">17</span><span class="token punctuation">]</span> <span class="token operator">|</span>∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  <span class="token number">101.181</span> <span class="token punctuation">[</span><span class="token number">16</span><span class="token punctuation">]</span> <span class="token operator">|</span>∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  <span class="token number">106.655</span> <span class="token punctuation">[</span><span class="token number">13</span><span class="token punctuation">]</span> <span class="token operator">|</span>∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎

Latency distribution:
  <span class="token number">10</span> % <span class="token keyword">in</span> <span class="token number">65.12</span> ms
  <span class="token number">25</span> % <span class="token keyword">in</span> <span class="token number">73.70</span> ms
  <span class="token number">50</span> % <span class="token keyword">in</span> <span class="token number">88.74</span> ms
  <span class="token number">75</span> % <span class="token keyword">in</span> <span class="token number">96.70</span> ms
  <span class="token number">90</span> % <span class="token keyword">in</span> <span class="token number">102.15</span> ms
  <span class="token number">95</span> % <span class="token keyword">in</span> <span class="token number">104.78</span> ms
  <span class="token number">99</span> % <span class="token keyword">in</span> <span class="token number">105.13</span> ms

Status code distribution:
  <span class="token punctuation">[</span>Unavailable<span class="token punctuation">]</span>   <span class="token number">10</span> responses
  <span class="token punctuation">[</span>OK<span class="token punctuation">]</span>            <span class="token number">100</span> responses

Error distribution:
  <span class="token punctuation">[</span><span class="token number">10</span><span class="token punctuation">]</span>   rpc error: code <span class="token operator">=</span> Unavailable desc <span class="token operator">=</span> concurrent connections over limit

<span class="token comment"># test api</span>
hey <span class="token parameter variable">-z</span> 1s <span class="token parameter variable">-c</span> <span class="token number">120</span> <span class="token parameter variable">-q</span> <span class="token number">1</span> <span class="token string">&#39;http://localhost:8001/api/v1/hello/you&#39;</span>

Summary:
  Total:	<span class="token number">1.0821</span> secs
  Slowest:	<span class="token number">0.0745</span> secs
  Fastest:	<span class="token number">0.0196</span> secs
  Average:	<span class="token number">0.0475</span> secs
  Requests/sec:	<span class="token number">110.8997</span>

  Total data:	<span class="token number">8880</span> bytes
  Size/request:	<span class="token number">74</span> bytes

Response <span class="token function">time</span> histogram:
  <span class="token number">0.020</span> <span class="token punctuation">[</span><span class="token number">1</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■
  <span class="token number">0.025</span> <span class="token punctuation">[</span><span class="token number">4</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■■■■■■■
  <span class="token number">0.031</span> <span class="token punctuation">[</span><span class="token number">11</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■■■■■■■■■■■■■■■■■■■■■
  <span class="token number">0.036</span> <span class="token punctuation">[</span><span class="token number">10</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■■■■■■■■■■■■■■■■■■■
  <span class="token number">0.042</span> <span class="token punctuation">[</span><span class="token number">15</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  <span class="token number">0.047</span> <span class="token punctuation">[</span><span class="token number">15</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  <span class="token number">0.053</span> <span class="token punctuation">[</span><span class="token number">19</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  <span class="token number">0.058</span> <span class="token punctuation">[</span><span class="token number">13</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■■■■■■■■■■■■■■■■■■■■■■■■■
  <span class="token number">0.064</span> <span class="token punctuation">[</span><span class="token number">20</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  <span class="token number">0.069</span> <span class="token punctuation">[</span><span class="token number">8</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■■■■■■■■■■■■■■■
  <span class="token number">0.075</span> <span class="token punctuation">[</span><span class="token number">4</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■■■■■■■


Latency distribution:
  <span class="token number">10</span>% <span class="token keyword">in</span> <span class="token number">0.0297</span> secs
  <span class="token number">25</span>% <span class="token keyword">in</span> <span class="token number">0.0379</span> secs
  <span class="token number">50</span>% <span class="token keyword">in</span> <span class="token number">0.0487</span> secs
  <span class="token number">75</span>% <span class="token keyword">in</span> <span class="token number">0.0584</span> secs
  <span class="token number">90</span>% <span class="token keyword">in</span> <span class="token number">0.0635</span> secs
  <span class="token number">95</span>% <span class="token keyword">in</span> <span class="token number">0.0670</span> secs
  <span class="token number">99</span>% <span class="token keyword">in</span> <span class="token number">0.0745</span> secs

Details <span class="token punctuation">(</span>average, fastest, slowest<span class="token punctuation">)</span>:
  DNS+dialup:	<span class="token number">0.0143</span> secs, <span class="token number">0.0196</span> secs, <span class="token number">0.0745</span> secs
  DNS-lookup:	<span class="token number">0.0072</span> secs, <span class="token number">0.0024</span> secs, <span class="token number">0.0108</span> secs
  req write:	<span class="token number">0.0005</span> secs, <span class="token number">0.0000</span> secs, <span class="token number">0.0031</span> secs
  resp wait:	<span class="token number">0.0250</span> secs, <span class="token number">0.0026</span> secs, <span class="token number">0.0457</span> secs
  resp read:	<span class="token number">0.0000</span> secs, <span class="token number">0.0000</span> secs, <span class="token number">0.0010</span> secs

Status code distribution:
  <span class="token punctuation">[</span><span class="token number">200</span><span class="token punctuation">]</span>	<span class="token number">120</span> responses


<span class="token comment"># test gateway</span>
<span class="token comment"># 用 hey 工具来进行压测，压测 120 个并发，执行 1 秒, 有 20 个被限流</span>
hey <span class="token parameter variable">-z</span> 1s <span class="token parameter variable">-c</span> <span class="token number">120</span> <span class="token parameter variable">-q</span> <span class="token number">1</span> <span class="token string">&#39;http://localhost:8001/api/v1/credential/version&#39;</span>

Summary:
  Total:	<span class="token number">1.1574</span> secs
  Slowest:	<span class="token number">0.1511</span> secs
  Fastest:	<span class="token number">0.0217</span> secs
  Average:	<span class="token number">0.1111</span> secs
  Requests/sec:	<span class="token number">103.6849</span>

  Total data:	<span class="token number">5800</span> bytes
  Size/request:	<span class="token number">48</span> bytes

Response <span class="token function">time</span> histogram:
  <span class="token number">0.022</span> <span class="token punctuation">[</span><span class="token number">1</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■
  <span class="token number">0.035</span> <span class="token punctuation">[</span><span class="token number">17</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■■■■■■■■■■■■■■
  <span class="token number">0.048</span> <span class="token punctuation">[</span><span class="token number">2</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■
  <span class="token number">0.061</span> <span class="token punctuation">[</span><span class="token number">0</span><span class="token punctuation">]</span>	<span class="token operator">|</span>
  <span class="token number">0.073</span> <span class="token punctuation">[</span><span class="token number">2</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■
  <span class="token number">0.086</span> <span class="token punctuation">[</span><span class="token number">2</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■
  <span class="token number">0.099</span> <span class="token punctuation">[</span><span class="token number">0</span><span class="token punctuation">]</span>	<span class="token operator">|</span>
  <span class="token number">0.112</span> <span class="token punctuation">[</span><span class="token number">5</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■■■■
  <span class="token number">0.125</span> <span class="token punctuation">[</span><span class="token number">24</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■■■■■■■■■■■■■■■■■■■■■
  <span class="token number">0.138</span> <span class="token punctuation">[</span><span class="token number">44</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  <span class="token number">0.151</span> <span class="token punctuation">[</span><span class="token number">23</span><span class="token punctuation">]</span>	<span class="token operator">|</span>■■■■■■■■■■■■■■■■■■■■■


Latency distribution:
  <span class="token number">10</span>% <span class="token keyword">in</span> <span class="token number">0.0330</span> secs
  <span class="token number">25</span>% <span class="token keyword">in</span> <span class="token number">0.1144</span> secs
  <span class="token number">50</span>% <span class="token keyword">in</span> <span class="token number">0.1264</span> secs
  <span class="token number">75</span>% <span class="token keyword">in</span> <span class="token number">0.1366</span> secs
  <span class="token number">90</span>% <span class="token keyword">in</span> <span class="token number">0.1424</span> secs
  <span class="token number">95</span>% <span class="token keyword">in</span> <span class="token number">0.1438</span> secs
  <span class="token number">99</span>% <span class="token keyword">in</span> <span class="token number">0.1511</span> secs

Details <span class="token punctuation">(</span>average, fastest, slowest<span class="token punctuation">)</span>:
  DNS+dialup:	<span class="token number">0.0128</span> secs, <span class="token number">0.0217</span> secs, <span class="token number">0.1511</span> secs
  DNS-lookup:	<span class="token number">0.0046</span> secs, <span class="token number">0.0009</span> secs, <span class="token number">0.0079</span> secs
  req write:	<span class="token number">0.0004</span> secs, <span class="token number">0.0000</span> secs, <span class="token number">0.0023</span> secs
  resp wait:	<span class="token number">0.0969</span> secs, <span class="token number">0.0056</span> secs, <span class="token number">0.1301</span> secs
  resp read:	<span class="token number">0.0000</span> secs, <span class="token number">0.0000</span> secs, <span class="token number">0.0002</span> secs

Status code distribution:
  <span class="token punctuation">[</span><span class="token number">200</span><span class="token punctuation">]</span>	<span class="token number">100</span> responses
  <span class="token punctuation">[</span><span class="token number">503</span><span class="token punctuation">]</span>	<span class="token number">20</span> responses
</code></pre><div class="line-numbers" aria-hidden="true"><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div><div class="line-number"></div></div></div>`,7),l=[t];function c(i,o){return s(),a("div",null,l)}const m=n(p,[["render",c],["__file","limit.html.vue"]]),d=JSON.parse('{"path":"/guide/config/limit.html","title":"限流配置","lang":"zh-CN","frontmatter":{"title":"限流配置","icon":"gears","star":true,"order":5,"category":"配置","tag":["Guide"],"description":"修改 config.yaml, 增加以下配置, 设置最大 qps 100 由于 jzero 集成了 go-zero 三个特性 rpc api gateway 我们依次测试这三种类型的接口 提示 https://github.com/zeromicro/go-zero/issues/4097 两种路由的限流都没生效 api 生成的 handler 注册到...","head":[["meta",{"property":"og:url","content":"https://jzero.jaronnie.com/guide/config/limit.html"}],["meta",{"property":"og:site_name","content":"Jzero Framework"}],["meta",{"property":"og:title","content":"限流配置"}],["meta",{"property":"og:description","content":"修改 config.yaml, 增加以下配置, 设置最大 qps 100 由于 jzero 集成了 go-zero 三个特性 rpc api gateway 我们依次测试这三种类型的接口 提示 https://github.com/zeromicro/go-zero/issues/4097 两种路由的限流都没生效 api 生成的 handler 注册到..."}],["meta",{"property":"og:type","content":"article"}],["meta",{"property":"og:locale","content":"zh-CN"}],["meta",{"property":"og:updated_time","content":"2024-05-06T10:44:50.000Z"}],["meta",{"property":"article:author","content":"jaronnie"}],["meta",{"property":"article:tag","content":"Guide"}],["meta",{"property":"article:modified_time","content":"2024-05-06T10:44:50.000Z"}],["script",{"type":"application/ld+json"},"{\\"@context\\":\\"https://schema.org\\",\\"@type\\":\\"Article\\",\\"headline\\":\\"限流配置\\",\\"image\\":[\\"\\"],\\"dateModified\\":\\"2024-05-06T10:44:50.000Z\\",\\"author\\":[{\\"@type\\":\\"Person\\",\\"name\\":\\"jaronnie\\",\\"url\\":\\"https://github.com/jaronnie\\"}]}"]]},"headers":[],"git":{"createdTime":1713525118000,"updatedTime":1714992290000,"contributors":[{"name":"jaronnie","email":"jaron@jaronnie.com","commits":4}]},"readingTime":{"minutes":1.77,"words":530},"filePathRelative":"guide/config/limit.md","localizedDate":"2024年4月19日","autoDesc":true}');export{m as comp,d as data};
