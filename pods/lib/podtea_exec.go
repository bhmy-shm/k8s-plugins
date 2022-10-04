package lib

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"os"
	"plugins/cache"
	"plugins/util"
)

// 远程pod shell 交互

func execPod(ns, pod, container string) remotecommand.Executor {
	option := &v1.PodExecOptions{
		Container: container,
		Command:   []string{"sh"},
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}

	req := cache.Client.CoreV1().RESTClient().Post().Resource("pods").
		Namespace(ns).
		Name(pod).
		SubResource("exec").
		Param("color", "false").
		VersionedParams(
			option, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(cache.RestConfig, "POST", req.URL())
	if err != nil {
		panic(err)
	}
	return exec
}

type execModel struct {
	items   []v1.Container
	index   int
	cmd     *cobra.Command
	podName string
	ns      string
}

func (this *execModel) Init() tea.Cmd {
	//根据podName 取出 container 列表
	this.ns = util.GetNameSpace(this.cmd)

	//根据传参获取指定的pod
	pod, err := cache.Client.CoreV1().Pods(this.ns).Get(context.Background(),
		this.podName, metav1.GetOptions{})
	if err != nil {
		return tea.Quit
	}

	//将这个pod传给v1.Containers
	this.items = pod.Spec.Containers
	return nil
}

func (this *execModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgTyp := msg.(type) {
	case tea.KeyMsg:
		switch msgTyp.String() {
		case "ctrl+c", "q":
			return this, tea.Quit
		case "up":
			if this.index > 0 {
				this.index--
			}
		case "down":
			if this.index < len(this.items)-1 {
				this.index++
			}
		case "enter":
			err := execPod(this.ns, this.podName, this.items[this.index].Name).
				Stream(remotecommand.StreamOptions{
					Stdin:  os.Stdin,
					Stdout: os.Stdout,
					Stderr: os.Stderr,
					Tty:    true,
				})
			if err != nil {
				if err != io.EOF {
					log.Println(err)
				}
			}
			return this, tea.Quit
		}
	}
	return this, nil
}

func (this *execModel) View() string {
	var s string
	for i, item := range this.items {
		selected := " "
		if this.index == i {
			selected = "»"
		}
		s = fmt.Sprintf("%s %s(镜像:%s)\n", selected,
			item.Name, item.Image)
	}

	s += "\n按Q退出\n"
	return s
}

func runTeaExec(args []string, cmd *cobra.Command) {
	if len(args) == 0 {
		log.Println("podName is required")
		return
	}
	var execmodel = &execModel{
		cmd:     cmd,
		podName: args[0],
	}
	//
	teaCmd := tea.NewProgram(execmodel)
	if err := teaCmd.Start(); err != nil {
		fmt.Println("start exec is failed", err)
		os.Exit(1)
	}
}
