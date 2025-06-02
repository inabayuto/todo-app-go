package main

import (
	"fmt"
	"todo-app/app/models"
)

func main() {
	/*
		fmt.Println(config.Config.Port)
		fmt.Println(config.Config.SQLDriver)
		fmt.Println(config.Config.DbHost)
		fmt.Println(config.Config.DbPort)
		fmt.Println(config.Config.DbUser)
		fmt.Println(config.Config.DbPassword)
		fmt.Println(config.Config.DbName)
		fmt.Println(config.Config.LogFile)

		log.Println("test")
	*/
	fmt.Println(models.Db)
}
