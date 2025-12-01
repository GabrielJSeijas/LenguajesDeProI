#include <iostream>
#include <string>
#include <vector>
#include <map>
#include <algorithm>
#include <numeric>

// este archivo define estructuras y funciones
// para registrar tipos (atómicos, struct, union) y calcular tamaños
// y alineaciones usando diferentes estrategias de empaquetado.

// Definición de un tipo: nombre, tamaño y alineación almacenados,
// los nombres de los tipos que lo componen (para struct/union) y
// la categoría del tipo ("ATOMIC", "STRUCT", "UNION").
struct TypeDefinition {
    std::string name ;
    size_t size;
    size_t alignment;
    std::vector<std::string> composition;
    std::string type_class ;
};

class TypeManager {
private:
    // Mapa que almacena las definiciones de tipos por nombre.
    std::map<std::string, TypeDefinition> defined_types;

    // Calcula el siguiente offset alineado a 'alignment' desde
    // 'current_offset'. Se usa para simular padding entre campos.
    size_t align_up(size_t current_offset, size_t alignment) {
        return (current_offset + alignment - 1) & ~(alignment - 1);
    }

    // Busca una definición de tipo en el mapa y lanza excepción si no existe.
    const TypeDefinition& find_type(const std::string& type_name) {
        if (defined_types.find(type_name) == defined_types.end()) {
            throw std::runtime_error("Error: Tipo '" + type_name + "' no definido.");
        }
        return defined_types.at(type_name);
    }

public:
    // Registra un tipo atómico con su tamaño y alineación.
    void execute_atomic(const std::string& name, size_t representation, size_t alignment);
    // Registra un struct compuesto por una lista de tipos (por nombre).
    void execute_struct(const std::string& name, const std::vector<std::string>& types);
    // Registra una unión compuesta por una lista de tipos (por nombre).
    void execute_union(const std::string& name, const std::vector<std::string>& types);
    // Muestra la descripción y los cálculos para un tipo dado.
    void execute_describe(const std::string& name);

    // Resultado de un cálculo: tamaño total, requisito de alineación y espacio desperdiciado.
    struct Result { size_t size; size_t alignment; size_t wasted; };
    
    // Estrategia 1: calcular tamaño considerando padding entre campos y al final para cumplir la alineación del struct.
    Result calculate_no_packing(const TypeDefinition& struct_type);
    
    // Estrategia 2: empaquetado sin padding; el tamaño es la suma de tamaños de campos.
    Result calculate_packed(const TypeDefinition& struct_type);

    // Estrategia 3: reordenar campos para minimizar el padding y calcular el tamaño resultante.
    Result calculate_optimal_packing(const TypeDefinition& struct_type);

    // Imprime los resultados de una estrategia con un título.
    void print_result(const std::string& strategy, const Result& res);
};