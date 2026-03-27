package console

import "sync"

type HistoryManager struct {
	mu      sync.RWMutex
	history []CommandRecord
	maxSize int
}

func NewHistoryManager(maxSize int) *HistoryManager {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &HistoryManager{
		history: make([]CommandRecord, 0, maxSize),
		maxSize: maxSize,
	}
}

func (m *HistoryManager) Add(record *CommandRecord) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.history = append(m.history, *record)

	if len(m.history) > m.maxSize {
		m.history = m.history[len(m.history)-m.maxSize:]
	}
}

func (m *HistoryManager) GetAll() []CommandRecord {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]CommandRecord, len(m.history))
	copy(result, m.history)
	return result
}

func (m *HistoryManager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.history = m.history[:0]
}

func (m *HistoryManager) GetSelected(ids []string) []CommandRecord {
	m.mu.RLock()
	defer m.mu.RUnlock()

	idSet := make(map[string]bool)
	for _, id := range ids {
		idSet[id] = true
	}

	var result []CommandRecord
	for _, record := range m.history {
		if idSet[record.ID] {
			result = append(result, record)
		}
	}
	return result
}
