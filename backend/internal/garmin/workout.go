package garmin

import (
	"strings"

	"github.com/5pirit5eal/swim-gen/internal/models"
)

// Garmin Connect Workout structs ported from python-garminconnect:
// https://github.com/cyberjunky/python-garminconnect/blob/master/garminconnect/workout.py

type SportTypeID int

const (
	SportTypeRunning          SportTypeID = 1
	SportTypeCycling          SportTypeID = 2
	SportTypeSwimming         SportTypeID = 4
	SportTypeStrengthTraining SportTypeID = 5
	SportTypeCardioTraining   SportTypeID = 6
	SportTypeYoga             SportTypeID = 7
	SportTypeHIIT             SportTypeID = 9
	SportTypeMultiSport       SportTypeID = 10
	SportTypeMobility         SportTypeID = 11
	SportTypeOther            SportTypeID = 3
)

type StepTypeID int

const (
	StepTypeWarmup   StepTypeID = 1
	StepTypeCooldown StepTypeID = 2
	StepTypeInterval StepTypeID = 3
	StepTypeRecovery StepTypeID = 4
	StepTypeRest     StepTypeID = 5
	StepTypeRepeat   StepTypeID = 6
	StepTypeOther    StepTypeID = 7
	StepTypeMain     StepTypeID = 8
)

type ConditionTypeID int

const (
	ConditionTypeDistance      ConditionTypeID = 3
	ConditionTypeTime          ConditionTypeID = 2
	ConditionTypeHeartRate     ConditionTypeID = 6
	ConditionTypeCalories      ConditionTypeID = 4
	ConditionTypePower         ConditionTypeID = 5
	ConditionTypeIterations    ConditionTypeID = 7
	ConditionTypeFixedRest     ConditionTypeID = 8
	ConditionTypeReps          ConditionTypeID = 10
	ConditionTypeSwimCSSOffset ConditionTypeID = 16
)

type TargetTypeID int

const (
	TargetTypeNoTarget        TargetTypeID = 1
	TargetTypePower           TargetTypeID = 2
	TargetTypeCadence         TargetTypeID = 3
	TargetTypeHeartRate       TargetTypeID = 4
	TargetTypeSpeed           TargetTypeID = 5
	TargetTypePaceZone        TargetTypeID = 6
	TargetTypeGrade           TargetTypeID = 7
	TargetTypeSwimStroke      TargetTypeID = 14
	TargetTypeSwimCSS         TargetTypeID = 17
	TargetTypeSwimInstruction TargetTypeID = 18
)

type SportTypeModel struct {
	SportTypeID  SportTypeID `json:"sportTypeId"`
	SportTypeKey string      `json:"sportTypeKey"`
	DisplayOrder int         `json:"displayOrder"`
}

type EndConditionModel struct {
	ConditionTypeID  ConditionTypeID `json:"conditionTypeId"`
	ConditionTypeKey string          `json:"conditionTypeKey"`
	DisplayOrder     int             `json:"displayOrder"`
	Displayable      bool            `json:"displayable"`
}

type TargetTypeModel struct {
	WorkoutTargetTypeID  TargetTypeID `json:"workoutTargetTypeId"`
	WorkoutTargetTypeKey string       `json:"workoutTargetTypeKey"`
	DisplayOrder         int          `json:"displayOrder"`
}

type StrokeTypeModel struct {
	StrokeTypeID int `json:"strokeTypeId"`
	DisplayOrder int `json:"displayOrder"`
}

type StrokeTypeID int

const (
	StrokeTypeAny          StrokeTypeID = 1
	StrokeTypeBackstroke   StrokeTypeID = 2
	StrokeTypeBreaststroke StrokeTypeID = 3
	StrokeTypeDrill        StrokeTypeID = 4
	StrokeTypeFly          StrokeTypeID = 5
	StrokeTypeFree         StrokeTypeID = 6
	StrokeTypeIM           StrokeTypeID = 7
	StrokeTypeMixed        StrokeTypeID = 8
	StrokeTypeIMByRound    StrokeTypeID = 9
	StrokeTypeReverseIM    StrokeTypeID = 10
)

type EquipmentTypeModel struct {
	EquipmentTypeID int `json:"equipmentTypeId"`
	DisplayOrder    int `json:"displayOrder"`
}

type EquipmentTypeID int

const (
	EquipmentFins      EquipmentTypeID = 1
	EquipmentKickboard EquipmentTypeID = 2
	EquipmentPaddles   EquipmentTypeID = 3
	EquipmentBuoy      EquipmentTypeID = 4
	EquipmentSnorkel   EquipmentTypeID = 5
)

type StepTypeModel struct {
	StepTypeID   StepTypeID `json:"stepTypeId"`
	StepTypeKey  string     `json:"stepTypeKey"`
	DisplayOrder int        `json:"displayOrder"`
}

// WorkoutStep represents either an ExecutableStep or RepeatGroup in the JSON.
// We use interface{} because a segment or repeat group can contain either.
type WorkoutStep interface {
	IsWorkoutStep()
}

type ExecutableStep struct {
	Type              string              `json:"type"` // "ExecutableStepDTO"
	StepOrder         int                 `json:"stepOrder"`
	StepType          *StepTypeModel      `json:"stepType,omitempty"`
	EndCondition      *EndConditionModel  `json:"endCondition,omitempty"`
	EndConditionValue *float64            `json:"endConditionValue,omitempty"`
	TargetType        *TargetTypeModel    `json:"targetType,omitempty"`
	StrokeType        *StrokeTypeModel    `json:"strokeType,omitempty"`
	EquipmentType     *EquipmentTypeModel `json:"equipmentType,omitempty"`
	ChildStepID       *int                `json:"childStepId,omitempty"`
}

func (e ExecutableStep) IsWorkoutStep() {}

type RepeatGroup struct {
	Type               string             `json:"type"` // "RepeatGroupDTO"
	StepOrder          int                `json:"stepOrder"`
	StepType           *StepTypeModel     `json:"stepType,omitempty"`
	NumberOfIterations int                `json:"numberOfIterations"`
	WorkoutSteps       []WorkoutStep      `json:"workoutSteps"`
	EndCondition       *EndConditionModel `json:"endCondition,omitempty"`
	EndConditionValue  *float64           `json:"endConditionValue,omitempty"`
	ChildStepID        *int               `json:"childStepId,omitempty"`
	SmartRepeat        bool               `json:"smartRepeat"`
}

func (r RepeatGroup) IsWorkoutStep() {}

// Equipment type key mapping from German to Garmin
var equipmentTypeIDMap = map[models.EquipmentType]EquipmentTypeID{
	models.EquipmentFins:      EquipmentFins,
	models.EquipmentKickboard: EquipmentKickboard,
	models.EquipmentPaddles:   EquipmentPaddles,
	models.EquipmentBuoy:      EquipmentBuoy,
	models.EquipmentSnorkel:   EquipmentSnorkel,
}

// detectStrokeType analyzes the row content and returns the appropriate stroke type
// TODO: Implement GenAI-based detection for accurate stroke type identification
// For now, returns StrokeTypeAny as placeholder
func detectStrokeType(row models.Row) *StrokeTypeModel {
	return &StrokeTypeModel{
		StrokeTypeID: int(StrokeTypeAny),
		DisplayOrder: 1,
	}
}

// detectEquipmentType maps the first equipment item to Garmin equipment type
func detectEquipmentType(equipment []models.EquipmentType) *EquipmentTypeModel {
	if len(equipment) == 0 {
		return nil
	}

	firstEquip := equipment[0]
	if id, ok := equipmentTypeIDMap[firstEquip]; ok {
		return &EquipmentTypeModel{
			EquipmentTypeID: int(id),
			DisplayOrder:    int(id),
		}
	}

	return nil
}

// determineStepType determines the step type based on row position and content
func determineStepType(row models.Row, rowIdx, warmupIdx, cooldownIdx int) StepTypeModel {
	contentLower := strings.ToLower(row.Content)
	intensityLower := strings.ToLower(row.Intensity)

	// Check explicit warmup/cooldown keywords first
	if strings.Contains(contentLower, "einschwimmen") || strings.Contains(intensityLower, "einschwimmen") ||
		strings.Contains(intensityLower, "warmup") {
		return StepTypeModel{
			StepTypeID:   StepTypeWarmup,
			StepTypeKey:  "warmup",
			DisplayOrder: 1,
		}
	}

	if strings.Contains(contentLower, "ausschwimmen") || strings.Contains(intensityLower, "ausschwimmen") ||
		strings.Contains(intensityLower, "cooldown") {
		return StepTypeModel{
			StepTypeID:   StepTypeCooldown,
			StepTypeKey:  "cooldown",
			DisplayOrder: 2,
		}
	}

	// Fall back to position-based detection
	if rowIdx == warmupIdx {
		return StepTypeModel{
			StepTypeID:   StepTypeWarmup,
			StepTypeKey:  "warmup",
			DisplayOrder: 1,
		}
	}

	if rowIdx == cooldownIdx {
		return StepTypeModel{
			StepTypeID:   StepTypeCooldown,
			StepTypeKey:  "cooldown",
			DisplayOrder: 2,
		}
	}

	return StepTypeModel{
		StepTypeID:   StepTypeInterval,
		StepTypeKey:  "interval",
		DisplayOrder: 3,
	}
}

// createExecutableStep creates an executable step from a row
func createExecutableStep(row models.Row, stepType StepTypeModel,
	strokeType *StrokeTypeModel, equipmentType *EquipmentTypeModel, stepOrder int) ExecutableStep {

	distance := float64(row.Distance)

	return ExecutableStep{
		Type:      "ExecutableStepDTO",
		StepOrder: stepOrder,
		StepType:  &stepType,
		EndCondition: &EndConditionModel{
			ConditionTypeID:  ConditionTypeDistance,
			ConditionTypeKey: "distance",
			DisplayOrder:     3,
			Displayable:      true,
		},
		EndConditionValue: &distance,
		TargetType: &TargetTypeModel{
			WorkoutTargetTypeID:  TargetTypeNoTarget,
			WorkoutTargetTypeKey: "no.target",
			DisplayOrder:         1,
		},
		StrokeType:    strokeType,
		EquipmentType: equipmentType,
	}
}

// createRestStep creates a rest step with time-based end condition
func createRestStep(breakSecs int, stepOrder int) ExecutableStep {
	duration := float64(breakSecs)

	return ExecutableStep{
		Type:      "ExecutableStepDTO",
		StepOrder: stepOrder,
		StepType: &StepTypeModel{
			StepTypeID:   StepTypeRest,
			StepTypeKey:  "rest",
			DisplayOrder: 5,
		},
		EndCondition: &EndConditionModel{
			ConditionTypeID:  ConditionTypeTime,
			ConditionTypeKey: "time",
			DisplayOrder:     2,
			Displayable:      true,
		},
		EndConditionValue: &duration,
		TargetType: &TargetTypeModel{
			WorkoutTargetTypeID:  TargetTypeNoTarget,
			WorkoutTargetTypeKey: "no.target",
			DisplayOrder:         1,
		},
	}
}

// convertSubRowsToSteps converts SubRows into sequential workout steps
func convertSubRowsToSteps(parentRow models.Row, stepType StepTypeModel,
	strokeType *StrokeTypeModel, equipmentType *EquipmentTypeModel) []WorkoutStep {

	var groupSteps []WorkoutStep
	subStepOrder := 1

	for _, subRow := range parentRow.SubRows {
		step := createExecutableStep(subRow, stepType, strokeType, equipmentType, subStepOrder)
		groupSteps = append(groupSteps, step)
		subStepOrder++
	}

	// Add break AFTER complete set (per user preference)
	if breakSecs := parentRow.BreakInSeconds(); breakSecs > 0 {
		restStep := createRestStep(breakSecs, subStepOrder)
		groupSteps = append(groupSteps, restStep)
	}

	return groupSteps
}

// calculateEstimatedDuration calculates estimated workout duration in seconds
func calculateEstimatedDuration(table models.Table) int {
	totalSecs := 0

	pacePer100m := map[string]int{
		"recovery":  120,
		"easy":      90,
		"warmup":    90,
		"moderate":  75,
		"hard":      60,
		"very hard": 50,
		"all out":   45,
		"fast":      55,
	}

	for _, row := range table {
		if strings.Contains(row.Content, "Gesamt") || strings.Contains(row.Content, "Total") {
			continue
		}

		intensity := strings.ToLower(row.Intensity)
		pace := pacePer100m[intensity]
		if pace == 0 {
			pace = 75
		}

		distanceMeters := row.Distance
		if len(row.SubRows) > 0 {
			for _, subRow := range row.SubRows {
				distanceMeters += subRow.Distance
			}
		}

		swimSecs := (distanceMeters / 100) * pace
		totalSecs += row.Amount * swimSecs

		if breakSecs := row.BreakInSeconds(); breakSecs > 0 {
			if len(row.SubRows) > 0 {
				totalSecs += row.Amount * breakSecs
			} else {
				totalSecs += breakSecs
			}
		}
	}

	return totalSecs
}

type WorkoutSegment struct {
	SegmentOrder int             `json:"segmentOrder"`
	SportType    *SportTypeModel `json:"sportType"`
	WorkoutSteps []WorkoutStep   `json:"workoutSteps"`
}

type BaseWorkout struct {
	WorkoutName             string           `json:"workoutName"`
	SportType               *SportTypeModel  `json:"sportType"`
	EstimatedDurationInSecs int              `json:"estimatedDurationInSecs"`
	WorkoutSegments         []WorkoutSegment `json:"workoutSegments"`
	Author                  map[string]any   `json:"author,omitempty"`
}

// NewSwimmingWorkout creates a basic swimming workout struct.
func NewSwimmingWorkout(name string, estimatedDurationSecs int, segments []WorkoutSegment) *BaseWorkout {
	return &BaseWorkout{
		WorkoutName: name,
		SportType: &SportTypeModel{
			SportTypeID:  SportTypeSwimming,
			SportTypeKey: "swimming",
			DisplayOrder: 3,
		},
		EstimatedDurationInSecs: estimatedDurationSecs,
		WorkoutSegments:         segments,
		Author:                  make(map[string]any),
	}
}

// Helper functions for common steps

func CreateWarmupStep(durationSeconds float64, stepOrder int, targetType *TargetTypeModel) ExecutableStep {
	if targetType == nil {
		targetType = &TargetTypeModel{WorkoutTargetTypeID: TargetTypeNoTarget, WorkoutTargetTypeKey: "no.target", DisplayOrder: 1}
	}
	return ExecutableStep{
		Type:      "ExecutableStepDTO",
		StepOrder: stepOrder,
		StepType: &StepTypeModel{
			StepTypeID:   StepTypeWarmup,
			StepTypeKey:  "warmup",
			DisplayOrder: 1,
		},
		EndCondition: &EndConditionModel{
			ConditionTypeID:  ConditionTypeTime,
			ConditionTypeKey: "time",
			DisplayOrder:     2,
			Displayable:      true,
		},
		EndConditionValue: &durationSeconds,
		TargetType:        targetType,
	}
}

func CreateIntervalStep(durationSeconds float64, stepOrder int, targetType *TargetTypeModel) ExecutableStep {
	if targetType == nil {
		targetType = &TargetTypeModel{WorkoutTargetTypeID: TargetTypeNoTarget, WorkoutTargetTypeKey: "no.target", DisplayOrder: 1}
	}
	return ExecutableStep{
		Type:      "ExecutableStepDTO",
		StepOrder: stepOrder,
		StepType: &StepTypeModel{
			StepTypeID:   StepTypeInterval,
			StepTypeKey:  "interval",
			DisplayOrder: 3,
		},
		EndCondition: &EndConditionModel{
			ConditionTypeID:  ConditionTypeTime,
			ConditionTypeKey: "time",
			DisplayOrder:     2,
			Displayable:      true,
		},
		EndConditionValue: &durationSeconds,
		TargetType:        targetType,
	}
}

func CreateRecoveryStep(durationSeconds float64, stepOrder int, targetType *TargetTypeModel) ExecutableStep {
	if targetType == nil {
		targetType = &TargetTypeModel{WorkoutTargetTypeID: TargetTypeNoTarget, WorkoutTargetTypeKey: "no.target", DisplayOrder: 1}
	}
	return ExecutableStep{
		Type:      "ExecutableStepDTO",
		StepOrder: stepOrder,
		StepType: &StepTypeModel{
			StepTypeID:   StepTypeRecovery,
			StepTypeKey:  "recovery",
			DisplayOrder: 4,
		},
		EndCondition: &EndConditionModel{
			ConditionTypeID:  ConditionTypeTime,
			ConditionTypeKey: "time",
			DisplayOrder:     2,
			Displayable:      true,
		},
		EndConditionValue: &durationSeconds,
		TargetType:        targetType,
	}
}

func CreateCooldownStep(durationSeconds float64, stepOrder int, targetType *TargetTypeModel) ExecutableStep {
	if targetType == nil {
		targetType = &TargetTypeModel{WorkoutTargetTypeID: TargetTypeNoTarget, WorkoutTargetTypeKey: "no.target", DisplayOrder: 1}
	}
	return ExecutableStep{
		Type:      "ExecutableStepDTO",
		StepOrder: stepOrder,
		StepType: &StepTypeModel{
			StepTypeID:   StepTypeCooldown,
			StepTypeKey:  "cooldown",
			DisplayOrder: 2,
		},
		EndCondition: &EndConditionModel{
			ConditionTypeID:  ConditionTypeTime,
			ConditionTypeKey: "time",
			DisplayOrder:     2,
			Displayable:      true,
		},
		EndConditionValue: &durationSeconds,
		TargetType:        targetType,
	}
}

func CreateRepeatGroup(iterations int, workoutSteps []WorkoutStep, stepOrder int) RepeatGroup {
	fv := float64(iterations)
	return RepeatGroup{
		Type:      "RepeatGroupDTO",
		StepOrder: stepOrder,
		StepType: &StepTypeModel{
			StepTypeID:   StepTypeRepeat,
			StepTypeKey:  "repeat",
			DisplayOrder: 6,
		},
		NumberOfIterations: iterations,
		WorkoutSteps:       workoutSteps,
		EndCondition: &EndConditionModel{
			ConditionTypeID:  ConditionTypeIterations,
			ConditionTypeKey: "iterations",
			DisplayOrder:     7,
			Displayable:      false,
		},
		EndConditionValue: &fv,
	}
}

func ConvertTrainingPlanToSwimWorkout(p *models.Plan, poolLength int) *BaseWorkout {
	var steps []WorkoutStep
	stepOrder := 1

	totalSecs := calculateEstimatedDuration(p.Table)

	warmupIndex := 0
	cooldownIndex := len(p.Table) - 1

	rowIdx := 0
	for _, row := range p.Table {
		if strings.Contains(row.Content, "Gesamt") || strings.Contains(row.Content, "Total") {
			continue
		}

		stepType := determineStepType(row, rowIdx, warmupIndex, cooldownIndex)
		strokeType := detectStrokeType(row)
		equipmentType := detectEquipmentType(row.Equipment)

		if len(row.SubRows) > 0 {
			groupSteps := convertSubRowsToSteps(row, stepType, strokeType, equipmentType)

			fv := float64(row.Amount)
			group := RepeatGroup{
				Type:      "RepeatGroupDTO",
				StepOrder: stepOrder,
				StepType: &StepTypeModel{
					StepTypeID:   StepTypeRepeat,
					StepTypeKey:  "repeat",
					DisplayOrder: 6,
				},
				NumberOfIterations: row.Amount,
				WorkoutSteps:       groupSteps,
				EndCondition: &EndConditionModel{
					ConditionTypeID:  ConditionTypeIterations,
					ConditionTypeKey: "iterations",
					DisplayOrder:     7,
					Displayable:      false,
				},
				EndConditionValue: &fv,
				SmartRepeat:       false,
			}

			steps = append(steps, group)
			stepOrder += len(groupSteps) + 1
		} else if row.Amount > 1 {
			step := createExecutableStep(row, stepType, strokeType, equipmentType, 1)

			var groupSteps []WorkoutStep
			groupSteps = append(groupSteps, step)

			if breakSecs := row.BreakInSeconds(); breakSecs > 0 {
				restStep := createRestStep(breakSecs, 2)
				groupSteps = append(groupSteps, restStep)
			}

			fv := float64(row.Amount)
			group := RepeatGroup{
				Type:      "RepeatGroupDTO",
				StepOrder: stepOrder,
				StepType: &StepTypeModel{
					StepTypeID:   StepTypeRepeat,
					StepTypeKey:  "repeat",
					DisplayOrder: 6,
				},
				NumberOfIterations: row.Amount,
				WorkoutSteps:       groupSteps,
				EndCondition: &EndConditionModel{
					ConditionTypeID:  ConditionTypeIterations,
					ConditionTypeKey: "iterations",
					DisplayOrder:     7,
					Displayable:      false,
				},
				EndConditionValue: &fv,
				SmartRepeat:       false,
			}

			steps = append(steps, group)
			stepOrder += len(groupSteps) + 1
		} else {
			step := createExecutableStep(row, stepType, strokeType, equipmentType, stepOrder)
			steps = append(steps, step)
			stepOrder++

			if breakSecs := row.BreakInSeconds(); breakSecs > 0 {
				restStep := createRestStep(breakSecs, stepOrder)
				steps = append(steps, restStep)
				stepOrder++
			}
		}

		rowIdx++
	}

	segment := WorkoutSegment{
		SegmentOrder: 1,
		SportType: &SportTypeModel{
			SportTypeID:  SportTypeSwimming,
			SportTypeKey: "swimming",
			DisplayOrder: 3,
		},
		WorkoutSteps: steps,
	}

	return &BaseWorkout{
		WorkoutName:             p.Title,
		SportType:               segment.SportType,
		EstimatedDurationInSecs: totalSecs,
		WorkoutSegments:         []WorkoutSegment{segment},
		Author:                  make(map[string]any),
	}
}
