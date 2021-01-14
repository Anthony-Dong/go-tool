package main

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/anthony-dong/go-tool/command"
	"github.com/anthony-dong/go-tool/command/upload"
	"github.com/anthony-dong/go-tool/logger"
	"github.com/juju/errors"
	"github.com/urfave/cli/v2"
)

var (
	log = logger.NewStdOutLogger(logger.NameOp("[GO-TOOL]"))
)

func main() {
	//os.Args = []string{os.Args[0], "-h"}
	//os.Args = []string{os.Args[0], "-v"}
	//os.Args = []string{os.Args[0], "upload", "-h"}
	//os.Args = []string{os.Args[0], "upload", "--log", "fatal", "-f", "./go.mod"}
	//os.Args = []string{os.Args[0], "upload", "-f", "./go.mod", "-d", "base64"}
	app := &cli.App{
		Name:         filepath.Base(os.Args[0]),
		HelpName:     filepath.Base(os.Args[0]),
		Usage:        "A go toll cli application for go",
		UsageText:    "",
		BashComplete: cli.DefaultAppComplete,
		Reader:       os.Stdin,
		Writer:       os.Stdout,
		ErrWriter:    os.Stdout,
		Version:      "v1.0.0",
		Flags:        []cli.Flag{},
		Commands: []*cli.Command{
			command.NewCommand("upload", "文件上传工具", upload.NewUploadCommand()),
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("run cli find err, err: %s", errors.ErrorStack(err))
	}
}
