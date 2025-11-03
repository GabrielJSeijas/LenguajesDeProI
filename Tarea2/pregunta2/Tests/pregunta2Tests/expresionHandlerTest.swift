import XCTest // Importa el framework de pruebas.
@testable import pregunta2

// Creamos una clase que hereda de XCTestCase para agrupar las pruebas.
class ExprHandlerTests: XCTestCase {
    
    // Prueba para EVAL (caso simple y complejo)
    func testEval() {
        // Prueba EVAL POST
        let postTokens = ["8", "3", "-", "8", "4", "4", "*", "+", "*"]
        let postResult = evalPOST(tokens: postTokens)
        // La expresion (8-3) * (8 + (4*4)) = 5 * (8 + 16) = 5 * 24 = 120
        XCTAssertEqual(postResult, 120, "La evaluacion POST no es correcta.")
        
        // Prueba EVAL PRE
        let preTokens = ["+", "+", "*", "3", "4", "5", "7"]
        let preResult = evalPRE(tokens: preTokens)
        // La expresion + (* 3 4) (5 7) no es la misma que la del ejemplo.
        // La expresion del ejemplo (+ * 3 4 5 7) es: ((3 * 4) + 5) + 7 = (12+5)+7 = 24
        // Usamos el ejemplo del requerimiento: EVAL PRE + * 3 4 5 7 (debería ser 24)
        XCTAssertEqual(preResult, 24, "La evaluacion PRE no es correcta.")
    }
    
    // Prueba para MOSTRAR (Infix)
    func testToInfix() {
        // Prueba MOSTRAR POST (Asociatividad y Precedencia)
        let postTokens = ["8", "3", "-", "8", "4", "4", "*", "+", "*"]
        let postInfix = displayPOST(tokens: postTokens)
        // La salida esperada es (8 - 3) * (8 + (4 * 4))
        XCTAssertEqual(postInfix, "(8 - 3) * (8 + 4 * 4)", "La conversion POST a Infix es incorrecta.")
        
        // Prueba MOSTRAR PRE (Asociatividad y Precedencia)
        let preTokens = ["+", "+", "*", "3", "4", "5", "7"]
        // La salida esperada es ((3 * 4) + 5) + 7
        let preInfix = displayPRE(tokens: preTokens)
XCTAssertEqual(preInfix, "((3 * 4) + 5) + 7", "La conversion PRE a Infix es incorrecta.")
    }
    
    // Prueba de manejo de errores (Expresiones incompletas)
    func testErrorHandling() {
        // POST con demasiados operadores
        let badPost = evalPOST(tokens: ["3", "4", "+", "*"])
        XCTAssertNil(badPost, "Deberia fallar con un operador extra.")
        
        // PRE con demasiados operandos (Error al procesar de derecha a izquierda)
        let badPre = evalPRE(tokens: ["+", "3", "4", "5"])
        XCTAssertNil(badPre, "Deberia fallar con un operando extra.")
    }
}