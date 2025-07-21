<script setup>
import { ref } from 'vue';
import { useChatStore } from '@/stores/chat';
import { Phone, Video, UserPlus, Info } from 'lucide-vue-next';
import CallModal from './modals/CallModal.vue';

const chatStore = useChatStore();
const showCallModal = ref(false);
const callType = ref('voice');

const mockAction = (action) => {
  alert(`${action} is a mock action for UI demonstration.`);
};

const startCall = (type) => {
  callType.value = type;
  showCallModal.value = true;
};
</script>

<template>
  <header class="chat-header">
    <div class="chat-info">
      <img :src="chatStore.activeChat.avatar" :alt="chatStore.activeChat.name" class="avatar">
      <span class="chat-name">{{ chatStore.activeChat.name }}</span>
    </div>
    <div class="chat-actions">
      <button @click="startCall('voice')" class="icon-button" title="Start Voice Call">
        <Phone :size="20" />
      </button>
      <button @click="startCall('video')" class="icon-button" title="Start Video Call">
        <Video :size="20" />
      </button>

      <button
          v-if="chatStore.activeChat.type === 'group'"
          @click="mockAction('Add member to group')"
          class="icon-button"
          title="Add Member">
        <UserPlus :size="20" />
      </button>
      <button @click="mockAction('Show chat info')" class="icon-button" title="Chat Info">
        <Info :size="20" />
      </button>
    </div>

    <teleport to="body">
      <CallModal
          v-if="showCallModal"
          :type="callType"
          :contact="chatStore.activeChat"
          @close="showCallModal = false" />
    </teleport>
  </header>
</template>

<style scoped>
.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 20px;
  border-bottom: 1px solid var(--border-color);
  background-color: var(--secondary-bg);
  flex-shrink: 0;
}
.chat-info { display: flex; align-items: center; gap: 15px; }
.chat-name { font-weight: 600; }
.chat-actions { display: flex; gap: 10px; }
</style>