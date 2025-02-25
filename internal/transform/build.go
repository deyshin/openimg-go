//go:build cgo
// +build cgo

package transform

// #cgo CFLAGS: -Wno-unused-result -Wno-xor-used-as-pow
import "C"