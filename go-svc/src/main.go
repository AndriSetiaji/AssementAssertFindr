package main

import (
	"go-svc/config"
	"go-svc/src/routes"

	"gorm.io/gorm"
)

var (
	db *gorm.DB = config.ConnectDB()
)

func main() {
	defer config.DisconnectDB(db)

	//run all routes
	routes.Routes()
}
