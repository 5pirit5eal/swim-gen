package genai

const ragTemplateStr = `
Du bist ein Schwimmtrainer und hilfst einem Schwimmer einen Trainingsplan zu erstellen.
Du bekommst eine Frage vom Schwimmer und du hast eine Liste von Trainingsplänen als Kontext.
Die Trainingspläne beinhalten Downloadlinks und nummerierte Titel. Mit diesen kann der Schwimmer nichts anfangen. Entsprechend sollten
nicht mit enthalten sein. Entferne deshalb "www.docswim.de" oder "EIN TRAININGSPLAN von DOC SWIM".
Erstelle dem Schwimmer einen passenden Trainingsplan basierend auf dem Kontext. Dafür kannst du die Trainingspläne im Kontext verwenden,
indem du sie selektierst, kombinierst, mischst, oder umformulierst, um sie an die Bedürfnisse des Schwimmers anzupassen.
Bei der Erstellung der Beschreibung gehe nur auf die Eigenschaften des Trainingsplans ein und erkläre dem Schwimmer, wofür der Trainingsplan geeignet ist.
Für den Schwimmer ist nicht relevant, ob der Plan aus mehreren oder einem anderen Trainingsplan erstellt wurde.
Die Antwort soll in Deutsch sein.
Die Antwort soll in JSON-Format sein.
Die Antwort soll die folgenden Felder enthalten:
{
	"title": "Ein passender, kurzer, prägnanter Titel des Trainings",
	"description": "Eine kurze Beschreibung, Kommentare oder Anmerkungen zu dem Trainingsplan, damit der Schwimmer den Plan besser versteht",
	"table": Eine Tabelle mit den Trainingsdaten nach dem unten stehenden Schema
}

table Schema:
%s

Die Antwort soll keine Fragen enthalten und auch nicht die Anweisung wiederholen.

Frage:
%s

Kontext:
%s
`

const choosePlanTemplateStr = `
Du bist ein Schwimmtrainer und hilfst einem Schwimmer einen Trainingsplan auszusuchen.
Du bekommst eine Frage vom Schwimmer und du hast eine Liste von Trainingsplänen als Kontext.
Wähle den besten Trainingsplan aus dem Kontext aus, der am besten zu der Frage passt.
Die Antwort soll in Deutsch sein.
Die Antwort soll in JSON-Format sein.
{
	"description": "Eine kurze Beschreibung, Kommentare oder Anmerkungen zu dem Trainingsplan, damit der Schwimmer den Plan besser versteht",
	"index": "Der Index des Trainingsplans in der Liste als integer",
}

Die Antwort soll keine Fragen enthalten und auch nicht die Anweisung wiederholen.

Frage:
%s
Kontext:
%s
`

const metadataTemplateStr = `
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
    "Be": "Beinarbeit",
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

const describeTemplateStr = `
Du bist ein Schwimmtrainer und bekommst einen Trainingsplan von deinem Assitenten vorgelegt.
Du sollst den Trainingsplan beschreiben und die wichtigsten Eigenschaften benennen.
Analysiere für welche Schwimmer und für welche Trainingsziele der Trainingsplan geeignet ist.
Die Antwort soll in Deutsch sein.
Die Antwort soll in JSON-Format sein und in folgendem Schema:
%s

Tabelle:
%s

Antwort:
`
