package console

import (
	"gopkg.in/yaml.v3"
)

type YamlFlow struct {
	Name  string        `yaml:"name"`
	AppID string        `yaml:"appId,omitempty"`
	Steps []YamlCommand `yaml:"steps"`
}

type YamlCommand struct {
	TapOn     *TapCommand   `yaml:"tapOn,omitempty"`
	InputText *InputCommand `yaml:"inputText,omitempty"`
	Swipe     *SwipeCommand `yaml:"swipe,omitempty"`
	Launch    *struct{}     `yaml:"launch,omitempty"`
	Terminate *struct{}     `yaml:"terminate,omitempty"`
}

type TapCommand struct {
	X    int    `yaml:"x,omitempty"`
	Y    int    `yaml:"y,omitempty"`
	Text string `yaml:"text,omitempty"`
}

type InputCommand struct {
	Text string `yaml:"text"`
}

type SwipeCommand struct {
	X1       int `yaml:"x1"`
	Y1       int `yaml:"y1"`
	X2       int `yaml:"x2"`
	Y2       int `yaml:"y2"`
	Duration int `yaml:"duration,omitempty"`
}

func ExportToYaml(records []CommandRecord, name string) ([]byte, error) {
	flow := YamlFlow{
		Name:  name,
		Steps: make([]YamlCommand, 0, len(records)),
	}

	for _, record := range records {
		if !record.Success {
			continue
		}

		cmd := YamlCommand{}
		switch record.Type {
		case "tap":
			x, _ := toInt(record.Params["x"])
			y, _ := toInt(record.Params["y"])
			cmd.TapOn = &TapCommand{X: x, Y: y}
		case "input":
			text, _ := toString(record.Params["text"])
			cmd.InputText = &InputCommand{Text: text}
		case "swipe":
			x1, _ := toInt(record.Params["x1"])
			y1, _ := toInt(record.Params["y1"])
			x2, _ := toInt(record.Params["x2"])
			y2, _ := toInt(record.Params["y2"])
			dur, _ := toInt(record.Params["duration"])
			cmd.Swipe = &SwipeCommand{X1: x1, Y1: y1, X2: x2, Y2: y2, Duration: dur}
		case "launch":
			cmd.Launch = &struct{}{}
		case "terminate":
			cmd.Terminate = &struct{}{}
		}

		flow.Steps = append(flow.Steps, cmd)
	}

	return yaml.Marshal(&flow)
}

func ExportToYamlString(records []CommandRecord, name string) (string, error) {
	data, err := ExportToYaml(records, name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
