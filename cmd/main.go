package main

import (
	"github.com/flyflyhe/dbMigrate/internal/db"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)

var yamlFile string

func main() {
	command := cobra.Command{
		Use:              "dbMigrate -c=config.yaml",
		TraverseChildren: true,
		Run: func(cmd *cobra.Command, args []string) {
			if yamlFileBytes, err := ioutil.ReadFile(yamlFile); err != nil {
				log.Println(err)
				return
			} else {
				if taskExternalConfigList, err := db.GetConfig(yamlFileBytes); err != nil {
					log.Println(err)
					return
				} else {
					log.Println(taskExternalConfigList)

					for _, v := range taskExternalConfigList {
						task := db.CreateTask(v.Dsn0, v.Dsn1, v.T0, v.T1, v.DsnType0, v.DsnType1)
						taskConfig := db.CreateTaskConfigByEConfig(v)

						task.SetFuncByConfig(taskConfig)
						log.Println(task.Migrate())
					}
				}
			}
		},
	}

	command.Flags().StringVarP(&yamlFile, "yaml", "y", "config.yaml", "yaml config file")
	_ = command.MarkFlagRequired("yaml")
	if err := command.Execute(); err != nil {
		log.Println(err)
	}
}
