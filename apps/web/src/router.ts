import { createRouter, createWebHistory } from 'vue-router'
import Home from './pages/Home.vue'
import Members from './pages/Members.vue'
import Achievements from './pages/Achievements.vue'
import Activities from './pages/Activities.vue'
import Contact from './pages/Contact.vue'
import Cloud from './pages/Cloud.vue'
import Login from './pages/Login.vue'

export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/',          component: Home      },
    { path: '/members',   component: Members   },
    { path: '/results',   component: Achievements },
    { path: '/activities',component: Activities },
    { path: '/contact',   component: Contact   },
    { path: '/cloud',     component: Cloud     },
    { path: '/_/login',   component: Login },
  ],
  scrollBehavior: () => ({ top: 0 }),
})

