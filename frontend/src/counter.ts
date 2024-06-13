import { MessageType, UpdateMessage, checkMessageType, packUpdateMessage, parseUpdateMessage, parsePullMessage } from "./message";

let socket: WebSocket | null = null;

export function setupButton(element: HTMLButtonElement) {
    let connected: Boolean = false
    const handleClick = () => {
        if (connected) {
            let updMessage: UpdateMessage = {
                messageType: "UPDT",
                username: "Pratham",
                timestamp: BigInt(Date.now()),
                pos: 0,
                color: 3
            }
            console.log("Sending")
            socket?.send(packUpdateMessage(updMessage));
            // socket?.close();
            // socket = null;
            // connected = false;
            // element.textContent = "Connect";
        } else {
            socket = new WebSocket("ws://localhost:8000/ws")
            socket.binaryType = "arraybuffer"
            connected = true;
            socket.onmessage = (event) => {
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
            element.textContent = "Disconnect"
        }
    }
    element.addEventListener('click', () => handleClick())
}
// export function setupCounter(element: HTMLButtonElement) {
//   let counter = 0
//   const setCounter = (count: number) => {
//     counter = count
//     element.innerHTML = `count is ${counter}`
//   }
//   element.addEventListener('click', () => setCounter(counter + 1))
//   setCounter(0)
// }
