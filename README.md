# (Unofficial) MyFitnessPal Golang API Wrapper

This is an unofficial MyFitnessPal API Wrapper which uses the "internal" MyFitnessPal APIs and not the official one.

## Installation

```shell
go get -u github.com/jnelle/MyFitnesspal-API-Golang
```

## How to use it?

### Example with credentials

```golang
import "github.com/jnelle/MyFitnesspal-API-Golang/app/client"

func main() {
	mfp := client.MFPClient{Username: "YOUR_EMAIL", Password: "YOUR_PASSWORD"}

	err := mfp.InitialLoad()
	if err != nil {
		log.Fatalln(err)
	}

	food, err := mfp.SearchFoodWithoutPagination("Pizza")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(food.Items[0])

}
```

### Example without credentials

```golang
import "github.com/jnelle/MyFitnesspal-API-Golang/app/client"

func main() {
	mfp := client.MFPClient{}

	// Example without using credentials and with pagination
	resultOffset, err := mfp.SearchFoodWithoutAuth("Pizza", 5, 1)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resultOffset)

	resultOffset, err = mfp.SearchFoodWithoutAuth("Pizza", 10, 2)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resultOffset)


}
```

## TODOS

- Add pagination for searchfood method with credentials
- For additional features just create an issue
