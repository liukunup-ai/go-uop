import React, { useState, useEffect } from 'react';
import { useDeviceStore } from '../stores/deviceStore';

export const CommandHistory: React.FC = () => {
  const { history, fetchHistory, clearHistory, exportYaml, connectedDevice } = useDeviceStore();
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());

  useEffect(() => {
    if (connectedDevice) {
      fetchHistory();
    }
  }, [connectedDevice]);

  const toggleSelect = (id: string) => {
    const newSet = new Set(selectedIds);
    if (newSet.has(id)) {
      newSet.delete(id);
    } else {
      newSet.add(id);
    }
    setSelectedIds(newSet);
  };

  const handleExport = () => {
    exportYaml(selectedIds.size > 0 ? Array.from(selectedIds) : undefined);
  };

  return (
    <div className="bg-white rounded-lg shadow p-4">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold">命令历史</h2>
        <div className="flex gap-2">
          <button
            onClick={handleExport}
            disabled={history.length === 0}
            className="px-3 py-1 text-sm bg-blue-100 text-blue-600 rounded disabled:bg-gray-100"
          >
            导出 YAML
          </button>
          <button
            onClick={clearHistory}
            disabled={history.length === 0}
            className="px-3 py-1 text-sm bg-red-100 text-red-600 rounded disabled:bg-gray-100"
          >
            清空
          </button>
        </div>
      </div>

      <div className="space-y-2 max-h-96 overflow-y-auto">
        {history.length === 0 ? (
          <p className="text-gray-400 text-center py-4">暂无命令历史</p>
        ) : (
          [...history].reverse().map((record) => (
            <div
              key={record.id}
              onClick={() => toggleSelect(record.id)}
              className={`p-3 rounded border cursor-pointer ${
                selectedIds.has(record.id)
                  ? 'border-blue-500 bg-blue-50'
                  : 'border-gray-200 hover:border-gray-300'
              }`}
            >
              <div className="flex items-center justify-between">
                <span className={`font-medium ${record.success ? 'text-green-600' : 'text-red-600'}`}>
                  {record.command.toUpperCase()}
                </span>
                <span className="text-xs text-gray-400">
                  {new Date(record.timestamp).toLocaleTimeString()}
                </span>
              </div>
              <div className="text-sm text-gray-600 mt-1">
                {formatParams(record.params)}
              </div>
              {record.duration && (
                <div className="text-xs text-gray-400 mt-1">
                  耗时: {record.duration}
                </div>
              )}
            </div>
          ))
        )}
      </div>

      {selectedIds.size > 0 && (
        <div className="mt-4 p-2 bg-blue-50 rounded text-sm text-blue-600">
          已选择 {selectedIds.size} 条命令
        </div>
      )}
    </div>
  );
};

function formatParams(params: Record<string, any>): string {
  const entries = Object.entries(params);
  if (entries.length === 0) return '-';
  return entries.map(([k, v]) => `${k}: ${v}`).join(', ');
}