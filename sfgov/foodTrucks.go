package sfgov

import (
	"encoding/json"
	"fmt"
	"foodtrucks/util"
	"io/ioutil"
	"net/http"
	"time"
)

const SF_TIMEZONE = "America/Los_Angeles"

type ByName []FoodTruck

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Applicant < a[j].Applicant }

type FoodTruck struct {
	DayOrder int `json:"dayorder,string"`
	DayOfWeekStr string `json:"dayofweekstr"`
	StartTime string `json:"starttime"`
	EndTime string `json:"endtime"`
	Start24 string `json:"start24"`
	End24 string `json:"end24"`
	Permit string `json:"permit"`
	Location string `json:"location"`
	LocationDescription string `json:"locationdesc"`
	OptionalText string `json:"optionaltext"`
	// Don't really need this to fulfill the requirements
	//LocationId string `json:"locationid"`
	//CNN string `json:"cnn"`
	//AddrDateCreated time.Time `json:"addr_date_create,string"`
	//AddrDateModified time.Time `json:"addr_date_modified,string"`
	//Block string `json:"block"`
	//Lot string `json:"lot"`
	ColdTruck util.YesNoBool `json:"coldtruck,string"`
	Applicant string         `json:"applicant"`  // Name of the food truck
	X float64                `json:"x,string"`
	Y float64                `json:"y,string"`
	Latitude float64         `json:"latitude,string"`
	Longitude float64        `json:"longitude,string"`
	// Location 2 in json looks like a copy of Latitude and Longitude with a broken human address
	// not going to count it for now

	// the computed regions are undocumented, not going to do anything with them for now
}

func (r *SfGoveApi) GetOpenNow() ([]FoodTruck, error) {
	return r.GetOpenAtTime(time.Now())
}

func (r *SfGoveApi) GetOpenAtTime(t time.Time) ([]FoodTruck, error) {
	foodTrucks, err := r.GetMobileFoodTrucksSchedule()
	if err != nil {
		return foodTrucks, err
	}

	// We need to use the timezone of SF
	location, err := time.LoadLocation(SF_TIMEZONE)
	if err != nil {
		// this might only happen if you have a typo in the timezone name
		panic(err)
	}
	timeInSF := t.In(location)

	currentTimeString := timeInSF.Format("15:04")
	currentDay := int(timeInSF.Weekday())
	var openFoodTrucks []FoodTruck
	for _, ft := range foodTrucks {
		if ft.DayOrder == currentDay && ft.Start24 <= currentTimeString && currentTimeString < ft.End24 {
			openFoodTrucks = append(openFoodTrucks, ft)
		}
	}

	return openFoodTrucks, nil
}


func (r *SfGoveApi) GetMobileFoodTrucksSchedule() ([]FoodTruck, error) {
	req, err := http.NewRequest("GET", r.MobileFoodScheduleEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	// don't forget to close a read so that we don't have a memory leak
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		// something is wrong
		return nil, fmt.Errorf("endpoint returned unexpected http status code: %d", resp.StatusCode)
	}
	var foodTrucks []FoodTruck
	err = json.Unmarshal(data, &foodTrucks)
	if err != nil {
		// bad json or field type mismatch
		return nil, err
	}

	return foodTrucks, nil
}