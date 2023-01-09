package cmd

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/shimmeris/SCFProxy/cmd/config"
	"github.com/shimmeris/SCFProxy/sdk"
)

var clearCmd = &cobra.Command{
	Use:       "clear [http|socks|reverse] -p providers -r regions",
	Short:     "Clear deployed module-specific proxies",
	ValidArgs: []string{"http", "socks", "reverse"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		providers, err := createProviders(cmd)
		if err != nil {
			return err
		}
		completely, _ := cmd.Flags().GetBool("completely")

		module := args[0]
		switch module {
		case "http":
			return clearHttp(providers, completely)
		case "socks":
			return clearSocks(providers)
		case "reverse":
			origin, _ := cmd.Flags().GetString("origin")
			if origin == "" {
				return errors.New("missing parameter [-o,--origin]")
			}
			return clearReverse(providers, origin)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)

	clearCmd.Flags().StringSliceP("provider", "p", nil, "specify which cloud providers to clear proxy")
	clearCmd.Flags().StringSliceP("region", "r", nil, "specify which regions of cloud providers clear proxy")
	clearCmd.Flags().StringP("config", "c", config.ProviderConfigPath, "path of provider credential file")

	// clear http needed
	clearCmd.Flags().BoolP("completely", "e", false, "[http] whether to completely clear up deployed proxies (by default only delete triggers)`[http | socks]`")

	// clear reverse needed
	clearCmd.Flags().StringP("origin", "o", "", "[reverset] Address of the reverse proxy back to the source")

	clearCmd.MarkFlagRequired("provider")
	clearCmd.MarkFlagRequired("region")
}

func clearHttp(providers []sdk.Provider, completely bool) error {
	conf, err := config.LoadHttpConfig()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(providers))

	for _, p := range providers {
		go func(p sdk.Provider) {
			defer wg.Done()
			hp := p.(sdk.HttpProxyProvider)
			provider, region := hp.Name(), hp.Region()

			if record, ok := conf.Get(provider, region); ok && record.Api == "" && !completely {
				logrus.Infof("%s %s trigger has already been cleared", provider, region)
				return
			}

			opts := &sdk.FunctionOpts{
				Namespace:    Namespace,
				FunctionName: HTTPFunctionName,
				TriggerName:  HTTPTriggerName,
				OnlyTrigger:  !completely,
			}

			err := hp.ClearHttpProxy(opts)
			if err != nil {
				logrus.Error(err)
				return
			}
			if completely {
				conf.Delete(provider, region)
				logrus.Printf("[success] cleared http function in %s.%s", provider, region)
			} else {
				conf.Set(provider, region, &config.HttpRecord{})
				logrus.Printf("[success] cleared http trigger in %s.%s", provider, region)
			}
		}(p)
	}
	wg.Wait()

	return conf.Save()
}

func clearSocks(providers []sdk.Provider) error {
	conf, err := config.LoadSocksConfig()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(providers))

	for _, p := range providers {
		go func(p sdk.Provider) {
			defer wg.Done()
			sp := p.(sdk.SocksProxyProvider)

			provider, region := sp.Name(), sp.Region()

			opts := &sdk.FunctionOpts{
				Namespace:    Namespace,
				FunctionName: SocksFunctionName,
			}
			err := sp.ClearSocksProxy(opts)
			if err != nil {
				logrus.Error(err)
				return
			}

			conf.Delete(provider, region)
			logrus.Printf("[success] cleared socks function in %s.%s", provider, region)
		}(p)
	}

	wg.Wait()
	return conf.Save()
}

func clearReverse(providers []sdk.Provider, origin string) error {
	conf, err := config.LoadReverseConfig()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, p := range providers {
		i := 0
		for _, record := range conf.Records {
			if record.Provider != p.Name() || record.Region != p.Region() || record.Origin != origin {
				conf.Records[i] = record
				i++
				continue
			}

			wg.Add(1)
			go func(p sdk.Provider, record *config.ReverseRecord) {
				defer wg.Done()

				rp := p.(sdk.ReverseProxyProvider)
				opts := &sdk.ReverseProxyOpts{
					ServiceId: record.ServiceId,
					ApiId:     record.ApiId,
					PluginId:  record.PluginId,
				}
				err := rp.ClearReverseProxy(opts)
				if err != nil {
					logrus.Error(err)
					return
				}

				logrus.Printf("[success] cleard reverse proxy for %s in %s.%s", origin, p.Name(), p.Region())
			}(p, record)
		}
		conf.Records = conf.Records[:i]
	}

	wg.Wait()
	return conf.Save()
}
