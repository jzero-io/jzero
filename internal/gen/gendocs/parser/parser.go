package parser

import (
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

type DocsParser struct {
	apiSpec *spec.ApiSpec
}

func NewDocsParser(apiSpec *spec.ApiSpec) *DocsParser {
	return &DocsParser{apiSpec: apiSpec}
}

// DocsSpec represents the hierarchical structure of the documentation.
type DocsSpec struct {
	Group     string
	Children  []*DocsSpec
	GroupSpec *spec.Group
}

// BuildDocsSpecHierarchy builds the hierarchical structure of DocsSpec based on the group strings.
func (dp *DocsParser) BuildDocsSpecHierarchy(groups []string) []*DocsSpec {
	root := &DocsSpec{
		Group:    "",
		Children: []*DocsSpec{},
	}

	for _, group := range groups {
		dp.addGroupToHierarchy(root, group)
	}

	return root.Children
}

// addGroupToHierarchy adds a group string to the hierarchical structure.
func (dp *DocsParser) addGroupToHierarchy(root *DocsSpec, group string) {
	parts := strings.Split(group, "/")
	current := root

	for _, part := range parts {
		child := dp.findOrCreateChild(current, part)
		current = child
	}
}

// findOrCreateChild finds an existing child with the given group name, or creates a new one if not found.
func (dp *DocsParser) findOrCreateChild(parent *DocsSpec, group string) *DocsSpec {
	for _, child := range parent.Children {
		if child.Group == group {
			return child
		}
	}

	newChild := &DocsSpec{
		Group:     filepath.ToSlash(filepath.Join(parent.Group, group)),
		Children:  []*DocsSpec{},
		GroupSpec: dp.findGroupSpec(filepath.ToSlash(filepath.Join(parent.Group, group))),
	}
	parent.Children = append(parent.Children, newChild)
	return newChild
}

func (dp *DocsParser) findGroupSpec(group string) *spec.Group {
	for _, v := range dp.apiSpec.Service.Groups {
		if v.GetAnnotation("group") == group {
			return &v
		}
	}
	return nil
}
