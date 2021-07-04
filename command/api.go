package command

import (
	"github.com/anthony-dong/go-tool/command/api"
	"github.com/anthony-dong/go-tool/command/log"
	"github.com/juju/errors"
	"github.com/urfave/cli/v2"
)

/**
所有的command 都需要实现这个.
*/
type Command interface {
	Run(context *cli.Context) error
	Flag() []cli.Flag
	InitConfig(context *cli.Context, config api.CommonConfig) ([]byte, error) // 获取配置信息
}

/**
创建 command.
*/
func NewCommand(name, desc string, command Command) *cli.Command {
	return &cli.Command{
		Name:  name,
		Usage: desc,
		Flags: command.Flag(),
		Action: func(context *cli.Context) error {
			config := api.GetCommonConfig(context)
			log.SetLevel(config.LogLevel)
			jb, err := command.InitConfig(context, config)
			if err != nil {
				return errors.Trace(err)
			}
			log.Infof("[%s] command load config:\n%s", context.Command.Name, jb)
			return errors.Trace(command.Run(context))
		},
		HelpName: name,
	}
}
