package models

type UserProfile struct {
	UserID            string   `db:"user_id"`
	UpdatedAt         string   `db:"updated_at"`
	Username          string   `db:"username"`
	Experience        *string  `db:"experience,omitempty"`
	PreferredLanguage *string  `db:"preferred_language,omitempty"`
	PreferredStrokes  []string `db:"preferred_strokes"`
	Categories        []string `db:"categories"`
}

type Feedback struct {
	UserID    string `db:"user_id"`
	PlanID    string `db:"plan_id"`
	Rating    int    `db:"rating"`
	Comment   string `db:"comment"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

type ChoiceResult struct {
	Idx         int    `json:"index" example:"1"`
	Description string `json:"description" example:"Selected plan based on your requirements"`
}
