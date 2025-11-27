#include "TypeManager.h"

// --- 2. Implementación de los Comandos de Definición ---

void TypeManager::execute_atomic(const std::string& name, size_t representation, size_t alignment) {
    defined_types[name] = {
        name,
        representation,
        alignment,
        {},
        "ATOMIC"
    };
    std::cout << "Tipo ATOMICO '" << name << "' definido (Tamaño: " << representation << ", Alineación: " << alignment << ").\n";
}

void TypeManager::execute_struct(const std::string& name, const std::vector<std::string>& types) {
    defined_types[name] = {
        name,
        0, // El tamaño se calculará al describirlo
        0, // La alineación se calculará al describirlo
        types,
        "STRUCT"
    };
    std::cout << "Tipo STRUCT '" << name << "' definido.\n";
}

void TypeManager::execute_union(const std::string& name, const std::vector<std::string>& types) {
    size_t max_size = 0;
    size_t max_alignment = 0;

    for (const auto& type_name : types) {
        const auto& component = find_type(type_name);
        max_size = std::max(max_size, component.size);
        max_alignment = std::max(max_alignment, component.alignment);
    }
    
    // El tamaño final de UNION debe ser un múltiplo de su alineación máxima (padding al final)
    size_t final_size = align_up(max_size, max_alignment);
    size_t wasted_bytes = final_size - max_size; // Padding al final

    defined_types[name] = {
        name,
        final_size,
        max_alignment,
        types,
        "UNION"
    };
    std::cout << "Tipo UNION '" << name << "' definido (Tamaño: " << final_size << ", Alineación: " << max_alignment << ").\n";
}

// --- 3. Funciones de Cálculo de Estrategias ---

// a) Sin Empaquetar (Estándar: Padding en campos + Padding final)
TypeManager::Result TypeManager::calculate_no_packing(const TypeDefinition& struct_type) {
    size_t current_offset = 0;
    size_t max_alignment = 0;
    size_t padding_internal = 0;

    for (const auto& type_name : struct_type.composition) {
        const auto& component = find_type(type_name);
        max_alignment = std::max(max_alignment, component.alignment);
        
        // Calcular padding antes de este campo para alinearlo
        size_t next_offset = align_up(current_offset, component.alignment);
        padding_internal += (next_offset - current_offset);
        
        current_offset = next_offset + component.size;
    }

    // Padding final: El tamaño total debe ser múltiplo de la alineación máxima
    size_t final_size = align_up(current_offset, max_alignment);
    size_t padding_final = final_size - current_offset;

    return {
        final_size,
        max_alignment,
        padding_internal + padding_final
    };
}

// b) Empaquetado (Packed: Sin padding, solo suma)
TypeManager::Result TypeManager::calculate_packed(const TypeDefinition& struct_type) {
    size_t total_size = 0;
    size_t max_alignment = 0;

    for (const auto& type_name : struct_type.composition) {
        const auto& component = find_type(type_name);
        
        size_t current_size = component.size;

        // Si el componente es un STRUCT, su tamaño empaquetado debe calcularse recursivamente.
        if (component.type_class == "STRUCT") {
            // El tamaño empaquetado de un STRUCT anidado es la suma estricta de sus propios campos.
            // Es necesario llamar a calculate_packed recursivamente para este componente anidado.
            // Esto asegura que Pair, que es un STRUCT, se calcule como 6 bytes (4+2), NO 8.
            current_size = calculate_packed(component).size;
        }

        total_size += current_size;
        max_alignment = std::max(max_alignment, component.alignment);
    }
    
    // ...
    return {
        total_size, // Ahora será 6 (Pair) + 1 (char) + 16 (byte_16) = 23 bytes
        max_alignment,
        0
    };
}

// c) Reordenando (Optimal: Ordenar por alineación descendente)
TypeManager::Result TypeManager::calculate_optimal_packing(const TypeDefinition& struct_type) {
    // 1. Obtener y preparar los tipos componentes con sus ALINEACIONES REALES
    // Usamos pair<size_t, size_t> para almacenar {Alignment, Size} para la clasificación
    // Almacenamos alineación primero porque es la clave de clasificación.
    std::vector<std::pair<size_t, size_t>> components_info;

    for (const auto& type_name : struct_type.composition) {
        const auto& defined_type = find_type(type_name);
        
        size_t current_size = defined_type.size;
        size_t current_align = defined_type.alignment;

        // Si el componente es un STRUCT anidado, debemos calcular su tamaño "no empaquetado"
        if (defined_type.type_class == "STRUCT") {
            Result res = calculate_no_packing(defined_type);
            current_size = res.size;
            current_align = res.alignment;
        }

        // Corregido: Inicializa std::pair con current_align y current_size
        components_info.push_back({current_align, current_size});
    }

    // 2. Ordenar por ALINEACIÓN DESCENDENTE (heurística óptima)
    std::sort(components_info.begin(), components_info.end(), 
        [](const auto& a, const auto& b) {
            // El primer elemento del pair (a.first) es la alineación
            return a.first > b.first;
        });

    // 3. Aplicar el cálculo de NO PACKING al ordenado
    size_t current_offset = 0;
    size_t max_alignment = 0;
    size_t padding_internal = 0;

    for (const auto& component : components_info) {
        // El primer elemento del pair (component.first) es la alineación
        max_alignment = std::max(max_alignment, component.first);
        
        // Calcular padding antes de este campo para alinearlo
        size_t component_alignment = component.first;// Corregido: Usa la alineación del pair
        size_t next_offset = align_up(current_offset, component_alignment);
        padding_internal += (next_offset - current_offset);
        
        // El segundo elemento del pair (component.second) es el tamaño
        current_offset = next_offset + component.second;
    }

    // Padding final
    size_t final_size = align_up(current_offset, max_alignment);
    size_t padding_final = final_size - current_offset;

    return {
        final_size,
        max_alignment,
        padding_internal + padding_final
    };
}

// Utilidad para imprimir resultados
void TypeManager::print_result(const std::string& strategy, const Result& res) {
    std::cout << "  - Estrategia " << strategy << ":\n";
    std::cout << "    * Tamaño Total: " << res.size << " bytes\n";
    std::cout << "    * Alineación Requerida: " << res.alignment << " bytes\n";
    std::cout << "    * Bytes Desperdiciados (Padding): " << res.wasted << " bytes\n";
}

// Comando DESCRIBIR
void TypeManager::execute_describe(const std::string& name) {
    try {
        const auto& type = find_type(name);
        std::cout << "\n--- DESCRIPCIÓN DEL TIPO '" << name << "' ---\n";
        std::cout << "Clase de Tipo: " << type.type_class << "\n";

        if (type.type_class == "ATOMIC") {
            std::cout << "  * Tamaño: " << type.size << " bytes\n";
            std::cout << "  * Alineacion: " << type.alignment << " bytes\n";
            std::cout << "  * Bytes Desperdiciados: 0\n";
        } else if (type.type_class == "UNION") {
            std::cout << "  * Tamaño: " << type.size << " bytes\n";
            std::cout << "  * Alineacion: " << type.alignment << " bytes\n";
            // Para la unión, el desperdicio es el padding final
            size_t max_size_component = 0;
            for (const auto& comp_name : type.composition) {
                max_size_component = std::max(max_size_component, find_type(comp_name).size);
            }
            std::cout << "  * Bytes Desperdiciados: " << (type.size - max_size_component) << " (padding final)\n";
        } else if (type.type_class == "STRUCT") {
            // Sin empaquetar
            Result res_no_packing = calculate_no_packing(type);
            print_result("Sin empaquetar", res_no_packing); // [cite: 65]

            // Empaquetado
            Result res_packed = calculate_packed(type);
            print_result("Empaquetado", res_packed); // [cite: 66]

            // Reordenando
            Result res_optimal = calculate_optimal_packing(type);
            print_result("Reordenando de manera optima", res_optimal); // [cite: 67]
        }
        std::cout << "---------------------------------------\n";

    } catch (const std::runtime_error& e) {
        std::cerr << e.what() << "\n";
    }
}