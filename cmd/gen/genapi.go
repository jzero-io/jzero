package gen

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/app/pkg/templatex"
	"github.com/rinchsan/gosimports"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	apiparser "github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"

	"github.com/jzero-io/jzero/embeded"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"golang.org/x/exp/constraints"
)

type JzeroApi struct {
	Wd     string
	Module string
}

func (ja *JzeroApi) Gen() error {
	apiDirName := filepath.Join(ja.Wd, "app", "desc", "api")
	if pathx.FileExists(apiDirName) {
		// format api dir
		command := fmt.Sprintf("goctl api format --dir %s", apiDirName)
		_, err := execx.Run(command, ja.Wd)
		if err != nil {
			return err
		}

		fmt.Printf("%s to generate api code.\n", color.WithColor("Start", color.FgGreen))
		mainApiFilePath := GetMainApiFilePath(apiDirName)

		err = generateApiCode(ja.Wd, mainApiFilePath)
		if err != nil {
			return err
		}
		fmt.Println(color.WithColor("Done", color.FgGreen))

		// 生成 api types
		apiFileTypes := make([]ApiFileTypes, 0)

		allApiFilePaths := getAllApiFilePath(apiDirName)

		for _, v := range allApiFilePaths {
			apiSpec, err := apiparser.Parse(filepath.Join(apiDirName, v))
			if err != nil {
				return err
			}
			apiFileTypes = append(apiFileTypes, ApiFileTypes{
				Filepath: v,
				ApiSpec:  *apiSpec,
			})
		}

		types := getFileTypes(apiFileTypes)
		for _, t := range types {
			if len(t.GenTypes) == 0 {
				continue
			}
			typesGoString, err := gogen.BuildTypes(t.GenTypes)
			if err != nil {
				return err
			}
			typesGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
				"Module": ja.Module,
				"Types":  typesGoString,
			}, []byte(TypesGoTpl))
			if err != nil {
				return err
			}

			typesGoFormatBytes, err := gosimports.Process("", typesGoBytes, &gosimports.Options{
				FormatOnly: true,
				Comments:   true,
			})
			if err != nil {
				return err
			}

			prefix := strings.ReplaceAll(filepath.Dir(t.Filepath), "/", "_") + "_"
			if len(strings.Split(t.Filepath, "/")) == 1 {
				prefix = ""
			}

			fileBase := filepath.Base(t.Filepath)
			typesGoFilePath := prefix + fileBase[0:len(fileBase)-len(path.Ext(fileBase))] + ".types.go"
			if t.Base {
				typesGoFilePath = "types.go"
			}
			err = os.WriteFile(filepath.Join("app", "internal", "types", typesGoFilePath), typesGoFormatBytes, 0o644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getFileTypes(apiFileTypes []ApiFileTypes) []ApiFileTypes {
	var newApiFileTypes []ApiFileTypes

	allTypesRawName := make([][]string, 0, len(apiFileTypes)+1)

	for _, apiFileType := range apiFileTypes {
		var typesRawName []string
		for _, name := range apiFileType.ApiSpec.Types {
			typesRawName = append(typesRawName, name.Name())
		}
		allTypesRawName = append(allTypesRawName, typesRawName)
	}

	elements := separateCommonElements(allTypesRawName...)

	for i, apiFileType := range apiFileTypes {
		newApiFileType := apiFileTypes[i]
		var genTypes []spec.Type
		elementArray := elements[i]
		for _, e := range elementArray {
			for _, t := range apiFileType.ApiSpec.Types {
				if t.Name() == e {
					genTypes = append(genTypes, t)
				}
			}
		}
		newApiFileType.GenTypes = genTypes
		newApiFileTypes = append(newApiFileTypes, newApiFileType)
	}

	// append base
	var genTypes []spec.Type

	typeSet := make(map[string]struct{})
	for _, apiFileType := range apiFileTypes {
		for _, genType := range apiFileType.ApiSpec.Types {
			for _, e := range elements[len(elements)-1] {
				if e == genType.Name() {
					if _, ok := typeSet[genType.Name()]; !ok {
						genTypes = append(genTypes, genType)
						typeSet[genType.Name()] = struct{}{}
					}
				}
			}
		}
	}
	newApiFileTypes = append(newApiFileTypes, ApiFileTypes{
		GenTypes: genTypes,
		Base:     true,
	})

	return newApiFileTypes
}

func getAllApiFilePath(apiDirName string) []string {
	var apiFiles []string
	_ = filepath.Walk(apiDirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".api" {
			rel, err := filepath.Rel(apiDirName, path)
			if err != nil {
				return err
			}
			apiFiles = append(apiFiles, filepath.ToSlash(rel))
		}
		return nil
	})
	return apiFiles
}

func getRouteApiFilePath(apiDirName string) []string {
	var apiFiles []string
	_ = filepath.Walk(apiDirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".api" {
			apiSpec, err := parser.Parse(path, nil)
			if err != nil {
				return err
			}
			if len(apiSpec.Service.Routes()) > 0 {
				rel, err := filepath.Rel(apiDirName, path)
				if err != nil {
					return err
				}
				apiFiles = append(apiFiles, filepath.ToSlash(rel))
			}
		}
		return nil
	})
	return apiFiles
}

func generateApiCode(wd string, mainApiFilePath string) error {
	if mainApiFilePath == "" {
		return errors.New("empty mainApiFilePath")
	}
	defer os.Remove(mainApiFilePath)

	fmt.Printf("%s api file %s\n", color.WithColor("Using", color.FgGreen), mainApiFilePath)
	command := fmt.Sprintf("goctl api go --api %s --dir ./app --home %s", mainApiFilePath, filepath.Join(embeded.Home, "go-zero"))
	if _, err := execx.Run(command, wd); err != nil {
		return err
	}
	return nil
}

func separateCommonElements(arrays ...[]string) [][]string {
	if len(arrays) == 0 {
		return nil
	}

	// 创建一个 map 来统计每个元素的出现次数
	elementCount := make(map[string]int)

	// 遍历所有数组,统计每个元素的出现次数
	for _, arr := range arrays {
		for _, elem := range arr {
			elementCount[elem]++
		}
	}

	// 创建一个切片来存储独一无二的元素数组
	uniqueArrays := make([][]string, len(arrays))

	// 创建一个切片来存储公共元素
	commonElements := make([]string, 0)

	// 遍历所有数组,将元素加入到对应的数组中
	for i, arr := range arrays {
		uniqueArr := make([]string, 0)
		for _, elem := range arr {
			// 如果元素只出现在该数组中,则加入独一无二的元素数组
			if elementCount[elem] == 1 {
				uniqueArr = append(uniqueArr, elem)
			} else if elementCount[elem] > 1 {
				// 如果元素出现在多个数组中,则加入公共元素数组
				commonElements = append(commonElements, elem)
			} else {
				// 如果元素不是独一无二的,也不是公共的,则加入独一无二的元素数组
				uniqueArr = append(uniqueArr, elem)
			}
		}
		uniqueArrays[i] = uniqueArr
	}

	// 去重公共元素数组
	commonElements = unifySlice(commonElements)

	// 将公共元素数组加到结果的最后
	uniqueArrays = append(uniqueArrays, commonElements)
	return uniqueArrays
}

func unifySlice[T constraints.Ordered](slice []T) []T {
	uniqueElements := make(map[T]struct{})
	for _, elem := range slice {
		uniqueElements[elem] = struct{}{}
	}
	result := make([]T, 0, len(uniqueElements))
	for elem := range uniqueElements {
		result = append(result, elem)
	}
	return result
}
