package gendb

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var gRenderTemplateFuncs template.FuncMap = template.FuncMap{
	"EscapePath": func(path string) string {
		return strings.Replace(path, "\\", "\\\\", -1)
	},
}

func CreateBootGoFile(env *Env, boot_tpl string, fillVarsFunc func(map[string]interface{})) error {
	vars := make(map[string]interface{})
	vars["Env"] = env
	fillVarsFunc(vars)
	bb, err := RenderTpl("", boot_tpl, vars)
	if err != nil {
		return err
	}
	nowstr := time.Now().Format("20060102150415")
	bootgofile := filepath.Join(env.TmpDir, "gogen_boot_"+nowstr+".go")
	err = ioutil.WriteFile(bootgofile, bb, os.ModePerm)
	if err != nil {
		return err
	}

	return runGoFile(bootgofile)
}

//根据模版生成内容
func RenderTpl(tplname string, tplcontent string, data interface{}) ([]byte, error) {
	out := new(bytes.Buffer)
	var err error
	var t *template.Template

	t = template.New(tplname)
	t.Funcs(gRenderTemplateFuncs)
	t, err = t.Parse(tplcontent)

	if err != nil {
		return nil, err
	}
	err = t.Execute(out, data)
	return out.Bytes(), err
}

func RenderByTplFile(tplfile string, data interface{}) ([]byte, error) {
	out := new(bytes.Buffer)

	t, err := template.ParseFiles(tplfile)
	if err != nil {
		return nil, err
	}

	t = t.Funcs(gRenderTemplateFuncs)

	err = t.Execute(out, data)
	return out.Bytes(), err
}

func execmd(cmd *exec.Cmd) error {
	var out bytes.Buffer
	var errOut bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("CMD: %s\nSTDOUT:\n%s\nSTDERR:\n%s\n",
			cmd.Args,
			string(out.Bytes()),
			string(errOut.Bytes()))
	} else {
		fmt.Printf("%s\n", out.String())
		return nil
	}
}

func runGoFile(bootGoFilePath string) error {
	if err := execmd(exec.Command("go", "run", "-a", bootGoFilePath)); err != nil {
		return err
	}

	defer func() {
		os.Remove(bootGoFilePath)
	}()

	return nil
}

func FmtGoFile(gofile string) error {
	return execmd(exec.Command("go", "fmt", "-x", gofile))
}
