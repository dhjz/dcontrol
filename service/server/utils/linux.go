//go:build linux

package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetCmdOutput(cmd *exec.Cmd, isCombine bool) (string, error) {
	var output []byte
	var err error
	if isCombine {
		output, err = cmd.CombinedOutput()
	} else {
		output, err = cmd.Output()
	}
	if err != nil {
		fmt.Println("执行cmd命令失败:", err)
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func OpenBrowser(url string) error {
	return nil
}

func GenTaskBarIcon() {

}
