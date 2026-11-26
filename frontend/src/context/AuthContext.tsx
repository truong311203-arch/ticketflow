import { createContext, useContext, useState, type ReactNode } from "react";
import type { AuthSession } from "../types";

const SESSION_KEY = "ticketflow_auth_session";

interface AuthContextValue {
  session: AuthSession | null;
  isAuthenticated: boolean;
  loginMock: (email: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

function loadStoredSession(): AuthSession | null {
  const raw = localStorage.getItem(SESSION_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as AuthSession;
  } catch {
    localStorage.removeItem(SESSION_KEY);
    return null;
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [session, setSession] = useState<AuthSession | null>(() => loadStoredSession());

  const loginMock = (email: string) => {
    // TODO (Phase Auth): Thay đoạn mock này bằng lời gọi API:
    // POST http://localhost:8080/api/v1/auth/login -> nhận về JWT token thật từ Go Backend
    const nextSession: AuthSession = {
      email: email.trim() || "demo@ticketflow.dev",
      role: "CUSTOMER",
      token: "mock-jwt-eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    };

    localStorage.setItem(SESSION_KEY, JSON.stringify(nextSession));
    setSession(nextSession);
  };

  const logout = () => {
    localStorage.removeItem(SESSION_KEY);
    setSession(null);
  };

  return (
    <AuthContext.Provider
      value={{
        session,
        isAuthenticated: session !== null,
        loginMock,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used inside AuthProvider");
  }
  return ctx;
}
