<template>
  <main class="hl-wrap">
    <!-- 管理工具栏（只有 admin/editor 可见） -->
    <div class="toolbar" v-if="me && isEditor(me.role)">
      <el-button type="primary" @click="openCreate">新增亮点</el-button>
      <el-button @click="refresh">刷新</el-button>
    </div>

    <!-- 列表 -->
    <el-empty v-if="!loading && highlights.length===0" description="暂无研究亮点" />
    <div v-else class="hl-list">
      <div
        v-for="h in highlights"
        :key="h.id"
        class="hl-card"
        :style="{'--img-w': (h.image_url ? (h.split || 30) : 0) + '%'}"
      >
        <!-- 左：文字 -->
        <div class="hl-left">
          <h3 class="hl-title">{{ h.title }}</h3>
          <div class="hl-summary">{{ h.summary }}</div>

          <div class="hl-link" v-if="h.paper_link">
            <el-link :href="normalizeLink(h.paper_link)" target="_blank" rel="noopener">
              {{ normalizeLink(h.paper_link) }}
            </el-link>
          </div>

          <!-- 行内管理（仅管理员可见） -->
          <div class="hl-admin-row" v-if="me && isEditor(me.role)">
            <el-input-number
              v-model="h.sort"
              :min="0"
              :max="999"
              size="small"
              @change="(v:number)=>saveSort(h, v)"
              :step="1"
              :precision="0"
              :controls="true"
              placeholder="排序"
            />
            <el-popconfirm title="确定删除这条亮点？" confirm-button-text="删除" cancel-button-text="取消"
                           @confirm="remove(h.id)">
              <template #reference>
                <el-button size="small" type="danger" plain>删除</el-button>
              </template>
            </el-popconfirm>
          </div>
        </div>

        <!-- 右：图片 -->
        <div class="hl-right" v-if="h.image_url">
          <img :src="h.image_url" alt="highlight" />
        </div>

          <!-- 滑块悬浮层 -->
        <div class="hl-split" v-if="me && isEditor(me.role)">
          <div class="ratio">{{ 100 - (h.split || 30) }}% / {{ h.split || 30 }}%</div>
          <el-slider
            v-model="h.split"
            :min="15"
            :max="60"
            :step="1"
            :format-tooltip="fmtRatio"
            @change="(val:number)=>saveSplit(h, val)"
            style="width:200px"
          />
        </div>
      </div>
    </div>

    <!-- ============== 新建弹窗 ============== -->
    <el-dialog v-model="showCreate" title="新增研究亮点" width="760px">
      <el-form :model="form" label-width="92px">
        <el-form-item label="标题">
          <el-input v-model="form.title" maxlength="120" show-word-limit />
        </el-form-item>

        <el-form-item label="摘要">
          <el-input v-model="form.summary" type="textarea" :rows="5" maxlength="1000" show-word-limit />
        </el-form-item>

        <el-form-item label="论文链接">
          <el-input v-model="form.paper_link" placeholder="https://doi.org/..." />
        </el-form-item>

        <el-form-item label="配图">
          <div class="upload-line">
            <el-image
              v-if="form.image_url"
              :src="form.image_url"
              fit="cover"
              style="width:240px;height:150px;border-radius:8px"
            />
            <el-upload
              :auto-upload="false"
              :show-file-list="false"
              accept="image/*"
              :on-change="onImgSelected"
            >
              <el-button>选择图片</el-button>
            </el-upload>
            <span class="hint">建议 4:3 或 16:10，≤ 1MB；可不选图</span>
          </div>
        </el-form-item>

        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" :max="999" />
          <span class="hint" style="margin-left:8px">数值越小越靠前</span>
        </el-form-item>

        <el-form-item label="左右比例">
          <el-slider
            v-model="form.split"
            :min="15" :max="60" :step="1"
            :format-tooltip="fmtRatio"
            style="width:260px"
          />
          <span class="hint" style="margin-left:8px">
            文字 / 图片 = {{ 100-(form.split||30) }}% / {{ form.split||30 }}%
          </span>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showCreate=false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="create">创建</el-button>
      </template>
    </el-dialog>
  </main>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { UploadFile } from 'element-plus'
import api from '../api'

type Role = 'admin'|'editor'|'viewer'
const me = ref<{ id:number; email:string; role:Role } | null>(null)
const isEditor = (r?:string) => r==='admin' || r==='editor'

type Highlight = {
  id: number
  title: string
  summary: string
  paper_link?: string
  image_url?: string
  sort: number
  split: number        // ← 后端持久化的比例（图片宽度百分比）
  // 后端还有 published_at，可忽略展示
}

const highlights = ref<Highlight[]>([])
const loading = ref(false)

/* 拉取数据 */
async function fetchMe(){ try{ const {data}=await api.get('/auth/me'); me.value=data } catch{ me.value=null } }
async function fetchList() {
  loading.value = true
  try {
    const { data } = await api.get<Highlight[]>('/highlights')
    // 没有 split 的老数据兜底成 30
    highlights.value = (data || []).map(h => ({ ...h, split: h.split || 30 }))
  } finally { loading.value = false }
}
function refresh(){ fetchList() }

/* 工具 */
function normalizeLink(url: string){
  if(!url) return ''
  if(/^https?:\/\//i.test(url)) return url
  return 'http://' + url
}
function fmtRatio(v: number){ return `${100 - v}% / ${v}%` }

/* 删除 / 保存排序 / 保存比例 */
async function remove(id:number){
  try{
    await api.delete(`/highlights/${id}`)
    highlights.value = highlights.value.filter(x => x.id !== id)
    ElMessage.success('已删除')
  }catch(e:any){ ElMessage.error(e?.response?.data?.error || e?.message || '删除失败') }
}
async function saveSort(h:Highlight, v:number){
  try{
    await api.patch(`/highlights/${h.id}`, { sort: v })
    ElMessage.success('排序已保存')
  }catch(e:any){ ElMessage.error(e?.response?.data?.error || e?.message || '保存失败') }
}
async function saveSplit(h:Highlight, v:number){
  try{
    await api.patch(`/highlights/${h.id}`, { split: v })
    ElMessage.success('比例已保存')
  }catch(e:any){ ElMessage.error(e?.response?.data?.error || e?.message || '保存失败') }
}

/* 新建 */
const showCreate = ref(false)
const submitting = ref(false)
const form = ref<Partial<Highlight>>({
  title: '', summary: '', paper_link: '', image_url: '', sort: 0, split: 30,
})
function openCreate(){
  if(!me.value || !isEditor(me.value.role)) return
  form.value = { title:'', summary:'', paper_link:'', image_url:'', sort:0, split:30 }
  showCreate.value = true
}
function onImgSelected(file: UploadFile){
  const raw = file.raw as File|undefined
  if(!raw) return
  if(raw.size > 1024*1024){ ElMessage.warning('请上传 ≤1MB 图片'); return }
  const rd = new FileReader()
  rd.onload = () => form.value.image_url = rd.result as string
  rd.readAsDataURL(raw)
}
async function create(){
  if(!form.value.title) return ElMessage.warning('请填写标题')
  submitting.value = true
  try{
    const payload = {
      title: form.value.title,
      summary: form.value.summary || '',
      paper_link: form.value.paper_link || '',
      image_url: form.value.image_url || '',
      sort: form.value.sort ?? 0,
      split: form.value.split ?? 30,
      published_at: '',   // 可留空，后端会兜底
    }
    await api.post('/highlights', payload)
    ElMessage.success('已创建')
    showCreate.value = false
    await fetchList()
  }catch(e:any){
    ElMessage.error(e?.response?.data?.error || e?.message || '创建失败')
  }finally{ submitting.value = false }
}

onMounted(async ()=>{ await fetchMe(); await fetchList() })
</script>

<style scoped>
.hl-wrap{ max-width:1100px; margin:24px auto; padding:0 16px; }
.toolbar{ margin-bottom:12px; display:flex; gap:10px; align-items:center; }

.hl-list{ display:flex; flex-direction:column; gap:18px; }

/* 卡片布局：用 CSS Grid 控制左右宽度，右列宽度由 --img-w 决定 */
.hl-card{
  --img-w: 30%;
  display:grid;
  grid-template-columns: calc(100% - var(--img-w)) var(--img-w);
  gap:18px;
  background:#fff;
  border-radius:14px;
  padding:18px;
  box-shadow: 0 2px 10px rgba(0,0,0,.05);
  align-items: stretch; /* 两列等高 */
  position: relative;
}

/* 左列（文字） */
.hl-left{ min-width:0; }
.hl-title{ margin:2px 0 10px; font-size:22px; line-height:1.35; font-weight:800; color:#1b2559; }
.hl-summary{ color:#2c2c2c; line-height:1.85; white-space: pre-line; overflow-wrap:anywhere; }
.hl-link{ margin-top:10px; }

/* 右列（图片） */
.hl-right{ position:relative; border-radius:10px; overflow:hidden; background:#f6f7f9; }
.hl-right img{ width:100%; height:100%; object-fit:cover; display:block; }

/* 比例滑块浮层 */
.hl-split {
  position: absolute;
  right: 16px;
  top: -14px; /* 向上浮出图片 */
  z-index: 3;
  background: rgba(255,255,255,0.95);
  backdrop-filter: blur(4px);
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 6px 10px;
  display: flex;
  align-items: center;
  gap: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,.08);
}
.ratio{ font-size:12px; color:#666; }

/* 行内管理区域 */
.hl-admin-row{ margin-top:12px; display:flex; gap:10px; align-items:center; }

/* 上传区 */
.upload-line{ display:flex; align-items:center; gap:12px; }
.hint{ color:#888; font-size:12px; }

/* 小屏改为纵向 */
@media (max-width: 900px){
  .hl-card{
    grid-template-columns: 1fr;
  }
  .hl-right{ height:220px; } /* 给图片一个合理高度 */
}
</style>
