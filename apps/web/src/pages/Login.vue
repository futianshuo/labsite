<template>
  <main style="max-width:420px;margin:60px auto;">
    <h2 style="font:700 22px/1.2 system-ui;margin-bottom:16px;">登录</h2>
    <el-form :model="form" label-width="72px">
      <el-form-item label="邮箱"><el-input v-model="form.email" /></el-form-item>
      <el-form-item label="密码"><el-input v-model="form.password" show-password /></el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="loading" @click="login">登录</el-button>
      </el-form-item>
    </el-form>
  </main>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import api from '../api'

const router = useRouter()
const form = ref({ email: 'admin@local', password: 'admin123' })
const loading = ref(false)

async function login() {
  loading.value = true
  try {
    await api.post('/auth/login', form.value)
    ElMessage.success('登录成功')
    router.push('/members') // 回到成员页
  } catch (e:any) {
    ElMessage.error(e?.response?.data?.error || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>
