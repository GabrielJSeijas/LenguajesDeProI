// Función que calcula el número de pasos almacenados en contador para alcanzar 1
func calculateCollatzCount(n: Int) -> Int {
    // Caso base si el número es 1, ya no hay más pasos
    if n <= 0 {
        print("El número debe ser positivo.")
        return 0
    }
    
    var numeroActual = n // El número que vamos a modificar en cada paso
    var contador = 0         // El contador de aplicaciones de f(n)
    
    // Repetimos mientras el número actual no sea 1
    while numeroActual != 1 {
        // La fórmula f(n)
        if numeroActual % 2 == 0 {
            // Caso par: n / 2
            numeroActual = numeroActual / 2
        } else {
            // Caso impar: 3n + 1
            numeroActual = 3 * numeroActual + 1
        }
        
        // Aumentamos el contador en cada paso
        contador += 1
    }
    
    //Devolvemos el total de pasos
    return count
}

// Ejemplo de uso (n = 42)
let numeroN = 42
let totalContador = calculateCollatzCount(n: numeroN)
print("El contador(\(numeroN)) es: \(totalContador)") 
