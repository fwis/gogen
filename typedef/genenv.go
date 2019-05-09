package typedef

import (
	"fmt"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Env struct {
	GOPATH    string
	ExecDir   string
	OutputDir string
	TmpDir    string
}

func (m *Env) Dump(out io.Writer) {
	fmt.Fprintf(out, `GOPATH=%s
ExecDir=%s
OutputDir=%s
TmpDir=%s
`,
		m.GOPATH,
		m.ExecDir,
		m.OutputDir,
		m.TmpDir,
	)
}

func NewEnv(outputPath string) (m *Env) {
	m = &Env{}
	m.ExecDir, _ = os.Getwd()
	m.TmpDir = os.TempDir()
	m.GOPATH = build.Default.GOPATH

	if outputPath == "" {
		m.OutputDir = m.ExecDir
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
	return
}
