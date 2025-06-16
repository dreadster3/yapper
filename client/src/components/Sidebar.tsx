import React from 'react';
import { Plus, MessageSquare, Trash2, Settings } from 'lucide-react';
import { ChatSession } from '../types';

interface SidebarProps {
  sessions: ChatSession[];
  activeSessionId: string | null;
  onNewSession: () => void;
  onSelectSession: (sessionId: string) => void;
  onDeleteSession: (sessionId: string) => void;
  onOpenSettings: () => void;
  isOpen: boolean;
  onToggle: () => void;
}

export const Sidebar: React.FC<SidebarProps> = ({
  sessions,
  activeSessionId,
  onNewSession,
  onSelectSession,
  onDeleteSession,
  onOpenSettings,
  isOpen,
}) => {
  if (!isOpen) return null;

  return (
    <div className="w-64 bg-gray-900/50 backdrop-blur-sm border-r border-gray-700 flex flex-col">
      <div className="p-4 border-b border-gray-700">
        <button
          onClick={onNewSession}
          className="w-full flex items-center space-x-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
        >
          <Plus className="w-4 h-4" />
          <span className="text-sm font-medium">New Chat</span>
        </button>
      </div>

      <div className="flex-1 overflow-y-auto">
        <div className="p-2 space-y-1">
          {sessions.map((session) => (
            <div
              key={session.id}
              className={`group flex items-center space-x-2 p-3 rounded-lg cursor-pointer transition-colors ${
                session.id === activeSessionId
                  ? 'bg-gray-700 text-white'
                  : 'text-gray-300 hover:bg-gray-800'
              }`}
              onClick={() => onSelectSession(session.id)}
            >
              <MessageSquare className="w-4 h-4 flex-shrink-0" />
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium truncate">{session.title}</p>
                <p className="text-xs text-gray-400">
                  {new Date(session.updatedAt).toLocaleDateString()}
                </p>
              </div>
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  onDeleteSession(session.id);
                }}
                className="opacity-0 group-hover:opacity-100 p-1 hover:bg-gray-600 rounded transition-opacity"
              >
                <Trash2 className="w-4 h-4 text-red-400" />
              </button>
            </div>
          ))}
        </div>
      </div>

      <div className="p-4 border-t border-gray-700">
        <button
          onClick={onOpenSettings}
          className="w-full flex items-center space-x-2 px-4 py-2 text-gray-300 hover:bg-gray-800 rounded-lg transition-colors"
        >
          <Settings className="w-4 h-4" />
          <span className="text-sm">Settings</span>
        </button>
      </div>
    </div>
  );
};