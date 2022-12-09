package k8aast

import "go/ast"

const (
	import_diff = iota + 70
)

func compareImports(oldPack, newPack *ast.Package) ([]*codeDifference, error) {

	cDifferences := make([]*codeDifference, 0)
	if oldPack.Imports == nil && newPack.Imports == nil {
		return cDifferences, nil
	}

	oldImports := oldPack.Imports
	newImports := newPack.Imports

	return compareHeaderObject(oldImports, newImports)

}
