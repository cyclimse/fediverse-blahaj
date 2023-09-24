import { RouteLocation } from 'vue-router'

import About from './pages/About.vue'
import NotFound from './pages/NotFound.vue'
import Instance from './pages/Instance.vue'
import Instances from './pages/Instances.vue'

export const routes = [
    {
        path: '/',
        component: Instances,
        props: (route: RouteLocation) => ({ page: route.query.page ? parseInt(route.query.page as string) : 1 }),
    },
    {
        path: '/instances/:id',
        component: Instance,
        props: true,
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