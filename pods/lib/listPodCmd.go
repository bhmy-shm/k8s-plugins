package lib

import (
	"context"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"plugins/cache"
	"plugins/util"
)

// ListCmd 交互
var ListCmd = &cobra.Command{
	Use:          "list",
	Short:        "list pods",
	Example:      "kubectl pods list [flags]",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		//执行针对pods的查询
		ns := util.GetNameSpace(cmd)
		list, err := cache.Client.CoreV1().Pods(ns).List(
			context.Background(), v1.ListOptions{
				LabelSelector: Labels,
				FieldSelector: Fields,
			})
		if err != nil {
			return err
		}

		//此函数内部过滤了按照名称查找符合的pod
		FilterListByJSON(list)

		//以上属于查询操作，找到所有符合内容，以下是构建成table表格的形式返回数据
		table := tablewriter.NewWriter(os.Stdout)

		//设置表格头
		table.SetHeader(PodsHeader())

		//设置表格体
		for _, pod := range list.Items {
			//按照行来生成记录
			row := []string{pod.Name, pod.Namespace, pod.Status.PodIP, string(pod.Status.Phase)}
			if ShowLabels {
				row = append(row, util.Map2String(pod.Labels))
			}
			table.Append(row)
		}

		//设置表格参数
		util.SetTable(table)
		table.Render()
		return nil
	},
}
