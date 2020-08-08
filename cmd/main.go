package main

import (
	"fmt"
	"github.com/anthony-dong/aliyun-oss-cli/logger"
	"github.com/anthony-dong/aliyun-oss-cli/oss"
)

func main() {
	client := oss.GetClient()
	bucket := client.ListBucket()
	fmt.Println(bucket)

	err := client.SelectBucket("tyut")
	if err != nil {
		logger.FatalF("%s", err)
	}

	//strings,err := client.ListFile()
	//if err != nil {
	//	panic(err)
	//}
	//for _, e := range strings {
	//	fmt.Println(e)
	//}
	client.GetMod("static/upload/js/vendor.fd1bf7ade4641f022b76.js")
	//kernel := shell.NewKernel()
	//kernel.Register(shell.NewCmd("ls", func(args []string) {
	//	fmt.Println("this is ls")
	//}))
	//kernel.Register(shell.NewCmd("ls2", func(args []string) {
	//	fmt.Println("this is ls")
	//}))
	//
	//sel := shell.NewShell(">")
	//scanner := bufio.NewScanner(os.Stdin)
	//sel.PrintLine()
	//for scanner.Scan() {
	//	args := util.SplitString(scanner.Bytes())
	//	kernel.Run(args)
	//	sel.PrintLine()
	//}
}
