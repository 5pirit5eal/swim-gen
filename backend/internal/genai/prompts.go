package genai

const ragTemplateStr string = `
Du bist ein Schwimmtrainer und hilfst einem Schwimmer einen Trainingsplan zu erstellen.
Du bekommst eine Frage vom Schwimmer und du hast eine Liste von Trainingsplänen als Referenz.
Die Trainingspläne beinhalten Downloadlinks und nummerierte Titel. Mit diesen kann der Schwimmer nichts anfangen. Entsprechend sollten
diese nicht mit enthalten sein. Entferne deshalb "www.docswim.de" oder "EIN TRAININGSPLAN von DOC SWIM".
Erstelle dem Schwimmer einen passenden Trainingsplan basierend auf dem Kontext und seiner Anfrage. Dafür kannst du die Referenztrainingspläne verwenden,
indem du sie selektierst, kombinierst, mischst, oder umformulierst, um sie an die Bedürfnisse des Schwimmers anzupassen.
Ziehe dabei auch die konfigurierte Poollänge in Betracht: %s. Die Standard-Poollänge ist 25m.
Achte darauf, dass die Gesamtdistanz des Trainingsplans möglichst genau zu der Anfrage des Schwimmers passt!
Erhöhe die Anzahl oder die Distanz der einzelnen Wiederholungen,
oder entferne oder füge weitere Übungen hinzu um die Gesamtdistanz anzupassen.
Verwende zwingend das SubRows-Feld für Untereinheiten, z.B. "8 x (800 + 200)"
Die technischen Übungen dürfen nur als Referenzen eingefügt werden. Das Format ist ein Markdown URL Link.
Dafür wird der slug als Linktext verwendet und die URL als Linkziel. Exemplarisch: [slug](URL).
Erstelle nur Referenzen für technische Übungen, die mit /drill erreichbar sind. Andere sind nicht als Referenz geeignet.
Diese Übungungen sollen im Plan nicht näher beschrieben werden. Wichtig ist das die Intensität auf TÜ gesetzt wird.
Der Titel des Trainingsplans sollte humorvoll, einprägsam und kurz sein. Wortspiele sind willkommen. Füge nicht die angefragte Gesamtlänge in den Titel ein.
Bei der Erstellung der kurzen Beschreibung gehe nur auf die Eigenschaften des Trainingsplans ein.
Nutze eine freundliche und motivierende Sprache. Grüße nicht den Schwimmer.
Für den Schwimmer ist nicht relevant, ob der Plan aus mehreren oder einem anderen Trainingsplan erstellt wurde.
Die Antwort soll in %s (Sprache) sein.

%s

Die Antwort soll keine Fragen enthalten und auch nicht die Anweisung wiederholen.

Anfrage:
%s

Kontext:
| Belastungszone | Charakteristik                                                                                            | Dauer                                                                     | Intensität v%% akt. BZ                         | Laktat                                 | HF                                        | VO2max         | Pause                                                |
|----------------|-----------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------|------------------------------------------------|----------------------------------------|-------------------------------------------|----------------|------------------------------------------------------|
| BZ 1 (Rekom)   | - Zur Regeneration und Nachbereitung von Belastungen (Übergang zu BZ 2) - zur Lockerung - zum Laktatabbau | <70 %% (bei Lockerung ohne Bedeutung)                                     | <2 mmol/l (<LT1)                               | <120 oder <60%%max 50-70 unter HFmax   | 45-65 %%                                  | ohne Pause     |                                                      |
| BZ 2 (GA1)     | - extensive aerobe Ausdauer bei dominanter Fettverbrennung - Überdistanz-bereich                          | >60 min                                                                   | >70 %% (zumeist F/R)                           | 2-3 mmol/l (>LT1, <LT2)                | 120-50 oder 75-80%%max 40-50 unter HFmax  | 65-80 %%       | - kurze Trinkpause bei Intervall - Nach TS 10''-60'' |
| BZ 3 (GA1-2)   | - intensive aerobe Ausdauer/Glykolyse - Schwimm-v bei 3 mmol/l                                            | 30-60 min                                                                 | ca. 75-80 %% (je nach Schwimmart/Streckenlänge) | 2,5-4 mmol/l (>LT1, <LT2)              | 140-180 oder 80-90%%max 30-40 unter HFmax  | 80-87 %%       | 10''-20'' (bis 60'' längere Strecken)                |
| BZ 4 (GA2)     | - Aerob-anaerober Übergangsbereich - GA-Entwicklung - Intensive Ausdauer - Nahe Distanzbereich            | 20-45 min                                                                 | >85 %% (Schmett/Sprinter >80%%)                 | 4-6 mmol/l (>LT2)                      | 150-180 85-95%%max 20-30 unter HFmax      | 87-94 %%       | Je nach TS 30''-60''                                 |
| BZ 5 (GA2)     | - Aerob/anaerobe Leistungsfähigkeit - Nahe Distanzbereich - max VO2                                       | 10-30 min                                                                 | 85-95 %% je nach Schwimmart/Streckenlänge      | >6 mmol/l (bis über 10 mmol/l möglich) | 170-200 oder 90-100%%max 10-20 unter HFmax | 94-100 %%      | 60''-90''                                            |
| BZ 6 (WA)      | - Anaerobe Ausdauer - Wettkampf-spezifisch - Distanz voll oder gebrochen - Laktatmobilisation             | 3-10 (15) min, Wettkampf-zeit (eine TE wie WK mit Vor- u. Nach-bereitung) | >100 %% (Zielzeit) (je nach WK-Strecke)        | >8 mmol/l (LZA BZ 5)                   | maximal                                   | nicht relevant | 10''/15''/20'' bei Wdhlg. >400 Ko                    |
| BZ 7 (SA)      | - Anaerobe Ausdauer - Übergang von GA 2 zu WA - Unterdistanz - Wettkampfnah                               | 10-20 min (20''-120'' je TS)                                              | Unterdistanz 100-105 %%                        | >7 mmol/l                              | >180 95-100%%max 0-10 unter HFmax         | nicht relevant | 1-3 min                                              |
| BZ 8 (S)       | - Sprintschnelligkeit - Weitgehend alaktazid - Start / Wende                                              | <15 min                                                                   | 105-110 %% von v100m (bis 8 mmol/l möglich)    | nicht von Bedeutung                    | nicht relevant                            | bis 4' (aktiv) | vollständige Erholung                                |

Legende:
BZ = Belastungszone, GA = Grundlagenausdauer, HF = Herzfrequenz, LT = lactate threshold, LZA = Langzeitausdauer, Rekom = Regenerations- und Kompensationsbereich, SA = Schnelligkeitsausdauer, S = Schnelligkeit, WA = Wettkampfspezifische Ausdauer, WK = Wettkampf

Pläne:
%s

Technische Übungen:
%s
`

const choosePlanTemplateStr string = `
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

const metadataTemplateStr string = `
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
| BZ 2 (GA1)     | - extensive aerobe Ausdauer bei dominanter Fettverbrennung - Überdistanz-bereich                          | >60 min                                                                   | >70 %% (zumeist F/R)                           | 2-3 mmol/l (>LT1, <LT2)                | 120-50 oder 75-80%%max 40-50 unter HFmax  | 65-80 %%       | - kurze Trinkpause bei Intervall - Nach TS 10''-60'' |
| BZ 3 (GA1-2)   | - intensive aerobe Ausdauer/Glykolyse - Schwimm-v bei 3 mmol/l                                            | 30-60 min                                                                 | ca. 75-80 %% (je nach Schwimmart/Streckenlänge)| 2,5-4 mmol/l (>LT1, <LT2)              | 140-180 oder 80-90%%max 30-40 unter HFmax | 80-87 %%       | 10''-20'' (bis 60'' längere Strecken)                |
| BZ 4 (GA2)     | - Aerob-anaerober Übergangsbereich - GA-Entwicklung - Intensive Ausdauer - Nahe Distanzbereich            | 20-45 min                                                                 | >85 %% (Schmett/Sprinter >80%%)                | 4-6 mmol/l (>LT2)                      | 150-180 85-95%%max 20-30 unter HFmax      | 87-94 %%       | Je nach TS 30''-60''                                 |
| BZ 5 (GA2)     | - Aerob/anaerobe Leistungsfähigkeit - Nahe Distanzbereich - max VO2                                       | 10-30 min                                                                 | 85-95 %% je nach Schwimmart/Streckenlänge      | >6 mmol/l (bis über 10 mmol/l möglich) | 170-200 oder 90-100%%max 10-20 unter HFmax| 94-100 %%      | 60''-90''                                            |
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

const describeTemplateStr string = `
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
| BZ 2 (GA1)     | - extensive aerobe Ausdauer bei dominanter Fettverbrennung - Überdistanz-bereich                          | >60 min                                                                   | >70 %% (zumeist F/R)                           | 2-3 mmol/l (>LT1, <LT2)                | 120-50 oder 75-80%%max 40-50 unter HFmax  | 65-80 %%       | - kurze Trinkpause bei Intervall - Nach TS 10''-60'' |
| BZ 3 (GA1-2)   | - intensive aerobe Ausdauer/Glykolyse - Schwimm-v bei 3 mmol/l                                            | 30-60 min                                                                 | ca. 75-80 %% (je nach Schwimmart/Streckenlänge) | 2,5-4 mmol/l (>LT1, <LT2)              | 140-180 oder 80-90%%max 30-40 unter HFmax  | 80-87 %%       | 10''-20'' (bis 60'' längere Strecken)                |
| BZ 4 (GA2)     | - Aerob-anaerober Übergangsbereich - GA-Entwicklung - Intensive Ausdauer - Nahe Distanzbereich            | 20-45 min                                                                 | >85 %% (Schmett/Sprinter >80%%)                 | 4-6 mmol/l (>LT2)                      | 150-180 85-95%%max 20-30 unter HFmax      | 87-94 %%       | Je nach TS 30''-60''                                 |
| BZ 5 (GA2)     | - Aerob/anaerobe Leistungsfähigkeit - Nahe Distanzbereich - max VO2                                       | 10-30 min                                                                 | 85-95 %% je nach Schwimmart/Streckenlänge      | >6 mmol/l (bis über 10 mmol/l möglich) | 170-200 oder 90-100%%max 10-20 unter HFmax | 94-100 %%      | 60''-90''                                            |
| BZ 6 (WA)      | - Anaerobe Ausdauer - Wettkampf-spezifisch - Distanz voll oder gebrochen - Laktatmobilisation             | 3-10 (15) min, Wettkampf-zeit (eine TE wie WK mit Vor- u. Nach-bereitung) | >100 %% (Zielzeit) (je nach WK-Strecke)        | >8 mmol/l (LZA BZ 5)                   | maximal                                   | nicht relevant | 10''/15''/20'' bei Wdhlg. >400 Ko                    |
| BZ 7 (SA)      | - Anaerobe Ausdauer - Übergang von GA 2 zu WA - Unterdistanz - Wettkampfnah                               | 10-20 min (20''-120'' je TS)                                              | Unterdistanz 100-105 %%                        | >7 mmol/l                              | >180 95-100%%max 0-10 unter HFmax         | nicht relevant | 1-3 min                                              |
| BZ 8 (S)       | - Sprintschnelligkeit - Weitgehend alaktazid - Start / Wende                                              | <15 min                                                                   | 105-110 %% von v100m (bis 8 mmol/l möglich)    | nicht von Bedeutung                    | nicht relevant                            | bis 4' (aktiv) | vollständige Erholung                                |

Legende:
BZ = Belastungszone, GA = Grundlagenausdauer, HF = Herzfrequenz, LT = lactate threshold, LZA = Langzeitausdauer, Rekom = Regenerations- und Kompensationsbereich, SA = Schnelligkeitsausdauer, S = Schnelligkeit, WA = Wettkampfspezifische Ausdauer, WK = Wettkampf

Tabelle:
%s

Antwort:
`

const generatePromptTemplateStr string = `
Du bist ein Assistent für einen Schwimmer der einen Trainingsplan von deinem Trainer erstellt bekommen möchte.
Du erstellst eine konkrete Anfrage für den Trainer, um einen Plan für ein einzelnes Training zu generieren.
Deine Antwort soll folgende Inhalte enthalten:
Ziele, Erfahrung, Zeitaufwand, ungefähre Gesamtdistanz, und Vorlieben.
Beginne die Anfrage mit "Erstelle einen Trainingplan mit ..." oder dem equivalenten in der jeweiligen Sprache.
Sei kreativ und halte dich kurz. Deine Antwort sollte nicht länger als 3 Sätze sein.
Deine Antwort sollte im Fließtext sein und keine Formattierung enthalten.
Make sure to respond in %s (language code).
`

const translateTemplateStr string = `
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

const chatRefineTemplateStr string = `
Du bist ein Schwimmtrainer und hilfst einem Schwimmer, einen Trainingsplan in einer Unterhaltung zu erstellen oder zu verfeinern.

Du hast Zugriff auf:
1. Den bisherigen Gesprächsverlauf zwischen dir und dem Schwimmer
2. Den aktuellen Trainingsplan (falls bereits einer existiert)
3. Eine Liste von ähnlichen Trainingsplänen als Referenz
4. Die neueste Nachricht des Schwimmers

Deine Aufgabe ist es:
- Den Trainingsplan basierend auf dem Feedback des Schwimmers zu erstellen oder anzupassen
- Eine freundliche, motivierende Antwort zu geben, die die Änderungen erklärt
- Konsistenz mit vorherigen Entscheidungen im Gespräch zu wahren
- Auf spezifische Fragen des Schwimmers zum Plan einzugehen

Der Schwimmer möchte für eine Beckenlänge von %v trainieren.
Die Antwort soll in %s (Sprache) sein.

GESPRÄCHSVERLAUF:
%s

AKTUELLER PLAN:
%s

REFERENZPLÄNE (zur Inspiration):
%s

NEUE NACHRICHT VOM SCHWIMMER:
%s

Bitte gib deine Antwort in folgendem JSON-Format zurück:
- "plan": Der aktualisierte oder neue Trainingsplan mit title, description und table
  - Wenn der Schwimmer nur eine Frage stellt ohne Änderungswunsch, gib den bestehenden Plan zurück (oder null falls keiner existiert)
  - Wenn der Schwimmer Änderungen wünscht, passe den Plan entsprechend an
- "response": Eine freundliche, konversationelle Antwort die:
  - Erklärt, was du gemacht hast oder beantwortet die Frage des Schwimmers
  - Auf spezifische Aspekte des Plans eingeht
  - Motivierend und hilfreich ist
  - KEINE Grußformeln enthält (kein "Hallo", kein "Viel Erfolg", etc.)

Antwort:
`

const ocrTemplateStr string = `
Analysiere diese Datei und extrahiere den darin enthaltenen Plan möglichst genau.
Falls das Schema für den Trainingsplan nicht genau passt, modifiziere den Plan entsprechend
und passe ihn an das Schema an. Gib das Ergebnis im JSON-Format zurück.

Wenn wiederholende Blöcke erkannt werden (z.B. "2 x (100m Kraul, 200m Brust, 100m Locker)"), verwende das SubRows-Feld

Antworte in der Sprache: %s.

Antwort:
`

const restructureTemplateStr string = `
Du bist ein Schwimmtrainer-Experte. Deine Aufgabe ist es, Trainingspläne in ein optimiertes Format zu überführen.

Der folgende Plan wurde aus einer flach strukturierten Quelle extrahiert. Analysiere den Plan und strukturiere ihn neu.
Nutze dazu das JSON schema. WENN verkapselte Blöcke erkannt werden, verwende IMMER das SubRows-Feld für Untereinheiten.
Passe die Content Spalte entsprechend an, um die Aufteilung in SubRows zu reflektieren. Behalte die ursprüngliche Bedeutung bei
damit der Trainingsplan inhaltlich gleich bleibt, aber strukturell verbessert wird.

VORMARKIERTE SUBROW-KANDIDATEN:
Manche Zeilen in der Tabelle sind mit dem Hinweis [→ SUBROW-KANDIDAT: '+' Zeichen gefunden] markiert.
Diese Zeilen MÜSSEN zwingend in SubRows aufgeteilt werden (z.B. "100m Kraul + 200m Brust" = zwei SubRows).
Trenne diese Zeile anhand der '+'-Trennstellen in separate SubRows.
Manche Zeilen in der Tabelle sind mit dem Hinweis [→ SUBROW-KANDIDAT: ' x ' Zeichen gefunden] markiert.
Prüfe kritisch, ob diese Zeilen in SubRows aufgeteilt werden müssen, oder das Aufspalten in mehrere Rows sinnvoller ist.
das eine wiederholende Übungseinheit signalisiert (z.B. "2 x (100m Kraul, 200m Brust)" = zwei SubRows).
Trenne diese Zeile anhand der ' x ' Trennstellen in separate SubRows auf.

Die Vormarkierung ist nicht abschließend. Es können weitere Zeilen SubRows oder Aufspaltung benötigen, auch wenn sie
nicht markiert sind - zum Beispiel bei nummerierten Sequenzen, Klammern, oder anderen Mustern, die
wiederholende Untereinheiten beschreiben. Auch sind tiefere Ebenen möglich. Prüfe den gesamten Plan auf solche Muster.

Beispiele:
- "Amount": 8, "Distance": 75, "Content": "1.-4. je pro 25m das Tempo erhöhen 5.-8. je 25m Spurt + 50m easy going", ...
  sollte in 2 Rows mit jeweils "Amount:" 4, "Distance": 75, und 2 SubRows umgewandelt werden, da es sich um eine wiederholende Sequenz handelt.
  Die zweite Row muss wiederum in 2 SubRows aufgeteilt werden, da sie zwei Untereinheiten enthält.
  Die neuen Rows, welche die bisherige ersetzen, würden dann ungefähr so aussehen:
  {"Amount": 4, "Distance": 75, "Content": "pro 25m das Tempo erhöhen", SubRows: [], ...},
  {"Amount": 4, "Distance": 75, "Content": "25m Spurt", SubRows: [
    {"Amount": 1, "Distance": 25, "Content": "Spurt", ...},
    {"Amount": 1, "Distance": 50, "Content": "easy going", ...},
  ], ...}

- "Amount": 5 "Distance": , "Content": "1. 100m Kraul, 2. 200m Brust, 3. 100m Locker", ...
  Dann umwandeln zu:
  {"Amount": 5, "Distance": 0, "Content": "Gemischte Serie", SubRows: [
    {"Amount": 1, "Distance": 100, "Content": "Kraul", ...},
    {"Amount": 1, "Distance": 200, "Content": "Brust", ...},
    {"Amount": 1, "Distance": 100, "Content": "Locker", ...},
  ], ...}


EQUIPMENT ERKENNUNG:
- Analysiere den Content und die Beschreibung auf Hinweise zu benötigter Ausrüstung. Das im Input möglicherweise bestehende Equipmentfeld ist immer leer.
- Verwende EXAKT diese Ausrüstungswerte (in Deutsch): Flossen, Kickboard, Handpaddles, Pull buoy, Schnorchel
- Beispiele für Equipment-Erkennung:
  * "Kraul-Beine" oder "Kick" → Kickboard
  * "Delfinbeine mit Flossen" oder "Kraul mit Flossen" → Flossen
  * "Paddles", "Fingerpaddles" oder "Handpaddles" → Handpaddles
  * "Pull buoy", "Arme mit Brett" oder "Kraularme" → Pull buoy
  * "Schnorchel" → Schnorchel
- Falls ein anderes Equipment notwendig ist, schreibe es in die Content Spalte
- "Brett" in der Nutzung mit Beine ist ein "Kickboard"
- "Brett" in der Nutzung mit Arme ist ein "Pull buoy"
- Parent rows (mit SubRows) sollten kein Equipment haben, es sei denn es gilt für alle Untereinheiten
- SubRows können auch mehrere Ausrüstungsgegenstände haben, wenn die z.B. Arme mit Paddles geschwommen wird --> "Equipment": "Pull buoy, Handpaddles"

**Antworte in folgendem Schema:**
%s

**Aktueller Plan:**
Titel: %s
Beschreibung: %s
Tabelle:
%s

**Antwort:**
`
