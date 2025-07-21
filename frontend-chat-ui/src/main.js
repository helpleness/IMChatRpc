import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import './assets/main.css'

const app = createApp(App)

// 使用 Pinia
app.use(createPinia())

app.mount('#app')