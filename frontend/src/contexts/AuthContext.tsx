import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";

import api from "../services/api";
import type {
  ApiResponse,
  LoginRequest,
  LoginResponse,
  MeResponse,
  User,
} from "../types/auth";
import {
  clearAuthStorage,
  getStoredUser,
  getToken,
  setStoredUser,
  setToken,
} from "../utils/storage";

interface AuthContextValue {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({
  children,
}: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(
    () => getStoredUser<User>(),
  );

  const [token, setTokenState] = useState<string | null>(
    () => getToken(),
  );

  const [isLoading, setIsLoading] = useState(true);

  const logout = useCallback(() => {
    clearAuthStorage();

    setUser(null);
    setTokenState(null);
  }, []);

  const login = useCallback(
    async (credentials: LoginRequest) => {
      const response = await api.post<
        ApiResponse<LoginResponse>
      >("/auth/login", credentials);

      const loginData = response.data.data;

      setToken(loginData.token);
      setStoredUser(loginData.user);

      setTokenState(loginData.token);
      setUser(loginData.user);
    },
    [],
  );

  useEffect(() => {
    let mounted = true;

    async function validateSession() {
      const currentToken = getToken();

      if (!currentToken) {
        if (mounted) {
          setIsLoading(false);
        }

        return;
      }

      try {
        const response = await api.get<
          ApiResponse<MeResponse>
        >("/auth/me");

        if (!mounted) {
          return;
        }

        const me = response.data.data;

        const storedUser = getStoredUser<User>();

        const authenticatedUser: User = {
          id: me.id,
          username: me.username,
          role: me.role,
          full_name: storedUser?.full_name,
        };

        setUser(authenticatedUser);
        setTokenState(currentToken);
        setStoredUser(authenticatedUser);
      } catch {
        if (mounted) {
          logout();
        }
      } finally {
        if (mounted) {
          setIsLoading(false);
        }
      }
    }

    validateSession();

    return () => {
      mounted = false;
    };
  }, [logout]);

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      token,
      isAuthenticated: Boolean(token && user),
      isLoading,
      login,
      logout,
    }),
    [
      user,
      token,
      isLoading,
      login,
      logout,
    ],
  );

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error(
      "useAuth must be used inside AuthProvider",
    );
  }

  return context;
}