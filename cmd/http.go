package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/shimmeris/SCFProxy/cmd/config"
	"github.com/shimmeris/SCFProxy/http"
)

var (
	listenAddr string
	certPath   string
	keyPath    string
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Start http proxy",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := config.LoadHttpConfig()
		if err != nil {
			return err
		}

		apis := conf.AvailableApis()
		if len(apis) < 1 {
			return errors.New("available HTTP proxy apis must be at least one")
		}
		opts := &http.Options{
			ListenAddr: listenAddr,
			CertPath:   certPath,
			KeyPath:    keyPath,
			Apis:       apis,
		}
		return http.ServeProxy(opts)
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
	httpCmd.Flags().StringVarP(&listenAddr, "listen", "l", "", "host:port of the proxy")
	httpCmd.Flags().StringVarP(&certPath, "certPath", "c", config.CertPath, "filepath to the CA certificate used to sign MITM certificates")
	httpCmd.Flags().StringVarP(&keyPath, "keyPath", "k", config.KeyPath, "filepath to the private key of the CA used to sign MITM certificates")

	httpCmd.MarkFlagRequired("listen")
}
