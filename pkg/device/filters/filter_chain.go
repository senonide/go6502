package filters

type FilterChain []Filter

func (fc FilterChain) Step(x float32) float32 {
	for i := range fc {
		x = fc[i].Step(x)
	}
	return x
}
