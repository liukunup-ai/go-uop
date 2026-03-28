import { create } from 'zustand';
import type { Device, CommandRecord } from '../types';
import { api } from '../api/client';

interface DeviceStore {
  devices: Device[];
  selectedDevice: Device | null;
  connectedDevice: Device | null;
  history: CommandRecord[];
  screenshot: string | null;
  isLoading: boolean;
  error: string | null;

  fetchDevices: () => Promise<void>;
  selectDevice: (device: Device | null) => void;
  connectDevice: (device: Device) => Promise<void>;
  disconnectDevice: () => Promise<void>;
  fetchScreenshot: () => Promise<void>;
  executeCommand: (command: string, params: Record<string, any>) => Promise<void>;
  fetchHistory: () => Promise<void>;
  clearHistory: () => void;
  exportYaml: (ids?: string[]) => Promise<void>;
}

export const useDeviceStore = create<DeviceStore>((set, get) => ({
  devices: [],
  selectedDevice: null,
  connectedDevice: null,
  history: [],
  screenshot: null,
  isLoading: false,
  error: null,

  fetchDevices: async () => {
    try {
      const devices = await api.listDevices();
      set({ devices, error: null });
    } catch (err: any) {
      set({ error: err.message });
    }
  },

  selectDevice: (device) => {
    set({ selectedDevice: device });
  },

  connectDevice: async (device) => {
    try {
      set({ isLoading: true, error: null });
      await api.connectDevice(device);
      set({ connectedDevice: device, isLoading: false });
      get().fetchScreenshot();
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
    }
  },

  disconnectDevice: async () => {
    set({ connectedDevice: null, screenshot: null });
  },

  fetchScreenshot: async () => {
    const { connectedDevice } = get();
    if (!connectedDevice) return;

    try {
      const blob = await api.getScreenshot(connectedDevice.id);
      const url = URL.createObjectURL(blob);
      set({ screenshot: url });
    } catch (err) {
      console.error('Screenshot failed:', err);
    }
  },

  executeCommand: async (command, params) => {
    const { connectedDevice } = get();
    if (!connectedDevice) {
      set({ error: 'No device connected' });
      return;
    }

    try {
      set({ isLoading: true, error: null });
      const record = await api.executeCommand(connectedDevice.id, { command, params });
      set((state) => ({
        history: [...state.history, record],
        isLoading: false,
      }));
      if (command === 'screenshot' || command === 'tap' || command === 'swipe') {
        get().fetchScreenshot();
      }
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
    }
  },

  fetchHistory: async () => {
    try {
      const history = await api.getHistory();
      set({ history });
    } catch (err) {
      console.error('Fetch history failed:', err);
    }
  },

  clearHistory: () => {
    set({ history: [] });
  },

  exportYaml: async (ids) => {
    try {
      const yaml = await api.exportYaml(ids);
      const blob = new Blob([yaml], { type: 'text/yaml' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'debug-session.yaml';
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      console.error('Export YAML failed:', err);
    }
  },
}));