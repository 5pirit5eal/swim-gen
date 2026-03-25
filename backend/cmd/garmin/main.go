package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/5pirit5eal/swim-gen/internal/garmin"
	"github.com/5pirit5eal/swim-gen/internal/models"
)

func main() {
	// Authentication flags
	email := flag.String("email", os.Getenv("GARMIN_EMAIL"), "Garmin Connect email")
	password := flag.String("password", os.Getenv("GARMIN_PASSWORD"), "Garmin Connect password")
	tokenDir := flag.String("tokens", "./tokens", "Directory to store/load tokens")

	// Training plan input flags
	planJSON := flag.String("plan-json", "", "Path to JSON file with training plan")
	planInline := flag.String("plan-inline", "", "Inline JSON training plan definition")
	poolLength := flag.Int("pool-length", 25, "Pool length in meters (default: 25)")

	// Test mode flags (mutually exclusive)
	testConvert := flag.Bool("test-convert", false, "Convert plan to JSON and print summary")
	testUpload := flag.Bool("test-upload", false, "Convert plan and upload to Garmin Connect")
	testFIT := flag.Bool("test-fit", false, "Convert plan and export as FIT file")
	flag.Parse()

	// Validate exactly one test mode is selected
	modeCount := 0
	if *testConvert {
		modeCount++
	}
	if *testUpload {
		modeCount++
	}
	if *testFIT {
		modeCount++
	}

	if modeCount == 0 {
		log.Fatal("Must specify exactly one test mode: -test-convert, -test-upload, or -test-fit")
	}
	if modeCount > 1 {
		log.Fatal("Test modes are mutually exclusive. Choose only one: -test-convert, -test-upload, or -test-fit")
	}

	// Authentication required for upload mode
	var client *garmin.Client
	if *testUpload {
		if *email == "" || *password == "" {
			log.Fatal("Garmin email and password must be provided for upload mode via flags or GARMIN_EMAIL/GARMIN_PASSWORD env vars")
		}

		var err error
		client, err = authenticate(*email, *password, *tokenDir)
		if err != nil {
			log.Fatalf("Authentication failed: %v", err)
		}
	}

	// Load training plan
	plan := loadTrainingPlan(*planJSON, *planInline)

	// Execute selected test mode
	ctx := context.Background()
	switch {
	case *testConvert:
		runConvertMode(plan, *poolLength)
	case *testUpload:
		runUploadMode(ctx, client, plan, *poolLength, *tokenDir)
	case *testFIT:
		runFITMode(ctx, plan, *poolLength)
	}
}

// authenticate handles Garmin Connect authentication and token management
func authenticate(email, password, tokenDir string) (*garmin.Client, error) {
	fmt.Printf("Attempting to load tokens from %s...\n", tokenDir)
	client := garmin.NewClientFromTokens(nil, nil)
	err := client.LoadTokens(tokenDir)

	if err != nil {
		fmt.Printf("Could not load tokens (%v), performing fresh login...\n", err)
		client, err = garmin.NewClient(email, password)
		if err != nil {
			return nil, fmt.Errorf("login failed: %w", err)
		}
		fmt.Println("Login successful!")

		// Always attempt to save tokens
		if saveErr := client.SaveTokens(tokenDir); saveErr != nil {
			fmt.Printf("Warning: failed to save tokens: %v\n", saveErr)
		} else {
			fmt.Printf("Tokens saved to %s\n", tokenDir)
		}
	} else {
		fmt.Println("Tokens loaded successfully!")
	}

	return client, nil
}

// loadTrainingPlan loads a training plan from inline JSON, file, or uses built-in example
func loadTrainingPlan(jsonFile, jsonInline string) *models.Plan {
	var plan *models.Plan
	var err error

	if jsonInline != "" {
		plan, err = parseInlinePlan(jsonInline)
		if err != nil {
			log.Fatalf("Failed to parse inline plan: %v", err)
		}
	} else if jsonFile != "" {
		plan, err = loadPlanFromFile(jsonFile)
		if err != nil {
			log.Fatalf("Failed to load plan from file: %v", err)
		}
	} else {
		fmt.Println("Using built-in example training plan...")
		plan = BuiltInExamplePlan
	}

	// Calculate sums for all rows
	plan.Table.UpdateSum()

	return plan
}

// parseInlinePlan parses a training plan from inline JSON string
func parseInlinePlan(jsonStr string) (*models.Plan, error) {
	var plan models.Plan
	if err := json.Unmarshal([]byte(jsonStr), &plan); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return &plan, nil
}

// loadPlanFromFile loads a training plan from a JSON file
func loadPlanFromFile(filePath string) (*models.Plan, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var plan models.Plan
	if err := json.Unmarshal(data, &plan); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return &plan, nil
}

// runConvertMode converts the plan and prints summary + JSON
func runConvertMode(plan *models.Plan, poolLength int) {
	fmt.Println("Converting training plan to Garmin workout format...")
	workout := garmin.ConvertTrainingPlanToSwimWorkout(plan, poolLength)

	printWorkoutSummary(workout, plan, poolLength)
	printWorkoutJSON(workout)
}

// runUploadMode converts the plan, uploads to Garmin, and saves tokens
func runUploadMode(ctx context.Context, client *garmin.Client, plan *models.Plan, poolLength int, tokenDir string) {
	fmt.Println("Converting training plan to Garmin workout format...")
	workout := garmin.ConvertTrainingPlanToSwimWorkout(plan, poolLength)

	printWorkoutSummary(workout, plan, poolLength)

	fmt.Println("\n=== UPLOADING TO GARMIN CONNECT ===")
	payload, err := json.Marshal(workout)
	if err != nil {
		log.Fatalf("Failed to marshal workout: %v", err)
	}

	result, err := client.UploadWorkout(ctx, payload)
	if err != nil {
		log.Fatalf("Upload failed: %v", err)
	}

	fmt.Printf("✓ Upload successful!\n")
	fmt.Printf("Result: %+v\n", result)

	// Always attempt to save tokens after successful operation
	if err := client.SaveTokens(tokenDir); err != nil {
		fmt.Printf("Warning: failed to save tokens after upload: %v\n", err)
	}
}

// runFITMode converts the plan and exports as FIT file
func runFITMode(ctx context.Context, plan *models.Plan, poolLength int) {
	fmt.Println("Converting training plan to FIT format...")

	fit, err := garmin.ConvertTrainingPlanToFit(plan, poolLength)
	if err != nil {
		log.Fatalf("FIT conversion failed: %v", err)
	}

	outFile, err := os.Create("workout.fit")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	if err := garmin.EncodeFITFile(ctx, outFile, fit); err != nil {
		log.Fatalf("FIT encoding failed: %v", err)
	}

	fmt.Println("✓ FIT file written to workout.fit")
}

// printWorkoutSummary prints a human-readable summary of the workout
func printWorkoutSummary(workout *garmin.BaseWorkout, plan *models.Plan, poolLength int) {
	fmt.Println("\n=== SWIM WORKOUT SUMMARY ===")
	fmt.Printf("Name: %s\n", workout.WorkoutName)
	fmt.Printf("Sport: Swimming\n")
	fmt.Printf("Estimated Duration: %s\n", formatDuration(workout.EstimatedDurationInSecs))
	fmt.Printf("Pool Length: %dm\n", poolLength)

	// Print training plan table
	fmt.Println("\n=== TRAINING PLAN ===")
	fmt.Println(plan.Table.String())

	// Calculate and display total distance
	totalDistance := calculateTotalDistance(plan.Table)
	fmt.Printf("Total Distance: %dm\n", totalDistance)

	// Print workout structure
	fmt.Println("\n=== WORKOUT STRUCTURE ===")
	for i, segment := range workout.WorkoutSegments {
		fmt.Printf("Segment %d (%s):\n", i+1, segment.SportType.SportTypeKey)
		printSteps(segment.WorkoutSteps, 0)
	}
}

// printSteps recursively prints workout steps with indentation
func printSteps(steps []garmin.WorkoutStep, depth int) {
	indent := strings.Repeat("  ", depth)
	for i, step := range steps {
		isLast := i == len(steps)-1
		prefix := "├──"
		if isLast {
			prefix = "└──"
		}

		switch s := step.(type) {
		case garmin.ExecutableStep:
			stepName := s.StepType.StepTypeKey
			distance := ""
			if s.EndConditionValue != nil {
				distance = fmt.Sprintf("%.0fm", *s.EndConditionValue)
			}

			fmt.Printf("%s %s %s (%s)\n", indent, prefix, strings.Title(stepName), distance)

			// Show equipment if present
			if s.EquipmentType != nil {
				equipName := getEquipmentName(s.EquipmentType.EquipmentTypeID)
				fmt.Printf("%s  └─ Equipment: %s (ID: %d)\n", indent, equipName, s.EquipmentType.EquipmentTypeID)
			}

			// Show stroke type if present
			if s.StrokeType != nil {
				strokeName := getStrokeName(s.StrokeType.StrokeTypeID)
				fmt.Printf("%s  └─ Stroke: %s (ID: %d)\n", indent, strokeName, s.StrokeType.StrokeTypeID)
			}

		case garmin.RepeatGroup:
			fmt.Printf("%s %s REPEAT x%d\n", indent, prefix, s.NumberOfIterations)
			printSteps(s.WorkoutSteps, depth+1)
		}
	}
}

// printWorkoutJSON prints the workout as formatted JSON
func printWorkoutJSON(workout *garmin.BaseWorkout) {
	fmt.Println("\n=== GARMIN WORKOUT JSON ===")
	payload, err := json.MarshalIndent(workout, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal workout: %v", err)
	}
	fmt.Println(string(payload))
}

// formatDuration converts seconds to a human-readable string (HH:MM:SS)
func formatDuration(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, secs)
	}
	return fmt.Sprintf("%d:%02d", minutes, secs)
}

// calculateTotalDistance calculates the total swimming distance from a plan table
func calculateTotalDistance(table models.Table) int {
	total := 0
	for _, row := range table {
		if strings.Contains(row.Content, "Gesamt") || strings.Contains(row.Content, "Total") {
			continue
		}
		total += row.Sum
	}
	return total
}

// getEquipmentName returns a human-readable name for equipment type ID
func getEquipmentName(id int) string {
	names := map[int]string{
		1: "Fins",
		2: "Kickboard",
		3: "Paddles",
		4: "Pull Buoy",
		5: "Snorkel",
	}
	if name, ok := names[id]; ok {
		return name
	}
	return fmt.Sprintf("Unknown (ID: %d)", id)
}

// getStrokeName returns a human-readable name for stroke type ID
func getStrokeName(id int) string {
	names := map[int]string{
		1:  "Any Stroke",
		2:  "Backstroke",
		3:  "Breaststroke",
		4:  "Drill",
		5:  "Fly",
		6:  "Freestyle",
		7:  "Individual Medley",
		8:  "Mixed",
		9:  "IM by Round",
		10: "Reverse IM by Round",
	}
	if name, ok := names[id]; ok {
		return name
	}
	return fmt.Sprintf("Unknown (ID: %d)", id)
}
