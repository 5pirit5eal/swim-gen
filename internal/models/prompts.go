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

type Schwierigkeitsgrad string

const (
	Nichtschwimmer     Schwierigkeitsgrad = "Nichtschwimmer"
	Anfaenger          Schwierigkeitsgrad = "Anfaenger"
	Fortgeschritten    Schwierigkeitsgrad = "Fortgeschritten"
	Leistungsschwimmer Schwierigkeitsgrad = "Leistungsschwimmer"
	Athlet             Schwierigkeitsgrad = "Top-Athlet"
)

type Trainingstyp string

const (
	Techniktraining       Trainingstyp = "Techniktraining"
	Leistungstest         Trainingstyp = "Leistungstest"
	Grundlagen            Trainingstyp = "Grundlagenausdauer"
	Recovery              Trainingstyp = "Recovery"
	Kurzstrecken          Trainingstyp = "Kurzstrecken"
	Langstrecken          Trainingstyp = "Langstrecken"
	Atemmangel            Trainingstyp = "Atemmangel"
	Wettkampfvorbereitung Trainingstyp = "Wettkampfvorbereitung"
)

type Metadata struct {
	Freistil           bool               `json:"freistil"`
	Brust              bool               `json:"brust"`
	Delfin             bool               `json:"delfin"`
	Ruecken            bool               `json:"ruecken"`
	Lagen              bool               `json:"lagen"`
	Schwierigkeitsgrad Schwierigkeitsgrad `json:"schwierigkeitsgrad"`
	Trainingstyp       Trainingstyp       `json:"trainingstyp"`
	Begründung         string             `json:"Begründung"`
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
				"enum": []Schwierigkeitsgrad{
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
				"enum": []Trainingstyp{
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
				"amount": map[string]interface{}{
					"type":        "integer",
					"description": "The number of repetitions or sets",
				},
				"multiplier": map[string]interface{}{
					"type":        "string",
					"description": "The multiplier for the distance (e.g., 'x', 'times')",
				},
				"distance": map[string]interface{}{
					"type":        "integer",
					"description": "The distance in meters",
				},
				"break": map[string]interface{}{
					"type":        "string",
					"description": "The break time in seconds",
				},
				"content": map[string]interface{}{
					"type":        "string",
					"description": "The content or description of the row",
				},
				"intensity": map[string]interface{}{
					"type":        "string",
					"description": "The intensity level of the activity",
				},
				"sum": map[string]interface{}{
					"type":        "integer",
					"description": "The total volume or sum for the row",
				},
			},
			"required": []string{
				"amount",
				"multiplier",
				"distance",
				"break",
				"content",
				"intensity",
				"sum",
			},
		},
	}

	jsonSchema, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON schema: %w", err)
	}

	return string(jsonSchema), nil
}
