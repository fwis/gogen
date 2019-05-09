package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"./gendb"
)

func main() {
	var structNames []string
	var MmPackage = ""
	var outputPath = ""

	MmPackageFlag := flag.String("i", "", "go import 路径，默认 biohit.com/mm")
	structNamesFlag := flag.String("s", "", "struct name，多个引号内逗号分隔")
	outputPathFlag := flag.String("o", "", "输出路径, 默认 src/biohit.com/dm")
	tablePrefixFlag := flag.String("p", "", "表名前缀, 默认空")
	flag.Parse()

	if *MmPackageFlag == "" {
		MmPackage = "biohit.com/mm"
	} else {
		MmPackage = *MmPackageFlag
	}

	if *structNamesFlag == "" {
		flag.Usage()
		return
	} else {
		structNames = strings.Split(*structNamesFlag, ",")
	}

	//m.Env.StructNames = structNames
	outputPath = *outputPathFlag
	_ = MmPackage

	env := gendb.NewEnv(outputPath)
	env.TmpDir = "D:/"
	if outputPath == "" {
		env.OutputDir = filepath.Join(env.GOPATH, "src", "biohit.com", "dm")
	}

	tablePrefix := ""
	if tablePrefixFlag != nil && *tablePrefixFlag != "" {
		tablePrefix = *tablePrefixFlag
	}
	err := gendb.CreateBootGoFile(env, dbgen_boot_tpl, func(vars map[string]interface{}) {
		List := []*AA{}
		for _, sn := range structNames {
			tableName := sn
			if tablePrefix != "" {
				tableName = tablePrefix + tableName
			}

			List = append(List, &AA{
				StructName:      sn,
				TableName:       tableName, //utils.ToSnakeCase(sn),
				LowerStructName: strings.ToLower(sn),
			})
		}

		vars["List"] = List
		vars["MmPackage"] = MmPackage
	})

	if err != nil {
		fmt.Print(err)
		return
	}
}

type AA struct {
	StructName      string
	TableName       string
	LowerStructName string
}

const dbgen_boot_tpl = `
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"{{$.MmPackage}}"

	"github.com/fwis/gogen/gendb"
)

func main() {
	outputdir := "{{EscapePath $.Env.OutputDir}}"
	var err error

	{{range $i, $n := .List -}}
	sqlgen{{$i}} := gendb.CreateSQLGenByObj(&mm.{{$n.StructName}}{}, "{{$n.TableName}}", "") //默认取结构体第一个字段名称为主键
	bb{{$i}}, err := gendb.RenderDMgo("{{EscapePath $.Env.TemplatePath}}", "{{$.MmPackage}}", sqlgen{{$i}})
	if err != nil {
		panic(err)
	}
	gengofile{{$i}} := filepath.Join(outputdir,"dm_{{$n.LowerStructName}}_gen.go")

	/*
	skip{{$i}}, exist{{$i}}, err := gendb.CheckGenConflict(gengofile{{$i}})
	_ = exist{{$i}}
	if err != nil {
		panic(err)
	}
	*/
	skip{{$i}} := false

	if !skip{{$i}} {
		err = ioutil.WriteFile(gengofile{{$i}}, bb{{$i}}, os.ModePerm)
		if err != nil {
			panic(err)
		}
		err = gendb.FmtGoFile(gengofile{{$i}})
		if err != nil {
			fmt.Printf("format go file, err=%v\n", err)
		}
	} else {
		fmt.Printf("SKIP gen: [%s]\n", gengofile{{$i}})
	}
	{{end -}}
}
`
