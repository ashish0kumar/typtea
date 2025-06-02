package game

import (
	"strings"
	"time"
	"fmt"
	"github.com/ashish0kumar/typtea/internal/database"
)

// TypingStats holds the statistics for a game session
type TypingStats struct {
	WPM               float64
	BestWPM			  float64
	Accuracy          float64
	BestAccuracy      float64
	CharactersTyped   int
	CorrectChars      int
	TotalChars        int
	TimeElapsed       time.Duration
	IsComplete        bool
	UncorrectedErrors int
}

// TypingGame represents the state of a game session
type TypingGame struct {
	AllWords        []string
	DisplayLines    []string
	UserInput       string
	CurrentPos      int
	GlobalPos       int
	StartTime       time.Time
	Duration        int
	IsStarted       bool
	IsFinished      bool
	Errors          map[int]bool
	TotalErrorsMade int
	LinesPerView    int
	CharsPerLine    int
	WordsTyped      int
}

// NewTypingGame initializes a new TypingGame instance with a specified duration
func NewTypingGame(duration int) *TypingGame {
	game := &TypingGame{
		AllWords:     GenerateWords(200),
		Duration:     duration,
		Errors:       make(map[int]bool),
		LinesPerView: 3,
		CharsPerLine: 50,
	}
	game.generateDisplayLines()
	return game
}

// Reset reinitializes the game to a fresh state
func (g *TypingGame) Reset() {
	*g = *NewTypingGame(g.Duration)
}

// generateDisplayLines creates the initial display lines based on the words available
func (g *TypingGame) generateDisplayLines() {
	lines := make([]string, 0, g.LinesPerView)
	wordIndex := g.WordsTyped

	// Generate exactly g.LinesPerView lines
	for lineNum := 0; lineNum < g.LinesPerView && wordIndex < len(g.AllWords); lineNum++ {
		var currentLine strings.Builder

		// Fill current line with words
		for wordIndex < len(g.AllWords) {
			word := g.AllWords[wordIndex]
			spaceNeeded := 0
			if currentLine.Len() > 0 {
				spaceNeeded = 1
			}

			// Check if word fits
			if currentLine.Len()+spaceNeeded+len(word) <= g.CharsPerLine {
				if currentLine.Len() > 0 {
					currentLine.WriteString(" ")
				}
				currentLine.WriteString(word)
				wordIndex++
			} else {
				// Word doesn't fit, break to next line
				break
			}
		}

		// Add the completed line
		if currentLine.Len() > 0 {
			lines = append(lines, currentLine.String())
		} else {
			// If no words fit, add empty line
			lines = append(lines, "")
		}
	}

	// Ensure we have exactly g.LinesPerView lines
	for len(lines) < g.LinesPerView {
		lines = append(lines, "")
	}

	// Truncate if somehow we have more than g.LinesPerView lines
	if len(lines) > g.LinesPerView {
		lines = lines[:g.LinesPerView]
	}

	g.DisplayLines = lines
}

// Start initializes the game session if it hasn't started yet
func (g *TypingGame) Start() {
	if !g.IsStarted {
		g.StartTime = time.Now()
		g.IsStarted = true
	}
}

// AddCharacter handles user input and updates game state
func (g *TypingGame) AddCharacter(char rune) {
	if !g.IsStarted {
		g.Start()
	}

	if g.IsFinished || g.IsTimeUp() {
		g.IsFinished = true
		return
	}

	g.UserInput += string(char)
	displayText := []rune(strings.Join(g.DisplayLines, " "))

	// Check if the character is within bounds
	if g.CurrentPos < len(displayText) && g.CurrentPos >= 0 {
		if displayText[g.CurrentPos] != char {
			g.Errors[g.GlobalPos] = true
			g.TotalErrorsMade++
		}
		g.CurrentPos++
		g.GlobalPos++

		// Check if the first line is completed
		if g.CurrentPos >= len([]rune(g.DisplayLines[0])) {
			g.shiftLines()
		}
	} else {
		// End of available text
		g.IsFinished = true
	}
}

// shiftLines moves to the next line in the game, updating the words typed and generating new lines
func (g *TypingGame) shiftLines() {
	// Move to next line
	g.WordsTyped += len(strings.Fields(g.DisplayLines[0]))
	g.CurrentPos = 0

	// Generate new lines
	g.generateDisplayLines()

	// Extend words if needed
	if g.WordsTyped > len(g.AllWords)-50 {
		newWords := GenerateWords(100)
		g.AllWords = append(g.AllWords, newWords...)
	}
}

// RemoveCharacter removes the last character from the user input and updates the position
func (g *TypingGame) RemoveCharacter() {
	if len(g.UserInput) > 0 && g.CurrentPos > 0 {
		g.UserInput = g.UserInput[:len(g.UserInput)-1]
		g.CurrentPos--
		g.GlobalPos--

		// Remove error mark if previously added
		delete(g.Errors, g.GlobalPos)
	}
}

// GetDisplayText returns the current text to be displayed in the game
func (g *TypingGame) GetDisplayText() string {
	return strings.Join(g.DisplayLines, " ")
}

// GetStats calculates and returns the typing statistics for the current game session
func (g *TypingGame) GetStats() TypingStats {
	if !g.IsStarted {
		return TypingStats{}
	}

	elapsed := time.Since(g.StartTime)
	minutes := elapsed.Minutes()

	// Calculate Gross WPM (all typed entries / 5 / time in minutes)
	grossWPM := 0.0
	if minutes > 0 {
		grossWPM = float64(g.GlobalPos) / 5 / minutes
	}

	// Calculate uncorrected errors (errors still present in the text)
	uncorrectedErrors := len(g.Errors)

	// Calculate Net WPM (Gross WPM - uncorrected errors per minute)
	errorRate := float64(uncorrectedErrors) / minutes
	netWPM := grossWPM - errorRate

	// Ensure Net WPM doesn't go below 0
	if netWPM < 0 {
		netWPM = 0
	}

	// Calculate accuracy (correct characters / total characters typed * 100)
	correctChars := g.GlobalPos - g.TotalErrorsMade
	accuracy := 0.0
	if g.GlobalPos > 0 {
		accuracy = float64(correctChars) / float64(g.GlobalPos) * 100
	}

	// Ensure accuracy doesn't go below 0
	if accuracy < 0 {
		accuracy = 0
	}

	if _, err := database.InsertStats(&database.Stats{WPM: netWPM, Accuracy: accuracy}); err != nil {
		fmt.Println("Insert failed:", err)
	}

	maxStats, err := database.MaxStats()
	if err != nil {
		fmt.Println(err)
		maxStats = &database.MaxStatsData{MaxWPM: 0, MaxAccuracy: 0}
	}


	return TypingStats{
		WPM:               netWPM,
		BestWPM:		   maxStats.MaxWPM,
		Accuracy:          accuracy,
		BestAccuracy:	   maxStats.MaxAccuracy,
		CharactersTyped:   g.GlobalPos,
		CorrectChars:      correctChars,
		TotalChars:        len([]rune(g.GetDisplayText())),
		TimeElapsed:       elapsed,
		IsComplete:        g.IsFinished,
		UncorrectedErrors: uncorrectedErrors,
	}
}

// IsTimeUp checks if the game time has exceeded the specified duration
func (g *TypingGame) IsTimeUp() bool {
	if !g.IsStarted {
		return false
	}
	return time.Since(g.StartTime).Seconds() >= float64(g.Duration)
}

// GetRemainingTime returns the remaining time in seconds for the game
func (g *TypingGame) GetRemainingTime() int {
	if !g.IsStarted {
		return g.Duration
	}
	elapsed := int(time.Since(g.StartTime).Seconds())
	remaining := g.Duration - elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}
