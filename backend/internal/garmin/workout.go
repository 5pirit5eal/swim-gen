package garmin

// Garmin Connect Workout structs ported from python-garminconnect:
// https://github.com/cyberjunky/python-garminconnect/blob/master/garminconnect/workout.py

type SportTypeID int

const (
	SportTypeRunning          SportTypeID = 1
	SportTypeCycling          SportTypeID = 2
	SportTypeSwimming         SportTypeID = 3
	SportTypeWalking          SportTypeID = 4
	SportTypeMultiSport       SportTypeID = 5
	SportTypeFitnessEquipment SportTypeID = 6
	SportTypeHiking           SportTypeID = 7
	SportTypeOther            SportTypeID = 8
)

type StepTypeID int

const (
	StepTypeWarmup   StepTypeID = 1
	StepTypeCooldown StepTypeID = 2
	StepTypeInterval StepTypeID = 3
	StepTypeRecovery StepTypeID = 4
	StepTypeRest     StepTypeID = 5
	StepTypeRepeat   StepTypeID = 6
)

type ConditionTypeID int

const (
	ConditionTypeDistance   ConditionTypeID = 1
	ConditionTypeTime       ConditionTypeID = 2
	ConditionTypeHeartRate  ConditionTypeID = 3
	ConditionTypeCalories   ConditionTypeID = 4
	ConditionTypeCadence    ConditionTypeID = 5
	ConditionTypePower      ConditionTypeID = 6
	ConditionTypeIterations ConditionTypeID = 7
)

type TargetTypeID int

const (
	TargetTypeNoTarget  TargetTypeID = 1
	TargetTypeHeartRate TargetTypeID = 2
	TargetTypeCadence   TargetTypeID = 3
	TargetTypeSpeed     TargetTypeID = 4
	TargetTypePower     TargetTypeID = 5
	TargetTypeOpen      TargetTypeID = 6
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

type EquipmentTypeModel struct {
	EquipmentTypeID int `json:"equipmentTypeId"`
	DisplayOrder    int `json:"displayOrder"`
}

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
