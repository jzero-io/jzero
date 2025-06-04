module {{ .Module }}

go {{ .GoVersion }}

{{if (VersionCompare .GoVersion ">=" "1.24")}}
tool (
	github.com/jzero-io/jzero/cmd/jzero
	github.com/zeromicro/go-zero/tools/goctl
)
{{end}}