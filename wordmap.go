package main

// Simple string set backed by a map, used to avoid selecting duplicate
// words from different walks of the board. Also keep track of the longest
// word, removing the need to sort to find the longest word.
type WordSet interface {
	Add(words ...string)
	Values() []string
	Longest() []string
	Count() int
}

func NewWordSet() (result WordSet) {
	result = &mapWordSet{
		words:      map[string]byte{},
		longestLen: 0,
		longest:    []string{},
	}
	return
}

// WordSet backed by a map
type mapWordSet struct {
	// Only the key is useful, ignore the value
	words      map[string]byte
	longestLen int
	longest    []string
}

// Add a word to the set
func (m *mapWordSet) Add(words ...string) {
	for _, word := range words {
		if _, exists := m.words[word]; exists {
			continue
		}

		m.words[word] = 0 // value is unused

		// Keep track of longest words
		wordLen := len(word)
		if wordLen > m.longestLen {
			// Word is longer than any found so far, toss old words
			m.longest = []string{word}
			m.longestLen = wordLen
		} else if wordLen == m.longestLen {
			// Word is a tie
			m.longest = append(m.longest, word)
		}
	}
}

// Returns the longest word in the set
func (m *mapWordSet) Longest() []string {
	return m.longest
}

// Extracts values as a new slice
func (m *mapWordSet) Values() []string {
	values := make([]string, len(m.words))
	i := 0
	for word := range m.words {
		values[i] = word
		i++
	}
	return values
}

func (m *mapWordSet) Count() int {
	return len(m.words)
}
