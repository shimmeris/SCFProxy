package function

import _ "embed"

// compile socks code with `GOOS=linux GOARCH=amd64 go build main.go`

var (
	//go:embed http/tencent.py
	TencentHttpCode []byte

	//go:embed http/alibaba.py
	AlibabaHttpCode []byte

	//go:embed http/huawei.py
	HuaweiHttpCode []byte

	//go:embed socks/tencent
	TencentSocksCode []byte

	//go:embed socks/alibaba
	AlibabaSocksCode []byte
)
