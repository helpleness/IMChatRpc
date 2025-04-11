// import {SendMessageRequest} from "./IMChat_pb.js"
import { proto } from './IMChat_pb.js';
const SendMessageRequest = proto.chat.SendMessageRequest
// import { ChatServiceClient } from "@/stores/rpc/IMChat_grpc_web_pb.js"; // 导入生成的 ChatServiceClient
import {defineStore} from "pinia";
import {useChatStore} from "@/stores/chatStore.js";
export const use_send_message = defineStore('send_message', {
    // State: Defines the initial state of the store
    state: () => ({
        message: '',
        userFrom: 'user1',  // 模拟的发送方
        userTarget: 'user2', // 模拟的接收方
        err: '',
        success: '',
    }),

    // Actions: Define methods that can modify the state
    actions: {

        Getmessage(message){
            this.message=message;
        },
        GetuserFrom(userFrom){
            this.userFrom=userFrom;
        },
        GetuserTarget(userTarget){
            this.userTarget=userTarget;
        },
        sendMessage(message,userFrom,userTarget,client) {
            this.Getmessage(message);
            this.GetuserFrom(userFrom);
            this.GetuserTarget(userTarget);
            // 创建 SendMessageRequest 请求对象
            let request = new proto.chat.SendMessageRequest();
            let my_message = new proto.chat.MyMessage(); // 创建 MyMessage 对象
            // 设置 MyMessage 对象的各个字段
            my_message.setMessageId('123'); // 设置唯一的消息 ID
            my_message.setUserFrom(this.userFrom); // 设置发送者
            my_message.setSendTarget(this.userTarget); // 设置接收者
            my_message.setContent(this.message); // 设置消息内容
            my_message.setType(0); // 设置消息类型为 TEXT (假设 0 表示文本消息)
            my_message.setSendTime(Math.floor(Date.now() / 1000)); // 设置发送时间（时间戳）

            request.setMessage(my_message);// 将 MyMessage 对象设置到 SendMessageRequest 中

            // 检查 gRPC 客户端是否已初始化
            // if (!this.client) {
            //     this.initClient();
            // }

            // 使用 gRPC-Web 客户端发送消息
            client.sendMessage(request, {}, (err , response) => {
                if (err) {
                    // 错误处理，记录错误信息
                    console.error('Error sending message:', err.message);
                    this.err = err.message;
                } else {
                    // 确保 response 是有效的，执行成功的逻辑
                    this.success = response.toObject().success;
                    console.log('Message sent successfully:', response.toObject());
                    if (this.success) {
                        const chatStore = useChatStore();
                        chatStore.sendMessage(chatStore.newMessage);
                        chatStore.newMessage = ""; // 清空输入框
                    } else {
                        console.error('Message sending failed, success flag is false.');
                    }

                }
            });
        }
    },
});
