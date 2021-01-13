package command

import (
	"github.com/anthony-dong/go-tool/logger"
	"github.com/juju/errors"
	"github.com/urfave/cli/v2"
)

var (
	log = logger.NewStdOutLogger(logger.NameOp("[GO-TOOL]"))
)

type Command interface {
	Run(context *cli.Context) error
	Flag() []cli.Flag
	InitConfig(context *cli.Context) ([]byte, error) // 获取配置信息
}

var (
	flag = &cli.StringFlag{
		Name:     "log",
		Usage:    "the log level of print logger",
		Required: false,
		Value:    "debug",
	}
)

func NewCommand(name, desc string, command Command) *cli.Command {
	return &cli.Command{
		Name:  name,
		Usage: desc,
		Flags: append(command.Flag(), flag),
		Action: func(context *cli.Context) error {
			log.SetLevel(context.String("log"))
			jb, err := command.InitConfig(context)
			if err != nil {
				return errors.Trace(err)
			}
			log.Infof("[%s] command load config:\n%s", context.Command.Name, jb)
			return errors.Trace(command.Run(context))
		},
		HelpName: name,
	}
}
