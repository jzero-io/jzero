---
title: 版本升级兼容
icon: fluent-mdl2:upgrade-analysis
star: true
order: 2
category: faq
tag:
  - faq
---

## from old to v0.30.0+

add `MustGetConfig` method for `ServiceContext` return `config.Config`

```go
func (sc *ServiceContext) MustGetConfig() config.Config {
	return sc.Config
}
```

