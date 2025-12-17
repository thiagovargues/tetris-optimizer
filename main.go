package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

type Point struct {
	X int
	Y int
}

type Piece struct {
	Blocks []Point // normalized (top-left shifted to 0,0)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR")
		return
	}

	content, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("ERROR")
		return
	}

	pieces, err := ParseAndValidate(string(content))
	if err != nil || len(pieces) == 0 {
		fmt.Println("ERROR")
		return
	}
	if len(pieces) > 26 { // spec implies A..Z; safest for audits
		fmt.Println("ERROR")
		return
	}

	size := MinStartSize(len(pieces))
	for {
		board := NewBoard(size)
		if Solve(board, pieces, 0) {
			PrintBoard(board)
			return
		}
		size++
	}
}

// ---------------- Parsing + Validation ----------------

func ParseAndValidate(input string) ([]Piece, error) {
	// Normalize Windows newlines
	input = strings.ReplaceAll(input, "\r\n", "\n")

	lines := strings.Split(input, "\n")

	var pieces []Piece
	i := 0

	for i < len(lines) {
		// If we hit an empty line where a piece should start, that's invalid
		// (also rejects multiple blank lines).
		if lines[i] == "" {
			// allow a single trailing newline at end? only if it's truly the last empty
			// e.g. file ends with "\n" -> last element "" -> accept only if nothing else.
			onlyTrailing := true
			for j := i; j < len(lines); j++ {
				if lines[j] != "" {
					onlyTrailing = false
					break
				}
			}
			if onlyTrailing {
				break
			}
			return nil, fmt.Errorf("bad format: unexpected empty line")
		}

		// Need 4 lines for a tetromino
		if i+3 >= len(lines) {
			return nil, fmt.Errorf("bad format: incomplete piece")
		}

		block := lines[i : i+4]
		i += 4

		p, err := ValidateAndNormalize(block)
		if err != nil {
			return nil, err
		}
		pieces = append(pieces, p)

		// After 4 lines: must be end OR exactly one empty line separator
		if i == len(lines) {
			break
		}
		if lines[i] == "" {
			i++ // consume the separator blank line
			// If after consuming separator we are at end, that's okay (file ending newline)
			// but multiple blank lines will be caught at top of loop.
			continue
		}

		// If not empty line separator, then format is wrong
		return nil, fmt.Errorf("bad format: missing blank separator")
	}

	if len(pieces) == 0 {
		return nil, fmt.Errorf("no pieces")
	}
	return pieces, nil
}

func ValidateAndNormalize(block []string) (Piece, error) {
	if len(block) != 4 {
		return Piece{}, fmt.Errorf("bad piece size")
	}

	grid := make([][]rune, 4)
	hashCount := 0

	for y := 0; y < 4; y++ {
		if len(block[y]) != 4 {
			return Piece{}, fmt.Errorf("bad line length")
		}
		row := []rune(block[y])
		grid[y] = row
		for x := 0; x < 4; x++ {
			if row[x] != '.' && row[x] != '#' {
				return Piece{}, fmt.Errorf("invalid char")
			}
			if row[x] == '#' {
				hashCount++
			}
		}
	}
	if hashCount != 4 {
		return Piece{}, fmt.Errorf("must have 4 blocks")
	}

	// Touch-count validation: total neighbor links must be 6 or 8
	touches := 0
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if grid[y][x] != '#' {
				continue
			}
			if y > 0 && grid[y-1][x] == '#' {
				touches++
			}
			if y < 3 && grid[y+1][x] == '#' {
				touches++
			}
			if x > 0 && grid[y][x-1] == '#' {
				touches++
			}
			if x < 3 && grid[y][x+1] == '#' {
				touches++
			}
		}
	}
	if touches != 6 && touches != 8 {
		return Piece{}, fmt.Errorf("not connected tetromino")
	}

	// Extract points
	points := make([]Point, 0, 4)
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if grid[y][x] == '#' {
				points = append(points, Point{X: x, Y: y})
			}
		}
	}

	// Normalize to top-left
	minX, minY := points[0].X, points[0].Y
	for _, p := range points {
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
	}
	for k := range points {
		points[k].X -= minX
		points[k].Y -= minY
	}

	// Sort for consistency (helps debugging; not required)
	sort.Slice(points, func(i, j int) bool {
		if points[i].Y == points[j].Y {
			return points[i].X < points[j].X
		}
		return points[i].Y < points[j].Y
	})

	return Piece{Blocks: points}, nil
}

// ---------------- Solver ----------------

func MinStartSize(nPieces int) int {
	blocks := 4 * nPieces
	root := math.Sqrt(float64(blocks))
	size := int(math.Ceil(root))
	if size < 2 {
		size = 2
	}
	return size
}

func NewBoard(size int) [][]rune {
	b := make([][]rune, size)
	for y := 0; y < size; y++ {
		b[y] = make([]rune, size)
		for x := 0; x < size; x++ {
			b[y][x] = '.'
		}
	}
	return b
}

func Solve(board [][]rune, pieces []Piece, idx int) bool {
	if idx == len(pieces) {
		return true
	}

	size := len(board)
	ch := rune('A' + idx)
	p := pieces[idx]

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if CanPlace(board, p, x, y) {
				Place(board, p, x, y, ch)
				if Solve(board, pieces, idx+1) {
					return true
				}
				Place(board, p, x, y, '.') // undo
			}
		}
	}

	return false
}

func CanPlace(board [][]rune, piece Piece, ox, oy int) bool {
	size := len(board)
	for _, b := range piece.Blocks {
		x := ox + b.X
		y := oy + b.Y
		if x < 0 || y < 0 || x >= size || y >= size {
			return false
		}
		if board[y][x] != '.' {
			return false
		}
	}
	return true
}

func Place(board [][]rune, piece Piece, ox, oy int, ch rune) {
	for _, b := range piece.Blocks {
		x := ox + b.X
		y := oy + b.Y
		board[y][x] = ch
	}
}

func PrintBoard(board [][]rune) {
	for _, row := range board {
		fmt.Println(string(row))
	}
}
