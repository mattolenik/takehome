## Prerequesites
Just Go.

## How to Run

### Run All

To find all words across all grid positions, execute the `run-all.sh` script in a terminal. The output will be something like this:

```
Result for row 2, col 4: witty

Result for row 2, col 5: orthography

Result for row 2, col 6: rare, rude

Result for row 2, col 7:

Result for row 2, col 8: situate, sitting
```

### Program Usage

The program has the grid hardcoded and takes the word list as stdin, with row and column for arguments.

```
$ go run *.go < words.txt 3 4
```

This will find the longest word starting at grid position [3, 4].
```
although
```


## Algorithm

1. Load words into array
2. Get letter at current grid position (call this 'prefix')
3. Narrow list of words to those beginning with said prefix
    1. Do binary search for prefix to find lower bound, the first value to start with prefix
    2. Do binary search for prefix but increment last character, e.g. "abc" becomes "abd" -- this selects the upper bound, the last value to start with prefix
    3. Return a slice narrowed to those boundaries, i.e. `words[lower:upper]`
    4. Go slices do not reallocate so this substep only takes 2log(n) time and no additional space
4. For each valid knight move, recursively repeat step 2 for each new position
5. The recursive call completes when the list has been narrowed to one item (the longest word, the result we want), or when the list has been narrowed to an empty list (no word can be formed)
    1. The recursion will also terminate at a max depth, which is set to the length of the longest word in the input

A small set type, backed by a map, was made to keep track of found words. This type keeps track of the longest found words and removes duplicates. When the search is done, this set will already contain the longest word (or words if there are ties). This is as opposed to e.g. adding all the matches to a list, sorting the list, then taking the top result.

### Time Complexity

Time complexity here is a function of grid size, number of words, and the length of the words

x = max depth (length of the longest word)
g = num grid spaces
n = number of words

Time complexity is:
O(g⋅x⋅2⋅log(n))

2⋅log(n) comes from the two binary searches described in step 2 above. That step is performed once for each level of recursive, which happens at most x times per spot, and then once for each spot on the grid. For a given grid size, assuming the words are of reasonble length, the complexity will roughly follow O(2⋅log(n)).


### Space Complexity

Space complexity is constant, requiring space for size of the grid g, number of words n. Practically speaking, O(n).

## Assumptions

Input words are all lowercase, stripped of all punctuation (all characters removed, not replaced with spaces or anything else).

```
```
