import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './index.css'

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as Icons from '@element-plus/icons-vue'

const app = createApp(App)
app.use(router)
app.use(ElementPlus)
// 可选：全局注册图标
Object.entries(Icons).forEach(([name, comp]) => app.component(name, comp as any))

app.mount('#app')
