#include "TypeManager.h"


// Implementación de TypeManager:
// - Comandos para definir tipos (ATOMIC, STRUCT, UNION).
// - Cálculos de tamaño y alineación según tres estrategias.
// - Comando DESCRIBIR para mostrar información y comparativas.

void TypeManager::execute_atomic(const std::string& name, size_t representation, size_t alignment) {
    // Define un tipo atómico: guardamos su nombre, tamaño y alineación.
    // Los tipos atómicos no tienen composición.
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
    // Define un tipo STRUCT: aquí solo almacenamos la composición.
    // El tamaño y alineación reales se calculan cuando se describe el tipo,
    // porque pueden depender de tipos anidados y de la estrategia elegida.

    defined_types[name] = {
        name,
        0,
        0,
        types,
        "STRUCT"
    };
    std::cout << "Tipo STRUCT '" << name << "' definido.\n";

}

void TypeManager::execute_union(const std::string& name,const std::vector<std::string>& types) {
    // Define un tipo UNION: el tamaño es el máximo de los tamaños de sus componentes,
    // y la alineación requerida es la máxima alineación entre ellos.
    // Además ajustamos el tamaño final para que sea múltiplo de la alineación (padding al final).
    size_t max_size = 0;
    size_t max_alignment = 0 ;

    for (const auto& type_name : types) {
        const auto& component = find_type(type_name);
        max_size = std::max(max_size, component.size);
        max_alignment = std::max(max_alignment, component.alignment);
    }
    
    size_t final_size = align_up(max_size, max_alignment);
    size_t wasted_bytes = final_size - max_size; // padding final

    defined_types[name] = {
        name,
        final_size,
        max_alignment,
        types,
        "UNION"
    };
    std::cout << "Tipo UNION '" << name << "' definido (Tamaño: " << final_size << ", Alineación: " << max_alignment << ").\n";
}

// Estrategias de cálculo de tamaño y alineación

// Sin empaquetar:
// Calcula como lo haría el compilador por defecto: respeta la alineación de cada campo,
// añade padding interno antes de cada campo según su alineación, y padding final
// para que el tamaño total sea múltiplo de la máxima alineación del struct.
TypeManager::Result TypeManager::calculate_no_packing(const TypeDefinition& struct_type) {
    size_t current_offset = 0;
    size_t max_alignment = 0;
    size_t padding_internal = 0;

    for (const auto& type_name : struct_type.composition) {
        const auto& component_definition = find_type(type_name);
        
        size_t current_size = component_definition.size;
        size_t current_align = component_definition.alignment;

        // Manejar STRUCTs anidados recursivamente 
        if (component_definition.type_class == "STRUCT") {
            // Si es un STRUCT anidado, calculamos su tamaño y alineación no-empaquetados.
            Result res_nested = calculate_no_packing(component_definition);
            current_size = res_nested.size;
            current_align = res_nested.alignment;
            //
        }

        max_alignment = std::max(max_alignment, current_align);
        
        size_t next_offset = align_up(current_offset, current_align);
        padding_internal += (next_offset - current_offset);
        
        current_offset = next_offset + current_size;
    }

    size_t final_size = align_up(current_offset, max_alignment);
    size_t padding_final = final_size - current_offset;

    return {
        final_size,
        max_alignment,
        padding_internal + padding_final
    };
}
// Empaquetado (packed):
// Suma estricta de los tamaños de los campos sin introducir padding.
// Si un componente es un STRUCT anidado, su tamaño empaquetado se obtiene recursivamente.
// Esto refleja la situación donde se fuerza packed y los structs anidados también se consideran sin padding.
TypeManager::Result TypeManager::calculate_packed(const TypeDefinition& struct_type) {
    size_t total_size = 0;
    size_t max_alignment = 0;

    for (const auto& type_name : struct_type.composition) {
        const auto& component = find_type(type_name);
        
        size_t current_size = component.size;

        if (component.type_class == "STRUCT") {
            // ParaSTRUCT anidado, obtenemos su tamaño usando la estrategia packed recursiva.
            current_size = calculate_packed(component).size;
        }

        total_size += current_size;
        max_alignment = std::max(max_alignment, component.alignment);
    }
    
    return {
        total_size,
        max_alignment,
        0
    };
}

// Reordenando (optimo):
// Ordenamos los campos por alineación descendente y aplicamos la lógica de "sin empaquetar"
// sobre ese orden. La idea es minimizar padding interno colocando primero los campos
// de mayor alineación.
TypeManager::Result TypeManager::calculate_optimal_packing(const TypeDefinition& struct_type) {
    // Recolectamos pares {alineación, tamaño} para poder ordenar por alineación.
    std::vector<std::pair<size_t, size_t>> components_info;

    for (const auto& type_name : struct_type.composition) {
        const auto& defined_type = find_type(type_name);
        
        size_t current_size = defined_type.size;
        size_t current_align = defined_type.alignment;

        // Para STRUCT anidado usamos su tamaño y alineación bajo la estrategia "no packing"
        // porque al reordenar queremos las propiedades reales de memoria de ese componente.
        if (defined_type.type_class == "STRUCT") {
            Result res = calculate_no_packing(defined_type);
            current_size = res.size;
            current_align = res.alignment;
        }

        components_info.push_back({current_align, current_size});
    }

    std::sort(components_info.begin(), components_info.end(), 
        [](const auto& a, const auto& b) {
            return a.first > b.first;
        });

    size_t current_offset = 0;
    size_t max_alignment = 0;
    size_t padding_internal = 0;

    for (const auto& component : components_info) {
        max_alignment = std::max(max_alignment, component.first);
        
        size_t component_alignment = component.first;
        size_t next_offset = align_up(current_offset, component_alignment);
        padding_internal += (next_offset - current_offset);
        
        current_offset = next_offset + component.second;
    }

    size_t final_size = align_up(current_offset, max_alignment);
    size_t padding_final = final_size - current_offset;

    return {
        final_size,
        max_alignment,
        padding_internal + padding_final
    };
}

// Imprime los resultados de una estrategia de cálculo de tamaño/alineación.
// Uso para comparar las tres estrategias cuando se describe un STRUCT.
void TypeManager::print_result(const std::string& strategy, const Result& res) {
    std::cout << "  - Estrategia " << strategy << ":\n";
    std::cout << "    * Tamaño Total: " << res.size << " bytes\n";
    std::cout << "    * Alineación Requerida: " << res.alignment << " bytes\n";
    std::cout << "    * Bytes Desperdiciados (Padding): " << res.wasted << " bytes\n";
}

// Comando DESCRIBIR:
// Muestra la información detallada de un tipo definido.
// - Para ATOMIC muestra tamaño y alineación.
// - Para UNION muestra tamaño, alineación y padding final.
// - Para STRUCT calcula y muestra las tres estrategias (no packing, packed, reordenando).
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
            size_t max_size_component = 0;
            for (const auto& comp_name : type.composition) {
                max_size_component = std::max(max_size_component, find_type(comp_name).size);
            }
            std::cout << "  * Bytes Desperdiciados: " << (type.size - max_size_component) << " (padding final)\n";
        } else if (type.type_class == "STRUCT") {
            // Calculamos las tres estrategias y las imprimimos para comparar.
            Result res_no_packing = calculate_no_packing(type);
            print_result("Sin empaquetar", res_no_packing);

            Result res_packed = calculate_packed(type);
            print_result("Empaquetado", res_packed);

            Result res_optimal = calculate_optimal_packing(type);
            print_result("Reordenando de manera optima", res_optimal);
        }
        std::cout << "---------------------------------------\n";

    } catch (const std::runtime_error& e) {
        std::cerr << e.what() << "\n";
    }
}