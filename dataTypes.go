package main

type atom struct {
	element string
	id int
	isIon bool

	pos []float64

	aminoAcid string
	residue string
}

type ionSystem struct {
	center []float64
	charge int
	atoms map[int]*atom
	residueList map[string]string
}

func copyAtom(thisAtom *atom) atom {
	var newAtom atom

	newAtom.element = thisAtom.element
	newAtom.id = thisAtom.id
	newAtom.isIon = thisAtom.isIon
	newAtom.pos = thisAtom.pos
	newAtom.residue = thisAtom.residue
	newAtom.aminoAcid = thisAtom.aminoAcid

	return newAtom
}
