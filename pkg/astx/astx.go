package astx

import (
	"go/ast"
	"go/token"
)

// HasImport checks if the given import path is already declared in the file.
func HasImport(f *ast.File, path string) bool {
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}
		for _, spec := range genDecl.Specs {
			importSpec, ok := spec.(*ast.ImportSpec)
			if !ok {
				continue
			}
			if importSpec.Path.Value == path {
				return true
			}
		}
	}
	return false
}

// AddImport adds the import declaration to the file if it's not already marked as added.
func AddImport(f *ast.File, path string, addedImports map[string]bool) {
	if !addedImports[path] {
		addedImports[path] = true
		f.Decls = append([]ast.Decl{&ast.GenDecl{
			Tok: token.IMPORT,
			Specs: []ast.Spec{
				&ast.ImportSpec{
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: path,
					},
				},
			},
		}}, f.Decls...)
	}
}

// DeleteImport removes the import declaration from the file if it exists.
func DeleteImport(f *ast.File, path string) {
	for i, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}
		for j, spec := range genDecl.Specs {
			importSpec, ok := spec.(*ast.ImportSpec)
			if !ok {
				continue
			}
			if importSpec.Path.Value == path {
				// Remove the import spec from the slice
				genDecl.Specs = append(genDecl.Specs[:j], genDecl.Specs[j+1:]...)
				// If there are no specs left, remove the declaration
				if len(genDecl.Specs) == 0 {
					f.Decls = append(f.Decls[:i], f.Decls[i+1:]...)
				}
				return
			}
		}
	}
}
