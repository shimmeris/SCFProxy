package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/shimmeris/SCFProxy/socks"
)

var socksCmd = &cobra.Command{
	Use:   "socks",
	Short: "Start socks proxy",
	RunE: func(cmd *cobra.Command, args []string) error {
		lp, _ := cmd.Flags().GetString("lp")
		sp, _ := cmd.Flags().GetString("sp")
		key, _ := cmd.Flags().GetString("key")
		if len(key) != socks.KeyLength {
			return errors.New(fmt.Sprintf("key must be %d bytes", socks.KeyLength))
		}
		socks.Serve(lp, sp, key)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(socksCmd)

	socksCmd.Flags().StringP("lp", "l", "", "listen port for client socks5 connection")
	socksCmd.Flags().StringP("sp", "s", "", "listen port for cloud function's connection")
	socksCmd.Flags().StringP("key", "k", "", "8-bytes string used to verify that the connection initiated to [-s port] is from the cloud function")
	socksCmd.MarkFlagRequired("lp")
	socksCmd.MarkFlagRequired("sp")
	socksCmd.MarkFlagRequired("key")

}
