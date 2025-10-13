// Gabriel Seijas 19-00036
package main

// RotarString rota la cadena w a la izquierda k veces usando recursión.
func RotarString(w string, k int) string {
	// Si k es 0, no hay que rotar nada
	if k == 0 {
		return w
	}

	// Si la cadena está vacía, no se puede rotar
	if len(w) == 0 {
		return w
	}

	// Si k es mayor que el largo de la cadena, lo ajusto
	k = k % len(w)
	if k == 0 {
		return w
	}

	// Tomo el primer caracter y el resto de la cadena
	primerCaracter := string(w[0])
	restoCadena := w[1:]

	// Roto la cadena una vez y llamo recursivamente
	cadenaRotadaUnaPosicion := restoCadena + primerCaracter
	return RotarString(cadenaRotadaUnaPosicion, k-1)
}
