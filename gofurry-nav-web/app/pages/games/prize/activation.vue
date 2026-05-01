<template>
  <div class="page-shell bg-orange-50">
    <section class="mx-auto w-full max-w-3xl px-4 py-16 text-center sm:px-6">
      <h1 class="text-4xl font-semibold text-slate-950">{{ isZh ? '抽奖确认' : 'Prize Activation' }}</h1>
      <p class="mt-4 text-sm leading-7 text-slate-600">
        {{ isZh ? '该页面保留为客户端交互入口，后续会继续兼容现有抽奖确认接口。' : 'This client-only page remains as the entry point for the existing prize activation flow.' }}
      </p>
    </section>
  </div>
</template>

<script setup lang="ts">
const { locale } = useI18n()
const isZh = computed(() => locale.value === 'zh-CN')

useSeoMeta({
  title: () => isZh.value ? 'GoFurry 抽奖确认' : 'GoFurry Prize Activation',
  description: () => isZh.value ? 'GoFurry 抽奖确认入口。' : 'GoFurry prize activation entry.'
})
</script>
