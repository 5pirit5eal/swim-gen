package models

import (
	"encoding/json"
	"fmt"
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
	HasFreestyle    bool            `json:"freistil"`
	HasBreaststroke bool            `json:"brust"`
	HasButterfly    bool            `json:"delfin"`
	HasBackstroke   bool            `json:"ruecken"`
	HasMeddley      bool            `json:"lagen"`
	Difficulty      DifficultyLevel `json:"schwierigkeitsgrad"`
	TrainingType    Training        `json:"trainingstyp"`
	Reasoning       string          `json:"Begründung"`
}

func MetadataSchema() (string, error) {
	schema := map[string]interface{}{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"title":   "Metadata",
		"type":    "object",
		"properties": map[string]interface{}{
			"freistil": map[string]interface{}{
				"type":        "boolean",
				"description": "Indicates if the training explicitely includes Freistil (Crawl)",
			},
			"brust": map[string]interface{}{
				"type":        "boolean",
				"description": "Indicates if the training explicitely includes Brust (Breaststroke)",
			},
			"delfin": map[string]interface{}{
				"type":        "boolean",
				"description": "Indicates if the training explicitely includes Delfin (Butterfly)",
			},
			"ruecken": map[string]interface{}{
				"type":        "boolean",
				"description": "Indicates if the training explicitely includes Ruecken (Backstroke)",
			},
			"lagen": map[string]interface{}{
				"type":        "boolean",
				"description": "Indicates if the training explicitely includes Lagen (Medley)",
			},
			"schwierigkeitsgrad": map[string]interface{}{
				"type": "string",
				"enum": []DifficultyLevel{
					Nichtschwimmer,
					Anfaenger,
					Fortgeschritten,
					Leistungsschwimmer,
					Athlet,
				},
				"description": "The difficulty level of the training. Consider total volume, swim techniques, intensity and breaks",
			},
			"trainingstyp": map[string]interface{}{
				"type": "string",
				"enum": []Training{
					Techniktraining,
					Leistungstest,
					Grundlagen,
					Recovery,
					Kurzstrecken,
					Langstrecken,
					Atemmangel,
					Wettkampfvorbereitung,
				},
				"description": "The type of training",
			},
			"Begründung": map[string]interface{}{
				"type":        "string",
				"description": "Reasoning for why the LLM added what it did",
			},
		},
		"required": []string{
			"freistil",
			"brust",
			"delfin",
			"ruecken",
			"lagen",
			"schwierigkeitsgrad",
			"trainingstyp",
			"Begründung",
		},
	}

	jsonSchema, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON schema: %w", err)
	}

	return string(jsonSchema), nil
}

func TableSchema() (string, error) {
	schema := map[string]interface{}{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"title":   "Table",
		"type":    "array",
		"items": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"Amount": map[string]interface{}{
					"type":        "integer",
					"description": "The number of repetitions or sets",
				},
				"Multiplier": map[string]interface{}{
					"type":        "string",
					"description": "The multiplier for the distance (e.g., 'x', 'times')",
				},
				"Distance": map[string]interface{}{
					"type":        "integer",
					"description": "The distance in meters",
				},
				"Break": map[string]interface{}{
					"type":        "string",
					"description": "The break time in seconds",
				},
				"Content": map[string]interface{}{
					"type":        "string",
					"description": "The content or description of the row",
				},
				"Intensity": map[string]interface{}{
					"type":        "string",
					"description": "The intensity level of the activity",
				},
				"Sum": map[string]interface{}{
					"type":        "integer",
					"description": "The total volume or sum for the row",
				},
			},
			"required": []string{
				"Amount",
				"Multiplier",
				"Distance",
				"Break",
				"Content",
				"Intensity",
				"Sum",
			},
		},
	}

	jsonSchema, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON schema: %w", err)
	}

	return string(jsonSchema), nil
}
