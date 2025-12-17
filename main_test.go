package main

import "testing"

func TestParseAndValidate_BadChar(t *testing.T) {
	input := "....\n..#.\n..X.\n....\n"
	_, err := ParseAndValidate(input)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseAndValidate_GoodOnePiece(t *testing.T) {
	input := "#...\n#...\n#...\n#...\n"
	p, err := ParseAndValidate(input)
	if err != nil || len(p) != 1 {
		t.Fatalf("expected 1 piece, got len=%d err=%v", len(p), err)
	}
}
