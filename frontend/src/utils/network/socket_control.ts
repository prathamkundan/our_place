import { View } from "../ui/canvas";
import { MessageType, UpdateMessage, checkMessageType, packUpdateMessage, parsePullMessage, parseUpdateMessage } from "./message";

export class WebSocketController {
    public view: View | null = null;
    public socket: WebSocket | null;
    public username: string | null;

    constructor() {
        this.socket = null
        this.username = null
    }

    onMessage = (event: MessageEvent) => {
        const messageType = checkMessageType(event.data as ArrayBuffer)
        // console.log(messageType)
        if (messageType == MessageType.PULL) {
            const data = parsePullMessage(event.data as ArrayBuffer)
            // console.log(data)
            this.view?.setGrid(data.imageData);
        } else if (messageType == MessageType.UPDT) {
            const data = parseUpdateMessage(event.data as ArrayBuffer)
            // console.log("Got Message: ", data);
            this.view!.updateGrid(data.pos, data.color)
        }
    }

    sendUpdate = (pos: number, color: number) => {
        const updateMessage: UpdateMessage = {
            messageType: MessageType.UPDT,
            timestamp: BigInt(Math.ceil(Date.now() / 1000)),
            username: this.username == null ? "" : this.username.slice(0, 64),
            color: color,
            pos: pos
        }
        // console.log("Sending:", updateMessage);
        this.socket?.send(packUpdateMessage(updateMessage));
    }

    init = (url: string | URL, view: View) => {
        this.socket = new WebSocket(url);
        this.socket.binaryType = "arraybuffer"
        this.view = view;

        this.socket.onmessage = this.onMessage
    }

    cleanup = () => {
        this.socket?.close();
        this.view = null;

    }

}
