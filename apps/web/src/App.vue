<template>
  <div class="layout">
    <!-- 顶部横幅 -->
    <header class="hero">
      <div class="hero-mask">
        <div class="hero-title">
          <h1>课题组主页</h1>
          <p>Correlated Oxide Thin Film & Interface Group</p>
        </div>
      </div>
    </header>

    <!-- 居中导航 -->
    <nav class="topnav">
      <el-menu
        :default-active="$route.path"
        mode="horizontal"
        router
        :ellipsis="false"
        class="menu"
      >
        <el-menu-item index="/">主页</el-menu-item>
        <el-menu-item index="/members">成员</el-menu-item>
        <el-menu-item index="/results">研究成果</el-menu-item>
        <el-menu-item index="/activities">组内活动</el-menu-item>
        <el-menu-item index="/contact">联系方式</el-menu-item>
        <el-menu-item index="/cloud">云计算</el-menu-item>
      </el-menu>
    </nav>

    <!-- 页面主体 -->
    <main class="page-main">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()

// 1) 键盘暗门：Ctrl + Alt + L
function secretGo(e: KeyboardEvent) {
  // 注意要在页面获得焦点的情况下按（点一下页面空白处再按）
  if (e.shiftKey && e.altKey && e.key.toLowerCase() === 'l') {
    ElMessage.success('进入管理登录…')
    router.push('/_/login')
  }
}

onMounted(() => {
  window.addEventListener('keydown', secretGo)
})
onUnmounted(() => {
  window.removeEventListener('keydown', secretGo)
})
</script>


<style scoped>
/* 外层布局：竖排，铺满屏高 */
.layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

/* 横幅：用 public/banner.jpg，或替换为你的图片路径 */
.hero {
  width: 100%;
  height: 220px;                 /* 需要更高就调这里 */
  background-image: url('/banner.jpg');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  border-bottom: 1px solid #e9e9e9;
}

/* 遮罩让标题更清晰，可按需减淡 */
.hero-mask {
  width: 100%;
  height: 100%;
  background: linear-gradient(180deg, rgba(255,255,255,0.12), rgba(255,255,255,0.92));
  display: flex;
  align-items: flex-end;
}

.hero-title {
  max-width: 1100px;             /* 控制左右留白 */
  margin: 0 auto;
  width: 100%;
  padding: 20px 24px;
}
.hero-title h1 {
  margin: 0;
  font-size: 32px;
  font-weight: 800;
}
.hero-title p {
  margin: 6px 0 0;
  color: #444;
}

/* 居中导航 + 左右空隙 */
.topnav {
  border-bottom: 1px solid #eee;
  background: #fff;
}
.menu {
  max-width: 1100px;             /* 同样控制居中容器宽度 */
  margin: 0 auto;
  border-bottom: none;           /* 去掉菜单下划线 */
  display: flex;
  justify-content: center;       /* 菜单项整体居中 */
  gap: 28px;                     /* 菜单项间距 */
}

/* 主体内容与容器对齐 */
.page-main {
  max-width: 1100px;
  margin: 0 auto;
  width: 100%;
  padding: 24px 24px 40px;
}
</style>
