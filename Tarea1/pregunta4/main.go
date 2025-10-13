// gabriel Seijas 19-00036
package main

import (
	"fmt"
	"pregunta4/vector3d"
)

func main() {
	// Definición de tres vectores 3D
	a := vector3d.Vector3D{X: 1, Y: 2, Z: 3}
	b := vector3d.Vector3D{X: 4, Y: 5, Z: 6}
	c := vector3d.Vector3D{X: 7, Y: 8, Z: 9}

	fmt.Println("--- Vectores Iniciales ---")
	fmt.Println("Vector a:", a)
	fmt.Println("Vector b:", b)
	fmt.Println("Vector c:", c)
	fmt.Println("--------------------------")

	fmt.Println(" Ejemplos de Operaciones con Vectores 3D:")

	// Suma de dos vectores
	res1 := b.Add(c)
	fmt.Println("b + c =", res1)

	// Producto cruz entre a y b, luego suma con c
	crossAB := a.CrossProduct(b)
	res2 := crossAB.Add(c)
	fmt.Println("a x b + c =", res2)

	// Suma de b consigo mismo, resta de c y a, y producto cruz entre los resultados
	sumBB := b.Add(b)
	subCA := c.Subtract(a)
	res3 := sumBB.CrossProduct(subCA)
	fmt.Println("(b + b) x (c - a) =", res3)

	// Producto cruz entre c y b, luego producto punto con a
	crossCB := c.CrossProduct(b)
	res4 := a.DotProduct(crossCB)
	fmt.Println("a . (c x b) =", res4)

	fmt.Println("\n Operaciones con Escalares ")

	// Suma de un escalar a cada componente de b
	res5 := b.AddScalar(3.0)
	fmt.Println("b + 3 =", res5)

	// Multiplicación escalar de a y suma con b
	scaledA := a.Scale(3.0)
	res6 := scaledA.Add(b)
	fmt.Println("a * 3.0 + b =", res6)

	// Suma de b consigo mismo, producto punto entre c y a, y multiplicación escalar
	sumBB_scalar := b.Add(b)
	dotCA_scalar := c.DotProduct(a)
	res7 := sumBB_scalar.Scale(dotCA_scalar)
	fmt.Println("(b + b) * (c . a) =", res7)
}
