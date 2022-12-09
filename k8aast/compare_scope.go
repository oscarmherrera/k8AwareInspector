package k8aast

import "go/ast"

const (
	scope_diff = iota + 50
	scope_decl
	scope_data
	scope_type
	scope_kind
)

func compareScope(oldPack, newPack *ast.Package) ([]*codeDifference, error) {

	cDifferences := make([]*codeDifference, 0)

	if oldPack.Scope == nil && newPack.Scope == nil {
		return cDifferences, nil
	}

	oldPackObject := oldPack.Scope.Objects
	newPackObject := newPack.Scope.Objects

	return compareHeaderObject(oldPackObject, newPackObject)

}

func compareHeaderObject(oldPackObject, newPackObject map[string]*ast.Object) ([]*codeDifference, error) {

	cDifferences := make([]*codeDifference, 0)

	for key, oldObject := range oldPackObject {
		if newObject, ok := newPackObject[key]; ok {
			// Found it now compare it
			if oldObject.Decl != newObject.Decl {
				diff := &codeDifference{
					DiffType: scope_decl,
					OldCode:  oldObject.Decl,
					NewCode:  newObject.Decl,
					Pos:      oldObject.Pos(),
				}
				cDifferences = append(cDifferences, diff)
			}
			if oldObject.Data != newObject.Data {
				diff := &codeDifference{
					DiffType: scope_data,
					OldCode:  oldObject.Data,
					NewCode:  newObject.Data,
					Pos:      oldObject.Pos(),
				}
				cDifferences = append(cDifferences, diff)
			}
			if oldObject.Type != newObject.Type {
				diff := &codeDifference{
					DiffType: scope_type,
					OldCode:  oldObject.Type,
					NewCode:  newObject.Type,
					Pos:      oldObject.Pos(),
				}
				cDifferences = append(cDifferences, diff)
			}

			if oldObject.Kind != newObject.Kind {
				diff := &codeDifference{
					DiffType: scope_kind,
					OldCode:  oldObject.Kind,
					NewCode:  newObject.Kind,
					Pos:      oldObject.Pos(),
				}
				cDifferences = append(cDifferences, diff)
			}
		}

	}
	return cDifferences, nil
}
