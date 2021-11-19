package main

import (
	"math"
)

func structure2Systems(atoms map[int]*atom, ion string, shellDist float64) []*ionSystem {
	var ionSystems []*ionSystem
	// initialize systems with center ion
	for id, thisAtom := range atoms {
		if thisAtom.element == ion {

			var newSystem ionSystem
			atomList := make(map[int]*atom)
			atomList[id] = thisAtom

			newSystem.atoms = atomList
			newSystem.center = thisAtom.pos

			ionSystems = append(ionSystems, &newSystem)
		}
	}

	// add atoms from all neighboring waters within distance specified
	for _, system := range ionSystems {
		approvedResidues := make(map[string]string)

		// mark residues in range
		for _, thisAtom := range atoms {
			if getDistance(thisAtom, system.center) < shellDist {

				approvedResidues[thisAtom.residue] = thisAtom.aminoAcid
			}
		}

		// add atoms in approved residues
		for id, thisAtom := range atoms {
			if _, ok := approvedResidues[thisAtom.residue]; ok {
				system.atoms[id] = thisAtom
			}
		}

		// add residue record
		system.residueList = approvedResidues

	}
	return ionSystems
}

func getDistance(atom1 *atom, pos []float64) float64 {
	dx2 := math.Pow(atom1.pos[0] - pos[0], 2)
	dy2 := math.Pow(atom1.pos[1] - pos[1], 2)
	dz2 := math.Pow(atom1.pos[2] - pos[2], 2)

	return math.Sqrt(dx2 + dy2 + dz2)
}


