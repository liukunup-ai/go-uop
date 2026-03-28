import { useState } from 'react';
import { useSerialStore } from '../stores/serialStore';

export function SerialCommandTable() {
  const { commands, activeConnection, sendByID, loadCommandTable, isLoading } = useSerialStore();
  const [yamlInput, setYamlInput] = useState('');
  const [showYamlInput, setShowYamlInput] = useState(false);

  const handleSendCommand = async (commandId: string) => {
    await sendByID(commandId);
  };

  const handleLoadYaml = async () => {
    if (!yamlInput.trim()) return;
    await loadCommandTable(undefined, yamlInput);
    setYamlInput('');
    setShowYamlInput(false);
  };

  const handleLoadFile = () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.yaml,.yml';
    input.onchange = async (e) => {
      const file = (e.target as HTMLInputElement).files?.[0];
      if (file) {
        await loadCommandTable(file.name);
      }
    };
    input.click();
  };

  if (!activeConnection) {
    return (
      <div className="bg-white rounded-lg shadow p-4">
        <h2 className="text-lg font-semibold mb-4">Command Table</h2>
        <div className="text-gray-500 text-center py-8">
          Connect to a serial port to view commands
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow p-4">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold">Command Table</h2>
        <div className="flex gap-2">
          <button
            onClick={() => setShowYamlInput(!showYamlInput)}
            className="text-sm px-2 py-1 bg-gray-100 rounded hover:bg-gray-200"
          >
            {showYamlInput ? 'Hide YAML Input' : 'Load YAML'}
          </button>
          <button
            onClick={handleLoadFile}
            className="text-sm px-2 py-1 bg-gray-100 rounded hover:bg-gray-200"
          >
            Load File
          </button>
        </div>
      </div>

      {showYamlInput && (
        <div className="mb-4 space-y-2">
          <textarea
            className="w-full border rounded p-2 font-mono text-sm h-32"
            placeholder="Paste YAML content here..."
            value={yamlInput}
            onChange={(e) => setYamlInput(e.target.value)}
          />
          <button
            onClick={handleLoadYaml}
            disabled={isLoading || !yamlInput.trim()}
            className="bg-blue-500 text-white rounded px-4 py-2 text-sm hover:bg-blue-600 disabled:opacity-50"
          >
            Load Command Table
          </button>
        </div>
      )}

      {commands.length === 0 ? (
        <div className="text-gray-500 text-center py-4">
          No commands loaded. Load a command table to see available commands.
        </div>
      ) : (
        <div className="space-y-2">
          {commands.map((cmd) => (
            <div
              key={cmd.id}
              className="flex items-center justify-between p-3 border rounded hover:bg-gray-50"
            >
              <div>
                <div className="font-medium text-sm">{cmd.name || cmd.id}</div>
                <div className="text-xs text-gray-500 font-mono">{cmd.command}</div>
                {cmd.log && (
                  <div className="text-xs text-gray-400 mt-1">
                    Expect: {cmd.log}
                  </div>
                )}
              </div>
              <button
                onClick={() => handleSendCommand(cmd.id)}
                disabled={isLoading}
                className="bg-green-500 text-white rounded px-3 py-1 text-sm hover:bg-green-600 disabled:opacity-50"
              >
                Send
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
