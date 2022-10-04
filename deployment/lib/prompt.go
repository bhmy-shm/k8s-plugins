package lib

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"os"
	"plugins/util"
	"strings"
)

var myConsoleWriter = prompt.NewStdoutWriter()

func clearConsole() {
	myConsoleWriter.EraseScreen()
	myConsoleWriter.CursorGoTo(0, 0)
	myConsoleWriter.Flush()
}

func executorCmd(cmd *cobra.Command) func(in string) {
	return func(in string) {
		in = strings.TrimSpace(in)
		blocks := strings.Split(in, " ")
		args := []string{}
		if len(blocks) > 1 {
			args = blocks[1:]
		}
		switch blocks[0] {
		case "exit":
			fmt.Println("Bye!")
			os.Exit(0)
		case "use":
			util.SetNameSpace(args, cmd)
		case "ns":
			fmt.Println("您当前所处的namespace是：", util.GetNameSpace(cmd))
		case "list":
			RenderDeploy(args, cmd)
			//if err := ListDeployCmd.RunE(cmd, []string{}); err != nil {
			//	log.Fatalln(err)
			//}
		case "clear":
			clearConsole()
		}
	}
}

var suggestions = []prompt.Suggest{
	// Command
	{"list", "显示Deployment列表"},
	{"clear", "清除屏幕"},
	{"use", "设置当前namespace,请填写名称"},
	{"ns", "显示当前命名空间"},
	{"exit", "退出交互式窗口"},
}

func completer(c *cobra.Command) func(prompt.Document) []prompt.Suggest {
	return func(in prompt.Document) []prompt.Suggest {
		w := in.GetWordBeforeCursor()
		if w == "" {
			return []prompt.Suggest{}
		}

		cmd, opt := util.ParseCmd(in.TextBeforeCursor())
		if util.InArray([]string{"get"}, cmd) {
			return prompt.FilterHasPrefix(getDeployList(util.GetNameSpace(c)),
				opt, true)
		}

		return prompt.FilterHasPrefix(suggestions, w, true)
	}
}

var promptCmd = &cobra.Command{
	Use:          "prompt",
	Short:        "prompt deployments",
	Example:      "kubectl deployments prompt",
	SilenceUsage: true,
	RunE: func(c *cobra.Command, args []string) error {
		p := prompt.New(
			executorCmd(c),
			completer(c),
			prompt.OptionPrefix(">>> "),
			prompt.OptionWriter(myConsoleWriter),
		)
		p.Run()
		return nil
	},
}
