package embedx

func WithDir(dir string) Opts {
	return func(config *embedxConfig) {
		config.Dir = dir
	}
}

func WithFileMatchFunc(fileFilter func(path string) bool) Opts {
	return func(config *embedxConfig) {
		config.FileMatchFunc = fileFilter
	}
}
