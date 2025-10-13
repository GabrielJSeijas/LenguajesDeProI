// Gabriel Seijas 19-00036
package main

import (
	"errors"
	"fmt"
	"math"
)

// BuddyAllocator maneja la memoria usando el método Buddy System
type BuddyAllocator struct {
	TotalMemorySize int               // Tamaño total de la memoria (debe ser potencia de 2)
	FreeLists       [][]*Block        // Listas de bloques libres por nivel
	AllocatedBlocks map[string]*Block // Bloques reservados identificados por tag
	RootBlock       *Block            // Bloque raíz que representa toda la memoria
}

// NewBuddyAllocator inicializa el sistema de memoria con el tamaño dado
func NewBuddyAllocator(totalBlocks int) (*BuddyAllocator, error) {
	if totalBlocks <= 0 {
		return nil, errors.New("el tamaño total de bloques debe ser positivo")
	}

	// Ajusta el tamaño a la siguiente potencia de 2
	totalMemorySize := 1
	for totalMemorySize < totalBlocks {
		totalMemorySize *= 2
	}

	maxLevel := int(math.Log2(float64(totalMemorySize))) + 1

	allocator := &BuddyAllocator{
		TotalMemorySize: totalMemorySize,
		FreeLists:       make([][]*Block, maxLevel),
		AllocatedBlocks: make(map[string]*Block),
	}

	// Crea el bloque raíz y lo pone en la lista de libres
	allocator.RootBlock = NewBlock(totalMemorySize, 0)
	allocator.addBlockToFreeList(allocator.RootBlock)

	return allocator, nil
}

// addBlockToFreeList agrega un bloque a la lista de libres según su tamaño
func (ba *BuddyAllocator) addBlockToFreeList(block *Block) {
	level := int(math.Log2(float64(block.Size)))
	if level < len(ba.FreeLists) {
		ba.FreeLists[level] = append(ba.FreeLists[level], block)
	}
}

// removeBlockFromFreeList quita un bloque de la lista de libres
func (ba *BuddyAllocator) removeBlockFromFreeList(block *Block) {
	level := int(math.Log2(float64(block.Size)))
	if level < len(ba.FreeLists) {
		for i, b := range ba.FreeLists[level] {
			if b == block {
				ba.FreeLists[level] = append(ba.FreeLists[level][:i], ba.FreeLists[level][i+1:]...)
				return
			}
		}
	}
}

// Reserve reserva un bloque de memoria del tamaño solicitado
func (ba *BuddyAllocator) Reserve(requestedSize int, tag string) error {
	if requestedSize <= 0 {
		return errors.New("el tamaño solicitado debe ser positivo")
	}
	if _, exists := ba.AllocatedBlocks[tag]; exists {
		return errors.New("ya existe un bloque con ese nombre")
	}

	// Busca el tamaño real (potencia de 2) que cubre la solicitud
	actualSize := 1
	for actualSize < requestedSize {
		actualSize *= 2
	}

	targetLevel := int(math.Log2(float64(actualSize)))

	var foundBlock *Block
	for level := targetLevel; level < len(ba.FreeLists); level++ {
		if len(ba.FreeLists[level]) > 0 {
			foundBlock = ba.FreeLists[level][0]
			ba.removeBlockFromFreeList(foundBlock)
			break
		}
	}

	if foundBlock == nil {
		return errors.New("no hay suficiente memoria disponible para la solicitud")
	}

	// Divide el bloque hasta llegar al tamaño necesario
	for foundBlock.Size > actualSize {
		leftChild, rightChild := foundBlock.Split()
		ba.addBlockToFreeList(rightChild)
		foundBlock = leftChild
	}

	foundBlock.Free = false
	foundBlock.Tag = tag
	ba.AllocatedBlocks[tag] = foundBlock
	return nil
}

// Free libera un bloque de memoria previamente reservado
func (ba *BuddyAllocator) Free(tag string) error {
	blockToFree, exists := ba.AllocatedBlocks[tag]
	if !exists {
		return errors.New("no existe un bloque con ese nombre")
	}

	blockToFree.Free = true
	blockToFree.Tag = ""
	delete(ba.AllocatedBlocks, tag)

	ba.addBlockToFreeList(blockToFree)

	// Intenta fusionar el bloque con su buddy si ambos están libres
	ba.coalesce(blockToFree)

	return nil
}

// coalesce fusiona un bloque con su buddy si ambos están libres, de forma recursiva
func (ba *BuddyAllocator) coalesce(block *Block) {
	if block.Parent == nil {
		return
	}

	buddy := ba.findBuddy(block)

	if buddy != nil && buddy.Free {
		ba.removeBlockFromFreeList(block)
		ba.removeBlockFromFreeList(buddy)

		parentBlock := block.Parent
		parentBlock.Free = true
		parentBlock.LeftChild = nil
		parentBlock.RightChild = nil

		ba.addBlockToFreeList(parentBlock)

		ba.coalesce(parentBlock)
	}
}

// findBuddy busca el bloque buddy de un bloque dado
func (ba *BuddyAllocator) findBuddy(block *Block) *Block {
	if block.Parent == nil {
		return nil
	}

	if block.Parent.LeftChild == block {
		return block.Parent.RightChild
	}
	return block.Parent.LeftChild
}

// Show muestra el estado actual de la memoria
func (ba *BuddyAllocator) Show() {
	fmt.Println("\n Estado de la Memoria ")
	ba.displayBlock(ba.RootBlock, 0)
	fmt.Println("---------------------------")
}

// displayBlock imprime la información de un bloque y sus hijos
func (ba *BuddyAllocator) displayBlock(block *Block, indent int) {
	if block == nil {
		return
	}

	for i := 0; i < indent; i++ {
		fmt.Print("  ")
	}

	status := "LIBRE"
	if !block.Free {
		status = fmt.Sprintf("OCUPADO (%s)", block.Tag)
	}
	fmt.Printf("├─ [Addr: %d, Size: %d, %s]\n", block.Address, block.Size, status)

	if block.LeftChild != nil {
		ba.displayBlock(block.LeftChild, indent+1)
	}
	if block.RightChild != nil {
		ba.displayBlock(block.RightChild, indent+1)
	}
}

// GetAllocatedBlocks regresa los bloques reservados (útil para pruebas)
func (ba *BuddyAllocator) GetAllocatedBlocks() map[string]*Block {
	return ba.AllocatedBlocks
}
