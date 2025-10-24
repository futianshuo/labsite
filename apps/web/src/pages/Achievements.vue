<template>
  <main class="ach-wrap">
    <el-page-header content="研究成果" icon="">
      <template #title><span class="hdr">课题组主页</span></template>
    </el-page-header>
    <el-divider />

    <el-tabs v-model="tab" class="ach-tabs" stretch>
      <el-tab-pane label="论文" name="papers" />
      <el-tab-pane label="竞赛奖项" name="awards" />
      <el-tab-pane label="项目" name="projects" />
    </el-tabs>

    <!-- 管理工具 -->
    <div class="toolbar" v-if="me && isEditor(me.role)">
      <el-button v-if="tab==='papers'" type="primary" @click="openCreate('paper')">新增论文</el-button>
      <el-button v-if="tab==='awards'" type="primary" @click="openCreate('award')">新增奖项</el-button>
      <el-button v-if="tab==='projects'" type="primary" @click="openCreate('project')">新增项目</el-button>
    </div>

    <!-- 论文 -->
    <section v-show="tab==='papers'">
      <div class="ach-card" v-for="p in papers" :key="p.id">
        <div class="ach-text">{{ p.text }}</div>

        <div class="ach-meta" v-if="p.link">
          <el-link :href="p.link" target="_blank" rel="noopener">{{ p.link }}</el-link>
        </div>

        <div class="ach-actions" v-if="me && isEditor(me.role)">
          <span class="sort">排序：</span>
          <el-input-number v-model="p.sort" :min="0" :max="999" size="small" />
          <el-button size="small" @click="saveSort('paper', p)">保存</el-button>
          <el-popconfirm title="删除这条论文？" confirm-button-text="删除" cancel-button-text="取消"
                         @confirm="removeItem('paper', p.id)">
            <template #reference>
              <el-button size="small" type="danger" plain>删除</el-button>
            </template>
          </el-popconfirm>
        </div>
      </div>

      <el-empty v-if="!loading && papers.length===0" description="暂无论文" />
    </section>

    <!-- 奖项 -->
    <section v-show="tab==='awards'">
      <div class="ach-card" v-for="a in awards" :key="a.id">
        <div class="ach-text">{{ a.text }}</div>
        <div v-if="a.work" class="ach-sub">作品名：{{ a.work }}</div>

        <div class="ach-actions" v-if="me && isEditor(me.role)">
          <span class="sort">排序：</span>
          <el-input-number v-model="a.sort" :min="0" :max="999" size="small" />
          <el-button size="small" @click="saveSort('award', a)">保存</el-button>
          <el-popconfirm title="删除该奖项？" confirm-button-text="删除" cancel-button-text="取消"
                         @confirm="removeItem('award', a.id)">
            <template #reference>
              <el-button size="small" type="danger" plain>删除</el-button>
            </template>
          </el-popconfirm>
        </div>
      </div>

      <el-empty v-if="!loading && awards.length===0" description="暂无奖项" />
    </section>

    <!-- 项目 -->
    <section v-show="tab==='projects'">
      <div class="ach-card" v-for="pr in projects" :key="pr.id">
        <div class="ach-title">{{ pr.title }}</div>

        <div class="ach-actions" v-if="me && isEditor(me.role)">
          <span class="sort">排序：</span>
          <el-input-number v-model="pr.sort" :min="0" :max="999" size="small" />
          <el-button size="small" @click="saveSort('project', pr)">保存</el-button>
          <el-popconfirm title="删除该项目？" confirm-button-text="删除" cancel-button-text="取消"
                         @confirm="removeItem('project', pr.id)">
            <template #reference>
              <el-button size="small" type="danger" plain>删除</el-button>
            </template>
          </el-popconfirm>
        </div>
      </div>

      <el-empty v-if="!loading && projects.length===0" description="暂无项目" />
    </section>

    <!-- 新建弹窗 -->
    <el-dialog v-model="showCreate" :title="createTitle" width="680px">
      <el-form :model="form" label-width="92px">
        <template v-if="mode==='paper'">
          <el-form-item label="论文描述">
            <el-input v-model="form.text" type="textarea" :rows="5" placeholder="整条论文信息（作者、题目、期刊等），支持换行" />
          </el-form-item>
          <el-form-item label="外部链接">
            <el-input v-model="form.link" placeholder="可留空（论文DOI/ArXiv/出版社）" />
          </el-form-item>
          <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" :max="999" /></el-form-item>
        </template>

        <template v-else-if="mode==='award'">
          <el-form-item label="奖项描述">
            <el-input v-model="form.text" type="textarea" :rows="4" placeholder="奖项名称、主办方、团队/个人等" />
          </el-form-item>
          <el-form-item label="作品名">
            <el-input v-model="form.work" placeholder="可留空" />
          </el-form-item>
          <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" :max="999" /></el-form-item>
        </template>

        <template v-else>
          <el-form-item label="项目标题">
            <el-input v-model="form.title" placeholder="只填标题即可" />
          </el-form-item>
          <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" :max="999" /></el-form-item>
        </template>
      </el-form>

      <template #footer>
        <el-button @click="showCreate=false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitCreate">保存</el-button>
      </template>
    </el-dialog>
  </main>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import api from '../api'

type Role = 'admin'|'editor'|'viewer'
const me = ref<{ id:number; email:string; role:Role } | null>(null)
const isEditor = (r?:string) => r==='admin' || r==='editor'

/** Tabs */
const tab = ref<'papers'|'awards'|'projects'>('papers')
const loading = ref(false)

/** 数据结构 */
type Paper   = { id:number; text:string; link?:string; sort:number }
type Award   = { id:number; text:string; work?:string; sort:number }
type Project = { id:number; title:string; sort:number }

const papers   = ref<Paper[]>([])
const awards   = ref<Award[]>([])
const projects = ref<Project[]>([])

/** 拉数据 */
async function fetchMe(){ try{ const {data}=await api.get('/auth/me'); me.value=data } catch{ me.value=null } }
async function fetchPapers(){ const {data}=await api.get<Paper[]>('/papers'); papers.value=data }
async function fetchAwards(){ const {data}=await api.get<Award[]>('/awards'); awards.value=data }
async function fetchProjects(){ const {data}=await api.get<Project[]>('/projects'); projects.value=data }

onMounted(async ()=>{
  loading.value = true
  try{
    await fetchMe()
    await Promise.all([fetchPapers(), fetchAwards(), fetchProjects()])
  } finally { loading.value = false }
})

/** 删除 & 排序保存 */
async function removeItem(kind:'paper'|'award'|'project', id:number){
  try{
    await api.delete(`/${kind}s/${id}`)
    if(kind==='paper') papers.value = papers.value.filter(x=>x.id!==id)
    if(kind==='award') awards.value = awards.value.filter(x=>x.id!==id)
    if(kind==='project') projects.value = projects.value.filter(x=>x.id!==id)
    ElMessage.success('已删除')
  }catch(e:any){ ElMessage.error(e?.response?.data?.error || e?.message || '删除失败') }
}

async function saveSort(kind:'paper'|'award'|'project', item:any){
  try{
    await api.patch(`/${kind}s/${item.id}`, { sort: item.sort ?? 0 })
    ElMessage.success('排序已保存')
  }catch(e:any){ ElMessage.error(e?.response?.data?.error || e?.message || '保存失败') }
}

/** 新建弹窗 */
type Mode = 'paper'|'award'|'project'
const mode = ref<Mode>('paper')
const showCreate = ref(false)
const submitting = ref(false)

const form = ref<any>({ sort: 0 })
const createTitle = computed(()=> mode.value==='paper' ? '新增论文' : mode.value==='award' ? '新增奖项' : '新增项目')

function openCreate(m:Mode){
  mode.value = m
  form.value = { sort: 0 }
  showCreate.value = true
}

async function submitCreate(){
  submitting.value = true
  try{
    if(mode.value==='paper'){
      if(!form.value.text) return ElMessage.warning('请填写论文描述')
      await api.post('/papers', { text: form.value.text, link: form.value.link || '', sort: form.value.sort||0 })
      await fetchPapers()
    }else if(mode.value==='award'){
      if(!form.value.text) return ElMessage.warning('请填写奖项描述')
      await api.post('/awards', { text: form.value.text, work: form.value.work || '', sort: form.value.sort||0 })
      await fetchAwards()
    }else{
      if(!form.value.title) return ElMessage.warning('请填写项目标题')
      await api.post('/projects', { title: form.value.title, sort: form.value.sort||0 })
      await fetchProjects()
    }
    ElMessage.success('已保存')
    showCreate.value = false
  }catch(e:any){
    ElMessage.error(e?.response?.data?.error || e?.message || '保存失败')
  }finally{ submitting.value = false }
}
</script>

<style scoped>
.ach-wrap{ max-width:1100px; margin:24px auto; padding:0 16px; }
.hdr{ font-weight:700; font-size:20px; }
.ach-tabs :deep(.el-tabs__item){ font-size:16px; }
.toolbar{ margin:8px 0 12px; display:flex; gap:10px; align-items:center; }

.ach-card{
  background:#fff;
  border:1px solid #eee;
  border-radius:12px;
  padding:18px 20px;
  margin-bottom:16px;
  box-shadow:0 1px 4px rgba(0,0,0,.04);
}
.ach-title{ font-size:20px; font-weight:800; line-height:1.5; }
.ach-text{ font-size:15px; line-height:1.9; white-space:pre-line; overflow-wrap:anywhere; }
.ach-sub{ margin-top:6px; color:#666; font-size:13px; }
.ach-meta{ margin-top:8px; }
.ach-actions{ margin-top:10px; display:flex; align-items:center; gap:10px; }
.sort{ color:#666; }
</style>
