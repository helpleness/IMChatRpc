<script setup>
import { ref } from 'vue';
import { useChatStore } from '@/stores/chat';
import { UserPlus } from 'lucide-vue-next';

// 定义组件可以发出的事件，这里是 'close' 事件
const emit = defineEmits(['close']);

const chatStore = useChatStore();
const friendName = ref('');
const error = ref('');

const handleAddFriend = () => {
  if (!friendName.value.trim()) {
    error.value = 'Friend name cannot be empty.';
    return;
  }

  // 在真实应用中，这里会调用 API 检查用户是否存在并发送好友请求
  // 我们这里只是模拟一个纯前端的添加操作
  console.log(`Mock Action: Sending friend request to "${friendName.value}"`);

  // 调用 store 中的 action 来更新状态
  chatStore.addFriend(friendName.value.trim());

  // 关闭模态框
  emit('close');
};
</script>

<template>
  <div class="modal-backdrop" @click.self="emit('close')">
    <div class="modal-content">
      <h2>Add New Friend</h2>
      <p class="modal-description">Enter the username of the person you want to add.</p>

      <form @submit.prevent="handleAddFriend">
        <div class="input-group">
          <UserPlus :size="20" class="input-icon" />
          <input
              type="text"
              v-model="friendName"
              placeholder="Username"
              required
              autofocus
          />
        </div>

        <p v-if="error" class="error-text">{{ error }}</p>

        <div class="modal-actions">
          <button type="button" class="btn-secondary" @click="emit('close')">Cancel</button>
          <button type="submit" class="btn-primary">Add Friend</button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
/* 使用通用的模态框样式 */
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