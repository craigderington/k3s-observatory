import { useEffect, useRef, useCallback, useState } from 'react';

interface WebSocketEvent {
  type: string;
  data: any;
}

interface UseWebSocketOptions {
  url: string;
  onMessage?: (event: WebSocketEvent) => void;
  onConnect?: () => void;
  onDisconnect?: () => void;
  onError?: (error: Event) => void;
  reconnectDelay?: number;
}

export function useWebSocket({
  url,
  onMessage,
  onConnect,
  onDisconnect,
  onError,
  reconnectDelay = 3000,
}: UseWebSocketOptions) {
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout>();
  const [isConnected, setIsConnected] = useState(false);
  const isConnectingRef = useRef(false);
  const shouldReconnectRef = useRef(true);

  // Store callbacks in refs to avoid reconnecting when they change
  const onMessageRef = useRef(onMessage);
  const onConnectRef = useRef(onConnect);
  const onDisconnectRef = useRef(onDisconnect);
  const onErrorRef = useRef(onError);

  // Update refs when callbacks change
  useEffect(() => {
    onMessageRef.current = onMessage;
    onConnectRef.current = onConnect;
    onDisconnectRef.current = onDisconnect;
    onErrorRef.current = onError;
  }, [onMessage, onConnect, onDisconnect, onError]);

  const connect = useCallback(() => {
    // Prevent multiple simultaneous connection attempts
    if (isConnectingRef.current) {
      console.log('Connection attempt already in progress, skipping...');
      return;
    }

    // Don't reconnect if we shouldn't
    if (!shouldReconnectRef.current) {
      console.log('Reconnection disabled, skipping...');
      return;
    }

    // Clear any pending reconnection attempts
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = undefined;
    }

    // Close existing connection if any
    if (wsRef.current && wsRef.current.readyState !== WebSocket.CLOSED) {
      wsRef.current.close();
      wsRef.current = null;
    }

    isConnectingRef.current = true;
    console.log(`[WebSocket] Connecting to ${url}`);
    const ws = new WebSocket(url);

    // Send periodic heartbeat to keep connection alive
    let heartbeatInterval: NodeJS.Timeout;

    ws.onopen = () => {
      console.log('[WebSocket] Connected successfully');
      isConnectingRef.current = false;
      setIsConnected(true);
      onConnectRef.current?.();

      // Send heartbeat every 30 seconds to keep connection alive
      heartbeatInterval = setInterval(() => {
        if (ws.readyState === WebSocket.OPEN) {
          console.log('[WebSocket] Sending heartbeat');
          ws.send(JSON.stringify({ type: 'ping' }));
        }
      }, 30000);
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log('[WebSocket] Received message:', data.type);
        onMessageRef.current?.(data);
      } catch (error) {
        console.error('[WebSocket] Failed to parse message:', error);
      }
    };

    ws.onerror = (error) => {
      console.error('[WebSocket] Error occurred:', error);
      isConnectingRef.current = false;
      onErrorRef.current?.(error);
      if (heartbeatInterval) {
        clearInterval(heartbeatInterval);
      }
    };

    ws.onclose = (event) => {
      console.log(`[WebSocket] Disconnected - Code: ${event.code}, Reason: ${event.reason || 'none'}, Clean: ${event.wasClean}`);
      isConnectingRef.current = false;
      setIsConnected(false);
      onDisconnectRef.current?.();

      if (heartbeatInterval) {
        clearInterval(heartbeatInterval);
      }

      // Only attempt to reconnect if we should
      if (shouldReconnectRef.current) {
        console.log(`[WebSocket] Scheduling reconnection in ${reconnectDelay}ms...`);
        reconnectTimeoutRef.current = setTimeout(() => {
          connect();
        }, reconnectDelay);
      }
    };

    wsRef.current = ws;
  }, [url, reconnectDelay]);

  useEffect(() => {
    shouldReconnectRef.current = true;
    connect();

    return () => {
      console.log('Cleaning up WebSocket connection...');
      shouldReconnectRef.current = false;
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
        reconnectTimeoutRef.current = undefined;
      }
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, [connect]);

  return {
    send: useCallback((data: any) => {
      if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify(data));
      }
    }, []),
    isConnected,
  };
}
