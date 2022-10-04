package lib

import (
	"github.com/spf13/cobra"
	"log"
	"plugins/cache"
)

func MergeFlags(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		cache.CfgFlags.AddFlags(cmd.Flags())
	}
}

var ShowLabels bool
var Labels string

func RunCmd() {
	cmd := &cobra.Command{
		Use:          "kubectl deploy [flags]",
		Short:        "list deploy",
		Example:      "kubectl deploy [flags]",
		SilenceUsage: true,
	}

	cache.InitClient() //初始化k8s client
	cache.InitCache()  //初始化本地 缓存---informer

	//合并主命令的参数
	MergeFlags(cmd, listDeployCmd, promptCmd)

	//加入子命令
	cmd.AddCommand(listDeployCmd, promptCmd)

	cmd.Flags().BoolVar(&ShowLabels, "show-labels", false, "kubectl deploys --show-labels")
	cmd.Flags().StringVar(&Labels, "labels", "", "kubectl deploys --labels=\"app=ngx,version=1\"")

	err := cmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
