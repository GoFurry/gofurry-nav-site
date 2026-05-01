<template>
  <div class="page-shell content-grid">
    <section class="mx-auto w-full max-w-7xl px-4 py-10 sm:px-6">
      <ErrorState v-if="error" :title="isZh ? '加载游戏详情失败' : 'Failed to load game detail'" :description="String(error.message || error)" />
      <LoadingState v-else-if="pending" :label="isZh ? '正在加载游戏详情...' : 'Loading game detail...'" />

      <article v-else-if="game" class="grid gap-8 lg:grid-cols-[minmax(0,1fr)_320px]">
        <main class="rounded-lg border border-orange-200 bg-white/90 p-6">
          <img :src="game.cover" :alt="game.name" class="aspect-[16/9] w-full rounded-lg object-cover" />
          <h1 class="mt-6 text-4xl font-semibold text-slate-950">{{ game.name }}</h1>
          <p class="mt-4 text-sm leading-7 text-slate-700">{{ game.info }}</p>

          <section class="mt-8">
            <h2 class="text-2xl font-semibold text-slate-950">{{ isZh ? '游戏介绍' : 'About' }}</h2>
            <p class="mt-3 whitespace-pre-line text-sm leading-7 text-slate-700">{{ aboutText }}</p>
          </section>

          <section v-if="game.news?.length" class="mt-8">
            <h2 class="text-2xl font-semibold text-slate-950">{{ isZh ? '相关新闻' : 'News' }}</h2>
            <div class="mt-4 space-y-3">
              <a v-for="news in game.news.slice(0, 5)" :key="news.url" :href="news.url" target="_blank" rel="noopener noreferrer" class="block rounded-md border border-orange-100 p-4 transition hover:border-orange-400">
                <h3 class="font-semibold text-slate-950">{{ news.headline }}</h3>
                <p class="mt-1 text-xs text-slate-500">{{ news.author }} · {{ displayDate(news.post_time) }}</p>
              </a>
            </div>
          </section>
        </main>

        <aside class="space-y-6">
          <section class="rounded-lg border border-orange-200 bg-white/90 p-5">
            <h2 class="text-lg font-semibold text-slate-950">{{ isZh ? '信息' : 'Info' }}</h2>
            <dl class="mt-4 space-y-3 text-sm">
              <div><dt class="text-slate-500">{{ isZh ? '发行日期' : 'Release' }}</dt><dd class="mt-1 font-medium">{{ game.release_date || '-' }}</dd></div>
              <div><dt class="text-slate-500">{{ isZh ? '开发者' : 'Developers' }}</dt><dd class="mt-1 font-medium">{{ game.developers?.join(', ') || '-' }}</dd></div>
              <div><dt class="text-slate-500">{{ isZh ? '在线人数' : 'Online' }}</dt><dd class="mt-1 font-medium">{{ game.online_count || '-' }}</dd></div>
            </dl>
          </section>

          <section v-if="recommended?.length" class="rounded-lg border border-orange-200 bg-white/90 p-5">
            <h2 class="text-lg font-semibold text-slate-950">{{ isZh ? '相似推荐' : 'Similar' }}</h2>
            <div class="mt-4 space-y-3">
              <NuxtLinkLocale v-for="item in recommended" :key="item.id" :to="`/games/${item.id}`" class="block rounded-md bg-orange-50 p-3 text-sm hover:bg-orange-100">
                <div class="font-semibold text-slate-950">{{ item.name }}</div>
                <div class="mt-1 line-clamp-2 text-slate-600">{{ item.info }}</div>
              </NuxtLinkLocale>
            </div>
          </section>

          <section v-if="remark" class="rounded-lg border border-orange-200 bg-white/90 p-5">
            <h2 class="text-lg font-semibold text-slate-950">{{ isZh ? '评论概览' : 'Reviews' }}</h2>
            <p class="mt-3 text-sm text-slate-600">{{ isZh ? '平均分' : 'Average score' }} {{ remark.avg_score || '-' }} · {{ remark.total || 0 }}</p>
          </section>
        </aside>
      </article>

      <EmptyState v-else :title="isZh ? '未找到游戏' : 'Game not found'" />
    </section>
  </div>
</template>

<script setup lang="ts">
import { getGameBaseInfo, getGameRemark, getRecommendedGame } from '~/services/game'
import { displayDate, stripHtml, truncate } from '~/utils/format'

const route = useRoute()
const { locale } = useI18n()
const isZh = computed(() => locale.value === 'zh-CN')
const lang = computed(() => locale.value === 'zh-CN' ? 'zh' : 'en')
const gameId = computed(() => String(route.params.id || ''))

const { data, pending, error } = await useAsyncData(
  () => `game-detail-${gameId.value}-${lang.value}`,
  async () => {
    const [game, remark, recommended] = await Promise.all([
      getGameBaseInfo(gameId.value, lang.value),
      getGameRemark(gameId.value).catch(() => null),
      getRecommendedGame(gameId.value, lang.value).catch(() => [])
    ])
    return { game, remark, recommended }
  },
  { watch: [lang] }
)

const game = computed(() => data.value?.game || null)
const remark = computed(() => data.value?.remark || null)
const recommended = computed(() => data.value?.recommended || [])
const aboutText = computed(() => stripHtml(game.value?.about_the_game || game.value?.detailed_description || game.value?.info))

useSeoMeta({
  title: () => game.value?.name ? `${game.value.name} - GoFurry` : 'GoFurry Game Detail',
  description: () => truncate(game.value?.info || aboutText.value, 150) || 'GoFurry game detail',
  ogTitle: () => game.value?.name || 'GoFurry',
  ogDescription: () => truncate(game.value?.info || aboutText.value, 150) || 'GoFurry game detail',
  ogImage: () => game.value?.cover || '/og-default.png',
  twitterCard: 'summary_large_image'
})
</script>
