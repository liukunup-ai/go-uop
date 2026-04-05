package watcher

import (
	"os"
	"testing"
)

func TestParseWatcherConfig(t *testing.T) {
	yamlContent := `
watcher:
  enabled: true
  rules:
    - name: "permission popup"
      priority: 10
      match:
        type: text
        text: "允许"
      actions:
        - tapOn: {x: 500, y: 800}
    
    - name: "upgrade popup"
      priority: 20
      match:
        type: image
        template: "upgrade.png"
        threshold: 0.85
      actions:
        - ref: "dismiss-upgrade"
      retry: 3
`
	tmpfile, err := os.CreateTemp("", "watcher-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(yamlContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	engine, err := LoadWatcherConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if !engine.Enabled() {
		t.Error("engine should be enabled from config")
	}

	rules := engine.Rules()
	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules))
	}

	if rules[0].Name != "permission popup" {
		t.Error("rules should be sorted by priority")
	}

	if rules[0].Retry != 0 {
		t.Error("first rule should have 0 retry")
	}

	if rules[1].Retry != 3 {
		t.Errorf("second rule should have 3 retry, got %d", rules[1].Retry)
	}
}

func TestParseWatcherConfigCompound(t *testing.T) {
	yamlContent := `
watcher:
  enabled: true
  rules:
    - name: "compound popup"
      priority: 10
      match:
        type: compound
        operator: or
        conditions:
          - type: text
            text: "确定"
          - type: text
            text: "取消"
      actions:
        - tapOn: {x: 500, y: 800}
`
	tmpfile, err := os.CreateTemp("", "watcher-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(yamlContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	engine, err := LoadWatcherConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	rules := engine.Rules()
	if len(rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(rules))
	}

	compound, ok := rules[0].Match.(*CompoundMatch)
	if !ok {
		t.Fatal("expected CompoundMatch")
	}

	if compound.Operator != "or" {
		t.Errorf("expected 'or' operator, got '%s'", compound.Operator)
	}

	if len(compound.Conditions) != 2 {
		t.Errorf("expected 2 conditions, got %d", len(compound.Conditions))
	}
}
