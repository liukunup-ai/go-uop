package diff

type DiffResult struct {
	HasDiff    bool
	Diffs      []DiffRegion
	OutputPath string
	Similarity float64
}

type DiffRegion struct {
	X, Y, Width, Height int
	Score               float64
	PixelCount          int
}
