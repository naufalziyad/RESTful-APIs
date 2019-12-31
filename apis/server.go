package apis

import (
	"fmt"
	"log"
	"os"

	"github.com/naufalziyad/RESTful-APIs/apis/action"
	"github.com/naufalziyad/RESTful-APIs/apis/dummy"
	"github.com/naufalziyad/godotenv"
)

var server = action.Server{}

func Run() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Erorr getting env, %v", err)
	} else {
		fmt.Println("We have get env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	//add dummy data
	dummy.Load(server.DB)

	server.Run(":8080")
}
