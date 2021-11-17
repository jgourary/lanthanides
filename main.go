package main

import (
	"fmt"
	"strconv"
)

var nProc = 8
var mem = 16
var disk = 200
var basisSet = "Def2TZVP"
var DFT = "PBE1PBE"
var charge = 3
var shellDist = 3.0
var ion = "LA"

// Program begins here
func main() {
	path := "C:\\Users\\jtgou\\lanthanides\\test\\7cco.pdb"
	outDir := "C:\\Users\\jtgou\\lanthanides\\test"
	fmt.Println("Processing file at: " + path)
	pdb2Systems(path, outDir)
}

func pdb2Systems(path string, dir string) {
	fmt.Println("Reading in file..")
	sysName, atoms := pdbReader(path)
	fmt.Println("Read in file with " + strconv.Itoa(len(atoms)) + " atoms.")
	fmt.Println("Finding ion systems in file...")
	systems := structure2Systems(atoms, ion, shellDist)
	fmt.Println("Found " + strconv.Itoa(len(systems)) + " systems.")
	fmt.Println("Writing systems...")
	for i, system := range systems {
		outName := sysName + "_" + strconv.Itoa(i)
		writeSystemGJF(system.atoms, ion, dir, outName)
	}

}