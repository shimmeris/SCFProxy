package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/shimmeris/SCFProxy/cmd/config"
	"github.com/shimmeris/SCFProxy/fileutil"
)

const version = "0.1.0"

var debug bool

var rootCmd = &cobra.Command{
	Use:   "scfproxy",
	Short: "scfproxy is a tool that implements multiple proxies based on cloud functions and API gateway functions provided by various cloud providers",
	Long: `
███████╗ ██████╗███████╗██████╗ ██████╗  ██████╗ ██╗  ██╗██╗   ██╗
██╔════╝██╔════╝██╔════╝██╔══██╗██╔══██╗██╔═══██╗╚██╗██╔╝╚██╗ ██╔╝
███████╗██║     █████╗  ██████╔╝██████╔╝██║   ██║ ╚███╔╝  ╚████╔╝ 
╚════██║██║     ██╔══╝  ██╔═══╝ ██╔══██╗██║   ██║ ██╔██╗   ╚██╔╝  
███████║╚██████╗██║     ██║     ██║  ██║╚██████╔╝██╔╝ ██╗   ██║   
╚══════╝ ╚═════╝╚═╝     ╚═╝     ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝   ╚═╝   
	https://github.com/shimmeris/SCFProxy

scfproxy is a tool that implements multiple proxies based on cloud functions and API gateway functions provided by various cloud providers
`,
	Version: version,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !fileutil.PathExists(config.ProviderConfigPath) {
			f, err := os.Create(config.ProviderConfigPath)
			defer f.Close()
			if err != nil {
				return err
			}
			if _, err := f.Write([]byte(config.ProviderConfigContent)); err != nil {
				return err
			}
			logrus.Printf("credential config file has been generated in %s", config.ProviderConfigPath)
		}

		return cmd.Help()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			logrus.SetLevel(logrus.TraceLevel)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "set debug log level")
	rootCmd.PersistentFlags().BoolP("help", "", false, "help for this command")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
