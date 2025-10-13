// Gabriel Seijas 19-00036
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("--- Simulador de Manejador de Memoria (Buddy System) ---")
	fmt.Print("Ingrese la cantidad total de bloques de memoria (potencia de 2 recomendada): ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	totalBlocks, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Error: Cantidad de bloques inválida. Debe ser un número entero.")
		return
	}

	allocator, err := NewBuddyAllocator(totalBlocks)
	if err != nil {
		fmt.Printf("Error al inicializar el Buddy System: %v\n", err)
		return
	}

	fmt.Printf("Sistema Buddy inicializado con %d unidades de memoria.\n", allocator.TotalMemorySize)

	for {
		fmt.Print("\nIngrese una acción (RESERVAR <cantidad> <nombre> | LIBERAR <nombre> | MOSTRAR | SALIR): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue
		}

		action := strings.ToUpper(parts[0])

		switch action {
		case "RESERVAR":
			if len(parts) != 3 {
				fmt.Println("Error: Formato incorrecto. Uso: RESERVAR <cantidad> <nombre>")
				continue
			}
			size, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Error: La cantidad debe ser un número entero.")
				continue
			}
			name := parts[2]
			err = allocator.Reserve(size, name)
			if err != nil {
				fmt.Printf("Error al reservar: %v\n", err)
			} else {
				fmt.Printf("Memoria de %d unidades reservada para '%s'.\n", size, name)
			}
		case "LIBERAR":
			if len(parts) != 2 {
				fmt.Println("Error: Formato incorrecto. Uso: LIBERAR <nombre>")
				continue
			}
			name := parts[1]
			err := allocator.Free(name)
			if err != nil {
				fmt.Printf("Error al liberar: %v\n", err)
			} else {
				fmt.Printf("Memoria para '%s' liberada.\n", name)
			}
		case "MOSTRAR":
			allocator.Show()
		case "SALIR":
			fmt.Println("Saliendo del simulador.")
			return
		default:
			fmt.Println("Error: Acción no reconocida. Acciones válidas: RESERVAR, LIBERAR, MOSTRAR, SALIR.")
		}
	}
}

// Pruebas unitarias para el BuddyAllocator
// Estas pruebas las hice para practicar cómo funciona el sistema de asignación de memoria.
// No cubren todos los casos posibles, pero ayudan a ver si lo básico funciona.

func TestBuddyAllocator(t *testing.T) {
	// Prueba de inicialización del sistema
	allocator, err := NewBuddyAllocator(16)
	if err != nil {
		t.Fatalf("No se pudo inicializar el BuddyAllocator: %v", err)
	}
	if allocator.TotalMemorySize != 16 {
		t.Errorf("Esperaba tamaño total de memoria 16, pero obtuve %d", allocator.TotalMemorySize)
	}
	if len(allocator.FreeLists[4]) != 1 || allocator.FreeLists[4][0].Size != 16 { // log2(16) = 4
		t.Errorf("La lista de bloques libres inicial no es correcta")
	}

	// Prueba para reservar memoria
	err = allocator.Reserve(5, "procesoA")
	if err != nil {
		t.Errorf("Error al reservar 5 para procesoA: %v", err)
	}
	if _, exists := allocator.GetAllocatedBlocks()["procesoA"]; !exists {
		t.Errorf("procesoA no fue asignado")
	}

	// Intentar reservar más memoria de la disponible
	err = allocator.Reserve(100, "procesoB")
	if err == nil {
		t.Errorf("Debería haber fallado al reservar 100")
	}

	// Reservar con un nombre que ya existe
	err = allocator.Reserve(2, "procesoA")
	if err == nil {
		t.Errorf("Debería haber fallado al reservar con nombre duplicado")
	}

	// Prueba para liberar memoria
	err = allocator.Free("procesoA")
	if err != nil {
		t.Errorf("Error al liberar procesoA: %v", err)
	}
	if _, exists := allocator.GetAllocatedBlocks()["procesoA"]; exists {
		t.Errorf("procesoA no fue liberado")
	}

	// Intentar liberar un bloque que no existe
	err = allocator.Free("procesoC")
	if err == nil {
		t.Errorf("Debería haber fallado al liberar un bloque inexistente")
	}

	// Prueba de coalescencia (fusionar bloques libres)
	allocatorCoalesce, _ := NewBuddyAllocator(8)
	allocatorCoalesce.Reserve(1, "p1") // Asigna un bloque de tamaño 1
	allocatorCoalesce.Reserve(1, "p2") // Asigna el siguiente buddy de tamaño 1
}
