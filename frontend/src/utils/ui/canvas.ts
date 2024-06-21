export const colorMap: { [key: number]: string } = {
    0: '#FFFFFF', // WHITE
    1: '#D3D3D3', // LGRAY
    2: '#808080', // GRAY
    3: '#000000', // BLACK
    4: '#FFB6C1', // KIRBY (Light Pink)
    5: '#FF0000', // RED
    6: '#FFA500', // ORANGE
    7: '#A52A2A', // BROWN
    8: '#FFFF00', // YELLOW
    9: '#90EE90', // lGREEN (Light Green)
    10: '#008000', // GREEN
    11: '#ADD8E6', // LBLUE (Light Blue)
    12: '#40E0D0', // TURQUOISE
    13: '#0000FF', // BLUE
    14: '#FF1893', // PINK
    15: '#8F00FF'  // VIOLET
};

/**
* # The Grid Canvas
* It is represented as a Square of width UNIVERSE_WIDTH where the width of each block is BLOCK_WIDTH.
* We define a viewport that has its top left corner at `(vp_ox, vp_oy)` wrt to the top left corner of the universe
* with width and height `vp_h` and `vp_w`.
*
* Only the cells inside the viewport are visible to the user and drawn by the `drawGrid()` funtion.
*
*/
export class View {
    public canvas: HTMLCanvasElement;
    public ctx: CanvasRenderingContext2D;
    public grid: Uint8Array | null;

    // The viewport properties position of the top left corner of the viewport, its width and height
    public vp_ox;
    public vp_oy;

    public vp_w: number;
    public vp_h: number;

    public UNIVERSE_WIDTH;
    public NUM_ROWS;
    public BLOCK_WIDTH;
    public MODE: string = "NORMAL"

    public BG_COLOR: string
    public BDR_COLOR: string
    public HIGHLIGHT_COLOR: string

    constructor(canvas: HTMLCanvasElement, block_size: number, rows: number) {
        this.canvas = canvas
        const parent = this.canvas.parentNode as HTMLElement
        this.canvas.height = parent.clientHeight
        this.canvas.width = parent.clientWidth;

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

        this.BG_COLOR = "#FFFFFF"
        this.BDR_COLOR = "#D9D9D9"
        this.HIGHLIGHT_COLOR = "#898E96"
    }

    setGrid = (grid: Uint8Array) => {
        this.grid = grid;
        this.render();
    }

    setCanvasDimensions(width: number, height: number) {
        this.canvas.width = width;
        this.canvas.height = height;
        this.vp_w = width;
        this.vp_h = height;
    }

    updateGrid(pos: number, color: number) {
        this.grid![pos] = color;
        this.render()
    }

    highLight(x: number, y: number) {
        const k = this.canvas.width / this.vp_w
        const base_x = Math.ceil(this.vp_ox / this.BLOCK_WIDTH)
        const base_y = Math.ceil(this.vp_oy / this.BLOCK_WIDTH)
        const last_x = base_x + Math.ceil(this.vp_w / this.BLOCK_WIDTH)
        const last_y = base_y + Math.ceil(this.vp_h / this.BLOCK_WIDTH)

        const pos_x = base_x * this.BLOCK_WIDTH - this.vp_ox + (x - base_x) * this.BLOCK_WIDTH
        const pos_y = base_y * this.BLOCK_WIDTH - this.vp_oy + (y - base_y) * this.BLOCK_WIDTH
        if (x < base_x || x > last_x || y < base_y || y > last_y) return
        else {
            this.ctx.strokeStyle = this.grid![this.toIndex(x, y)] === 3 ? "white" : "black"
            this.ctx.lineWidth = 1 * k;
            this.ctx.strokeRect(k * (pos_x + 1), k * (pos_y + 1), k * (this.BLOCK_WIDTH - 2), k * (this.BLOCK_WIDTH - 2))
        }

        this.ctx.strokeStyle = "white";
        this.ctx.lineWidth = 0;
    }

    clearRect = (x: number, y: number) => {
        const k = this.canvas.width / this.vp_w
        const base_x = Math.ceil(this.vp_ox / this.BLOCK_WIDTH)
        const base_y = Math.ceil(this.vp_oy / this.BLOCK_WIDTH)
        const last_x = base_x + Math.ceil(this.vp_w / this.BLOCK_WIDTH)
        const last_y = base_y + Math.ceil(this.vp_h / this.BLOCK_WIDTH)

        const pos_x = base_x * this.BLOCK_WIDTH - this.vp_ox + (x - base_x) * this.BLOCK_WIDTH
        const pos_y = base_y * this.BLOCK_WIDTH - this.vp_oy + (y - base_y) * this.BLOCK_WIDTH
        if (x < base_x || x > last_x || y < base_y || y > last_y) return

        else {
            this.ctx.beginPath();
            this.ctx.clearRect(k * pos_x, k * pos_y, k * this.BLOCK_WIDTH, k * this.BLOCK_WIDTH);
            this.ctx.fillStyle = colorMap[this.grid![this.toIndex(x, y)]];
            this.ctx.fillRect(k * pos_x, k * pos_y, k * this.BLOCK_WIDTH, k * this.BLOCK_WIDTH);
            this.ctx.closePath();
        }
    }

    render = () => {
        let x_ind = Math.floor(this.vp_ox / this.BLOCK_WIDTH);
        const x_ind_start = x_ind
        const x_coord_start = x_ind * this.BLOCK_WIDTH - this.vp_ox;
        // we subtract vp_ox as it is the place we must consider it wrt the screen

        let y_ind = Math.floor(this.vp_oy / this.BLOCK_WIDTH);
        const y_ind_start = y_ind;
        const y_coord_start = y_ind * this.BLOCK_WIDTH - this.vp_oy;

        const k = this.canvas.width / this.vp_w;
        const transformed_block_width = this.BLOCK_WIDTH * k;
        // This is the ratio by which the width of blocks is scaled when they appear on screen

        this.ctx.beginPath()
        for (let c = 0; c < 16; c++) {
            this.ctx.fillStyle = colorMap[c]
            for (let x = k * x_coord_start; x < this.vp_w * k; x += transformed_block_width) {
                for (let y = k * y_coord_start; y < this.vp_h * k; y += transformed_block_width) {
                    if (this.grid![this.toIndex(x_ind, y_ind)] === c)
                        this.ctx.fillRect(x, y, k * this.BLOCK_WIDTH, k * this.BLOCK_WIDTH);
                    y_ind++;
                }
                y_ind = y_ind_start;
                x_ind++;
            }
            x_ind = x_ind_start
        }
        this.ctx.closePath()
    }

    locToIndex = (x: number, y: number) => {
        const k = this.canvas.width / this.vp_w;
        const transformed_block_width = k * this.BLOCK_WIDTH;
        // This is the ratio by which the width of blocks is scaled when they appear on screen

        const base_x = this.vp_ox / this.BLOCK_WIDTH;
        const base_y = this.vp_oy / this.BLOCK_WIDTH;
        // This is the index of the top-left-est block

        const diff_x = x / transformed_block_width
        const diff_y = y / transformed_block_width

        return [Math.floor(base_x + diff_x), Math.floor(base_y + diff_y)];
    }

    toIndex(x: number, y: number) {
        return x * this.NUM_ROWS + y;
    }

}
