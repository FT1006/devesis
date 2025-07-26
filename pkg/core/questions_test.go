package core

import (
	"testing"
)

func TestGetRandomQuestion_ReturnsValidQuestion(t *testing.T) {
	state := newQuestionTestGameState()
	question, _ := GetRandomQuestion(state, 42)
	
	if question.Text == "" {
		t.Error("Question text should not be empty")
	}
	if len(question.Options) != 4 {
		t.Errorf("Question should have exactly 4 options, got %d", len(question.Options))
	}
	if question.CorrectAnswer < 0 || question.CorrectAnswer >= 4 {
		t.Error("Correct answer should be valid option index (0-3)")
	}
	if question.ID < 0 || question.ID >= 50 {
		t.Error("Question ID should be between 0-49")
	}
}

func TestGetRandomQuestion_MarksAsUsed(t *testing.T) {
	state := newQuestionTestGameState()
	question, newState := GetRandomQuestion(state, 42)
	
	// Question should be marked as used
	found := false
	for _, usedID := range newState.UsedQuestions {
		if usedID == question.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("Question should be marked as used")
	}
}

func TestGetRandomQuestion_NoRepeats(t *testing.T) {
	state := newQuestionTestGameState()
	usedQuestions := make(map[int]bool)
	
	// Ask multiple questions
	for i := 0; i < 10; i++ {
		question, newState := GetRandomQuestion(state, int64(i))
		
		if usedQuestions[question.ID] {
			t.Errorf("Question %d was repeated", question.ID)
		}
		usedQuestions[question.ID] = true
		state = newState
	}
}

func TestGetRandomQuestion_ExhaustsBank(t *testing.T) {
	state := newQuestionTestGameState()
	
	// Use all 50 questions
	for i := 0; i < 50; i++ {
		question, newState := GetRandomQuestion(state, int64(i))
		if question.ID == -1 {
			t.Errorf("Should not run out of questions at iteration %d", i)
		}
		state = newState
	}
	
	// 51st question should indicate no more questions
	question, _ := GetRandomQuestion(state, 99)
	if question.ID != -1 {
		t.Error("Should return empty question when bank is exhausted")
	}
}

func TestCheckAnswer_CorrectAnswer(t *testing.T) {
	state := newQuestionTestGameState()
	question, _ := GetRandomQuestion(state, 42)
	
	isCorrect := CheckAnswer(question, question.CorrectAnswer)
	
	if !isCorrect {
		t.Error("Correct answer should return true")
	}
}

func TestCheckAnswer_WrongAnswer(t *testing.T) {
	state := newQuestionTestGameState()
	question, _ := GetRandomQuestion(state, 42)
	wrongAnswer := (question.CorrectAnswer + 1) % 4
	
	isCorrect := CheckAnswer(question, wrongAnswer)
	
	if isCorrect {
		t.Error("Wrong answer should return false")
	}
}

// Helper for question tests
func newQuestionTestGameState() GameState {
	return GameState{
		UsedQuestions: []int{}, // Empty - no questions used yet
	}
}