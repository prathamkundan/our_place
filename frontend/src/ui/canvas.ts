export class View {
    public canvas: HTMLCanvasElement;
    public ctx: CanvasRenderingContext2D;
    public grid: Uint8Array | null;

    public vp_ox;
    public vp_oy;
    public vp_w: number;
    public vp_h: number;

    public UNIVERSE_WIDTH;
    public NUM_ROWS;
    public BLOCK_WIDTH;
    public MODE: string = "NORMAL"

    constructor(id: string, width: number, height: number, block_size: number, rows: number) {
        this.canvas = document.getElementById(id)! as HTMLCanvasElement
        this.canvas.height = height;
        this.canvas.width = width;

        this.ctx = this.canvas.getContext('2d')!;
        this.grid = null;

        this.vp_ox = 0;
        this.vp_oy = 0;
        this.vp_w = this.canvas.width;
        this.vp_h = this.canvas.height;

        this.BLOCK_WIDTH = block_size;
        this.NUM_ROWS = rows
        this.UNIVERSE_WIDTH = block_size * this.NUM_ROWS;
        this.grid = new Uint8Array(this.NUM_ROWS * this.NUM_ROWS);

        this.updateGrid();
    }

    setCanvasDimensions(width: number, height: number) {
        this.canvas.width = width;
        this.canvas.height = height;
        this.vp_w = width;
        this.vp_h = height;
    }

    updateGrid() {
    }


    drawGrid() {
        this.updateGrid()
        this.ctx.fillStyle = "#FFFFFF";
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
        this.ctx.beginPath();

        let x_start = Math.ceil(this.vp_ox / this.BLOCK_WIDTH) * this.BLOCK_WIDTH - this.vp_ox;
        let y_start = Math.ceil(this.vp_oy / this.BLOCK_WIDTH) * this.BLOCK_WIDTH - this.vp_oy;

        let k = this.canvas.width / this.vp_w;

        let transformed_block_width = this.BLOCK_WIDTH * k;
        for (let x = k * x_start; x < this.vp_w * k; x += transformed_block_width) {
            this.ctx.moveTo(x, 0);
            this.ctx.lineTo(x, k * this.vp_h);
        }

        for (let y = k * y_start; y < this.vp_h * k; y += transformed_block_width) {
            this.ctx.moveTo(0, y);
            this.ctx.lineTo(k * this.vp_w, y);
        }
        this.ctx.strokeStyle = "#212020";
        this.ctx.stroke();
        this.drawCells()
    }

    drawCells() {
        let x_ind = Math.floor(this.vp_ox / this.BLOCK_WIDTH);
        let x_coord_start = x_ind * this.BLOCK_WIDTH - this.vp_ox;

        let y_ind = Math.floor(this.vp_oy / this.BLOCK_WIDTH);
        let y_ind_start = y_ind;
        let y_coord_start = y_ind * this.BLOCK_WIDTH - this.vp_oy;

        let k = this.canvas.width / this.vp_w;
        let transformed_block_width = this.BLOCK_WIDTH * k;
        this.ctx.fillStyle = '#FFFFFF';

        for (let x = k * x_coord_start; x < this.vp_w * k; x += transformed_block_width) {
            for (let y = k * y_coord_start; y < this.vp_h * k; y += transformed_block_width) {
                if (this.grid![this.toIndex(x_ind, y_ind)] == 1) {
                    this.ctx.fillRect(x, y, k * this.BLOCK_WIDTH, k * this.BLOCK_WIDTH);
                    // let centerX = x + radius;
                    // let centerY = y + radius;
                    // // Draw filled circle
                    // this.ctx.beginPath();
                    // this.ctx.arc(centerX, centerY, radius * 1.33333, 0, 2 * Math.PI);
                    // this.ctx.fill();
                }
                y_ind++;
            }
            y_ind = y_ind_start;
            x_ind++;
        }
    }

    locToIndex(x: number, y: number) {
        let k = this.canvas.width / this.vp_w;
        let transformed_block_width = k * this.BLOCK_WIDTH;
        let base_x = this.vp_ox / this.BLOCK_WIDTH;
        let base_y = this.vp_oy / this.BLOCK_WIDTH;

        let diff_x = x / transformed_block_width
        let diff_y = y / transformed_block_width

        return [Math.floor(base_x + diff_x), Math.floor(base_y + diff_y)];
    }

    toIndex(x: number, y: number) {
        return x * this.NUM_ROWS + y;
    }

    public render() {
        this.drawGrid();
    }
}
