const decoder = new TextDecoder('utf-8');
const encoder = new TextEncoder();

export enum MessageType {
    PULL = "PULL",
    UPDT = "UPDT"
}

export interface PullMessage {
    messageType: string;
    height: number;
    width: number;
    imageData: Uint8Array;
}

export interface UpdateMessage {
    messageType: string;
    username: string;
    timestamp: bigint;
    pos: number;
    color: number;
}

export function parsePullMessage(arrayBuffer: ArrayBuffer): PullMessage {
    const dataView = new DataView(arrayBuffer);

    // Extract the first 4 bytes for the message type
    const messageTypeBytes = new Uint8Array(arrayBuffer, 0, 4);
    const messageType = decoder.decode(messageTypeBytes);

    // Extract the next 4 bytes for the height
    const height = dataView.getUint32(4, true);

    // Extract the next 4 bytes for the width
    const width = dataView.getUint32(8, true);

    // Calculate the starting point of the image data
    const imageDataStart = 12;
    const imageDataLength = height * width;

    const imageData = new Uint8Array(arrayBuffer, imageDataStart, imageDataLength);
    return {
        messageType,
        height,
        width,
        imageData,
    }
}

export function parseUpdateMessage(arrayBuffer: ArrayBuffer): UpdateMessage {
    const dataView = new DataView(arrayBuffer);

    // Extract the first 4 bytes for the message type
    const messageTypeBytes = new Uint8Array(arrayBuffer, 0, 4);
    const messageType = decoder.decode(messageTypeBytes);

    // Extract the next 64 bytes for the username
    const usernameBytes = new Uint8Array(arrayBuffer, 4, 64);
    const username = decoder.decode(usernameBytes).replace(/\0/g, ''); // Remove null bytes

    // Extract the next 8 bytes for the timestamp
    const timestamp = dataView.getBigUint64(68, true);

    // Extract the next 4 bytes for the position
    const pos = dataView.getUint32(76, true);

    // Extract the color
    const color = dataView.getUint8(80);

    return {
        messageType,
        username,
        timestamp,
        pos,
        color
    };
}

export function packUpdateMessage(message: UpdateMessage): ArrayBuffer {
    const buffer = new ArrayBuffer(81); // Total size is 4 + 64 + 8 + 4 + 1 = 81 bytes
    const dataView = new DataView(buffer);

    // Pack the first 4 bytes with the message type
    const messageTypeBytes = encoder.encode(message.messageType);
    for (let i = 0; i < 4; i++) {
        dataView.setUint8(i, i < messageTypeBytes.length ? messageTypeBytes[i] : 0); // Pad with null bytes if needed
    }

    // Pack the next 64 bytes with the username
    const usernameBytes = encoder.encode(message.username);
    for (let i = 0; i < 64; i++) {
        dataView.setUint8(4 + i, i < usernameBytes.length ? usernameBytes[i] : 0); // Pad with null bytes if needed
    }

    // Pack the next 8 bytes with the timestamp
    dataView.setBigUint64(68, message.timestamp, true); // true for little-endian

    // Pack the next 4 bytes with the uint32
    dataView.setUint32(76, message.pos, true); // true for little-endian

    // Pack the final byte with the uint8
    dataView.setUint8(80, message.color);

    return buffer;
}

export function checkMessageType(arrayBuffer: ArrayBuffer): MessageType {
    // Extract the first 4 bytes for the message type
    const messageTypeBytes = new Uint8Array(arrayBuffer, 0, 4);
    const messageType = decoder.decode(messageTypeBytes);
    return messageType as MessageType
}
