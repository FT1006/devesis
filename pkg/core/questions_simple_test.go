package core

import (
	"testing"
)

func TestGetRandomQuestion_ReturnsValidQuestion(t *testing.T) {
	state := initializeGameState(42, Frontend)
	question, _ := GetRandomQuestion(state)
	
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
		t.Errorf("Question ID should be 0-49, got %d", question.ID)
	}
}

func TestGetRandomQuestion_AdvancesCounter(t *testing.T) {
	state := initializeGameState(42, Frontend)
	
	// Initial state
	if state.NextQuestion != 0 {
		t.Errorf("Initial NextQuestion should be 0, got %d", state.NextQuestion)
	}
	
	// First question
	question1, newState1 := GetRandomQuestion(state)
	if newState1.NextQuestion != 1 {
		t.Errorf("After first question, NextQuestion should be 1, got %d", newState1.NextQuestion)
	}
	
	// Second question
	question2, newState2 := GetRandomQuestion(newState1)
	if newState2.NextQuestion != 2 {
		t.Errorf("After second question, NextQuestion should be 2, got %d", newState2.NextQuestion)
	}
	
	// Questions should be different
	if question1.ID == question2.ID {
		t.Error("Sequential questions should be different")
	}
}

func TestCheckAnswer_CorrectAnswer(t *testing.T) {
	state := initializeGameState(42, Frontend)
	question, _ := GetRandomQuestion(state)
	
	// Test correct answer
	if !CheckAnswer(question, question.CorrectAnswer) {
		t.Error("CheckAnswer should return true for correct answer")
	}
}

func TestCheckAnswer_WrongAnswer(t *testing.T) {
	state := initializeGameState(42, Frontend)
	question, _ := GetRandomQuestion(state)
	
	// Test wrong answer (assuming correct answer is not 3)
	wrongAnswer := (question.CorrectAnswer + 1) % 4
	if CheckAnswer(question, wrongAnswer) {
		t.Error("CheckAnswer should return false for wrong answer")
	}
}