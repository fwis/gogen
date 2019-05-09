package main

import (
	"gen/typedef"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
)

type ConsoleWriterWrapper struct {
	prompt.ConsoleWriter
}

func (this *ConsoleWriterWrapper) Write(data []byte) (int, error) {
	this.ConsoleWriter.Write(data)
	return len(data), nil
}

type CmdSession struct {
	State  string
	Prefix string
}

var (
	Session        = &CmdSession{}
	GEvn           = typedef.NewEnv("")
	DefaultSuggest = []prompt.Suggest{
		{Text: "env", Description: "显示环境变量"},
		{Text: "mm", Description: "显示models"},
		{Text: "erpmm", Description: "显示erpmm"},
	}

	MMSuggest = []prompt.Suggest{
		{Text: "LadingBill", Description: "提单"},
		{Text: "LadingBillItem", Description: "Store the article text posted by user"},
		{Text: "User", Description: "Store the text commented to articles"},
		{Text: "Dept", Description: "Combine users with specific rules"},
	}
)

func completer(in prompt.Document) []prompt.Suggest {
	if Session.State == "mm" {
		return prompt.FilterHasPrefix(MMSuggest, in.GetWordBeforeCursor(), true)
	} else {
		return prompt.FilterHasPrefix(DefaultSuggest, in.GetWordBeforeCursor(), true)
	}
}

func executor(in string) {
	switch in {
	case "exit", "quit", "q":
		os.Exit(0)
	case "env":
		www.SetColor(prompt.DefaultColor, prompt.DefaultColor, false)
		GEvn.Dump(www)
	case "?":
		if Session.State == "" {
			www.SetColor(prompt.DefaultColor, prompt.DefaultColor, false)
			www.WriteRawStr("当前支持的命令: ")
			www.SetColor(prompt.Red, prompt.DefaultColor, true)
			www.WriteRawStr("mm erpmm\n")
			www.SetColor(prompt.DefaultColor, prompt.DefaultColor, false)
		}
	case "mm":
		//Session.State = "mm"
		//Session.Prefix = Session.State + "> "
		www.SetColor(prompt.DefaultColor, prompt.DefaultColor, false)
		var a []string = []string{"LadingBill", "LadingBillItem", "UserCore", "Dept", "BizGroup", "GeoRegion", "MoGrant", "MoBase"}
		www.WriteRawStr(strings.Join(a, " ") + "\n")
	}

}

func changeLivePrefix() (string, bool) {
	if Session.Prefix == "" {
		return "> ", false
	}
	return Session.Prefix, true
}

var www = &ConsoleWriterWrapper{prompt.NewStdoutWriter()}

func main() {
	p := prompt.New(executor, completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionWriter(www.ConsoleWriter),
		prompt.OptionLivePrefix(changeLivePrefix),
		prompt.OptionTitle("live-prefix-example"),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionDescriptionBGColor(prompt.LightGray),
	)

	p.Run()
}
