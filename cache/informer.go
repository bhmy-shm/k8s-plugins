package cache

import "k8s.io/client-go/informers"

//使用informer功能必须要实现接口

var Fact informers.SharedInformerFactory

func InitCache() {
	Fact = informers.NewSharedInformerFactoryWithOptions(Client, 0)

	//Pods
	Fact.Core().V1().Pods().Informer().AddEventHandler(NewPodHandler())

	//Event
	Fact.Core().V1().Events().Informer().AddEventHandler(NewPodHandler())

	//Deployment
	Fact.Apps().V1().Deployments().Informer().AddEventHandler(NewDeployHandler())

	//启动一个channel接收informer数据
	ch := make(chan struct{})
	Fact.Start(ch)
	Fact.WaitForCacheSync(ch)
}
