package client_test

import (
	"testing"

	"github.com/jnelle/MyFitnesspal-API-Golang/app/client"
	"github.com/stretchr/testify/require"
)

var testClient = &client.MFPClient{Username: "saxigen930@delowd.com", Password: "fitnessproject"}

func TestInitialLoad(t *testing.T) {
	err := testClient.InitialLoad()
	if err != nil {
		t.Fatal(err)
	}
	require.NoError(t, err)
	require.NotEmpty(t, testClient.AuthSignKey)
	require.NotNil(t, testClient.AccessToken)
	require.NotNil(t, testClient.IDTokenResponse)
}

// func TestSearchFood(t *testing.T) {
// 	food, err := testClient.SearchFoodWithoutPagination("Pizza")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	require.NotNil(t, food)
// 	require.NoError(t, err)
// }

// func TestRefreshToken(t *testing.T) {
// 	err := testClient.RefreshToken()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	require.NoError(t, err)
// 	require.NotNil(t, testClient.AccessToken)
// }

// func TestFoodSearchWithoutAuth(t *testing.T) {
// 	result, err := testClient.SearchFoodWithoutAuth("Pizza", 10, 1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	require.NoError(t, err)
// 	require.NotNil(t, result)
// }
