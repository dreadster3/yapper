import React, { useState } from 'react';
import { Zap, Settings, Eye, EyeOff, LogIn, AlertCircle } from 'lucide-react';
import { AuthConfig } from '../types';

interface LoginScreenProps {
  authConfig: AuthConfig;
  onUpdateConfig: (config: AuthConfig) => void;
  onLogin: () => void;
  isLoading: boolean;
  error: string | null;
  onClearError: () => void;
}

export const LoginScreen: React.FC<LoginScreenProps> = ({
  authConfig,
  onUpdateConfig,
  onLogin,
  isLoading,
  error,
  onClearError,
}) => {
  const [showSettings, setShowSettings] = useState(false);
  const [localConfig, setLocalConfig] = useState(authConfig);
  const [showClientId, setShowClientId] = useState(false);

  const handleSaveConfig = () => {
    onUpdateConfig(localConfig);
    setShowSettings(false);
  };

  const handleLogin = () => {
    onClearError();
    onLogin();
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        {/* Logo and Title */}
        <div className="text-center mb-8">
          <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-2xl">
            <Zap className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-white mb-2">Welcome to Yapper</h1>
          <p className="text-gray-400">Sign in to start chatting with AI models</p>
        </div>

        {/* Error Display */}
        {error && (
          <div className="mb-6 p-4 bg-red-500/20 border border-red-500/50 rounded-lg flex items-center space-x-3">
            <AlertCircle className="w-5 h-5 text-red-400 flex-shrink-0" />
            <div className="flex-1">
              <p className="text-red-400 text-sm">{error}</p>
            </div>
            <button
              onClick={onClearError}
              className="text-red-400 hover:text-red-300 transition-colors"
            >
              ×
            </button>
          </div>
        )}

        {/* Login Card */}
        <div className="bg-gray-800/50 backdrop-blur-sm border border-gray-700 rounded-2xl p-8 shadow-2xl">
          {!showSettings ? (
            <>
              <div className="space-y-6">
                <div className="text-center">
                  <h2 className="text-xl font-semibold text-white mb-2">Authentication Required</h2>
                  <p className="text-gray-400 text-sm">
                    Connect with your Authentik instance to continue
                  </p>
                </div>

                <div className="space-y-4">
                  <div className="p-4 bg-gray-700/30 rounded-lg border border-gray-600">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-300">Authentik URL</span>
                      <button
                        onClick={() => setShowSettings(true)}
                        className="text-blue-400 hover:text-blue-300 transition-colors"
                      >
                        <Settings className="w-4 h-4" />
                      </button>
                    </div>
                    <p className="text-gray-400 text-sm truncate">{authConfig.authentikUrl}</p>
                  </div>

                  <button
                    onClick={handleLogin}
                    disabled={isLoading}
                    className="w-full flex items-center justify-center space-x-2 px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 disabled:from-gray-600 disabled:to-gray-600 text-white font-medium rounded-lg transition-all duration-200 shadow-lg hover:shadow-xl disabled:cursor-not-allowed"
                  >
                    {isLoading ? (
                      <>
                        <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                        <span>Connecting...</span>
                      </>
                    ) : (
                      <>
                        <LogIn className="w-5 h-5" />
                        <span>Sign in with Authentik</span>
                      </>
                    )}
                  </button>
                </div>
              </div>
            </>
          ) : (
            <>
              <div className="space-y-6">
                <div className="flex items-center justify-between">
                  <h2 className="text-xl font-semibold text-white">Authentication Settings</h2>
                  <button
                    onClick={() => setShowSettings(false)}
                    className="text-gray-400 hover:text-white transition-colors"
                  >
                    ×
                  </button>
                </div>

                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-300 mb-2">
                      Authentik URL
                    </label>
                    <input
                      type="url"
                      value={localConfig.authentikUrl}
                      onChange={(e) => setLocalConfig(prev => ({ ...prev, authentikUrl: e.target.value }))}
                      placeholder="https://auth.example.com"
                      className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-300 mb-2">
                      Client ID
                    </label>
                    <div className="relative">
                      <input
                        type={showClientId ? 'text' : 'password'}
                        value={localConfig.clientId}
                        onChange={(e) => setLocalConfig(prev => ({ ...prev, clientId: e.target.value }))}
                        placeholder="yapper-client"
                        className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent pr-12"
                      />
                      <button
                        type="button"
                        onClick={() => setShowClientId(!showClientId)}
                        className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-300"
                      >
                        {showClientId ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
                      </button>
                    </div>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-300 mb-2">
                      Redirect URI
                    </label>
                    <input
                      type="url"
                      value={localConfig.redirectUri}
                      onChange={(e) => setLocalConfig(prev => ({ ...prev, redirectUri: e.target.value }))}
                      placeholder={window.location.origin + '/callback'}
                      className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-300 mb-2">
                      Scope
                    </label>
                    <input
                      type="text"
                      value={localConfig.scope}
                      onChange={(e) => setLocalConfig(prev => ({ ...prev, scope: e.target.value }))}
                      placeholder="openid profile email"
                      className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    />
                  </div>
                </div>

                <div className="flex space-x-3">
                  <button
                    onClick={() => setShowSettings(false)}
                    className="flex-1 px-4 py-3 text-gray-300 hover:text-white border border-gray-600 hover:border-gray-500 rounded-lg transition-colors"
                  >
                    Cancel
                  </button>
                  <button
                    onClick={handleSaveConfig}
                    className="flex-1 px-4 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                  >
                    Save
                  </button>
                </div>
              </div>
            </>
          )}
        </div>

        {/* Footer */}
        <div className="text-center mt-8">
          <p className="text-gray-500 text-sm">
            Secure authentication powered by Authentik
          </p>
        </div>
      </div>
    </div>
  );
};