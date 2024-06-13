import './style.css'
import { View } from './ui/canvas'
import { MouseHandler } from './ui/handler'

document.querySelector<HTMLDivElement>('#app')!.innerHTML = `
  <div id = "canvas-container" class="h-screen w-screen">
  <canvas id="canvas"></canvas>
  </div>
`
const BLOCK_WIDTH = 10;


const container = document.getElementById("canvas-container") as HTMLElement
const view = new View('canvas', container.clientWidth, container.clientHeight, BLOCK_WIDTH, 512);
const handler = new MouseHandler(view)

window.addEventListener('resize', () => {
    let width = Math.round(container.clientWidth / 10) * BLOCK_WIDTH;
    let height = Math.round(container.clientHeight / 10) * BLOCK_WIDTH;
    view.setCanvasDimensions(width, height);
    view.drawGrid();
});

view.drawGrid()
handler.init()

