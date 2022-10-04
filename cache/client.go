package cache

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
)

var (
	CfgFlags     *genericclioptions.ConfigFlags
	Client       *kubernetes.Clientset
	RestConfig   *rest.Config
	MetricClient *versioned.Clientset
)

func InitClient() {

	CfgFlags = genericclioptions.NewConfigFlags(true)
	config, err := CfgFlags.ToRawKubeConfigLoader().ClientConfig()
	if err != nil {
		log.Fatalln(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	// 初始化全局变量
	Client = client
	RestConfig = config
	MetricClient = versioned.NewForConfigOrDie(config)
}
