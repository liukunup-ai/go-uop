package serial

import (
	"os"
	"sync"
	"testing"
	"time"
)

func TestMonitor_MatchOnce(t *testing.T) {
	m := NewMonitor()

	count := 0
	m.AddRule("OK", MatchOnce, func(e Event) {
		count++
	})

	// 第一次匹配，触发
	m.OnData(Event{Data: []byte("OK")})
	if count != 1 {
		t.Errorf("expected count=1, got %d", count)
	}

	// 第二次匹配，不触发（已禁用）
	m.OnData(Event{Data: []byte("OK")})
	if count != 1 {
		t.Errorf("expected count=1 after second match, got %d", count)
	}
}

func TestMonitor_MatchContinuous(t *testing.T) {
	m := NewMonitor()

	count := 0
	m.AddRule("OK", MatchContinuous, func(e Event) {
		count++
	})

	m.OnData(Event{Data: []byte("OK")})
	m.OnData(Event{Data: []byte("OK")})
	m.OnData(Event{Data: []byte("OK")})

	if count != 3 {
		t.Errorf("expected count=3, got %d", count)
	}
}

func TestMonitor_MatchRateLimited(t *testing.T) {
	m := NewMonitor()

	count := 0
	m.AddRateLimitedRule("OK", 100*time.Millisecond, func(e Event) {
		count++
	})

	m.OnData(Event{Data: []byte("OK")})
	m.OnData(Event{Data: []byte("OK")}) // 应被限频

	if count != 1 {
		t.Errorf("expected count=1, got %d", count)
	}

	time.Sleep(150 * time.Millisecond)
	m.OnData(Event{Data: []byte("OK")}) // 超过间隔，再次触发

	if count != 2 {
		t.Errorf("expected count=2, got %d", count)
	}
}

func TestMonitor_EnableDisable(t *testing.T) {
	m := NewMonitor()

	count := 0
	m.AddRule("OK", MatchOnce, func(e Event) {
		count++
	})

	m.OnData(Event{Data: []byte("OK")})
	if count != 1 {
		t.Errorf("expected count=1, got %d", count)
	}

	// 禁用后无法触发
	m.DisableRule("OK")
	m.OnData(Event{Data: []byte("OK")})
	if count != 1 {
		t.Errorf("expected count=1 after disable, got %d", count)
	}

	// 启用后恢复
	m.EnableRule("OK")
	m.OnData(Event{Data: []byte("OK")})
	if count != 2 {
		t.Errorf("expected count=2 after re-enable, got %d", count)
	}
}

func TestMonitor_Concurrent(t *testing.T) {
	m := NewMonitor()

	var mu sync.Mutex
	count := 0

	m.AddRule("OK", MatchContinuous, func(e Event) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				m.OnData(Event{Data: []byte("OK")})
			}
		}()
	}

	wg.Wait()

	mu.Lock()
	expected := 1000
	if count != expected {
		t.Errorf("expected count=%d, got %d", expected, count)
	}
	mu.Unlock()
}

func TestFileHandler(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "serial_test_*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	handler := NewFileHandler(tmpfile.Name())

	event := Event{
		Data:      []byte("test message"),
		Timestamp: time.Now(),
	}

	handler.Handle(event)

	// 验证文件内容
	data, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Error("expected file to have content")
	}
}

func TestMonitor_NotMatched(t *testing.T) {
	m := NewMonitor()

	count := 0
	m.AddRule("OK", MatchOnce, func(e Event) {
		count++
	})

	// 不匹配关键字
	m.OnData(Event{Data: []byte("NOT_MATCH")})

	if count != 0 {
		t.Errorf("expected count=0, got %d", count)
	}
}
