package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func ResetSTTY() {
	cc := exec.Command("stty", "echo")
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
//  将输入的命令按照空格分开
func ParseCmd(w string) (string, string) {
	w = regexp.MustCompile("\\s+").ReplaceAllString(w, " ")
	l := strings.Split(w, " ")
	if len(l) >= 2 {
		return l[0], strings.Join(l[1:], " ")
	}
	return w, ""
}
