Tetris Optimizer (Tetromino Square Solver)
Overview

This program reads a list of tetrominoes from a text file and assembles them in input order into the smallest possible square.
Each tetromino is labeled in the output using uppercase letters:

A for the first tetromino

B for the second

etc.

Empty cells are printed as ..

If the input file format is invalid or a tetromino is invalid, the program prints:

ERROR

Requirements

Go (standard library only)

The program expects exactly one argument: a path to the input text file.

Usage

From the project folder:

go run . <path_to_file>


Example:

go run . sample.txt

Input File Format

Each tetromino is a 4×4 block

A tetromino block is 4 lines

Each line is exactly 4 characters

Allowed characters: . and #

Tetrominoes are separated by one empty line

Each tetromino must contain exactly 4 #, connected by edges (not diagonals)

Example input (sample.txt)
#...
#...
#...
#...

....
.##.
.##.
....

Output Format

Prints the smallest square solution found

Uses letters A, B, C, ...

Uses . for empty spaces

Example output (one possible valid result):

AB..
AB..
A...
A...

Error Handling

The program prints ERROR if any of the following happen:

wrong number of arguments

file cannot be read

bad file format (wrong line lengths, wrong separators, invalid characters, incomplete blocks, etc.)

invalid tetromino (not exactly 4 blocks, not connected)

How It Works

Read and parse tetrominoes from the file

Validate and normalize each tetromino into coordinates

Start from the smallest possible square size (based on total blocks)

Use backtracking to place pieces in order

If it doesn’t fit, increase the square size and try again

Tests

Run unit tests:

go test ./...


You can also manually test with files in tests/:

go run . tests/good02.txt
go run . tests/bad00.txt

Project Structure (typical)

main.go — entry point, reading args, running solver

main_test.go — unit tests

tests/ — sample good/bad input files

README.md — documentation

Notes

Only standard Go packages are used.

Output solutions may differ (multiple solutions can exist), but the square size should be minimal and tetromino labels must match their input order.
