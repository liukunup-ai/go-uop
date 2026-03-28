package vision

type loftrMatcher struct {
	config *Config
}

func newLoFTRMatcher(cfg *Config) *loftrMatcher {
	return &loftrMatcher{config: cfg}
}

func (m *loftrMatcher) Name() string {
	return "loftr"
}

func (m *loftrMatcher) Find(screenshot, template []byte) ([]*MatchResult, error) {
	return nil, nil
}

func (m *loftrMatcher) DebugRender(screenshot []byte, results []*MatchResult) []byte {
	return nil
}
