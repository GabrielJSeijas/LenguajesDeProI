// Gabriel Seijas 19-00036
package main

import (
	"errors"
	"fmt"
)

// TransponerMatriz toma una matriz cuadrada A de tamaño N x N y devuelve su transpuesta.
// Si la matriz no es cuadrada, retorna un error.
func TransponerMatriz(A [][]float64) ([][]float64, error) {
	n := len(A)
	if n == 0 {
		return nil, errors.New("la matriz no puede estar vacía")
	}

	// Verifico que todas las filas tengan la misma longitud que n
	for i := 0; i < n; i++ {
		if len(A[i]) != n {
			return nil, errors.New("la matriz debe ser cuadrada")
		}
	}

	// Creo la matriz transpuesta con las mismas dimensiones
	AT := make([][]float64, n)
	for i := range AT {
		AT[i] = make([]float64, n)
	}

	// Intercambio filas por columnas
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			AT[i][j] = A[j][i]
		}
	}

	return AT, nil
}

// MultiplicarMatrices toma dos matrices cuadradas A y B de tamaño N x N y devuelve su producto C = A x B.
// Si las matrices no son cuadradas o no tienen el mismo tamaño, retorna un error.
func MultiplicarMatrices(A, B [][]float64) ([][]float64, error) {
	n := len(A)
	if n == 0 || len(B) == 0 {
		return nil, errors.New("las matrices no pueden estar vacías")
	}

	// Verifico que ambas matrices sean cuadradas y del mismo tamaño
	for i := 0; i < n; i++ {
		if len(A[i]) != n || len(B[i]) != n {
			return nil, errors.New("ambas matrices deben ser cuadradas y de la misma dimensión N x N")
		}
	}

	// Creo la matriz resultado C
	C := make([][]float64, n)
	for i := range C {
		C[i] = make([]float64, n)
	}

	// Multiplico las matrices usando la fórmula clásica
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				C[i][j] += A[i][k] * B[k][j]
			}
		}
	}

	return C, nil
}

// ProductoMatrizPorTranspuesta calcula el producto de una matriz cuadrada A por su transpuesta.
// Devuelve la matriz resultado o un error si ocurre algún problema.
func ProductoMatrizPorTranspuesta(A [][]float64) ([][]float64, error) {
	// Primero obtengo la transpuesta de A
	AT, err := TransponerMatriz(A)
	if err != nil {
		return nil, fmt.Errorf("error al transponer la matriz: %w", err)
	}

	// Multiplico A por su transpuesta
	resultado, err := MultiplicarMatrices(A, AT)
	if err != nil {
		return nil, fmt.Errorf("error al multiplicar la matriz por su transpuesta: %w", err)
	}

	return resultado, nil
}

// ImprimirMatriz imprime una matriz de números flotantes en consola.
// Si la matriz está vacía, imprime solo corchetes.
func ImprimirMatriz(m [][]float64) {
	if len(m) == 0 {
		fmt.Println("[]")
		return
	}
	for _, row := range m {
		fmt.Println(row)
	}
}
