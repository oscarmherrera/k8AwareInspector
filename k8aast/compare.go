package k8aast

import (
	"github.com/pkg/errors"
	"go/ast"
	"go/token"
)

type codeDifference struct {
	DiffType int
	OldCode  interface{}
	NewCode  interface{}
	Pos      token.Pos
}

func comparePackageMaps(oldPackages, newPackages *map[string]*ast.Package) ([]*codeDifference, error) {

	differences := make([]*codeDifference, 0)

	for key, value := range *oldPackages {
		oldPack := value
		newPacks := *newPackages
		if newPack, ok := newPacks[key]; ok {
			diff, err := comparePackage(oldPack, newPack)
			if err != nil {
				return nil, errors.Wrap(err, "error compare old and new packages")
			}
			if diff != nil {
				differences = append(differences, diff...)
			}

		} else {
			//Todo add new package pointer
			diff := &codeDifference{
				DiffType: new_package,
				NewCode:  newPack,
			}
			differences = append(differences, diff)
		}
	}

	return differences, nil
}

func comparePackage(oldPack, newPack *ast.Package) ([]*codeDifference, error) {
	cDifferences := make([]*codeDifference, 0)

	scopeDifferences, err := compareScope(oldPack, newPack)
	if err != nil {
		return nil, errors.Wrap(err, "error comparing package maps: scope")
	}

	cDifferences = append(cDifferences, scopeDifferences...)

	importDifferences, err := compareImports(oldPack, newPack)
	if err != nil {
		return nil, errors.Wrap(err, "error comparing package maps: import")
	}

	cDifferences = append(cDifferences, importDifferences...)
	return cDifferences, nil
}
