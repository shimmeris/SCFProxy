package function

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"encoding/base64"

	"github.com/sirupsen/logrus"
)

type File struct {
	Name string
	Content []byte
	HighPriv bool
}

// compile socks code with `GOOS=linux GOARCH=amd64 go build main.go`
var (
	//go:embed http/tencent.py
	tencentHttpCode []byte
	TencentHttpCodeZip = CreateZipBase64([]File{{Name: "index.py", Content: tencentHttpCode}})

	//go:embed http/alibaba.py
	alibabaHttpCode []byte
	AlibabaHttpCodeZip = CreateZipBase64([]File{{Name: "index.py", Content: alibabaHttpCode}})

	//go:embed http/huawei.py
	huaweiHttpCode []byte
	HuaweiHttpCodeZip =CreateZipBase64([]File{{Name: "index.py", Content: huaweiHttpCode}})

	//go:embed socks/tencent
	tencentSocksCode []byte
	TencentSocksCodeZip = CreateZipBase64([]File{{Name: "index.py", Content: tencentSocksCode, HighPriv: true}})

	//go:embed socks/alibaba
	alibabaSocksCode []byte
	AlibabaSocksCodeZip = CreateZipBase64([]File{{Name: "index.py", Content: alibabaSocksCode, HighPriv: true}})
)



func CreateZipBase64(files []File) string {
	buf := new(bytes.Buffer)

	zw := zip.NewWriter(buf)

	for _, f := range files{
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
	return base64.StdEncoding.EncodeToString(buf.Bytes())

}
