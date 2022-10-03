package lib

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"os"
	"plugins/util"
)

//使用informer功能必须要实现接口

type PodHandler struct {
}

func NewPodHandler() *PodHandler {
	return &PodHandler{}
}

func (this *PodHandler) OnAdd(obj interface{})               {}
func (this *PodHandler) OnUpdate(oldObj, newObj interface{}) {}
func (this *PodHandler) OnDelete(obj interface{})            {}

var fact informers.SharedInformerFactory

func InitCache() {
	fact = informers.NewSharedInformerFactoryWithOptions(Client, 0)

	//Pods
	fact.Core().V1().Pods().Informer().AddEventHandler(NewPodHandler())

	//Event
	fact.Core().V1().Events().Informer().AddEventHandler(NewPodHandler())

	//启动一个channel接收informer数据
	ch := make(chan struct{})
	fact.Start(ch)
	fact.WaitForCacheSync(ch)
}

var cacheCmd = &cobra.Command{
	Use:    "cache",
	Short:  "pods by cache",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ns, err := cmd.Flags().GetString("namespace")
		if err != nil {
			return err
		}
		if ns == "" {
			ns = "default"
		}

		//注意这里是从informer中读取列表数据，而不再是clientSet
		pods, err := fact.Core().V1().Pods().Lister().Pods(ns).List(labels.Everything())
		if err != nil {
			return err
		}
		fmt.Println("从缓存取数据成功")
		table := tablewriter.NewWriter(os.Stdout)

		//设置头
		table.SetHeader(PodsHeader())
		for _, pod := range pods {
			podRow := []string{pod.Name, pod.Namespace, pod.Status.PodIP, string(pod.Status.Phase)}
			if ShowLabels {
				podRow = append(podRow, util.Map2String(pod.Labels))
			}
			table.Append(podRow)
		}
		//设置表格参数
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetTablePadding("\t") // pad with tabs
		table.SetNoWhiteSpace(true)

		//执行表格
		table.Render()
		return nil
	},
}
