<template>
  <div class="page-shell bg-orange-50">
    <section class="mx-auto w-full max-w-5xl px-4 py-10 sm:px-6">
      <p class="text-sm font-semibold uppercase tracking-[0.24em] text-orange-700">Updates</p>
      <h1 class="mt-3 text-4xl font-semibold text-slate-950">{{ isZh ? '更新记录' : 'Updates' }}</h1>
      <p class="mt-3 max-w-2xl text-sm leading-7 text-slate-700">
        {{ isZh ? 'GoFurry 公开更新记录，内容来自现有后端 changelog 接口。' : 'Public GoFurry update records from the existing changelog API.' }}
      </p>

      <ErrorState v-if="error" class="mt-8" :title="isZh ? '加载更新记录失败' : 'Failed to load updates'" :description="String(error.message || error)" />
      <LoadingState v-else-if="pending" :label="isZh ? '正在加载更新记录...' : 'Loading updates...'" />
      <div v-else class="mt-8 space-y-4">
        <a v-for="item in list" :key="item.url" :href="item.url" target="_blank" rel="noopener noreferrer" class="block rounded-lg border border-orange-200 bg-white p-5 transition hover:border-orange-400">
          <div class="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
            <h2 class="text-xl font-semibold text-slate-950">{{ item.title }}</h2>
            <time class="text-sm text-slate-500">{{ displayDate(item.create_time) }}</time>
          </div>
          <p class="mt-3 text-sm text-slate-600">{{ isZh ? '打开 Markdown 更新详情' : 'Open Markdown update detail' }}</p>
        </a>
        <EmptyState v-if="list.length === 0" :title="isZh ? '暂无更新记录' : 'No updates yet'" />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { getChangeLog } from '~/services/update'
import { displayDate } from '~/utils/format'

const { locale } = useI18n()
const isZh = computed(() => locale.value === 'zh-CN')

useSeoMeta({
  title: () => isZh.value ? 'GoFurry 更新记录' : 'GoFurry Updates',
  description: () => isZh.value ? '查看 GoFurry 的公开更新记录和公告。' : 'Read public GoFurry update records and announcements.'
})

const { data, pending, error } = await useAsyncData('updates-page', () => getChangeLog())
const list = computed(() => data.value || [])
</script>
