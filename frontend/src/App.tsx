import Canvas from "./Canvas"
import WSProvider from "./context/WSContext"

function App() {
    return (
        <>
            <WSProvider>
                <div className="h-screen w-screen">
                    <Canvas />
                </div>
            </WSProvider>
        </>
    )
}

export default App
