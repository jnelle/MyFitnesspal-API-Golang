package utils

import (
	"net/url"

	"github.com/valyala/fasthttp"
)

func BuildFoodSearchURL(foodName string, searchFoodBaseURL string) string {
	uri := &fasthttp.URI{
		DisablePathNormalizing: true,
	}

	uri.SetScheme("https")
	uri.SetHost(searchFoodBaseURL)
	t := &url.URL{Path: foodName}

	uri.QueryArgs().Add("q", t.String())
	uri.QueryArgs().Add("scope", "all")
	uri.QueryArgs().Add("max_items", "1000")
	uri.QueryArgs().Add("resource_type[]", "foods")
	uri.QueryArgs().Add("fields[]", "id")
	uri.QueryArgs().Add("fields[]", "nutritional_contents")
	uri.QueryArgs().Add("fields[]", "serving_sizes")
	uri.QueryArgs().Add("fields[]", "version")
	uri.QueryArgs().Add("fields[]", "brand_name")
	uri.QueryArgs().Add("fields[]", "description")
	uri.QueryArgs().Add("resource_type[]", "venues")

	return uri.String()

}
