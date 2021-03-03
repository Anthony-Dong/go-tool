package markdown

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/anthony-dong/go-tool/command"
	"github.com/anthony-dong/go-tool/command/api"
	"github.com/anthony-dong/go-tool/command/log"
	"github.com/anthony-dong/go-tool/util"
	"github.com/juju/errors"
	"github.com/urfave/cli/v2"
)

const (
	readmeFileName = "README.md"
)

type markdownCommand struct {
	api.CommonConfig
	Dir          string `json:"dir"`
	TemplateFile string `json:"template"`
}

type readmeFileInfo struct {
	UrlPath string
	Total   int
}

func NewMarkdownCommand() command.Command {
	return new(markdownCommand)
}

func (m *markdownCommand) InitConfig(context *cli.Context, config api.CommonConfig) ([]byte, error) {
	m.CommonConfig = config
	file, err := util.GetFileAbsPath(m.Dir)
	if err != nil {
		return nil, errors.Trace(err)
	}
	m.Dir = file
	return util.ToJsonString(m), nil
}

func (m *markdownCommand) Flag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "dir",
			Aliases:     []string{"d"},
			Destination: &m.Dir,
			Required:    true,
			Usage:       "The markdown project dir",
		},
		&cli.StringFlag{
			Name:        "template",
			Aliases:     []string{"t"},
			Destination: &m.TemplateFile,
			Required:    true,
			Usage:       fmt.Sprintf("The markdown template file path, go template struct: %+v", new(readmeFileInfo)),
		},
	}
}

func (m *markdownCommand) Run(context *cli.Context) error {
	info, err := m.getReadmeFileInfo()
	if err != nil {
		return errors.Trace(err)
	}
	log.Infof("Get ReadmeFileInfo success, UrlPath: Not show, Total: %d", info.Total)

	readmeFile := filepath.Join(m.Dir, readmeFileName)
	file, err := os.OpenFile(readmeFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Trace(err)
	}
	defer file.Close()
	log.Infof("Create %s file success !!", readmeFile)

	parser, err := m.getParser()
	if err != nil {
		return errors.Trace(err)
	}
	log.Infof("New parser success, template file: %s", m.TemplateFile)

	if err := parser.Execute(file, info); err != nil {
		return errors.Trace(err)
	}
	log.Infof("Write README file success !!")
	return nil
}

func (m *markdownCommand) getParser() (*template.Template, error) {
	templateFile, err := os.Open(m.TemplateFile)
	if err != nil {
		log.Errorf("open %s file err: %v", m.TemplateFile, err)
		return nil, errors.Trace(err)
	}
	templateBody, err := ioutil.ReadAll(templateFile)
	if err != nil {
		return nil, errors.Trace(err)
	}
	temp := template.New("readme")
	parse, err := temp.Parse(string(templateBody))
	if err != nil {
		return nil, errors.Trace(err)
	}
	return parse, nil
}

func (m *markdownCommand) getReadmeFileInfo() (*readmeFileInfo, error) {
	builder := strings.Builder{}
	// 获取全部文件
	files, err := util.GetAllFiles(m.Dir, func(fileName string) bool {
		return strings.HasSuffix(fileName, ".md") || strings.HasSuffix(fileName, ".markdown")
	})
	if err != nil {
		return nil, errors.Trace(err)
	}
	// 转成 markdown的url写法
	url := func(file string) string {
		file = strings.TrimPrefix(file, m.Dir)
		return fmt.Sprintf("- [%s](.%s)\n", file, util.Base64Encode(file))
	}
	// 遍历写
	for _, elem := range files {
		_, err := builder.WriteString(url(elem))
		if err != nil {
			return nil, errors.Trace(err)
		}
	}
	return &readmeFileInfo{
		UrlPath: builder.String(),
		Total:   len(files),
	}, nil
}