#include <iostream>
#include <string>
#include <vector>
#include <map>
#include <algorithm>
#include <numeric>

// --- 1. Definición Base del Tipo ---
struct TypeDefinition {
    std::string name;
    // El tamaño y la alineación ALMACENADOS.
    size_t size; 
    size_t alignment;
    std::vector<std::string> composition; // Nombres de los tipos que lo componen (para STRUCT/UNION)
    std::string type_class; // "ATOMIC", "STRUCT", "UNION"
};

class TypeManager {
private:
    std::map<std::string, TypeDefinition> defined_types;

    // --- Funciones de Ayuda para el cálculo de alineación ---
    
    // Función de alineación clave: calcula el próximo múltiplo
    // Alignment es el requisito de alineación del campo actual.
    size_t align_up(size_t current_offset, size_t alignment) {
        // Fórmula: (offset + alignment - 1) / alignment * alignment
        return (current_offset + alignment - 1) & ~(alignment - 1); 
    }

    // Encuentra el tipo en el mapa
    const TypeDefinition& find_type(const std::string& type_name) {
        if (defined_types.find(type_name) == defined_types.end()) {
            throw std::runtime_error("Error: Tipo '" + type_name + "' no definido.");
        }
        return defined_types.at(type_name);
    }

public:
    // --- 2. Implementación de los Comandos ---

    void execute_atomic(const std::string& name, size_t representation, size_t alignment);
    void execute_struct(const std::string& name, const std::vector<std::string>& types);
    void execute_union(const std::string& name, const std::vector<std::string>& types);
    void execute_describe(const std::string& name);

    // --- 3. Funciones de Cálculo de Estrategias ---

    // a) Sin empaquetar (Estrategia estándar con padding)
    struct Result { size_t size; size_t alignment; size_t wasted; };
    Result calculate_no_packing(const TypeDefinition& struct_type);
    
    // b) Empaquetado (No hay padding, solo suma de tamaños)
    Result calculate_packed(const TypeDefinition& struct_type);

    // c) Reordenando los campos (Optimización para minimizar padding)
    Result calculate_optimal_packing(const TypeDefinition& struct_type);

    // Utilidad para imprimir resultados
    void print_result(const std::string& strategy, const Result& res);
};