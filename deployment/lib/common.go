package lib

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"log"
	"os"
	"plugins/cache"
	"plugins/util"
	"sort"
)

func DeployHeader() []string {
	res := []string{"名称", "命名空间", "READY", "有效副本数", "起始日期"}
	if ShowLabels {
		res = append(res, "标签")
	}
	return res
}

func listByNs(ns string) []*v1.Deployment {
	list, err := cache.Fact.Apps().V1().Deployments().Lister().Deployments(ns).
		List(labels.Everything())
	if err != nil {
		log.Println(err)
		return nil
	}
	//用创建时间进行排序
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreationTimestamp.String() > list[j].CreationTimestamp.String()
	})
	return list
}

func getDeployList(ns string) (ret []prompt.Suggest) {
	depList := listByNs(ns)
	if depList == nil {
		return
	}
	for _, dep := range depList {
		ret = append(ret, prompt.Suggest{
			Text: dep.Name,
			Description: fmt.Sprintf("副本:%d/%d", dep.Status.Replicas,
				dep.Status.Replicas),
		})
	}
	return
}

// RenderDeploy 渲染 deploys 列表
func RenderDeploy(args []string, cmd *cobra.Command) {
	depList := listByNs(util.GetNameSpace(cmd))
	if depList == nil {
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	//设置头
	table.SetHeader([]string{"名称"})
	for _, dep := range depList {
		depRow := []string{dep.Name}

		table.Append(depRow)
	}
	util.SetTable(table)
	table.Render()
}
