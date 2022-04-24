package apu

import (
	"github.com/se-nonide/go6502/pkg/device/filters"
)

type APU struct {
	channel     chan float32
	sampleRate  float64
	pulse1      Pulse
	pulse2      Pulse
	triangle    Triangle
	noise       Noise
	dmc         DMC
	cycle       uint64
	framePeriod byte
	frameValue  byte
	frameIRQ    bool
	filterChain filters.FilterChain
}
