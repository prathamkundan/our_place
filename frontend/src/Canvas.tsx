import { useEffect, useRef } from "react"
import { View, colorMap } from "./utils/ui/canvas"
import { MouseHandler } from "./utils/ui/handler";

const Canvas = function() {
    const BLOCK_WIDTH = 10;
    const canvasRef = useRef<HTMLCanvasElement>(null);

    useEffect(() => {
        const canvas = canvasRef.current as HTMLCanvasElement
        const view = new View(canvas, BLOCK_WIDTH, 512);
        const mouseHandler = new MouseHandler(view)
        mouseHandler.init()
        view.render()

        return () => {
            mouseHandler.cleanup()
        }
    }, [])
    return (
        <>
            <div id="container" className="relative h-screen w-screen">
                <canvas ref={canvasRef} id="canvas"></canvas>
                <div className="absolute bottom-0 w-screen">
                    <div className="grid grid-rows-2 lg:grid-rows-1 grid-flow-col gap-2 bg-white w-screen md:w-fit mx-auto p-3 justify-center rounded-t-xl shadow-2xl">
                        {Object.entries(colorMap).map(([key, value]) =>
                            <button
                                key={key}
                                className="h-10 w-10 rounded-xl border-black border-2 hover:animate-bounce"
                                style={{ background: value }}
                            />
                        )}
                    </div>
                </div>
            </div>
        </>
    )
}

export default Canvas
