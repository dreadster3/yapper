export interface Message {
  id: string;
  role: 'user' | 'assistant' | 'system';
  content: string;
  timestamp: number;
  model?: string;
  provider?: string;
}

export interface ChatSession {
  id: string;
  title: string;
  messages: Message[];
  createdAt: number;
  updatedAt: number;
  model: string;
  provider: string;
}

export interface ModelProvider {
  id: string;
  name: string;
  baseUrl?: string;
  apiKey?: string;
  models: Model[];
  requiresApiKey: boolean;
}

export interface Model {
  id: string;
  name: string;
  description?: string;
  contextLength?: number;
  provider: string;
}

export interface AppSettings {
  theme: 'dark' | 'light';
  providers: ModelProvider[];
  selectedModel: string;
  selectedProvider: string;
  streamingEnabled: boolean;
}

export interface AuthConfig {
  authentikUrl: string;
  clientId: string;
  redirectUri: string;
  scope: string;
}

export interface User {
  id: string;
  username: string;
  email: string;
  name: string;
  avatar?: string;
}

export interface AuthState {
  isAuthenticated: boolean;
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
}