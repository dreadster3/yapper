import { useState, useCallback, useEffect } from 'react';
import { Message, ChatSession, AppSettings } from '../types';
import { sendMessage, APIError } from '../utils/api';
import { getSessions, saveSessions, getActiveSessionId, saveActiveSessionId } from '../utils/storage';

export const useChat = (settings: AppSettings) => {
  const [sessions, setSessions] = useState<ChatSession[]>([]);
  const [activeSessionId, setActiveSessionId] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const activeSession = sessions.find(s => s.id === activeSessionId);

  useEffect(() => {
    const loadedSessions = getSessions();
    setSessions(loadedSessions);
    
    const savedActiveId = getActiveSessionId();
    if (savedActiveId && loadedSessions.find(s => s.id === savedActiveId)) {
      setActiveSessionId(savedActiveId);
    } else if (loadedSessions.length > 0) {
      setActiveSessionId(loadedSessions[0].id);
    }
  }, []);

  const createNewSession = useCallback(() => {
    const newSession: ChatSession = {
      id: Date.now().toString(),
      title: 'New Chat',
      messages: [],
      createdAt: Date.now(),
      updatedAt: Date.now(),
      model: settings.selectedModel,
      provider: settings.selectedProvider,
    };

    const updatedSessions = [newSession, ...sessions];
    setSessions(updatedSessions);
    setActiveSessionId(newSession.id);
    saveSessions(updatedSessions);
    saveActiveSessionId(newSession.id);
  }, [sessions, settings]);

  const deleteSession = useCallback((sessionId: string) => {
    const updatedSessions = sessions.filter(s => s.id !== sessionId);
    setSessions(updatedSessions);
    saveSessions(updatedSessions);

    if (sessionId === activeSessionId) {
      const newActiveId = updatedSessions.length > 0 ? updatedSessions[0].id : null;
      setActiveSessionId(newActiveId);
      if (newActiveId) {
        saveActiveSessionId(newActiveId);
      }
    }
  }, [sessions, activeSessionId]);

  const switchSession = useCallback((sessionId: string) => {
    setActiveSessionId(sessionId);
    saveActiveSessionId(sessionId);
  }, []);

  const sendMessageToAPI = useCallback(async (content: string) => {
    if (!activeSession) {
      createNewSession();
      return;
    }

    const provider = settings.providers.find(p => p.id === settings.selectedProvider);
    if (!provider) {
      setError('Provider not found');
      return;
    }

    if (provider.requiresApiKey && !provider.apiKey) {
      setError(`API key required for ${provider.name}`);
      return;
    }

    setIsLoading(true);
    setError(null);

    const userMessage: Message = {
      id: Date.now().toString(),
      role: 'user',
      content,
      timestamp: Date.now(),
    };

    const assistantMessage: Message = {
      id: (Date.now() + 1).toString(),
      role: 'assistant',
      content: '',
      timestamp: Date.now(),
      model: settings.selectedModel,
      provider: settings.selectedProvider,
    };

    // Update session with user message
    const updatedMessages = [...activeSession.messages, userMessage, assistantMessage];
    const updatedSession = {
      ...activeSession,
      messages: updatedMessages,
      updatedAt: Date.now(),
      title: activeSession.title === 'New Chat' ? content.slice(0, 50) + '...' : activeSession.title,
    };

    const updatedSessions = sessions.map(s => s.id === activeSession.id ? updatedSession : s);
    setSessions(updatedSessions);

    try {
      let fullResponse = '';
      
      const onChunk = (chunk: string) => {
        fullResponse += chunk;
        const newAssistantMessage = { ...assistantMessage, content: fullResponse };
        const newMessages = [...activeSession.messages, userMessage, newAssistantMessage];
        const newSession = { ...updatedSession, messages: newMessages, updatedAt: Date.now() };
        const newSessions = sessions.map(s => s.id === activeSession.id ? newSession : s);
        setSessions(newSessions);
      };

      const response = await sendMessage(
        content,
        activeSession.messages,
        provider,
        settings.selectedModel,
        settings.streamingEnabled ? onChunk : undefined
      );

      if (!settings.streamingEnabled) {
        fullResponse = response;
        onChunk(response);
      }

      // Save final state
      const finalSession = {
        ...updatedSession,
        messages: [...activeSession.messages, userMessage, { ...assistantMessage, content: fullResponse }],
        updatedAt: Date.now(),
      };
      const finalSessions = sessions.map(s => s.id === activeSession.id ? finalSession : s);
      setSessions(finalSessions);
      saveSessions(finalSessions);
    } catch (err) {
      if (err instanceof APIError) {
        setError(err.message);
      } else {
        setError('An unexpected error occurred');
      }
      
      // Remove the empty assistant message on error
      const errorSessions = sessions.map(s => 
        s.id === activeSession.id 
          ? { ...s, messages: [...activeSession.messages, userMessage] }
          : s
      );
      setSessions(errorSessions);
      saveSessions(errorSessions);
    } finally {
      setIsLoading(false);
    }
  }, [activeSession, sessions, settings, createNewSession]);

  return {
    sessions,
    activeSession,
    isLoading,
    error,
    createNewSession,
    deleteSession,
    switchSession,
    sendMessage: sendMessageToAPI,
    clearError: () => setError(null),
  };
};