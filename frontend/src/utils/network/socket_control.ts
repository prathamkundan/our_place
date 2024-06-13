import { MessageType, checkMessageType, parsePullMessage, parseUpdateMessage } from "./message";

export class WebSocketController {
    public socket: WebSocket | null;
    public username: string | null;

    constructor() {
        this.socket = null
        this.username = null
    }

    onMessage(event: MessageEvent) {
        let messageType = checkMessageType(event.data)
        console.log(messageType)
        if (messageType == MessageType.PULL) {
            let data = parsePullMessage(event.data as ArrayBuffer)
            console.log(data)
        } else if (messageType == MessageType.UPDT) {
            let data = parseUpdateMessage(event.data as ArrayBuffer)
            console.log(data)
        }
    }
}
