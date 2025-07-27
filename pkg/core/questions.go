package core


// GetRandomQuestion returns the next question in the pre-shuffled order
func GetRandomQuestion(state GameState) (Question, GameState) {
	// Create a copy of the state to avoid mutations
	newState := deepCopyGameState(state)
	
	// Check if we've exhausted all questions
	if newState.NextQuestion >= len(newState.QuestionOrder) {
		return Question{ID: -1}, newState
	}
	
	// Get the next question ID from pre-shuffled order
	questionID := newState.QuestionOrder[newState.NextQuestion]
	newState.NextQuestion++
	
	// Return the question
	question := getQuestionBank()[questionID]
	return question, newState
}

// CheckAnswer checks if the provided answer index is correct
func CheckAnswer(question Question, answerIndex int) bool {
	return question.CorrectAnswer == answerIndex
}

// getQuestionBank returns the hardcoded question bank
func getQuestionBank() [50]Question {
	return [50]Question{
		{ID: 0, Text: "Which keyword is used to declare a variable in Go?", Options: [4]string{"var", "let", "const", "declare"}, CorrectAnswer: 0},
		{ID: 1, Text: "What is the zero value of an int in Go?", Options: [4]string{"nil", "0", "1", "-1"}, CorrectAnswer: 1},
		{ID: 2, Text: "Which of these is a valid Go slice declaration?", Options: [4]string{"[]int{1,2,3}", "int[]{1,2,3}", "slice<int>", "Array[int]"}, CorrectAnswer: 0},
		{ID: 3, Text: "What does 'go' keyword do in Go?", Options: [4]string{"Import package", "Start goroutine", "Create variable", "Define function"}, CorrectAnswer: 1},
		{ID: 4, Text: "Which is the correct way to define a function in Go?", Options: [4]string{"function myFunc() {}", "func myFunc() {}", "def myFunc():", "fn myFunc() {}"}, CorrectAnswer: 1},
		{ID: 5, Text: "What is the time complexity of binary search?", Options: [4]string{"O(n)", "O(log n)", "O(n²)", "O(1)"}, CorrectAnswer: 1},
		{ID: 6, Text: "Which data structure uses LIFO principle?", Options: [4]string{"Queue", "Stack", "Array", "Linked List"}, CorrectAnswer: 1},
		{ID: 7, Text: "What does API stand for?", Options: [4]string{"Application Programming Interface", "Advanced Program Integration", "Automated Process Implementation", "Application Process Interface"}, CorrectAnswer: 0},
		{ID: 8, Text: "Which HTTP method is idempotent?", Options: [4]string{"POST", "GET", "PATCH", "All of the above"}, CorrectAnswer: 1},
		{ID: 9, Text: "What is the output of: fmt.Println(5 / 2) in Go?", Options: [4]string{"2.5", "2", "3", "Error"}, CorrectAnswer: 1},
		{ID: 10, Text: "Which is NOT a primitive data type in Go?", Options: [4]string{"int", "string", "bool", "array"}, CorrectAnswer: 3},
		{ID: 11, Text: "What is a deadlock?", Options: [4]string{"Infinite loop", "Memory leak", "Circular wait condition", "Stack overflow"}, CorrectAnswer: 2},
		{ID: 12, Text: "Which sorting algorithm has O(n log n) average time complexity?", Options: [4]string{"Bubble sort", "Quick sort", "Selection sort", "Insertion sort"}, CorrectAnswer: 1},
		{ID: 13, Text: "What is polymorphism?", Options: [4]string{"Multiple inheritance", "Function overloading", "One interface, multiple implementations", "Code reuse"}, CorrectAnswer: 2},
		{ID: 14, Text: "Which is a NoSQL database?", Options: [4]string{"MySQL", "PostgreSQL", "MongoDB", "SQLite"}, CorrectAnswer: 2},
		{ID: 15, Text: "What does SQL stand for?", Options: [4]string{"Structured Query Language", "Simple Query Language", "Standard Query Language", "System Query Language"}, CorrectAnswer: 0},
		{ID: 16, Text: "Which is the correct way to create a channel in Go?", Options: [4]string{"make(chan int)", "chan int{}", "new(chan int)", "channel<int>"}, CorrectAnswer: 0},
		{ID: 17, Text: "What is the Big O notation for accessing an element in an array?", Options: [4]string{"O(n)", "O(log n)", "O(1)", "O(n²)"}, CorrectAnswer: 2},
		{ID: 18, Text: "Which design pattern ensures only one instance of a class?", Options: [4]string{"Factory", "Observer", "Singleton", "Strategy"}, CorrectAnswer: 2},
		{ID: 19, Text: "What is the purpose of a constructor?", Options: [4]string{"Destroy objects", "Initialize objects", "Copy objects", "Compare objects"}, CorrectAnswer: 1},
		{ID: 20, Text: "Which is NOT a version control system?", Options: [4]string{"Git", "SVN", "Docker", "Mercurial"}, CorrectAnswer: 2},
		{ID: 21, Text: "What does REST stand for?", Options: [4]string{"Representational State Transfer", "Remote State Transfer", "Relational State Transfer", "Resource State Transfer"}, CorrectAnswer: 0},
		{ID: 22, Text: "Which HTTP status code indicates 'Not Found'?", Options: [4]string{"200", "404", "500", "301"}, CorrectAnswer: 1},
		{ID: 23, Text: "What is the purpose of a hash table?", Options: [4]string{"Sort data", "Store key-value pairs", "Implement recursion", "Handle concurrency"}, CorrectAnswer: 1},
		{ID: 24, Text: "Which is a characteristic of functional programming?", Options: [4]string{"Mutable state", "Side effects", "Pure functions", "Global variables"}, CorrectAnswer: 2},
		{ID: 25, Text: "What is the difference between compile-time and runtime?", Options: [4]string{"No difference", "Compile-time is before execution", "Runtime is before compilation", "Both happen simultaneously"}, CorrectAnswer: 1},
		{ID: 26, Text: "Which data structure is best for implementing BFS?", Options: [4]string{"Stack", "Queue", "Array", "Tree"}, CorrectAnswer: 1},
		{ID: 27, Text: "What is encapsulation in OOP?", Options: [4]string{"Data hiding", "Multiple inheritance", "Function overloading", "Memory management"}, CorrectAnswer: 0},
		{ID: 28, Text: "Which is NOT a JavaScript data type?", Options: [4]string{"undefined", "number", "character", "boolean"}, CorrectAnswer: 2},
		{ID: 29, Text: "What is the purpose of unit testing?", Options: [4]string{"Test entire system", "Test individual components", "Test user interface", "Test database"}, CorrectAnswer: 1},
		{ID: 30, Text: "Which algorithm is used for finding shortest path?", Options: [4]string{"Binary search", "Merge sort", "Dijkstra's algorithm", "Quick sort"}, CorrectAnswer: 2},
		{ID: 31, Text: "What is a race condition?", Options: [4]string{"Fast algorithm", "Concurrent access issue", "Memory overflow", "Infinite recursion"}, CorrectAnswer: 1},
		{ID: 32, Text: "Which is a characteristic of agile development?", Options: [4]string{"Waterfall model", "Iterative development", "No documentation", "Fixed requirements"}, CorrectAnswer: 1},
		{ID: 33, Text: "What does CPU stand for?", Options: [4]string{"Central Processing Unit", "Computer Processing Unit", "Central Program Unit", "Computer Program Unit"}, CorrectAnswer: 0},
		{ID: 34, Text: "Which is NOT a programming paradigm?", Options: [4]string{"Object-oriented", "Functional", "Procedural", "Debugging"}, CorrectAnswer: 3},
		{ID: 35, Text: "What is the purpose of a firewall?", Options: [4]string{"Speed up internet", "Block unauthorized access", "Store data", "Compile code"}, CorrectAnswer: 1},
		{ID: 36, Text: "Which is a valid Go interface declaration?", Options: [4]string{"interface Reader", "type Reader interface", "interface{} Reader", "Reader interface{}"}, CorrectAnswer: 1},
		{ID: 37, Text: "What is the time complexity of inserting at the end of a dynamic array?", Options: [4]string{"O(1) amortized", "O(n)", "O(log n)", "O(n²)"}, CorrectAnswer: 0},
		{ID: 38, Text: "Which is NOT a relational database?", Options: [4]string{"MySQL", "Redis", "PostgreSQL", "Oracle"}, CorrectAnswer: 1},
		{ID: 39, Text: "What is the purpose of middleware in web development?", Options: [4]string{"Store data", "Handle requests between client and server", "Compile code", "Design UI"}, CorrectAnswer: 1},
		{ID: 40, Text: "Which is a mutable data structure in most languages?", Options: [4]string{"String", "Array", "Integer", "Boolean"}, CorrectAnswer: 1},
		{ID: 41, Text: "What does DRY principle stand for?", Options: [4]string{"Don't Repeat Yourself", "Do Repeat Yourself", "Don't Run Yet", "Dynamic Resource Yielding"}, CorrectAnswer: 0},
		{ID: 42, Text: "Which is the correct way to handle errors in Go?", Options: [4]string{"try-catch", "if err != nil", "throw exception", "error.handle()"}, CorrectAnswer: 1},
		{ID: 43, Text: "What is the purpose of a load balancer?", Options: [4]string{"Store data", "Distribute incoming requests", "Compile code", "Encrypt data"}, CorrectAnswer: 1},
		{ID: 44, Text: "Which is NOT a software testing type?", Options: [4]string{"Unit testing", "Integration testing", "Compilation testing", "System testing"}, CorrectAnswer: 2},
		{ID: 45, Text: "What is the difference between HTTP and HTTPS?", Options: [4]string{"No difference", "HTTPS is encrypted", "HTTP is faster", "HTTPS is older"}, CorrectAnswer: 1},
		{ID: 46, Text: "Which data structure uses FIFO principle?", Options: [4]string{"Stack", "Queue", "Tree", "Graph"}, CorrectAnswer: 1},
		{ID: 47, Text: "What is the purpose of version control?", Options: [4]string{"Speed up code", "Track code changes", "Compile code", "Design UI"}, CorrectAnswer: 1},
		{ID: 48, Text: "Which is a characteristic of microservices architecture?", Options: [4]string{"Single large application", "Loosely coupled services", "Shared database", "Monolithic deployment"}, CorrectAnswer: 1},
		{ID: 49, Text: "What is the time complexity of linear search?", Options: [4]string{"O(1)", "O(log n)", "O(n)", "O(n²)"}, CorrectAnswer: 2},
	}
}