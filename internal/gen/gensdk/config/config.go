package config

type Config struct {
	Language     string
	APP          string
	Module       string
	Output       string // output dir
	ApiDir       string
	ProtoDir     string
	WrapResponse bool
}
