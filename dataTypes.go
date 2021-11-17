package main

type atom struct {
	element string
	id int
	isIon bool

	pos []float64

	residue string
}

type ionSystem struct {
	center []float64
	charge int
	atoms map[int]*atom
}

func copyAtom(thisAtom *atom) atom {
	var newAtom atom

	newAtom.element = thisAtom.element
	newAtom.id = thisAtom.id
	newAtom.isIon = thisAtom.isIon
	newAtom.pos = thisAtom.pos
	newAtom.residue = thisAtom.residue

	return newAtom
}
