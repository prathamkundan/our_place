import { createContext, useContext, useEffect, useState } from "react";
import { WebSocketController } from "../utils/network/socket_control";

export const WSContext = createContext<WebSocketController | null>(null);

export const useWebSocket = () => {
    return useContext(WSContext);
}

const WSProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [wsController, setWsController] = useState<WebSocketController | null>(null)
    useEffect(() => {
        setWsController(new WebSocketController());
        return () => {
            wsController?.cleanup()
        }
    }, []);

    return (
        <WSContext.Provider value={wsController}>
            {children}
        </WSContext.Provider>
    )
}

export default WSProvider;
