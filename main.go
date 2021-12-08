package main

import (
	"fmt"
	"io/ioutil"
	"log"
	filepath2 "path/filepath"
	"strconv"
)

var nProc = 8
var mem = 16
var disk = 200
var basisSet = "Def2TZVP"
var DFT = "PBE1PBE"
var charge = 3
var shellDist = 2.8
var ion = "LA"

var aminoAcidTally = make(map[string]int)

// Program begins here
func main() {
	inDir := "C:\\Users\\jtgou\\OneDrive\\Documents\\UT_Austin\\ren_lab\\lanthanides\\lanthanides\\input"
	outDir := "C:\\Users\\jtgou\\OneDrive\\Documents\\UT_Austin\\ren_lab\\lanthanides\\lanthanides\\output"
	fmt.Println("Processing directory at: " + inDir)

	totalSystems := 0
	totalResidues := 0
	totalStructures := 0

	// Read in all files in dir
	fileInfo, err := ioutil.ReadDir(inDir)
	if err != nil {
		fmt.Println("failed to read directory: " + inDir)
		log.Fatal(err)
	}

	for _, file := range fileInfo {
		if filepath2.Ext(file.Name()) == ".pdb" {
			sysNum, resNum := pdb2Systems(filepath2.Join(inDir, file.Name()), outDir)
			totalSystems += sysNum
			totalResidues += resNum
			totalStructures++
		}
	}
	finalizeAATallies()
	fmt.Println("Processed " + strconv.Itoa(totalSystems) + " Systems from " + strconv.Itoa(totalStructures) + " Structures")
	fmt.Println("Total Residues: " + strconv.Itoa(totalResidues))
}

func pdb2Systems(path string, dir string) (int, int) {
	residueCount := 0
	fmt.Println("Reading in file at: " + path)
	sysName, atoms := pdbReader(path)
	fmt.Println("Read in file with " + strconv.Itoa(len(atoms)) + " atoms.")
	fmt.Println("Finding ion systems in file...")
	systems := structure2Systems(atoms, ion, shellDist)
	fmt.Println("Found " + strconv.Itoa(len(systems)) + " systems.")
	fmt.Println("Writing systems...")
	for i, system := range systems {
		outName := sysName + "_" + strconv.Itoa(i)
		residueNum := writeSystemGJF(*system, ion, dir, outName)
		residueCount += residueNum
	}
	addToAATallies(systems)
	return len(systems), residueCount
}

func addToAATallies(systems []*ionSystem) {
	for _, system := range systems {
		for _, v := range system.residueList {
			if _, ok := aminoAcidTally[v]; ok {
				aminoAcidTally[v] += 1
			} else {
				aminoAcidTally[v] = 1
			}
		}
		aminoAcidTally["LA"]--
	}

	// remove duplicate La

}

func finalizeAATallies() {

	total := 0
	for _, v := range aminoAcidTally {
		total += v
	}
	floatTotal := float64(total)
	fmt.Println("\nAmino Acids in Binding Pockets (n = " + strconv.Itoa(total) + " residues): ")
	for k, v := range aminoAcidTally {
		floatV := float64(v)
		fmt.Println(k + " = " + fmt.Sprintf("%.3f", floatV / floatTotal))
	}
}