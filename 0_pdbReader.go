package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	filepath2 "path/filepath"
	"strconv"
	"strings"
)

func pdbReader(filePath string) (string, map[int]*atom) {
	// open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open molecule file: " + filePath)
		log.Fatal(err)
	}
	structureName := strings.Split(filepath2.Base(filePath),".")[0]

	// Create structure to store atoms
	atoms := make(map[int]*atom)

	// Initialize scanner
	scanner := bufio.NewScanner(file)
	// ignore first line
	scanner.Scan()
	// create line counter
	i := 1
	// iterate over all other lines
	for scanner.Scan() {
		// get next line
		line := scanner.Text()

		// split by whitespace
		tokens := strings.Fields(line)

		// check line length before proceeding
		if len(tokens) >= 8 && (tokens[0] == "ATOM" || tokens[0] == "HETATM") {

			// create new atom
			var newAtom atom

			// get number of atom from file
			atomNum, err := strconv.Atoi(tokens[1])
			if err != nil {
				newErr := errors.New("Failed to convert token" + "\"" + tokens[0] + "\"" + " in position 0 on line " + strconv.Itoa(i) + " to an integer")
				log.Fatal(newErr)
			}
			newAtom.id = atomNum

			// assign element
			newAtom.element = tokens[2]
			newAtom.aminoAcid = tokens[3]
			newAtom.residue = tokens[5]

			// assign positions
			pos := make([]float64,3)
			for j := 6; j < 9; j++ {
				pos[j-6], err = strconv.ParseFloat(tokens[j],64)
				if err != nil {
					newErr := errors.New("Failed to convert \"" + tokens[j] + "\" in position 0 on line " + strconv.Itoa(i) + " to a float64")
					log.Fatal(newErr)
				}
			}
			newAtom.pos = pos

			// add atom to map
			atoms[atomNum] = &newAtom
		}
		i++
	}

	return structureName, atoms
}
