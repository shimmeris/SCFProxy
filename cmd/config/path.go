package config

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var (
	configPath         = filepath.Join(userHomeDir(), ".config/scfproxy")
	CertPath           = filepath.Join(configPath, "cert/scfproxy.cer")
	KeyPath            = filepath.Join(configPath, "cert/scfproxy.key")
	HttpProxyPath      = filepath.Join(configPath, "http.json")
	SocksProxyPath     = filepath.Join(configPath, "socks.json")
	ReverseProxyPath   = filepath.Join(configPath, "reverse.json")
	ProviderConfigPath = filepath.Join(configPath, "sdk.toml")
)

func init() {
	os.MkdirAll(filepath.Join(configPath, "cert"), os.ModePerm)
}

func userHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		logrus.Fatal("Could not get user home directory: %s\n", err)
	}
	return usr.HomeDir
}
