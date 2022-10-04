package cache

type PodHandler struct{}

func NewPodHandler() *PodHandler {
	return &PodHandler{}
}
func (this *PodHandler) OnAdd(obj interface{})               {}
func (this *PodHandler) OnUpdate(oldObj, newObj interface{}) {}
func (this *PodHandler) OnDelete(obj interface{})            {}

//

type DeployHandler struct{}

func NewDeployHandler() *DeployHandler {
	return &DeployHandler{}
}
func (this *DeployHandler) OnAdd(obj interface{})               {}
func (this *DeployHandler) OnUpdate(oldObj, newObj interface{}) {}
func (this *DeployHandler) OnDelete(obj interface{})            {}
