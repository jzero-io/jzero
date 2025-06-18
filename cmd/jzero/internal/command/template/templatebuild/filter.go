package templatebuild

import (
	"os"
	"path/filepath"

	"github.com/moby/patternmatcher"
)

func filter(dir, name string, matcher *patternmatcher.PatternMatcher) bool {
	pwd, err := os.Getwd()
	if err != nil {
		return true
	}
	target := filepath.Join(dir, name)
	relFilePath, err := filepath.Rel(pwd, target)
	if err != nil {
		return true
	}
	skip, err := filepathMatches(matcher, relFilePath)
	if err != nil || skip {
		return true
	}
	return false
}

func filepathMatches(matcher *patternmatcher.PatternMatcher, file string) (bool, error) {
	file = filepath.Clean(file)
	if file == "." {
		// Don't let them exclude everything, kind of silly.
		return false, nil
	}
	return matcher.MatchesOrParentMatches(file)
}
