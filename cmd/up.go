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
	"github.com/thmeitz/ksqldb-go"
	"github.com/thmeitz/ksqldb-go/net"
	"github.com/thmeitz/ksqldb-go/parser"
	"github.com/thmeitz/ksqldb-migrate/internal"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "up reads the migration yaml file and executes the steps",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := rootCmd.MarkPersistentFlagRequired("file"); err != nil {
			log.Current.Fatal(err)
		}
		return nil
	},
}

func init() {
	upCmd.Run = up
	rootCmd.AddCommand(upCmd)
}

func up(cmd *cobra.Command, args []string) {
	setLogger()

	options := net.Options{
		Credentials: net.Credentials{Username: username, Password: password},
		BaseUrl:     host,
		AllowHTTP:   true,
	}

	client, err := ksqldb.NewClientWithOptions(options)
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
			if err := parser.ParseSql(step.Exec); err != nil {
				log.Fatalf("error in step:%v, %v", currentIndex, errors.Unwrap(err))
			}
			log.Current.Infow("preflight check", log.Fields{"step": currentIndex, "name": step.Name, "status": "ok"})
		}
		var execResult *ksqldb.KsqlResponseSlice
		if execResult, err = client.Execute(ksqldb.ExecOptions{KSql: step.Exec}); err != nil {
			log.Current.Error(err)
			os.Exit(-1)
		}
		log.Current.Infof("%+v", execResult)
		log.Current.Infow("processed", log.Fields{"status": "ok", "step": currentIndex, "name": step.Name})
	}

	client.Close()

	if err := rootCmd.MarkPersistentFlagRequired("file"); err != nil {
		log.Current.Fatal(err)
	}
}
