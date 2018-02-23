package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var Grid [][]string = [][]string{
	{"e", "x", "t", "r", "a", "h", "o", "p"},
	{"n", "e", "t", "w", "o", "r", "k", "s"},
	{"q", "i", "h", "a", "c", "i", "q", "t"},
	{"l", "f", "u", "n", "u", "r", "x", "b"},
	{"b", "w", "d", "i", "l", "a", "t", "v"},
	{"o", "s", "s", "y", "n", "a", "c", "k"},
	{"q", "w", "o", "p", "m", "t", "c", "p"},
	{"k", "i", "p", "a", "c", "k", "e", "t"},
}

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
	fmt.Println(strings.Join(results, ", "))
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
	prefix += Grid[p.Row][p.Col]
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
