package config

type Config struct {
	Language     string
	APP          string
	GenModule    bool
	Module       string
	Output       string // output dir
	ApiDir       string
	ProtoDir     string
	WrapResponse bool
}
