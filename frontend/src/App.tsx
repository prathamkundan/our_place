import Canvas from "./components/Canvas"
import LoginButton from "./components/LoginButton"
import WSProvider from "./context/WSContext"

function App() {
    return (
        <>
            <WSProvider>
                <div className="relative h-screen w-screen font-jetbrains">
                    <div className="absolute top-0 right-0 z-10 p-4">
                        <LoginButton />
                    </div>
                    <Canvas />
                </div>
            </WSProvider>
        </>
    )
}

export default App
