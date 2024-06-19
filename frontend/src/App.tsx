import Canvas from "./Canvas"
import WSProvider from "./context/WSContext"

function App() {
    return (
        <>
            <WSProvider>
                <div className="relative h-screen w-screen">
                    <div className="absolute top-0 right-0 z-10 p-4">
                        <button className="rounded-full ring-2 ring-black px-3 hover:bg-black hover:text-white">
                            Signin / Register
                        </button>
                    </div>
                    <Canvas />
                </div>
            </WSProvider>
        </>
    )
}

export default App
