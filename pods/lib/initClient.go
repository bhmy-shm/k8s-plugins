package lib

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"log"
)

var CfgFlags *genericclioptions.ConfigFlags

var Client = initClient()

func initClient() *kubernetes.Clientset {

	CfgFlags = genericclioptions.NewConfigFlags(true)
	config, err := CfgFlags.ToRawKubeConfigLoader().ClientConfig()
	if err != nil {
		log.Fatalln(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

func MergeFlags(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		CfgFlags.AddFlags(cmd.Flags())
	}
}
