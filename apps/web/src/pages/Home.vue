<template>
  <main class="home-wrap">
    <el-tabs v-model="active" class="home-tabs" stretch>
      <el-tab-pane label="课题组动态" name="news" />
      <el-tab-pane label="课题组概况" name="about" />
      <el-tab-pane label="研究亮点" name="highlights" />
    </el-tabs>

    <!-- =================== 课题组动态 =================== -->
    <section v-show="active==='news'">
      <!-- 管理工具（仅管理员/编辑可见） -->
      <div class="toolbar" v-if="me && isEditor(me.role)">
        <el-button type="primary" @click="openBanner">新增轮播图</el-button>
        <el-button @click="showBannerMgr = true">管理轮播图</el-button>
        <el-divider direction="vertical" />
        <el-button type="success" @click="openUploadNewsImage">上传新闻图片</el-button>
        <el-button @click="showNewsImgMgr = true">管理新闻图库</el-button>
        <el-divider direction="vertical" />
        <el-button type="warning" @click="openPublish">发布新闻</el-button>
      </div>

      <!-- 顶部轮播：只展示 Banner 图库（独立系统） -->
      <el-carousel height="380px" trigger="click" indicator-position="outside" v-if="banners.length">
        <el-carousel-item v-for="b in banners.slice(0,4)" :key="b.id">
          <div class="banner-wrap">
            <a v-if="b.link" :href="b.link" target="_blank" rel="noopener">
              <img :src="b.image_url" alt="banner" />
            </a>
            <img v-else :src="b.image_url" alt="banner" />
            <el-button
              v-if="me && isEditor(me.role)"
              class="banner-del" type="danger" plain size="small"
              @click.stop="removeBanner(b.id)"
            >删除</el-button>
          </div>
        </el-carousel-item>
      </el-carousel>

      <!-- 下方两列：左=新闻文字流（独立），右=新闻图片库（独立） -->
      <div class="news-two-col">
        <!-- 左列：新闻文字 -->
        <div class="news-col-left">
          <div class="news-line" v-for="n in news" :key="n.id">
            <h3 class="n-title">{{ n.title }}</h3>
            <p  class="n-summary">{{ n.summary }}</p>
            <div class="n-meta">
              <span>{{ fmtDate(n.published_at) }}</span>
              <el-link v-if="n.link" :href="n.link" target="_blank">详情</el-link>
              <template v-if="me && isEditor(me.role)">
                <el-divider direction="vertical" />
                <el-popconfirm title="删除这条新闻？" confirm-button-text="删除" cancel-button-text="取消"
                               @confirm="removeNews(n.id)">
                  <template #reference>
                    <el-button type="danger" plain size="small">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </div>
          </div>
          <el-empty v-if="!loading && news.length===0" description="暂无新闻" />
        </div>

        <!-- 右列：新闻图片库（独立的流，不与左侧配对） -->
        <div class="news-col-right">
          <div class="img-item" v-for="img in newsImages" :key="img.id">
            <div class="imgbox">
              <img :src="img.image_url" alt="news-image" />
            </div>
            <div class="img-actions" v-if="me && isEditor(me.role)">
              <el-popconfirm title="删除这张图片？" confirm-button-text="删除" cancel-button-text="取消"
                             @confirm="removeNewsImage(img.id)">
                <template #reference>
                  <el-button type="danger" plain size="small">删除</el-button>
                </template>
              </el-popconfirm>
            </div>
          </div>
          <el-empty v-if="!loading && newsImages.length===0" description="暂无图片" />
        </div>
      </div>

      <el-empty v-if="!loading && banners.length===0 && news.length===0 && newsImages.length===0"
                description="暂无内容" />
    </section>

    <!-- =================== 课题组概况（占位） =================== -->
    <section v-show="active==='about'">
      <div class="placeholder">这里放课题组概况（等你定版式再做）。</div>
    </section>

    <!-- =================== 研究亮点（占位） =================== -->
    <section v-show="active==='highlights'">
      <div class="placeholder">这里放研究亮点（等你定结构再做）。</div>
    </section>

    <!-- ========== 弹层：新增 Banner ========= -->
    <el-dialog v-model="showBanner" title="新增轮播图" width="620px">
    <el-form :model="bannerForm" label-width="92px">
        <el-form-item label="图片">
            <div class="upload-line">
            <el-image v-if="bannerForm.image_url" :src="bannerForm.image_url" fit="cover"
                        style="width:240px;height:150px;border-radius:8px;" />
            <el-upload :auto-upload="false" :show-file-list="false" accept="image/*" :on-change="onBannerSelected">
                <el-button>选择图片</el-button>
            </el-upload>
            <span class="hint">建议 1200×700，≤ 1MB</span>
            </div>
        </el-form-item>

        <el-form-item label="跳转链接">
            <el-input v-model="bannerForm.link" placeholder="可留空" />
        </el-form-item>

        <el-form-item label="排序">
            <el-input-number v-model="bannerForm.sort" :min="0" :max="999" />
            <span class="hint" style="margin-left:8px;">数值越小越靠前</span>
        </el-form-item>
    </el-form>
      <template #footer>
        <el-button @click="showBanner=false">取消</el-button>
        <el-button type="primary" :loading="submittingBanner" @click="submitBanner">保存</el-button>
      </template>
    </el-dialog>

    <!-- ========== 抽屉：管理 Banner ========= -->
    <el-drawer v-model="showBannerMgr" title="管理轮播图" size="50%">
      <el-table :data="banners" style="width:100%">
        <el-table-column label="预览" width="180">
          <template #default="{ row }">
            <img :src="row.image_url" style="width:160px;height:100px;object-fit:cover;border-radius:8px" />
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="120" />
        <el-table-column prop="link" label="链接" />
        <el-table-column label="操作" width="140">
          <template #default="{ row }">
            <el-popconfirm title="确定删除这张轮播图？" confirm-button-text="删除" cancel-button-text="取消"
                           @confirm="removeBanner(row.id)">
              <template #reference>
                <el-button type="danger" plain size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>

    <!-- ========== 弹层：上传 新闻图片 ========= -->
    <el-dialog v-model="showNewsImg" title="上传新闻图片" width="520px">
      <div class="upload-line">
        <el-image v-if="newsImgForm.image_url" :src="newsImgForm.image_url" fit="cover"
                  style="width:200px;height:120px;border-radius:8px;" />
        <el-upload :auto-upload="false" :show-file-list="false" accept="image/*" :on-change="onNewsImgSelected">
          <el-button>选择图片</el-button>
        </el-upload>
        <span class="hint">建议 3:2，≤ 1MB</span>
      </div>
      <template #footer>
        <el-button @click="showNewsImg=false">取消</el-button>
        <el-button type="primary" :loading="submittingNewsImg" @click="submitNewsImg">上传</el-button>
      </template>
    </el-dialog>

    <!-- ========== 抽屉：管理 新闻图库 ========= -->
    <el-drawer v-model="showNewsImgMgr" title="管理新闻图库" size="50%">
      <el-table :data="newsImages" style="width:100%">
        <el-table-column label="预览" width="180">
          <template #default="{ row }">
            <img :src="row.image_url" style="width:160px;height:100px;object-fit:cover;border-radius:8px" />
          </template>
        </el-table-column>
        <el-table-column prop="id" label="ID" width="100" />
        <el-table-column label="操作" width="140">
          <template #default="{ row }">
            <el-popconfirm title="删除这张图片？" confirm-button-text="删除" cancel-button-text="取消"
                           @confirm="removeNewsImage(row.id)">
              <template #reference>
                <el-button type="danger" plain size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>

    <!-- ========== 弹层：发布新闻（纯文字信息） ========= -->
    <el-dialog v-model="showPublish" title="发布新闻" width="760px">
      <el-form :model="form" label-width="96px">
        <el-form-item label="标题"><el-input v-model="form.title" maxlength="100" show-word-limit /></el-form-item>
        <el-form-item label="摘要">
          <el-input v-model="form.summary" type="textarea" :rows="4" maxlength="300" show-word-limit />
        </el-form-item>
        <el-form-item label="外部链接"><el-input v-model="form.link" placeholder="可填论文/新闻链接" /></el-form-item>
        <el-form-item label="发布时间">
          <el-date-picker v-model="form.published_at" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPublish=false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitNews">发布</el-button>
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
const active = ref<'news'|'about'|'highlights'>('news')
const loading = ref(false)
const isEditor = (r?:string) => r==='admin' || r==='editor'

/* 顶部 Banner 图库（独立系统） */
type Banner = { id:number; image_url:string; link?:string; sort:number }
const banners = ref<Banner[]>([])
const showBanner = ref(false)
const showBannerMgr = ref(false)
const submittingBanner = ref(false)
const bannerForm = ref<Partial<Banner>>({ sort: 0 })
function openBanner(){ if(!me.value || !isEditor(me.value.role)) return; bannerForm.value={ sort:0 }; showBanner.value=true }
function onBannerSelected(file: UploadFile){
  const raw = file.raw as File | undefined
  if(!raw) return
  if(raw.size > 1024*1024){ ElMessage.warning('请上传 ≤1MB 图片'); return }
  const rd = new FileReader(); rd.onload=()=>bannerForm.value.image_url = rd.result as string; rd.readAsDataURL(raw)
}
async function submitBanner(){
  if(!bannerForm.value.image_url){ ElMessage.warning('请先选择图片'); return }
  submittingBanner.value = true
  try{ await api.post('/banners', bannerForm.value); ElMessage.success('已保存'); showBanner.value=false; await fetchBanners() }
  catch(e:any){ ElMessage.error(e?.response?.data?.error || e?.message || '保存失败') }
  finally{ submittingBanner.value=false }
}
async function removeBanner(id:number){
  try{ await api.delete(`/banners/${id}`); ElMessage.success('已删除'); await fetchBanners() }
  catch(e:any){ ElMessage.error(e?.response?.data?.error || e?.message || '删除失败') }
}

/* 新闻图片库（独立系统，与 Banner 与新闻文字均解耦） */
type NewsImage = { id:number; image_url:string }
const newsImages = ref<NewsImage[]>([])
const showNewsImg = ref(false)
const showNewsImgMgr = ref(false)
const submittingNewsImg = ref(false)
const newsImgForm = ref<{ image_url?: string }>({})
function openUploadNewsImage(){ if(!me.value || !isEditor(me.value.role)) return; newsImgForm.value={}; showNewsImg.value=true }
function onNewsImgSelected(file: UploadFile){
  const raw = file.raw as File | undefined
  if(!raw) return
  if(raw.size > 1024*1024){ ElMessage.warning('请上传 ≤1MB 图片'); return }
  const rd = new FileReader(); rd.onload=()=>newsImgForm.value.image_url = rd.result as string; rd.readAsDataURL(raw)
}
async function submitNewsImg(){
  if(!newsImgForm.value.image_url) return ElMessage.warning('请先选择图片')
  submittingNewsImg.value = true
  try{ await api.post('/news_images', newsImgForm.value); ElMessage.success('已上传到新闻图库'); showNewsImg.value=false; await fetchNewsImages() }
  catch(e:any){ ElMessage.error(e?.response?.data?.error || e?.message || '上传失败') }
  finally{ submittingNewsImg.value=false }
}
async function removeNewsImage(id:number){
  try{ await api.delete(`/news_images/${id}`); ElMessage.success('已删除'); await fetchNewsImages() }
  catch(e:any){ ElMessage.error(e?.response?.data?.error || e?.message || '删除失败') }
}

/* 新闻文字流（独立系统） */
type News = { id:number; title:string; summary:string; link?:string; published_at:string }
const news = ref<News[]>([])
const showPublish = ref(false)
const submitting = ref(false)
const form = ref<Partial<News>>({ published_at: new Date().toISOString().slice(0,10) })
function openPublish(){ if(!me.value || !isEditor(me.value.role)) return; form.value={ published_at: new Date().toISOString().slice(0,10) }; showPublish.value=true }
async function submitNews(){
  if(!form.value.title)   return ElMessage.warning('请填写标题')
  if(!form.value.summary) return ElMessage.warning('请填写摘要')
  submitting.value = true
  try{ await api.post('/news', form.value); ElMessage.success('发布成功'); showPublish.value=false; await fetchNews() }
  catch(e:any){ ElMessage.error(e?.response?.data?.error || e?.message || '发布失败') }
  finally{ submitting.value = false }
}
async function removeNews(id:number){
  try{ await api.delete(`/news/${id}`); ElMessage.success('已删除'); await fetchNews() }
  catch(e:any){ ElMessage.error(e?.response?.data?.error || e?.message || '删除失败') }
}

/* 工具 */
function fmtDate(s?: string){ if(!s) return ''; return s.replace('T',' ').slice(0,10) }

/* 拉数据 */
async function fetchMe(){ try{ const {data}=await api.get('/auth/me'); me.value=data } catch{ me.value=null } }
async function fetchBanners(){ const {data}=await api.get<Banner[]>('/banners'); banners.value=data }
async function fetchNewsImages(){ const {data}=await api.get<NewsImage[]>('/news_images'); newsImages.value=data }
async function fetchNews(){
  loading.value = true
  try{ const {data}=await api.get<News[]>('/news'); news.value=data }
  finally{ loading.value = false }
}
onMounted(async () => { await fetchMe(); await fetchBanners(); await fetchNewsImages(); await fetchNews() })
</script>

<style scoped>
.home-wrap{ max-width:1100px; margin:24px auto; padding:0 16px; }
.home-tabs :deep(.el-tabs__item){ font-size:16px; }
.toolbar{ margin:-6px 0 10px; display:flex; gap:10px; align-items:center; }

/* 顶部轮播（纯图片） */
.banner-wrap{ position:relative; width:100%; height:100%; border-radius:12px; overflow:hidden; box-shadow:0 2px 8px rgba(0,0,0,.06); }
.banner-wrap img{ width:100%; height:100%; object-fit:cover; display:block; }
.banner-del{ position:absolute; right:12px; top:12px; z-index:10; }

/* 两列区域：左 60%（新闻文字），右 40%（新闻图库） */
.news-two-col{ display:flex; gap:24px; margin-top:16px; }
.news-col-left{ flex:0 0 60%; min-width:0; }
.news-col-right{ flex:0 0 40%; }

/* 新闻文字行：分割线 */
.news-line{ padding:18px 0; }
.news-line + .news-line{ border-top:1px solid #eee; }
.n-title{ font-size:20px; font-weight:800; margin:0 0 8px; line-height:1.3; }
.n-summary{ color:#444; font-size:15px; line-height:1.8; margin:0; }
.n-meta{ margin-top:10px; color:#777; display:flex; gap:12px; align-items:center; }

/* 右侧图片库：纵向流 */
.img-item{ margin-bottom:16px; }
.imgbox{ width:100%; aspect-ratio:16/9; background:#f6f7f9; border-radius:10px; overflow:hidden; }
.imgbox img{ width:100%; height:100%; object-fit:cover; display:block; }
.img-actions{ margin-top:6px; }

/* 占位样式 */
.placeholder{ padding:40px; color:#666; background:#fafafa; border:1px dashed #ddd; border-radius:10px; text-align:center; }

@media (max-width: 900px){
  .news-two-col{ flex-direction:column; }
  .news-col-left, .news-col-right{ flex:1 1 auto; }
}
.upload-line{ display:flex; align-items:center; gap:12px; }
.hint{ color:#888; font-size:12px; }
</style>
