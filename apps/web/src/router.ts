import { createRouter, createWebHistory } from 'vue-router'

// 首页容器 + 三个子页
import Home from './pages/Home.vue'
import HomeNews from './pages/HomeNews.vue'
import HomeAbout from './pages/HomeAbout.vue'
import HomeHighlights from './pages/HomeHighlights.vue'

// 其他独立页面
import Members from './pages/Members.vue'
import Achievements from './pages/Achievements.vue'
import Activities from './pages/Activities.vue'
import Contact from './pages/Contact.vue'
import Cloud from './pages/Cloud.vue'
import Login from './pages/Login.vue'

export default createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: Home,
      children: [
        // 进入根路径时默认跳到 news
        { path: '', redirect: { name: 'home-news' } },
        { path: 'news',       name: 'home-news',       component: HomeNews },
        { path: 'about',      name: 'home-about',      component: HomeAbout },
        { path: 'highlights', name: 'home-highlights', component: HomeHighlights },
      ],
    },

    // 其它一级路由保持不变
    { path: '/members',    component: Members },
    { path: '/results',    component: Achievements },
    { path: '/activities', component: Activities },
    { path: '/contact',    component: Contact },
    { path: '/cloud',      component: Cloud },
    { path: '/_/login',    component: Login },

    // 兜底：未知路径回到首页
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
  scrollBehavior(_to, _from, saved) {
    return saved || { top: 0 }
  },
})
