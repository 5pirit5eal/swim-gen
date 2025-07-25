package models

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
)

// DifficultyLevel represents the difficulty level of a training plan
// @Description The difficulty level categorization for swimming training plans
type DifficultyLevel string

const (
	Nichtschwimmer     DifficultyLevel = "Nichtschwimmer"     // @Description Non-swimmer level
	Anfaenger          DifficultyLevel = "Anfaenger"          // @Description Beginner level
	Fortgeschritten    DifficultyLevel = "Fortgeschritten"    // @Description Advanced level
	Leistungsschwimmer DifficultyLevel = "Leistungsschwimmer" // @Description Competitive swimmer level
	Athlet             DifficultyLevel = "Top-Athlet"         // @Description Elite athlete level
)

// Training represents the type of training focus
// @Description The specific training type or focus area for the swimming plan
type Training string

const (
	Techniktraining       Training = "Techniktraining"       // @Description Technique training focus
	Leistungstest         Training = "Leistungstest"         // @Description Performance testing
	Grundlagen            Training = "Grundlagenausdauer"    // @Description Basic endurance training
	Recovery              Training = "Recovery"              // @Description Recovery/easy training
	Kurzstrecken          Training = "Kurzstrecken"          // @Description Short distance training
	Langstrecken          Training = "Langstrecken"          // @Description Long distance training
	Atemmangel            Training = "Atemmangel"            // @Description Hypoxic/breath control training
	Wettkampfvorbereitung Training = "Wettkampfvorbereitung" // @Description Competition preparation
)

// Metadata represents the metadata associated with a training plan
// @Description Detailed metadata and categorization for swimming training plans
type Metadata struct {
	HasFreestyle    bool            `json:"freistil" jsonschema_description:"Indicates if the training explicitely includes Freistil (Crawl)" example:"true"`
	HasBreaststroke bool            `json:"brust" jsonschema_description:"Indicates if the training explicitely includes Brust (Breaststroke)" example:"false"`
	HasButterfly    bool            `json:"delfin" jsonschema_description:"Indicates if the training explicitely includes Delfin (Butterfly)" example:"false"`
	HasBackstroke   bool            `json:"ruecken" jsonschema_description:"Indicates if the training explicitely includes Ruecken (Backstroke)" example:"false"`
	HasMeddley      bool            `json:"lagen" jsonschema_description:"Indicates if the training explicitely includes Lagen (Medley)" example:"false"`
	Difficulty      DifficultyLevel `json:"schwierigkeitsgrad" jsonschema:"description=The difficulty level of the training,enum=Nichtschwimmer,enum=Anfaenger,enum=Fortgeschritten,enum=Leistungsschwimmer,enum=Top-Athlet" example:"Fortgeschritten"`
	TrainingType    Training        `json:"trainingstyp" jsonschema:"description=The type of training,enum=Techniktraining,enum=Leistungstest,enum=Grundlagen,enum=Recovery,enum=Kurzstrecken,enum=Langstrecken,enum=Atemmangel,enum=Wettkampfvorbereitung" example:"Techniktraining"`
	Reasoning       string          `json:"Begründung" jsonschema_description:"Reasoning for why the attributes were chosen" example:"This plan focuses on freestyle technique improvement with moderate intensity"`
	Equipment       []string        `json:"Ausrüstung" jsonschema_description:"Equipment needed for the training" example:"[\"Kickboard\", \"Pull buoy\"]"`
}

type Description struct {
	Title string    `json:"title" jsonschema_description:"Title of the training plan"`
	Text  string    `json:"text" jsonschema_description:"Description of the training plan"`
	Meta  *Metadata `json:"metadata" jsonschema_description:"Metadata for the training plan"`
}

func MetadataSchema() (string, error) {
	schema := jsonschema.Reflect(&Metadata{})
	jsonSchema, err := json.Marshal(schema)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON schema: %w", err)
	}
	return string(jsonSchema), nil
}

func TableSchema() (string, error) {
	schema := jsonschema.Reflect(&Table{})

	jsonSchema, err := json.Marshal(schema)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON schema: %w", err)
	}

	return string(jsonSchema), nil
}

func DescriptionSchema() (string, error) {
	schema := jsonschema.Reflect(&Description{})
	jsonSchema, err := json.Marshal(schema)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON schema: %w", err)
	}
	return string(jsonSchema), nil
}
