package fileutil

import (
	"archive/zip"
	"bytes"
	"encoding/base64"

	"github.com/sirupsen/logrus"
)

func CreateZipBase64(filename string, content []byte) string {
	buf := new(bytes.Buffer)

	zw := zip.NewWriter(buf)
	fw, err := zw.CreateHeader(&zip.FileHeader{
		CreatorVersion: 3 << 8,     // indicates Unix
		ExternalAttrs:  0777 << 16, // -rwxrwxrwx file permissions
		Name:           filename,
		Method:         zip.Deflate,
	})
	if err != nil {
		logrus.Error(err)
	}
	_, err = fw.Write(content)
	if err != nil {
		logrus.Error(err)
	}
	//fw, err := zw.Create(filename)
	//if err != nil {
	//	logrus.Error(err)
	//}
	//n, err := fw.Write(content)
	//if err != nil {
	//	logrus.Error(n, err)
	//}

	zw.Close()
	return base64.StdEncoding.EncodeToString(buf.Bytes())

}
