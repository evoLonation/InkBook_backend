package main

import (
	"backend/api"
	"backend/entity"
)

func main() {
	api.Start("127.0.0.1:80")
	print(entity.Db.Error)
}
