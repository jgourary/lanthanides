package main

var radialPath = "C:\\sandia_lib\\radial.exe"

// Program begins here
func main() {
	// charmm2amoebaBatch("C:\\Users\\jtgou\\go\\src\\bilayerProcessorManualFinal\\input", "C:\\Users\\jtgou\\go\\src\\bilayerProcessorManualFinal\\output")

	path := "C:\\sandia_lib\\kcl.arc"
	ion := "K+"
	count := 10
	outDir := "C:\\sandia_lib\\k+"
	file2systems(path, ion, count, outDir)

	ion = "Cl-"
	outDir = "C:\\sandia_lib\\cl-"
	file2systems(path, ion, count, outDir)
}