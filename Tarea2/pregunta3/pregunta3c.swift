struct OrdenadorIterador: IteratorProtocol {
    //La lista mutable 
    private var listaMutable: [Int]

    // Inicializador
    init(listaOriginal: [Int]) {
        self.listaMutable = listaOriginal
    }

    // El corazón del iterador calcula y devuelve el siguiente elemento
    mutating func next() -> Int? {
        // continúa mientras queden elementos
        guard !listaMutable.isEmpty else {
            return nil 
        }

        // Encontramos el elemento mínimo en la lista restante
        var indiceMinimo = 0
        var valorMinimo = listaMutable[0]

        for i in 1..<listaMutable.count {
            if listaMutable[i] < valorMinimo {
                valorMinimo = listaMutable[i]
                indiceMinimo = i
            }
        }

        // Eliminamos el elemento mínimo de la lista mutable
        listaMutable.remove(at: indiceMinimo)

        // Devolvemos el elemento mínimo encontrado 
        return valorMinimo
    }
}

struct ListaOrdenada: Sequence {
    private let lista: [Int]

    init(listaOriginal: [Int]) {
        self.lista = listaOriginal
    }

    // definimos el iterador a usar para recorrer la secuencia
    func makeIterator() -> OrdenadorIterador {
        return OrdenadorIterador(listaOriginal: self.lista)
    }
}
let listaDesordenada = [1, 3, 3, 2, 1]

// Creamos la secuencia ordenada (el iterador se crea internamente)
let secuenciaOrdenada = ListaOrdenada(listaOriginal: listaDesordenada)

print("La lista es: \(listaDesordenada)")
print("Elementos en ordenados en orden ascendente:")
for elemento in secuenciaOrdenada {
    print(elemento, terminator: " ")
}
// Salida: 1 1 2 3 3 