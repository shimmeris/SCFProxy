package cmd

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/shimmeris/SCFProxy/cmd/config"
	"github.com/shimmeris/SCFProxy/sdk"
	"github.com/shimmeris/SCFProxy/socks"
)

var socksCmd = &cobra.Command{
	Use:   "socks",
	Short: "Start socks proxy",
	RunE: func(cmd *cobra.Command, args []string) error {
		lp, _ := cmd.Flags().GetString("lp")
		sp, _ := cmd.Flags().GetString("sp")
		host, _ := cmd.Flags().GetString("host")
		auth, _ := cmd.Flags().GetString("auth")
		key := randomString(socks.KeyLength)

		var wg sync.WaitGroup
		go func() {
			wg.Add(1)
			defer wg.Done()
			socks.Serve(lp, sp, key)
		}()

		providerConfigPath, _ := cmd.Flags().GetString("config")
		message := &Message{
			Key:  key,
			Addr: fmt.Sprintf("%s:%s", host, sp),
			Auth: auth,
		}
		invoke(providerConfigPath, message.Json())

		wg.Wait()
		return nil
	},
}

func invoke(providerConfigPath, message string) {
	providerConfig, err := config.LoadProviderConfig(providerConfigPath)
	if err != nil {
		logrus.Fatalf("Loading provider config failed")
	}

	conf, err := config.LoadSocksConfig()
	if err != nil {
		logrus.Fatalf("Loading socks config failed")
	}

	for provider, rmap := range conf.Records {
		for region := range rmap {
			go func(provider, region string) {
				p, err := createProvider(provider, region, providerConfig)
				if err != nil {
					logrus.Error(err)
					return
				}
				sp, ok := p.(sdk.SocksProxyProvider)
				if !ok {
					logrus.Errorf("%s can't deploy reverse proxy", provider)
					return
				}

				opts := &sdk.FunctionOpts{
					Namespace:    Namespace,
					FunctionName: SocksFunctionName,
				}
				err = sp.InvokeFunction(opts, message)
				if err != nil {
					logrus.Error(err)
				}
			}(provider, region)
		}
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())

	rootCmd.AddCommand(socksCmd)
	socksCmd.Flags().StringP("lp", "l", "", "listen port for client socks5 connection")
	socksCmd.Flags().StringP("sp", "s", "", "listen port for cloud function's connection")
	socksCmd.Flags().StringP("host", "h", "", "host:port address of the cloud function callback")
	socksCmd.Flags().StringP("config", "c", config.ProviderConfigPath, "path of provider credential file")
	socksCmd.Flags().String("auth", "", "username:password for socks proxy authentication")
	socksCmd.MarkFlagRequired("lp")
	socksCmd.MarkFlagRequired("sp")
	socksCmd.MarkFlagRequired("host")

}

type Message struct {
	Key  string
	Addr string
	Auth string
}

func (m *Message) Json() string {
	b, _ := json.Marshal(m)
	return string(b)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
