import { proto } from './IMChat_pb.js';
const SendFriendRequest = proto.chat.SendFriendRequest
// import { ChatServiceClient } from "@/stores/rpc/IMChat_grpc_web_pb.js"; // 导入生成的 gRPC 客户端
import {defineStore} from "pinia";

export const use_send_friend_request = defineStore('new_friend', {
    // State: Defines the initial state of the store
    state: () => ({
        friendId: "",
        userFrom: 'user1',  // 当前用户 ID
        reason: "",

    }),

    // Actions: Define methods that can modify the state
    actions: {

        GetFriendId(friendid){
            this.friendId=friendid;
        },
        GetUserFrom(userFrom){
            this.userFrom=userFrom;
        },
        GetReason(reason){
            this.reason=reason;
        },
        sendFriendRequest(friendid,userFrom,reason,client) {
            this,this.GetFriendId(friendid);
            this.GetUserFrom(userFrom);
            this.GetReason(reason);
            const request = new SendFriendRequest(); // 创建 gRPC 请求对象
            const FriendRequest =new proto.chat.FriendRequest();
            FriendRequest.setApplicantFrom(this.userFrom);  // 设置申请人 ID
            FriendRequest.setApplicantTarget(this.friendId);  // 设置好友目标 ID
            FriendRequest.setApplicantReason(this.reason);  // 设置申请理由
            FriendRequest.setApplicantTime(Math.floor(Date.now() / 1000)); // 设置时间戳

            // 调用 gRPC 的 sendFriendRequest 方法
            client.sendFriendRequest(FriendRequest, {}, (err, response) => {
                if (err) {
                    console.error('Error sending friend request:', err.message);
                } else {
                    console.log('Friend request sent successfully:', response.toObject());
                }
            });
        }

    },
});
