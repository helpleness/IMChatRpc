import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useUserStore } from "@/stores/user.js";
import { useChatStore } from "@/stores/chatStore.js";

// 构建 WebSocket 连接的 URL
const urlCallback = () => {
    const chatStore = useChatStore();
    if (!chatStore.isLoggedIn) {
        window.location.reload();
    }

    const socketApiUrl = import.meta.env.VITE_SOCKET_API;

    return `${socketApiUrl}/ws/default.io`;
};

export const useWebSocketStore = defineStore('websocket', () => {
    const ws = ref(null);
    const isConnected = ref(false);
    // const messages = ref([]);

    // 连接 WebSocket
    const connect = () => {
        if (ws.value) {
            console.warn('WebSocket 已经连接.');
            return;
        }

        ws.value = new WebSocket(urlCallback()); // 使用动态生成的 URL

        ws.value.onopen = () => {
            isConnected.value = true;
            console.log('WebSocket 连接成功.');
            sendMessage('连接已建立');
        };

        ws.value.onclose = () => {
            isConnected.value = false;
            console.log('WebSocket 连接关闭.');
            reconnectWebSocket(); // 自动重连
        };

        ws.value.onerror = (error) => {
            console.error('WebSocket 出现错误:', error);
            isConnected.value = false;
            reconnectWebSocket(); // 出现错误时自动重连
        };

        ws.value.onmessage = (event) => {
            console.log('收到消息:', event.data);
            messages.value.push(event.data);
        };
    };

    const disconnect = () => {
        if (!ws.value) {
            console.warn('WebSocket 未连接.');
            return;
        }
        ws.value.close();
        ws.value = null;
        isConnected.value = false;
    };

    const sendMessage = (message) => {
        if (!ws.value || ws.value.readyState !== WebSocket.OPEN) {
            console.warn('无法发送消息，WebSocket 未连接.');
            return;
        }
        ws.value.send(message);
        console.log('发送消息:', message);
    };

    // 自动重连
    let reconnectInterval;
    const reconnectWebSocket = () => {
        clearInterval(reconnectInterval); // 清除之前的重连尝试
        reconnectInterval = setInterval(() => {
            console.log('尝试重新连接 WebSocket...');
            connect(); // 调用 WebSocket 连接方法
        }, 5000); // 每5秒尝试重连
    };

    return {
        ws,
        isConnected,
        connect,
        disconnect,
        sendMessage,
    };
});
