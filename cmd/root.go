// Copyright © 2017 Jack Zampolin <jack.zampolin@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "blockstack-twitter",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.blockstack-twitter.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	RootCmd.PersistentFlags().StringP("search", "s", "Verifying my Blockstack", "blockstack-core node to run rpc comands against")
	RootCmd.PersistentFlags().StringP("consumerKey", "k", "", "blockstack-core node to run rpc comands against")
	RootCmd.PersistentFlags().StringP("consumerSecret", "e", "", "blockstack-core node to run rpc comands against")
	RootCmd.PersistentFlags().StringP("accessToken", "t", "", "blockstack-core node to run rpc comands against")
	RootCmd.PersistentFlags().StringP("accessSecret", "c", "", "blockstack-core node to run rpc comands against")
	RootCmd.PersistentFlags().StringP("port", "p", ":8080", "blockstack-core node to run rpc comands against")
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose mode")
	viper.BindPFlag("search", RootCmd.PersistentFlags().Lookup("search"))
	viper.BindPFlag("consumerKey", RootCmd.PersistentFlags().Lookup("consumerKey"))
	viper.BindPFlag("consumerSecret", RootCmd.PersistentFlags().Lookup("consumerSecret"))
	viper.BindPFlag("accessToken", RootCmd.PersistentFlags().Lookup("accessToken"))
	viper.BindPFlag("accessSecret", RootCmd.PersistentFlags().Lookup("accessSecret"))
	viper.BindPFlag("port", RootCmd.PersistentFlags().Lookup("port"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".blockstack-twitter" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".blockstack-twitter")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
