package main

import (
	"fmt"

	"./models"
	"./utils"
)

func main() {
	users := []models.Account{
		models.Account{Name: "test1"},
		models.Account{Name: "test2"},
	}

	ms, _ := utils.ToMapList(users, []string{"Name"}, utils.FilterModeInclude)

	fmt.Println(ms)
}
