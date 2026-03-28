package vision

type siftMatcher struct {
	config *Config
}

func newSIFTMatcher(cfg *Config) *siftMatcher {
	return &siftMatcher{config: cfg}
}

func (m *siftMatcher) Name() string {
	return "sift"
}

func (m *siftMatcher) Find(screenshot, template []byte) ([]*MatchResult, error) {
	return nil, nil
}

func (m *siftMatcher) DebugRender(screenshot []byte, results []*MatchResult) []byte {
	return nil
}
