#include "TypeManager.h"
#include <sstream>

// Programa principal del simulador de tipos de datos.
// Contiene pruebas iniciales para asegurar cobertura y un bucle interactivo
// que permite al usuario definir atomicos, structs, unions y pedir descripciones.

void run_tests(TypeManager& manager) {
    // Ejecuta un conjunto de pruebas automatizadas que ejercitan las funciones
    // principales del TypeManager. Estas pruebas crean tipos atómicos, structs y
    // unions y llaman a la función de descripción para verificar la salida.
    std::cout << "\n--- EJECUTANDO PRUEBAS UNITARIAS ---\n";
    
    // Crear tipos atómicos con representación y alineación.
    manager.execute_atomic("char", 1, 1);
    manager.execute_atomic("short",2, 2);
    manager.execute_atomic("int",4, 4);
    manager.execute_atomic("double", 8, 8);

    manager.execute_atomic("byte", 1, 2);

    // Crear una union grande para probar tamaño y alineación.
    manager.execute_union("big_union", {"char", "int", "double"});

    // Crear structs con distintos ordenamientos para probar padding y optimización.
    manager.execute_struct("S1", {"char", "int", "char"});
    manager.execute_struct("S2_optimal", {"int", "char", "char"});
    manager.execute_struct("S3_padding_final", {"char", "int", "short"});

    // LLamar a describir para mostrar detalles internos de los tipos creados.
    std::cout << "\n--- DESCRIPCIONES DE PRUEBA ---\n";
    manager.execute_describe("big_union");
    manager.execute_describe("S1");
    manager.execute_describe("S2_optimal");

    std::cout << "--- PRUEBAS COMPLETADAS ---\n\n";
}

int main() {

    // Ajuste de localización en Windows para manejar correctamente UTF-8 en la salida.
    #ifdef _WIN32
    std::setlocale(LC_ALL, "es_ES.UTF-8");
    #endif

    // Instancia del gestor de tipos que mantiene la tabla de tipos y operaciones.
    TypeManager manager;
    std::string line;
    bool running = true;

    // Ejecutar pruebas automáticas al inicio para garantizar que las funciones
    // básicas ya están registradas y probadas antes de la interacción con el usuario.
    run_tests(manager);

    // Mensajes iniciales y ayuda mínima para el usuario.
    std::cout << "--- MANEJADOR DE TIPOS DE DATOS ---\n";
    std::cout << "Comandos: ATOMICO, STRUCT, UNION, DESCRIBIR, SALIR.\n";

    // Bucle principal de línea de comando: lee una línea, tokeniza y ejecuta.
    while (running) {
        std::cout << "Accion: ";
        std::getline(std::cin, line);

        std::stringstream ss(line);
        std::string command;
        ss >> command;

        try {
            if (command == "ATOMICO") {
                // Sintaxis: ATOMICO <nombre> <representación> <alineación>
                std::string name;
                size_t rep, align;
                ss >> name >> rep >> align;
                if (name.empty() || ss.fail()) throw std::runtime_error("Uso: ATOMICO <nombre> <representación> <alineación>");
                manager.execute_atomic(name, rep, align);

            } else if (command == "STRUCT" || command == "UNION") {
                // Recolecta todos los tipos componentes y delega en TypeManager.
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

            } else if (command =="DESCRIBIR") {
                // Sintaxis: DESCRIBIR <nombre>
                std::string name;
                ss >> name;
                if (name.empty()) throw std::runtime_error("Uso: DESCRIBIR <nombre>");
                manager.execute_describe(name);

            } else if (command == "SALIR") {
                // Mensaje de salida y terminación del programa.
                std::cout << "Saliendo del simulador.\n";
                running = false;

            } else if (!command.empty()) {
                // Manejo de comandos desconocidos.
                std::cout << "Comando desconocido.\n";

            }
        } catch (const std::runtime_error& e) {
            // Mostrar errores de sintaxis o uso incorrecto de comandos.
            std::cerr << "Error de Sintaxis: " << e.what() << "\n";
        }
    }

    return 0;
}