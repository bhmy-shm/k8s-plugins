package lib

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"plugins/pods/lib"
)

var ListDeployCmd = &cobra.Command{
	Use:          "list",
	Short:        "list deployments",
	Example:      "kubectl deployments list [flags]",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ns, err := cmd.Flags().GetString("namespace")
		if err != nil {
			return err
		}
		if ns == "" {
			ns = "default"
		}

		//针对deployment的查询
		list, err := lib.Client.AppsV1().Deployments(ns).List(
			context.Background(), v1.ListOptions{
				LabelSelector: lib.Labels,
				FieldSelector: lib.Fields,
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
				fmt.Sprintf("%d", d.Status.UpdatedReplicas),
				fmt.Sprintf("%d", d.Status.AvailableReplicas),
				d.CreationTimestamp.String(),
			}
			table.Append(row)
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

		table.Render()
		return nil
	},
}
