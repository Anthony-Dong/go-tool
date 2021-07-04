package hexo

import (
	"fmt"

	"github.com/anthony-dong/go-tool/command"
	"github.com/anthony-dong/go-tool/command/api"
	"github.com/anthony-dong/go-tool/commons/codec/gjson"
	"github.com/anthony-dong/go-tool/commons/gfile"
	"github.com/juju/errors"
	"github.com/urfave/cli/v2"
)

type hexoCommand struct {
	api.CommonConfig
	Dir         string          `json:"dir"`
	Keyword     cli.StringSlice `json:"-"`
	ShowKeyword []string        `json:"keyword"`
	TargetDir   string          `json:"target_dir"`
}

func NewCommand() command.Command {
	return new(hexoCommand)
}

func (m *hexoCommand) InitConfig(context *cli.Context, config api.CommonConfig) ([]byte, error) {
	m.CommonConfig = config
	var err error
	m.Dir, err = gfile.GetFileAbsPath(m.Dir)
	if err != nil {
		return nil, errors.Trace(err)
	}
	m.TargetDir, err = gfile.GetFileAbsPath(m.TargetDir)
	if err != nil {
		return nil, errors.Trace(err)
	}
	m.ShowKeyword = m.Keyword.Value()
	return gjson.ToJsonString(m), nil
}

func (m *hexoCommand) Flag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "dir",
			Aliases:     []string{"d"},
			Destination: &m.Dir,
			Required:    true,
			Usage:       "The markdown project dir",
		},
		&cli.StringFlag{
			Name:        "target_dir",
			Aliases:     []string{"t"},
			Destination: &m.TargetDir,
			Required:    true,
			Usage:       fmt.Sprintf("The hexo post dir"),
		},
		&cli.StringSliceFlag{
			Name:        "keyword",
			Aliases:     []string{"k"},
			Destination: &m.Keyword,
			Required:    false,
			Usage:       fmt.Sprintf("The keyword need clear, eg: baidu => xxxxx"),
		},
	}
}

func (m *hexoCommand) Run(ctx *cli.Context) error {
	if err := Run(ctx.Context, m.Dir, m.TargetDir, m.Keyword.Value()); err != nil {
		return err
	}
	return nil
}
