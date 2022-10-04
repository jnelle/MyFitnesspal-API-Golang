package main

import (
	"log"

	"github.com/jnelle/MyFitnesspal-API-Golang/app/client"
)

func main() {

	// temp mail
	mfp := &client.MFPClient{Username: "laxewo9737@cebaike.com", Password: "fitnessproject"}
	err := mfp.InitialLoad()
	if err != nil {
		log.Fatalln(err)
	}

	// result, err := mfp.SearchFoodWithoutPagination("Pizza")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// Example without using credentials and pagination
	// resultOffset, err := mfp.SearchFoodWithoutAuth("Pizza", 5, 1)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println(resultOffset)

	// resultOffset, _ = mfp.SearchFoodWithoutAuth("Pizza", 10, 2)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println(result)

}
