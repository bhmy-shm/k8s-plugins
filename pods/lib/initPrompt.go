package lib

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/labels"
	"log"
	"os"
	"plugins/util"
	"strings"
)

var MyConsoleWrite = prompt.NewStdoutWriter()

//清除屏幕内容
func clearConsole() {
	MyConsoleWrite.EraseScreen()    //清屏
	MyConsoleWrite.CursorGoTo(0, 0) //光标位置0行0列
	MyConsoleWrite.Flush()          //刷新
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
		case "list":
			//informer缓存方式，可以一直获取
			err := cacheCmd.RunE(cmd, args)
			if err != nil {
				log.Fatalln(err)
			}
		//直接查询的方式，但是只能查询一次
		//if err := ListCmd.RunE(cmd, []string{}); err != nil {
		//	log.Fatalln(err)
		//}
		case "get":
			//getPodDetail(args)
			clearConsole()
			runTea(args, cmd)
		case "clear":
			clearConsole()
		}
	}

}

var suggestions = []prompt.Suggest{
	// Command
	{"get", "pod detail"},
	{"list", "pods list"},
	{"exit", "Exit prompt"},
	{"clear", "清除屏幕"},
}

func getPodsList() (ret []prompt.Suggest) {
	pods, err := fact.Core().V1().Pods().Lister().Pods("default").List(labels.Everything())
	if err != nil {
		return
	}
	for _, pod := range pods {
		ret = append(ret, prompt.Suggest{
			Text:        pod.Name,
			Description: "节点:" + pod.Spec.NodeName + "状态：" + string(pod.Status.Phase) + "IP:" + pod.Status.PodIP,
		})
	}
	return
}

func completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	cmd, opt := util.ParseCmd(in.TextBeforeCursor())
	if cmd == "get" {
		return prompt.FilterHasPrefix(getPodsList(), opt, true)
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}

var promptCmd = &cobra.Command{
	Use:          "prompt",
	Short:        "prompt pods ",
	Example:      "kubectl pods prompt",
	SilenceUsage: true,
	RunE: func(c *cobra.Command, args []string) error {
		InitCache() //初始化缓存，在进入交互之后立即初始化
		p := prompt.New(
			executorCmd(c),
			completer,
			prompt.OptionPrefix(">>> "),
			prompt.OptionWriter(MyConsoleWrite), //设置自己的writer
		)
		p.Run()
		return nil
	},
}
