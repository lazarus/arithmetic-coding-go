package ac

import (
	"bytes"
)

const (
	codeBits = 16
	one      = (1 << codeBits) - 1
	half     = 1 << (codeBits - 1)
	quarter  = half / 2

	frequencyBits = 14
	maxFreq       = (1 << frequencyBits) - 1

	EOF          = 256
	totalSymbols = 257
)

// Context holds the state for the arithmetic encoder/decoder.
type Context struct {
	Probability
	value          uint32
	pending        int
	buffer         *bytes.Buffer
	frequencies    [totalSymbols + 1]uint32
	totalFrequency uint32
	frozen         bool
	*InputBuffer
	*OutputBuffer
}

// NewContext creates a new arithmetic coding context.
func NewContext() *Context {
	context := &Context{
		Probability: Probability{
			low:  0,
			high: one,
		},
		buffer:       new(bytes.Buffer),
		InputBuffer:  NewInputBuffer(),
		OutputBuffer: NewOutputBuffer(),
	}
	context.Reset()
	return context
}

type Probability struct {
	low, high uint32
}

func (ac *Context) Reset() {
	ac.totalFrequency = totalSymbols
	ac.frozen = false

	for i := uint32(0); i <= totalSymbols; i++ {
		ac.frequencies[i] = i
	}
}

func (ac *Context) update(symbol uint32) {
	for i := symbol + 1; i <= totalSymbols; i++ {
		ac.frequencies[i]++
	}
	ac.totalFrequency++
	if ac.totalFrequency >= maxFreq {
		ac.frozen = true
	}
}

func (ac *Context) getProbability(symbol uint32) (Probability, uint32) {
	prob := Probability{
		low:  ac.frequencies[symbol],
		high: ac.frequencies[symbol+1],
	}
	total := ac.totalFrequency
	if !ac.frozen {
		ac.update(symbol)
	}
	return prob, total
}

func (ac *Context) getChar(offset uint32) (prob Probability, total uint32, char uint32) {
	for i := uint32(0); i < totalSymbols; i++ {
		if offset < ac.frequencies[i+1] {
			char = i
			prob = Probability{
				low:  ac.frequencies[i],
				high: ac.frequencies[i+1],
			}
			total = ac.totalFrequency
			if !ac.frozen {
				ac.update(char)
			}
			return
		}
	}
	return
}

// Encode takes a slice of bytes and encodes it using arithmetic coding.
func (ac *Context) Encode(data []byte) []byte {
	for _, b := range data {
		ac.encodeSymbol(uint32(b))
	}
	ac.encodeSymbol(EOF)

	ac.finalizeEncoding()
	return ac.buffer.Bytes()
}

// Decode takes an encoded byte slice and decodes it.
func (ac *Context) Decode(encoded []byte) []byte {
	ac.buffer = bytes.NewBuffer(encoded)
	output := bytes.Buffer{}

	ac.value = ac.initializeDecoding()

	for {
		if b := ac.decodeSymbol(); b != EOF {
			output.WriteByte(byte(b))
		} else {
			break
		}
	}

	return output.Bytes()
}

// Encode a single symbol.
func (ac *Context) encodeSymbol(symbol uint32) {
	prob, total := ac.getProbability(symbol)

	rangeVal := ac.high - ac.low + 1
	ac.high = ac.low + (rangeVal * prob.high / total) - 1
	ac.low = ac.low + (rangeVal * prob.low / total)

	for {
		if ac.high < half {
			ac.putBit(0)
		} else if ac.low >= half {
			ac.putBit(1)
		} else if ac.low >= quarter && ac.high < 3*quarter {
			ac.pending++
			ac.low -= quarter
			ac.high -= quarter
		} else {
			break
		}
		ac.low = (ac.low << 1) & one
		ac.high = ((ac.high << 1) + 1) & one
	}
}

// Finalize encoding by flushing any remaining bits.
func (ac *Context) finalizeEncoding() {
	ac.pending++
	if ac.low < quarter {
		ac.putBit(0)
	} else {
		ac.putBit(1)
	}
	ac.flushRemainingBits()
}

// Initialize decoding.
func (ac *Context) initializeDecoding() uint32 {
	var value uint32 = 0
	for i := 0; i < codeBits; i++ {
		value <<= 1
		if ac.inputBit() == 1 {
			value++
		}
	}
	return value
}

// Decode a single symbol.
func (ac *Context) decodeSymbol() uint32 {
	rangeVal := ac.high - ac.low + 1
	offset := ((ac.value-ac.low+1)*ac.totalFrequency - 1) / rangeVal

	prob, total, symbol := ac.getChar(offset)

	if symbol == EOF || symbol == 0 {
		return symbol
	}

	ac.high = ac.low + (rangeVal*prob.high)/total - 1
	ac.low = ac.low + (rangeVal*prob.low)/total

	for {
		if ac.high < half {
			// Nothing to do
		} else if ac.low >= half {
			ac.value -= half
			ac.low -= half
			ac.high -= half
		} else if ac.low >= quarter && ac.high < 3*quarter {
			ac.value -= quarter
			ac.low -= quarter
			ac.high -= quarter
		} else {
			break
		}
		ac.low <<= 1
		ac.high = (ac.high << 1) + 1
		ac.value = (ac.value << 1) | ac.inputBit()
	}
	return symbol
}

type InputBuffer struct {
	lastMask      int
	remainingBits int
	currentByte   int
}

func NewInputBuffer() *InputBuffer {
	return &InputBuffer{
		remainingBits: codeBits,
		lastMask:      1,
	}
}

// Input a single bit.
func (ac *Context) inputBit() uint32 {
	if ac.lastMask == 1 {
		b, err := ac.buffer.ReadByte()
		if err != nil {
			if ac.remainingBits <= 0 {
				return EOF
			}
			ac.currentByte = -1
			ac.remainingBits -= 8
		} else {
			ac.currentByte = int(b)
		}
		ac.lastMask = 0x80
	} else {
		ac.lastMask >>= 1
	}

	if (ac.currentByte & ac.lastMask) != 0 {
		return 1
	}
	return 0
}

type OutputBuffer struct {
	nextByte int
	byteMask int
}

func NewOutputBuffer() *OutputBuffer {
	return &OutputBuffer{
		byteMask: 0x80,
	}
}

// Output a bit and flush any that are pending
func (ac *Context) putBit(bit uint8) {
	ac.outputBit(bit)
	for i := 0; i < ac.pending; i++ {
		ac.outputBit(1 - bit)
	}
	ac.pending = 0
}

// Output a single bit.
func (ac *Context) outputBit(bit uint8) {
	if bit == 1 {
		ac.nextByte |= ac.byteMask
	}
	ac.byteMask >>= 1
	if ac.byteMask == 0 {
		ac.buffer.WriteByte(byte(ac.nextByte))
		ac.nextByte = 0
		ac.byteMask = 0x80
	}
}

// Flush remaining bits (used in final encoding).
func (ac *Context) flushRemainingBits() {
	if ac.byteMask != 0x80 {
		ac.buffer.WriteByte(byte(ac.nextByte))
	}
}
