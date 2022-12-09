package k8aast

import (
	goplantuml "github.com/jfeliu007/goplantuml/parser"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var (
	zLog  *zap.Logger
	sugar *zap.SugaredLogger
)

func setLogging(zl *zap.Logger) {
	if zLog == nil {
		zLog = zl
		sugar = zl.Sugar()
	}
}

func GenerateGoAST(tempDir string) (*map[string]*ast.Package, error) {
	fset := token.NewFileSet()

	packageList := make([]string, 0)
	//packageASTList := make([]interface{}, 20)

	err := filepath.Walk(tempDir,

		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}
			if info.IsDir() == true {

				if strings.HasPrefix(info.Name(), ".") || info.Name() == "vendor" {
					return filepath.SkipDir
				}

				list, err := os.ReadDir(path)
				if err != nil {
					return err
				}
				for _, d := range list {
					if d.IsDir() || !strings.HasSuffix(d.Name(), ".go") {
						continue
					} else {
						packageList = append(packageList, path)
						break
					}

				}
			}
			return err
		})

	if err != nil {
		return nil, errors.Wrap(err, "error acquiring package list ")
	}

	packageMap := map[string]*ast.Package{}

	for _, dirValue := range packageList {
		//var b bytes.Buffer
		f, err := parser.ParseDir(fset, dirValue, nil, parser.AllErrors)
		if err != nil {
			return nil, errors.Wrap(err, "error parsing package ")
		}
		for k, value := range f {
			if _, ok := packageMap[k]; ok {
				zLog.Error("key already exists in package map", zap.String("key", k))
				continue
			}
			packageMap[k] = value
		}
		//err = ast.Fprint(&b, fset, f, nil)
		//if err != nil {
		//	return errors.Wrap(err, "error printing ast for package ")
		//}
		//prettyString := fmt.Sprintln(b.String())
		//sugar.Debugln("%s", prettyString)

	}
	return &packageMap, nil
}

func GeneratePlantUML(destinationPath, archiveFilePath string, zl *zap.Logger) (*string, error) {
	setLogging(zl)

	err := unzipSource(archiveFilePath, destinationPath)
	if err != nil {
		return nil, err
	}

	var directoryPaths []string
	directoryPaths = append(directoryPaths, destinationPath)

	options := &goplantuml.ClassDiagramOptions{
		Directories:        directoryPaths,
		IgnoredDirectories: nil,
		Recursive:          true,
		RenderingOptions:   map[goplantuml.RenderingOption]interface{}{},
		FileSystem:         afero.NewOsFs(),
	}
	p, err := goplantuml.NewClassDiagramWithOptions(options)
	if err != nil {
		return nil, err
	}

	err = p.SetRenderingOptions(map[goplantuml.RenderingOption]interface{}{
		goplantuml.RenderPrivateMembers:    true,
		goplantuml.RenderAggregations:      true,
		goplantuml.AggregatePrivateMembers: true,
		goplantuml.RenderCompositions:      true,
		goplantuml.RenderImplementations:   true,
		goplantuml.RenderAliases:           true,
		goplantuml.RenderFields:            true,
		goplantuml.RenderMethods:           true,
	})
	if err != nil {
		return nil, err
	}
	diag := p.Render()
	return &diag, nil

}
