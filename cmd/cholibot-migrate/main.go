package main

import "github.com/heshoots/cholibot/pkg/models"

func main() {
	models.Create()
	models.Migrate()
}
