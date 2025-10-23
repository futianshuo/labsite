<template>
  <section>
    <div class="about-toolbar" v-if="me && isEditor(me.role)">
      <el-button type="primary" @click="openEditAbout">编辑概况</el-button>
      <el-button @click="openEditPics">管理图片</el-button>
    </div>

    <el-card shadow="never" class="about-text">
      <div class="about-typo" v-html="formatAbout(aboutText)"></div>
    </el-card>

    <el-row :gutter="16" class="about-pics">
      <el-col :xs="24" :sm="12" :md="8" v-for="p in aboutPics" :key="p.id">
        <el-card shadow="hover" class="about-pic-card">
          <el-image :src="p.image_url" fit="cover" class="about-pic" />
          <div class="about-pic-title" v-if="p.title">{{ p.title }}</div>
          <div class="about-pic-actions" v-if="me && isEditor(me.role)">
            <el-popconfirm title="删除这张图片？" confirm-button-text="删除" cancel-button-text="取消"
                           @confirm="removeAboutPic(p.id)">
              <template #reference>
                <el-button type="danger" text>删除</el-button>
              </template>
            </el-popconfirm>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 编辑正文 -->
    <el-dialog v-model="showEditAbout" title="编辑课题组概况" width="720px">
      <el-input v-model="aboutDraft" type="textarea" :rows="12" placeholder="支持换行；可粘贴纯文本" />
      <template #footer>
        <el-button @click="showEditAbout=false">取消</el-button>
        <el-button type="primary" @click="saveAbout">保存</el-button>
      </template>
    </el-dialog>

    <!-- 管理图片 -->
    <el-dialog v-model="showEditPics" title="管理概况图片（最多 3 张）" width="720px">
      <div class="pics-grid">
        <div v-for="p in aboutPics" :key="p.id" class="pics-item">
          <el-image :src="p.image_url" fit="cover" style="width:180px;height:120px;border-radius:8px;" />
          <el-input v-model="p.title" placeholder="可选：标题/说明" style="margin-top:6px" />
          <el-input-number v-model="p.sort" :min="0" :max="99" style="margin-top:6px" />
        </div>

        <div class="pics-item">
          <el-upload :auto-upload="false" :show-file-list="false" accept="image/*" :on-change="onAboutPicSelect">
            <el-button type="primary" :disabled="aboutPics.length>=3">添加图片</el-button>
          </el-upload>
          <div class="hint">建议 3:2，≤ 1MB，最多 3 张</div>
        </div>
      </div>
      <template #footer>
        <el-button @click="showEditPics=false">关闭</el-button>
        <el-button type="primary" @click="saveAboutPicsMeta">保存修改</el-button>
      </template>
    </el-dialog>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { UploadFile } from 'element-plus'
import api from '../api'

type Role = 'admin'|'editor'|'viewer'
const me = ref<{ id:number; email:string; role:Role } | null>(null)
const isEditor = (r?:string) => r==='admin' || r==='editor'

type AboutPic = { id:number; image_url:string; title?:string; sort:number }
type AboutResp = { body:string; images: AboutPic[] }

const aboutText = ref('')
const aboutDraft = ref('')
const aboutPics = ref<AboutPic[]>([])
const showEditAbout = ref(false)
const showEditPics  = ref(false)

function formatAbout(t:string){ return (t||'').replace(/\n/g,'<br>') }
function openEditAbout(){ aboutDraft.value = aboutText.value; showEditAbout.value = true }
async function saveAbout(){
  try { await api.post('/about', { body: aboutDraft.value }); aboutText.value = aboutDraft.value; ElMessage.success('已保存概况正文'); showEditAbout.value = false }
  catch (e:any) { ElMessage.error(e?.response?.data?.error || e?.message || '保存失败') }
}
function openEditPics(){ showEditPics.value = true }

async function onAboutPicSelect(file: UploadFile){
  const raw = file.raw as File | undefined
  if(!raw) return
  if(aboutPics.value.length >= 3){ ElMessage.warning('最多 3 张'); return }
  if(raw.size > 1024*1024){ ElMessage.warning('请上传 ≤1MB 图片'); return }
  const rd = new FileReader()
  rd.onload = async () => {
    try{
      const payload = { image_url: rd.result as string, title: '', sort: aboutPics.value.length }
      const { data } = await api.post('/about_images', payload)
      aboutPics.value.push({ id: data.id, image_url: data.image_url, title: data.title||'', sort: data.sort??0 })
      ElMessage.success('已添加')
    }catch(e:any){ ElMessage.error(e?.response?.data?.error || e?.message || '上传失败') }
  }
  rd.readAsDataURL(raw)
}
async function removeAboutPic(id:number){
  try{ await api.delete(`/about_images/${id}`); aboutPics.value = aboutPics.value.filter(x=>x.id!==id); ElMessage.success('已删除') }
  catch(e:any){ ElMessage.error(e?.response?.data?.error || e?.message || '删除失败') }
}
async function saveAboutPicsMeta(){
  try{
    await Promise.all(aboutPics.value.map(p => api.patch(`/about_images/${p.id}`, { title: p.title??'', sort: p.sort??0 })))
    ElMessage.success('图片信息已更新')
  }catch(e:any){ ElMessage.error(e?.response?.data?.error || e?.message || '保存失败') }
}

async function fetchAboutAll(){
  const { data } = await api.get<AboutResp>('/about')
  aboutText.value = data.body || ''
  aboutPics.value = (data.images||[]).map(x => ({ id:x.id, image_url:x.image_url, title:x.title||'', sort:x.sort??0 }))
}
async function fetchMe(){ try{ const {data}=await api.get('/auth/me'); me.value=data } catch{ me.value=null } }

onMounted(async () => { await fetchMe(); await fetchAboutAll() })
</script>

<style scoped>
.about-toolbar{ margin-bottom:10px; display:flex; gap:10px; align-items:center; }
.about-text{ padding:16px 20px; }
.about-typo{ font-size:16px; line-height:1.9; color:#2b2b2b; }
.about-pics{ margin-top:16px; }
.about-pic-card{ overflow:hidden; border-radius:12px; }
.about-pic{ width:100%; height:180px; display:block; border-radius:8px; object-fit:cover; }
.about-pic-title{ margin-top:8px; color:#555; }
.pics-grid{ display:flex; flex-wrap:wrap; gap:16px; }
.pics-item{ display:flex; flex-direction:column; align-items:flex-start; }
.hint{ color:#888; font-size:12px; }
</style>
