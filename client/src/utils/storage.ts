import { ChatSession, AppSettings, ModelProvider } from "../types";

const STORAGE_KEYS = {
  SETTINGS: "yapper_settings",
  SESSIONS: "yapper_sessions",
  ACTIVE_SESSION: "yapper_active_session",
};

export const defaultProviders: ModelProvider[] = [
  {
    id: "openai",
    name: "OpenAI",
    requiresApiKey: true,
    models: [
      { id: "gpt-4", name: "GPT-4", provider: "openai", contextLength: 8192 },
      {
        id: "gpt-4-turbo",
        name: "GPT-4 Turbo",
        provider: "openai",
        contextLength: 128000,
      },
      {
        id: "gpt-3.5-turbo",
        name: "GPT-3.5 Turbo",
        provider: "openai",
        contextLength: 4096,
      },
    ],
  },
  {
    id: "anthropic",
    name: "Anthropic",
    requiresApiKey: true,
    models: [
      {
        id: "claude-3-opus",
        name: "Claude 3 Opus",
        provider: "anthropic",
        contextLength: 200000,
      },
      {
        id: "claude-3-sonnet",
        name: "Claude 3 Sonnet",
        provider: "anthropic",
        contextLength: 200000,
      },
      {
        id: "claude-3-haiku",
        name: "Claude 3 Haiku",
        provider: "anthropic",
        contextLength: 200000,
      },
    ],
  },
  {
    id: "ollama",
    name: "Ollama",
    baseUrl: "http://localhost:11434",
    requiresApiKey: false,
    models: [
      {
        id: "deepseek-coder-v2",
        name: "DeepSeek Coder",
        provider: "ollama",
        contextLength: 16384,
      },
      {
        id: "deepseek-r1",
        name: "DeepSeek R1",
        provider: "ollama",
        contextLength: 16384,
      },
      {
        id: "llama2",
        name: "Llama 2",
        provider: "ollama",
        contextLength: 4096,
      },
      {
        id: "codellama",
        name: "Code Llama",
        provider: "ollama",
        contextLength: 16384,
      },
    ],
  },
];

export const defaultSettings: AppSettings = {
  theme: "dark",
  providers: defaultProviders,
  selectedModel: "gpt-3.5-turbo",
  selectedProvider: "openai",
  streamingEnabled: true,
};

export const getSettings = (): AppSettings => {
  const stored = localStorage.getItem(STORAGE_KEYS.SETTINGS);
  if (stored) {
    try {
      return { ...defaultSettings, ...JSON.parse(stored) };
    } catch (error) {
      console.error("Failed to parse settings:", error);
    }
  }
  return defaultSettings;
};

export const saveSettings = (settings: AppSettings): void => {
  localStorage.setItem(STORAGE_KEYS.SETTINGS, JSON.stringify(settings));
};

export const getSessions = (): ChatSession[] => {
  const stored = localStorage.getItem(STORAGE_KEYS.SESSIONS);
  if (stored) {
    try {
      return JSON.parse(stored);
    } catch (error) {
      console.error("Failed to parse sessions:", error);
    }
  }
  return [];
};

export const saveSessions = (sessions: ChatSession[]): void => {
  localStorage.setItem(STORAGE_KEYS.SESSIONS, JSON.stringify(sessions));
};

export const getActiveSessionId = (): string | null => {
  return localStorage.getItem(STORAGE_KEYS.ACTIVE_SESSION);
};

export const saveActiveSessionId = (sessionId: string): void => {
  localStorage.setItem(STORAGE_KEYS.ACTIVE_SESSION, sessionId);
};

