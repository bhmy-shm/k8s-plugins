package lib

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"plugins/cache"
	"plugins/util"
)

var listDeployCmd = &cobra.Command{
	Use:          "list",
	Short:        "list deployments",
	Example:      "kubectl deployments list [flags]",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		ns := util.GetNameSpace(cmd)

		//针对deployment的查询
		list, err := cache.Client.AppsV1().Deployments(ns).List(
			context.Background(), v1.ListOptions{
				LabelSelector: Labels,
			})
		if err != nil {
			return err
		}

		//设置表格
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(DeployHeader())

		//设置表格体
		for _, d := range list.Items {
			row := []string{d.Name, d.Namespace,
				fmt.Sprintf("%d/%d", d.Status.ReadyReplicas, d.Status.Replicas),
				fmt.Sprintf("%d", d.Status.AvailableReplicas),
				d.CreationTimestamp.Format("2006-01-02 15:04:05"),
			}
			if ShowLabels {
				row = append(row, util.Map2String(d.Labels))
			}
			table.Append(row)
		}

		//设置表格参数
		util.SetTable(table)
		table.Render()
		return nil
	},
}
