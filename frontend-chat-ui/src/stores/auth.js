import { defineStore } from 'pinia'
import { ref } from 'vue'

// 使用 defineStore 创建一个 store
export const useAuthStore = defineStore('auth', () => {
    // state
    const user = ref(null) // null表示未登录, 对象表示已登录的用户信息
    const isAuthenticated = ref(false)

    // actions
    function login(username, password) {
        // --- 模拟登录 ---
        // 在真实应用中，这里会发送 API 请求到后端进行验证
        if (username && password) {
            console.log(`Attempting to log in with ${username}`);
            user.value = {
                id: `u_${Date.now()}`,
                name: username,
                avatar: `https://api.dicebear.com/8.x/initials/svg?seed=${username}`
            };
            isAuthenticated.value = true;
            return true;
        }
        return false;
    }

    function register(username, password) {
        // --- 模拟注册 ---
        console.log(`Registering new user: ${username}`);
        // 注册成功后直接登录
        return login(username, password);
    }

    function logout() {
        console.log('Logging out.');
        user.value = null;
        isAuthenticated.value = false;
    }

    return { user, isAuthenticated, login, register, logout }
})