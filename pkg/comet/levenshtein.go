// Copied from https://github.com/texttheater/golang-levenshtein/blob/3d00ed8319236cf52f93aa60714ce1902e864c4a/levenshtein/levenshtein.go
// The MIT License (MIT)
//
// # Copyright (c) 2013 Kilian Evang and contributors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package comet

// Implements the Levenshtein algorithm for computing the
// similarity between two strings. The central function is MatrixForStrings,
// which computes the Levenshtein matrix. The functions DistanceForMatrix,
// EditScriptForMatrix and RatioForMatrix read various interesting properties
// off the matrix. The package also provides the convenience functions
// DistanceForStrings, EditScriptForStrings and RatioForStrings for going
// directly from two strings to the property of interest.

type EditOperation int

const (
	Ins = iota
	Del
	Sub
	Match
)

type EditScript []EditOperation

type MatchFunction func(rune, rune) bool

// IdenticalRunes is the default MatchFunction: it checks whether two runes are
// identical.
func IdenticalRunes(a rune, b rune) bool {
	return a == b
}

type Options struct {
	InsCost int
	DelCost int
	SubCost int
	Matches MatchFunction
}

// DefaultOptions is the default options without substitution: insertion cost
// is 1, deletion cost is 1, substitution cost is 2 (meaning insert and delete
// will be used instead), and two runes match iff they are identical.
var DefaultOptions = Options{
	InsCost: 1,
	DelCost: 1,
	SubCost: 2,
	Matches: IdenticalRunes,
}

func (operation EditOperation) String() string {
	if operation == Match {
		return "match"
	} else if operation == Ins {
		return "ins"
	} else if operation == Sub {
		return "sub"
	}
	return "del"
}

// DistanceForStrings returns the edit distance between source and target.
//
// It has a runtime proportional to len(source) * len(target) and memory use
// proportional to len(target).
func DistanceForStrings(source []rune, target []rune, op Options) int {
	// Note: This algorithm is a specialization of MatrixForStrings.
	// MatrixForStrings returns the full edit matrix. However, we only need a
	// single value (see DistanceForMatrix) and the main loop of the algorithm
	// only uses the current and previous row. As such we create a 2D matrix,
	// but with height 2 (enough to store current and previous row).
	height := len(source) + 1
	width := len(target) + 1
	matrix := make([][]int, 2)

	// Initialize trivial distances (from/to empty string). That is, fill
	// the left column and the top row with row/column indices multiplied
	// by deletion/insertion cost.
	for i := 0; i < 2; i++ {
		matrix[i] = make([]int, width)
		matrix[i][0] = i * op.DelCost
	}
	for j := 1; j < width; j++ {
		matrix[0][j] = j * op.InsCost
	}

	// Fill in the remaining cells: for each prefix pair, choose the
	// (edit history, operation) pair with the lowest cost.
	for i := 1; i < height; i++ {
		cur := matrix[i%2]
		prev := matrix[(i-1)%2]
		cur[0] = i * op.DelCost
		for j := 1; j < width; j++ {
			delCost := prev[j] + op.DelCost
			matchSubCost := prev[j-1]
			if !op.Matches(source[i-1], target[j-1]) {
				matchSubCost += op.SubCost
			}
			insCost := cur[j-1] + op.InsCost
			cur[j] = min(delCost, min(matchSubCost, insCost))
		}
	}
	return matrix[(height-1)%2][width-1]
}
