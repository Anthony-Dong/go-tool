package file

import (
	"fmt"
	"path/filepath"

	"github.com/anthony-dong/upload-file-cli/config"
)

type UploadFileInfo struct {
	LocalPath string
	Prefix    string
	FileName  string
	Suffix    string
}

// image/2019-08-29/38564c69-85ba-4415-93d8-cb05e783c4b6.jpg
func (u *UploadFileInfo) GetPutPath(config *config.Config) string {
	return filepath.Clean(fmt.Sprintf("%s/%s/%s%s", config.PathPrefix, u.Prefix, u.FileName, u.Suffix))
}

// https://tyut.oss-accelerate.aliyuncs.com/image/2020/9-14/d21baa6d76a14aa8b70db1c033891990.png
func (u *UploadFileInfo) GetOSSUrl(config *config.Config) string {
	path := u.GetPutPath(config)
	return fmt.Sprintf("https://%s/%s", config.UrlEndpoint, path)
}
