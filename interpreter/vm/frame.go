package vm

import (
	"eventloop/interpreter/interpreter/code"
	"eventloop/interpreter/interpreter/object"
)

type Frame struct {
	fn *object.CompiledFunction
	ip int // instruction pointer for this frame, for this function
}

func NewFrame(fn *object.CompiledFunction) *Frame {
	return &Frame{fn: fn, ip: -1}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
