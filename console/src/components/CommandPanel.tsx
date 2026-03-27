import React, { useState } from 'react';
import { useDeviceStore } from '../stores/deviceStore';

interface CommandButton {
  label: string;
  command: string;
  icon: string;
  fields?: { name: string; label: string; type: string; default?: any }[];
}

const commands: CommandButton[] = [
  { label: 'Tap', command: 'tap', icon: '👆', fields: [
    { name: 'x', label: 'X', type: 'number' },
    { name: 'y', label: 'Y', type: 'number' },
  ]},
  { label: 'Input', command: 'input', icon: '⌨️', fields: [
    { name: 'text', label: '文本', type: 'text' },
  ]},
  { label: 'Swipe', command: 'swipe', icon: '👆', fields: [
    { name: 'x1', label: 'X1', type: 'number' },
    { name: 'y1', label: 'Y1', type: 'number' },
    { name: 'x2', label: 'X2', type: 'number' },
    { name: 'y2', label: 'Y2', type: 'number' },
  ]},
  { label: 'Launch', command: 'launch', icon: '🚀', fields: [] },
  { label: 'Terminate', command: 'terminate', icon: '❌', fields: [] },
];

export const CommandPanel: React.FC = () => {
  const { connectedDevice, executeCommand, error } = useDeviceStore();
  const [selectedCommand, setSelectedCommand] = useState<string | null>(null);
  const [params, setParams] = useState<Record<string, any>>({});

  const handleCommandSelect = (cmd: CommandButton) => {
    setSelectedCommand(cmd.command);
    const defaults: Record<string, any> = {};
    cmd.fields?.forEach((f) => {
      defaults[f.name] = f.default || '';
    });
    setParams(defaults);
  };

  const handleExecute = async () => {
    if (!selectedCommand) return;
    await executeCommand(selectedCommand, params);
    setSelectedCommand(null);
    setParams({});
  };

  const selected = commands.find((c) => c.command === selectedCommand);

  return (
    <div className="bg-white rounded-lg shadow p-4">
      <h2 className="text-lg font-semibold mb-4">命令面板</h2>

      <div className="grid grid-cols-2 gap-2 mb-4">
        {commands.map((cmd) => (
          <button
            key={cmd.command}
            onClick={() => handleCommandSelect(cmd)}
            disabled={!connectedDevice}
            className={`px-4 py-2 rounded border text-left disabled:opacity-50 ${
              selectedCommand === cmd.command
                ? 'border-blue-500 bg-blue-50'
                : 'border-gray-200 hover:border-gray-300'
            }`}
          >
            <span className="mr-2">{cmd.icon}</span>
            {cmd.label}
          </button>
        ))}
      </div>

      {selected && selected.fields && selected.fields.length > 0 && (
        <div className="space-y-3 mb-4">
          {selected.fields.map((field) => (
            <div key={field.name}>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                {field.label}
              </label>
              <input
                type={field.type}
                value={params[field.name] || ''}
                onChange={(e) =>
                  setParams({
                    ...params,
                    [field.name]: field.type === 'number' ? Number(e.target.value) : e.target.value,
                  })
                }
                className="w-full px-3 py-2 border rounded"
              />
            </div>
          ))}
          <button
            onClick={handleExecute}
            disabled={!connectedDevice}
            className="w-full px-4 py-2 bg-green-500 text-white rounded disabled:bg-gray-300"
          >
            执行
          </button>
        </div>
      )}

      {selected && (!selected.fields || selected.fields.length === 0) && (
        <button
          onClick={handleExecute}
          disabled={!connectedDevice}
          className="w-full px-4 py-2 bg-green-500 text-white rounded disabled:bg-gray-300"
        >
          执行 {selected.label}
        </button>
      )}

      {error && (
        <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded text-red-600 text-sm">
          {error}
        </div>
      )}
    </div>
  );
};