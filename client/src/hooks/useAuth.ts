import { useState, useEffect, useCallback } from 'react';
import { AuthState, AuthConfig } from '../types';
import {
  getAuthState,
  saveAuthState,
  getAuthConfig,
  saveAuthConfig,
  initiateLogin,
  handleCallback,
  logout,
  refreshAccessToken,
} from '../utils/auth';

export const useAuth = () => {
  const [authState, setAuthState] = useState<AuthState>(getAuthState());
  const [authConfig, setAuthConfig] = useState<AuthConfig>(getAuthConfig());
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    saveAuthState(authState);
  }, [authState]);

  useEffect(() => {
    saveAuthConfig(authConfig);
  }, [authConfig]);

  const login = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      await initiateLogin(authConfig);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      setIsLoading(false);
    }
  }, [authConfig]);

  const handleAuthCallback = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const newAuthState = await handleCallback(authConfig);
      setAuthState(newAuthState);
      // Redirect to main app
      window.history.replaceState({}, document.title, '/');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Authentication failed');
    } finally {
      setIsLoading(false);
    }
  }, [authConfig]);

  const logoutUser = useCallback(async () => {
    if (authState.accessToken) {
      await logout(authConfig, authState.accessToken);
    }
    setAuthState({
      isAuthenticated: false,
      user: null,
      accessToken: null,
      refreshToken: null,
    });
  }, [authConfig, authState.accessToken]);

  const refreshToken = useCallback(async () => {
    if (!authState.refreshToken) return;
    
    setIsLoading(true);
    try {
      const newAuthState = await refreshAccessToken(authConfig, authState.refreshToken);
      setAuthState(newAuthState);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Token refresh failed');
      // If refresh fails, logout user
      setAuthState({
        isAuthenticated: false,
        user: null,
        accessToken: null,
        refreshToken: null,
      });
    } finally {
      setIsLoading(false);
    }
  }, [authConfig, authState.refreshToken]);

  const updateAuthConfig = useCallback((config: AuthConfig) => {
    setAuthConfig(config);
  }, []);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  return {
    authState,
    authConfig,
    isLoading,
    error,
    login,
    logout: logoutUser,
    handleAuthCallback,
    refreshToken,
    updateAuthConfig,
    clearError,
  };
};