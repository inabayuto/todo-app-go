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
	/*
		// Create User
			fmt.Println(models.Db)

			u := &models.User{}
			u.Name = "test"
			u.Email = "test@example.com"
			u.PassWord = "4JpRPxcwrATa"

			u.CreateUser()
	*/
	/*
		// Get User
			u, _ := models.GetUser(1)
			fmt.Println(u)
	*/

	/*
		// Update user
		u, _ := models.GetUser(1)

		fmt.Println(u)

		u.Name = "test2"
		u.Email = "test2@example.com"
		u.UpdateUser()

		u, _ = models.GetUser(1)

		fmt.Println(u)
	*/

	// Delete User
	u, _ := models.GetUser(1)
	u.DeleteUser()

	u, _ = models.GetUser(1)
	fmt.Println(u)

}
