// Gabriel Seijas 19-00036
package main

import "fmt"

// Block representa un bloque de memoria en el sistema buddy
type Block struct {
	Size       int    // Tamaño del bloque (debe ser potencia de 2)
	Address    int    // Dirección inicial del bloque
	Free       bool   // Indica si el bloque está libre
	Tag        string // Etiqueta para identificar el bloque si está ocupado
	Parent     *Block // Referencia al bloque padre
	LeftChild  *Block // Referencia al hijo izquierdo
	RightChild *Block // Referencia al hijo derecho
}

// NewBlock crea un bloque nuevo con el tamaño y dirección dados
func NewBlock(size, address int) *Block {
	return &Block{
		Size:    size,
		Address: address,
		Free:    true,
		Tag:     "",
	}
}

// String devuelve una cadena con la información del bloque
func (b *Block) String() string {
	status := "LIBRE"
	if !b.Free {
		status = fmt.Sprintf("OCUPADO (%s)", b.Tag)
	}
	return fmt.Sprintf("Dirección: %d, Tamaño: %d, Estado: %s", b.Address, b.Size, status)
}

// Split divide el bloque en dos bloques de la mitad del tamaño
func (b *Block) Split() (*Block, *Block) {
	if b.Size <= 1 {
		return nil, nil
	}

	halfSize := b.Size / 2
	b.Free = false // El bloque original ya no está disponible como unidad

	leftChild := NewBlock(halfSize, b.Address)
	leftChild.Parent = b

	rightChild := NewBlock(halfSize, b.Address+halfSize)
	rightChild.Parent = b

	b.LeftChild = leftChild
	b.RightChild = rightChild

	return leftChild, rightChild
}

// Merge une dos bloques si ambos están libres y son hermanos
func (b *Block) Merge(buddy *Block) *Block {
	if b.Parent == nil || !b.Free || !buddy.Free || b.Parent.LeftChild != b || b.Parent.RightChild != buddy {
		return nil
	}

	parent := b.Parent
	parent.Free = true
	parent.LeftChild = nil
	parent.RightChild = nil
	return parent
}
