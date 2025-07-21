<script setup>
import { ref } from 'vue';
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();
const isRegistering = ref(false);
const username = ref('');
const password = ref('');
const error = ref('');

const handleSubmit = () => {
  error.value = '';
  let success = false;
  if (isRegistering.value) {
    success = authStore.register(username.value, password.value);
  } else {
    success = authStore.login(username.value, password.value);
  }
  if (!success) {
    error.value = 'Invalid username or password (mock).';
  }
};
</script>

<template>
  <div class="login-container">
    <div class="login-box">
      <h1>{{ isRegistering ? 'Register' : 'Welcome Back!' }}</h1>
      <p>{{ isRegistering ? 'Create a new account' : 'Sign in to continue' }}</p>
      <form @submit.prevent="handleSubmit">
        <input type="text" v-model="username" placeholder="Username" required>
        <input type="password" v-model="password" placeholder="Password" required>
        <p v-if="error" class="error-text">{{ error }}</p>
        <button type="submit">{{ isRegistering ? 'Register' : 'Login' }}</button>
      </form>
      <a href="#" @click.prevent="isRegistering = !isRegistering">
        {{ isRegistering ? 'Already have an account? Login' : 'Don\'t have an account? Register' }}
      </a>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: var(--bg-color);
}
.login-box {
  background: var(--primary-bg);
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  width: 100%;
  max-width: 400px;
  text-align: center;
}
h1 { font-size: 2em; margin-bottom: 10px; }
p { color: var(--text-secondary); margin-bottom: 30px; }
form { display: flex; flex-direction: column; gap: 15px; }
input {
  padding: 15px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 1em;
}
button {
  padding: 15px;
  background-color: var(--accent-blue);
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 1em;
  font-weight: bold;
  cursor: pointer;
  transition: background-color 0.2s;
}
button:hover { background-color: #0056b3; }
a {
  margin-top: 20px;
  display: inline-block;
  color: var(--accent-blue);
  text-decoration: none;
}
.error-text {
  color: #dc3545;
  margin-bottom: 0;
  margin-top: -5px;
}
</style>