import { useState, useRef, useEffect } from 'react';
import { useSerialStore } from '../stores/serialStore';

export function SerialTerminal() {
  const { logs, sendRaw, clearLogs, activeConnection, isLoading } = useSerialStore();
  const [input, setInput] = useState('');
  const [hexMode, setHexMode] = useState(false);
  const logsEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    logsEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [logs]);

  const handleSend = async () => {
    if (!input.trim()) return;
    await sendRaw(input);
    setInput('');
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  if (!activeConnection) {
    return (
      <div className="bg-white rounded-lg shadow p-4">
        <h2 className="text-lg font-semibold mb-4">Serial Terminal</h2>
        <div className="text-gray-500 text-center py-8">
          Connect to a serial port to start
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow p-4 flex flex-col h-96">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold">Serial Terminal</h2>
        <div className="flex gap-2">
          <button
            onClick={() => setHexMode(!hexMode)}
            className={`text-sm px-2 py-1 rounded ${hexMode ? 'bg-blue-100 text-blue-700' : 'bg-gray-100'}`}
          >
            {hexMode ? 'Hex Mode' : 'Text Mode'}
          </button>
          <button
            onClick={clearLogs}
            className="text-sm px-2 py-1 bg-gray-100 rounded hover:bg-gray-200"
          >
            Clear
          </button>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto border rounded p-2 bg-gray-50 font-mono text-sm mb-4">
        {logs.length === 0 ? (
          <div className="text-gray-400 text-center py-4">No data yet</div>
        ) : (
          logs.map((log) => (
            <div
              key={log.id}
              className={`py-1 ${
                log.direction === 'out'
                  ? 'text-blue-600'
                  : log.success === false
                  ? 'text-red-600'
                  : 'text-green-600'
              }`}
            >
              <span className="text-gray-400 text-xs">
                [{new Date(log.timestamp).toLocaleTimeString()}]
              </span>{' '}
              <span className="text-gray-600">{log.direction === 'out' ? '>>>' : '<<<'}</span>{' '}
              <span className="break-all">{log.data}</span>
            </div>
          ))
        )}
        <div ref={logsEndRef} />
      </div>

      <div className="flex gap-2">
        <input
          type="text"
          className="flex-1 border rounded px-3 py-2 font-mono"
          placeholder={hexMode ? 'Enter hex (e.g., 0D 0A)' : 'Enter command...'}
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          disabled={isLoading}
        />
        <button
          onClick={handleSend}
          disabled={isLoading || !input.trim()}
          className="bg-blue-500 text-white rounded px-4 py-2 hover:bg-blue-600 disabled:opacity-50"
        >
          Send
        </button>
      </div>
    </div>
  );
}
