import { createApp } from 'vue'
import './style.scss'
import App from './App.vue'
import { createRouter, createWebHistory } from 'vue-router';
import { routes } from './routes';
import { FedClient } from './api/generated/FedClient';

const router = createRouter({
    history: createWebHistory(),
    routes
})

const api = new FedClient({
    BASE: import.meta.env.VITE_API_URL,
});

const app = createApp(App)
app.config.globalProperties.$api = api.default;

app.use(router)
app.mount('#app')
