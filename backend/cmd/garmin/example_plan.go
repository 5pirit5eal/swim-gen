package main

import "github.com/5pirit5eal/swim-gen/internal/models"

// BuiltInExamplePlan provides a comprehensive training plan demonstrating
// all features: SubRows, Equipment, warmup/cooldown detection
var BuiltInExamplePlan = &models.Plan{
	Title:       "Advanced Freestyle Training",
	Description: "Comprehensive workout with compound sets and equipment drills",
	Table: models.Table{
		{
			Amount:    1,
			Distance:  400,
			Content:   "Einschwimmen",
			Intensity: "Warmup",
			Break:     "",
		},
		{
			Amount:    8,
			Distance:  50,
			Content:   "Kick Drill",
			Intensity: "Moderate",
			Break:     "15s",
			Equipment: []models.EquipmentType{models.EquipmentKickboard},
		},
		{
			Amount:    4,
			Distance:  1000,
			Content:   "Main Set",
			Intensity: "Hard",
			Break:     "20s",
			SubRows: []models.Row{
				{
					Distance:  800,
					Content:   "Freestyle",
					Intensity: "Hard",
				},
				{
					Distance:  200,
					Content:   "IM",
					Intensity: "Moderate",
				},
			},
		},
		{
			Amount:    6,
			Distance:  100,
			Content:   "Pull Set",
			Intensity: "Moderate",
			Break:     "20s",
			Equipment: []models.EquipmentType{models.EquipmentBuoy},
		},
		{
			Amount:    1,
			Distance:  400,
			Content:   "Ausschwimmen",
			Intensity: "Cooldown",
			Break:     "",
		},
	},
}
