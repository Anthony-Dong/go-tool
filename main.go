package main

import (
	"github.com/anthony-dong/go-tool/command/wrk"
	"os"
	"path/filepath"
	"sort"

	logger2 "github.com/anthony-dong/go-tool/commons/logger"

	"github.com/anthony-dong/go-tool/command/hexo"

	"github.com/anthony-dong/go-tool/command/markdown"

	"github.com/anthony-dong/go-tool/command/api"

	"github.com/anthony-dong/go-tool/command"
	"github.com/anthony-dong/go-tool/command/upload"
	"github.com/juju/errors"
	"github.com/urfave/cli/v2"
)

var (
	log = logger2.NewStdOutLogger(logger2.NameOp("[GO-TOOL]"))
)

func main() {
	//os.Args = []string{os.Args[0], "-h"}
	//os.Args = []string{os.Args[0], "-v"}
	//os.Args = []string{os.Args[0], "upload", "-h"}
	//os.Args = []string{os.Args[0], "upload", "--log", "fatal", "-f", "./go.mod"}
	//os.Args = []string{os.Args[0], "--config", "/Users/fanhaodong/go/bin/upload-config.json", "upload", "-f", "./go.mod", "-d", "base64"}
	//os.Args = []string{os.Args[0], "--config", "/Users/fanhaodong/go/bin/upload-config.json", "markdown", "-h"}
	//os.Args = []string{os.Args[0], "--config", "/Users/fanhaodong/note/note/.config/go-tool.json", "markdown", "-d", "/Users/fanhaodong/note/note", "-t", "/Users/fanhaodong/note/note/.config/README-template.md", "-i", "/hexo-home"}
	//os.Args = []string{os.Args[0], "hexo", "-h"}
	//os.Args = []string{os.Args[0], "hexo", "-d", "/Users/fanhaodong/go/code/go-tool/test/test", "-t", "/Users/fanhaodong/go/code/go-tool/test/post", "-k", "baidu", "-k", "alibaba"}
	//os.Args = []string{os.Args[0], "wrk", "-h"}
	//os.Args = []string{os.Args[0], "wrk", "-d", "1s", "-t","5","-u", "http://localhost:8888"}
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
		Flags:        append([]cli.Flag{}, api.GlobalFlag...),
		Commands: []*cli.Command{
			command.NewCommand("upload", "文件上传工具", upload.NewUploadCommand()),
			command.NewCommand("markdown", "生成markdown项目的README文件", markdown.NewMarkdownCommand()),
			command.NewCommand("hexo", "生成hexo项目的Markdown文件", hexo.NewCommand()),
			command.NewCommand("wrk", "压测接口工具", wrk.NewCommand()),
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("run cli find err:\n%s", errors.ErrorStack(err))
	}
}
