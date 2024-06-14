import { useEffect, useRef, useState } from "react"
import { View, colorMap } from "./utils/ui/canvas"
import { MouseHandler } from "./utils/ui/handler";
import { WebSocketController } from "./utils/network/socket_control";

const Canvas = function() {
    const BLOCK_WIDTH = 10;
    const canvasRef = useRef<HTMLCanvasElement>(null);
    const [selectedPos, setSelectedPos] = useState<number>(-1);
    const [seletedColor, setSelectedColor] = useState<number>(-1);
    const [wsController, _] = useState(new WebSocketController());

    useEffect(() => {
        const canvas = canvasRef.current as HTMLCanvasElement
        const view = new View(canvas, BLOCK_WIDTH, 512);
        const mouseHandler = new MouseHandler()

        mouseHandler.init(view, setSelectedPos)
        wsController.init("ws://localhost:8000/ws", view)
        view.render()

        return () => {
            mouseHandler.cleanup()
            wsController.cleanup()
        }
    }, [])

    const onSubmit = () => {
        wsController.sendUpdate(selectedPos, seletedColor);
    }

    return (
        <>
            <div id="container" className="relative h-screen w-screen">
                <canvas ref={canvasRef} id="canvas"></canvas>
                <div className="absolute bottom-0 w-screen">
                    {selectedPos > -1 && <div className="grid grid-flow-col grid-rows-2 lg:grid-rows-1 gap-2 bg-white w-fit mx-auto p-3 justify-center rounded-t-xl shadow-2xl">
                        {Object.entries(colorMap).map(([key, value], index) =>
                            <button
                                key={key}
                                className={"h-10 w-10 rounded-xl border-black border-2 hover:animate-bounce " + (seletedColor === index ? "ring-4 ring-black ring-offset-2" : "")}
                                style={{ background: value }}
                                onClick={() => setSelectedColor(index)}
                            />
                        )}
                        {seletedColor < 16 && seletedColor >= 0 && selectedPos > 0 &&
                            <button onClick={onSubmit} className="rounded-full mx-3 border-black border-2 px-2"> OK </button>}
                    </div>
                    }
                </div>
            </div>
        </>
    )
}

export default Canvas
