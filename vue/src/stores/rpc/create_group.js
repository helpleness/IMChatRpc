import { defineStore } from "pinia";
// import { proto } from './IMChat_grpc_web_pb.js';
// const ChatServiceClient = proto.ChatServiceClient
import { proto } from './IMChat_pb.js';
const CreateGroupRequest = proto.chat.CreateGroupRequest
const GroupCreation = proto.chat.GroupCreation


export const use_create_group = defineStore('newgroup', {
    // State: Defines the initial state of the store
    state: () => ({
        groupName: '',
        userFrom: 'user1',  // 当前用户 ID
        members: ['user2', 'user3'], // 群成员
    }),

    // Actions: Define methods that can modify the state
    actions: {
        GetGroupName(groupName) {
            this.groupName = groupName;
        },
        GetUserFrom(userFrom) {
            this.userFrom = userFrom;
        },
        GetMembers(members) {
            this.members = members;
        },

        // 创建群聊
        createGroup(groupName,userFrom,members,client) {
            this.GetGroupName(groupName);
            this.GetUserFrom(userFrom);
            this.GetMembers(members);
            // 创建 GroupCreation 消息
            const groupCreation = new GroupCreation();
            groupCreation.setCreatorId(this.userFrom);
            groupCreation.setGroupName(this.groupName);
            groupCreation.setCreationTime(Math.floor(Date.now() / 1000));
            this.members.forEach(member => groupCreation.addMembers(member)); // 添加成员

            // 创建 CreateGroupRequest 消息
            const createGroupRequest = new CreateGroupRequest();
            createGroupRequest.setGroupCreation(groupCreation);

            // 使用 gRPC 客户端发送请求
            client.createGroup(createGroupRequest, {}, (err, response) => {
                if (err) {
                    console.error('Error creating group:', err);
                    return;
                }
                console.log('Group created:', response.getSuccess());
                // 可以处理成功回调
            });
        }
    },
});
