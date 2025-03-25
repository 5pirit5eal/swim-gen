package rag

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

type Config struct {
	ProjectID string `env:"PROJECT_ID"`
	Region    string `env:"REGION"`
	Model     string `env:"MODEL"`
	Embedding struct {
		Name  string `env:"EMBEDDING_NAME"`
		Model string `env:"EMBEDDING_MODEL"`
		Size  int    `env:"EMBEDDING_SIZE"`
	}

	DB struct {
		Name         string   `env:"DB_NAME"`
		Instance     string   `env:"DB_INSTANCE"`
		Port         string   `env:"DB_PORT"`
		User         string   `env:"DB_USER"`
		Pass         string   `env:"DB_PASS"`
		PassLocation string   `env:"DB_PASS_LOCATION"`
		Method       DBMethod `env:"DB_METHOD"`
	}
}

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
	Begründung         string
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
