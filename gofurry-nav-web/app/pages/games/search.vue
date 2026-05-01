<template>
  <div class="page-shell content-grid">
    <section class="mx-auto w-full max-w-6xl px-4 py-10 sm:px-6">
      <h1 class="text-4xl font-semibold text-slate-950">{{ isZh ? '游戏搜索' : 'Game Search' }}</h1>
      <form class="mt-6 flex flex-col gap-3 rounded-lg border border-orange-200 bg-white/90 p-4 md:flex-row" @submit.prevent="search">
        <input v-model="keyword" class="min-w-0 flex-1 rounded-md border border-orange-200 px-4 py-3 text-sm outline-none focus:border-orange-500" :placeholder="isZh ? '输入游戏名或关键词' : 'Game name or keyword'" />
        <button class="rounded-md bg-orange-500 px-5 py-3 text-sm font-semibold text-white hover:bg-orange-600" type="submit">
          {{ isZh ? '搜索' : 'Search' }}
        </button>
      </form>

      <ErrorState v-if="errorMsg" class="mt-6" :title="isZh ? '搜索失败' : 'Search failed'" :description="errorMsg" />
      <LoadingState v-if="loading" :label="isZh ? '正在搜索...' : 'Searching...'" />
      <div v-else class="mt-8 grid gap-4 md:grid-cols-2 xl:grid-cols-3">
        <NuxtLinkLocale v-for="item in results" :key="item.id" :to="`/games/${item.id}`" class="rounded-lg border border-orange-200 bg-white/90 p-4 transition hover:border-orange-400">
          <img :src="item.cover" :alt="item.name" class="aspect-[16/9] w-full rounded-md object-cover" loading="lazy" />
          <h2 class="mt-4 text-lg font-semibold text-slate-950">{{ item.name }}</h2>
          <p class="mt-2 line-clamp-2 text-sm leading-6 text-slate-600">{{ item.info }}</p>
        </NuxtLinkLocale>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import type { SearchItemModel } from '~/types/game'
import { getSearchSimple } from '~/services/game'

const { locale } = useI18n()
const isZh = computed(() => locale.value === 'zh-CN')
const lang = computed(() => locale.value === 'zh-CN' ? 'zh' : 'en')
const keyword = ref('')
const results = ref<SearchItemModel[]>([])
const loading = ref(false)
const errorMsg = ref('')

useSeoMeta({
  title: () => isZh.value ? 'GoFurry 游戏搜索' : 'GoFurry Game Search',
  description: () => isZh.value ? '搜索 GoFurry 收录的兽游。' : 'Search furry games indexed by GoFurry.'
})

async function search() {
  if (!keyword.value.trim()) return
  loading.value = true
  errorMsg.value = ''
  try {
    results.value = await getSearchSimple(lang.value, keyword.value.trim())
  } catch (error) {
    errorMsg.value = error instanceof Error ? error.message : String(error)
  } finally {
    loading.value = false
  }
}
</script>
