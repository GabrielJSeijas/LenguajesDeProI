import Foundation

// Constantes del problema: parámetros alpha, beta y el tamaño del caso base.
let ALPHA = 6
let BETA = 7
let casoBase = ALPHA * BETA // 42

// Calcula F_alpha,beta(n) con la traducción directa de la definición recursiva.
// Explicación (estudiante): uso la definición tal cual; si n < 42 devuelvo n.
// Para n >= 42 sumo F(n - 7*i) para i = 1..6.
func fDirecta(n: Int) -> Int {
    if n < casoBase {
        return n
    }
    
    var sumatoria = 0
    for i in 1...ALPHA {
        let argumento = n - (BETA * i) // n - 7*i
        sumatoria += fDirecta(n: argumento)
    }
    
    return sumatoria
}

// Implementación que intenta usar recursión de cola manteniendo una ventana con los últimos 42 valores necesarios.
// guardo los valores F(0..41) en un arreglo y voy generando F(42), F(43), ... actualizando la ventana como una cola.
func fCola(n: Int) -> Int {
    if n < casoBase {
        return n
    }
    
    // Ventana inicial con F(0..41) = 0..41
    let valoresRecientes: [Int] = Array(0..<casoBase)
    
    // Función auxiliar que simula la iteración mediante recursión de cola.
    func tailRecursion(actual: Int, ventana: [Int]) -> Int {
        // Si llegamos al objetivo, el resultado estará en el último element de la ventana (por cómo se actualiza la ventana en cada paso).
        if actual == n {
            return ventana.last!
        }
        
        // Calculo F(actual) sumando los elementos desplazados por 7*i
        var siguienteValor = 0
        for i in 1...ALPHA {
            let indice = casoBase - (BETA * i)
            siguienteValor += ventana[indice]
        }
        
        // Actualizo la ventana: quito el más antiguo y añado el nuevo
        var nuevaVentana = ventana
        nuevaVentana.removeFirst()
        nuevaVentana.append(siguienteValor)
        
        return tailRecursion(actual: actual + 1, ventana: nuevaVentana)
    }
    
    return tailRecursion(actual: casoBase , ventana: valoresRecientes)
}

// Enfoque iterativo usando programación dinámica y una ventana deslizante.
// es la implementación más eficiente de las subrutinas mantengo sólo los últimos 42 valores y avanzo de 42 hasta n calculando cada nuevo F con O(1) trabajo (sumando 6 entradas de la ventana).
func fIterativa(n: Int) -> Int {
    if n < casoBase {
        return n
    }
    
    var valoresRecientes: [Int] = Array(0..<casoBase)
    
    // Itero desde CASO_BASE hasta n, actualizando la ventana cada vez.
    for _ in casoBase...n {
        var siguienteValor = 0
        for i in 1...ALPHA {
            let indice = casoBase - (BETA * i)
            siguienteValor += valoresRecientes[indice]
        }
        
        valoresRecientes.removeFirst()
        valoresRecientes.append(siguienteValor)
    }
    
    return valoresRecientes.last!
}
// Medimos el tiempo de ejecución de un cierre que devuelve Int.
//  uso clock() y convierto la diferencia a segundos
func measureTime(closure: () -> Int) -> Double {
    let start = clock() // Inicia el reloj
    let _ = closure() // Ejecuta el cierre
    let end = clock() // Detiene el reloj
    
    let elapsed = Double(end - start) / Double(CLOCKS_PER_SEC)
    return elapsed
}

let valoresenN: [Int] = [42,43,50,100,1000,5000,10000]

print("Analisis comparativo")
print("N | Tiempo Directa (segundos) | Tiempo Cola (segundos) | Tiempo Iterativa (segundos)")
var resultados: [(n: Int, directa: Double, cola: Double, iterativa: Double)] = []

for n in valoresenN {
    print("Calculando para N=\(n)...")
    
    var tiempoDirecta: Double = -1.0
    if n <= 45 { // limitamos la versión directa para que no tarde demasiado
        tiempoDirecta = measureTime {
            return fDirecta(n: n)
        }
    }

    let tiempoCola: Double = measureTime {
        return fCola(n: n)
    }

    let tiempoIterativa: Double = measureTime {
        return fIterativa(n: n)
    }
    
    let recDirecta = tiempoDirecta >= 0 ? String(format: "%.6f", tiempoDirecta) : "N/A"
    let recCola = String(format: "%.6f", tiempoCola)
    let recIterativa = String(format: "%.6f", tiempoIterativa)

    print("\(n) | \(recDirecta) | \(recCola) | \(recIterativa)")

    resultados.append((n: n, directa: tiempoDirecta, cola: tiempoCola, iterativa: tiempoIterativa))
}


