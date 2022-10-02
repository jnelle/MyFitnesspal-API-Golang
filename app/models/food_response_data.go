package models

type FoodSearchResponse struct {
	Items []Items `json:"items"`
}

type Energy struct {
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
}

type NutritionalContents struct {
	Calcium            float64 `json:"calcium"`
	Carbohydrates      float64 `json:"carbohydrates"`
	Cholesterol        float64 `json:"cholesterol"`
	Energy             Energy  `json:"energy"`
	Fat                float64 `json:"fat"`
	Fiber              float64 `json:"fiber"`
	Iron               float64 `json:"iron"`
	MonounsaturatedFat float64 `json:"monounsaturated_fat"`
	PolyunsaturatedFat float64 `json:"polyunsaturated_fat"`
	Potassium          float64 `json:"potassium"`
	Protein            float64 `json:"protein"`
	SaturatedFat       float64 `json:"saturated_fat"`
	Sodium             float64 `json:"sodium"`
	Sugar              float64 `json:"sugar"`
	TransFat           float64 `json:"trans_fat"`
	VitaminA           float64 `json:"vitamin_a"`
	VitaminC           float64 `json:"vitamin_c"`
}

type ServingSizes struct {
	ID                  string  `json:"id"`
	Index               float64 `json:"index"`
	NutritionMultiplier float64 `json:"nutrition_multiplier"`
	Unit                string  `json:"unit"`
	Value               float64 `json:"value"`
}

type Item struct {
	BrandName           string              `json:"brand_name"`
	CountryCode         string              `json:"country_code"`
	Deleted             bool                `json:"deleted"`
	Description         string              `json:"description"`
	ID                  string              `json:"id"`
	NutritionalContents NutritionalContents `json:"nutritional_contents"`
	Public              bool                `json:"public"`
	ServingSizes        []ServingSizes      `json:"serving_sizes"`
	Type                string              `json:"type"`
	UserID              string              `json:"user_id"`
	Verified            bool                `json:"verified"`
	Version             string              `json:"version"`
}

type Items struct {
	Item Item          `json:"item"`
	Tags []interface{} `json:"tags"`
	Type string        `json:"type"`
}

type FoodSearchResponseWithoutAuth struct {
	Items []FoodSearchResponseWithoutAuthItems `json:"items"`
	Meta  Meta                                 `json:"meta"`
}

type FoodSearchResponseWithoutAuthItems struct {
	Item FoodSearchResponseWithoutAuthItem `json:"item"`
	Tags []interface{}                     `json:"tags"`
	Type string                            `json:"type"`
}

type FoodSearchResponseWithoutAuthItem struct {
	BrandName           string              `json:"brand_name"`
	CountryCode         string              `json:"country_code"`
	Deleted             bool                `json:"deleted"`
	Description         string              `json:"description"`
	ID                  int                 `json:"id"`
	NutritionalContents NutritionalContents `json:"nutritional_contents"`
	Public              bool                `json:"public"`
	ServingSizes        []ServingSizes      `json:"serving_sizes"`
	Type                string              `json:"type"`
	UserID              string              `json:"user_id"`
	Verified            bool                `json:"verified"`
	Version             string              `json:"version"`
}

type Meta struct {
	TotalEntries int `json:"total_entries"`
}
