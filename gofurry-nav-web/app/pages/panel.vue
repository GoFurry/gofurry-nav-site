<template>
  <div class="page-shell bg-orange-50">
    <section class="mx-auto w-full max-w-4xl px-4 py-16 sm:px-6">
      <h1 class="text-4xl font-semibold text-slate-950">{{ isZh ? '数据面板' : 'Metrics Panel' }}</h1>
      <p class="mt-4 text-sm leading-7 text-slate-600">
        {{ isZh ? '该入口依赖图表和浏览器交互，迁移到 Nuxt 后作为客户端页面保留。' : 'This entry depends on charts and browser interactions, so it remains client-rendered after the Nuxt migration.' }}
      </p>
    </section>
  </div>
</template>

<script setup lang="ts">
const { locale } = useI18n()
const isZh = computed(() => locale.value === 'zh-CN')
</script>
