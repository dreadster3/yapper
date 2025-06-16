import React from "react";
import { Menu, Zap, AlertCircle, LogOut, User } from "lucide-react";
import { AppSettings, User as UserType } from "../types";
import { ModelSelector } from "./ModelSelector";

interface HeaderProps {
  settings: AppSettings;
  user: UserType;
  onToggleSidebar: () => void;
  onSettingsChange: (settings: AppSettings) => void;
  onLogout: () => void;
  error: string | null;
  onClearError: () => void;
}

export const Header: React.FC<HeaderProps> = ({
  settings,
  user,
  onToggleSidebar,
  onSettingsChange,
  onLogout,
  error,
  onClearError,
}) => {
  const [showUserMenu, setShowUserMenu] = React.useState(false);

  const handleModelChange = (provider: string, model: string) => {
    onSettingsChange({
      ...settings,
      selectedProvider: provider,
      selectedModel: model,
    });
  };

  return (
    <div className="flex items-center justify-between p-4 bg-gray-900/50 backdrop-blur-sm border-b border-gray-700">
      <div className="flex items-center space-x-4">
        <button
          onClick={onToggleSidebar}
          className="p-2 hover:bg-gray-700 rounded-lg transition-colors"
        >
          <Menu className="w-5 h-5 text-gray-400" />
        </button>
        <div className="flex items-center space-x-2">
          <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
            <Zap className="w-5 h-5 text-white" />
          </div>
          <h1 className="text-xl font-bold text-white">Yapper</h1>
        </div>
      </div>

      <div className="flex items-center space-x-4">
        {error && (
          <div className="flex items-center space-x-2 px-3 py-2 bg-red-500/20 border border-red-500/50 rounded-lg">
            <AlertCircle className="w-4 h-4 text-red-400" />
            <span className="text-sm text-red-400">{error}</span>
            <button
              onClick={onClearError}
              className="text-red-400 hover:text-red-300 ml-2"
            >
              ×
            </button>
          </div>
        )}

        <ModelSelector settings={settings} onModelChange={handleModelChange} />

        {/* User Menu */}
        <div className="relative">
          <button
            onClick={() => setShowUserMenu(!showUserMenu)}
            className="flex items-center space-x-2 p-2 hover:bg-gray-700 rounded-lg transition-colors"
          >
            {user.avatar ? (
              <img
                src={user.avatar}
                alt={user.name}
                className="w-8 h-8 rounded-full"
              />
            ) : (
              <div className="w-8 h-8 bg-gradient-to-br from-green-500 to-blue-500 rounded-full flex items-center justify-center">
                <User className="w-4 h-4 text-white" />
              </div>
            )}
            <span className="text-sm text-gray-300 hidden sm:block">
              {user.name}
            </span>
          </button>

          {showUserMenu && (
            <div className="absolute top-full right-0 mt-2 w-64 bg-gray-800/95 backdrop-blur-sm border border-gray-700 rounded-lg shadow-xl z-50">
              <div className="p-4 border-b border-gray-700">
                <div className="flex items-center space-x-3">
                  {user.avatar ? (
                    <img
                      src={user.avatar}
                      alt={user.name}
                      className="w-10 h-10 rounded-full"
                    />
                  ) : (
                    <div className="w-10 h-10 bg-gradient-to-br from-green-500 to-blue-500 rounded-full flex items-center justify-center">
                      <User className="w-5 h-5 text-white" />
                    </div>
                  )}
                  <div>
                    <p className="font-medium text-white">{user.name}</p>
                    <p className="text-sm text-gray-400">{user.email}</p>
                  </div>
                </div>
              </div>
              <div className="p-2">
                <button
                  onClick={() => {
                    setShowUserMenu(false);
                    onLogout();
                  }}
                  className="w-full flex items-center space-x-2 px-3 py-2 text-red-400 hover:bg-red-500/20 rounded-lg transition-colors"
                >
                  <LogOut className="w-4 h-4" />
                  <span>Sign Out</span>
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

