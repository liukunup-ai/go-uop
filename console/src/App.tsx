import { DeviceList } from './components/DeviceList';
import { ScreenPreview } from './components/ScreenPreview';
import { CommandPanel } from './components/CommandPanel';
import { CommandHistory } from './components/CommandHistory';
import { useDeviceStore } from './stores/deviceStore';

function App() {
  const { exportYaml } = useDeviceStore();

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <h1 className="text-2xl font-bold text-gray-800">go-uop Console</h1>
          <button
            onClick={() => exportYaml()}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            导出全部 YAML
          </button>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 py-6">
        <div className="grid grid-cols-12 gap-6">
          <div className="col-span-3">
            <DeviceList />
          </div>

          <div className="col-span-5">
            <ScreenPreview />
          </div>

          <div className="col-span-4 space-y-4">
            <CommandPanel />
            <CommandHistory />
          </div>
        </div>
      </main>
    </div>
  );
}

export default App;