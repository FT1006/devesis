package core

import (
	"testing"
)

func TestQuestionExhaustion_50Questions(t *testing.T) {
	// Initialize game state with pre-shuffled questions
	state := initializeGameState(42, Frontend)
	
	// Verify we start with 50 questions available
	if len(state.QuestionOrder) != 50 {
		t.Fatalf("Expected 50 questions in order, got %d", len(state.QuestionOrder))
	}
	
	if state.NextQuestion != 0 {
		t.Fatalf("Expected NextQuestion to start at 0, got %d", state.NextQuestion)
	}
	
	// Ask 50 questions - should all succeed
	usedQuestionIDs := make(map[int]bool)
	currentState := state
	
	for i := 0; i < 50; i++ {
		question, newState := GetRandomQuestion(currentState)
		
		if question.ID == -1 {
			t.Fatalf("Question %d returned exhausted signal (ID: -1)", i+1)
		}
		
		// Verify question ID is valid (0-49)
		if question.ID < 0 || question.ID >= 50 {
			t.Fatalf("Question %d has invalid ID: %d", i+1, question.ID)
		}
		
		// Verify question ID hasn't been used before
		if usedQuestionIDs[question.ID] {
			t.Fatalf("Question %d returned duplicate ID: %d", i+1, question.ID)
		}
		usedQuestionIDs[question.ID] = true
		
		// Verify NextQuestion counter incremented
		if newState.NextQuestion != i+1 {
			t.Fatalf("Question %d: expected NextQuestion = %d, got %d", i+1, i+1, newState.NextQuestion)
		}
		
		currentState = newState
	}
	
	// Verify we've used exactly 50 unique questions
	if len(usedQuestionIDs) != 50 {
		t.Fatalf("Expected 50 unique questions, got %d", len(usedQuestionIDs))
	}
}

func TestQuestionExhaustion_51stQuestion(t *testing.T) {
	// Initialize game state with pre-shuffled questions
	state := initializeGameState(42, Frontend)
	currentState := state
	
	// Ask 50 questions (exhaust the pool)
	for i := 0; i < 50; i++ {
		question, newState := GetRandomQuestion(currentState)
		
		if question.ID == -1 {
			t.Fatalf("Question %d unexpectedly returned exhausted signal", i+1)
		}
		
		currentState = newState
	}
	
	// Verify we're at the end
	if currentState.NextQuestion != 50 {
		t.Fatalf("After 50 questions, expected NextQuestion = 50, got %d", currentState.NextQuestion)
	}
	
	// Now ask the 51st question - should return exhausted signal
	question, finalState := GetRandomQuestion(currentState)
	
	if question.ID != -1 {
		t.Fatalf("51st question should return exhausted signal (ID: -1), got ID: %d", question.ID)
	}
	
	// State should remain unchanged when exhausted
	if finalState.NextQuestion != currentState.NextQuestion {
		t.Fatalf("Exhausted state should not change NextQuestion, was %d, became %d", 
			currentState.NextQuestion, finalState.NextQuestion)
	}
}

func TestQuestionOrder_Deterministic(t *testing.T) {
	// Same seed should produce same question order
	state1 := initializeGameState(42, Frontend)
	state2 := initializeGameState(42, Frontend)
	
	if len(state1.QuestionOrder) != len(state2.QuestionOrder) {
		t.Fatalf("Same seed produced different question order lengths")
	}
	
	for i := 0; i < len(state1.QuestionOrder); i++ {
		if state1.QuestionOrder[i] != state2.QuestionOrder[i] {
			t.Fatalf("Same seed produced different question order at index %d: %d vs %d", 
				i, state1.QuestionOrder[i], state2.QuestionOrder[i])
		}
	}
}

func TestQuestionOrder_DifferentSeeds(t *testing.T) {
	// Different seeds should produce different question orders
	state1 := initializeGameState(42, Frontend)
	state2 := initializeGameState(100, Frontend)
	
	// Check that at least some positions are different
	differences := 0
	for i := 0; i < len(state1.QuestionOrder); i++ {
		if state1.QuestionOrder[i] != state2.QuestionOrder[i] {
			differences++
		}
	}
	
	if differences == 0 {
		t.Fatalf("Different seeds produced identical question orders")
	}
	
	// Should have significant differences (expect at least 20% different)
	if differences < 10 {
		t.Fatalf("Different seeds produced too similar orders, only %d differences out of 50", differences)
	}
}