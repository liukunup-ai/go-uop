import axios from 'axios';
import type { Device, CommandRecord, CommandRequest, SerialPortInfo, SerialConfig, SerialConnection, SerialCommand, SerialCommandResult } from '../types';

const client = axios.create({
  baseURL: '/api',
  timeout: 30000,
});

export interface IOSDevice {
  udid: string;
  name?: string;
  model?: string;
  iosVersion?: string;
  status: 'available' | 'forwarding' | 'wda_running';
}

export const api = {
  listDevices: async (): Promise<Device[]> => {
    try {
      const res = await client.get<{ devices: Device[] }>('/devices');
      return res.data.devices || [];
    } catch {
      return [];
    }
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

  listSerialPorts: async (): Promise<SerialPortInfo[]> => {
    const res = await client.get<{ ports: SerialPortInfo[] }>('/serial/ports');
    return res.data.ports;
  },

  connectSerial: async (config: SerialConfig): Promise<{ connection: SerialConnection }> => {
    const res = await client.post<{ connection: SerialConnection }>('/serial/connect', { config });
    return res.data;
  },

  disconnectSerial: async (connId: string): Promise<void> => {
    await client.post(`/serial/${connId}/disconnect`);
  },

  sendSerialRaw: async (connId: string, data: string): Promise<SerialCommandResult> => {
    const res = await client.post<SerialCommandResult>(`/serial/${connId}/send`, { data });
    return res.data;
  },

  sendSerialByID: async (connId: string, commandId: string): Promise<SerialCommandResult> => {
    const res = await client.post<SerialCommandResult>(`/serial/${connId}/sendByID`, { commandId });
    return res.data;
  },

  listSerialCommands: async (connId: string): Promise<SerialCommand[]> => {
    const res = await client.get<{ commands: SerialCommand[] }>(`/serial/${connId}/commands`);
    return res.data.commands;
  },

  loadSerialCommandTable: async (connId: string, filePath?: string, yamlContent?: string): Promise<void> => {
    await client.post(`/serial/${connId}/loadTable`, { filePath, yamlContent });
  },

  listIOSDevices: async (): Promise<IOSDevice[]> => {
    try {
      const res = await client.get<{ devices: IOSDevice[] }>('/ios/devices');
      console.log('ios devices response:', res.data);
      return res.data.devices || [];
    } catch (e) {
      console.error('listIOSDevices failed:', e);
      return [];
    }
  },

  iosForward: async (udid: string, ports: number[]): Promise<{ success: boolean; command: string }> => {
    const res = await client.post<{ success: boolean; command: string }>('/ios/forward', { udid, ports });
    return res.data;
  },

  iosWdaStart: async (udid: string, bundleId?: string): Promise<{ success: boolean; command: string }> => {
    const res = await client.post<{ success: boolean; command: string }>('/ios/wda/start', { udid, bundleId });
    return res.data;
  },

  iosWdaStop: async (udid: string): Promise<void> => {
    await client.post('/ios/wda/stop', { udid });
  },
};