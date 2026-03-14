package garmin

// 1. Implement conversion of training plan to FIT factory
// 2. Implement encoding of fit files
// 3. Implement endpoints for backend for fit file creation

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/muktihari/fit/encoder"
	"github.com/muktihari/fit/profile/filedef"
	"github.com/muktihari/fit/profile/mesgdef"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/fit/proto"
)

func ConvertTrainingPlanToFit(p *models.Plan, poolLength int) (*proto.FIT, error) {
	now := time.Now()
	workout := filedef.NewWorkout()
	workout.FileId.
		SetTimeCreated(now).
		SetManufacturer(typedef.ManufacturerDevelopment).
		SetSerialNumber(uint32(now.Unix()))

	workout.Workout = mesgdef.NewWorkout(nil).
		SetSport(typedef.SportSwimming).
		SetSubSport(typedef.SubSportLapSwimming).
		SetNumValidSteps(uint16(len(p.Table) - 1)).
		SetWktName(p.Title).
		SetWktDescription(p.Description).
		SetPoolLength(uint16(poolLength * 100)).
		SetPoolLengthUnit(typedef.DisplayMeasureMetric)

	// Iterate over rows to add workout steps
	for _, row := range p.Table {
		if strings.Contains(row.Content, "Gesamt") || strings.Contains(row.Content, "Total") {
			continue
		}
		// First two words in the content
		words := strings.Fields(row.Content)
		short_content := ""
		if len(words) >= 2 {
			short_content = words[0] + " " + words[1]
		} else {
			short_content = row.Content
		}
		// Create a break time in s by parsing the break string
		break_time := row.BreakInSeconds()
		startIndex := len(workout.WorkoutSteps)

		workout.WorkoutSteps = append(
			workout.WorkoutSteps,
			// working set
			mesgdef.NewWorkoutStep(nil).
				SetMessageIndex(typedef.MessageIndex(startIndex)).
				SetWktStepName(short_content).
				SetDurationType(typedef.WktStepDurationDistance).
				SetDurationValue(uint32(row.Distance*100)).
				SetTargetType(typedef.WktStepTargetSwimStroke).
				SetTargetValue(uint32(typedef.SwimStrokeFreestyle)).
				SetIntensity(typedef.IntensityActive).
				SetEquipment(typedef.WorkoutEquipmentSwimPullBuoy).
				SetNotes(row.Content),
			// rest set
			mesgdef.NewWorkoutStep(nil).
				SetMessageIndex(typedef.MessageIndex(startIndex+1)).
				SetWktStepName("Pause").
				SetDurationType(typedef.WktStepDurationTime).
				SetDurationValue(uint32(break_time)).
				SetTargetType(typedef.WktStepTargetSwimStroke).
				SetTargetValue(uint32(typedef.SwimStrokeFreestyle)).
				SetIntensity(typedef.IntensityRest).
				SetNotes(row.Content),
		)
		// repeat if necessary
		if row.Amount > 1 {
			workout.WorkoutSteps = append(
				workout.WorkoutSteps,
				mesgdef.NewWorkoutStep(nil).
					SetMessageIndex(typedef.MessageIndex(startIndex+2)).
					SetWktStepName("Repeat").
					SetDurationType(typedef.WktStepDurationRepeatUntilStepsCmplt).
					SetDurationValue(uint32(startIndex)).
					SetTargetValue(uint32(row.Amount)),
			)
		}
	}

	fit := workout.ToFIT(nil)
	return &fit, nil
}

func EncodeFITFile(ctx context.Context, w io.Writer, fit *proto.FIT) error {
	enc := encoder.New(w, encoder.WithProtocolVersion(proto.V2))
	return enc.EncodeWithContext(ctx, fit)
}
