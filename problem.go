package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var grid [][]string = [][]string{
	{"e", "x", "t", "r", "a", "h", "o", "p"},
	{"n", "e", "t", "w", "o", "r", "k", "s"},
	{"q", "i", "h", "a", "c", "i", "q", "t"},
	{"l", "f", "u", "n", "u", "r", "x", "b"},
	{"b", "w", "d", "i", "l", "a", "t", "v"},
	{"o", "s", "s", "y", "n", "a", "c", "k"},
	{"q", "w", "o", "p", "m", "t", "c", "p"},
	{"k", "i", "p", "a", "c", "k", "e", "t"},
}

//var testgrid [][]string = [][]string{
//	{"q", "w", "e", "r", "t", "n", "u", "i"},
//	{"o", "p", "a", "a", "d", "f", "g", "h"},
//	{"t", "k", "l", "z", "x", "c", "v", "b"},
//	{"n", "m", "r", "w", "f", "r", "t", "y"},
//	{"u", "i", "o", "p", "a", "s", "d", "f"},
//	{"g", "h", "j", "o", "l", "z", "x", "c"},
//	{"v", "b", "n", "m", "q", "w", "e", "r"},
//	{"t", "y", "u", "i", "o", "p", "a", "s"},
//}

// Candidate knight moves, relative values
var knightSpots []Pos = []Pos{Pos{-2, -1}, Pos{-1, -2}, Pos{1, -2}, Pos{2, -1}, Pos{2, 1}, Pos{1, 2}, Pos{-1, 2}, Pos{-2, 1}}

// Memo map for valid moves
var moves map[Pos][]Pos = make(map[Pos][]Pos)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Must provide starting row and column numbers as arguments, starting at 1")
		os.Exit(1)
	}
	row, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("First arg must be a row number")
		os.Exit(1)
	}
	row-- // change to 0 index
	col, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Second arg must be a column number")
		os.Exit(1)
	}
	col-- // change to 0 index
	reader := bufio.NewReader(os.Stdin)
	largestLen := 0
	var words []string
	var line string
	for err == nil {
		line, err = reader.ReadString('\n')
		line = strings.Trim(line, " \r\n\t")
		if line == "" {
			continue
		}
		if len(line) > largestLen {
			largestLen = len(line)
		}
		words = append(words, line)
	}
	results := findLongest(words, row, col, largestLen)
	fmt.Println(results)
}

func findLongest(words []string, row, col, maxDepth int) (results []string) {
	resultSet := NewWordSet()
	find(words, Pos{row, col}, "", 0, maxDepth, resultSet)
	return resultSet.Longest()
}

// Finds a list of words that can be created from a given starting position.
// maxDepth prevents the function from walking the board forever,
// should be set to the longest word in the source data
func find(words []string, p Pos, prefix string, depth, maxDepth int, results WordSet) {
	if depth > maxDepth {
		return
	}
	prefix += grid[p.Row][p.Col]
	narrowed := narrow(words, prefix)
	if len(narrowed) < 2 {
		results.Add(narrowed...)
		return
	}
	for _, pos := range p.GetMoves() {
		find(narrowed, pos, prefix, depth+1, maxDepth, results)
	}
	return
}

// This function takes a slice of words, a prefix string, and returns a subslice of only words
// that begin with that prefix. It takes advantage of the lexical ordering to find the upper bound of the slice.
// Rather than seek to find the last word that begins with prefix, it performs another binary search with prefix
// but the last letter incremented. E.g. a prefix of "ab" will be incremented to "ac"
func narrow(words []string, prefix string) []string {
	prefix2 := prefix[:len(prefix)-1] + string(prefix[len(prefix)-1]+1)
	bottom := sort.Search(len(words), func(i int) bool { return strings.Compare(words[i], prefix) >= 0 })
	top := sort.Search(len(words), func(i int) bool { return strings.Compare(words[i], prefix2) >= 0 })
	return words[bottom:top]
}

// Encapsulates row and column positions
type Pos struct {
	Row int
	Col int
}

// Returns a list of valid moves for a given position, pruning ones not on the board
// Moves are relative [r, c] coordinates
func (p Pos) GetMoves() (results []Pos) {
	var ok bool
	results, ok = moves[p]
	if ok {
		return
	}
	for _, pos := range knightSpots {
		if p.IsMoveInGrid(pos) {
			results = append(results, p.Add(pos))
		}
	}
	moves[p] = results
	return
}

// Checks whether or not a given move will be inside the grid
func (p Pos) IsMoveInGrid(move Pos) bool {
	new := p.Add(move)
	return new.Row > 0 && new.Col > 0 && new.Row < len(grid) && new.Col < len(grid[0])
}

// Adds two positions together, just like adding vectors
func (p Pos) Add(other Pos) Pos {
	return Pos{p.Row + other.Row, p.Col + other.Col}
}

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

// Pretty print the set
func (m *mapWordSet) String() string {
	return strings.Join(m.Values(), ",")
}
