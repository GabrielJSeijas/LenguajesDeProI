// Gabriel Seijas 19-00036
package main

import (
	"testing"
)

// Prueba la función String del bloque
func TestBlockString(t *testing.T) {
	// Caso cuando el bloque está libre
	blockFree := NewBlock(8, 0)
	expectedFree := "Dirección: 0, Tamaño: 8, Estado: LIBRE"
	if blockFree.String() != expectedFree {
		t.Errorf("String para bloque libre incorrecto. Esperaba '%s', obtuve '%s'", expectedFree, blockFree.String())
	}

	// Caso cuando el bloque está ocupado
	blockOccupied := NewBlock(4, 16)
	blockOccupied.Free = false
	blockOccupied.Tag = "MiProceso"
	expectedOccupied := "Dirección: 16, Tamaño: 4, Estado: OCUPADO (MiProceso)"
	if blockOccupied.String() != expectedOccupied {
		t.Errorf("String para bloque ocupado incorrecto. Esperaba '%s', obtuve '%s'", expectedOccupied, blockOccupied.String())
	}
}

// Prueba la función Split del bloque
func TestBlockSplit(t *testing.T) {
	// Split normal, el bloque se puede dividir
	parent := NewBlock(16, 0)
	leftChild, rightChild := parent.Split()

	if leftChild == nil || rightChild == nil {
		t.Fatalf("Split devolvió nil para un bloque dividible")
	}
	if leftChild.Size != 8 || leftChild.Address != 0 || leftChild.Parent != parent {
		t.Errorf("Hijo izquierdo incorrecto: %+v", leftChild)
	}
	if rightChild.Size != 8 || rightChild.Address != 8 || rightChild.Parent != parent {
		t.Errorf("Hijo derecho incorrecto: %+v", rightChild)
	}
	if parent.Free {
		t.Errorf("El bloque padre debería estar ocupado después de Split")
	}
	if parent.LeftChild != leftChild || parent.RightChild != rightChild {
		t.Errorf("Enlaces de hijo del padre incorrectos")
	}

	// Split cuando el tamaño es menor o igual a 1, no se puede dividir
	blockSmall := NewBlock(1, 0)
	lc, rc := blockSmall.Split()
	if lc != nil || rc != nil {
		t.Errorf("Split debería devolver nil, nil para bloques pequeños. Obtenido: %+v, %+v", lc, rc)
	}
}

// Prueba la función Merge del bloque
func TestBlockMerge(t *testing.T) {
	// Caso donde la fusión es exitosa
	parent := NewBlock(16, 0)
	parent.Free = false // El padre ya fue dividido
	leftChild := NewBlock(8, 0)
	leftChild.Parent = parent
	rightChild := NewBlock(8, 8)
	rightChild.Parent = parent
	parent.LeftChild = leftChild
	parent.RightChild = rightChild

	mergedParent := leftChild.Merge(rightChild)
	if mergedParent != parent {
		t.Errorf("Merge falló. Esperaba el bloque padre, obtuve: %+v", mergedParent)
	}
	if !parent.Free {
		t.Errorf("El bloque padre debería estar libre después de Merge")
	}
	if parent.LeftChild != nil || parent.RightChild != nil {
		t.Errorf("Los hijos del padre deberían ser nil después de Merge")
	}

	// --- Casos donde Merge debe fallar (probando las condiciones del if) ---

	// El bloque no está libre
	blockOccupied := NewBlock(8, 0)
	blockOccupied.Free = false // bloque ocupado
	buddyFree := NewBlock(8, 8)
	buddyFree.Parent = NewBlock(16, 0) // padre ficticio
	merged := blockOccupied.Merge(buddyFree)
	if merged != nil {
		t.Errorf("Merge debería fallar si b no está libre")
	}

	// El buddy no está libre
	blockFree := NewBlock(8, 0)
	buddyOccupied := NewBlock(8, 8)
	buddyOccupied.Free = false // buddy ocupado
	blockFree.Parent = NewBlock(16, 0)
	merged = blockFree.Merge(buddyOccupied)
	if merged != nil {
		t.Errorf("Merge debería fallar si buddy no está libre")
	}

	// El bloque no tiene padre
	blockNoParent := NewBlock(8, 0)
	blockNoParent.Free = true
	buddyNoParent := NewBlock(8, 8)
	buddyNoParent.Free = true
	merged = blockNoParent.Merge(buddyNoParent)
	if merged != nil {
		t.Errorf("Merge debería fallar si b.Parent es nil")
	}

	// El bloque no es el hijo izquierdo de su padre
	p := NewBlock(16, 0)
	b := NewBlock(8, 0)
	b.Free = true
	buddy := NewBlock(8, 8)
	buddy.Free = true
	p.LeftChild = buddy // orden invertido
	p.RightChild = b
	b.Parent = p
	buddy.Parent = p
	merged = b.Merge(buddy)
	if merged != nil {
		t.Errorf("Merge debería fallar si b no es el hijo izquierdo del padre")
	}

	// El buddy no es el hijo derecho de su padre
	p2 := NewBlock(16, 0)
	b2 := NewBlock(8, 0)
	b2.Free = true
	buddy2 := NewBlock(8, 8)
	buddy2.Free = true
	p2.LeftChild = b2
	p2.RightChild = NewBlock(8, 24) // otro bloque, no buddy2
	b2.Parent = p2
	buddy2.Parent = p2
	merged = b2.Merge(buddy2)
	if merged != nil {
		t.Errorf("Merge debería fallar si buddy no es el hijo derecho del padre")
	}
}
