package main

type Quiz struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Questions []Question `json:"questions"`
}

type Question struct {
	ID            string   `json:"id"`
	Text          string   `json:"text"`
	Options       []string `json:"options"`
	CorrectOption int      `json:"-"`
}

type Answer struct {
	QuestionID     string `json:"question_id"`
	SelectedOption int    `json:"selected_option"`
	IsCorrect      bool   `json:"is_correct"`
}

type Result struct {
	QuizID  string   `json:"quiz_id"`
	UserID  string   `json:"user_id"`
	Score   int      `json:"score"`
	Answers []Answer `json:"answers"`
}
