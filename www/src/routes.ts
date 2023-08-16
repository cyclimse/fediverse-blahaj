import Home from './pages/Home.vue'
import About from './pages/About.vue'
import NotFound from './pages/NotFound.vue'
import Server from './pages/Server.vue'

export const routes = [
    {
        path: '/',
        component: Home,
    },
    {
        path: '/servers/:id',
        component: Server,
    },
    {
        path: '/about',
        component: About,
    },
    {
        path: '/:catchAll(.*)',
        component: NotFound,
    }
]