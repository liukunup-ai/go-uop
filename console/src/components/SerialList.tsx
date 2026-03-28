import { useSerialStore } from '../stores/serialStore';

export function SerialList() {
  const { ports, config, setConfig, connect, activeConnection, disconnect, isLoading, error } = useSerialStore();

  const handleConnect = async () => {
    if (!config.name) {
      alert('Please select a port');
      return;
    }
    await connect(config);
  };

  const handleDisconnect = async () => {
    if (activeConnection) {
      await disconnect(activeConnection.id);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow p-4">
      <h2 className="text-lg font-semibold mb-4">Serial Ports</h2>

      {error && (
        <div className="mb-4 p-2 bg-red-100 text-red-700 rounded">{error}</div>
      )}

      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-1">Port</label>
          <select
            className="w-full border rounded px-3 py-2"
            value={config.name}
            onChange={(e) => setConfig({ name: e.target.value })}
            disabled={!!activeConnection}
          >
            <option value="">Select port...</option>
            {ports.map((port) => (
              <option key={port.name} value={port.name}>
                {port.name}
              </option>
            ))}
          </select>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-1">Baud Rate</label>
            <select
              className="w-full border rounded px-3 py-2"
              value={config.baud}
              onChange={(e) => setConfig({ baud: parseInt(e.target.value) })}
              disabled={!!activeConnection}
            >
              <option value={9600}>9600</option>
              <option value={19200}>19200</option>
              <option value={38400}>38400</option>
              <option value={57600}>57600</option>
              <option value={115200}>115200</option>
              <option value={230400}>230400</option>
              <option value={460800}>460800</option>
              <option value={921600}>921600</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Data Bits</label>
            <select
              className="w-full border rounded px-3 py-2"
              value={config.dataBits}
              onChange={(e) => setConfig({ dataBits: parseInt(e.target.value) })}
              disabled={!!activeConnection}
            >
              <option value={5}>5</option>
              <option value={6}>6</option>
              <option value={7}>7</option>
              <option value={8}>8</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Parity</label>
            <select
              className="w-full border rounded px-3 py-2"
              value={config.parity}
              onChange={(e) => setConfig({ parity: e.target.value })}
              disabled={!!activeConnection}
            >
              <option value="N">None</option>
              <option value="O">Odd</option>
              <option value="E">Even</option>
              <option value="M">Mark</option>
              <option value="S">Space</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Stop Bits</label>
            <select
              className="w-full border rounded px-3 py-2"
              value={config.stopBits}
              onChange={(e) => setConfig({ stopBits: parseInt(e.target.value) })}
              disabled={!!activeConnection}
            >
              <option value={1}>1</option>
              <option value={2}>2</option>
            </select>
          </div>
        </div>

        <div className="flex gap-2">
          {!activeConnection ? (
            <button
              onClick={handleConnect}
              disabled={isLoading || !config.name}
              className="flex-1 bg-blue-500 text-white rounded py-2 px-4 hover:bg-blue-600 disabled:opacity-50"
            >
              {isLoading ? 'Connecting...' : 'Connect'}
            </button>
          ) : (
            <button
              onClick={handleDisconnect}
              disabled={isLoading}
              className="flex-1 bg-red-500 text-white rounded py-2 px-4 hover:bg-red-600 disabled:opacity-50"
            >
              {isLoading ? 'Disconnecting...' : 'Disconnect'}
            </button>
          )}
          <button
            onClick={() => useSerialStore.getState().fetchPorts()}
            className="bg-gray-200 rounded py-2 px-4 hover:bg-gray-300"
          >
            Refresh
          </button>
        </div>

        {activeConnection && (
          <div className="mt-4 p-3 bg-green-100 rounded">
            <p className="text-sm text-green-800">
              Connected to {activeConnection.config.name} @ {activeConnection.config.baud} baud
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
