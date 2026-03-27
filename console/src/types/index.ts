export interface Device {
  id: string;
  platform: 'ios' | 'android';
  name: string;
  serial: string;
  status: 'available' | 'connected' | 'error';
  model?: string;
  address?: string;
  packageName?: string;
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