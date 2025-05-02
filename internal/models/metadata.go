package models

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
)

type DBMethod string

const (
	MethodUnixSocket DBMethod = "unix"
	MethodTCP        DBMethod = "tcp"
	MethodURL        DBMethod = "url"
)

type DifficultyLevel string

const (
	Nichtschwimmer     DifficultyLevel = "Nichtschwimmer"
	Anfaenger          DifficultyLevel = "Anfaenger"
	Fortgeschritten    DifficultyLevel = "Fortgeschritten"
	Leistungsschwimmer DifficultyLevel = "Leistungsschwimmer"
	Athlet             DifficultyLevel = "Top-Athlet"
)

type Training string

const (
	Techniktraining       Training = "Techniktraining"
	Leistungstest         Training = "Leistungstest"
	Grundlagen            Training = "Grundlagenausdauer"
	Recovery              Training = "Recovery"
	Kurzstrecken          Training = "Kurzstrecken"
	Langstrecken          Training = "Langstrecken"
	Atemmangel            Training = "Atemmangel"
	Wettkampfvorbereitung Training = "Wettkampfvorbereitung"
)

type Metadata struct {
	HasFreestyle    bool            `json:"freistil" jsonschema_description:"Indicates if the training explicitely includes Freistil (Crawl)"`
	HasBreaststroke bool            `json:"brust" jsonschema_description:"Indicates if the training explicitely includes Brust (Breaststroke)"`
	HasButterfly    bool            `json:"delfin" jsonschema_description:"Indicates if the training explicitely includes Delfin (Butterfly)"`
	HasBackstroke   bool            `json:"ruecken" jsonschema_description:"Indicates if the training explicitely includes Ruecken (Backstroke)"`
	HasMeddley      bool            `json:"lagen" jsonschema_description:"Indicates if the training explicitely includes Lagen (Medley)"`
	Difficulty      DifficultyLevel `json:"schwierigkeitsgrad" jsonschema:"description=The difficulty level of the training,enum=Nichtschwimmer,enum=Anfaenger,enum=Fortgeschritten,enum=Leistungsschwimmer,enum=Top-Athlet"`
	TrainingType    Training        `json:"trainingstyp" jsonschema:"description=The type of training,enum=Techniktraining,enum=Leistungstest,enum=Grundlagen,enum=Recovery,enum=Kurzstrecken,enum=Langstrecken,enum=Atemmangel,enum=Wettkampfvorbereitung"`
	Reasoning       string          `json:"Begründung" jsonschema_description:"Reasoning for why the attributes were chosen"`
	Equipment       []string        `json:"Ausrüstung" jsonschema_description:"Equipment needed for the training"`
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
