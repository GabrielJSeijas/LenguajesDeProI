// Gabriel Seijas 19-00036
package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"testing"
)

// Función para capturar la salida de Show (esto lo uso para verificar lo que imprime Show)
func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = oldStdout
	log.SetOutput(os.Stderr)
	return buf.String() + string(out)
}

// Test para crear el BuddyAllocator y verificar que inicializa bien
func TestNewBuddyAllocator(t *testing.T) {
	allocator, err := NewBuddyAllocator(16)
	if err != nil {
		t.Fatalf("Error al crear el allocator: %v", err)
	}
	if allocator.TotalMemorySize != 16 {
		t.Errorf("El tamaño total no es el esperado")
	}
	if len(allocator.FreeLists[4]) != 1 || allocator.FreeLists[4][0].Size != 16 {
		t.Errorf("La lista de libres no tiene el bloque esperado")
	}

	allocator2, err := NewBuddyAllocator(10)
	if err != nil {
		t.Fatalf("Error al crear el allocator con tamaño no potencia de 2: %v", err)
	}
	if allocator2.TotalMemorySize != 16 {
		t.Errorf("No redondeó correctamente a la siguiente potencia de 2")
	}

	_, err = NewBuddyAllocator(0)
	if err == nil {
		t.Errorf("No dio error con tamaño cero")
	}
	_, err = NewBuddyAllocator(-5)
	if err == nil {
		t.Errorf("No dio error con tamaño negativo")
	}
}

// Test para reservar memoria y checar que los bloques se asignan bien
func TestReserve(t *testing.T) {
	allocator, _ := NewBuddyAllocator(32)

	err := allocator.Reserve(5, "procesoA")
	if err != nil {
		t.Errorf("No se pudo reservar memoria para procesoA")
	}
	if _, exists := allocator.AllocatedBlocks["procesoA"]; !exists {
		t.Errorf("No se asignó el bloque a procesoA")
	}
	if allocator.AllocatedBlocks["procesoA"].Size != 8 {
		t.Errorf("El tamaño asignado no es el esperado")
	}

	err = allocator.Reserve(1, "procesoB")
	if err != nil {
		t.Errorf("No se pudo reservar memoria para procesoB")
	}
	if allocator.AllocatedBlocks["procesoB"].Size != 1 {
		t.Errorf("El tamaño asignado a procesoB no es correcto")
	}

	// Prueba que no se pueda reservar con nombre duplicado
	err = allocator.Reserve(2, "procesoA")
	if err == nil || !strings.Contains(err.Error(), "ya existe un bloque") {
		t.Errorf("No detectó nombre duplicado")
	}

	// Prueba que detecta cuando no hay suficiente memoria
	allocatorFull, _ := NewBuddyAllocator(16)
	_ = allocatorFull.Reserve(10, "p1")
	err = allocatorFull.Reserve(2, "p2")
	if err == nil || !strings.Contains(err.Error(), "no hay suficiente memoria") {
		t.Errorf("No detectó falta de memoria")
	}

	// Prueba que no deja reservar si ya está todo ocupado
	allocatorFill, _ := NewBuddyAllocator(32)
	_ = allocatorFill.Reserve(16, "fill16")
	_ = allocatorFill.Reserve(8, "fill8")
	_ = allocatorFill.Reserve(4, "fill4")
	_ = allocatorFill.Reserve(2, "fill2")
	_ = allocatorFill.Reserve(1, "fill1")

	err = allocatorFill.Reserve(1, "fill1-fail")
	if err == nil || !strings.Contains(err.Error(), "no hay suficiente memoria") {
		t.Errorf("No detectó falta de memoria al llenar todo")
	}

	// Prueba reservar todo el bloque de un jalón
	allocatorBig, _ := NewBuddyAllocator(16)
	err = allocatorBig.Reserve(16, "big_proc")
	if err != nil {
		t.Errorf("No se pudo reservar todo el bloque")
	}
	if allocatorBig.AllocatedBlocks["big_proc"].Size != 16 {
		t.Errorf("El tamaño asignado no es el esperado para big_proc")
	}
}

// Test para liberar memoria y checar que los bloques se juntan bien
func TestFree(t *testing.T) {
	allocator, _ := NewBuddyAllocator(32)

	_ = allocator.Reserve(5, "procesoA")
	err := allocator.Free("procesoA")
	if err != nil {
		t.Errorf("No se pudo liberar procesoA")
	}
	if _, exists := allocator.AllocatedBlocks["procesoA"]; exists {
		t.Errorf("No se liberó procesoA")
	}

	// Prueba que no deja liberar un bloque que no existe
	err = allocator.Free("procesoZ")
	if err == nil || !strings.Contains(err.Error(), "no existe un bloque") {
		t.Errorf("No detectó liberar un bloque inexistente")
	}

	// Prueba que los bloques se pueden fusionar después de liberar
	allocatorCoalesce, _ := NewBuddyAllocator(8)
	_ = allocatorCoalesce.Reserve(1, "p1")
	_ = allocatorCoalesce.Reserve(1, "p2")
	_ = allocatorCoalesce.Reserve(2, "p3")

	_ = allocatorCoalesce.Free("p1")
	_ = allocatorCoalesce.Free("p2")

	_ = allocatorCoalesce.Free("p3")

	_ = allocatorCoalesce.Reserve(8, "final_block")
	if _, exists := allocatorCoalesce.AllocatedBlocks["final_block"]; !exists {
		t.Errorf("No se restauró el bloque grande después de liberar todo")
	}
	if allocatorCoalesce.AllocatedBlocks["final_block"].Address != 0 {
		t.Errorf("El bloque grande no está en la dirección 0")
	}
}

// Test para checar que Show imprime bien los bloques ocupados
func TestShow(t *testing.T) {
	allocator, _ := NewBuddyAllocator(16)
	_ = allocator.Reserve(3, "app1")
	_ = allocator.Reserve(7, "app2")

	output := captureOutput(func() {
		allocator.Show()
	})

	if !strings.Contains(output, "OCUPADO (app1)") {
		t.Errorf("No muestra que app1 está ocupado")
	}
	if !strings.Contains(output, "OCUPADO (app2)") {
		t.Errorf("No muestra que app2 está ocupado")
	}
}

// Test para Split y Merge de bloques, checa que los hijos y el padre se enlazan bien
func TestBlockSplitMerge(t *testing.T) {
	parent := NewBlock(16, 0)
	left, right := parent.Split()
	if left.Size != 8 || left.Address != 0 {
		t.Errorf("El hijo izquierdo no es correcto")
	}
	if right.Size != 8 || right.Address != 8 {
		t.Errorf("El hijo derecho no es correcto")
	}
	if parent.LeftChild != left || parent.RightChild != right {
		t.Errorf("Los enlaces del padre no son correctos")
	}

	allocator, _ := NewBuddyAllocator(16)
	_ = allocator.Reserve(1, "a")
	_ = allocator.Reserve(1, "b")
	_ = allocator.Reserve(1, "c")
	_ = allocator.Reserve(1, "d")

	allocator.Free("a")
	allocator.Free("b")

	if _, exists := allocator.GetAllocatedBlocks()["a"]; exists {
		t.Errorf("No se liberó el bloque 'a'")
	}
}

// Test para casos raros al agregar y quitar bloques de la free list
func TestAddRemoveBlockFromFreeListEdgeCases(t *testing.T) {
	allocator, _ := NewBuddyAllocator(4) // Max level 2 (log2(4))

	// Caso donde el nivel del bloque es mayor que las listas (no debe fallar)
	largeBlock := NewBlock(1024, 0)          // log2(1024) = 10
	allocator.addBlockToFreeList(largeBlock) // No debe causar pánico

	// Caso donde el nivel es mayor al quitar (tampoco debe fallar)
	allocator.removeBlockFromFreeList(largeBlock)

	// Caso donde el bloque no está en la lista (no debe modificar nada)
	blockNotFound := NewBlock(1, 99)
	level1Block := NewBlock(1, 0)
	allocator.addBlockToFreeList(level1Block)
	allocator.removeBlockFromFreeList(blockNotFound)
	if len(allocator.FreeLists[0]) != 1 {
		t.Errorf("removeBlockFromFreeList modificó la lista aunque el bloque no existía")
	}
}

// Test para reservar el bloque completo de memoria
func TestReserveEdgeCases(t *testing.T) {
	allocator, _ := NewBuddyAllocator(16)

	err := allocator.Reserve(16, "total")
	if err != nil {
		t.Errorf("Error al reservar el tamaño total: %v", err)
	}
	if _, exists := allocator.AllocatedBlocks["total"]; !exists {
		t.Errorf("Bloque 'total' no asignado")
	}
}

// Test para casos raros de coalesce (fusionar bloques)
func TestCoalesceEdgeCases(t *testing.T) {
	allocator, _ := NewBuddyAllocator(8) // Levels: 0, 1, 2, 3 (sizes 1, 2, 4, 8)

	// Caso donde el bloque es root y no tiene padre (no debe hacer nada)
	allocator.coalesce(allocator.RootBlock)

	// Caso donde el buddy está ocupado y no se debe fusionar
	_ = allocator.Reserve(1, "p1")
	_ = allocator.Reserve(1, "p2")

	blockP1 := allocator.AllocatedBlocks["p1"]

	blockP1.Free = true
	allocator.addBlockToFreeList(blockP1)
	allocator.coalesce(blockP1)

	levelP1 := int(math.Log2(float64(blockP1.Size)))
	foundInFreeList := false
	for _, b := range allocator.FreeLists[levelP1] {
		if b == blockP1 {
			foundInFreeList = true
			break
		}
	}
	if !foundInFreeList {
		t.Errorf("Bloque p1 se fusionó a pesar de que su buddy estaba ocupado")
	}
	blockP1.Free = false

	_ = allocator.Free("p1")
	_ = allocator.Free("p2")

	// Prueba que después de liberar todo se fusiona el bloque raíz
	_ = allocator.Reserve(2, "procA")
	_ = allocator.Reserve(2, "procB")
	_ = allocator.Reserve(2, "procC")
	_ = allocator.Reserve(2, "procD")

	_ = allocator.Free("procA")
	_ = allocator.Free("procB")
	_ = allocator.Free("procC")
	_ = allocator.Free("procD")

	if len(allocator.FreeLists[int(math.Log2(float64(8)))]) != 1 {
		t.Errorf("Después de liberar todos los bloques, el root no se fusionó completamente")
	}
	if allocator.FreeLists[int(math.Log2(float64(8)))][0].Size != 8 {
		t.Errorf("El bloque raíz fusionado no tiene el tamaño esperado")
	}
}

// Test para checar que findBuddy regresa el bloque correcto
func TestFindBuddy(t *testing.T) {
	allocator, _ := NewBuddyAllocator(8)

	// El root no tiene buddy
	buddy := allocator.findBuddy(allocator.RootBlock)
	if buddy != nil {
		t.Errorf("findBuddy para el RootBlock debería ser nil")
	}

	// Divido el root y checo los buddies
	leftChild, rightChild := allocator.RootBlock.Split()
	leftChild.Parent = allocator.RootBlock
	rightChild.Parent = allocator.RootBlock
	allocator.RootBlock.LeftChild = leftChild
	allocator.RootBlock.RightChild = rightChild

	buddy = allocator.findBuddy(leftChild)
	if buddy != rightChild {
		t.Errorf("Buddy de leftChild incorrecto. Esperaba %+v, obtuve %+v", rightChild, buddy)
	}

	buddy = allocator.findBuddy(rightChild)
	if buddy != leftChild {
		t.Errorf("Buddy de rightChild incorrecto. Esperaba %+v, obtuve %+v", leftChild, buddy)
	}
}

// Test para checar que displayBlock imprime bien cuando el bloque es nil o está libre/ocupado
func TestDisplayBlockNilAndFreeState(t *testing.T) {
	// Caso block == nil
	output := captureOutput(func() {
		ba := &BuddyAllocator{}
		ba.displayBlock(nil, 0)
	})
	if output != "" {
		t.Errorf("displayBlock para nil debería producir salida vacía")
	}

	// Caso block.Free = true
	blockFree := NewBlock(4, 0)
	output = captureOutput(func() {
		ba := &BuddyAllocator{}
		ba.displayBlock(blockFree, 0)
	})
	if !strings.Contains(output, "Estado: LIBRE") {
		t.Errorf("displayBlock no muestra estado LIBRE")
	}

	// Caso block.Free = false
	blockOccupied := NewBlock(4, 0)
	blockOccupied.Free = false
	blockOccupied.Tag = "test"
	output = captureOutput(func() {
		ba := &BuddyAllocator{}
		ba.displayBlock(blockOccupied, 0)
	})
	if !strings.Contains(output, "Estado: OCUPADO (test)") {
		t.Errorf("displayBlock no muestra estado OCUPADO")
	}
}
