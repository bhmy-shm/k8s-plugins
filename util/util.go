package util

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func ResetSTTY() {
	cc := exec.Command("stty", "-F", "/dev/tty", "echo")
	cc.Stdout = os.Stdout
	cc.Stderr = os.Stderr
	if err := cc.Run(); err != nil {
		log.Println(err)
	}
}

func InArray(arr []string, item string) bool {
	for _, p := range arr {
		if p == item {
			return true
		}
	}
	return false
}

// 将Map转换成string

func Map2String(data map[string]string) (ret string) {
	for k, v := range data {
		ret += fmt.Sprintf("%s=%s\n", k, v)
	}
	return
}

// ParseCmd 正则匹配，将多个空格替换成一个空格
//
//	将输入的命令按照空格分开
func ParseCmd(w string) (string, string) {
	w = regexp.MustCompile("\\s+").ReplaceAllString(w, " ")
	l := strings.Split(w, " ")
	if len(l) >= 2 {
		return l[0], strings.Join(l[1:], " ")
	}
	return w, ""
}

// SetTable 设置table的样式，不重要 。看看就好
func SetTable(table *tablewriter.Table) {
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
}

//

// 获取当前命令行输入的命名空间

const DefaultNameSpace = "default"

func GetNameSpace(cmd *cobra.Command) string {
	//从传入的命令中获取
	ns, err := cmd.Flags().GetString("namespace")
	if err != nil {
		log.Println("error ns param")
		return DefaultNameSpace
	}
	if ns == "" {
		ns = DefaultNameSpace
	}
	return ns
}

func ShowNameSpace(cmd *cobra.Command) {
	ns := GetNameSpace(cmd)
	fmt.Println("您当前所处的namespace是：", ns)
}

func SetNameSpace(args []string, cmd *cobra.Command) {
	if len(args) == 0 {
		log.Println("namespace name is required")
		return
	}
	//将当前命令行中的namespace 参数进行替换
	err := cmd.Flags().Set("namespace", args[0])
	if err != nil {
		log.Println("设置namespace失败:", err.Error())
		return
	}
	fmt.Println("设置namespace成功")
}
