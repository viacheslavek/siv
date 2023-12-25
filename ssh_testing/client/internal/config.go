package internal

type Waterfall struct {
	FlowSize      int
	WaterfallSize int
	MaxAmplitude  int
	Period        int
	BatchCount    int
}

type Flashing struct {
	FlowSize     int
	MaxAmplitude int
	Period       int
	BatchCount   int
}
