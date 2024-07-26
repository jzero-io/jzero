package config

type Config struct {
	Language     string
	Scope        string
	GenModule    bool
	GoPackage    string
	GoModule     string
	Output       string
	ApiDir       string
	ProtoDir     string
	WrapResponse bool
}
