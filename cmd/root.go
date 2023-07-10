/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/NguyenLe1605/gop4mini/pkg/controller"
	"github.com/NguyenLe1605/gop4mini/pkg/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gop4mini",
	Short: "A prototype for a P4 controller",
	Long:  `gop4mini is a P4 runtime controller that can install P4 rule into the bmv2 switches.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p4Info, err := cmd.Flags().GetString("p4info")
		if err != nil {
			return err
		}

		bmv2Json, err := cmd.Flags().GetString("bmv2-json")
		if err != nil {
			return err
		}

		if err := utils.FileExists(p4Info); err != nil {
			return fmt.Errorf("\np4info file not found: %s\nHave you run 'make'?", p4Info)
		}

		if err := utils.FileExists(bmv2Json); err != nil {
			return fmt.Errorf("\nBMv2 JSON file not found: %s\nHave you run 'make'?", bmv2Json)
		}

		if err := controller.Run(p4Info, bmv2Json); err != nil {
			return err
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("p4info", "./build/advanced_tunnel.p4.p4info.txt", "p4info proto in text format from p4c")
	rootCmd.PersistentFlags().String("bmv2-json", "./build/advanced_tunnel.json", "BMv2 JSON file from p4c")
}
