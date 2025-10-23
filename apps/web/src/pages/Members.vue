<template>
  <main class="container" style="max-width: 1100px; margin: 24px auto;">
    <el-page-header content="团队成员" icon="">
      <template #title>
        <span style="font-weight:700;font-size:20px;">课题组主页</span>
      </template>
    </el-page-header>

    <el-divider />

    <!-- 工具栏：只在已登录管理员/编辑时显示新增与退出 -->
    <div style="display:flex; gap:8px; margin-bottom:16px;">
      <el-button
        v-if="me && isEditor(me.role)"
        type="primary"
        @click="openAdd"
      >
        新增成员
      </el-button>
      <el-button
        v-if="me"
        type="warning"
        plain
        @click="logout"
      >
        退出管理员
      </el-button>
      <!-- 按要求移除了“刷新”按钮；新增/删除后会自动刷新列表 -->
    </div>

    <!-- —— 按角色分组 —— -->
    <!-- 1) 课题组负责人（保留：职称、研究方向） -->
    <section v-if="piList.length">
      <el-divider content-position="left">课题组负责人（PI）</el-divider>
      <el-row :gutter="16">
        <el-col :xs="24" :md="12" v-for="m in piList" :key="'pi-'+m.id" style="margin-bottom:16px;">
          <el-card class="member-card" :body-style="{padding:'16px 18px'}" shadow="hover">
            <div class="member">
              <el-avatar :size="64" :src="m.avatar_url || ''">{{ m.name?.[0] || 'P' }}</el-avatar>
              <div class="meta">
                <div class="name">{{ m.name }}</div>
                <div class="sub">{{ m.title }}</div>
                <div class="sub" v-if="m.research_areas">研究方向：{{ m.research_areas }}</div>
              </div>

              <el-popconfirm
                v-if="me && isEditor(me.role)"
                title="确定删除该成员？"
                confirm-button-text="删除"
                cancel-button-text="取消"
                @confirm="remove(m.id)"
              >
                <template #reference>
                  <el-button type="danger" plain size="small">删除</el-button>
                </template>
              </el-popconfirm>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </section>

    <!-- 2) 硕士研究生（保留：研究方向） -->
    <section v-if="msList.length">
      <el-divider content-position="left">硕士研究生（MS）</el-divider>
      <el-row :gutter="16">
        <el-col :xs="24" :md="12" v-for="m in msList" :key="'ms-'+m.id" style="margin-bottom:16px;">
          <el-card class="member-card" :body-style="{padding:'16px 18px'}" shadow="hover">
            <div class="member">
              <el-avatar :size="64" :src="m.avatar_url || ''">{{ m.name?.[0] || 'M' }}</el-avatar>
              <div class="meta">
                <div class="name">{{ m.name }}</div>
                <div class="sub" v-if="m.research_areas">研究方向：{{ m.research_areas }}</div>
              </div>

              <el-popconfirm
                v-if="me && isEditor(me.role)"
                title="确定删除该成员？"
                confirm-button-text="删除"
                cancel-button-text="取消"
                @confirm="remove(m.id)"
              >
                <template #reference>
                  <el-button type="danger" plain size="small">删除</el-button>
                </template>
              </el-popconfirm>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </section>

    <!-- 3) 本科生（无需研究方向等） -->
    <section v-if="ugList.length">
      <el-divider content-position="left">本科生（UG）</el-divider>
      <el-row :gutter="16">
        <el-col :xs="24" :md="12" v-for="m in ugList" :key="'ug-'+m.id" style="margin-bottom:16px;">
          <el-card class="member-card" :body-style="{padding:'16px 18px'}" shadow="hover">
            <div class="member">
              <el-avatar :size="64" :src="m.avatar_url || ''">{{ m.name?.[0] || 'U' }}</el-avatar>
              <div class="meta">
                <div class="name">{{ m.name }}</div>
              </div>

              <el-popconfirm
                v-if="me && isEditor(me.role)"
                title="确定删除该成员？"
                confirm-button-text="删除"
                cancel-button-text="取消"
                @confirm="remove(m.id)"
              >
                <template #reference>
                  <el-button type="danger" plain size="small">删除</el-button>
                </template>
              </el-popconfirm>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </section>

    <!-- 4) 校友（保留：去向） -->
    <section v-if="alumniList.length">
      <el-divider content-position="left">校友（Alumni）</el-divider>
      <el-row :gutter="16">
        <el-col :xs="24" :md="12" v-for="m in alumniList" :key="'alumni-'+m.id" style="margin-bottom:16px;">
          <el-card class="member-card" :body-style="{padding:'16px 18px'}" shadow="hover">
            <div class="member">
              <el-avatar :size="64" :src="m.avatar_url || ''">{{ m.name?.[0] || 'A' }}</el-avatar>
              <div class="meta">
                <div class="name">{{ m.name }}</div>
                <div class="sub" v-if="m.destination">去向：{{ m.destination }}</div>
              </div>

              <el-popconfirm
                v-if="me && isEditor(me.role)"
                title="确定删除该成员？"
                confirm-button-text="删除"
                cancel-button-text="取消"
                @confirm="remove(m.id)"
              >
                <template #reference>
                  <el-button type="danger" plain size="small">删除</el-button>
                </template>
              </el-popconfirm>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </section>

    <!-- 空态：四组都没有 -->
    <el-empty v-if="!loading && !piList.length && !msList.length && !ugList.length && !alumniList.length"
              description="暂无成员" />

    <!-- 新增成员（根据角色动态显示字段），仅管理员可打开 -->
    <el-dialog v-model="showAdd" title="新增成员" width="560px">
      <el-form :model="form" label-width="96px" style="padding-right:8px;">
        <el-form-item label="姓名">
          <el-input v-model="form.name" placeholder="如：小付" />
        </el-form-item>

        <el-form-item label="角色">
          <el-select v-model="form.role" placeholder="选择角色">
            <el-option label="课题组负责人（PI）" value="PI" />
            <el-option label="硕士研究生（MS）" value="MS" />
            <el-option label="本科生（UG）" value="UG" />
            <el-option label="校友（Alumni）" value="alumni" />
          </el-select>
        </el-form-item>

        <!-- 头像上传：本地转 Base64 存入 avatar_url -->
        <el-form-item label="头像">
          <div style="display:flex; align-items:center; gap:12px;">
            <el-avatar :size="48" :src="form.avatar_url || ''" />
            <el-upload
              :auto-upload="false"
              :show-file-list="false"
              accept="image/*"
              :on-change="handleAvatarSelected"
            >
              <el-button>选择图片</el-button>
            </el-upload>
            <span style="color:#888; font-size:12px;">（建议 ≤ 1MB）</span>
          </div>
        </el-form-item>

        <!-- 仅 PI：职称 -->
        <el-form-item v-if="form.role==='PI'" label="职称/头衔">
          <el-input v-model="form.title" placeholder="如：教授 / 研究员" />
        </el-form-item>

        <!-- PI / MS：研究方向 -->
        <el-form-item v-if="form.role==='PI' || form.role==='MS'" label="研究方向">
          <el-input v-model="form.research_areas" placeholder="逗号分隔，如：3D Vision, Geophysics" />
        </el-form-item>

        <!-- 仅 Alumni：去向 -->
        <el-form-item v-if="form.role==='alumni'" label="去向">
          <el-input v-model="form.destination" placeholder="如：字节-算法工程师 / 某大学-PhD" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showAdd=false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submit">提交</el-button>
      </template>
    </el-dialog>

    <el-backtop :right="20" :bottom="20" />
  </main>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from "vue"
import { ElMessage } from "element-plus"
import type { UploadFile } from "element-plus"
import api from "../api"

interface Member {
  id: number
  name: string
  role: "PI" | "MS" | "UG" | "alumni"
  title?: string
  avatar_url?: string
  research_areas?: string
  destination?: string   // ← 新增字段：校友去向
}

const list = ref<Member[]>([])
const loading = ref(false)
const me = ref<{ id: number; email: string; role: string } | null>(null)

const isEditor = (role?: string) => role === "admin" || role === "editor"

// 分组
const piList     = computed(() => list.value.filter(m => m.role === "PI"))
const msList     = computed(() => list.value.filter(m => m.role === "MS"))
const ugList     = computed(() => list.value.filter(m => m.role === "UG"))
const alumniList = computed(() => list.value.filter(m => m.role === "alumni"))

async function fetchMe() {
  try {
    const { data } = await api.get("/auth/me")
    me.value = data
  } catch { me.value = null }
}

async function fetchList() {
  loading.value = true
  try {
    const { data } = await api.get<Member[]>("/members")
    list.value = data
  } catch (e: any) {
    ElMessage.error(`获取成员失败：${e?.message || e}`)
  } finally {
    loading.value = false
  }
}

// 退出管理员
async function logout() {
  try {
    await api.post("/auth/logout")
    me.value = null
    ElMessage.success("已退出管理员")
  } catch (e:any) {
    ElMessage.error(e?.response?.data?.error || "退出失败")
  }
}

// ==== 新增 ====
const showAdd = ref(false)
const submitting = ref(false)
const form = ref<Partial<Member>>({ role: "MS" })

function openAdd() {
  if (!me.value || !isEditor(me.value.role)) return
  form.value = { role: "MS" }
  showAdd.value = true
}

// 本地选择头像 → 读为 Base64 存到 avatar_url
function handleAvatarSelected(file: UploadFile) {
  const raw = file.raw as File | undefined
  if (!raw) return
  if (raw.size > 1024 * 1024) { // 1MB
    ElMessage.warning("请上传小于 1MB 的图片")
    return
  }
  const reader = new FileReader()
  reader.onload = () => { form.value.avatar_url = reader.result as string }
  reader.readAsDataURL(raw)
}

async function submit() {
  if (!form.value.name) return ElMessage.warning("请填写姓名")
  // 角色必填且限定
  if (!form.value.role || !["PI","MS","UG","alumni"].includes(form.value.role)) {
    return ElMessage.warning("请选择角色")
  }
  submitting.value = true
  try {
    await api.post("/members", form.value)
    ElMessage.success("新增成功")
    showAdd.value = false
    await fetchList()
  } catch (e: any) {
    ElMessage.error(`新增失败：${e?.response?.data?.error || e?.message || e}`)
  } finally {
    submitting.value = false
  }
}

// ==== 删除 ====
async function remove(id: number) {
  try {
    await api.delete(`/members/${id}`)
    ElMessage.success("已删除")
    await fetchList()
  } catch (e: any) {
    ElMessage.error(`删除失败：${e?.response?.data?.error || e?.message || e}`)
  }
}

onMounted(async () => {
  await fetchMe()
  await fetchList()
})
</script>

<style scoped>
.member-card { border-radius: 14px; }
.member {
  display: flex;
  align-items: center;
  gap: 14px;
}
.member .meta { flex: 1; min-width: 0; }
.member .name { font-size: 18px; font-weight: 700; }
.member .sub  { font-size: 13px; color: #666; margin-top: 4px; }
</style>
