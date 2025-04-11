import { proto } from './IMChat_pb.js';
const HandleGroupJoinRequest = proto.chat.HandleGroupJoinRequest
const GroupJoinRequest = proto.chat.GroupJoinRequest
// import { ChatServiceClient } from "@/stores/rpc/IMChat_grpc_web_pb.js"; // 导入生成的 ChatServiceClient
import {defineStore} from "pinia";

export const use_handle_group_join_request = defineStore('handle_group', {
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

        handleGroupJoinRequest(groupId,userFrom,reason,accept,client) {
            this.GetGroupId(groupId);
            this.GetUserFrom(userFrom);
            this.GetReason(reason);
            // 创建GroupJoinRequest对象
            let groupJoinRequest = new GroupJoinRequest();
            groupJoinRequest.setApplicantFrom(this.userFrom);  // 设置申请人ID
            groupJoinRequest.setGroupId(this.groupId);         // 设置群聊ID
            groupJoinRequest.setAccept(accept);                // 是否同意加入
            groupJoinRequest.setApplicantReason(this.reason);  // 设置申请原因
            groupJoinRequest.setApplicantTime(Math.floor(Date.now() / 1000)); // 设置申请时间戳

            // 创建HandleGroupJoinRequest对象
            let handleRequest = new HandleGroupJoinRequest();
            handleRequest.setGroupJoinRequest(groupJoinRequest);

            // 使用gRPC客户端发送请求
            // const client = new ChatServiceClient('http://localhost:9080', null, null);

            client.handleGroupJoinRequest(handleRequest, {}, (err, response) => {
                if (err) {
                    console.error('Error sending group join request:', err);
                } else {
                    console.log('Group join request sent:', response.toObject());
                }
            });
        }
    },
});
