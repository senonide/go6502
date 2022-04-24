package filters

import "math"

type Filter interface {
	Step(x float32) float32
}

func LowPassFilter(sampleRate float32, cutoffFreq float32) Filter {
	c := sampleRate / math.Pi / cutoffFreq
	a0i := 1 / (1 + c)
	return &FirstOrderFilter{
		B0: a0i,
		B1: a0i,
		A1: (1 - c) * a0i,
	}
}

func HighPassFilter(sampleRate float32, cutoffFreq float32) Filter {
	c := sampleRate / math.Pi / cutoffFreq
	a0i := 1 / (1 + c)
	return &FirstOrderFilter{
		B0: c * a0i,
		B1: -c * a0i,
		A1: (1 - c) * a0i,
	}
}
