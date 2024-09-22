package vm

import (
	"eventloop/interpreter/interpreter/code"
	"eventloop/interpreter/interpreter/object"
)

type Frame struct {
	fn          *object.CompiledFunction
	ip          int // instruction pointer for this frame, for this function
	basePointer int
}

func NewFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	return &Frame{fn: fn, ip: -1, basePointer: basePointer}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
