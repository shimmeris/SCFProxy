package function

import (
	"archive/zip"
	"bytes"
	"embed"
	"encoding/base64"
	"io/fs"
	"strings"

	"github.com/sirupsen/logrus"
)

type File struct {
	Name     string
	Content  []byte
	HighPriv bool
}

// compile socks code with `GOOS=linux GOARCH=amd64 go build main.go`
var (
	//go:embed http/tencent.py
	tencentHttpCode    []byte
	TencentHttpCodeZip = CreateZipBase64([]File{{Name: "index.py", Content: tencentHttpCode}})

	//go:embed http/alibaba.py
	alibabaHttpCode    []byte
	AlibabaHttpCodeZip = CreateZipBase64([]File{{Name: "index.py", Content: alibabaHttpCode}})

	//go:embed http/huawei.py
	huaweiHttpCode    []byte
	HuaweiHttpCodeZip = CreateZipBase64([]File{{Name: "index.py", Content: huaweiHttpCode}})

	//go:embed http/aws.py
	awsHttpCode    []byte
	AwsHttpCodeZip = awsHttpCodeZip()

	//go:embed socks/tencent
	tencentSocksCode    []byte
	TencentSocksCodeZip = CreateZipBase64([]File{{Name: "main", Content: tencentSocksCode, HighPriv: true}})

	//go:embed socks/alibaba
	alibabaSocksCode    []byte
	AlibabaSocksCodeZip = CreateZipBase64([]File{{Name: "main", Content: alibabaSocksCode, HighPriv: true}})

	//go:embed socks/aws
	awsSocksCode    []byte
	AwsSocksCodeZip = CreateZip([]File{{Name: "main", Content: awsSocksCode, HighPriv: true}})

	//go:embed http/package
	urllib3 embed.FS
)

func awsHttpCodeZip() []byte {
	// aws Python runtime does not have urllib3 dependency, need to be uploaded along with the code
	files := []File{}
	fs.WalkDir(urllib3, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		content, err := fs.ReadFile(urllib3, path)
		files = append(files, File{Name: strings.SplitN(path, "/", 3)[2], Content: content})
		return nil
	})

	files = append(files, File{Name: "index.py", Content: awsHttpCode})
	return CreateZip(files)
}

func CreateZip(files []File) []byte {
	buf := new(bytes.Buffer)

	zw := zip.NewWriter(buf)

	for _, f := range files {
		if f.HighPriv {
			fw, err := zw.CreateHeader(&zip.FileHeader{
				CreatorVersion: 3 << 8,     // indicates Unix
				ExternalAttrs:  0777 << 16, // -rwxrwxrwx file permissions
				Name:           f.Name,
				Method:         zip.Deflate,
			})
			if err != nil {
				logrus.Error(err)
			}

			_, err = fw.Write(f.Content)
			if err != nil {
				logrus.Error(err)
			}
		} else {
			fw, err := zw.Create(f.Name)
			if err != nil {
				logrus.Error(err)
			}
			_, err = fw.Write(f.Content)
			if err != nil {
				logrus.Error(err)
			}
		}
	}

	zw.Close()
	return buf.Bytes()
}

func CreateZipBase64(files []File) string {
	b := CreateZip(files)
	return base64.StdEncoding.EncodeToString(b)
}
