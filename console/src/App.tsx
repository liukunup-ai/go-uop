import { useState } from 'react';
import { DeviceList } from './components/DeviceList';
import { ScreenPreview } from './components/ScreenPreview';
import { CommandPanel } from './components/CommandPanel';
import { CommandHistory } from './components/CommandHistory';
import { SerialList } from './components/SerialList';
import { SerialTerminal } from './components/SerialTerminal';
import { SerialCommandTable } from './components/SerialCommandTable';
import { useDeviceStore } from './stores/deviceStore';

type Tab = 'devices' | 'serial';

function App() {
  const [tab, setTab] = useState<Tab>('devices');
  const { exportYaml } = useDeviceStore();

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-6">
            <h1 className="text-2xl font-bold text-gray-800">go-uop Console</h1>
            <nav className="flex gap-1">
              <button
                onClick={() => setTab('devices')}
                className={`px-4 py-2 rounded ${
                  tab === 'devices'
                    ? 'bg-blue-500 text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                设备
              </button>
              <button
                onClick={() => setTab('serial')}
                className={`px-4 py-2 rounded ${
                  tab === 'serial'
                    ? 'bg-blue-500 text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                串口
              </button>
            </nav>
          </div>
          {tab === 'devices' && (
            <button
              onClick={() => exportYaml()}
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              导出全部 YAML
            </button>
          )}
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 py-6">
        {tab === 'devices' && (
          <div className="grid grid-cols-12 gap-6">
            <div className="col-span-5">
              <ScreenPreview />
            </div>

            <div className="col-span-7 space-y-4">
              <DeviceList />
              <CommandPanel />
              <CommandHistory />
            </div>
          </div>
        )}

        {tab === 'serial' && (
          <div className="grid grid-cols-12 gap-6">
            <div className="col-span-3">
              <SerialList />
            </div>

            <div className="col-span-5">
              <SerialTerminal />
            </div>

            <div className="col-span-4">
              <SerialCommandTable />
            </div>
          </div>
        )}
      </main>
    </div>
  );
}

export default App;
