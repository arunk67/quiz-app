package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generate unique ID for new quizzes or questions
func generateID() string {
	return strconv.Itoa(rand.Intn(100000))
}

// Handler to create a new quiz
func createQuiz(w http.ResponseWriter, r *http.Request) {
	var quiz Quiz
	if err := json.NewDecoder(r.Body).Decode(&quiz); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	quiz.ID = generateID()
	for i := range quiz.Questions {
		quiz.Questions[i].ID = generateID()
	}

	quizMutex.Lock()
	quizzes[quiz.ID] = quiz
	quizMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(quiz)
}

// Handler to get a quiz by ID (without revealing correct answers)
func getQuiz(w http.ResponseWriter, r *http.Request) {
	quizID := r.URL.Query().Get("id")
	quizMutex.Lock()
	quiz, found := quizzes[quizID]
	quizMutex.Unlock()

	if !found {
		http.Error(w, "Quiz not found", http.StatusNotFound)
		return
	}

	// Remove correct answers before sending response
	sanitizedQuiz := quiz
	for i := range sanitizedQuiz.Questions {
		sanitizedQuiz.Questions[i].CorrectOption = -1
	}

	json.NewEncoder(w).Encode(sanitizedQuiz)
}

// Handler to submit an answer to a question
func submitAnswer(w http.ResponseWriter, r *http.Request) {
	var answer Answer
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	quizID := r.URL.Query().Get("quiz_id")
	quizMutex.Lock()
	quiz, found := quizzes[quizID]
	quizMutex.Unlock()

	if !found {
		http.Error(w, "Quiz not found", http.StatusNotFound)
		return
	}

	var correct bool
	for _, question := range quiz.Questions {
		if question.ID == answer.QuestionID {
			correct = answer.SelectedOption == question.CorrectOption
			break
		}
	}

	answer.IsCorrect = correct
	json.NewEncoder(w).Encode(answer)
}

// Handler to get quiz results
func getResults(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	quizID := r.URL.Query().Get("quiz_id")
	resultID := quizID + "_" + userID

	quizMutex.Lock()
	result, found := results[resultID]
	quizMutex.Unlock()

	if !found {
		http.Error(w, "Results not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/quiz/create", createQuiz)
	http.HandleFunc("/quiz", getQuiz)
	http.HandleFunc("/quiz/answer", submitAnswer)
	http.HandleFunc("/quiz/results", getResults)

	http.ListenAndServe(":8080", nil)
}
