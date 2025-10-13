// Gabriel Seijas 19-00036
package vector3d_test

import (
	"math"
	"pregunta4/vector3d"
	"testing"
)

const epsilon = 1e-9 // Tolerancia para comparar números flotantes

// areFloatsEqual compara dos valores float64 considerando una pequeña tolerancia.
func areFloatsEqual(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

// Prueba la suma de dos vectores 3D.
func TestVector3D_Add(t *testing.T) {
	v1 := vector3d.Vector3D{X: 1, Y: 2, Z: 3}
	v2 := vector3d.Vector3D{X: 4, Y: 5, Z: 6}
	expected := vector3d.Vector3D{X: 5, Y: 7, Z: 9}
	result := v1.Add(v2)

	if !areFloatsEqual(result.X, expected.X) || !areFloatsEqual(result.Y, expected.Y) || !areFloatsEqual(result.Z, expected.Z) {
		t.Errorf("Add falló: esperado %v, obtenido %v", expected, result)
	}
}

// Prueba la resta de dos vectores 3D.
func TestVector3D_Subtract(t *testing.T) {
	v1 := vector3d.Vector3D{X: 5, Y: 7, Z: 9}
	v2 := vector3d.Vector3D{X: 4, Y: 5, Z: 6}
	expected := vector3d.Vector3D{X: 1, Y: 2, Z: 3}
	result := v1.Subtract(v2)

	if !areFloatsEqual(result.X, expected.X) || !areFloatsEqual(result.Y, expected.Y) || !areFloatsEqual(result.Z, expected.Z) {
		t.Errorf("Subtract falló: esperado %v, obtenido %v", expected, result)
	}
}

// Prueba el producto cruzado entre dos vectores 3D.
func TestVector3D_CrossProduct(t *testing.T) {
	v1 := vector3d.Vector3D{X: 1, Y: 0, Z: 0}       // Vector i
	v2 := vector3d.Vector3D{X: 0, Y: 1, Z: 0}       // Vector j
	expected := vector3d.Vector3D{X: 0, Y: 0, Z: 1} // Vector k
	result := v1.CrossProduct(v2)

	if !areFloatsEqual(result.X, expected.X) || !areFloatsEqual(result.Y, expected.Y) || !areFloatsEqual(result.Z, expected.Z) {
		t.Errorf("CrossProduct falló para i x j: esperado %v, obtenido %v", expected, result)
	}

	v3 := vector3d.Vector3D{X: 2, Y: 3, Z: 4}
	v4 := vector3d.Vector3D{X: 5, Y: 6, Z: 7}
	expected2 := vector3d.Vector3D{X: -3, Y: 6, Z: -3}
	result2 := v3.CrossProduct(v4)

	if !areFloatsEqual(result2.X, expected2.X) || !areFloatsEqual(result2.Y, expected2.Y) || !areFloatsEqual(result2.Z, expected2.Z) {
		t.Errorf("CrossProduct falló para vectores generales: esperado %v, obtenido %v", expected2, result2)
	}
}

// Prueba el producto punto entre dos vectores 3D.
func TestVector3D_DotProduct(t *testing.T) {
	v1 := vector3d.Vector3D{X: 1, Y: 2, Z: 3}
	v2 := vector3d.Vector3D{X: 4, Y: 5, Z: 6}
	expected := 32.0 // (1*4) + (2*5) + (3*6) = 32
	result := v1.DotProduct(v2)

	if !areFloatsEqual(result, expected) {
		t.Errorf("DotProduct falló: esperado %f, obtenido %f", expected, result)
	}
}

// Prueba el cálculo de la magnitud de un vector 3D.
func TestVector3D_Magnitude(t *testing.T) {
	v := vector3d.Vector3D{X: 3, Y: 4, Z: 0}
	expected := 5.0 // sqrt(3^2 + 4^2 + 0^2) = 5
	result := v.Magnitude()

	if !areFloatsEqual(result, expected) {
		t.Errorf("Magnitude falló: esperado %f, obtenido %f", expected, result)
	}

	v2 := vector3d.Vector3D{X: 0, Y: 0, Z: 0}
	expected2 := 0.0
	result2 := v2.Magnitude()
	if !areFloatsEqual(result2, expected2) {
		t.Errorf("Magnitude falló para vector cero: esperado %f, obtenido %f", expected2, result2)
	}
}

// Prueba la multiplicación de un vector por un escalar.
func TestVector3D_Scale(t *testing.T) {
	v := vector3d.Vector3D{X: 1, Y: 2, Z: 3}
	scalar := 2.0
	expected := vector3d.Vector3D{X: 2, Y: 4, Z: 6}
	result := v.Scale(scalar)

	if !areFloatsEqual(result.X, expected.X) || !areFloatsEqual(result.Y, expected.Y) || !areFloatsEqual(result.Z, expected.Z) {
		t.Errorf("Scale falló: esperado %v, obtenido %v", expected, result)
	}
}

// Prueba la suma de un escalar a cada componente de un vector.
func TestVector3D_AddScalar(t *testing.T) {
	v := vector3d.Vector3D{X: 1, Y: 2, Z: 3}
	scalar := 10.0
	expected := vector3d.Vector3D{X: 11, Y: 12, Z: 13}
	result := v.AddScalar(scalar)

	if !areFloatsEqual(result.X, expected.X) || !areFloatsEqual(result.Y, expected.Y) || !areFloatsEqual(result.Z, expected.Z) {
		t.Errorf("AddScalar falló: esperado %v, obtenido %v", expected, result)
	}
}
