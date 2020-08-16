package fitbit

import (
	"encoding/json"
	"strconv"
)

// FoodFavorites returns favorite food
// TODO: not tested, seems to be not implemented within the app
func (m *Session) FoodFavorites() (FoodCollectionList, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/foods/log/favorite.json")
	if err != nil {
		return FoodCollectionList{}, err
	}

	favs := FoodCollectionList{}
	if err := json.Unmarshal(contents, &favs); err != nil {
		return FoodCollectionList{}, err
	}

	return favs, nil
}

type FoodCollectionList []struct {
	AccessLevel        string `json:"accessLevel"`
	Amount             int    `json:"amount,omitempty"`
	Brand              string `json:"brand"`
	Calories           int    `json:"calories"`
	CreatorEncodedID   string `json:"creatorEncodedId,omitempty"`
	DateLastEaten      string `json:"dateLastEaten,omitempty"`
	DefaultServingSize int    `json:"defaultServingSize"`
	DefaultUnit        struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Plural string `json:"plural"`
	} `json:"defaultUnit"`
	FoodID     uint64 `json:"foodId"`
	MealTypeID uint64 `json:"mealTypeId,omitempty"`
	Name       string `json:"name"`
	Servings   []struct {
		Multiplier  int `json:"multiplier"`
		ServingSize int `json:"servingSize"`
		UnitID      int `json:"unitId"`
		Unit        struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Plural string `json:"plural"`
		} `json:"unit,omitempty"`
	} `json:"servings,omitempty"`
	Unit struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Plural string `json:"plural"`
	} `json:"unit"`
	Units             []int  `json:"units"`
	Locale            string `json:"locale,omitempty"`
	NutritionalValues struct {
		Biotin            int     `json:"biotin"`
		Calcium           int     `json:"calcium"`
		Calories          int     `json:"calories"`
		CaloriesFromFat   int     `json:"caloriesFromFat"`
		Cholesterol       int     `json:"cholesterol"`
		Copper            int     `json:"copper"`
		DietaryFiber      int     `json:"dietaryFiber"`
		FolicAcid         int     `json:"folicAcid"`
		Iodine            int     `json:"iodine"`
		Iron              int     `json:"iron"`
		Magnesium         int     `json:"magnesium"`
		Niacin            int     `json:"niacin"`
		PantothenicAcid   int     `json:"pantothenicAcid"`
		Phosphorus        int     `json:"phosphorus"`
		Potassium         int     `json:"potassium"`
		Protein           int     `json:"protein"`
		Riboflavin        int     `json:"riboflavin"`
		SaturatedFat      float64 `json:"saturatedFat"`
		Sodium            int     `json:"sodium"`
		Sugars            int     `json:"sugars"`
		Thiamin           int     `json:"thiamin"`
		TotalCarbohydrate int     `json:"totalCarbohydrate"`
		TotalFat          int     `json:"totalFat"`
		TransFat          int     `json:"transFat"`
		VitaminA          int     `json:"vitaminA"`
		VitaminB12        int     `json:"vitaminB12"`
		VitaminB6         int     `json:"vitaminB6"`
		VitaminC          int     `json:"vitaminC"`
		VitaminD          int     `json:"vitaminD"`
		VitaminE          int     `json:"vitaminE"`
		Zinc              int     `json:"zinc"`
	} `json:"nutritionalValues,omitempty"`
}

// FoodFrequent returns frequently eaten food
func (m *Session) FoodFrequent() (FoodCollectionList, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/foods/log/frequent.json")
	if err != nil {
		return FoodCollectionList{}, err
	}

	frequent := FoodCollectionList{}
	if err := json.Unmarshal(contents, &frequent); err != nil {
		return FoodCollectionList{}, err
	}

	return frequent, nil
}

// FoodRecent returns recently eaten food
func (m *Session) FoodRecent() (FoodCollectionList, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/foods/recent.json")
	if err != nil {
		return FoodCollectionList{}, err
	}

	favs := FoodCollectionList{}
	if err := json.Unmarshal(contents, &favs); err != nil {
		return FoodCollectionList{}, err
	}

	return favs, nil
}

// AddFoodFavorite adds food by id to favorites
func (m *Session) AddFoodFavorite(id uint64) error {
	_, err := m.makePOSTRequest("https://api.fitbit.com/1/user/-/foods/log/favorite/"+strconv.FormatUint(id, 10)+".json", map[string]string{})
	if err != nil {
		return err
	}

	return nil
}

// RemoveFoodFavorite removes food by id from facorites
func (m *Session) RemoveFoodFavorite(id uint64) error {
	_, err := m.makeDELETERequest("https://api.fitbit.com/1/user/-/foods/log/favorite/" + strconv.FormatUint(id, 10) + ".json")
	if err != nil {
		return err
	}

	return nil
}
