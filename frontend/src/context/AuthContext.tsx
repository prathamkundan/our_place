import { createContext, useContext, useEffect, useState } from "react";
import { User, parseJwt } from "../utils/auth/auth";
import Cookies from "js-cookie";

export interface AuthData {
    user: User | null;
    login: () => void;
    logout: () => void;
}

const login = () => {
    const baseUrl = "https://accounts.google.com/o/oauth2/v2/auth?"
    const options = {
        client_id: import.meta.env.VITE_GOOGLE_CLIENT_ID as string,
        redirect_uri: import.meta.env.VITE_GOOGLE_OAUTH_REDIRECT as string,
        scope: ['https://www.googleapis.com/auth/userinfo.profile', 'https://www.googleapis.com/auth/userinfo.email'].join(" "),
        include_granted_scopes: 'true',
        response_type: 'code',
        prompt: 'consent'
    }

    const qs = new URLSearchParams(options);
    console.log(qs.toString())
    window.location.href = baseUrl + qs.toString();
}

const logout = () => {
    Cookies.remove("token")
    window.document.location.reload();
}

const AuthContext = createContext<AuthData>({user: null, login, logout})

export const useAuth = () => {
    return useContext(AuthContext);
}

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);

    useEffect(() => {
        const token = Cookies.get("token")
        if (token === undefined) {
            setUser(null)
        } else {
            const user = parseJwt(token)
            setUser(user)
        }
    }, [])

    return (
        <AuthContext.Provider value={{ user, login, logout }}>
            {children}
        </AuthContext.Provider>
    )
}
