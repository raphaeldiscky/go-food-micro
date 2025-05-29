// Package typemapper provides a type mapper.
package typemapper

import "unsafe"

// typelinks2 is a link name for reflect.typelinks.
//
//go:linkname typelinks2 reflect.typelinks
func typelinks2() (sections []unsafe.Pointer, offset [][]int32)

// resolveTypeOff is a link name for reflect.resolveTypeOff.
//
//go:linkname resolveTypeOff reflect.resolveTypeOff
func resolveTypeOff(rtype unsafe.Pointer, off int32) unsafe.Pointer

// emptyInterface is a struct for empty interface.
type emptyInterface struct {
	typ  unsafe.Pointer
	data unsafe.Pointer
}
