import { create } from 'zustand';
import type { SerialPortInfo, SerialConfig, SerialConnection, SerialCommand } from '../types';
import { api } from '../api/client';

interface SerialStore {
  ports: SerialPortInfo[];
  connections: SerialConnection[];
  activeConnection: SerialConnection | null;
  commands: SerialCommand[];
  logs: SerialLogEntry[];
  isLoading: boolean;
  error: string | null;
  config: SerialConfig;

  fetchPorts: () => Promise<void>;
  connect: (config: SerialConfig) => Promise<void>;
  disconnect: (connId: string) => Promise<void>;
  sendRaw: (data: string) => Promise<void>;
  sendByID: (commandId: string) => Promise<void>;
  loadCommandTable: (filePath?: string, yamlContent?: string) => Promise<void>;
  fetchCommands: () => Promise<void>;
  setConfig: (config: Partial<SerialConfig>) => void;
  clearLogs: () => void;
}

export interface SerialLogEntry {
  id: string;
  timestamp: string;
  direction: 'in' | 'out';
  data: string;
  success?: boolean;
}

export const useSerialStore = create<SerialStore>((set, get) => ({
  ports: [],
  connections: [],
  activeConnection: null,
  commands: [],
  logs: [],
  isLoading: false,
  error: null,
  config: {
    name: '',
    baud: 115200,
    dataBits: 8,
    parity: 'N',
    stopBits: 1,
  },

  fetchPorts: async () => {
    try {
      set({ isLoading: true, error: null });
      const ports = await api.listSerialPorts();
      set({ ports, isLoading: false });
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
    }
  },

  connect: async (config) => {
    try {
      set({ isLoading: true, error: null });
      const { connection } = await api.connectSerial(config);
      set((state) => ({
        connections: [...state.connections, connection],
        activeConnection: connection,
        isLoading: false,
      }));
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
    }
  },

  disconnect: async (connId) => {
    try {
      set({ isLoading: true, error: null });
      await api.disconnectSerial(connId);
      set((state) => ({
        connections: state.connections.filter((c) => c.id !== connId),
        activeConnection: state.activeConnection?.id === connId ? null : state.activeConnection,
        isLoading: false,
      }));
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
    }
  },

  sendRaw: async (data) => {
    const { activeConnection } = get();
    if (!activeConnection) {
      set({ error: 'No connection selected' });
      return;
    }

    const logEntry: SerialLogEntry = {
      id: `log-${Date.now()}`,
      timestamp: new Date().toISOString(),
      direction: 'out',
      data,
    };

    set((state) => ({ logs: [...state.logs, logEntry] }));

    try {
      const result = await api.sendSerialRaw(activeConnection.id, data);
      const responseLog: SerialLogEntry = {
        id: `log-${Date.now()}-resp`,
        timestamp: result.timestamp,
        direction: 'in',
        data: result.sent,
        success: result.success,
      };
      set((state) => ({ logs: [...state.logs, responseLog] }));
    } catch (err: any) {
      const errorLog: SerialLogEntry = {
        id: `log-${Date.now()}-err`,
        timestamp: new Date().toISOString(),
        direction: 'in',
        data: err.message,
        success: false,
      };
      set((state) => ({ logs: [...state.logs, errorLog] }));
    }
  },

  sendByID: async (commandId) => {
    const { activeConnection } = get();
    if (!activeConnection) {
      set({ error: 'No connection selected' });
      return;
    }

    const logEntry: SerialLogEntry = {
      id: `log-${Date.now()}`,
      timestamp: new Date().toISOString(),
      direction: 'out',
      data: `CMD: ${commandId}`,
    };

    set((state) => ({ logs: [...state.logs, logEntry] }));

    try {
      const result = await api.sendSerialByID(activeConnection.id, commandId);
      const responseLog: SerialLogEntry = {
        id: `log-${Date.now()}-resp`,
        timestamp: result.timestamp,
        direction: 'in',
        data: result.sent,
        success: result.success,
      };
      set((state) => ({ logs: [...state.logs, responseLog] }));
    } catch (err: any) {
      const errorLog: SerialLogEntry = {
        id: `log-${Date.now()}-err`,
        timestamp: new Date().toISOString(),
        direction: 'in',
        data: err.message,
        success: false,
      };
      set((state) => ({ logs: [...state.logs, errorLog] }));
    }
  },

  loadCommandTable: async (filePath, yamlContent) => {
    const { activeConnection } = get();
    if (!activeConnection) {
      set({ error: 'No connection selected' });
      return;
    }

    try {
      set({ isLoading: true, error: null });
      await api.loadSerialCommandTable(activeConnection.id, filePath, yamlContent);
      await get().fetchCommands();
      set({ isLoading: false });
    } catch (err: any) {
      set({ error: err.message, isLoading: false });
    }
  },

  fetchCommands: async () => {
    const { activeConnection } = get();
    if (!activeConnection) return;

    try {
      const commands = await api.listSerialCommands(activeConnection.id);
      set({ commands });
    } catch (err: any) {
      console.error('Fetch commands failed:', err);
    }
  },

  setConfig: (config) => {
    set((state) => ({
      config: { ...state.config, ...config },
    }));
  },

  clearLogs: () => {
    set({ logs: [] });
  },
}));
