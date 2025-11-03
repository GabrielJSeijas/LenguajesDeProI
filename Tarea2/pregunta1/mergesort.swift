// Fusiona dos arrays ordenados en uno solo ordenado.
// Implementación basada en dos índices que recorren los arrays de entrada.
func merge(left: [Int], right: [Int]) -> [Int] {
    var mergedArray: [Int] = []
    var leftIndex = 0
    var rightIndex = 0

    // Mientras haya elementos en ambos arrays seleccionamos el menor 
    while leftIndex < left.count && rightIndex < right.count {
        if left[leftIndex] < right[rightIndex] {
            mergedArray.append(left[leftIndex])
            leftIndex += 1
        } else {
            mergedArray.append(right[rightIndex])
            rightIndex += 1
        }
    }

    // Añadir cualquier resto restante,  solo uno de los dos puede tener elementos
    // Comprobaciones para evitar índices fuera de rango al crear particiones 
    if leftIndex < left.count {
        mergedArray.append(contentsOf: left[leftIndex...])
    }
    if rightIndex < right.count {
        mergedArray.append(contentsOf: right[rightIndex...])
    }

    return mergedArray
}

// Ordenación por mezcla (merge sort) usando divide and conquer.
// Caso base: arrays de tamaño 0 o 1 ya están ordenados divide el array en dos mitades, ordena recursivamente y fusiona.
// Se crean subarrays con Array(array[...]) en cada llamada, lo que implica copias adicionales.

func mergeSort(array: [Int]) -> [Int] {
    guard array.count > 1 else { return array }

    let middleIndex = array.count / 2

    // Creamos subarrays para las mitades 
    let leftHalf = mergeSort(array: Array(array[0..<middleIndex]))
    let rightHalf = mergeSort(array: Array(array[middleIndex..<array.count]))

    return merge(left: leftHalf, right: rightHalf)
}

let unsortedArray = [12, 5, 1, 9, 4, 8]
let sortedArray = mergeSort(array: unsortedArray)
print("Array Desordenado: \(unsortedArray)")
print("Array Ordenado: \(sortedArray)")
// Salida: [1, 4, 5, 8, 9, 12]