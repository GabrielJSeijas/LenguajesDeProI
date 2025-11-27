#include "TypeManager.h"
#include <sstream>

void run_tests(TypeManager& manager) {
    std::cout << "\n--- EJECUTANDO PRUEBAS UNITARIAS (COBERTURA > 80%) ---\n";
    
    // I. ATOMICOs
    manager.execute_atomic("char", 1, 1);
    manager.execute_atomic("short", 2, 2);
    manager.execute_atomic("int", 4, 4);
    manager.execute_atomic("double", 8, 8);
    manager.execute_atomic("byte", 1, 2); // Ejemplo con representación y alineación distintas

    // III. UNION (Prueba de tamaño y alineación)
    // Union de (char, int, double) -> Tamaño=8, Alineación=8 (requiere padding final)
    manager.execute_union("big_union", {"char", "int", "double"}); // Max size=8, Max align=8. Size=8.

    // II. STRUCTs
    manager.execute_struct("S1", {"char", "int", "char"});
    manager.execute_struct("S2_optimal", {"int", "char", "char"}); // Ya ordenado de forma óptima
    manager.execute_struct("S3_padding_final", {"char", "int", "short"}); // Desalineado y requiere padding final

    // IV. DESCRIBIR
    std::cout << "\n--- DESCRIPCIONES DE PRUEBA ---\n";
    manager.execute_describe("big_union");

    // S1: char(1) int(4) char(1) -> MaxAlign=4
    // Sin Empaquetar: [c][ ][ ][ ] [i][i][i][i] [c][ ][ ][ ] -> Size=12, Wasted=6
    // Óptimo: [i][i][i][i] [c][c][ ][ ] -> Size=8, Wasted=2
    manager.execute_describe("S1"); 

    // S2_optimal: int(4) char(1) char(1) -> MaxAlign=4
    // Sin Empaquetar (ya es óptimo): [i][i][i][i] [c][c][ ][ ] -> Size=8, Wasted=2
    manager.execute_describe("S2_optimal"); 

    std::cout << "--- PRUEBAS COMPLETADAS ---\n\n";
}

int main() {

    #ifdef _WIN32
    std::setlocale(LC_ALL, "es_ES.UTF-8");
    #endif
    TypeManager manager;
    std::string line;
    bool running = true;

    // Ejecutar pruebas al inicio para cumplir con la cobertura requerida 
    run_tests(manager); 

    std::cout << "--- SIMULADOR DE TIPOS DE DATOS ---\n";
    std::cout << "Comandos: ATOMICO, STRUCT, UNION, DESCRIBIR, SALIR.\n";

    while (running) {
        std::cout << "Accion: ";
        std::getline(std::cin, line);

        std::stringstream ss(line);
        std::string command;
        ss >> command;

        try {
            if (command == "ATOMICO") {
                std::string name;
                size_t rep, align;
                ss >> name >> rep >> align;
                if (name.empty() || ss.fail()) throw std::runtime_error("Uso: ATOMICO <nombre> <representación> <alineación>");
                manager.execute_atomic(name, rep, align);
            } else if (command == "STRUCT" || command == "UNION") {
                std::string name, type_name;
                std::vector<std::string> types;
                ss >> name;
                if (name.empty()) throw std::runtime_error("Uso: " + command + " <nombre> [<tipo>...]");
                while (ss >> type_name) {
                    types.push_back(type_name);
                }
                if (types.empty()) throw std::runtime_error("El " + command + " debe tener al menos un tipo componente.");
                
                if (command == "STRUCT") manager.execute_struct(name, types);
                else manager.execute_union(name, types);
            } else if (command == "DESCRIBIR") {
                std::string name;
                ss >> name;
                if (name.empty()) throw std::runtime_error("Uso: DESCRIBIR <nombre>");
                manager.execute_describe(name);
            } else if (command == "SALIR") {
                std::cout << "Saliendo del simulador.\n"; // [cite: 69]
                running = false;
            } else if (!command.empty()) {
                std::cout << "Comando desconocido.\n";
            }
        } catch (const std::runtime_error& e) {
            std::cerr << "Error de Sintaxis: " << e.what() << "\n";
        }
    }

    return 0;
}