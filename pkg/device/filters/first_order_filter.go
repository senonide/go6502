package filters

type FirstOrderFilter struct {
	B0    float32
	B1    float32
	A1    float32
	prevX float32
	prevY float32
}

func (f *FirstOrderFilter) Step(x float32) float32 {
	y := f.B0*x + f.B1*f.prevX - f.A1*f.prevY
	f.prevY = y
	f.prevX = x
	return y
}
