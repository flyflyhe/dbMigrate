package main

import (
	"github.com/flyflyhe/dbMigrate/internal/db"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"sync"
)

var yamlFile string
var debug bool

func main() {
	command := cobra.Command{
		Use:              "dbMigrate -c=config.yaml",
		TraverseChildren: true,
		Run: func(cmd *cobra.Command, args []string) {
			// Default level for this example is info, unless debug flag is present
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
			if debug {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}
			if yamlFileBytes, err := ioutil.ReadFile(yamlFile); err != nil {
				log.Println(err)
				return
			} else {
				if taskExternalConfigList, err := db.GetConfig(yamlFileBytes); err != nil {
					log.Println(err)
					return
				} else {
					log.Println(taskExternalConfigList)

					wg := sync.WaitGroup{}
					wg.Add(len(taskExternalConfigList))
					for k, v := range taskExternalConfigList {
						log.Println("任务", k, "--", v.T0, "===>", v.T1)
						task := db.CreateTask(v.Dsn0, v.Dsn1, v.T0, v.T1, v.DsnType0, v.DsnType1)
						taskConfig := db.CreateTaskConfigByEConfig(v)

						go func() {
							defer wg.Done()
							task.SetFuncByConfig(taskConfig)
							log.Println(task.Migrate())
						}()
					}
					wg.Wait()
				}
			}
		},
	}

	command.Flags().StringVarP(&yamlFile, "yaml", "y", "config.yaml", "yaml config file")
	command.Flags().BoolVarP(&debug, "debug", "d", false, "是否开启debug")
	_ = command.MarkFlagRequired("yaml")
	if err := command.Execute(); err != nil {
		log.Println(err)
	}
}
