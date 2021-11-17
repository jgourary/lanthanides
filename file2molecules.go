package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	filepath2 "path/filepath"
	"strconv"
	"strings"
	"unicode"
)


func file2systems(path string, ion string, count int, outDir string) {
	sysName, frames := arc2frames(path)
	fmt.Println("Read in " + strconv.Itoa(len(frames)) + " frames from " + sysName)

	// get radial info
	fmt.Println("Starting radial")
	shell1Dist, shell2Dist := runRadial(path, ion, len(frames))
	fmt.Println("Ending radial")
	fmt.Println(ion + " Shell 1 = " + fmt.Sprintf("%0.2f", shell1Dist) + " angstroms")
	fmt.Println(ion + " Shell 2 = " + fmt.Sprintf("%0.2f", shell2Dist) + " angstroms")

	shell1Systems := frames2Systems(frames, ion, count, shell1Dist)
	shell2Systems := frames2Systems(frames, ion, count, shell2Dist)

	fmt.Println("Read in " + strconv.Itoa(len(shell1Systems[0])) + " systems per frame with " + strconv.Itoa(len(shell1Systems[0][0].atoms)) )

	for i, frame := range shell1Systems {
		for j, system := range frame {
			renumberedAtoms := copySystemWithRenumbering(system.atoms)
			writeSystemXYZ(renumberedAtoms, filepath2.Join(outDir, "shell_1"), sysName + "_" + ion + "_" + strconv.Itoa(i) + "_" + strconv.Itoa(j))
			writeSystemINP(renumberedAtoms, system.charge, filepath2.Join(outDir, "shell_1"), sysName + "_" + ion + "_" + strconv.Itoa(i) + "_" + strconv.Itoa(j))
		}
	}
	for i, frame := range shell2Systems {
		for j, system := range frame {
			renumberedAtoms := copySystemWithRenumbering(system.atoms)
			writeSystemXYZ(renumberedAtoms, filepath2.Join(outDir, "shell_2"), sysName + "_" + ion + "_" + strconv.Itoa(i) + "_" + strconv.Itoa(j))
			writeSystemINP(renumberedAtoms, system.charge, filepath2.Join(outDir, "shell_2"), sysName + "_" + ion + "_" + strconv.Itoa(i) + "_" + strconv.Itoa(j))
		}
	}
	fmt.Println("Printed Systems")
}

func frames2Systems(frames []map[int]*atom, ion string, count int, shellDist float64) [][]ionSystem {
	var frameSystems [][]ionSystem
	for _, frame := range frames {
		systems := frame2Systems(frame, ion, count, shellDist)
		frameSystems = append(frameSystems, systems)
	}
	return frameSystems
}

func frame2Systems(atoms map[int]*atom, ion string, count int, shellDist float64) []ionSystem {
	var ionSystems []ionSystem
	counter := 0

	// initialize systems with center ion
	for id, thisAtom := range atoms {
		if counter < count && thisAtom.element == ion {

			var newSystem ionSystem
			atomList := make(map[int]*atom)
			atomList[id] = thisAtom

			newSystem.atoms = atomList
			newSystem.center = thisAtom.pos

			if ion == "Cl-" {
				newSystem.charge = -1
			} else {
				newSystem.charge = 1
			}

			ionSystems = append(ionSystems, newSystem)
			counter++
		}
	}

	// add atoms from all neighboring waters within distance specified
	for _, system := range ionSystems {
		for id, thisAtom := range atoms {
			// add oxygen
			if thisAtom.element == "O" && getDistance(thisAtom, system.center) < shellDist {
				system.atoms[id] = thisAtom
				// add bound hydrogens
				for _, bondedAtomID := range atoms[id].bondedAtoms {
					system.atoms[bondedAtomID] = atoms[bondedAtomID]
				}
			}
		}
	}

	return ionSystems
}

func getDistance(atom1 *atom, pos []float64) float64 {
	dx2 := math.Pow(atom1.pos[0] - pos[0], 2)
	dy2 := math.Pow(atom1.pos[1] - pos[1], 2)
	dz2 := math.Pow(atom1.pos[2] - pos[2], 2)

	return math.Sqrt(dx2 + dy2 + dz2)
}




// Includes renumbering atom map keys from 1 to len(molecule)-1
func copySystemWithRenumbering(atomsA map[int]*atom) []*atom {

	atomIDOldToNewMap := make(map[int]int)
	i := 1
	for _, atom := range atomsA {
		if atom.isIon {
			atomIDOldToNewMap[atom.id] = i
			i++
		}
	}
	for _, atom := range atomsA {
		if atom.element == "O" {
			atomIDOldToNewMap[atom.id] = i
			i++
			for _, bondedAtomID := range atom.bondedAtoms {
				atomIDOldToNewMap[bondedAtomID] = i
				i++
			}
		}
	}

	atomsB := make([]*atom, len(atomsA))
	for oldIndex, thisAtom := range atomsA {
		newAtom := copyAtom(thisAtom)
		newIndex := atomIDOldToNewMap[oldIndex]
		newAtom.id = newIndex
		for j := 0; j < len(newAtom.bondedAtoms); j++ {
			newAtom.bondedAtoms[j] = atomIDOldToNewMap[newAtom.bondedAtoms[j]]
		}
		atomsB[newIndex] = &newAtom
	}
	return atomsB
}

func arc2frames(filePath string) (string, []map[int]*atom) {



	// open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open molecule file: " + filePath)
		log.Fatal(err)
	}
	structureName := strings.Split(filepath2.Base(filePath),".")[0]

	var frames []map[int]*atom

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
		if len(tokens) >= 6 && unicode.IsLetter(rune(tokens[1][0])) {
			// create new atom
			var newAtom atom

			// get number of atom from file
			atomNum, err := strconv.Atoi(tokens[0])
			if err != nil {
				newErr := errors.New("Failed to convert token" + "\"" + tokens[0] + "\"" + " in position 0 on line " + strconv.Itoa(i) + " to an integer")
				log.Fatal(newErr)
			}
			newAtom.id = atomNum

			// assign element
			newAtom.element = tokens[1]

			// assign positions
			pos := make([]float64,3)
			for j := 2; j < 5; j++ {
				pos[j-2], err = strconv.ParseFloat(tokens[j],64)
				if err != nil {
					newErr := errors.New("Failed to convert token in position 0 on line " + strconv.Itoa(j) + " to a float64")
					log.Fatal(newErr)
				}
			}
			newAtom.pos = pos

			// assign atomType from file
			newAtom.atomType, err = strconv.Atoi(tokens[5])
			if err != nil {
				newErr := errors.New("Failed to convert token in position 5 on line " + strconv.Itoa(i) + " to an integer")
				log.Fatal(newErr)
			}

			// assign bonds from file
			bonds := make([]int,len(tokens)-6)
			for j := 6; j < len(tokens); j++ {
				bonds[j-6], err = strconv.Atoi(tokens[j])
				if err != nil {
					newErr := errors.New("Failed to convert token in position " + strconv.Itoa(j) + " on line " + strconv.Itoa(i) + " to an integer")
					log.Fatal(newErr)
				}
			}
			newAtom.bondedAtoms = bonds

			// add atom to map
			atoms[atomNum] = &newAtom
		} else {
			//fmt.Println("Warning: line " + strconv.Itoa(i) + " has insufficient tokens. Program is skipping this " +
			//	"line when reading your input file.")
		}
		i++
	}

	return structureName, frames
}

func copyAtom(thisAtom *atom) atom {
	var newAtom atom


	newAtom.element = thisAtom.element
	newAtom.mass = thisAtom.mass
	newAtom.atomType = thisAtom.atomType
	newAtom.id = thisAtom.id
	newAtom.bondedAtoms = make([]int, len(thisAtom.bondedAtoms))
	for i := range thisAtom.bondedAtoms {
		newAtom.bondedAtoms[i] = thisAtom.bondedAtoms[i]
	}

	newAtom.pos = thisAtom.pos
	newAtom.charmmID = thisAtom.charmmID


	return newAtom
}