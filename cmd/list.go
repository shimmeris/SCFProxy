package cmd

import (
	"errors"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/shimmeris/SCFProxy/cmd/config"
)

var listCmd = &cobra.Command{
	Use:   "list [provider|region|http|socks|reverse] [flags]",
	Short: "Display all kinds of data",
	Long: "Display all kinds of data\n" +
		"`list provider` accepts `-m module` flag to filter out providers for a specific module\n" +
		"`list region` accepts `-p providers` flag to specify which providers supported regions to view\n" +
		"remain arguments like `http`, `socks`, `reverse` are used to list the proxies that are currently deployed",
	ValidArgs: []string{"provider", "region", "http", "socks", "reverse"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoMergeCells(true)
		table.SetRowLine(true)

		switch args[0] {
		case "provider":
			logrus.Debug("test")
			module, _ := cmd.Flags().GetString("module")
			table.Append(listProviders(module))
			table.Render()

			return nil
		case "region":
			provider, _ := cmd.Flags().GetStringSlice("provider")
			if provider == nil {
				return errors.New("missing parameter [-p/--provider]")
			}

			data := [][]string{}
			for _, p := range provider {
				for _, r := range listRegions(p) {
					data = append(data, []string{p, r})
				}
			}
			table.AppendBulk(data)
			table.Render()
		case "http":
			table.SetHeader([]string{"Provider", "Region", "Api"})
			conf, err := config.LoadHttpConfig()
			if err != nil {
				return err
			}
			data := conf.ToDoubleArray()
			table.AppendBulk(data)
			table.Render()
		case "socks":
			table.SetHeader([]string{"Provider", "Region", "Host", "Port", "Key"})
			conf, err := config.LoadSocksConfig()
			if err != nil {
				return err
			}
			data := conf.ToDoubleArray()
			table.AppendBulk(data)
			table.Render()
		case "reverse":
			table.SetHeader([]string{"Provider", "Region", "Origin", "Api"})
			conf, err := config.LoadReverseConfig()
			if err != nil {
				return err
			}
			data := conf.ToDoubleArray()
			table.AppendBulk(data)
			table.Render()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("module", "m", "", "filter out providers for a specific module `[provider]`")
	listCmd.Flags().StringSliceP("provider", "p", nil, "specify which providers supported regions to view `[region]`")
}
