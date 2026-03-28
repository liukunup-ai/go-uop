package serial

import (
	"os"
	"testing"
)

func TestCommandTable_New(t *testing.T) {
	ct := NewCommandTable()
	if ct == nil {
		t.Fatal("expected non-nil CommandTable")
	}
}

func TestCommandTable_LoadFromFile(t *testing.T) {
	// 创建临时 YAML 文件
	content := `
commands:
  - id: reset
    name: 设备复位
    command: "AT+RST\r\n"
    log: "(?i)ok"
    timeout: 5s
  - id: version
    name: 查询版本
    command: "AT+VERSION?\r\n"
    timeout: 3s
`
	tmpfile, err := os.CreateTemp("", "commands_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	ct := NewCommandTable()
	if err := ct.LoadFromFile(tmpfile.Name()); err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	// 测试 GetByID
	cmd, ok := ct.GetByID("reset")
	if !ok {
		t.Fatal("expected to find command 'reset'")
	}
	if cmd.Name != "设备复位" {
		t.Errorf("expected name '设备复位', got '%s'", cmd.Name)
	}
	if cmd.Command != "AT+RST\r\n" {
		t.Errorf("expected command 'AT+RST\\r\\n', got '%s'", cmd.Command)
	}

	// 测试 GetByName
	cmd, ok = ct.GetByName("设备复位")
	if !ok {
		t.Fatal("expected to find command by name '设备复位'")
	}
	if cmd.ID != "reset" {
		t.Errorf("expected id 'reset', got '%s'", cmd.ID)
	}

	// 测试不存在的命令
	_, ok = ct.GetByID("notexist")
	if ok {
		t.Error("expected not to find 'notexist'")
	}
}

func TestCommandTable_GetByDefaultName(t *testing.T) {
	// 创建临时 YAML 文件
	content := `
commands:
  - id: screenshot
    name: 截屏
    default_name: screenshot
    command: "AT+SCREENSHOT\r\n"
    log: "OK"
    timeout: 5s
  - id: reset
    name: 复位
    default_name: reset
    command: "AT+RST\r\n"
    timeout: 3s
`
	tmpfile, err := os.CreateTemp("", "commands_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	ct := NewCommandTable()
	if err := ct.LoadFromFile(tmpfile.Name()); err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	// 测试 GetByDefaultName
	cmd, ok := ct.GetByDefaultName("screenshot")
	if !ok {
		t.Fatal("expected to find command 'screenshot'")
	}
	if cmd.ID != "screenshot" {
		t.Errorf("expected id 'screenshot', got '%s'", cmd.ID)
	}
	if cmd.Name != "截屏" {
		t.Errorf("expected name '截屏', got '%s'", cmd.Name)
	}

	// 测试不存在的 default name
	_, ok = ct.GetByDefaultName("notexist")
	if ok {
		t.Error("expected not to find 'notexist'")
	}
}

func TestCommandTable_List(t *testing.T) {
	ct := NewCommandTable()
	// 直接添加命令进行测试
	ct.mu.Lock()
	ct.commands["cmd1"] = &Command{ID: "cmd1", Name: "Command 1"}
	ct.commands["cmd2"] = &Command{ID: "cmd2", Name: "Command 2"}
	ct.mu.Unlock()

	cmds := ct.List()
	if len(cmds) != 2 {
		t.Errorf("expected 2 commands, got %d", len(cmds))
	}
}

func TestCommandTable_LoadFromFile_Invalid(t *testing.T) {
	// 测试不存在的文件
	ct := NewCommandTable()
	err := ct.LoadFromFile("/not/exist/file.yaml")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestCommandTable_LoadFromFile_MissingID(t *testing.T) {
	// 创建缺少 ID 的 YAML 文件
	content := `
commands:
  - name: 设备复位
    command: "AT+RST\\r\\n"
`
	tmpfile, err := os.CreateTemp("", "commands_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	ct := NewCommandTable()
	err = ct.LoadFromFile(tmpfile.Name())
	if err == nil {
		t.Error("expected error for missing id")
	}
}
