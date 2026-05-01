<template>
  <div class="page-shell content-grid">
    <section class="mx-auto w-full max-w-7xl px-4 py-10 sm:px-6">
      <div class="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
        <div>
          <p class="text-sm font-semibold uppercase tracking-[0.24em] text-orange-700">Games</p>
          <h1 class="mt-3 text-4xl font-semibold text-slate-950">{{ isZh ? '兽游情报' : 'Furry Game Intel' }}</h1>
          <p class="mt-3 max-w-2xl text-sm leading-7 text-slate-700">
            {{ isZh ? '聚合最新、热门、免费和近期兽游内容。' : 'Latest, hot, free, and recent furry games in one place.' }}
          </p>
        </div>
        <NuxtLinkLocale to="/games/search" class="rounded-md bg-orange-500 px-4 py-3 text-sm font-semibold text-white transition hover:bg-orange-600">
          {{ isZh ? '高级搜索' : 'Advanced Search' }}
        </NuxtLinkLocale>
      </div>

      <ErrorState v-if="error" class="mt-8" :title="isZh ? '加载游戏失败' : 'Failed to load games'" :description="String(error.message || error)" />
      <LoadingState v-else-if="pending" :label="isZh ? '正在加载游戏...' : 'Loading games...'" />

      <div v-else class="mt-10 space-y-10">
        <section v-for="section in gameSections" :key="section.key">
          <h2 class="mb-4 text-2xl font-semibold text-slate-950">{{ section.title }}</h2>
          <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
            <NuxtLinkLocale v-for="game in section.items" :key="game.game_id" :to="`/games/${game.game_id}`" class="group rounded-lg border border-orange-200/70 bg-white/85 p-4 transition hover:-translate-y-0.5 hover:border-orange-400 hover:shadow-lg">
              <img :src="game.header" :alt="game.name" class="aspect-[16/9] w-full rounded-md object-cover" loading="lazy" />
              <h3 class="mt-4 line-clamp-1 text-lg font-semibold text-slate-950 group-hover:text-orange-700">{{ displayGameName(game) }}</h3>
              <p class="mt-2 line-clamp-2 text-sm leading-6 text-slate-600">{{ displayGameInfo(game) }}</p>
              <div class="mt-4 flex items-center justify-between text-xs text-slate-500">
                <span>{{ isZh ? '评分' : 'Score' }} {{ game.avg_score || '-' }}</span>
                <span>{{ game.comment_count || 0 }} {{ isZh ? '评论' : 'reviews' }}</span>
              </div>
            </NuxtLinkLocale>
          </div>
        </section>

        <section v-if="latestNews.length">
          <div class="mb-4 flex items-center justify-between">
            <h2 class="text-2xl font-semibold text-slate-950">{{ isZh ? '最新更新' : 'Latest Updates' }}</h2>
            <NuxtLinkLocale to="/games/news/more" class="text-sm font-semibold text-orange-700 hover:text-orange-900">{{ isZh ? '更多' : 'More' }}</NuxtLinkLocale>
          </div>
          <div class="grid gap-4 md:grid-cols-2">
            <a v-for="news in latestNews" :key="news.id" :href="news.url" target="_blank" rel="noopener noreferrer" class="rounded-lg border border-orange-200/70 bg-white/85 p-5 transition hover:border-orange-400">
              <h3 class="font-semibold text-slate-950">{{ news.headline }}</h3>
              <p class="mt-2 text-sm text-slate-500">{{ news.name }} · {{ displayDate(news.post_time) }}</p>
            </a>
          </div>
        </section>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import type { BaseGameInfoRecord } from '~/types/game'
import { getGameMainInfo, getLatestGameNews } from '~/services/game'
import { displayDate } from '~/utils/format'

const { locale } = useI18n()
const isZh = computed(() => locale.value === 'zh-CN')
const lang = computed(() => locale.value === 'zh-CN' ? 'zh' : 'en')

useSeoMeta({
  title: () => isZh.value ? 'GoFurry 兽游情报' : 'GoFurry Furry Games',
  description: () => isZh.value ? '浏览 GoFurry 收录的兽游、更新、评分和热门内容。' : 'Browse furry games, updates, scores, and trending content indexed by GoFurry.'
})

const { data, pending, error } = await useAsyncData(
  'games-page',
  async () => {
    const [main, news] = await Promise.all([
      getGameMainInfo(),
      getLatestGameNews().catch(() => null)
    ])
    return { main, news }
  }
)

const gameSections = computed(() => {
  const main = data.value?.main
  if (!main) return []
  return [
    { key: 'latest', title: isZh.value ? '最新收录' : 'Latest', items: main.latest || [] },
    { key: 'hot', title: isZh.value ? '热门游戏' : 'Hot', items: main.hot || [] },
    { key: 'free', title: isZh.value ? '免费游戏' : 'Free', items: main.free || [] },
    { key: 'recent', title: isZh.value ? '近期更新' : 'Recent', items: main.recent || [] }
  ].filter((section) => section.items.length)
})

const latestNews = computed(() => {
  const news = data.value?.news
  if (!news) return []
  return (lang.value === 'zh' ? news.news_zh : news.news_en).slice(0, 6)
})

function displayGameName(game: BaseGameInfoRecord) {
  return isZh.value ? game.name : (game.name_en || game.name)
}

function displayGameInfo(game: BaseGameInfoRecord) {
  return isZh.value ? game.info : (game.info_en || game.info)
}
</script>
