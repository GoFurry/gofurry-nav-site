<template>
  <div class="page-shell bg-orange-50">
    <section class="mx-auto w-full max-w-5xl px-4 py-10 sm:px-6">
      <h1 class="text-4xl font-semibold text-slate-950">{{ isZh ? '抽奖活动' : 'Prize Events' }}</h1>
      <ErrorState v-if="errorMsg" class="mt-6" :title="isZh ? '加载失败' : 'Load failed'" :description="errorMsg" />
      <LoadingState v-if="loading" :label="isZh ? '正在加载...' : 'Loading...'" />
      <div v-else class="mt-8 grid gap-4">
        <article v-for="item in lottery?.active || []" :key="item.lottery.id" class="rounded-lg border border-orange-200 bg-white p-5">
          <h2 class="text-xl font-semibold text-slate-950">{{ item.lottery.title }}</h2>
          <p class="mt-2 text-sm leading-7 text-slate-600">{{ item.lottery.desc }}</p>
          <p class="mt-3 text-sm text-slate-500">{{ item.count }} {{ isZh ? '人参与' : 'participants' }}</p>
        </article>
        <EmptyState v-if="!lottery?.active?.length" :title="isZh ? '暂无活动' : 'No active events'" />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import type { LotteryResp } from '~/types/game'
import { getLottery } from '~/services/game'

const { locale } = useI18n()
const isZh = computed(() => locale.value === 'zh-CN')
const lottery = ref<LotteryResp | null>(null)
const loading = ref(true)
const errorMsg = ref('')

useSeoMeta({
  title: () => isZh.value ? 'GoFurry 抽奖活动' : 'GoFurry Prize Events',
  description: () => isZh.value ? '查看 GoFurry 游戏抽奖活动。' : 'View GoFurry game prize events.'
})

onMounted(async () => {
  try {
    lottery.value = await getLottery()
  } catch (error) {
    errorMsg.value = error instanceof Error ? error.message : String(error)
  } finally {
    loading.value = false
  }
})
</script>
