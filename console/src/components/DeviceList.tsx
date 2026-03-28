import React, { useState } from 'react';
import { useDeviceStore } from '../stores/deviceStore';
import { useIOSStore } from '../stores/iosStore';

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

  const {
    devices: iosDevices,
    selectedDevice: selectedIOSDevice,
    selectDevice: selectIOSDevice,
    setupForward,
    startWDA,
    stopWDA,
    fetchDevices: fetchIOSDevices,
    isLoading: iosLoading,
    forwardedPorts,
    forwardCommand,
    wdaCommand,
    error,
  } = useIOSStore();

  const [showWdaModal, setShowWdaModal] = useState(false);
  const [bundleId, setBundleId] = useState('com.facebook.WebDriverAgentRunner.xctrunner');

  React.useEffect(() => {
    fetchDevices();
    fetchIOSDevices();
    const interval = setInterval(() => {
      fetchDevices();
      fetchIOSDevices();
    }, 10000);
    return () => clearInterval(interval);
  }, []);

  React.useEffect(() => {
    console.log('selectedIOSDevice changed:', selectedIOSDevice?.udid);
  }, [selectedIOSDevice]);

  const androidDevices = devices.filter((d) => d.platform === 'android');

  const handleIOSClick = (device: typeof iosDevices[0]) => {
    console.log('iOS device clicked:', device.udid);
    selectIOSDevice(device);
    selectDevice(null);
  };

  const handleAndroidClick = (device: typeof androidDevices[0]) => {
    selectDevice(device);
    selectIOSDevice(null);
  };

  const handle8100Forward = async (udid: string) => {
    try {
      await setupForward(udid, [8100]);
    } catch (e) {
      console.error('8100 forward failed:', e);
    }
  };

  const handle9100Forward = async (udid: string) => {
    console.log('9100 forward clicked, udid:', udid);
    if (!udid) {
      console.error('UDID is empty!');
      return;
    }
    try {
      await setupForward(udid, [9100]);
    } catch (e) {
      console.error('9100 forward failed:', e);
    }
  };

  const tunnelCommand = 'sudo ios tunnel start';

  return (
    <div className="bg-white rounded-lg shadow p-4">
      <h2 className="text-lg font-semibold mb-4">设备列表</h2>

      {error && (
        <div className="mb-3 p-2 bg-red-50 border border-red-200 rounded text-xs text-red-600">
          {error}
        </div>
      )}

      <div className="mb-4">
        <h3 className="text-sm font-medium text-gray-500 mb-2">iOS</h3>
        <div className="mb-3 p-2 bg-yellow-50 border border-yellow-200 rounded text-xs">
          <div className="font-medium text-yellow-700 mb-1">首次使用需要启动 tunnel：</div>
          <code className="bg-yellow-100 px-1 py-0.5 rounded">{tunnelCommand}</code>
          <span className="text-yellow-600 ml-2">在终端执行</span>
        </div>
        <div className="space-y-2">
          {iosDevices.map((device) => (
            <div
              key={device.udid}
              onClick={() => handleIOSClick(device)}
              className={`p-3 rounded border cursor-pointer ${
                selectedIOSDevice?.udid === device.udid
                  ? 'border-blue-500 bg-blue-50'
                  : 'border-gray-200 hover:border-gray-300'
              }`}
            >
              <div className="flex items-center">
                <span className="text-lg mr-2">📱</span>
                <div>
                  <div className="font-medium">{device.name || 'iOS Device'}</div>
                  <div className="text-xs text-gray-500">{device.udid}</div>
                  {device.iosVersion && (
                    <div className="text-xs text-gray-400">iOS {device.iosVersion}</div>
                  )}
                </div>
                {device.status === 'wda_running' && (
                  <span className="ml-auto text-green-500 text-xs">WDA 运行中</span>
                )}
                {device.status === 'forwarding' && (
                  <span className="ml-auto text-blue-500 text-xs">转发中</span>
                )}
              </div>

              {selectedIOSDevice?.udid === device.udid && (
                <div className="mt-3 pt-2 border-t border-gray-200">
                  <div className="flex gap-2">
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        handle8100Forward(device.udid);
                      }}
                      disabled={iosLoading}
                      className="flex-1 px-2 py-1.5 bg-blue-500 text-white text-xs rounded hover:bg-blue-600 disabled:bg-gray-300"
                    >
                      {forwardedPorts.includes(8100) ? '✓ 8100' : '8100'}
                    </button>
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        handle9100Forward(device.udid);
                      }}
                      disabled={iosLoading}
                      className="flex-1 px-2 py-1.5 bg-blue-500 text-white text-xs rounded hover:bg-blue-600 disabled:bg-gray-300"
                    >
                      {forwardedPorts.includes(9100) ? '✓ 9100' : '9100'}
                    </button>
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        if (device.status === 'wda_running') {
                          stopWDA(device.udid);
                        } else {
                          setBundleId('com.facebook.WebDriverAgentRunner.xctrunner');
                          setShowWdaModal(true);
                        }
                      }}
                      disabled={iosLoading}
                      className={`flex-1 px-2 py-1.5 text-white text-xs rounded disabled:bg-gray-300 ${
                        device.status === 'wda_running'
                          ? 'bg-red-500 hover:bg-red-600'
                          : 'bg-green-500 hover:bg-green-600'
                      }`}
                    >
                      {device.status === 'wda_running' ? '停止WDA' : '启动WDA'}
                    </button>
                  </div>

                  {(forwardCommand || wdaCommand) && (
                    <div className="mt-2 p-2 bg-gray-50 rounded text-xs">
                      {forwardCommand && (
                        <div className="mb-1">
                          <div className="text-gray-500">端口转发命令：</div>
                          <code className="text-blue-600">{forwardCommand}</code>
                        </div>
                      )}
                      {wdaCommand && (
                        <div>
                          <div className="text-gray-500">WDA 启动命令：</div>
                          <code className="text-green-600">{wdaCommand}</code>
                        </div>
                      )}
                    </div>
                  )}
                </div>
              )}
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
              onClick={() => handleAndroidClick(device)}
              className={`p-3 rounded border cursor-pointer ${
                selectedDevice?.id === device.id
                  ? 'border-blue-500 bg-blue-50'
                  : 'border-gray-200 hover:border-gray-300'
              }`}
            >
              <div className="flex items-center">
                <span className="text-lg mr-2">🤖</span>
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
          onClick={() => {
            if (selectedIOSDevice) {
              handle9100Forward(selectedIOSDevice.udid);
            } else if (selectedDevice) {
              connectDevice(selectedDevice);
            }
          }}
          disabled={(!selectedDevice && !selectedIOSDevice) || isLoading || iosLoading}
          className="flex-1 px-4 py-2 bg-blue-500 text-white rounded disabled:bg-gray-300"
        >
          {selectedIOSDevice ? '预览屏幕' : '连接设备'}
        </button>
        <button
          onClick={disconnectDevice}
          disabled={!connectedDevice || isLoading}
          className="flex-1 px-4 py-2 bg-red-500 text-white rounded disabled:bg-gray-300"
        >
          断开
        </button>
      </div>

      {showWdaModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-lg p-6 w-[480px]">
            <h3 className="text-lg font-semibold mb-4">启动 WDA</h3>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Bundle ID
              </label>
              <input
                type="text"
                value={bundleId}
                onChange={(e) => setBundleId(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
                placeholder="com.facebook.WebDriverAgentRunner.xctrunner"
              />
            </div>
            <div className="flex gap-2 justify-end">
              <button
                onClick={() => setShowWdaModal(false)}
                className="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded"
              >
                取消
              </button>
              <button
                onClick={async () => {
                  if (selectedIOSDevice) {
                    await startWDA(selectedIOSDevice.udid, bundleId);
                    setShowWdaModal(false);
                  }
                }}
                className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
              >
                启动
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
