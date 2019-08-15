package main

import (
	"fmt"
	"foodtrucks/sfgov"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"sort"
	"text/tabwriter"
)

const (
	MOBILE_FOOD_SCHEDULE = "https://data.sfgov.org/resource/jjew-r69b.json"
	PAGE_SIZE = 10
)

var (
	app = kingpin.New("SFFoodTrucks", "A command-line app to find food trucks")
	openNow = app.Command("openNow", "list ones open now")
	verbose = openNow.Flag("verbose", "Verbose mode.").Short('v').Bool()
)

func main() {
	SfGov := &sfgov.SfGovApi{
		MobileFoodScheduleEndpoint: MOBILE_FOOD_SCHEDULE,
	}

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case openNow.FullCommand():
		foodTrucks, err := SfGov.GetOpenNow()
		if err != nil {
			panic(err)
		}

		sort.Sort(sfgov.ByName(foodTrucks))

		header := "Name\tAddress"
		if *verbose {
			header = "Number\tName\tAddress\tDay\tStart\tEnd"
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug)
		fmt.Fprintln(w,header)

		end := len(foodTrucks)
		for i, ft := range foodTrucks {
			if  i % PAGE_SIZE == 0 && i > 0 && i < end {
				err := w.Flush()
				if err != nil {
					panic(err)
				}
				fmt.Println("\nPress the Enter for next page of results...")
				var input string
				fmt.Scanln(&input)
				fmt.Fprintln(w,header)
			}
			if *verbose {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n", i+1, ft.Applicant, ft.Location, ft.DayOfWeekStr, ft.StartTime, ft.EndTime)
			} else {
				fmt.Fprintf(w, "%s\t%s\n", ft.Applicant, ft.Location)
			}
		}
		// write anything else left in the buffer
		err = w.Flush()
		if err != nil {
			panic(err)
		}
		fmt.Println("No More Results")
	}

}
