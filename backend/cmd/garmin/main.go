package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/5pirit5eal/swim-gen/internal/garmin"
)

func main() {
	email := flag.String("email", os.Getenv("GARMIN_EMAIL"), "Garmin Connect email")
	password := flag.String("password", os.Getenv("GARMIN_PASSWORD"), "Garmin Connect password")
	tokenDir := flag.String("tokens", "./tokens", "Directory to store/load tokens")
	testUpload := flag.Bool("upload", false, "Whether to upload a test workout")
	flag.Parse()

	if *email == "" || *password == "" {
		log.Fatal("Garmin email and password must be provided via flags or GARMIN_EMAIL/GARMIN_PASSWORD env vars")
	}

	ctx := context.Background()

	// 1. Try to load existing tokens
	fmt.Printf("Attempting to load tokens from %s...\n", *tokenDir)
	client := garmin.NewClientFromTokens(nil, nil)
	err := client.LoadTokens(*tokenDir)

	if err != nil {
		fmt.Printf("Could not load tokens (%v), performing fresh login...\n", err)
		// 2. Perform fresh login
		client, err = garmin.NewClient(*email, *password)
		if err != nil {
			log.Fatalf("Login failed: %v", err)
		}
		fmt.Println("Login successful!")

		// 3. Save tokens for future use
		if err := client.SaveTokens(*tokenDir); err != nil {
			fmt.Printf("Warning: failed to save tokens: %v\n", err)
		} else {
			fmt.Printf("Tokens saved to %s\n", *tokenDir)
		}
	} else {
		fmt.Println("Tokens loaded successfully!")
	}

	// 4. Test GetWorkouts
	fmt.Println("\n--- Fetching Workouts (first 5) ---")
	workouts, err := client.GetWorkouts(ctx, 0, 5)
	if err != nil {
		log.Fatalf("Failed to fetch workouts: %v", err)
	}

	if len(workouts) == 0 {
		fmt.Println("No workouts found.")
	} else {
		for i, w := range workouts {
			fmt.Printf("%d: [%v] %s\n", i+1, w["workoutId"], w["workoutName"])
		}

		// 5. Test GetWorkoutByID (using the first one found)
		wID := int64(workouts[0]["workoutId"].(float64))
		fmt.Printf("\n--- Fetching Workout Details for ID %d ---\n", wID)
		details, err := client.GetWorkoutByID(ctx, wID)
		if err != nil {
			fmt.Printf("Failed to fetch workout details: %v\n", err)
		} else {
			fmt.Printf("Fetched details for: %s\n", details["workoutName"])
		}
	}

	// 6. Optional: Test Upload
	if *testUpload {
		fmt.Println("\n--- Uploading Test Swimming Workout ---")

		warmup := garmin.CreateWarmupStep(600, 1, nil)
		interval := garmin.CreateIntervalStep(1200, 2, nil)
		cooldown := garmin.CreateCooldownStep(300, 3, nil)

		segment := garmin.WorkoutSegment{
			SegmentOrder: 1,
			SportType: &garmin.SportTypeModel{
				SportTypeID:  garmin.SportTypeSwimming,
				SportTypeKey: "swimming",
				DisplayOrder: 3,
			},
			WorkoutSteps: []garmin.WorkoutStep{warmup, interval, cooldown},
		}

		swimmingWorkout := garmin.NewSwimmingWorkout("CLI Test Swim", 2100, []garmin.WorkoutSegment{segment})

		payload, err := json.Marshal(swimmingWorkout)
		if err != nil {
			log.Fatalf("Failed to marshal workout: %v", err)
		}

		result, err := client.UploadWorkout(ctx, payload)
		if err != nil {
			log.Fatalf("Failed to upload workout: %v", err)
		}

		fmt.Printf("Workout uploaded successfully! Result: %+v\n", result)
	}
}
