package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/5pirit5eal/swim-gen/internal/garmin"
	"github.com/5pirit5eal/swim-gen/internal/models"
)

func main() {
	plan := models.Plan{
		Table: models.Table{
			{
				Amount:     2,
				Multiplier: "x",
				Distance:   100,
				Break:      "30",
				Content:    "Kraul-Beine",
				Intensity:  "GA1",
				Sum:        200,
			},
			{
				Amount:     6,
				Multiplier: "x",
				Distance:   50,
				Break:      "20",
				Content:    "Unterwasser-Sculling",
				Intensity:  "TÜ",
				Sum:        100,
			},
			{
				Amount:     0,
				Multiplier: "",
				Distance:   0,
				Break:      "",
				Content:    "Gesamt",
				Intensity:  "",
				Sum:        300,
			},
		},
		Title:       "test plan",
		Description: "test description",
		PlanID:      "test123",
	}

	fit, err := garmin.ConvertTrainingPlanToFit(&plan, 25)
	if err != nil {
		log.Fatalf("Error converting plan to FIT: %v", err)
	}

	outFile, err := os.Create("test_workout.fit")
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outFile.Close()

	if err := garmin.EncodeFITFile(context.Background(), outFile, fit); err != nil {
		log.Fatalf("Error encoding FIT file: %v", err)
	}

	fmt.Println("FIT file written to test_workout.fit")
}
