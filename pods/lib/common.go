package lib

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"os"
	"plugins/cache"
	"plugins/util"
	"regexp"
	"sigs.k8s.io/yaml"
)

// FilterListByJSON 实现按照名称筛选pod Name
func FilterListByJSON(list *v1.PodList) {

	jsonStr, _ := json.Marshal(list)

	//思路:
	//1.先找出所有符合的名称集合
	//2.然后根据名称集合找到相对应的pod

	podSet := []string{}
	isSearch := false
	if SearchPodName != "" {
		isSearch = true
		ret := gjson.Get(string(jsonStr), "items.#.metadata.name")
		for _, pod := range ret.Array() {
			if m, err := regexp.MatchString(SearchPodName, pod.String()); err == nil && m {
				podSet = append(podSet, pod.String())
			}
		}
	}
	if !isSearch {
		return //没有设置搜索，原样返回
	}

	//再通过podSet找到符合的pod,追加到列表中
	podList := []v1.Pod{}
	for _, v := range list.Items {
		if util.InArray(podSet, v.Name) {
			podList = append(podList, v)
		}
	}
	//最后将列表返回
	list.Items = podList
}

// PodsHeader 初始化表格头,将所有需要的表格标题生成一个[]string切片
func PodsHeader() []string {
	headers := []string{"名称", "命名空间", "IP", "状态"}
	if ShowLabels {
		headers = append(headers, "标签")
	}
	return headers
}

// EventsHeader 初始化事件头
var eventHeaders = []string{"事件类型", "REASON", "所属对象", "消息"}

// EventsBody 遍历出事件的消息体
func printEvent(events []*v1.Event) {
	table := tablewriter.NewWriter(os.Stdout)
	//设置头
	table.SetHeader(eventHeaders)
	for _, e := range events {
		podRow := []string{e.Type, e.Reason,
			fmt.Sprintf("%s/%s", e.InvolvedObject.Kind, e.InvolvedObject.Name), e.Message}

		table.Append(podRow)
	}
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
	table.Render()
}

// ------------------------------- 获取k8s 资源 -----------------------------------

// 获取pod详情
func getPodDetail(cmd *cobra.Command) {

	ns := util.GetNameSpace(cmd)

	pods, err := cache.Fact.Core().V1().Pods().Lister().Pods(ns).
		List(labels.Everything())
	if err != nil {
		log.Println("fact coreV1 pods List is failed", err)
		return
	}

	fmt.Println("从缓存取数据")
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
	util.SetTable(table)
	table.Render()
}

// 获取pod详情以yaml的方式输出
func getPodDetailByJson(podName, path string, cmd *cobra.Command) {
	//获取当前命令行输入的命名空间
	ns := util.GetNameSpace(cmd)

	//查询具体的某一个pod
	pod, err := cache.Fact.Core().V1().Pods().Lister().
		Pods(ns).Get(podName)
	if err != nil {
		log.Println("fact pods get is failed", err)
		return
	}

	//Event 事件
	if path == PodEventType {
		//获取informer中的事件信息
		eventList, err := cache.Fact.Core().V1().Events().Lister().List(labels.Everything())
		if err != nil {
			log.Println(err)
			return
		}
		//找到对应pod的事件信息
		podEvents := []*v1.Event{}
		for _, e := range eventList {
			//事件对象
			if e.InvolvedObject.UID == pod.UID {
				podEvents = append(podEvents, e)
			}
		}
		//直接打印
		printEvent(podEvents)
		return
	}

	//Log 日志
	if path == PodLogType {
		req := cache.Client.CoreV1().Pods(ns).GetLogs(pod.Name, &v1.PodLogOptions{})
		ret := req.Do(context.Background())
		b, err := ret.Raw()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(b))
		return
	}

	//将这个pod的数据转换成map 然后变成string
	jsonStr, _ := json.Marshal(pod)
	ret := gjson.Get(string(jsonStr), path)
	if !ret.Exists() {
		log.Println("无法找到对应的内容:" + path)
		return
	}

	//不是对象不是 数组，直接打印
	if !ret.IsObject() && !ret.IsArray() {
		fmt.Println(ret.Raw)
		return
	}

	//最终以map的形式进行反序列化和输出
	tempMap := make(map[string]interface{}, 0)
	err = yaml.Unmarshal([]byte(ret.Raw), &tempMap)
	if err != nil {
		log.Println(err)
		return
	}
	b, _ := yaml.Marshal(tempMap)
	fmt.Println(string(b))
}

// 删除pod
func delPod(args []string, cmd *cobra.Command) {
	if len(args) == 0 {
		log.Println("podName is required")
		return
	}

	//所有的命名空间均来自于 cobra命令行传入
	ns := util.GetNameSpace(cmd)
	err := cache.Client.CoreV1().Pods(ns).Delete(context.Background(), args[0], metav1.DeleteOptions{})
	if err != nil {
		log.Println("delete pod error:", err.Error())
		return
	}
	log.Println("删除POD:", args[0], "成功")
}

// 获取pod的指标列表
func getPodMetrics(ns string) {

	mlist, err := cache.MetricClient.MetricsV1beta1().PodMetricses(ns).
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("pod metrics failed", err)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"名称", "cpu/内存"})

	data := [][]string{}
	for _, p := range mlist.Items {
		for _, c := range p.Containers {
			podRow := []string{}
			if c.Name == "POD" {
				continue
			}
			mem := c.Usage.Memory().Value() / 1024 / 1024
			podRow = append(podRow, p.Name,
				fmt.Sprintf("%s(%sm/%dM)", c.Name, c.Usage.Cpu().String(), mem))
			data = append(data, podRow)
		}

	}
	table.AppendBulk(data)
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	table.Render()
}
