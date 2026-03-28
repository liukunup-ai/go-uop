package vision

type multiscaleMatcher struct {
	config *Config
}

func newMultiscaleMatcher(cfg *Config) *multiscaleMatcher {
	return &multiscaleMatcher{config: cfg}
}

func (m *multiscaleMatcher) Name() string {
	return "multiscale"
}

func (m *multiscaleMatcher) Find(screenshot, template []byte) ([]*MatchResult, error) {
	return nil, nil
}

func (m *multiscaleMatcher) DebugRender(screenshot []byte, results []*MatchResult) []byte {
	return nil
}
