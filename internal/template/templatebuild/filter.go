package templatebuild

func filter(name string) bool {
	switch name {
	case ".git":
		return true
	}
	return false
}
