package fitbit

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

// FoodGoal contains the food goal of a user
type FoodGoal struct {
	FoodPlan struct {
		EstimatedDate string `json:"estimatedDate"`
		Intensity     string `json:"intensity"`
		Personalized  bool   `json:"personalized"`
	} `json:"foodPlan"`
	Goals struct {
		Calories int `json:"calories"`
	} `json:"goals"`
}

// FoodGoal requests the food goal of the user
func (m *Session) FoodGoal() (FoodGoal, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/foods/log/goal.json")
	if err != nil {
		return FoodGoal{}, err
	}

	foods := FoodGoal{}
	if err := json.Unmarshal(contents, &foods); err != nil {
		return FoodGoal{}, err
	}

	return foods, nil
}

// SetFoodGoal sets the users food goal
func (m *Session) SetFoodGoal(goals map[string]string) (FoodGoal, error) {
	contents, err := m.makePOSTRequest("https://api.fitbit.com/1/user/-/foods/log/goal.json", goals)
	if err != nil {
		return FoodGoal{}, err
	}

	foodgoal := FoodGoal{}
	if err := json.Unmarshal(contents, &foodgoal); err != nil {
		return FoodGoal{}, err
	}

	return foodgoal, nil
}

// FoodLog contains a users food log
type FoodLog struct {
	Foods []struct {
		IsFavorite bool   `json:"isFavorite"`
		LogDate    string `json:"logDate"`
		LogID      uint64 `json:"logId"`
		LoggedFood struct {
			AccessLevel      string  `json:"accessLevel"`
			Amount           float64 `json:"amount"`
			Brand            string  `json:"brand"`
			Calories         uint64  `json:"calories"`
			CreatorEncodedID string  `json:"creatorEncodedId,omitempty"`
			FoodID           uint64  `json:"foodId"`
			Locale           string  `json:"locale"`
			MealTypeID       uint64  `json:"mealTypeId"`
			Name             string  `json:"name"`
			Unit             struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Plural string `json:"plural"`
			} `json:"unit"`
			Units []int `json:"units"`
		} `json:"loggedFood,omitempty"`
		NutritionalValues struct {
			Calories int     `json:"calories"`
			Carbs    float64 `json:"carbs"`
			Fat      float64 `json:"fat"`
			Fiber    float64 `json:"fiber"`
			Protein  float64 `json:"protein"`
			Sodium   int     `json:"sodium"`
		} `json:"nutritionalValues"`
	} `json:"foods"`
	Goals struct {
		Calories             int `json:"calories"`
		EstimatedCaloriesOut int `json:"estimatedCaloriesOut"`
	} `json:"goals"`
	Summary struct {
		Calories int     `json:"calories"`
		Carbs    float64 `json:"carbs"`
		Fat      float64 `json:"fat"`
		Fiber    float64 `json:"fiber"`
		Protein  float64 `json:"protein"`
		Sodium   float64 `json:"sodium"`
		Water    int     `json:"water"`
	} `json:"summary"`
}

// FoodLogByDay returns the food log by a given date
// date must be in the format yyyy-MM-dd
func (m *Session) FoodLogByDay(day string) (FoodLog, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/foods/log/date/" + day + ".json")
	if err != nil {
		return FoodLog{}, err
	}

	foods := FoodLog{}
	if err := json.Unmarshal(contents, &foods); err != nil {
		return FoodLog{}, err
	}

	return foods, nil
}

// FoodWaterLogDateRange contains a summary of calories or water for a given date range
type FoodWaterLogDateRange struct {
	FoodsLogCaloriesIn []struct {
		DateTime string `json:"dateTime"`
		Value    string `json:"value"`
	} `json:"foods-log-caloriesIn,omitempty"`
	FoodsLogWater []struct {
		DateTime string `json:"dateTime"`
		Value    string `json:"value"`
	} `json:"foods-log-water,omitempty"`
}

// FoodLogByDateRange returns the calories log of a given time range by date
// date must be in the format yyyy-MM-dd
func (m *Session) FoodLogByDateRange(startDay string, endDay string) (FoodWaterLogDateRange, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/foods/log/caloriesIn/date/" + startDay + "/" + endDay + ".json")
	if err != nil {
		return FoodWaterLogDateRange{}, err
	}

	foods := FoodWaterLogDateRange{}
	if err := json.Unmarshal(contents, &foods); err != nil {
		return FoodWaterLogDateRange{}, err
	}

	return foods, nil
}

// WaterLogByDateRange returns the calories log of a given time range by date
// date must be in the format yyyy-MM-dd
func (m *Session) WaterLogByDateRange(startDay string, endDay string) (FoodWaterLogDateRange, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/foods/log/water/date/" + startDay + "/" + endDay + ".json")
	if err != nil {
		return FoodWaterLogDateRange{}, err
	}

	foods := FoodWaterLogDateRange{}
	if err := json.Unmarshal(contents, &foods); err != nil {
		return FoodWaterLogDateRange{}, err
	}

	return foods, nil
}

// WaterLog contains the water log of a user
type WaterLog struct {
	Summary struct {
		Water int `json:"water"`
	} `json:"summary"`
	Water []struct {
		Amount int    `json:"amount"`
		LogID  uint64 `json:"logId"`
	} `json:"water"`
}

// WaterLogByDay returns the water log by a given date
// date must be in the format yyyy-MM-dd
func (m *Session) WaterLogByDay(day string) (WaterLog, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/foods/log/water/date/" + day + ".json")
	if err != nil {
		return WaterLog{}, err
	}

	water := WaterLog{}
	if err := json.Unmarshal(contents, &water); err != nil {
		return WaterLog{}, err
	}

	return water, nil
}

// WaterGoal describes the water goal of the user
type WaterGoal struct {
	Goal struct {
		Goal      float64 `json:"goal"`
		StartDate string  `json:"startDate"`
	} `json:"goal"`
}

// WaterGoal get's the users water goal
func (m *Session) WaterGoal() (WaterGoal, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/user/-/foods/log/water/goal.json")
	if err != nil {
		return WaterGoal{}, err
	}

	water := WaterGoal{}
	if err := json.Unmarshal(contents, &water); err != nil {
		return WaterGoal{}, err
	}

	return water, nil
}

// SetWaterGoal sets a new water goal
func (m *Session) SetWaterGoal(goal float64) (WaterGoal, error) {
	contents, err := m.makePOSTRequest("https://api.fitbit.com/1/user/-/foods/log/water/goal.json", map[string]string{"target": strconv.FormatFloat(goal, 'E', -1, 64)})
	if err != nil {
		return WaterGoal{}, err
	}

	water := WaterGoal{}
	if err := json.Unmarshal(contents, &water); err != nil {
		return WaterGoal{}, err
	}

	return water, nil
}

// AddWater adds a new water log entry for the user
// date has to be in the format yyyy-MM-dd
// amount contains the amount of water consumed
// unit can be ml, fl oz or cup
func (m *Session) AddWater(date string, amount float64, unit string) (WaterLog, error) {
	if date == "" {
		return WaterLog{}, errors.New("date must be defined")
	}
	if amount <= 0 {
		return WaterLog{}, errors.New("amount must me greater than 0")
	}
	if unit != "ml" && unit != "fl oz" && unit != "cup" {
		return WaterLog{}, errors.New("unit must be ml, fl oz or cup")
	}

	contents, err := m.makePOSTRequest("https://api.fitbit.com/1/user/-/foods/log/water.json", map[string]string{
		"amount": strconv.FormatFloat(amount, 'f', 1, 10),
		"date":   date,
		"unit":   unit,
	})
	if err != nil {
		return WaterLog{}, err
	}

	water := WaterLog{}
	if err := json.Unmarshal(contents, &water); err != nil {
		return WaterLog{}, err
	}

	return water, nil
}

// UpdateWater updates an existing water log entry with new values
// amount contains the amount of water consumed
// unit can be ml, fl oz or cup
func (m *Session) UpdateWater(id uint64, amount float64, unit string) (WaterLog, error) {
	if id == 0 {
		return WaterLog{}, errors.New("id must be defined")
	}
	if amount <= 0 {
		return WaterLog{}, errors.New("amount must me greater than 0")
	}
	if unit != "ml" && unit != "fl oz" && unit != "cup" {
		return WaterLog{}, errors.New("unit must be ml, fl oz or cup")
	}

	contents, err := m.makePOSTRequest("https://api.fitbit.com/1/user/-/foods/log/water/"+strconv.FormatUint(id, 10)+".json", map[string]string{
		"amount": strconv.FormatFloat(amount, 'f', 1, 10),
		"unit":   unit,
	})
	if err != nil {
		return WaterLog{}, err
	}

	water := WaterLog{}
	if err := json.Unmarshal(contents, &water); err != nil {
		return WaterLog{}, err
	}

	return water, nil
}

// RemoveWater removes an existing water log entry
func (m *Session) RemoveWater(id uint64) error {
	if id == 0 {
		return errors.New("id must be defined")
	}

	_, err := m.makeDELETERequest("https://api.fitbit.com/1/user/-/foods/log/water/" + strconv.FormatUint(id, 10) + ".json")
	if err != nil {
		return err
	}

	return nil
}

// FoodLocales contains a list of supported food locales
type FoodLocales []struct {
	Barcode     bool   `json:"barcode"`
	ImageUpload bool   `json:"imageUpload"`
	Label       string `json:"label"`
	Value       string `json:"value"`
}

// FoodLocales returns a list of supported food locales
func (m *Session) FoodLocales() (FoodLocales, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/foods/locales.json")
	if err != nil {
		return FoodLocales{}, err
	}

	foods := FoodLocales{}
	if err := json.Unmarshal(contents, &foods); err != nil {
		return FoodLocales{}, err
	}

	return foods, nil
}

// FoodSearchResult contains a list of food found by FoodSearch
type FoodSearchResult struct {
	Foods []struct {
		AccessLevel        string `json:"accessLevel"`
		Brand              string `json:"brand"`
		Calories           int    `json:"calories"`
		DefaultServingSize int    `json:"defaultServingSize"`
		DefaultUnit        struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Plural string `json:"plural"`
		} `json:"defaultUnit"`
		FoodID    int    `json:"foodId"`
		IsGeneric bool   `json:"isGeneric"`
		Locale    string `json:"locale"`
		Name      string `json:"name"`
		Units     []int  `json:"units"`
	} `json:"foods"`
}

// FoodSearch searches for food within the fitbit database for matching food
func (m *Session) FoodSearch(value string) (FoodSearchResult, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/foods/search.json?query=" + url.QueryEscape(value))
	if err != nil {
		return FoodSearchResult{}, err
	}

	foods := FoodSearchResult{}
	if err := json.Unmarshal(contents, &foods); err != nil {
		return FoodSearchResult{}, err
	}

	return foods, nil
}

func (m *Session) FoodByID(id uint64) (FoodEntry, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/foods/" + strconv.FormatUint(id, 10) + ".json")
	if err != nil {
		return FoodEntry{}, err
	}

	foods := FoodEntry{}
	if err := json.Unmarshal(contents, &foods); err != nil {
		return FoodEntry{}, err
	}

	return foods, nil
}

type FoodEntry struct {
	Food struct {
		AccessLevel        string `json:"accessLevel"`
		Brand              string `json:"brand"`
		Calories           int    `json:"calories"`
		DefaultServingSize int    `json:"defaultServingSize"`
		DefaultUnit        struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Plural string `json:"plural"`
		} `json:"defaultUnit"`
		FoodID    int    `json:"foodId"`
		IsGeneric bool   `json:"isGeneric"`
		Locale    string `json:"locale"`
		Name      string `json:"name"`
		Servings  []struct {
			Multiplier  int `json:"multiplier"`
			ServingSize int `json:"servingSize"`
			Unit        struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Plural string `json:"plural"`
			} `json:"unit"`
			UnitID int `json:"unitId"`
		} `json:"servings"`
		Units []int `json:"units"`
	} `json:"food"`
}

func (m *Session) FoodUnits() (FoodUnits, error) {
	contents, err := m.makeRequest("https://api.fitbit.com/1/foods/units.json")
	if err != nil {
		return FoodUnits{}, err
	}

	foods := FoodUnits{}
	if err := json.Unmarshal(contents, &foods); err != nil {
		return FoodUnits{}, err
	}

	return foods, nil
}

type FoodUnits []struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Plural string `json:"plural"`
}

// AddFood logs a new consumed food entry
func (m *Session) AddFood(data NewFoodLog) (AddFoodLogResponse, error) {
	dataToPost := make(map[string]string)
	if data.MealTypeID < 1 || data.MealTypeID > 7 {
		return AddFoodLogResponse{}, errors.New("mealTypeID must be given and between 1 and 7")
	}
	if data.UnitID != 0 {
		return AddFoodLogResponse{}, errors.New("unitid must be given")
	}
	if data.Amount <= 0 {
		return AddFoodLogResponse{}, errors.New("amount must be given")
	}
	if data.Date == "" {
		return AddFoodLogResponse{}, errors.New("date must be given")
	}

	dataToPost["mealTypeId"] = strconv.Itoa(data.MealTypeID)
	dataToPost["unitId"] = strconv.FormatUint(data.UnitID, 10)
	dataToPost["amount"] = strconv.FormatFloat(data.Amount, 'f', 2, 10)
	dataToPost["date"] = data.Date

	if data.FoodID != 0 {
		dataToPost["foodId"] = strconv.FormatUint(data.FoodID, 10)
	} else if data.FoodName != "" {
		dataToPost["foodName"] = data.FoodName
		if data.Favorite {
			dataToPost["favorite"] = "true"
		}
		if data.BrandName != "" {
			dataToPost["brandName"] = data.BrandName
		}
		if data.Calories > 0 {
			dataToPost["calories"] = strconv.FormatUint(data.Calories, 10)
		}
	} else {
		return AddFoodLogResponse{}, errors.New("either foodId or foodName must be given")
	}

	contents, err := m.makePOSTRequest("https://api.fitbit.com/1/user/-/foods/log.json", dataToPost)
	if err != nil {
		return AddFoodLogResponse{}, err
	}

	foods := AddFoodLogResponse{}
	if err := json.Unmarshal(contents, &foods); err != nil {
		return AddFoodLogResponse{}, err
	}

	return foods, nil
}

// UpdateFood changes a stored food log entry
func (m *Session) UpdateFood(id uint64, data NewFoodLog) (AddFoodLogResponse, error) {
	dataToPost := make(map[string]string)
	if data.MealTypeID < 1 || data.MealTypeID > 7 {
		return AddFoodLogResponse{}, errors.New("mealTypeID must be given and between 1 and 7")
	}

	dataToPost["mealTypeId"] = strconv.Itoa(data.MealTypeID)
	if data.UnitID != 0 {
		dataToPost["unitId"] = strconv.FormatUint(data.UnitID, 10)
	}
	if data.Amount != 0 {
		dataToPost["amount"] = strconv.FormatFloat(data.Amount, 'f', 2, 10)
	}
	if data.Calories > 0 {
		dataToPost["calories"] = strconv.FormatUint(data.Calories, 10)
	}

	contents, err := m.makePOSTRequest("https://api.fitbit.com/1/user/-/foods/log/"+strconv.FormatUint(id, 10)+".json", dataToPost)
	if err != nil {
		return AddFoodLogResponse{}, err
	}

	foods := AddFoodLogResponse{}
	if err := json.Unmarshal(contents, &foods); err != nil {
		return AddFoodLogResponse{}, err
	}

	return foods, nil
}

// RemoveFood removes an existing food log entry
func (m *Session) RemoveFood(id uint64) error {
	if id == 0 {
		return errors.New("id must be defined")
	}

	_, err := m.makeDELETERequest("https://api.fitbit.com/1/user/-/foods/log/" + strconv.FormatUint(id, 10) + ".json")
	if err != nil {
		return err
	}

	return nil
}

type NewFoodLog struct {
	FoodID     uint64  `json:"foodId,omitempty"`
	FoodName   string  `json:"foodName,omitempty"`
	MealTypeID int     `json:"mealTypeId"`
	UnitID     uint64  `json:"unitId"` // given by unit search and food serch which units are captable here
	Amount     float64 `json:"amount"`
	Date       string  `json:"date"`
	Favorite   bool    `json:"favorite,omitempty"`
	BrandName  string  `json:"brandName,omitempty"`
	Calories   uint64  `json:"calories,omitempty"`
	// TODO: Additional Nutrition Information
}

type AddFoodLogResponse struct {
	FoodDay struct {
		Date    string `json:"date"`
		Summary struct {
			Calories int     `json:"calories"`
			Carbs    float64 `json:"carbs"`
			Fat      float64 `json:"fat"`
			Fiber    float64 `json:"fiber"`
			Protein  float64 `json:"protein"`
			Sodium   float64 `json:"sodium"`
			Water    int     `json:"water"`
		} `json:"summary"`
	} `json:"foodDay"`
	FoodLog struct {
		IsFavorite bool   `json:"isFavorite"`
		LogDate    string `json:"logDate"`
		LogID      int64  `json:"logId"`
		LoggedFood struct {
			AccessLevel string `json:"accessLevel"`
			Amount      int    `json:"amount"`
			Brand       string `json:"brand"`
			Calories    int    `json:"calories"`
			FoodID      int    `json:"foodId"`
			Locale      string `json:"locale"`
			MealTypeID  int    `json:"mealTypeId"`
			Name        string `json:"name"`
			Unit        struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Plural string `json:"plural"`
			} `json:"unit"`
			Units []int `json:"units"`
		} `json:"loggedFood"`
		NutritionalValues struct {
			Calories int `json:"calories"`
			Carbs    int `json:"carbs"`
			Fat      int `json:"fat"`
			Fiber    int `json:"fiber"`
			Protein  int `json:"protein"`
			Sodium   int `json:"sodium"`
		} `json:"nutritionalValues"`
	} `json:"foodLog"`
}
