import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './style.css'

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)
app.use(router)

// Init theme store so it applies the class immediately
import { useThemeStore } from './stores/theme'
useThemeStore()

app.mount('#app')
