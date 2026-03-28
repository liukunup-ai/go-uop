package vision

type templateMatcher struct {
	config *Config
}

func newTemplateMatcher(cfg *Config) *templateMatcher {
	return &templateMatcher{config: cfg}
}

func (m *templateMatcher) Name() string {
	return "template"
}

func (m *templateMatcher) Find(screenshot, template []byte) ([]*MatchResult, error) {
	return nil, nil
}

func (m *templateMatcher) DebugRender(screenshot []byte, results []*MatchResult) []byte {
	return nil
}
