// Gabriel Seijas 19-00036
package main

import (
	"fmt"
)

func main() {
	// Pruebas para la función RotarString
	fmt.Println(" Ejemplos de Rotación de Cadenas ")
	fmt.Println("rotar(\"hola\", 0) =", RotarString("hola", 0)) // No rota nada
	fmt.Println("rotar(\"hola\", 1) =", RotarString("hola", 1)) // Rota una posición
	fmt.Println("rotar(\"hola\", 2) =", RotarString("hola", 2)) // Rota dos posiciones
	// Más ejemplos de rotación
	fmt.Println("rotar(\"palabra\", 9) =", RotarString("palabra", 9)) // Rota más que la longitud
	fmt.Println("\n---------------------------------------")

	// Pruebas para la función ProductoMatrizPorTranspuesta
	fmt.Println(" Ejemplos de Producto Matriz por su Transpuesta")

	// Matriz 2x2 para probar
	matrizA2x2 := [][]float64{
		{1, 2},
		{3, 4},
	}
	fmt.Println("Matriz A (2x2):")
	ImprimirMatriz(matrizA2x2)

	resultado2x2, err := ProductoMatrizPorTranspuesta(matrizA2x2)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Println("Producto A x A^T (2x2):")
		ImprimirMatriz(resultado2x2)
	}
	fmt.Println()

	// Matriz 3x3 para probar
	matrizB3x3 := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	fmt.Println("Matriz B (3x3):")
	ImprimirMatriz(matrizB3x3)

	resultado3x3, err := ProductoMatrizPorTranspuesta(matrizB3x3)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Println("Producto B x B^T (3x3):")
		ImprimirMatriz(resultado3x3)
		// El resultado debería ser:
		// [[14 32 50]
		//  [32 77 122]
		//  [50 122 194]]
	}
	fmt.Println()

	// Prueba con matriz que no es cuadrada (debería dar error)
	matrizNoCuadrada := [][]float64{
		{1, 2, 3},
		{4, 5},
	}
	fmt.Println("Matriz No Cuadrada:")
	ImprimirMatriz(matrizNoCuadrada)
	_, err = ProductoMatrizPorTranspuesta(matrizNoCuadrada)
	if err != nil {
		fmt.Printf("Esperado error: %s\n", err)
	}
	fmt.Println("\n---------------------------------------")
}
