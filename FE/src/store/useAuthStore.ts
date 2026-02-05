import { create } from 'zustand';
import { createJSONStorage, persist } from 'zustand/middleware';

interface JwtPayload {
    role: string;
    user_id: string;
    exp: number;
    iat?: number;
}

interface AuthState {
    token: string | null;
    role: string | null;
    userId: string | null;
    isAuthenticated: boolean;

    // Actions
    login: (token: string) => void;
    logout: () => void;

    // Utilities (Optional, if we want to expose raw check)
    isTokenExpired: () => boolean;
}

const parseJwt = (token: string): JwtPayload | null => {
    try {
        const base64Url = token.split('.')[1];
        if (!base64Url) return null;

        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload = decodeURIComponent(
            window.atob(base64)
                .split('')
                .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
                .join('')
        );
        return JSON.parse(jsonPayload) as JwtPayload;
    } catch (e) {
        console.error("JWT Parse Error", e);
        return null;
    }
};

export const useAuthStore = create<AuthState>()(
    persist(
        (set, get) => ({
            token: null,
            role: null,
            userId: null,
            isAuthenticated: false,

            login: (token: string) => {
                const decoded = parseJwt(token);
                if (decoded) {
                    set({
                        token,
                        role: decoded.role,
                        userId: decoded.user_id,
                        isAuthenticated: true,
                    });
                } else {
                    console.error("Attempted to login with invalid token");
                }
            },

            logout: () => {
                set({ token: null, role: null, userId: null, isAuthenticated: false });
                localStorage.removeItem('auth-storage'); // Clean cleanup if needed
            },

            isTokenExpired: () => {
                const { token } = get();
                if (!token) return true;
                const decoded = parseJwt(token);
                if (!decoded) return true;
                return decoded.exp * 1000 < Date.now();
            }
        }),
        {
            name: 'auth-storage', // Key in localStorage
            storage: createJSONStorage(() => localStorage),
            onRehydrateStorage: () => (state) => {
                // Optional: Check expiration on load
                if (state && state.token) {
                    const decoded = parseJwt(state.token);
                    if (!decoded || decoded.exp * 1000 < Date.now()) {
                        console.warn("Restored token is expired, logging out.");
                        state.logout();
                    }
                }
            },
        }
    )
);
