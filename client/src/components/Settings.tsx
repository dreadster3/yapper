import React, { useState } from 'react';
import { X, Eye, EyeOff, Save, AlertCircle } from 'lucide-react';
import { AppSettings, ModelProvider } from '../types';

interface SettingsProps {
  settings: AppSettings;
  onSave: (settings: AppSettings) => void;
  onClose: () => void;
  isOpen: boolean;
}

export const Settings: React.FC<SettingsProps> = ({ settings, onSave, onClose, isOpen }) => {
  const [localSettings, setLocalSettings] = useState<AppSettings>(settings);
  const [showApiKeys, setShowApiKeys] = useState<Record<string, boolean>>({});

  if (!isOpen) return null;

  const handleProviderChange = (providerId: string, updates: Partial<ModelProvider>) => {
    setLocalSettings(prev => ({
      ...prev,
      providers: prev.providers.map(p => 
        p.id === providerId ? { ...p, ...updates } : p
      )
    }));
  };

  const handleSave = () => {
    onSave(localSettings);
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-gray-800 rounded-lg w-full max-w-2xl max-h-[90vh] overflow-hidden">
        <div className="flex items-center justify-between p-6 border-b border-gray-700">
          <h2 className="text-xl font-semibold text-white">Settings</h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-700 rounded-lg transition-colors"
          >
            <X className="w-5 h-5 text-gray-400" />
          </button>
        </div>

        <div className="p-6 overflow-y-auto max-h-[calc(90vh-120px)]">
          <div className="space-y-6">
            <div>
              <h3 className="text-lg font-medium text-white mb-4">General</h3>
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <label className="text-sm font-medium text-gray-300">
                    Enable Streaming
                  </label>
                  <input
                    type="checkbox"
                    checked={localSettings.streamingEnabled}
                    onChange={(e) => setLocalSettings(prev => ({ 
                      ...prev, 
                      streamingEnabled: e.target.checked 
                    }))}
                    className="rounded border-gray-600 bg-gray-700 text-blue-600 focus:ring-blue-500"
                  />
                </div>
              </div>
            </div>

            <div>
              <h3 className="text-lg font-medium text-white mb-4">API Providers</h3>
              <div className="space-y-4">
                {localSettings.providers.map((provider) => (
                  <div key={provider.id} className="p-4 bg-gray-700/50 rounded-lg">
                    <div className="flex items-center justify-between mb-3">
                      <h4 className="font-medium text-white">{provider.name}</h4>
                      {provider.requiresApiKey && !provider.apiKey && (
                        <div className="flex items-center space-x-1 text-yellow-400">
                          <AlertCircle className="w-4 h-4" />
                          <span className="text-xs">API Key Required</span>
                        </div>
                      )}
                    </div>

                    {provider.requiresApiKey && (
                      <div className="mb-3">
                        <label className="block text-sm font-medium text-gray-300 mb-2">
                          API Key
                        </label>
                        <div className="relative">
                          <input
                            type={showApiKeys[provider.id] ? 'text' : 'password'}
                            value={provider.apiKey || ''}
                            onChange={(e) => handleProviderChange(provider.id, { apiKey: e.target.value })}
                            placeholder="Enter your API key"
                            className="w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                          />
                          <button
                            type="button"
                            onClick={() => setShowApiKeys(prev => ({ 
                              ...prev, 
                              [provider.id]: !prev[provider.id] 
                            }))}
                            className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-300"
                          >
                            {showApiKeys[provider.id] ? (
                              <EyeOff className="w-4 h-4" />
                            ) : (
                              <Eye className="w-4 h-4" />
                            )}
                          </button>
                        </div>
                      </div>
                    )}

                    {provider.baseUrl !== undefined && (
                      <div>
                        <label className="block text-sm font-medium text-gray-300 mb-2">
                          Base URL
                        </label>
                        <input
                          type="url"
                          value={provider.baseUrl}
                          onChange={(e) => handleProviderChange(provider.id, { baseUrl: e.target.value })}
                          placeholder="http://localhost:11434"
                          className="w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                      </div>
                    )}

                    <div className="mt-3 text-xs text-gray-400">
                      Models: {provider.models.map(m => m.name).join(', ')}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>

        <div className="p-6 border-t border-gray-700">
          <div className="flex items-center justify-end space-x-3">
            <button
              onClick={onClose}
              className="px-4 py-2 text-gray-300 hover:text-white transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleSave}
              className="flex items-center space-x-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
            >
              <Save className="w-4 h-4" />
              <span>Save</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};