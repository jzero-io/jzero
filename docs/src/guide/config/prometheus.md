---
title: prometheus 配置
icon: vscode-icons:file-type-prometheus
star: true
order: 4
category: 配置
tag:
  - Guide
---

## rest

修改 etc/etc.yaml, 增加以下配置

```yaml
rest:
  devServer:
    enabled: true
```

## zrpc

修改 etc/etc.yaml, 增加以下配置

```shell
zrpc:
  devServer:
    enabled: true
```

## gateway

修改 etc/etc.yaml, 增加以下配置

```shell
zrpc:
  devServer:
    enabled: true
gateway:
  devServer:
    enabled: true
```

## Grafana

下面给出一个适配jzero的Grafana控制面板的JSON Model.

:::tip 
1、panels 中各个指标的 datasrouce 指代 prometheus 数据源，请更改为自己的数据源 uid
2、templating.list 中配置的是 dashboard 的变量信息，请根据需要修改
:::

```json
{
  "annotations": {
    "list": [
      {
        "builtIn": 1,
        "datasource": {
          "type": "grafana",
          "uid": "-- Grafana --"
        },
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
      }
    ]
  },
  "description": "专为Jzero框架设计的监控面板",
  "editable": true,
  "fiscalYearStartMonth": 0,
  "gnetId": 19909,
  "graphTooltip": 1,
  "id": 17,
  "links": [],
  "liveNow": false,
  "panels": [
    {
      "collapsed": false,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 0
      },
      "id": 9,
      "panels": [],
      "title": "系统信息",
      "type": "row"
    },
    {
      "datasource": {
        "default": true,
        "type": "prometheus",
        "uid": "fe0mra1jcnklca"
      },
      "description": "",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisBorderShow": false,
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "barWidthFactor": 0.6,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 1
      },
      "id": 8,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "fe0mra1jcnklca"
          },
          "editorMode": "code",
          "expr": "avg(go_goroutines)",
          "instant": false,
          "legendFormat": "goroutine",
          "range": true,
          "refId": "go_goroutines"
        },
        {
          "datasource": {
            "type": "prometheus",
            "uid": "fe0mra1jcnklca"
          },
          "editorMode": "code",
          "expr": "avg(go_threads)",
          "hide": false,
          "instant": false,
          "legendFormat": "threads",
          "range": true,
          "refId": "go_threads"
        }
      ],
      "title": "CPU",
      "type": "timeseries"
    },
    {
      "datasource": {
        "default": true,
        "type": "prometheus",
        "uid": "fe0mra1jcnklca"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisBorderShow": false,
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "barWidthFactor": 0.6,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "decbytes"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 1
      },
      "id": 12,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "fe0mra1jcnklca"
          },
          "editorMode": "code",
          "expr": "go_memstats_heap_inuse_bytes{instance='$instance'}+go_memstats_heap_idle_bytes{instance='$instance'}",
          "instant": false,
          "legendFormat": "heap",
          "range": true,
          "refId": "使用的堆内存"
        }
      ],
      "title": "memory",
      "type": "timeseries"
    },
    {
      "collapsed": false,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 9
      },
      "id": 11,
      "panels": [],
      "title": "RPC请求信息",
      "type": "row"
    },
    {
      "datasource": {
        "default": true,
        "type": "prometheus",
        "uid": "fe0mra1jcnklca"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisBorderShow": false,
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "barWidthFactor": 0.6,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 10
      },
      "id": 14,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "right",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "fe0mra1jcnklca"
          },
          "editorMode": "code",
          "expr": "round(increase(rpc_server_requests_code_total{instance=\"$instance\"}[12h]) )>=0",
          "instant": false,
          "legendFormat": "{{app}} | 状态码 {{code}} | {{method}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "RPC Server请求个数[12h]",
      "type": "timeseries"
    },
    {
      "datasource": {
        "default": true,
        "type": "prometheus",
        "uid": "fe0mra1jcnklca"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisBorderShow": false,
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "barWidthFactor": 0.6,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 10
      },
      "id": 15,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "right",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "fe0mra1jcnklca"
          },
          "editorMode": "code",
          "expr": "round(increase(rpc_client_requests_code_total{instance=\"$instance\"}[12h]) )>=0",
          "instant": false,
          "legendFormat": "{{app}} | 状态码 {{code}} | {{method}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "RPC Client请求个数[12h]",
      "type": "timeseries"
    },
    {
      "datasource": {
        "default": true,
        "type": "prometheus",
        "uid": "fe0mra1jcnklca"
      },
      "description": "",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisBorderShow": false,
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "barWidthFactor": 0.6,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green"
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "reqps"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 18
      },
      "id": 4,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "right",
          "showLegend": true,
          "width": 300
        },
        "tooltip": {
          "mode": "multi",
          "sort": "desc"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "fe0mra1jcnklca"
          },
          "editorMode": "code",
          "expr": "sum(rate(rpc_server_requests_code_total{namespace=~\"$Namespace\",app=~\"$Service\",pod=~\"$Pod\"}[5m])) by (app, code, method)",
          "instant": false,
          "legendFormat": "{{app}} | 状态码 {{code}} | {{method}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "RPC Server QPS",
      "type": "timeseries"
    },
    {
      "datasource": {
        "default": true,
        "type": "prometheus",
        "uid": "fe0mra1jcnklca"
      },
      "description": "",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisBorderShow": false,
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "barWidthFactor": 0.6,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green"
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "reqps"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 18
      },
      "id": 1,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "right",
          "showLegend": true,
          "width": 300
        },
        "tooltip": {
          "mode": "multi",
          "sort": "desc"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "fe0mra1jcnklca"
          },
          "editorMode": "code",
          "expr": "sum(rate(rpc_client_requests_code_total{namespace=~\"$Namespace\",app=~\"$Service\",pod=~\"$Pod\"}[5m])) by (app, code, method)",
          "instant": false,
          "legendFormat": "{{app}} | 状态码 {{code}} | {{method}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "RPC Client QPS",
      "type": "timeseries"
    },
    {
      "datasource": {
        "default": true,
        "type": "prometheus",
        "uid": "fe0mra1jcnklca"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisBorderShow": false,
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "barWidthFactor": 0.6,
            "drawStyle": "line",
            "fillOpacity": 15,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineStyle": {
              "fill": "solid"
            },
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "dashed"
            }
          },
          "decimals": 0,
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green"
              },
              {
                "color": "#EAB839",
                "value": 500
              },
              {
                "color": "red",
                "value": 1000
              }
            ]
          },
          "unit": "ms"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 26
      },
      "id": 3,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "right",
          "showLegend": true,
          "width": 300
        },
        "tooltip": {
          "mode": "multi",
          "sort": "desc"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "fe0mra1jcnklca"
          },
          "editorMode": "code",
          "expr": "sum(rate(rpc_server_requests_duration_ms_sum{namespace=~\"$Namespace\",app=~\"$Service\",pod=~\"$Pod\"}[5m]) / rate(rpc_server_requests_duration_ms_count{namespace=~\"$Namespace\",app=~\"$Service\",pod=~\"$Pod\"}[5m])) by (app, method)",
          "instant": false,
          "legendFormat": "{{app}} | {{method}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "RPC Server请求耗时",
      "type": "timeseries"
    },
    {
      "datasource": {
        "default": true,
        "type": "prometheus",
        "uid": "fe0mra1jcnklca"
      },
      "description": "",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisBorderShow": false,
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "barWidthFactor": 0.6,
            "drawStyle": "line",
            "fillOpacity": 15,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineStyle": {
              "fill": "solid"
            },
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "dashed"
            }
          },
          "decimals": 0,
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green"
              },
              {
                "color": "#EAB839",
                "value": 500
              },
              {
                "color": "red",
                "value": 1000
              }
            ]
          },
          "unit": "ms"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 26
      },
      "id": 2,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "right",
          "showLegend": true,
          "width": 300
        },
        "tooltip": {
          "mode": "multi",
          "sort": "desc"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "fe0mra1jcnklca"
          },
          "editorMode": "code",
          "expr": "sum(rate(rpc_client_requests_duration_ms_sum{namespace=~\"$Namespace\",app=~\"$Service\",pod=~\"$Pod\"}[5m]) / rate(rpc_client_requests_duration_ms_count{namespace=~\"$Namespace\",app=~\"$Service\",pod=~\"$Pod\"}[5m])) by (app, method)",
          "instant": false,
          "legendFormat": "{{app}} | {{method}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "RPC Client 请求耗时",
      "type": "timeseries"
    },
    {
      "collapsed": false,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 34
      },
      "id": 13,
      "panels": [],
      "title": "HTTP请求信息",
      "type": "row"
    },
    {
      "datasource": {
        "default": true,
        "type": "prometheus",
        "uid": "fe0mra1jcnklca"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisBorderShow": false,
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "barWidthFactor": 0.6,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green"
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 35
      },
      "id": 16,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "right",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "fe0mra1jcnklca"
          },
          "editorMode": "code",
          "expr": "round(increase(http_server_requests_code_total{instance=\"$instance\"}[12h]) )>=0",
          "instant": false,
          "legendFormat": "{{code}} | {{method}} | {{path}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Http Server 请求个数[12h]",
      "type": "timeseries"
    },
    {
      "datasource": {
        "default": true,
        "type": "prometheus",
        "uid": "fe0mra1jcnklca"
      },
      "description": "",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisBorderShow": false,
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "barWidthFactor": 0.6,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green"
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "reqps"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 35
      },
      "id": 6,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "right",
          "showLegend": true,
          "width": 300
        },
        "tooltip": {
          "mode": "multi",
          "sort": "desc"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "fe0mra1jcnklca"
          },
          "editorMode": "code",
          "expr": "sum(rate(http_server_requests_code_total{namespace=~\"$Namespace\",app=~\"$Service\",pod=~\"$Pod\"}[5m])) by (app, code, method, path)",
          "instant": false,
          "legendFormat": "{{app}} | {{code}} | {{method}} | {{path}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Http Server 请求QPS",
      "type": "timeseries"
    },
    {
      "datasource": {
        "default": true,
        "type": "prometheus",
        "uid": "fe0mra1jcnklca"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisBorderShow": false,
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "barWidthFactor": 0.6,
            "drawStyle": "line",
            "fillOpacity": 15,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "insertNulls": false,
            "lineInterpolation": "linear",
            "lineStyle": {
              "fill": "solid"
            },
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "dashed"
            }
          },
          "decimals": 0,
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green"
              },
              {
                "color": "#EAB839",
                "value": 500
              },
              {
                "color": "red",
                "value": 1000
              }
            ]
          },
          "unit": "ms"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 43
      },
      "id": 5,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "right",
          "showLegend": true,
          "width": 300
        },
        "tooltip": {
          "mode": "multi",
          "sort": "desc"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "fe0mra1jcnklca"
          },
          "editorMode": "code",
          "expr": "sum(rate(http_server_requests_duration_ms_sum{namespace=~\"$Namespace\",app=~\"$Service\",pod=~\"$Pod\"}[5m]) / rate(http_server_requests_duration_ms_count{namespace=~\"$Namespace\",app=~\"$Service\",pod=~\"$Pod\"}[5m])) by (app, method, path)",
          "instant": false,
          "legendFormat": "{{app}} | {{method}} | {{path}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Http Server 请求耗时",
      "type": "timeseries"
    }
  ],
  "refresh": "10s",
  "schemaVersion": 39,
  "tags": [
    "go-zero",
    "portal"
  ],
  "templating": {
    "list": [
      {
        "allValue": ".*",
        "current": {
          "selected": false,
          "text": "All",
          "value": "$__all"
        },
        "datasource": {
          "type": "prometheus",
          "uid": "fe0mra1jcnklca"
        },
        "definition": "label_values(instance)",
        "hide": 0,
        "includeAll": false,
        "label": "portal实例",
        "multi": false,
        "name": "instance",
        "options": [],
        "query": {
          "qryType": 1,
          "query": "label_values(instance)",
          "refId": "PrometheusVariableQueryEditor-VariableQuery"
        },
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "sort": 0,
        "type": "query"
      },
      {
        "allValue": ".*",
        "current": {
          "selected": false,
          "text": "All",
          "value": "$__all"
        },
        "datasource": {
          "type": "prometheus",
          "uid": "fe0mra1jcnklca"
        },
        "definition": "label_values(kube_pod_info,namespace)",
        "hide": 0,
        "includeAll": true,
        "label": "Namespace",
        "multi": false,
        "name": "Namespace",
        "options": [],
        "query": {
          "query": "label_values(kube_pod_info,namespace)",
          "refId": "PrometheusVariableQueryEditor-VariableQuery"
        },
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "sort": 1,
        "type": "query"
      },
      {
        "allValue": ".*",
        "current": {
          "selected": false,
          "text": "All",
          "value": "$__all"
        },
        "datasource": {
          "type": "prometheus",
          "uid": "fe0mra1jcnklca"
        },
        "definition": "label_values(kube_service_info{namespace=~\"$Namespace\"},service)",
        "hide": 0,
        "includeAll": true,
        "label": "Service",
        "multi": true,
        "name": "Service",
        "options": [],
        "query": {
          "query": "label_values(kube_service_info{namespace=~\"$Namespace\"},service)",
          "refId": "PrometheusVariableQueryEditor-VariableQuery"
        },
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "sort": 1,
        "type": "query"
      },
      {
        "allValue": "",
        "current": {
          "selected": false,
          "text": "All",
          "value": "$__all"
        },
        "datasource": {
          "type": "prometheus",
          "uid": "fe0mra1jcnklca"
        },
        "definition": "label_values(kube_pod_container_info{namespace=~\"$Namespace\"},pod)",
        "hide": 0,
        "includeAll": true,
        "label": "Pod",
        "multi": true,
        "name": "Pod",
        "options": [],
        "query": {
          "query": "label_values(kube_pod_container_info{namespace=~\"$Namespace\"},pod)",
          "refId": "PrometheusVariableQueryEditor-VariableQuery"
        },
        "refresh": 1,
        "regex": "/^$Service.*/",
        "skipUrlSync": false,
        "sort": 1,
        "type": "query"
      }
    ]
  },
  "time": {
    "from": "now-1h",
    "to": "now"
  },
  "timepicker": {},
  "timezone": "",
  "title": "Dashboard-for-Jzero",
  "uid": "a7727282-f5ed-49f8-baf4-29d24fb1de49",
  "version": 19,
  "weekStart": ""
}
```
