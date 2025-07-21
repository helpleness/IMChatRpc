<script setup>
import { ref } from 'vue';
import { useChatStore } from '@/stores/chat';
import { Users } from 'lucide-vue-next';

const emit = defineEmits(['close']);
const chatStore = useChatStore();
const groupName = ref('');
const error = ref('');

const handleCreateGroup = () => {
  if (!groupName.value.trim()) {
    error.value = 'Group name cannot be empty.';
    return;
  }

  // 模拟纯前端的创建群组操作
  console.log(`Mock Action: Creating new group "${groupName.value}"`);
  chatStore.createGroup(groupName.value.trim());
  emit('close');
};
</script>

<template>
  <div class="modal-backdrop" @click.self="emit('close')">
    <div class="modal-content">
      <h2>Create New Group</h2>
      <p class="modal-description">Give your new group a name.</p>

      <form @submit.prevent="handleCreateGroup">
        <div class="input-group">
          <Users :size="20" class="input-icon" />
          <input
              type="text"
              v-model="groupName"
              placeholder="Group Name"
              required
              autofocus
          />
        </div>

        <p v-if="error" class="error-text">{{ error }}</p>

        <div class="modal-actions">
          <button type="button" class="btn-secondary" @click="emit('close')">Cancel</button>
          <button type="submit" class="btn-primary">Create Group</button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
/* 样式与 AddFriendModal 完全相同，可以复用 */
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.6);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  width: 90%;
  max-width: 420px;
  box-shadow: 0 5px 15px rgba(0,0,0,0.3);
}

.modal-content h2 {
  margin-top: 0;
  margin-bottom: 0.5rem;
  color: #111;
}

.modal-description {
  margin-top: 0;
  margin-bottom: 1.5rem;
  color: #667781;
}

.input-group {
  position: relative;
  margin-bottom: 1rem;
}

.input-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #667781;
}

.input-group input {
  width: 100%;
  padding: 12px 12px 12px 40px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  font-size: 1rem;
  box-sizing: border-box;
}

.error-text {
  color: #dc3545;
  font-size: 0.9rem;
  margin-top: -0.5rem;
  margin-bottom: 1rem;
}

.modal-actions {
  margin-top: 1.5rem;
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

button {
  padding: 10px 20px;
  border-radius: 6px;
  border: none;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-primary {
  background-color: var(--accent-blue, #007bff);
  color: white;
}
.btn-primary:hover {
  background-color: #0056b3;
}

.btn-secondary {
  background-color: #e9edef;
  color: #111b21;
}
.btn-secondary:hover {
  background-color: #d1d7db;
}
</style>