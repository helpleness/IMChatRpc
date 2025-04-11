import { proto } from './IMChat_pb.js';
const SendGroupJoinRequest = proto.chat.SendGroupJoinRequest
// import { ChatServiceClient } from "@/stores/rpc/IMChat_grpc_web_pb.js"; // 导入生成的 ChatServiceClient
import {defineStore} from "pinia";

export const use_send_group_join_request = defineStore('join_group', {
    // State: Defines the initial state of the store
    state: () => ({
        groupId: "",
        userFrom: 'user1',  // 当前用户 ID
        reason: "",
    }),

    // Actions: Define methods that can modify the state
    actions: {

        GetGroupId(groupId){
            this.groupId=groupId;
        },
        GetUserFrom(userFrom){
            this.userFrom=userFrom;
        },
        GetReason(reason){
            this.reason=reason;
        },
        // 发送加入群组请求
        sendGroupJoinRequest(groupId,userFrom,reason,client) {
            this.GetGroupId(groupId);
            this.GetUserFrom(userFrom);
            this.GetReason(reason);
            const request = new SendGroupJoinRequest(); // 创建 gRPC 请求对象
            request.setApplicantFrom(this.userFrom);  // 设置申请人 ID
            request.setGroupId(this.groupId);  // 设置群组 ID
            request.setApplicantReason(this.reason);  // 设置申请理由
            request.setApplicantTime(Math.floor(Date.now() / 1000)); // 设置时间戳

            // // 检查 gRPC 客户端是否已初始化
            // if (!this.client) {
            //     this.initClient();
            // }

            // 调用 gRPC 的 sendGroupJoinRequest 方法
            client.sendGroupJoinRequest(request, {}, (err, response) => {
                if (err) {
                    console.error('Error sending group join request:', err.message);
                } else {
                    console.log('Group join request sent successfully:', response.toObject());
                }
            });
        }
    },
});
