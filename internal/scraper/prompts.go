package scraper

const scrapeTemplateStr = `
Deine Aufgabe ist die Klassifikation von Schwimmtrainingsplänen. 
Du bekommst eine Reihe von Informationen über den Trainingsplan und 
sollst eine strukturierte Antwort mit Deutschem Text wiedergeben.
Typischerweise wird bei Trainingsplänen mit GA (Grundlagenausdauer) in
den jeweiligen Abschnitten Kraul geschwommen.
Dabei sind folgende Abkürzungen relevant:
{
    "K": "Kraulschwimmen",
    "Kr": "Kraulschwimmen",
    "Freistil": "Kraulschwimmen",
    "F": "Kraulschwimmen",
    "Fr": "Kraulschwimmen",
    "R": "Rückenschwimmen",
    "B": "Brustschwimmen",
    "Br": "Brustschwimmen",
    "Be", "Beinarbeit",
    "S": "Schmetterling/Delfinschwimmen",
    "D": "Schmetterling/Delfinschwimmen",
    "Lagen": "Lagenstaffel",
    "DL": "Dauerlauf",
    "SL": "Sprintlauf",
    "GA": "Grundlagenausdauer",
    "SA": "Schnelligkeitsausdauer",
    "TA": "Technikausdauer",
    "TÜ": "Technische Übung",
    "TS": "Technisch Sauber"
}

Die folgenden Informationen sind zu beachten:

Titel:
%s

Beschreibung:
%s

Tabelle:
%s

Extrahiere die folgenden Informationen aus der Tabelle entsprechend dieses JSON-Schemas
und gib deine Antwort als JSON zurück:

%s

Antwort:
`
