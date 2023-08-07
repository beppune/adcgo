/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "adcgo",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.AllSettings())
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

var cfgFile string

func init() {
	cobra.OnInitialize(loadConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .adcgo.yml)")

}

func loadConfig() {

	if cfgFile != "" {

		viper.SetConfigFile(cfgFile)
		return
	}

	d, _ := os.Getwd()
	viper.AddConfigPath(d)
	viper.SetConfigName(".adcgo")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err.Error())
	}

}
