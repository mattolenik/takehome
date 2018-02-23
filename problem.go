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

type Pos struct {
	Row int
	Col int
}

func (p Pos) IsMoveInGrid(move Pos) bool {
	new := p.Add(move)
	return new.Row > 0 && new.Col > 0 && new.Row < len(grid) && new.Col < len(grid[0])
}

func (p Pos) Add(other Pos) Pos {
	return Pos{p.Row + other.Row, p.Col + other.Col}
}

// Candidate knight moves, relative values
var knightSpots []Pos = []Pos{Pos{-2, -1}, Pos{-1, -2}, Pos{1, -2}, Pos{2, -1}, Pos{2, 1}, Pos{1, 2}, Pos{-1, 2}, Pos{-2, 1}}

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
	row--
	col, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Second arg must be a column number")
		os.Exit(1)
	}
	col--
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
	if len(words) < 2 {
		return words
	}
	longest := find(words, Pos{row, col}, "", 0, maxDepth)
	if len(longest) < 2 {
		return longest
	}
	sort.Slice(longest, func(i, j int) bool { return len(longest[i]) > len(longest[j]) })
	length := len(longest[0])
	// Possibly multiple ties in length
	for _, word := range longest {
		if len(word) == length {
			results = append(results, word)
		}
	}
	return results
}

func find(words []string, p Pos, prefix string, depth, maxDepth int) (results []string) {
	charAtP := grid[p.Row][p.Col]
	prefix += charAtP
	narrowed := narrow(words, prefix)
	if len(narrowed) < 2 {
		return narrowed
	}
	if depth > maxDepth {
		return []string{}
	}
	for _, pos := range getMoves(p) {
		//fmt.Printf("%+v, %+v\n", p, pos)
		words := find(narrowed, pos, prefix, depth+1, maxDepth)
		results = append(results, words...)
	}
	return results
}

// This function takes a slice of words, a prefix string, and returns a subslice of only words
// that begin with that prefix. It takes advantage of the lexical ordering to find the upper bound of the slice.
// Rather than seek to find the last word that begins with prefix, it performs another binary search with prefix
// but the last letter incremented. E.g. a prefix of "ab" will be incremented to "ac"
func narrow(words []string, prefix string) []string {
	prefix2 := prefix[:len(prefix)-1] + string(prefix[len(prefix)-1]+1)
	bottom := sort.Search(len(words), func(i int) bool { return strings.Compare(words[i], prefix) >= 0 })
	top := sort.Search(len(words), func(i int) bool { return strings.Compare(words[i], prefix2) >= 0 })
	//fmt.Printf("%s %s %d %d %+v\n", prefix, prefix2, bottom, top, words[bottom:top])
	return words[bottom:top]
}

func getMoves(p Pos) (results []Pos) {
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
