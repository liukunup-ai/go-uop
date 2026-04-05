package watcher

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type watcherConfig struct {
	Watcher struct {
		Enabled bool         `yaml:"enabled"`
		Rules   []ruleConfig `yaml:"rules"`
	} `yaml:"watcher"`
}

type ruleConfig struct {
	Name     string      `yaml:"name"`
	Priority int         `yaml:"priority"`
	Match    matchConfig `yaml:"match"`
	Actions  []any       `yaml:"actions"`
	Retry    int         `yaml:"retry"`
}

type matchConfig struct {
	Type       string        `yaml:"type"`
	Text       string        `yaml:"text"`
	Pattern    string        `yaml:"pattern"`
	Template   string        `yaml:"template"`
	Threshold  float64       `yaml:"threshold"`
	Operator   string        `yaml:"operator"`
	Conditions []matchConfig `yaml:"conditions"`
}

func LoadWatcherConfig(path string) (*WatcherEngine, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg watcherConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	engine := NewWatcherEngine()
	if cfg.Watcher.Enabled {
		engine.Enable()
	}

	for _, rc := range cfg.Watcher.Rules {
		match, err := parseMatchConfig(rc.Match)
		if err != nil {
			return nil, fmt.Errorf("parse match config for rule %s: %w", rc.Name, err)
		}

		actions, err := parseActionsConfig(rc.Actions)
		if err != nil {
			return nil, fmt.Errorf("parse actions config for rule %s: %w", rc.Name, err)
		}

		engine.AddRule(Rule{
			Name:     rc.Name,
			Priority: rc.Priority,
			Match:    match,
			Actions:  actions,
			Retry:    rc.Retry,
		})
	}

	return engine, nil
}

func parseMatchConfig(cfg matchConfig) (MatchCondition, error) {
	switch cfg.Type {
	case "text":
		return NewTextMatch(cfg.Text), nil

	case "regex":
		return NewRegexMatch(cfg.Pattern), nil

	case "image":
		threshold := cfg.Threshold
		if threshold == 0 {
			threshold = 0.8
		}
		return NewImageMatch(cfg.Template, threshold), nil

	case "compound":
		operator := cfg.Operator
		if operator == "" {
			operator = "and"
		}

		conditions := make([]MatchCondition, 0, len(cfg.Conditions))
		for _, condCfg := range cfg.Conditions {
			cond, err := parseMatchConfig(condCfg)
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, cond)
		}

		return NewCompoundMatch(operator, conditions), nil

	default:
		return nil, fmt.Errorf("unknown match type: %s", cfg.Type)
	}
}

func parseActionsConfig(actionsCfg []any) ([]Action, error) {
	actions := make([]Action, 0, len(actionsCfg))

	for _, actionCfg := range actionsCfg {
		switch a := actionCfg.(type) {
		case map[string]any:
			for name, args := range a {
				var argsMap map[string]any
				if args != nil {
					var ok bool
					argsMap, ok = args.(map[string]any)
					if !ok {
						argsMap = nil
					}
				}
				actions = append(actions, NewInlineCommand(name, argsMap))
			}

		case string:
			actions = append(actions, NewReferenceFlow(a))

		default:
			return nil, fmt.Errorf("unsupported action type: %T", actionCfg)
		}
	}

	return actions, nil
}
