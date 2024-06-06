package templatebuild

func filter(name string) bool {
	switch name {
	case ".git":
		return true
	case ".github":
		return true
	case ".goreleaser.yaml", ".goreleaser.yml":
		return true
	case "Taskfile.yml":
		return true
	}
	return false
}
