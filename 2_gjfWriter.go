package main

import (
	"fmt"
	"log"
	"os"
	filepath2 "path/filepath"
	"strconv"
)

func writeSystemGJF(system ionSystem, ion string, dir string, fileName string) {
	os.MkdirAll(dir, 0755)
	thisPath := filepath2.Join(dir, fileName + ".gjf")
	fmt.Println("Creating file at " + thisPath)
	thisFile, err := os.Create(thisPath)
	if err != nil {
		fmt.Println("Failed to create new fragment file: " + thisPath)
		log.Fatal(err)
	}

	// write header
	_, err = thisFile.WriteString("%chk=" + fileName + ".chk\n\n")
	_, err = thisFile.WriteString("%nproc=" + strconv.Itoa(nProc) + "\n")
	_, err = thisFile.WriteString("%mem=" + strconv.Itoa(mem) + "gb\n")
	_, err = thisFile.WriteString("#p " + basisSet + " " + DFT + " maxDisk=" + strconv.Itoa(disk) + " nosymm Counterpoise=2\n\n")
	_, err = thisFile.WriteString("Counterpoise\n\n")
	_, err = thisFile.WriteString("! Num Residues = " + strconv.Itoa(len(system.residueList)) + "\n")
	line := "! Residue List = "
	for k, v := range system.residueList {
		line += k + ":" + v + ", "
	}
	_, err = thisFile.WriteString(line + "\n\n")



	// write body
	_, err = thisFile.WriteString(strconv.Itoa(charge) + ",1 0,1 " + strconv.Itoa(charge) + ",1\n")
	j := 1
	for _, thisAtom := range system.atoms {

		line := thisAtom.element
		if thisAtom.element == ion {
			line += "(Fragment=2) "
		} else {
			line += "(Fragment=1) "
		}

		line += "\t" + fmt.Sprintf("%.6f", thisAtom.pos[0]) + "\t" +
			fmt.Sprintf("%.6f", thisAtom.pos[1]) + "\t" + fmt.Sprintf("%.6f", thisAtom.pos[2])

		// add comment
		line += "\t\t! " + thisAtom.residue + " " + thisAtom.aminoAcid

		// write to file
		_, err := thisFile.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Failed to write atom")
			log.Fatal(err)
		}
		j++
	}
	_, err = thisFile.WriteString("\n")
}
