package k8aast

import (
	"fmt"
	"github.com/google/go-github/v48/github"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go/ast"
	"k8Aware/k8agithub"
	"strings"
)

const (
	new_tag = iota
	old_tag

	new_package
	new_difference
	changed_difference

	file_difference
)

func ProcessTags(ghClient *github.Client, urls []*string, accessToken *string, tempDir string, plantUML bool, zl *zap.Logger) error {
	setLogging(zl)

	filenameNew, err := k8agithub.GetGithubTagZip(ghClient, *urls[0], accessToken, tempDir, zLog)
	if err != nil {
		return err
	}

	zLog.Info("file newTag", zap.String("newTag", *filenameNew))

	filenameOld, err := k8agithub.GetGithubTagZip(ghClient, *urls[1], accessToken, tempDir, zLog)
	if err != nil {
		return err
	}
	zLog.Info("file oldTag", zap.String("oldTag", *filenameOld))

	tempDirNew := tempDir + "/new"
	tempDirOld := tempDir + "/old"

	err = unzipSource(*filenameNew, tempDirNew)
	if err != nil {
		return errors.Wrap(err, "expanding newtag archive")
	}

	err = unzipSource(*filenameOld, tempDirOld)
	if err != nil {
		return errors.Wrap(err, "expanding oldtag archive")
	}

	if plantUML == true {

		newAst, err := GeneratePlantUML(tempDir, *filenameNew, zLog)
		if err != nil {
			return errors.Wrap(err, "generating plant uml for newtag")

		}

		oldAst, err := GeneratePlantUML(tempDir, *filenameOld, zLog)
		if err != nil {
			return errors.Wrap(err, "generating plant uml for oldtag")
		}

		printNewAst := strings.ReplaceAll(*newAst, "\n", "\n\r")
		printOldAst := strings.ReplaceAll(*oldAst, "\n", "\n\r")

		sugar.Debug("New PlantUML: %s", printNewAst)
		sugar.Debug("Old PlantUML: %s", printOldAst)
	}

	packageMapList := make([]*map[string]*ast.Package, 2)
	newPackageMap, err := GenerateGoAST(tempDirNew)
	if err != nil {
		return errors.Wrap(err, "error processing GoAST newTag")
	}
	packageMapList[new_tag] = newPackageMap

	oldPackageMap, err := GenerateGoAST(tempDirOld)
	if err != nil {
		return errors.Wrap(err, "error processing GoAST newTag")
	}
	packageMapList[old_tag] = oldPackageMap

	differences, err := comparePackageMaps(oldPackageMap, newPackageMap)
	if err != nil {
		return errors.Wrap(err, "error comparing package maps")
	}

	printDifferences(differences)

	return nil
}

func printDifferences(differences []*codeDifference) {

	if len(differences) == 0 {
		zLog.Info("no differences")
	}

	for _, item := range differences {
		//value := *item
		fmt.Sprintln("Diff Type %d Position:%d", item.DiffType, item.Pos)
		fmt.Sprintln("Old Code %v", item.OldCode)
		fmt.Sprintln("New Code %v", item.NewCode)
	}

}
