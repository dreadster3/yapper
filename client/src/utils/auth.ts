import { AuthConfig, AuthState, User } from '../types';

const STORAGE_KEYS = {
  AUTH_STATE: 'yapper_auth_state',
  AUTH_CONFIG: 'yapper_auth_config',
};

export const defaultAuthConfig: AuthConfig = {
  authentikUrl: 'https://auth.example.com',
  clientId: 'yapper-client',
  redirectUri: window.location.origin + '/callback',
  scope: 'openid profile email',
};

export const getAuthConfig = (): AuthConfig => {
  const stored = localStorage.getItem(STORAGE_KEYS.AUTH_CONFIG);
  if (stored) {
    try {
      return { ...defaultAuthConfig, ...JSON.parse(stored) };
    } catch (error) {
      console.error('Failed to parse auth config:', error);
    }
  }
  return defaultAuthConfig;
};

export const saveAuthConfig = (config: AuthConfig): void => {
  localStorage.setItem(STORAGE_KEYS.AUTH_CONFIG, JSON.stringify(config));
};

export const getAuthState = (): AuthState => {
  const stored = localStorage.getItem(STORAGE_KEYS.AUTH_STATE);
  if (stored) {
    try {
      const state = JSON.parse(stored);
      // Check if token is expired (basic check)
      if (state.accessToken && isTokenExpired(state.accessToken)) {
        return {
          isAuthenticated: false,
          user: null,
          accessToken: null,
          refreshToken: null,
        };
      }
      return state;
    } catch (error) {
      console.error('Failed to parse auth state:', error);
    }
  }
  return {
    isAuthenticated: false,
    user: null,
    accessToken: null,
    refreshToken: null,
  };
};

export const saveAuthState = (state: AuthState): void => {
  localStorage.setItem(STORAGE_KEYS.AUTH_STATE, JSON.stringify(state));
};

export const clearAuthState = (): void => {
  localStorage.removeItem(STORAGE_KEYS.AUTH_STATE);
};

export const generateCodeVerifier = (): string => {
  const array = new Uint8Array(32);
  crypto.getRandomValues(array);
  return btoa(String.fromCharCode.apply(null, Array.from(array)))
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=/g, '');
};

export const generateCodeChallenge = async (verifier: string): Promise<string> => {
  const encoder = new TextEncoder();
  const data = encoder.encode(verifier);
  const digest = await crypto.subtle.digest('SHA-256', data);
  return btoa(String.fromCharCode.apply(null, Array.from(new Uint8Array(digest))))
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=/g, '');
};

export const initiateLogin = async (config: AuthConfig): Promise<void> => {
  const codeVerifier = generateCodeVerifier();
  const codeChallenge = await generateCodeChallenge(codeVerifier);
  
  // Store code verifier for later use
  sessionStorage.setItem('code_verifier', codeVerifier);
  
  const params = new URLSearchParams({
    response_type: 'code',
    client_id: config.clientId,
    redirect_uri: config.redirectUri,
    scope: config.scope,
    code_challenge: codeChallenge,
    code_challenge_method: 'S256',
    state: generateCodeVerifier(), // Use as random state
  });
  
  const authUrl = `${config.authentikUrl}/application/o/authorize/?${params.toString()}`;
  window.location.href = authUrl;
};

export const handleCallback = async (config: AuthConfig): Promise<AuthState> => {
  const urlParams = new URLSearchParams(window.location.search);
  const code = urlParams.get('code');
  const error = urlParams.get('error');
  
  if (error) {
    throw new Error(`OAuth error: ${error}`);
  }
  
  if (!code) {
    throw new Error('No authorization code received');
  }
  
  const codeVerifier = sessionStorage.getItem('code_verifier');
  if (!codeVerifier) {
    throw new Error('No code verifier found');
  }
  
  // Exchange code for tokens
  const tokenResponse = await fetch(`${config.authentikUrl}/application/o/token/`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    body: new URLSearchParams({
      grant_type: 'authorization_code',
      client_id: config.clientId,
      code,
      redirect_uri: config.redirectUri,
      code_verifier: codeVerifier,
    }),
  });
  
  if (!tokenResponse.ok) {
    throw new Error('Failed to exchange code for tokens');
  }
  
  const tokens = await tokenResponse.json();
  
  // Get user info
  const userResponse = await fetch(`${config.authentikUrl}/application/o/userinfo/`, {
    headers: {
      'Authorization': `Bearer ${tokens.access_token}`,
    },
  });
  
  if (!userResponse.ok) {
    throw new Error('Failed to get user info');
  }
  
  const userInfo = await userResponse.json();
  
  const user: User = {
    id: userInfo.sub,
    username: userInfo.preferred_username || userInfo.email,
    email: userInfo.email,
    name: userInfo.name || userInfo.preferred_username,
    avatar: userInfo.picture,
  };
  
  const authState: AuthState = {
    isAuthenticated: true,
    user,
    accessToken: tokens.access_token,
    refreshToken: tokens.refresh_token,
  };
  
  // Clean up
  sessionStorage.removeItem('code_verifier');
  
  return authState;
};

export const logout = async (config: AuthConfig, accessToken: string): Promise<void> => {
  try {
    await fetch(`${config.authentikUrl}/application/o/revoke/`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
        'Authorization': `Bearer ${accessToken}`,
      },
      body: new URLSearchParams({
        token: accessToken,
        client_id: config.clientId,
      }),
    });
  } catch (error) {
    console.error('Failed to revoke token:', error);
  }
  
  clearAuthState();
  window.location.href = `${config.authentikUrl}/application/o/end-session/?post_logout_redirect_uri=${encodeURIComponent(window.location.origin)}`;
};

const isTokenExpired = (token: string): boolean => {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.exp * 1000 < Date.now();
  } catch {
    return true;
  }
};

export const refreshAccessToken = async (config: AuthConfig, refreshToken: string): Promise<AuthState> => {
  const response = await fetch(`${config.authentikUrl}/application/o/token/`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    body: new URLSearchParams({
      grant_type: 'refresh_token',
      client_id: config.clientId,
      refresh_token: refreshToken,
    }),
  });
  
  if (!response.ok) {
    throw new Error('Failed to refresh token');
  }
  
  const tokens = await response.json();
  
  // Get updated user info
  const userResponse = await fetch(`${config.authentikUrl}/application/o/userinfo/`, {
    headers: {
      'Authorization': `Bearer ${tokens.access_token}`,
    },
  });
  
  const userInfo = await userResponse.json();
  
  const user: User = {
    id: userInfo.sub,
    username: userInfo.preferred_username || userInfo.email,
    email: userInfo.email,
    name: userInfo.name || userInfo.preferred_username,
    avatar: userInfo.picture,
  };
  
  return {
    isAuthenticated: true,
    user,
    accessToken: tokens.access_token,
    refreshToken: tokens.refresh_token || refreshToken,
  };
};