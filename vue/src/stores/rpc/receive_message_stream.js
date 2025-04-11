// import { MessageStreamRequest } from "@/stores/rpc/IMChat_pb.js";
import { proto } from './IMChat_pb.js';
const MessageStreamRequest = proto.chat.MessageStreamRequest
// import { ChatServiceClient } from "@/stores/rpc/IMChat_grpc_web_pb.js"; // 导入生成的 gRPC 客户端
import {defineStore} from "pinia";
import {useChatStore} from "@/stores/chatStore.js";


export const use_receive_message_stream = defineStore('receive_message', {
    // State: Defines the initial state of the store
    state: () => ({
        stream: null,
    }),

    // Actions: Define methods that can modify the state
    actions: {
        // 开始接收消息流
        startReceivingMessages(client) {
            const chatStore = useChatStore();
            // 创建 gRPC 客户端
            // const client = new ChatServiceClient('http://localhost:9080', null, null);

            // 创建消息流请求对象
            let messageStreamRequest = new MessageStreamRequest();
            messageStreamRequest.setUserId(chatStore.userID);  // 设置用户 ID

            // 使用 gRPC 客户端的流式方法来接收消息流
            this.stream = client.receiveMessageStream(messageStreamRequest, {});

            this.stream.on('data', (response) => {
                const message = response.toObject();  // 将 gRPC 响应转为 JavaScript 对象
                // 只处理发送给当前用户的消息
                if (message.sendTarget === chatStore.userID) {
                    chatStore.receiveMessage(message.user_from,message)
                }
                // chatStore.receiveMessage(this.userId,message)
                // this.messages.push(message);          // 将新消息添加到消息列表中
            });

            this.stream.on('error', (error) => {
                console.error('Stream error:', error);
            });

            this.stream.on('end', () => {
                console.log('Stream ended');
            });
        },

        // 关闭消息流
        stopReceivingMessages() {
            if (this.stream) {
                this.stream.cancel();  // 取消流
                this.stream = null;
            }
        },
    },
    // 销毁时关闭消息流
    beforeDestroy() {
        this.stopReceivingMessages();
    },

});
