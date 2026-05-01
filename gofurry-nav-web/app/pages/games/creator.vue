<template>
  <div class="page-shell bg-orange-50">
    <section class="mx-auto w-full max-w-7xl px-4 py-10 sm:px-6">
      <p class="text-sm font-semibold uppercase tracking-[0.24em] text-orange-700">Creators</p>
      <h1 class="mt-3 text-4xl font-semibold text-slate-950">{{ isZh ? '游戏创作者' : 'Game Creators' }}</h1>
      <p class="mt-3 max-w-2xl text-sm leading-7 text-slate-700">
        {{ isZh ? '展示与 GoFurry 游戏内容相关的创作者与社交链接。' : 'Creators and social links related to GoFurry game content.' }}
      </p>

      <ErrorState v-if="error" class="mt-8" :title="isZh ? '加载创作者失败' : 'Failed to load creators'" :description="String(error.message || error)" />
      <LoadingState v-else-if="pending" :label="isZh ? '正在加载创作者...' : 'Loading creators...'" />
      <div v-else class="mt-8 grid gap-5 md:grid-cols-2 xl:grid-cols-3">
        <article v-for="creator in creators" :key="creator.id" class="rounded-lg border border-orange-200 bg-white p-5">
          <div class="flex items-center gap-4">
            <img :src="creator.avatar" :alt="creator.name" class="h-16 w-16 rounded-lg object-cover" loading="lazy" />
            <div>
              <h2 class="text-xl font-semibold text-slate-950">{{ creator.name }}</h2>
              <a v-if="creator.url" :href="creator.url" target="_blank" rel="noopener noreferrer" class="text-sm text-orange-700 hover:text-orange-900">{{ creator.url }}</a>
            </div>
          </div>
          <p class="mt-4 line-clamp-4 text-sm leading-7 text-slate-600">{{ creator.info }}</p>
        </article>
        <EmptyState v-if="creators.length === 0" :title="isZh ? '暂无创作者' : 'No creators yet'" />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { getGameCreator } from '~/services/game'

const { locale } = useI18n()
const isZh = computed(() => locale.value === 'zh-CN')
const lang = computed(() => locale.value === 'zh-CN' ? 'zh' : 'en')

useSeoMeta({
  title: () => isZh.value ? 'GoFurry 游戏创作者' : 'GoFurry Game Creators',
  description: () => isZh.value ? '浏览 GoFurry 收录的游戏创作者信息。' : 'Browse game creator profiles indexed by GoFurry.'
})

const { data, pending, error } = await useAsyncData(
  () => `game-creators-${lang.value}`,
  () => getGameCreator(lang.value),
  { watch: [lang] }
)
const creators = computed(() => data.value || [])
</script>
