// Gabriel Seijas 19-00036
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Estructura para guardar la info de un programa y su lenguaje
type Programa struct {
	Nombre   string
	Lenguaje string
}

// Estructura para guardar la info de un intérprete (qué lenguaje interpreta y en cuál está hecho)
type Interprete struct {
	LenguajeBase string
	Lenguaje     string
}

// Estructura para guardar la info de un traductor (en qué lenguaje está hecho, de cuál a cuál traduce)
type Traductor struct {
	LenguajeBase    string
	LenguajeOrigen  string
	LenguajeDestino string
}

// Diccionarios para guardar los programas, intérpretes y traductores definidos
var programas = make(map[string]Programa)
var interpretes = make(map[string]Interprete)
var traductores = make(map[string]Traductor)

// El lenguaje LOCAL representa el que la compu puede ejecutar directamente
const LENGUAJE_LOCAL = "LOCAL"

func main() {
	fmt.Println("Simulador de Programas, Intérpretes y Traductores")
	fmt.Println("Comandos:")
	fmt.Println("  DEFINIR PROGRAMA <nombre> <lenguaje>")
	fmt.Println("  DEFINIR INTERPRETE <lenguaje_base> <lenguaje>")
	fmt.Println("  DEFINIR TRADUCTOR <lenguaje_base> <lenguaje_origen> <lenguaje_destino>")
	fmt.Println("  EJECUTABLE <nombre>")
	fmt.Println("  SALIR")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		parts := strings.Fields(input) // Separa el input por espacios

		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])

		switch command {
		case "DEFINIR":
			handleDefinir(parts[1:])
		case "EJECUTABLE":
			if len(parts) == 2 {
				handleEjecutable(parts[1])
			} else {
				fmt.Println("Error: Uso EJECUTABLE <nombre>")
			}
		case "SALIR":
			fmt.Println("Saliendo del simulador.")
			return
		default:
			fmt.Println("Comando no permitido. Tienes que usar DEFINIR, EJECUTABLE o SALIR.")
		}
	}
}

// Procesa el comando DEFINIR y guarda la info según el tipo
func handleDefinir(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Uso DEFINIR <tipo> [argumentos]")
		return
	}

	tipo := strings.ToUpper(args[0])
	switch tipo {
	case "PROGRAMA":
		if len(args) == 3 {
			nombre := args[1]
			lenguaje := args[2]
			programas[nombre] = Programa{Nombre: nombre, Lenguaje: lenguaje}
			fmt.Printf("Programa '%s' en %s definido.\n", nombre, lenguaje)
		} else {
			fmt.Println("Error: Se usa de la siguiente forma 'DEFINIR PROGRAMA <nombre> <lenguaje>'")
		}
	case "INTERPRETE":
		if len(args) == 3 {
			lenguajeBase := args[1]
			lenguaje := args[2]
			interpretes[lenguajeBase] = Interprete{LenguajeBase: lenguajeBase, Lenguaje: lenguaje}
			fmt.Printf("Intérprete de %s en %s definido.\n", lenguajeBase, lenguaje)
		} else {
			fmt.Println("Error: Se usa de la siguiente forma 'DEFINIR INTERPRETE <lenguaje_base> <lenguaje>'")
		}
	case "TRADUCTOR":
		if len(args) == 4 {
			lenguajeBase := args[1]
			lenguajeOrigen := args[2]
			lenguajeDestino := args[3]
			// La clave del traductor es una combinación de sus lenguajes
			traductores[fmt.Sprintf("%s-%s-%s", lenguajeBase, lenguajeOrigen, lenguajeDestino)] = Traductor{
				LenguajeBase:    lenguajeBase,
				LenguajeOrigen:  lenguajeOrigen,
				LenguajeDestino: lenguajeDestino,
			}
			fmt.Printf("Traductor de %s de %s a %s definido.\n", lenguajeBase, lenguajeOrigen, lenguajeDestino)
		} else {
			fmt.Println("Error: Se usa de la siguiente forma ' DEFINIR TRADUCTOR <lenguaje_base> <lenguaje_origen> <lenguaje_destino>'")
		}
	default:
		fmt.Println("Error: Tipo de definición desconocido. Usa PROGRAMA, INTERPRETE o TRADUCTOR.")
	}
}

// Busca si el programa puede ejecutarse en LOCAL y muestra la ruta si existe
func handleEjecutable(nombrePrograma string) {
	prog, ok := programas[nombrePrograma]
	if !ok {
		fmt.Printf("Error: El Programa '%s' no definido.\n", nombrePrograma)
		return
	}

	fmt.Printf("Intentando hacer ejecutable el programa '%s' (lenguaje: %s).\n", prog.Nombre, prog.Lenguaje)

	// Busca la ruta para ejecutar el programa en LOCAL
	path, found := encontrarRutaEjecucion(prog.Lenguaje, LENGUAJE_LOCAL)

	if found {
		fmt.Printf("El programa '%s' puede ser ejecutado en LOCAL siguiendo la ruta: %s\n", prog.Nombre, strings.Join(path, " -> "))
	} else {
		fmt.Printf("No se encontró una ruta para ejecutar el programa '%s' en LOCAL.\n", prog.Nombre)
	}
}

// Verifica si el lenguaje se puede ejecutar directamente en LOCAL o con algún intérprete
func puedeSerEjecutado(lenguaje string) bool {
	if lenguaje == LENGUAJE_LOCAL {
		return true
	}
	for _, interp := range interpretes {
		if interp.LenguajeBase == lenguaje && interp.Lenguaje == LENGUAJE_LOCAL {
			return true
		}
	}
	return false
}

// Busca una ruta para convertir el lenguaje de origen al destino usando intérpretes y traductores
func encontrarRutaEjecucion(origen, destino string) ([]string, bool) {
	visited := make(map[string]bool)
	var path []string

	return dfsEjecucion(origen, destino, path, visited)
}

// DFS para encontrar una ruta de ejecución entre lenguajes
func dfsEjecucion(actual, destino string, path []string, visited map[string]bool) ([]string, bool) {
	path = append(path, actual)
	visited[actual] = true

	if actual == destino {
		return path, true
	}

	// Busca intérpretes que puedan transformar el lenguaje actual
	for _, interp := range interpretes {
		if interp.LenguajeBase == actual && !visited[interp.Lenguaje] {
			if newPath, found := dfsEjecucion(interp.Lenguaje, destino, path, visited); found {
				return newPath, true
			}
		}
	}

	// Busca traductores que puedan transformar el lenguaje actual
	for _, trad := range traductores {
		if trad.LenguajeOrigen == actual && !visited[trad.LenguajeDestino] {
			if newPath, found := dfsEjecucion(trad.LenguajeDestino, destino, path, visited); found {
				return newPath, true
			}
		}
	}

	// Si no hay forma, regresa nil
	return nil, false
}
