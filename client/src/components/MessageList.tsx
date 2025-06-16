import React, { useEffect, useRef } from 'react';
import { User, Bot, Copy, Check } from 'lucide-react';
import { Message } from '../types';

interface MessageListProps {
  messages: Message[];
  isLoading: boolean;
}

export const MessageList: React.FC<MessageListProps> = ({ messages, isLoading }) => {
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const [copiedMessageId, setCopiedMessageId] = React.useState<string | null>(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const copyToClipboard = async (content: string, messageId: string) => {
    try {
      await navigator.clipboard.writeText(content);
      setCopiedMessageId(messageId);
      setTimeout(() => setCopiedMessageId(null), 2000);
    } catch (err) {
      console.error('Failed to copy text: ', err);
    }
  };

  const formatMessage = (content: string) => {
    // Simple markdown-like formatting
    return content
      .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
      .replace(/\*(.*?)\*/g, '<em>$1</em>')
      .replace(/`(.*?)`/g, '<code class="bg-gray-700 px-1 py-0.5 rounded text-sm">$1</code>')
      .replace(/\n/g, '<br />');
  };

  if (messages.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center p-8">
        <div className="text-center">
          <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center mx-auto mb-4">
            <Bot className="w-8 h-8 text-white" />
          </div>
          <h3 className="text-xl font-semibold text-gray-200 mb-2">Welcome to Yapper</h3>
          <p className="text-gray-400">Start a conversation with your AI assistant</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 overflow-y-auto p-4 space-y-4">
      {messages.map((message) => (
        <div
          key={message.id}
          className={`flex space-x-3 ${message.role === 'user' ? 'flex-row-reverse space-x-reverse' : ''}`}
        >
          <div className={`flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center ${
            message.role === 'user' 
              ? 'bg-gradient-to-br from-green-500 to-blue-500' 
              : 'bg-gradient-to-br from-purple-500 to-pink-500'
          }`}>
            {message.role === 'user' ? (
              <User className="w-4 h-4 text-white" />
            ) : (
              <Bot className="w-4 h-4 text-white" />
            )}
          </div>
          <div className={`flex-1 max-w-3xl ${message.role === 'user' ? 'text-right' : ''}`}>
            <div className={`inline-block p-4 rounded-2xl ${
              message.role === 'user'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-800 text-gray-100'
            }`}>
              <div 
                className="prose prose-invert max-w-none"
                dangerouslySetInnerHTML={{ __html: formatMessage(message.content) }}
              />
              {message.role === 'assistant' && (
                <div className="flex items-center justify-between mt-3 pt-3 border-t border-gray-700">
                  <div className="flex items-center space-x-2 text-xs text-gray-400">
                    {message.model && (
                      <span className="bg-gray-700 px-2 py-1 rounded">{message.model}</span>
                    )}
                    {message.provider && (
                      <span className="bg-gray-700 px-2 py-1 rounded">{message.provider}</span>
                    )}
                  </div>
                  <button
                    onClick={() => copyToClipboard(message.content, message.id)}
                    className="p-1 rounded hover:bg-gray-700 transition-colors"
                  >
                    {copiedMessageId === message.id ? (
                      <Check className="w-4 h-4 text-green-400" />
                    ) : (
                      <Copy className="w-4 h-4 text-gray-400" />
                    )}
                  </button>
                </div>
              )}
            </div>
            <div className="text-xs text-gray-500 mt-1">
              {new Date(message.timestamp).toLocaleTimeString()}
            </div>
          </div>
        </div>
      ))}
      {isLoading && (
        <div className="flex space-x-3">
          <div className="flex-shrink-0 w-8 h-8 rounded-full bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center">
            <Bot className="w-4 h-4 text-white" />
          </div>
          <div className="flex-1">
            <div className="inline-block p-4 bg-gray-800 rounded-2xl">
              <div className="flex space-x-1">
                <div className="w-2 h-2 bg-gray-400 rounded-full animate-pulse"></div>
                <div className="w-2 h-2 bg-gray-400 rounded-full animate-pulse" style={{ animationDelay: '0.2s' }}></div>
                <div className="w-2 h-2 bg-gray-400 rounded-full animate-pulse" style={{ animationDelay: '0.4s' }}></div>
              </div>
            </div>
          </div>
        </div>
      )}
      <div ref={messagesEndRef} />
    </div>
  );
};