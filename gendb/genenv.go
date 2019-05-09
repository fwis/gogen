package gendb

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
)

type Env struct {
	GOPATH       string
	CurDir       string
	OutputDir    string
	TmpDir       string
	TemplatePath string
	//EscapeOutputDir    string
	//EscapeCurDir       string
}

func (m *Env) PrintDebug() {
	fmt.Printf(`
GOPATH=%v
curdir=%s
outputdir=%s
tmpdir=%s
template=%s
`,
		m.GOPATH,
		m.CurDir,
		m.OutputDir,
		m.TmpDir,
		m.TemplatePath,
	)
}

func NewEnv(outputPath string) (m *Env) {
	m = &Env{}
	m.CurDir, _ = os.Getwd()
	//m.EscapeCurDir = strings.Replace(m.CurDir, "\\", "\\\\", -1) //using  EscapePath in template
	m.TmpDir = os.TempDir()
	m.GOPATH = build.Default.GOPATH

	if outputPath == "" {
		m.OutputDir = m.CurDir
	} else {
		var err error
		outputdir := outputPath
		if strings.HasPrefix(outputdir, "src") {
			outputdir = filepath.Join(m.GOPATH, outputdir)
		} else if strings.HasPrefix(outputdir, ".") || strings.HasPrefix(outputdir, "..") {
			outputdir, err = filepath.Abs(outputdir)
			if err != nil {
				fmt.Printf("err=%v, outputdir=%s\n", err, outputdir)
			}
		} else if filepath.IsAbs(outputdir) {
			outputdir, _ = filepath.Abs(outputdir)
		} else {
			//
		}

		if outputdir != "" {
			m.OutputDir = outputdir
		}
	}

	//m.EscapeOutputDir = strings.Replace(m.OutputDir, "\\", "\\\\", -1) //using  EscapePath in template
	m.TemplatePath = filepath.Join(m.GOPATH, "src", "gen", "gendb", "gendb.template")
	return
}
