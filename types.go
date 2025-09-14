package main

type PatronFile struct {
	Using map[string]bool
}

func NewPatronFile() *PatronFile {
	return &PatronFile{
		Using: make(map[string]bool),
	}
}
