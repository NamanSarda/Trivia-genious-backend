package models

type Quiz struct {
	ID           int32         `json:"id" db:"id"`
	AuthorID     int32         `json:"author_id" db:"author_id"`
	Participants []Participant `json:"participants" db:"participants"`
	Questions    []Question    `json:"questions" db:"questions"`
	TimeLimit    int32         `json:"time_limit" db:"time_limit"`
	Tags         []string      `json:"tags" db:"tags"`
	Difficulty   string        `json:"difficulty" db:"difficulty"`
}
type Participant struct {
	ID     int32 `json:"id" db:"id"`
	QuizID int32 `json:"quiz_id" db:"quiz_id"`
	UserID int32 `json:"user_id" db:"user_id"`
	Score  int32 `json:"score" db:"score"`
}

type Question struct {
	ID          int32    `json:"id" db:"id"`
	QuizID      int32    `json:"quiz_id" db:"quiz_id"`
	Description string   `json:"description" db:"description"`
	Score       int32    `json:"score" db:"score"`
	Answer      string   `json:"answer" db:"answer"`
	Options     []string `json:"options" db:"options"`
}
