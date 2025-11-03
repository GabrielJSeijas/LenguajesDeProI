// MARK: - 1. Funciones Auxiliares

// Verificamos si un string es un operador aritmetico de la lista permitida.
public func isOperator(_ token: String) -> Bool {
    return ["+", "-", "*", "/"].contains(token)
}

// Devuelve la precedencia de un operador mayor numero, mayor precedencia
public func precedence(of op: String) -> Int {
    switch op {
    case "+", "-":
        return 1
    case "*", "/":
        return 2
    default:
        return 0 // Para cualquier otra cosa
    }
}
// Función para EVAL POST
public func evalPOST(tokens: [String]) -> Int? {
    var stack: [Int] = [] // La pila de operandos usamos Int
    
    for token in tokens {
        if let number = Int(token) {
            // Caso 1: Es un número el que opera entoncesLo metemos a la pila.
            stack.append(number)
        } else if isOperator(token) {
            // Caso 2: Es un operador, por ende necesita dos numeros 
            
            // Primero verificamos que haya al menos 2 elementos
            guard stack.count >= 2 else { return nil } //expresion mal formada (si no tiene 2)
            
            // Sacamos el segundo operando, el ultimo que entro 
            let operand2 = stack.removeLast()
            // Sacamos el primer operando, el penultimo
            let operand1 = stack.removeLast()
            
            var result: Int
            
            // Realizamos la operacion
            switch token {
            case "+":
                result = operand1 + operand2
            case "-":
                result = operand1 - operand2
            case "*":
                result = operand1 * operand2
            case "/":
                // Division entera
                result = operand1 / operand2
            default:
                return nil // Operador desconocido, no pasa
            }
            
            // Metemos el resultado de la operacion de nuevo a la pila
            stack.append(result)
        } else {
            // valor desconocido
            return nil
        }
    }
    
    // Al final, la pila solo tiene un valor el resultado final
    return stack.count == 1 ? stack.first : nil
}

public func evalPRE(tokens: [String]) -> Int? {
    // Invertimos el orden de los tokens para procesarla como POST de derecha a izquierda.
    let reversedTokens = tokens.reversed()
    var stack: [Int] = []
    
    for token in reversedTokens {
        if let number = Int(token) {
            // Caso 1: Es un número Lo metemos a la pila.
            stack.append(number)
        } else if isOperator(token) {
            // Caso 2: Es un operador, necesita dos numeros            
           guard stack.count >= 2 else { return nil }
        
        //Para evaluar PRE leyendo de derecha a izquierda
        // el primero que sacas (removeLast) es el numero DERECHO
        let operand2 = stack.removeLast()
        let operand1 = stack.removeLast() // El segundo que sacas es el IZQUIERDO
        
        var result: Int
            
            // Realizamos la operacion
            switch token {
            case "+":
                result = operand1 + operand2
            case "-":
                result = operand1 - operand2
            case "*":
                result = operand1 * operand2
            case "/":
                result = operand1 / operand2
            default:
                return nil
            }
            
            stack.append(result)
        } else {
            return nil
        }
    }
    
    return stack.count == 1 ? stack.first : nil
}

// Estructura para almacenar la expresion INFIX y la precedencia del operador RAIZ (todos asocian a la izquierda)
public struct Expression {
    public let infixString: String
    public let precedence: Int
    public let isLeftAssociative: Bool = true  
}

// Función que genera la expresion INFIX a partir de valores y un orden de procesamiento
public func toInfix(tokens: [String], isPrefix: Bool) -> String? {
    var stack: [Expression] = [] // Pila de expresiones o cadenas 
    
    // Si es pre procesamos de derecha a izquierda
    let processedTokens = isPrefix ? tokens.reversed() : tokens
    
    for token in processedTokens {
        if let _ = Int(token) {
            // Caso 1: Es un numero Lo empujamos a la pila como una expresion simple. Le damos precedencia maxima (3) para que nunca se le pongan parentesis.
            stack.append(Expression(infixString: token, precedence: 3))
        } else if isOperator(token) {
            
            guard stack.count >= 2 else { return nil } // Error
            
            // Sacamos los dos elementos de la pila expr2 es el que sale primero
            let expr2 = stack.removeLast()
            let expr1 = stack.removeLast()
            
            let currentOp = token
            let currentPrecedence = precedence(of: currentOp)

            var leftExpr: Expression
            var rightExpr: Expression

            if isPrefix {
                // Si la expresión original era PRE 
                // el primer elemento sacado (expr2) es el IZQUIERDO para la sintaxis Infija final.
                leftExpr = expr2   // IZQUIERDO
                rightExpr = expr1  // DERECHO
            } else {
                // Si la expresión original era POST y procesamos de izq a der 
                // el primer elemento sacado (expr2) es el DERECHO para la sintaxis Infija final.
                leftExpr = expr1   // IZQUIERDO
                rightExpr = expr2  // DERECHO
            }

           var left = leftExpr.infixString
// Condicion para Paréntesis en IZQUIERDA:
// Si la precedencia del hijo no es la máxima (3), significa que es el resultado de una operación anterior y debe ser encerrado para reflejar la estructura del árbol PRE.
if leftExpr.precedence < 3 { 
    left = "(\(left))"
}

// Determinar el numero DERECHO con parentesis
var right = rightExpr.infixString
// Condicion para Paréntesis en DERECHA:
// Se mantiene la lógica original, ya que el test para POST ya pasó.
if rightExpr.precedence <= currentPrecedence {
    right = "(\(right))"
}
            // Construimos la nueva expresion INFIX
            let newInfix = "\(left) \(currentOp) \(right)"
            
            // Empujamos la nueva expresion a la pila
            stack.append(Expression(infixString: newInfix, precedence: currentPrecedence))
        } else {
            return nil // Valor desconocido 
        }
    }
    
    // El resultado final es la expresion INFIX en la cima de la pila
    return stack.count == 1 ? stack.first?.infixString : nil
}

// Funciones para envolver 
public func displayPRE(tokens: [String]) -> String? {
    return toInfix(tokens: tokens, isPrefix: true)
}

public func displayPOST(tokens: [String]) -> String? {
    return toInfix(tokens: tokens, isPrefix: false)
}



