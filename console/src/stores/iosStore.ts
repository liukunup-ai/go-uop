import { create } from 'zustand';
import { api, IOSDevice } from '../api/client';

let iosPreviewUrl: string | null = null;

export const setIosPreviewUrl = (url: string | null) => {
  iosPreviewUrl = url;
};

export const getIosPreviewUrl = () => iosPreviewUrl;

interface IOSStore {
  devices: IOSDevice[];
  selectedDevice: IOSDevice | null;
  isLoading: boolean;
  error: string | null;
  forwardedPorts: number[];
  wdaCommand: string | null;
  forwardCommand: string | null;

  fetchDevices: () => Promise<void>;
  selectDevice: (device: IOSDevice | null) => void;
  setupForward: (udid: string, ports: number[]) => Promise<string>;
  removeForward: (udid: string) => Promise<void>;
  startWDA: (udid: string, bundleId?: string) => Promise<string>;
  stopWDA: (udid: string) => Promise<void>;
}

export const useIOSStore = create<IOSStore>((set, get) => ({
  devices: [],
  selectedDevice: null,
  isLoading: false,
  error: null,
  forwardedPorts: [],
  wdaCommand: null,
  forwardCommand: null,

  fetchDevices: async () => {
    try {
      const devices = await api.listIOSDevices();
      set({ devices, error: null });
    } catch (err: any) {
      set({ error: err.message });
    }
  },

  selectDevice: (device) => {
    set({ selectedDevice: device, forwardedPorts: [], wdaCommand: null, forwardCommand: null });
    setIosPreviewUrl(null);
  },

  setupForward: async (udid, ports) => {
    set({ isLoading: true, error: null });
    try {
      const result = await api.iosForward(udid, ports);
      const newForwardedPorts = [...new Set([...get().forwardedPorts, ...ports])];
      set({
        isLoading: false,
        forwardedPorts: newForwardedPorts,
        forwardCommand: result.command
      });
      if (ports.includes(9100)) {
        setIosPreviewUrl(`http://localhost:3333`);
      }
      return result.command;
    } catch (err: any) {
      set({ isLoading: false, error: err.message || 'Forward failed' });
      throw err;
    }
  },

  removeForward: async (udid) => {
    set({ isLoading: true, error: null });
    try {
      await api.iosForward(udid, []);
      set({ isLoading: false, forwardedPorts: [], forwardCommand: null });
      setIosPreviewUrl(null);
    } catch (err: any) {
      set({ isLoading: false, error: err.message });
    }
  },

  startWDA: async (udid, bundleId) => {
    set({ isLoading: true, error: null });
    try {
      const result = await api.iosWdaStart(udid, bundleId);
      set({
        isLoading: false,
        wdaCommand: result.command
      });
      return result.command;
    } catch (err: any) {
      set({ isLoading: false, error: err.message || 'WDA start failed' });
      throw err;
    }
  },

  stopWDA: async (udid) => {
    set({ isLoading: true, error: null });
    try {
      await api.iosWdaStop(udid);
      set({ isLoading: false, wdaCommand: null });
    } catch (err: any) {
      set({ isLoading: false, error: err.message });
    }
  },
}));
