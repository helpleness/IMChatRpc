<script setup>
import Sidebar from '@/components/Sidebar.vue';
import ChatHeader from '@/components/ChatHeader.vue';
import ChatWindow from '@/components/ChatWindow.vue';
import MessageInput from '@/components/MessageInput.vue';
import { useChatStore } from '@/stores/chat';
import { onMounted } from 'vue';

const chatStore = useChatStore();
const API_URL = 'http://localhost:8080/api/messages';

// 长轮询逻辑，与之前的 demo 类似
const startLongPolling = async () => {
  console.log('Starting long poll for messages...');
  while (true) {
    try {
      const response = await fetch(`${API_URL}?since=${chatStore.messages.length}`);
      if (response.status === 200) {
        const newMessages = await response.json();
        if (newMessages && newMessages.length > 0) {
          // 在真实应用中，你需要根据 activeChatId 来判断消息是否属于当前窗口
          chatStore.messages.push(...newMessages.map(msg => ({ text: msg, sender: 'other' })));
        }
      } else if (response.status !== 204) {
        throw new Error(`Server responded with status: ${response.status}`);
      }
    } catch (e) {
      console.error('Long polling error:', e);
      await new Promise(resolve => setTimeout(resolve, 5000));
    }
  }
};

// 发送消息的方法，传递给 MessageInput 组件
const handleSendMessage = async (text, toUser) => {
  const fromUser = authStore.user.name; // 假设用户名就是用户ID
  try {
    // 假设用户可能连接到不同端口的服务器
    const API_BASE_URL = 'http://localhost:8080'; // 或者 8081
    const response = await fetch(`${API_BASE_URL}/api/send`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        from: fromUser,
        to: toUser,
        content: text
      }),
    });
    // ...
  } catch(e){
    console.error('Post message error:', e);
  }
};

onMounted(() => {
  startLongPolling();
});
</script>

<template>
  <div class="chat-layout">
    <Sidebar />
    <main class="chat-main">
      <ChatHeader />
      <ChatWindow :messages="chatStore.messages" />
      <MessageInput @sendMessage="handleSendMessage" />
    </main>
  </div>
</template>

<style scoped>
.chat-layout {
  display: flex;
  height: 100vh;
  width: 100vw;
  background-color: var(--primary-bg);
}
.chat-main {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  height: 100vh;
}
</style>