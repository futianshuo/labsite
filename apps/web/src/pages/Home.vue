<template>
  <main class="home-wrap">
    <el-tabs :model-value="active" @tab-change="goto" class="home-tabs" stretch>
      <el-tab-pane label="课题组动态" name="news" />
      <el-tab-pane label="课题组概况" name="about" />
      <el-tab-pane label="研究亮点" name="highlights" />
    </el-tabs>

    <router-view />
  </main>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const active = computed(() => route.path.split('/')[1] || 'news')
const goto = (name: string) => router.push('/' + name)
</script>

<style scoped>
.home-wrap{ max-width:1100px; margin:24px auto; padding:0 16px; }
.home-tabs :deep(.el-tabs__item){ font-size:16px; }
</style>
