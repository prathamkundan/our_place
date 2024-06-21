import Canvas from "./Canvas"
import WSProvider from "./context/WSContext"

function App() {
    const googleOAuth = () => {
        const baseUrl = "https://accounts.google.com/o/oauth2/v2/auth"
        const options = {
            client_id: import.meta.env.VITE_GOOGLE_CLIENT_ID,
            redirect_uri: import.meta.env.VITE_GOOGLE_OAUTH_REDIRECT,
            scope: ['https://www.googleapis.com/auth/userinfo.profile', 'https://www.googleapis.com/auth/userinfo.email'].join(" "),
            include_granted_scopes: 'true',
            response_type: 'code',
            prompt: 'consent'
        }

        const qs = new URLSearchParams(options);
        return `${baseUrl}?${qs.toString()}`
    }
    return (
        <>
            <WSProvider>
                <div className="relative h-screen w-screen">
                    <div className="absolute top-0 right-0 z-10 p-4">
                        <a href={googleOAuth()} className="rounded-full ring-2 ring-black px-3 hover:bg-black hover:text-white">
                            Signin / Register
                        </a>
                    </div>
                    <Canvas />
                </div>
            </WSProvider>
        </>
    )
}

export default App
