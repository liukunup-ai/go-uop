import React from 'react';
import { useDeviceStore } from '../stores/deviceStore';

export const DeviceList: React.FC = () => {
  const {
    devices,
    selectedDevice,
    connectedDevice,
    selectDevice,
    connectDevice,
    disconnectDevice,
    fetchDevices,
    isLoading,
  } = useDeviceStore();

  React.useEffect(() => {
    fetchDevices();
    const interval = setInterval(fetchDevices, 5000);
    return () => clearInterval(interval);
  }, []);

  const iosDevices = devices.filter((d) => d.platform === 'ios');
  const androidDevices = devices.filter((d) => d.platform === 'android');

  const handleConnect = async () => {
    if (!selectedDevice) return;
    await connectDevice(selectedDevice);
  };

  return (
    <div className="bg-white rounded-lg shadow p-4">
      <h2 className="text-lg font-semibold mb-4">设备列表</h2>

      <div className="mb-4">
        <h3 className="text-sm font-medium text-gray-500 mb-2">iOS</h3>
        <div className="space-y-2">
          {iosDevices.map((device) => (
            <div
              key={device.id}
              onClick={() => selectDevice(device)}
              className={`p-3 rounded border cursor-pointer ${
                selectedDevice?.id === device.id
                  ? 'border-blue-500 bg-blue-50'
                  : 'border-gray-200 hover:border-gray-300'
              }`}
            >
              <div className="flex items-center">
                <span className="text-lg mr-2">📱</span>
                <div>
                  <div className="font-medium">{device.name}</div>
                  <div className="text-xs text-gray-500">{device.serial}</div>
                </div>
                {connectedDevice?.id === device.id && (
                  <span className="ml-auto text-green-500 text-sm">已连接</span>
                )}
              </div>
            </div>
          ))}
          {iosDevices.length === 0 && (
            <p className="text-sm text-gray-400">未发现 iOS 设备</p>
          )}
        </div>
      </div>

      <div className="mb-4">
        <h3 className="text-sm font-medium text-gray-500 mb-2">Android</h3>
        <div className="space-y-2">
          {androidDevices.map((device) => (
            <div
              key={device.id}
              onClick={() => selectDevice(device)}
              className={`p-3 rounded border cursor-pointer ${
                selectedDevice?.id === device.id
                  ? 'border-blue-500 bg-blue-50'
                  : 'border-gray-200 hover:border-gray-300'
              }`}
            >
              <div className="flex items-center">
                <span className="text-lg mr-2">📱</span>
                <div>
                  <div className="font-medium">{device.name}</div>
                  <div className="text-xs text-gray-500">{device.serial}</div>
                </div>
                {connectedDevice?.id === device.id && (
                  <span className="ml-auto text-green-500 text-sm">已连接</span>
                )}
              </div>
            </div>
          ))}
          {androidDevices.length === 0 && (
            <p className="text-sm text-gray-400">未发现 Android 设备</p>
          )}
        </div>
      </div>

      <div className="flex gap-2">
        <button
          onClick={handleConnect}
          disabled={!selectedDevice || isLoading}
          className="flex-1 px-4 py-2 bg-blue-500 text-white rounded disabled:bg-gray-300"
        >
          连接设备
        </button>
        <button
          onClick={disconnectDevice}
          disabled={!connectedDevice || isLoading}
          className="flex-1 px-4 py-2 bg-red-500 text-white rounded disabled:bg-gray-300"
        >
          断开
        </button>
      </div>
    </div>
  );
};