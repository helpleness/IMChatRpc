<script setup>
import { ref } from 'vue';
import { Paperclip, Mic, Smile, Send } from 'lucide-vue-next';

const emit = defineEmits(['sendMessage']);
const messageText = ref('');

const handleSubmit = () => {
  if (messageText.value.trim()) {
    emit('sendMessage', messageText.value.trim());
    messageText.value = '';
  }
};

const mockAction = (action) => {
  alert(`${action} is a mock action for UI demonstration.\nIn a real app, this would open a file picker, start recording, or show an emoji keyboard.`);
}
</script>

<template>
  <form class="message-input-form" @submit.prevent="handleSubmit">
    <div class="input-actions">
      <button type="button" @click="mockAction('Send Emoji')" class="icon-button" title="Emoji">
        <Smile :size="22" />
      </button>
      <button type="button" @click="mockAction('Send File')" class="icon-button" title="Attach File">
        <Paperclip :size="22" />
      </button>
    </div>
    <input
        type="text"
        v-model="messageText"
        placeholder="Type a message..."
    />
    <div class="input-actions">
      <button type="button" @click="mockAction('Send Voice Message')" class="icon-button" title="Record Voice">
        <Mic :size="22" />
      </button>
      <button type="submit" class="icon-button send-button" title="Send">
        <Send :size="22" />
      </button>
    </div>
  </form>
</template>

<style scoped>
.message-input-form {
  display: flex;
  align-items: center;
  padding: 10px 20px;
  gap: 10px;
  background-color: var(--secondary-bg);
}
.input-actions { display: flex; align-items: center; gap: 5px; }
.message-input-form input {
  flex-grow: 1;
  border: none;
  background-color: white;
  border-radius: 20px;
  padding: 12px 18px;
  font-size: 1em;
}
.message-input-form input:focus { outline: none; }
.send-button { color: var(--accent-blue); }
</style>