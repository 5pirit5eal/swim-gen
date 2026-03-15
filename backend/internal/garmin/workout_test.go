package garmin_test

import (
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/garmin"
	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestConvertTrainingPlanToSwimWorkout(t *testing.T) {
	plan := &models.Plan{
		Title: "Test Plan",
		Table: models.Table{
			{Amount: 1, Distance: 400, Content: "Einschwimmen", Intensity: "Warmup"},
			{Amount: 4, Distance: 100, Content: "Main", Intensity: "Hard", Break: "20s"},
			{Amount: 1, Distance: 200, Content: "Ausschwimmen", Intensity: "Cooldown"},
			{Content: "Gesamt", Sum: 1000},
		},
	}

	workout := garmin.ConvertTrainingPlanToSwimWorkout(plan, 25)

	assert.NotNil(t, workout)
	assert.Equal(t, "Test Plan", workout.WorkoutName)
	assert.Equal(t, garmin.SportTypeSwimming, workout.SportType.SportTypeID)
	assert.Len(t, workout.WorkoutSegments, 1)

	segments := workout.WorkoutSegments[0]
	assert.Len(t, segments.WorkoutSteps, 3)

	// Step 1: Warmup (no repeat, no break)
	step1, ok := segments.WorkoutSteps[0].(garmin.ExecutableStep)
	assert.True(t, ok)
	assert.Equal(t, garmin.StepTypeWarmup, step1.StepType.StepTypeID)
	assert.Equal(t, float64(400), *step1.EndConditionValue)
	assert.Equal(t, 1, step1.StepOrder)

	// Step 2: Main grouped (repeat)
	step2, ok := segments.WorkoutSteps[1].(garmin.RepeatGroup)
	assert.True(t, ok)
	assert.Equal(t, garmin.StepTypeRepeat, step2.StepType.StepTypeID)
	assert.Equal(t, 4, step2.NumberOfIterations)
	assert.Len(t, step2.WorkoutSteps, 2)
	assert.Equal(t, 2, step2.StepOrder)

	groupStep1, ok := step2.WorkoutSteps[0].(garmin.ExecutableStep)
	assert.True(t, ok)
	assert.Equal(t, garmin.StepTypeInterval, groupStep1.StepType.StepTypeID)
	assert.Equal(t, float64(100), *groupStep1.EndConditionValue)
	assert.Equal(t, 3, groupStep1.StepOrder)

	groupStep2, ok := step2.WorkoutSteps[1].(garmin.ExecutableStep)
	assert.True(t, ok)
	assert.Equal(t, garmin.StepTypeRest, groupStep2.StepType.StepTypeID)
	assert.Equal(t, float64(20), *groupStep2.EndConditionValue)
	assert.Equal(t, 4, groupStep2.StepOrder)

	// Step 3: Cooldown (no repeat, no break)
	step3, ok := segments.WorkoutSteps[2].(garmin.ExecutableStep)
	assert.True(t, ok)
	assert.Equal(t, garmin.StepTypeCooldown, step3.StepType.StepTypeID)
	assert.Equal(t, float64(200), *step3.EndConditionValue)
	assert.Equal(t, 5, step3.StepOrder)
}

func TestWorkoutHelperFunctions(t *testing.T) {
	t.Run("CreateWarmupStep", func(t *testing.T) {
		step := garmin.CreateWarmupStep(300, 1, nil)
		assert.Equal(t, garmin.StepTypeWarmup, step.StepType.StepTypeID)
		assert.Equal(t, float64(300), *step.EndConditionValue)
		assert.Equal(t, garmin.ConditionTypeTime, step.EndCondition.ConditionTypeID)
		assert.Equal(t, 1, step.StepOrder)
	})

	t.Run("CreateIntervalStep", func(t *testing.T) {
		step := garmin.CreateIntervalStep(60, 2, nil)
		assert.Equal(t, garmin.StepTypeInterval, step.StepType.StepTypeID)
		assert.Equal(t, float64(60), *step.EndConditionValue)
		assert.Equal(t, 2, step.StepOrder)
	})

	t.Run("CreateRecoveryStep", func(t *testing.T) {
		step := garmin.CreateRecoveryStep(15, 3, nil)
		assert.Equal(t, garmin.StepTypeRecovery, step.StepType.StepTypeID)
		assert.Equal(t, float64(15), *step.EndConditionValue)
		assert.Equal(t, 3, step.StepOrder)
	})

	t.Run("CreateCooldownStep", func(t *testing.T) {
		step := garmin.CreateCooldownStep(120, 4, nil)
		assert.Equal(t, garmin.StepTypeCooldown, step.StepType.StepTypeID)
		assert.Equal(t, float64(120), *step.EndConditionValue)
		assert.Equal(t, 4, step.StepOrder)
	})

	t.Run("CreateRepeatGroup", func(t *testing.T) {
		steps := []garmin.WorkoutStep{
			garmin.CreateIntervalStep(60, 2, nil),
			garmin.CreateRecoveryStep(15, 3, nil),
		}
		group := garmin.CreateRepeatGroup(5, steps, 1)
		assert.Equal(t, 5, group.NumberOfIterations)
		assert.Equal(t, float64(5), *group.EndConditionValue)
		assert.Len(t, group.WorkoutSteps, 2)
		assert.Equal(t, 1, group.StepOrder)
	})
}
