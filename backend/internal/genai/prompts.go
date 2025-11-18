package genai

const ragTemplateStr = `
Du bist ein Schwimmtrainer und hilfst einem Schwimmer einen Trainingsplan zu erstellen.
Du bekommst eine Frage vom Schwimmer und du hast eine Liste von Trainingsplänen als Referenz.
Die Trainingspläne beinhalten Downloadlinks und nummerierte Titel. Mit diesen kann der Schwimmer nichts anfangen. Entsprechend sollten
diese nicht mit enthalten sein. Entferne deshalb "www.docswim.de" oder "EIN TRAININGSPLAN von DOC SWIM".
Erstelle dem Schwimmer einen passenden Trainingsplan basierend auf dem Kontext und seiner Anfrage. Dafür kannst du die Referenztrainingspläne verwenden,
indem du sie selektierst, kombinierst, mischst, oder umformulierst, um sie an die Bedürfnisse des Schwimmers anzupassen.
Ziehe dabei auch die konfigurierte Poollänge in Betracht: %s. Die Standard-Poollänge ist 25m.
Achte darauf, dass die Gesamtdistanz des Trainingsplans zu der Anfrage des Schwimmers passt.
Erhöhe die Anzahl der Wiederholungen, oder die Distanz der einzelnen Wiederholungen, um die Gesamtdistanz anzupassen.
Bei der Erstellung der kurzen Beschreibung gehe nur auf die Eigenschaften des Trainingsplans ein.
Nutze eine freundliche und motivierende Sprache. Grüße nicht den Schwimmer.
Für den Schwimmer ist nicht relevant, ob der Plan aus mehreren oder einem anderen Trainingsplan erstellt wurde.
Die Antwort soll in %s (Sprache) sein.

Die Antwort soll keine Fragen enthalten und auch nicht die Anweisung wiederholen.

Anfrage:
%s

Kontext:
| Belastungszone | Charakteristik                                                                                            | Dauer                                                                     | Intensität v%% akt. BZ                         | Laktat                                 | HF                                        | VO2max         | Pause                                                |
|----------------|-----------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------|------------------------------------------------|----------------------------------------|-------------------------------------------|----------------|------------------------------------------------------|
| BZ 1 (Rekom)   | - Zur Regeneration und Nachbereitung von Belastungen (Übergang zu BZ 2) - zur Lockerung - zum Laktatabbau | <70 %% (bei Lockerung ohne Bedeutung)                                     | <2 mmol/l (<LT1)                               | <120 oder <60%%max 50-70 unter HFmax   | 45-65 %%                                  | ohne Pause     |                                                      |
| BZ 2 (GA1ext)  | - extensive aerobe Ausdauer bei dominanter Fettverbrennung - Überdistanz-bereich                          | >60 min                                                                   | >70 %% (zumeist F/R)                           | 2-3 mmol/l (>LT1, <LT2)                | 120-50 oder 75-80%%max 40-50 unter HFmax  | 65-80 %%       | - kurze Trinkpause bei Intervall - Nach TS 10''-60'' |
| BZ 3 (GA1int)  | - intensive aerobe Ausdauer/Glykolyse - Schwimm-v bei 3 mmol/l                                            | 30-60 min                                                                 | ca. 75-80 %% (je nach Schwimmart/Streckenlänge) | 2,5-4 mmol/l (>LT1, <LT2)              | 140-180 oder 80-90%%max 30-40 unter HFmax  | 80-87 %%       | 10''-20'' (bis 60'' längere Strecken)                |
| BZ 4 (GA2Öko)  | - Aerob-anaerober Übergangsbereich - GA-Entwicklung - Intensive Ausdauer - Nahe Distanzbereich            | 20-45 min                                                                 | >85 %% (Schmett/Sprinter >80%%)                 | 4-6 mmol/l (>LT2)                      | 150-180 85-95%%max 20-30 unter HFmax      | 87-94 %%       | Je nach TS 30''-60''                                 |
| BZ 5 (GA2Entw) | - Aerob/anaerobe Leistungsfähigkeit - Nahe Distanzbereich - max VO2                                       | 10-30 min                                                                 | 85-95 %% je nach Schwimmart/Streckenlänge      | >6 mmol/l (bis über 10 mmol/l möglich) | 170-200 oder 90-100%%max 10-20 unter HFmax | 94-100 %%      | 60''-90''                                            |
| BZ 6 (WA)      | - Anaerobe Ausdauer - Wettkampf-spezifisch - Distanz voll oder gebrochen - Laktatmobilisation             | 3-10 (15) min, Wettkampf-zeit (eine TE wie WK mit Vor- u. Nach-bereitung) | >100 %% (Zielzeit) (je nach WK-Strecke)        | >8 mmol/l (LZA BZ 5)                   | maximal                                   | nicht relevant | 10''/15''/20'' bei Wdhlg. >400 Ko                    |
| BZ 7 (SA)      | - Anaerobe Ausdauer - Übergang von GA 2 zu WA - Unterdistanz - Wettkampfnah                               | 10-20 min (20''-120'' je TS)                                              | Unterdistanz 100-105 %%                        | >7 mmol/l                              | >180 95-100%%max 0-10 unter HFmax         | nicht relevant | 1-3 min                                              |
| BZ 8 (S)       | - Sprintschnelligkeit - Weitgehend alaktazid - Start / Wende                                              | <15 min                                                                   | 105-110 %% von v100m (bis 8 mmol/l möglich)    | nicht von Bedeutung                    | nicht relevant                            | bis 4' (aktiv) | vollständige Erholung                                |

Legende:
BZ = Belastungszone, GA = Grundlagenausdauer, HF = Herzfrequenz, LT = lactate threshold, LZA = Langzeitausdauer, Rekom = Regenerations- und Kompensationsbereich, SA = Schnelligkeitsausdauer, S = Schnelligkeit, WA = Wettkampfspezifische Ausdauer, WK = Wettkampf

Pläne:
%s
`

const choosePlanTemplateStr = `
Du bist ein Schwimmtrainer und hilfst einem Schwimmer einen Trainingsplan auszusuchen.
Du bekommst eine Frage vom Schwimmer und du hast eine Liste von Trainingsplänen als Kontext.
Wähle den besten Trainingsplan aus dem Kontext aus, der am besten zu der Frage und der gewünschten Beckenart, %s, passt.
Die Antwort soll in %s (Sprache) sein.
Die Antwort soll in JSON-Format sein.
{
	"description": "Eine kurze Beschreibung, Kommentare oder Anmerkungen zu dem Trainingsplan, damit der Schwimmer den Plan besser versteht",
	"index": "Der Index des Trainingsplans in der Liste als integer",
}

Die Antwort soll keine Fragen enthalten und auch nicht die Anweisung wiederholen.
Grüße nicht den Schwimmer, sondern beschreibe einfach den Trainingplan.

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

%v

Legende:
| Belastungszone | Charakteristik                                                                                            | Dauer                                                                     | Intensität v%% akt. BZ                         | Laktat                                 | HF                                        | VO2max         | Pause                                                |
|----------------|-----------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------|------------------------------------------------|----------------------------------------|-------------------------------------------|----------------|------------------------------------------------------|
| BZ 1 (Rekom)   | - Zur Regeneration und Nachbereitung von Belastungen (Übergang zu BZ 2) - zur Lockerung - zum Laktatabbau | <70 %% (bei Lockerung ohne Bedeutung)                                     | <2 mmol/l (<LT1)                               | <120 oder <60%%max 50-70 unter HFmax   | 45-65 %%                                  | ohne Pause     |                                                      |
| BZ 2 (GA1ext)  | - extensive aerobe Ausdauer bei dominanter Fettverbrennung - Überdistanz-bereich                          | >60 min                                                                   | >70 %% (zumeist F/R)                           | 2-3 mmol/l (>LT1, <LT2)                | 120-50 oder 75-80%%max 40-50 unter HFmax  | 65-80 %%       | - kurze Trinkpause bei Intervall - Nach TS 10''-60'' |
| BZ 3 (GA1int)  | - intensive aerobe Ausdauer/Glykolyse - Schwimm-v bei 3 mmol/l                                            | 30-60 min                                                                 | ca. 75-80 %% (je nach Schwimmart/Streckenlänge)| 2,5-4 mmol/l (>LT1, <LT2)              | 140-180 oder 80-90%%max 30-40 unter HFmax | 80-87 %%       | 10''-20'' (bis 60'' längere Strecken)                |
| BZ 4 (GA2Öko)  | - Aerob-anaerober Übergangsbereich - GA-Entwicklung - Intensive Ausdauer - Nahe Distanzbereich            | 20-45 min                                                                 | >85 %% (Schmett/Sprinter >80%%)                | 4-6 mmol/l (>LT2)                      | 150-180 85-95%%max 20-30 unter HFmax      | 87-94 %%       | Je nach TS 30''-60''                                 |
| BZ 5 (GA2Entw) | - Aerob/anaerobe Leistungsfähigkeit - Nahe Distanzbereich - max VO2                                       | 10-30 min                                                                 | 85-95 %% je nach Schwimmart/Streckenlänge      | >6 mmol/l (bis über 10 mmol/l möglich) | 170-200 oder 90-100%%max 10-20 unter HFmax| 94-100 %%      | 60''-90''                                            |
| BZ 6 (WA)      | - Anaerobe Ausdauer - Wettkampf-spezifisch - Distanz voll oder gebrochen - Laktatmobilisation             | 3-10 (15) min, Wettkampf-zeit (eine TE wie WK mit Vor- u. Nach-bereitung) | >100 %% (Zielzeit) (je nach WK-Strecke)        | >8 mmol/l (LZA BZ 5)                   | maximal                                   | nicht relevant | 10''/15''/20'' bei Wdhlg. >400 Ko                    |
| BZ 7 (SA)      | - Anaerobe Ausdauer - Übergang von GA 2 zu WA - Unterdistanz - Wettkampfnah                               | 10-20 min (20''-120'' je TS)                                              | Unterdistanz 100-105 %%                        | >7 mmol/l                              | >180 95-100%%max 0-10 unter HFmax         | nicht relevant | 1-3 min                                              |
| BZ 8 (S)       | - Sprintschnelligkeit - Weitgehend alaktazid - Start / Wende                                              | <15 min                                                                   | 105-110 %% von v100m (bis 8 mmol/l möglich)    | nicht von Bedeutung                    | nicht relevant                            | bis 4' (aktiv) | vollständige Erholung                                |


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
Die Antwort soll in Deutsch sein. Vermeide konkrete Zahlenangaben zur Gesamtdistanz.
Die Antwort soll in JSON-Format sein und in folgendem Schema:
%s

Kontext:
| Belastungszone | Charakteristik                                                                                            | Dauer                                                                     | Intensität v%% akt. BZ                         | Laktat                                 | HF                                        | VO2max         | Pause                                                |
|----------------|-----------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------|------------------------------------------------|----------------------------------------|-------------------------------------------|----------------|------------------------------------------------------|
| BZ 1 (Rekom)   | - Zur Regeneration und Nachbereitung von Belastungen (Übergang zu BZ 2) - zur Lockerung - zum Laktatabbau | <70 %% (bei Lockerung ohne Bedeutung)                                     | <2 mmol/l (<LT1)                               | <120 oder <60%%max 50-70 unter HFmax   | 45-65 %%                                  | ohne Pause     |                                                      |
| BZ 2 (GA1ext)  | - extensive aerobe Ausdauer bei dominanter Fettverbrennung - Überdistanz-bereich                          | >60 min                                                                   | >70 %% (zumeist F/R)                           | 2-3 mmol/l (>LT1, <LT2)                | 120-50 oder 75-80%%max 40-50 unter HFmax  | 65-80 %%       | - kurze Trinkpause bei Intervall - Nach TS 10''-60'' |
| BZ 3 (GA1int)  | - intensive aerobe Ausdauer/Glykolyse - Schwimm-v bei 3 mmol/l                                            | 30-60 min                                                                 | ca. 75-80 %% (je nach Schwimmart/Streckenlänge) | 2,5-4 mmol/l (>LT1, <LT2)              | 140-180 oder 80-90%%max 30-40 unter HFmax  | 80-87 %%       | 10''-20'' (bis 60'' längere Strecken)                |
| BZ 4 (GA2Öko)  | - Aerob-anaerober Übergangsbereich - GA-Entwicklung - Intensive Ausdauer - Nahe Distanzbereich            | 20-45 min                                                                 | >85 %% (Schmett/Sprinter >80%%)                 | 4-6 mmol/l (>LT2)                      | 150-180 85-95%%max 20-30 unter HFmax      | 87-94 %%       | Je nach TS 30''-60''                                 |
| BZ 5 (GA2Entw) | - Aerob/anaerobe Leistungsfähigkeit - Nahe Distanzbereich - max VO2                                       | 10-30 min                                                                 | 85-95 %% je nach Schwimmart/Streckenlänge      | >6 mmol/l (bis über 10 mmol/l möglich) | 170-200 oder 90-100%%max 10-20 unter HFmax | 94-100 %%      | 60''-90''                                            |
| BZ 6 (WA)      | - Anaerobe Ausdauer - Wettkampf-spezifisch - Distanz voll oder gebrochen - Laktatmobilisation             | 3-10 (15) min, Wettkampf-zeit (eine TE wie WK mit Vor- u. Nach-bereitung) | >100 %% (Zielzeit) (je nach WK-Strecke)        | >8 mmol/l (LZA BZ 5)                   | maximal                                   | nicht relevant | 10''/15''/20'' bei Wdhlg. >400 Ko                    |
| BZ 7 (SA)      | - Anaerobe Ausdauer - Übergang von GA 2 zu WA - Unterdistanz - Wettkampfnah                               | 10-20 min (20''-120'' je TS)                                              | Unterdistanz 100-105 %%                        | >7 mmol/l                              | >180 95-100%%max 0-10 unter HFmax         | nicht relevant | 1-3 min                                              |
| BZ 8 (S)       | - Sprintschnelligkeit - Weitgehend alaktazid - Start / Wende                                              | <15 min                                                                   | 105-110 %% von v100m (bis 8 mmol/l möglich)    | nicht von Bedeutung                    | nicht relevant                            | bis 4' (aktiv) | vollständige Erholung                                |

Legende:
BZ = Belastungszone, GA = Grundlagenausdauer, HF = Herzfrequenz, LT = lactate threshold, LZA = Langzeitausdauer, Rekom = Regenerations- und Kompensationsbereich, SA = Schnelligkeitsausdauer, S = Schnelligkeit, WA = Wettkampfspezifische Ausdauer, WK = Wettkampf

Tabelle:
%s

Antwort:
`

const generatePromptTemplateStr = `
Du bist ein Assistent für einen Schwimmer und der einen Trainingsplan von deinem Trainer erstellt bekommen möchte.
Du erstellst eine konkrete Anfrage für den Trainer, um einen Plan für ein einzelnes Training zu generieren.
Deine Antwort soll folgende Inhalte enthalten:
Ziele, Erfahrung, Zeitaufwand, ungefähre Gesamtdistanz, und Vorlieben.
Beginne die Anfrage mit "Erstelle einen Trainingplan mit ..." oder dem equivalenten in der jeweiligen Sprache.
Sei kreativ und halte dich kurz. Die Anfrage sollte nicht länger als 50 Wörter sein.
Die Antwort sollte im Fließtext sein und keine Formattierung enthalten.
Make sure to respond in %s (language code).
`

const translateTemplateStr = `
You are a professional translator. You are tasked with translating a training plan into a specified language.
The training plan is provided in JSON format and includes a title, a description, and a table of training data.
Your response must be in the same JSON format as the input.
The plan may contain abbreviations and specialized terms related to swimming training. Stay within the intent
of the original text and abbreviate, where appropriate, using common swimming terminology in the target language.

If the training plan is already in the target language, apply only minor adjustments if necessary and where appropriate.
Ensure that the structure of the training plan remains unchanged, including the table format.

Translate the following training plan into %s (language).

These abbreviations are relevant for the translation:
%v

Training Plan to Translate:

%s
%s

%s

Response:
`
