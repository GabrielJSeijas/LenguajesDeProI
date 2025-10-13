//Gabriel Seijas 19-00036
package vector3d

import "math"

// Estructura para representar un vector en 3D
type Vector3D struct {
	X, Y, Z float64
}

// Suma dos vectores y regresa el resultado
func (v Vector3D) Add(other Vector3D) Vector3D {
	return Vector3D{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

// Resta otro vector al actual y regresa el resultado
func (v Vector3D) Subtract(other Vector3D) Vector3D {
	return Vector3D{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

// Calcula el producto cruz entre dos vectores
func (v Vector3D) CrossProduct(other Vector3D) Vector3D {
	return Vector3D{
		v.Y*other.Z - v.Z*other.Y,
		v.Z*other.X - v.X*other.Z,
		v.X*other.Y - v.Y*other.X,
	}
}

// Calcula el producto punto entre dos vectores
func (v Vector3D) DotProduct(other Vector3D) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Calcula la magnitud del vector
func (v Vector3D) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Multiplica el vector por un escalar
func (v Vector3D) Scale(scalar float64) Vector3D {
	return Vector3D{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

// Suma un escalar a cada componente del vector
func (v Vector3D) AddScalar(n float64) Vector3D {
	return Vector3D{v.X + n, v.Y + n, v.Z + n}
}
