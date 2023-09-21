package main

import (
	"os"
)

func main() {
	if err := InitDB(os.Getenv("DATABASE_URL")); err != nil {
		panic(err)
	}

	// if err := PopulateBusLineInformation(context.Background()); err != nil {
	// 	log.Println(err)
	// }

	// if err := UpdateBusLocation(context.Background()); err != nil {
	// 	log.Println(err)
	// }

	StartAPIServer()
}
