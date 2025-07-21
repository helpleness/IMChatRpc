<script setup>
import { ref } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useChatStore } from '@/stores/chat';
import { User, Users, PlusCircle, LogOut } from 'lucide-vue-next';
import AddFriendModal from './modals/AddFriendModal.vue';
import CreateGroupModal from './modals/CreateGroupModal.vue';

const authStore = useAuthStore();
const chatStore = useChatStore();

const showAddFriendModal = ref(false);
const showCreateGroupModal = ref(false);
</script>

<template>
  <aside class="sidebar">
    <div class="sidebar-header">
      <img :src="authStore.user.avatar" :alt="authStore.user.name" class="avatar">
      <span class="username">{{ authStore.user.name }}</span>
      <button @click="authStore.logout()" class="icon-button logout-button" title="Logout">
        <LogOut :size="20" />
      </button>
    </div>

    <div class="sidebar-content">
      <div class="chat-list">
        <div class="list-header">
          <h2>Groups</h2>
          <button @click="showCreateGroupModal = true" class="icon-button" title="Create Group">
            <PlusCircle :size="18"/>
          </button>
        </div>
        <ul>
          <li
              v-for="group in chatStore.groups"
              :key="group.id"
              @click="chatStore.setActiveChat(group)"
              :class="{ active: chatStore.activeChatId === group.id }">
            <Users :size="20" />
            <span>{{ group.name }}</span>
          </li>
        </ul>
      </div>

      <div class="chat-list">
        <div class="list-header">
          <h2>Friends</h2>
          <button @click="showAddFriendModal = true" class="icon-button" title="Add Friend">
            <PlusCircle :size="18"/>
          </button>
        </div>
        <ul>
          <li
              v-for="friend in chatStore.friends"
              :key="friend.id"
              @click="chatStore.setActiveChat(friend)"
              :class="{ active: chatStore.activeChatId === friend.id }">
            <User :size="20" />
            <span>{{ friend.name }}</span>
          </li>
        </ul>
      </div>
    </div>

    <teleport to="body">
      <AddFriendModal v-if="showAddFriendModal" @close="showAddFriendModal = false" />
      <CreateGroupModal v-if="showCreateGroupModal" @close="showCreateGroupModal = false" />
    </teleport>
  </aside>
</template>

<style scoped>
.sidebar {
  width: 350px;
  flex-shrink: 0;
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  background-color: var(--secondary-bg);
}
.sidebar-header {
  padding: 10px 16px;
  display: flex;
  align-items: center;
  gap: 15px;
  border-bottom: 1px solid var(--border-color);
}
.username { font-weight: 600; flex-grow: 1; }
.logout-button { color: #dc3545; }
.sidebar-content { flex-grow: 1; overflow-y: auto; }
.chat-list .list-header {
  padding: 15px 16px 5px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.list-header h2 { font-size: 1em; color: var(--text-secondary); margin: 0; }
.chat-list ul { list-style: none; padding: 0; margin: 0; }
.chat-list li {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background-color 0.2s;
}
.chat-list li:hover { background-color: #f0f0f0; }
.chat-list li.active { background-color: var(--border-color); }
</style>