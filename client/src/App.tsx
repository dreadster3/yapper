import { useState, useEffect } from "react";
import { Header } from "./components/Header";
import { Sidebar } from "./components/Sidebar";
import { MessageList } from "./components/MessageList";
import { MessageInput } from "./components/MessageInput";
import { Settings } from "./components/Settings";
import { LoginScreen } from "./components/LoginScreen";
import { CallbackHandler } from "./components/CallbackHandler";
import { useChat } from "./hooks/useChat";
import { useAuth } from "./hooks/useAuth";
import { getSettings, saveSettings } from "./utils/storage";
import { AppSettings } from "./types";

function App() {
  const [settings, setSettings] = useState<AppSettings>(getSettings());
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [settingsOpen, setSettingsOpen] = useState(false);

  const {
    authState,
    authConfig,
    isLoading: authLoading,
    error: authError,
    login,
    logout,
    handleAuthCallback,
    updateAuthConfig,
    clearError: clearAuthError,
  } = useAuth();

  const {
    sessions,
    activeSession,
    isLoading,
    error,
    createNewSession,
    deleteSession,
    switchSession,
    sendMessage,
    clearError,
  } = useChat(settings);

  useEffect(() => {
    saveSettings(settings);
  }, [settings]);

  const handleSettingsChange = (newSettings: AppSettings) => {
    setSettings(newSettings);
  };

  // Check if we're on the callback route
  const isCallback =
    window.location.pathname === "/callback" ||
    window.location.search.includes("code=");

  // Handle OAuth callback
  if (isCallback) {
    return (
      <CallbackHandler
        onHandleCallback={handleAuthCallback}
        isLoading={authLoading}
        error={authError}
      />
    );
  }

  // Show login screen if not authenticated
  if (!authState.isAuthenticated) {
    return (
      <LoginScreen
        authConfig={authConfig}
        onUpdateConfig={updateAuthConfig}
        onLogin={login}
        isLoading={authLoading}
        error={authError}
        onClearError={clearAuthError}
      />
    );
  }

  const currentProvider = settings.providers.find(
    (p) => p.id === settings.selectedProvider,
  );
  const canSendMessage =
    !currentProvider?.requiresApiKey || !!currentProvider?.apiKey;

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex">
      <Sidebar
        sessions={sessions}
        activeSessionId={activeSession?.id || null}
        onNewSession={createNewSession}
        onSelectSession={switchSession}
        onDeleteSession={deleteSession}
        onOpenSettings={() => setSettingsOpen(true)}
        isOpen={sidebarOpen}
        onToggle={() => setSidebarOpen(!sidebarOpen)}
      />

      <div className="flex-1 flex flex-col">
        <Header
          settings={settings}
          user={authState.user!}
          onToggleSidebar={() => setSidebarOpen(!sidebarOpen)}
          onSettingsChange={handleSettingsChange}
          onLogout={logout}
          error={error}
          onClearError={clearError}
        />

        <MessageList
          messages={activeSession?.messages || []}
          isLoading={isLoading}
        />

        <MessageInput
          onSendMessage={sendMessage}
          isLoading={isLoading}
          disabled={!canSendMessage}
        />
      </div>

      <Settings
        settings={settings}
        onSave={handleSettingsChange}
        onClose={() => setSettingsOpen(false)}
        isOpen={settingsOpen}
      />
    </div>
  );
}

export default App;

