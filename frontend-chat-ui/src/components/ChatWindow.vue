<script setup>
import { ref, watch, nextTick } from 'vue';
import { useAuthStore } from '@/stores/auth';

const props = defineProps({
  messages: {
    type: Array,
    required: true,
  }
});

const authStore = useAuthStore();
const messageListEl = ref(null);

// 自动滚动到底部
watch(() => props.messages, async () => {
  await nextTick();
  if (messageListEl.value) {
    messageListEl.value.scrollTop = messageListEl.value.scrollHeight;
  }
}, { deep: true });

// 模拟渲染不同类型的消息
const renderMessage = (msg) => {
  if (msg.type === 'file') return `[File: ${msg.fileName}]`;
  if (msg.type === 'voice') return `[Voice Message: ${msg.duration}s]`;
  return msg.text;
};

</script>

<template>
  <div class="chat-window" ref="messageListEl">
    <div
        v-for="(msg, index) in messages"
        :key="index"
        class="message-wrapper"
        :class="msg.sender === 'me' ? 'sent' : 'received'">
      <div class="message">
        {{ renderMessage(msg) }}
      </div>
    </div>
    <div v-if="messages.length === 0" class="no-messages">
      Start the conversation!
    </div>
  </div>
</template>

<style scoped>
.chat-window {
  flex-grow: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #e5ddd5; /* WhatsApp-like background */
}
.message-wrapper {
  display: flex;
  margin-bottom: 10px;
}
.message {
  padding: 8px 12px;
  border-radius: 8px;
  max-width: 70%;
  word-wrap: break-word;
}
.message-wrapper.sent { justify-content: flex-end; }
.message-wrapper.sent .message {
  background-color: var(--sent-bg);
}
.message-wrapper.received { justify-content: flex-start; }
.message-wrapper.received .message {
  background-color: var(--received-bg);
  box-shadow: 0 1px 1px rgba(0,0,0,0.05);
}
.no-messages { text-align: center; color: var(--text-secondary); }
</style>