export interface Device {
  id: string;
  platform: 'ios' | 'android';
  name: string;
  serial: string;
  status: 'available' | 'connected' | 'error';
  model?: string;
  address?: string;
  packageName?: string;
  skipSession?: boolean;
}

export interface CommandRecord {
  id: string;
  timestamp: string;
  command: string;
  params: Record<string, any>;
  success: boolean;
  output?: string;
  duration: string;
}

export interface CommandRequest {
  command: string;
  params: Record<string, any>;
}

export interface SerialPortInfo {
  name: string;
  description?: string;
}

export interface SerialConfig {
  name: string;
  baud: number;
  dataBits: number;
  parity: string;
  stopBits: number;
  timeout?: number;
}

export interface SerialConnection {
  id: string;
  config: SerialConfig;
  status: 'open' | 'closed' | 'error';
  commands?: SerialCommand[];
}

export interface SerialCommand {
  id: string;
  name: string;
  command: string;
  log?: string;
  timeout?: number;
}

export interface SerialCommandResult {
  id: string;
  success: boolean;
  sent: string;
  matched?: boolean;
  timestamp: string;
}