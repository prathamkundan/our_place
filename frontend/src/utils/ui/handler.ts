import { View } from "./canvas";

function clamp(x: number, max: number, min: number) {
    return Math.max(min, Math.min(x, max));
}

export class MouseHandler {
    private isDragging: boolean = false;
    private dragStartX: number = 0;
    private dragStartY: number = 0;
    private noHighlight: boolean = false;

    public setPos: React.Dispatch<React.SetStateAction<number>> | null = null;
    public selectedPos: number[] | null = null;

    private view: View | null = null;
    private curr_coord: number[] = [0, 0]

    constructor() {
        console.log("Handler new")
    }

    handleMouseMove = (event: MouseEvent) => {
        const view = this.view!;
        if (!this.noHighlight) {
            const pos = this.view?.locToIndex(event.clientX, event.clientY)!

            this.view?.clearRect(this.curr_coord[0], this.curr_coord[1])
            this.view?.highLight(pos[0], pos[1])
            this.curr_coord = pos
        }

        if (this.isDragging) {
            // const time = Date.now()
            const deltaX = event.clientX - this.dragStartX;
            const deltaY = event.clientY - this.dragStartY;
            const k = view.canvas.width / view.vp_w;
            view.vp_ox = clamp(Math.ceil(view.vp_ox - deltaX / k), view.UNIVERSE_WIDTH - view.vp_w, 0);
            view.vp_oy = clamp(Math.ceil(view.vp_oy - deltaY / k), view.UNIVERSE_WIDTH - view.vp_h, 0);
            this.dragStartX = event.clientX;
            this.dragStartY = event.clientY;
            view.render();
        }
    }

    handleMouseDown = (event: MouseEvent) => {
        // this.timeDragStart = Date.now()
        this.isDragging = true;
        this.dragStartX = event.clientX;
        this.dragStartY = event.clientY;
    }

    handleWheel = (event: WheelEvent) => {
        if (this.view!.vp_w >= 3700 && event.deltaY > 0) return;
        const canvas = this.view?.canvas!;
        const view = this.view!;
        const wheelDelta = event.deltaY > 0 ? 1.1 : 0.9;
        const x = (event.clientX / canvas.width) * view.vp_w;
        const y = (event.clientY / canvas.height) * view.vp_h;
        if (view.vp_w * wheelDelta > this.view!.UNIVERSE_WIDTH || view.vp_w * wheelDelta < 10 * this.view!.BLOCK_WIDTH
            || view.vp_h * wheelDelta > view.UNIVERSE_WIDTH || view.vp_h * wheelDelta < 10 * view.BLOCK_WIDTH) {
            return
        }
        else {
            view.vp_w = Math.ceil(view.vp_w * wheelDelta)
            view.vp_h = Math.ceil(view.vp_h * wheelDelta)
            const ny = (event.clientY / canvas.height) * view.vp_h;
            const nx = (event.clientX / canvas.width) * view.vp_w;
            const dx = Math.ceil(x - nx);
            const dy = Math.ceil(y - ny);
            view.vp_ox = clamp(view.vp_ox + dx, view.UNIVERSE_WIDTH - view.vp_w, 0);
            view.vp_oy = clamp(view.vp_oy + dy, view.UNIVERSE_WIDTH - view.vp_h, 0);
        }
        view.render();
    }

    handleMouseUp = () => {
        this.isDragging = false;
    }

    handleMouseClick = (event: MouseEvent) => {
        const index = this.view?.locToIndex(event.clientX, event.clientY)!
        const linearIndex = this.view?.toIndex(index[0], index[1])!;

        if (this.selectedPos === null) {
            this.noHighlight = true;
            this.view?.highLight(index[0], index[1]);
            this.selectedPos = index;
            this.setPos!(linearIndex);
        } else {
            this.noHighlight = false;
            this.view?.clearRect(this.selectedPos[0], this.selectedPos[1])
            this.selectedPos = null;
            this.setPos!(-1);
        }
    }

    init(view: View, setPos: React.Dispatch<React.SetStateAction<number>>) {
        console.log("Handler init")
        this.view = view;
        this.setPos = setPos;
        this.view?.canvas.addEventListener("wheel", this.handleWheel);
        this.view?.canvas.addEventListener("mouseup", this.handleMouseUp);
        this.view?.canvas.addEventListener("mousedown", this.handleMouseDown);
        this.view?.canvas.addEventListener("mousemove", this.handleMouseMove);
        this.view?.canvas.addEventListener("click", this.handleMouseClick);
    }

    cleanup() {
        console.log("Handler cleanup")
        this.view?.canvas.removeEventListener("wheel", this.handleWheel)
        this.view?.canvas.removeEventListener("mouseup", this.handleMouseUp);
        this.view?.canvas.removeEventListener("mousedown", this.handleMouseDown);
        this.view?.canvas.removeEventListener("mousemove", this.handleMouseMove);
        this.view?.canvas.removeEventListener("click", this.handleMouseClick);
        this.setPos!(-1)
    }
}
