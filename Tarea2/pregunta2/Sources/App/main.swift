import pregunta2
public func runProgram() {
    print("Calculadora de Expresiones (Swift)")
    print("Comandos: EVAL PRE/POST <expr>, MOSTRAR PRE/POST <expr>, SALIR")
    
    // Bucle infinito que se rompe con el comando SALIR
    while true {
        print("\nIngrese acción: ", terminator: "") // terminator para que el cursor se quede en la misma linea
        
        // Leemos la linea completa del usuario.
        guard let input = readLine() else { continue }
        
        let components = input.uppercased().split(separator: " ").map { String($0) }
        
        guard let action = components.first else { continue }
        
        if action == "SALIR" {
            print("Saliendo del programa.")
            break
        }
        
        guard components.count >= 3 else {
            print("Error: Formato de comando incorrecto. Use EVAL/MOSTRAR <ORDEN> <EXPR>")
            continue
        }
        
        let order = components[1] // PRE o POST
        let exprTokens = Array(components[2...]) // El resto son los tokens de la expresion
        
        // Validacion del orden
        guard order == "PRE" || order == "POST" else {
            print("Error: Orden no reconocido. Debe ser PRE o POST.")
            continue
        }
        
        let isPrefix = (order == "PRE")
        
        switch action {
        case "EVAL":
            let result: Int?
            if isPrefix {
                result = evalPRE(tokens: exprTokens)
            } else {
                result = evalPOST(tokens: exprTokens)
            }
            
            if let finalResult = result {
                print("Resultado: \(finalResult)")
            } else {
                print("Error: La expresión es inválida o mal formada para la evaluación.")
            }
            
        case "MOSTRAR":
            let result: String?
            if isPrefix {
                result = displayPRE(tokens: exprTokens)
            } else {
                result = displayPOST(tokens: exprTokens)
            }
            
            if let finalResult = result {
                print("Infix: \(finalResult)")
            } else {
                print("Error: La expresión es inválida o mal formada para la conversión Infix.")
            }
            
        default:
            print("Error: Acción no reconocida. Use EVAL, MOSTRAR o SALIR.")
        }
    }
}

// Ejecutar el programa principal
runProgram()