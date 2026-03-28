import React, { useRef, useState, useEffect } from 'react';
import { useDeviceStore } from '../stores/deviceStore';
import { useIOSStore, getIosPreviewUrl } from '../stores/iosStore';

export const ScreenPreview: React.FC = () => {
  const { connectedDevice, screenshot, fetchScreenshot } = useDeviceStore();
  const { selectedDevice: iosDevice } = useIOSStore();
  const [mousePos, setMousePos] = useState<{ x: number; y: number } | null>(null);
  const [refreshKey, setRefreshKey] = useState(0);
  const imgRef = useRef<HTMLImageElement>(null);

  const isIOSPreview = iosDevice && getIosPreviewUrl();

  useEffect(() => {
    const interval = setInterval(() => {
      if (connectedDevice) {
        fetchScreenshot();
      }
      if (isIOSPreview) {
        setRefreshKey(k => k + 1);
      }
    }, 3000);
    return () => clearInterval(interval);
  }, [connectedDevice, isIOSPreview]);

  const handleMouseMove = (e: React.MouseEvent) => {
    if (!imgRef.current) return;
    const rect = imgRef.current.getBoundingClientRect();
    const scaleX = imgRef.current.naturalWidth / rect.width;
    const scaleY = imgRef.current.naturalHeight / rect.height;
    setMousePos({
      x: Math.round((e.clientX - rect.left) * scaleX),
      y: Math.round((e.clientY - rect.top) * scaleY),
    });
  };

  return (
    <div className="bg-white rounded-lg shadow p-4">
      <h2 className="text-lg font-semibold mb-4">屏幕预览</h2>

      {!connectedDevice && !isIOSPreview ? (
        <div className="aspect-[9/16] bg-gray-100 rounded flex items-center justify-center">
          <p className="text-gray-400">请先连接设备</p>
        </div>
      ) : isIOSPreview ? (
        <div className="relative">
          <img
            key={refreshKey}
            ref={imgRef}
            src={getIosPreviewUrl()!}
            alt="iOS Screen"
            className="max-w-full mx-auto rounded border"
            onMouseMove={handleMouseMove}
            onMouseLeave={() => setMousePos(null)}
          />
          {mousePos && (
            <div className="absolute top-2 right-2 bg-black/70 text-white px-2 py-1 rounded text-sm">
              ({mousePos.x}, {mousePos.y})
            </div>
          )}
        </div>
      ) : screenshot ? (
        <div className="relative">
          <img
            ref={imgRef}
            src={screenshot}
            alt="Device Screen"
            className="max-w-full mx-auto rounded border"
            onMouseMove={handleMouseMove}
            onMouseLeave={() => setMousePos(null)}
          />
          {mousePos && (
            <div className="absolute top-2 right-2 bg-black/70 text-white px-2 py-1 rounded text-sm">
              ({mousePos.x}, {mousePos.y})
            </div>
          )}
        </div>
      ) : (
        <div className="aspect-[9/16] bg-gray-100 rounded flex items-center justify-center">
          <p className="text-gray-400">加载中...</p>
        </div>
      )}

      <div className="flex gap-2 mt-4">
        <button
          onClick={() => connectedDevice && fetchScreenshot()}
          disabled={!connectedDevice && !isIOSPreview}
          className="flex-1 px-4 py-2 bg-gray-100 border rounded disabled:bg-gray-50"
        >
          📸 刷新截图
        </button>
      </div>

      {mousePos && (
        <div className="mt-2 text-sm text-gray-600">
          点击坐标: X={mousePos.x}, Y={mousePos.y}
        </div>
      )}
    </div>
  );
};
