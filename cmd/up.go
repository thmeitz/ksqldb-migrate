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
	"errors"
	"os"

	"github.com/Masterminds/log-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thmeitz/ksqldb-go"
	"github.com/thmeitz/ksqldb-migrate/internal"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "up reads the migration yaml file and executes the steps",
}

func init() {
	upCmd.Run = up
	rootCmd.AddCommand(upCmd)

	upCmd.Flags().StringP("file", "f", "", "migration file")
	if err := viper.BindPFlag("file", upCmd.Flags().Lookup("file")); err != nil {
		log.Current.Fatal(err)
	}

	if err := upCmd.MarkFlagRequired("file"); err != nil {
		log.Current.Fatal(err)
	}

	upCmd.Flags().BoolP("preflight", "p", false, "preflight migration steps before executing (sql syntax check)")
	if err := viper.BindPFlag("preflight", upCmd.Flags().Lookup("preflight")); err != nil {
		log.Current.Fatal(err)
	}
}

func up(cmd *cobra.Command, args []string) {
	setLogger()

	host := viper.GetString("host")
	user := viper.GetString("username")
	password := viper.GetString("password")
	file := viper.GetString("file")
	preflight := viper.GetBool("preflight")

	options := ksqldb.Options{
		Credentials: ksqldb.Credentials{Username: user, Password: password},
		BaseUrl:     host,
		AllowHTTP:   true,
	}

	client, err := ksqldb.NewClient(options, log.Current)
	if err != nil {
		log.Current.Fatal(err)
	}

	log.Current.Debugw("file", log.Fields{"filename": file})

	migrate, err := internal.NewMigration(file)
	if err != nil {
		log.Current.Fatal(err)
	}

	for idx, step := range migrate.Up {
		currentIndex := idx + 1
		log.Current.Infow("processing", log.Fields{"step": currentIndex, "name": step.Name})
		if preflight {
			if err := client.ParseKSQL(step.Exec); err != nil {
				log.Fatalf("error in step:%v, %v", currentIndex, errors.Unwrap(err))
			}

			log.Current.Infow("preflight check", log.Fields{"step": currentIndex, "name": step.Name, "status": "ok"})
			if err := ksqldb.Execute(client, step.Exec); err != nil {
				log.Current.Error(err)
				os.Exit(-1)
			}
			log.Current.Infow("processed", log.Fields{"status": "ok", "step": currentIndex, "name": step.Name})
		}
	}

	client.Close()
}
