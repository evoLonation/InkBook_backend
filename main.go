package main

import (
	"backend/api"
	"backend/entity"
)

func main() {
	api.Start("127.0.0.1:8080")
	print(entity.Db.Error)
}
