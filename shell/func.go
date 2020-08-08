package shell

import (
	"fmt"
	"github.com/anthony-dong/aliyun-oss-cli/logger"
	"strings"
)

type Operation func(args []string)

type kernel struct {
	cmd         map[string]*Cmd
	cmdFlagList []string
}

type Cmd struct {
	flag string
	op   Operation // 分隔符命令
}

func (this *kernel) Register(cmd *Cmd) {
	if cmd.flag == "" || cmd.op == nil {
		logger.FatalF("cmd flag can not be null")
	}
	_, isExist := this.cmd[cmd.flag]
	if isExist {
		logger.FatalF("cmd flag: %s is exist", cmd.flag)
	}
	this.cmd[cmd.flag] = cmd
	this.cmdFlagList = append(this.cmdFlagList, cmd.flag)
}

func (this *kernel) Run(args []string) {
	fmt.Println(args)
}

func (this *kernel) Help(flag string) []string {
	if flag == "help" {
		return this.cmdFlagList
	}
	// 匹配
	result := make([]string, 0)
	for _, elem := range this.cmdFlagList {
		if strings.HasPrefix(elem, flag) {
			result = append(result, elem)
		}
	}
	return result
}

func NewKernel() *kernel {
	return &kernel{
		cmd:         map[string]*Cmd{},
		cmdFlagList: []string{},
	}
}

func NewCmd(flag string, op Operation) *Cmd {
	return &Cmd{flag: flag, op: op}
}
