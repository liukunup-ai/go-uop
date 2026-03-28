package serial

import (
	"fmt"
	"os"
	"sync"
)

// FileHandler 文件写入器
type FileHandler struct {
	Path string
	mu   sync.Mutex
}

// NewFileHandler 创建文件写入器
func NewFileHandler(path string) FileHandler {
	return FileHandler{Path: path}
}

// Handle 处理事件，写入文件
func (h FileHandler) Handle(e Event) {
	h.mu.Lock()
	defer h.mu.Unlock()

	f, err := os.OpenFile(h.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "serial: open file %s: %v\n", h.Path, err)
		return
	}
	defer f.Close()

	ts := e.Timestamp.Format("2006-01-02 15:04:05.000")
	line := fmt.Sprintf("[%s] %s\n", ts, string(e.Data))
	_, err = f.WriteString(line)
	if err != nil {
		fmt.Fprintf(os.Stderr, "serial: write file %s: %v\n", h.Path, err)
	}
}
