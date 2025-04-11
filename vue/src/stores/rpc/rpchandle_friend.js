// import {HandleFriendRequest} from "./IMChat_pb.js"
import { proto } from './IMChat_pb.js';
const HandleFriendRequest = proto.chat.HandleFriendRequest
// import { ChatServiceClient } from "@/stores/rpc/IMChat_grpc_web_pb.js"; // 导入生成的 ChatServiceClient
import {defineStore} from "pinia";

export const use_handle_friend_request = defineStore('handle_friend', {
    // State: Defines the initial state of the store
    state: () => ({
        applicantId: '',
        userTarget: 'user1' // 当前用户 ID
    }),

    // Actions: Define methods that can modify the state
    actions: {

        GetapplicantId(applicantId){
            this.applicantId=applicantId;
        },
        GetuserTarget(userTarget){
            this.userTarget=userTarget;
        },
        // 处理好友请求
        handleFriendRequest(applicantId,userTarget,accept,client) {
            this.GetapplicantId(applicantId);
            this.GetuserTarget(userTarget);
            const request = new HandleFriendRequest(); // 创建 gRPC 请求对象
            request.setApplicantFrom(this.applicantId); // 设置申请人 ID
            request.setApplicantTarget(this.userTarget); // 设置目标用户 ID
            request.setAccept(accept); // 设置是否接受

            // 检查 gRPC 客户端是否已初始化
            // if (!this.client) {
            //     this.initClient();
            // }

            // 调用 gRPC-Web 的 handleFriendRequest 方法
            client.handleFriendRequest(request, {}, (err, response) => {
                if (err) {
                    console.error('Error handling friend request:', err.message);
                } else {
                    console.log('Friend request handled successfully:', response.toObject());
                }
            });
        },
    },
});
