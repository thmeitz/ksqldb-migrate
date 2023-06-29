/*
Copyright © 2021 Thomas Meitz <thme219@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "ksqldb-migrate",
	Short: "ksqldb-migrate migration tool for ksqlDB",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.ksqldb-migrate.yaml)")
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "migration file")
	rootCmd.PersistentFlags().BoolVarP(&preflight, "preflight", "p", true, "preflight migration steps before executing (sql syntax check)")
	rootCmd.PersistentFlags().StringVarP(&logformat, "logformat", "o", "text", "set log format [text|json]")
	rootCmd.PersistentFlags().StringVarP(&loglevel, "loglevel", "v", "debug", "set log level [info|debug|error|trace]")
	rootCmd.PersistentFlags().StringVarP(&host, "host", "t", "http://localhost:8088", "set the ksqldb host")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "set the ksqldb user name")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "w", "", "set the ksqldb user password")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatal(err)
	}

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

		// Search config in home directory with name ".ksqldb-migrate" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ksqldb-migrate")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
