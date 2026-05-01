<template>
  <div class="page-shell content-grid">
    <section class="mx-auto w-full max-w-5xl px-4 py-10 sm:px-6">
      <p class="text-sm font-semibold uppercase tracking-[0.24em] text-orange-700">Game News</p>
      <h1 class="mt-3 text-4xl font-semibold text-slate-950">{{ isZh ? '更多游戏更新' : 'More Game News' }}</h1>

      <ErrorState v-if="error" class="mt-8" :title="isZh ? '加载新闻失败' : 'Failed to load news'" :description="String(error.message || error)" />
      <LoadingState v-else-if="pending" :label="isZh ? '正在加载新闻...' : 'Loading news...'" />
      <div v-else class="mt-8 space-y-4">
        <a v-for="news in list" :key="news.id" :href="news.url" target="_blank" rel="noopener noreferrer" class="block rounded-lg border border-orange-200 bg-white/90 p-5 transition hover:border-orange-400">
          <h2 class="text-xl font-semibold text-slate-950">{{ news.headline }}</h2>
          <p class="mt-2 text-sm text-slate-500">{{ news.name }} · {{ news.author }} · {{ displayDate(news.post_time) }}</p>
          <p class="mt-3 line-clamp-3 text-sm leading-7 text-slate-700">{{ stripHtml(news.content) }}</p>
        </a>
        <EmptyState v-if="list.length === 0" :title="isZh ? '暂无新闻' : 'No news yet'" />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { getMoreLatestGameNews } from '~/services/game'
import { displayDate, stripHtml } from '~/utils/format'

const { locale } = useI18n()
const isZh = computed(() => locale.value === 'zh-CN')
const lang = computed(() => locale.value === 'zh-CN' ? 'zh' : 'en')

useSeoMeta({
  title: () => isZh.value ? 'GoFurry 更多游戏更新' : 'GoFurry More Game News',
  description: () => isZh.value ? '浏览更多 GoFurry 游戏更新与新闻。' : 'Browse more GoFurry game updates and news.'
})

const { data, pending, error } = await useAsyncData(
  () => `more-game-news-${lang.value}`,
  () => getMoreLatestGameNews(lang.value),
  { watch: [lang] }
)
const list = computed(() => data.value || [])
</script>
