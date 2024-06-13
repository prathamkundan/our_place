import { View } from "./canvas";

function clamp(x: number, max: number, min: number) {
    return Math.max(min, Math.min(x, max));
}

export class MouseHandler {
    private isDragging: Boolean = false;
    private dragStartX: number = 0;
    private dragStartY: number = 0;

    private view: View;
    private curr_coord: number[] = [0, 0]

    constructor(view: View) {
        this.view = view;
    }

    handleMouseMove = (event: MouseEvent) => {
        const view = this.view;
        let pos = this.view.locToIndex(event.clientX, event.clientY)

        if (this.curr_coord !== pos) {
            this.view.clearRect(this.curr_coord[0], this.curr_coord[1])
            this.view.highLight(pos[0], pos[1])
            this.curr_coord = pos
        }

        if (this.isDragging) {
            const deltaX = event.clientX - this.dragStartX;
            const deltaY = event.clientY - this.dragStartY;
            const k = view.canvas.width / view.vp_w;
            view.vp_ox = clamp(view.vp_ox - deltaX / k, view.UNIVERSE_WIDTH - view.vp_w, 0);
            view.vp_oy = clamp(view.vp_oy - deltaY / k, view.UNIVERSE_WIDTH - view.vp_h, 0);
            this.dragStartX = event.clientX;
            this.dragStartY = event.clientY;
            view.render();
        }
    }

    handleMouseDown = (event: MouseEvent) => {
        this.isDragging = true;
        this.dragStartX = event.clientX;
        this.dragStartY = event.clientY;
    }

    handleWheel = (event: WheelEvent) => {
        const canvas = this.view.canvas;
        const view = this.view;
        const wheelDelta = event.deltaY > 0 ? 1.1 : 0.9;
        const x = (event.clientX / canvas.width) * view.vp_w;
        const y = (event.clientY / canvas.height) * view.vp_h;
        if (view.vp_w * wheelDelta > this.view.UNIVERSE_WIDTH || view.vp_w * wheelDelta < 10 * this.view.BLOCK_WIDTH
            || view.vp_h * wheelDelta > view.UNIVERSE_WIDTH || view.vp_h * wheelDelta < 10 * view.BLOCK_WIDTH) {
        }
        else {
            view.vp_w *= wheelDelta;
            view.vp_h *= wheelDelta;
            const ny = (event.clientY / canvas.height) * view.vp_h;
            const nx = (event.clientX / canvas.width) * view.vp_w;
            const dx = x - nx;
            const dy = y - ny;
            view.vp_ox = clamp(view.vp_ox + dx, view.UNIVERSE_WIDTH - view.vp_w, 0);
            view.vp_oy = clamp(view.vp_oy + dy, view.UNIVERSE_WIDTH - view.vp_h, 0);
        }
        view.render();
    }

    handleMouseUp = () => {
        this.isDragging = false;
    }

    init() {
        this.view.canvas.addEventListener("wheel", this.handleWheel);
        this.view.canvas.addEventListener("mouseup", this.handleMouseUp);
        this.view.canvas.addEventListener("mousedown", this.handleMouseDown);
        this.view.canvas.addEventListener("mousemove", this.handleMouseMove);
    }

    cleanup() {
        this.view.canvas.removeEventListener("wheel", this.handleWheel)
        this.view.canvas.removeEventListener("mouseup", this.handleMouseUp);
        this.view.canvas.removeEventListener("mousedown", this.handleMouseDown);
        this.view.canvas.removeEventListener("mousemove", this.handleMouseMove);
    }
}
