<script setup>
import { useAuthStore } from './stores/auth';
import ChatView from './views/ChatView.vue';
import LoginView from './views/LoginView.vue';

const authStore = useAuthStore();
</script>

<template>
  <div id="app-container">
    <ChatView v-if="authStore.isAuthenticated" />
    <LoginView v-else />
  </div>
</template>

<style>
/* 全局样式可以放在 main.css 中，这里为了方便演示放在一起 */
#app-container {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}
</style>