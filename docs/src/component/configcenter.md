---
title: configcenter(配置中心)
icon: catppuccin:astro-config
order: 1
---

## 基于 fsnotify(默认)

* 支持动态加载配置并能设置回调
* 支持环境变量(参考 [envsubst](https://github.com/a8m/envsubst))

### 1. 初始化 configcenter

```go
cc := configcenter.MustNewConfigCenter[config.Config](
	   configcenter.Config{Type: "yaml"}, 
     subscriber.MustNewFsnotifySubscriber("etc/etc.yaml"),
	 )

// 支持环境变量
cc := configcenter.MustNewConfigCenter[config.Config](
     configcenter.Config{Type: "yaml"},
     subscriber.MustNewFsnotifySubscriber("etc/etc.yaml", subscriber.WithUseEnv(true)),
   )
```

### 2. 获取配置

```go
// 获取配置
cfg, err := cc.GetConfig()

// 必须获取配置
cfg := cc.MustGetConfig()
```

### 3. 设置环境变量

```yaml
sqlx:
    # 从 DATASOURCE 获取 sqlx 的 datasource 配置, 未配置则为 jzero-admin.db
    datasource: "${DATASOURCE:-jzero-admin.db}"
    # 从 DRIVER_NAME 获取 sqlx 的 driverName 配置, 未配置则为 sqlite
    driverName: "${DRIVER_NAME:-sqlite}"
```

### 4. 设置动态配置回调

```go
cc.AddListener(func() {})
```