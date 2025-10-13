// Gabriel Seijas 19-00036
package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// Esta función reinicia los mapas globales antes de cada test
func resetGlobalState() {
	programas = make(map[string]Programa)
	interpretes = make(map[string]Interprete)
	traductores = make(map[string]Traductor)
}

// Prueba para definir un programa y verificar que se guarda bien
func TestDefinirPrograma(t *testing.T) {
	resetGlobalState()
	handleDefinir([]string{"PROGRAMA", "miApp", "GO"})
	if _, ok := programas["miApp"]; !ok {
		t.Errorf("Programa 'miApp' no fue definido correctamente.")
	}
}

// Prueba para definir un intérprete y ver que se guarda en el mapa
func TestDefinirInterprete(t *testing.T) {
	resetGlobalState()
	handleDefinir([]string{"INTERPRETE", "GO", "LOCAL"})
	if _, ok := interpretes["GO"]; !ok {
		t.Errorf("Intérprete de GO no fue definido correctamente.")
	}
}

// Prueba para definir un traductor y checar que la clave esté en el mapa
func TestDefinirTraductor(t *testing.T) {
	resetGlobalState()
	handleDefinir([]string{"TRADUCTOR", "COMPILADOR_C", "C", "BINARIO"})
	key := fmt.Sprintf("%s-%s-%s", "COMPILADOR_C", "C", "BINARIO")
	if _, ok := traductores[key]; !ok {
		t.Errorf("Traductor de C a BINARIO no fue definido correctamente.")
	}
}

// Prueba para cuando no hay ruta para ejecutar un programa
func TestHandleEjecutable_NoRouteFound(t *testing.T) {
	resetGlobalState()
	programas["unreachableProg"] = Programa{Nombre: "unreachableProg", Lenguaje: "FANTASY_LANG"}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	handleEjecutable("unreachableProg")

	w.Close()
	out, _ := io.ReadAll(r)
	output := string(out)

	if !strings.Contains(output, "No se encontró una ruta para ejecutar el programa 'unreachableProg' en LOCAL.") {
		t.Errorf("Esperaba mensaje de 'no se encontró ruta' para un programa inalcanzable, obtuve: %s", output)
	}
}

// Prueba para simular la entrada por stdin y checar varios comandos
func TestMainInputHandling(t *testing.T) {
	resetGlobalState()

	interpretes["GO"] = Interprete{LenguajeBase: "GO", Lenguaje: LENGUAJE_LOCAL}

	input := `DEFINIR PROGRAMA miApp GO
EJECUTABLE miApp
EJECUTABLE
UNKNOWN_COMMAND
SALIR`
	oldStdin := os.Stdin
	rIn, wIn, _ := os.Pipe()
	os.Stdin = rIn
	go func() {
		defer wIn.Close()
		fmt.Fprint(wIn, input)
	}()

	oldStdout := os.Stdout
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut
	defer func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
	}()

	main()

	wOut.Close()
	out, _ := io.ReadAll(rOut)
	output := string(out)

	if !strings.Contains(output, "Simulador de Programas, Intérpretes y Traductores") {
		t.Errorf("Falta mensaje de bienvenida. Output: %s", output)
	}
	if !strings.Contains(output, "Programa 'miApp' en GO definido.") {
		t.Errorf("No se definió el programa miApp. Output: %s", output)
	}
	if !strings.Contains(output, "El programa 'miApp' puede ser ejecutado en LOCAL siguiendo la ruta: GO -> LOCAL") {
		t.Errorf("Salida EJECUTABLE miApp incorrecta. Output: %s", output)
	}
	if !strings.Contains(output, "Error: Uso EJECUTABLE <nombre>") {
		t.Errorf("No se manejó el error de EJECUTABLE sin nombre. Output: %s", output)
	}
	if !strings.Contains(output, "Comando desconocido. Por favor, usa DEFINIR, EJECUTABLE o SALIR.") {
		t.Errorf("No se manejó el comando desconocido. Output: %s", output)
	}
	if !strings.Contains(output, "Saliendo del simulador.") {
		t.Errorf("No se salió correctamente. Output: %s", output)
	}
}

// Prueba para checar los errores al definir cosas mal
func TestDefinirErrores(t *testing.T) {
	resetGlobalState()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	handleDefinir([]string{})
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old
	if !strings.Contains(string(out), "Error: Uso DEFINIR <tipo> [argumentos]") {
		t.Errorf("Esperaba error para DEFINIR sin argumentos, obtuve: %s", string(out))
	}

	resetOutput := func() {
		old = os.Stdout
		r, w, _ = os.Pipe()
		os.Stdout = w
	}
	captureOutput := func() string {
		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = old
		return string(out)
	}

	resetOutput()
	handleDefinir([]string{})
	if outStr := captureOutput(); !strings.Contains(outStr, "Error: Uso DEFINIR <tipo> [argumentos]") {
		t.Errorf("Esperaba error para DEFINIR sin argumentos, obtuve: %s", outStr)
	}

	resetOutput()
	handleDefinir([]string{"PROGRAMA", "nombre"})
	if outStr := captureOutput(); !strings.Contains(outStr, "Error: Uso DEFINIR PROGRAMA <nombre> <lenguaje>") {
		t.Errorf("Esperaba error para DEFINIR PROGRAMA sin lenguaje, obtuve: %s", outStr)
	}

	resetOutput()
	handleDefinir([]string{"INTERPRETE", "GO"})
	if outStr := captureOutput(); !strings.Contains(outStr, "Error: Uso DEFINIR INTERPRETE <lenguaje_base> <lenguaje>") {
		t.Errorf("Esperaba error para DEFINIR INTERPRETE con pocos argumentos, obtuve: %s", outStr)
	}

	resetOutput()
	handleDefinir([]string{"TRADUCTOR", "C++", "C"})
	if outStr := captureOutput(); !strings.Contains(outStr, "Error: Uso DEFINIR TRADUCTOR <lenguaje_base> <lenguaje_origen> <lenguaje_destino>") {
		t.Errorf("Esperaba error para DEFINIR TRADUCTOR con pocos argumentos, obtuve: %s", outStr)
	}

	resetOutput()
	handleDefinir([]string{"OTROTIPO", "algo", "mas"})
	if outStr := captureOutput(); !strings.Contains(outStr, "Error: Tipo de definición desconocido. Usa PROGRAMA, INTERPRETE o TRADUCTOR.") {
		t.Errorf("Esperaba error para tipo de definición desconocido, obtuve: %s", outStr)
	}
}

// Prueba para ver si puedeSerEjecutado funciona bien
func TestPuedeSerEjecutado(t *testing.T) {
	resetGlobalState()
	if !puedeSerEjecutado(LENGUAJE_LOCAL) {
		t.Errorf("LOCAL debería ser siempre ejecutable.")
	}
	if puedeSerEjecutado("GO") {
		t.Errorf("GO no debería ser ejecutable sin un intérprete.")
	}
	interpretes["GO"] = Interprete{LenguajeBase: "GO", Lenguaje: LENGUAJE_LOCAL}
	if !puedeSerEjecutado("GO") {
		t.Errorf("GO debería ser ejecutable con un intérprete GO -> LOCAL.")
	}
}

// Prueba para checar que encontrarRutaEjecucion encuentra la ruta correcta
func TestEncontrarRutaEjecucion(t *testing.T) {
	resetGlobalState()

	path, found := encontrarRutaEjecucion(LENGUAJE_LOCAL, LENGUAJE_LOCAL)
	if !found || strings.Join(path, " -> ") != "LOCAL" {
		t.Errorf("Ruta para LOCAL a LOCAL incorrecta. Obtenido: %v, %v", path, found)
	}

	interpretes["GO"] = Interprete{LenguajeBase: "GO", Lenguaje: LENGUAJE_LOCAL}
	path, found = encontrarRutaEjecucion("GO", LENGUAJE_LOCAL)
	if !found || strings.Join(path, " -> ") != "GO -> LOCAL" {
		t.Errorf("Ruta para GO a LOCAL incorrecta. Obtenido: %v, %v", path, found)
	}

	interpretes["JS"] = Interprete{LenguajeBase: "JS", Lenguaje: LENGUAJE_LOCAL}
	traductores[fmt.Sprintf("%s-%s-%s", "TS_COMPILER", "TS", "JS")] = Traductor{
		LenguajeBase: "TS_COMPILER", LenguajeOrigen: "TS", LenguajeDestino: "JS"}
	path, found = encontrarRutaEjecucion("TS", LENGUAJE_LOCAL)
	if !found || strings.Join(path, " -> ") != "TS -> JS -> LOCAL" {
		t.Errorf("Ruta para TS a LOCAL incorrecta. Obtenido: %v, %v", path, found)
	}

	interpretes["BINARIO"] = Interprete{LenguajeBase: "BINARIO", Lenguaje: LENGUAJE_LOCAL}
	traductores[fmt.Sprintf("%s-%s-%s", "PY_TO_C", "PYTHON", "C")] = Traductor{
		LenguajeBase: "PY_TO_C", LenguajeOrigen: "PYTHON", LenguajeDestino: "C"}
	traductores[fmt.Sprintf("%s-%s-%s", "C_COMPILER", "C", "BINARIO")] = Traductor{
		LenguajeBase: "C_COMPILER", LenguajeOrigen: "C", LenguajeDestino: "BINARIO"}
	path, found = encontrarRutaEjecucion("PYTHON", LENGUAJE_LOCAL)
	if !found || strings.Join(path, " -> ") != "PYTHON -> C -> BINARIO -> LOCAL" {
		t.Errorf("Ruta para PYTHON a LOCAL incorrecta. Obtenido: %v, %v", path, found)
	}

	found = false
	_, found = encontrarRutaEjecucion("UNKNOWN_LANG", LENGUAJE_LOCAL)
	if found {
		t.Errorf("Se encontró una ruta para un lenguaje desconocido, lo cual es incorrecto.")
	}
}

// Prueba para checar que handleEjecutable imprime lo correcto
func TestHandleEjecutable(t *testing.T) {
	resetGlobalState()
	programas["testProg"] = Programa{Nombre: "testProg", Lenguaje: "GO"}
	interpretes["GO"] = Interprete{LenguajeBase: "GO", Lenguaje: LENGUAJE_LOCAL}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	handleEjecutable("testProg")

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old

	if !strings.Contains(string(out), "El programa 'testProg' puede ser ejecutado en LOCAL siguiendo la ruta: GO -> LOCAL") {
		t.Errorf("Salida de EJECUTABLE incorrecta. Obtenido: %s", string(out))
	}

	old = os.Stdout
	r, w, _ = os.Pipe()
	os.Stdout = w
	handleEjecutable("nonExistentProg")
	w.Close()
	out, _ = io.ReadAll(r)
	os.Stdout = old
	if !strings.Contains(string(out), "Error: Programa 'nonExistentProg' no definido.") {
		t.Errorf("Esperaba error para programa no existente, obtuve: %s", string(out))
	}
}
