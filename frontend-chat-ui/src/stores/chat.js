import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useChatStore = defineStore('chat', () => {
    // --- 这是一个模拟数据库 ---
    const friends = ref([
        { id: 'friend1', name: 'Alice', avatar: 'https://api.dicebear.com/8.x/initials/svg?seed=Alice' },
        { id: 'friend2', name: 'Bob', avatar: 'https://api.dicebear.com/8.x/initials/svg?seed=Bob' },
    ]);

    const groups = ref([
        { id: 'group1', name: 'GoLang Enthusiasts', members: ['user', 'friend1'], avatar: 'https://api.dicebear.com/8.x/icons/svg?seed=Go' },
        { id: 'group2', name: 'Vue.js Fans', members: ['user', 'friend2'], avatar: 'https://api.dicebear.com/8.x/icons/svg?seed=Vue' },
    ]);

    // 当前激活的聊天窗口 (可以是好友ID或群组ID)
    const activeChatId = ref('group1');

    const activeChat = ref({
        id: 'group1',
        name: 'GoLang Enthusiasts',
        type: 'group',
        avatar: 'https://api.dicebear.com/8.x/icons/svg?seed=Go'
    });

    // 消息 (真实应用中每个聊天都应有自己的消息列表)
    // 这里我们简化，所有消息都放在一个列表里，并连接到之前的Go后端
    const messages = ref([]);

    function setActiveChat(chat) {
        activeChatId.value = chat.id;
        activeChat.value = { ...chat, type: chat.members ? 'group' : 'friend' };
        // 在真实应用中，切换聊天时应加载该聊天的历史消息
        // messages.value = loadMessagesFor(chat.id);
        console.log(`Switched to chat: ${chat.name}`);
    }

    function addFriend(name) {
        // 模拟添加好友
        const newFriend = {
            id: `friend_${Date.now()}`,
            name,
            avatar: `https://api.dicebear.com/8.x/initials/svg?seed=${name}`
        };
        friends.value.push(newFriend);
        alert(`Friend "${name}" added! (UI only)`);
    }

    function createGroup(name) {
        // 模拟创建群组
        const newGroup = {
            id: `group_${Date.now()}`,
            name,
            members: ['user'],
            avatar: `https://api.dicebear.com/8.x/icons/svg?seed=${name}`
        };
        groups.value.push(newGroup);
        alert(`Group "${name}" created! (UI only)`);
    }

    return { friends, groups, activeChatId, activeChat, messages, setActiveChat, addFriend, createGroup };
});