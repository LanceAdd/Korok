//go:build js

package webav

import (
	"syscall/js"
	"unsafe"
)

// Uint8Array returns a JavaScript Uint8Array containing a copy of src.
func Uint8Array(src []uint8) js.Value {
	return newTypedArray("Uint8Array", src)
}

// Int32Array returns a JavaScript Int32Array containing a copy of src.
func Int32Array(src []int32) js.Value {
	if len(src) == 0 {
		return newTypedArray("Int32Array", nil)
	}
	size := int(unsafe.Sizeof(src[0]))
	return newTypedArray("Int32Array", unsafe.Slice((*byte)(unsafe.Pointer(&src[0])), len(src)*size))
}

// Float32Array returns a JavaScript Float32Array containing a copy of src.
func Float32Array(src []float32) js.Value {
	if len(src) == 0 {
		return newTypedArray("Float32Array", nil)
	}
	size := int(unsafe.Sizeof(src[0]))
	return newTypedArray("Float32Array", unsafe.Slice((*byte)(unsafe.Pointer(&src[0])), len(src)*size))
}

func newTypedArray(ctor string, src []byte) js.Value {
	bytes := js.Global().Get("Uint8Array").New(len(src))
	if len(src) > 0 {
		js.CopyBytesToJS(bytes, src)
	}
	return js.Global().Get(ctor).New(bytes.Get("buffer"))
}
