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
	"github.com/Masterminds/log-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thmeitz/ksqldb-go"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "down reads the migration yaml file and executes the steps",
}

func init() {
	downCmd.Run = down
	rootCmd.AddCommand(downCmd)

	downCmd.Flags().StringP("file", "f", "", "migration file")
	if err := viper.BindPFlag("file", downCmd.Flags().Lookup("file")); err != nil {
		log.Current.Fatal(err)
	}

	if err := downCmd.MarkFlagRequired("file"); err != nil {
		log.Current.Fatal(err)
	}

	downCmd.Flags().BoolP("parse", "p", true, "parse migration steps before executing")
	if err := viper.BindPFlag("parse", downCmd.Flags().Lookup("parse")); err != nil {
		log.Current.Fatal(err)
	}
}

func down(cmd *cobra.Command, args []string) {
	setLogger()

	host := viper.GetString("host")
	user := viper.GetString("username")
	password := viper.GetString("password")

	options := ksqldb.Options{
		Credentials: ksqldb.Credentials{Username: user, Password: password},
		BaseUrl:     host,
		AllowHTTP:   true,
	}

	client, err := ksqldb.NewClient(options, log.Current)
	if err != nil {
		log.Current.Fatal(err)
	}

	// create the DOGS_BY_SIZE table
	// if err := ksqldb.Execute(client,
	// 	`
	// 		CREATE TABLE IF NOT EXISTS DOGS_BY_SIZE AS
	// 			SELECT DOGSIZE AS DOG_SIZE, COUNT(*) AS DOGS_CT
	// 			FROM DOGS WINDOW TUMBLING (SIZE 15 MINUTE)
	// 			GROdown BY DOGSIZE;
	// `); err != nil {
	// 	log.Current.Error(err)
	// 	os.Exit(-1)
	// }
	client.Close()
}
