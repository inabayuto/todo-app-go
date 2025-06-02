package main

import (
	"log"
	"todo-app/app/models/controllers"
)

func main() {

	err := controllers.StartMainServer()
	if err != nil {
		log.Println(err)
	}

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
		u.Name = "test2"
		u.Email = "test2@example.com"
		u.PassWord = "C8zuFu44LcQH"

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

	/*
		// Delete User
		u, _ := models.GetUser(1)
		u.DeleteUser()

		u, _ = models.GetUser(1)
		fmt.Println(u)
	*/

	/*
		// Create Todo
		u, _ := models.GetUser(3)
		u.CreateTodo("Third Todo")
	*/

	/*
		// Get Todo
		t, _ := models.GetTodo(1)
		fmt.Println(t)
	*/

	/*
		// Get Todos
		todos, _ := models.GetTodos()
		for k, v := range todos {
			fmt.Println(k, v)
		}
	*/

	/*
		// Get Todos ByUser
		user2, _ := models.GetUser(3)
		todos, _ := user2.GetTodosByUser()
		for k, v := range todos {
			fmt.Println(k, v)
		}
	*/

	/*
		// Update Todo
		t, _ := models.GetTodo(1)
		t.Content = "Update Todo"
		t.UpdateTodo()
	*/

	/*
		// Delete Todo
		t, _ := models.GetTodo(3)
		t.DeleteTodo()
	*/

}
