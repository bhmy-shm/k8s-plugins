package lib

import (
	"github.com/spf13/cobra"
	"log"
	"plugins/cache"
)

var ShowLabels bool
var Labels string
var Fields string
var SearchPodName string

func MergeFlags(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		cache.CfgFlags.AddFlags(cmd.Flags())
	}
}

func RunCmd() {
	cmd := &cobra.Command{
		Use:          "kubectl pods [flags]",
		Short:        "list Pods",
		Example:      "kubectl pods [flags]",
		SilenceUsage: true,
	}

	//初始化k8s
	cache.InitClient()

	//合并所有cmd
	MergeFlags(cmd, ListCmd, promptCmd)

	//加入子命令
	cmd.AddCommand(ListCmd, promptCmd)

	//接收所有终端指令
	cmd.Flags().BoolVar(&ShowLabels, "show-labels", false, "kubectl pods --show-labels")
	cmd.Flags().StringVar(&Labels, "labels", "", "kubectl pods --labels=\"app=ngx,version=1\"")
	cmd.Flags().StringVar(&Fields, "fields", "", "kubectl pods --fields=\"status.phase=Running\"")
	cmd.Flags().StringVar(&SearchPodName, "name", "", "kubectl pods --name=\"^myngx\"")

	err := cmd.Execute()
	if err != nil {
		log.Fatalln("execute is failed:", err)
	}
}
