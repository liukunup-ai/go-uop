import axios from 'axios';
import type { Device, CommandRecord, CommandRequest } from '../types';

const client = axios.create({
  baseURL: '/api',
  timeout: 30000,
});

export const api = {
  listDevices: async (): Promise<Device[]> => {
    const res = await client.get<{ devices: Device[] }>('/devices');
    return res.data.devices;
  },

  connectDevice: async (device: Partial<Device>): Promise<{ device: Device }> => {
    const res = await client.post<{ device: Device }>('/devices/connect', device);
    return res.data;
  },

  getScreenshot: async (deviceId: string): Promise<Blob> => {
    const res = await client.get(`/devices/${deviceId}/screenshot`, {
      responseType: 'blob',
    });
    return res.data;
  },

  getDeviceInfo: async (deviceId: string): Promise<any> => {
    const res = await client.get(`/devices/${deviceId}/info`);
    return res.data;
  },

  executeCommand: async (
    deviceId: string,
    cmd: CommandRequest
  ): Promise<CommandRecord> => {
    const res = await client.post<CommandRecord>(
      `/devices/${deviceId}/commands`,
      cmd
    );
    return res.data;
  },

  getHistory: async (): Promise<CommandRecord[]> => {
    const res = await client.get<{ history: CommandRecord[] }>('/commands/history');
    return res.data.history;
  },

  exportYaml: async (ids?: string[], name?: string): Promise<string> => {
    const params = new URLSearchParams();
    if (ids && ids.length > 0) {
      ids.forEach((id) => params.append('ids', id));
    }
    if (name) {
      params.append('name', name);
    }
    const res = await client.get(`/export/yaml?${params.toString()}`, {
      responseType: 'text',
    });
    return res.data;
  },
};