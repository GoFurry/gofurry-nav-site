<template>
  <div class="page-shell bg-[#f5ead8]">
    <section class="mx-auto w-full max-w-5xl px-4 py-14 sm:px-6">
      <p class="text-sm font-semibold uppercase tracking-[0.24em] text-orange-600">About GoFurry</p>
      <h1 class="mt-4 text-4xl font-semibold tracking-normal text-slate-950">
        {{ isZh ? '一个持续整理兽圈公共入口的项目' : 'A project for organizing public furry entry points' }}
      </h1>
      <p class="mt-5 max-w-3xl text-base leading-8 text-slate-700">
        {{ isZh ? 'GoFurry 希望把导航、游戏情报、开放平台与 DeepFurry 等内容放在更清晰、稳定、便于搜索的前台体验里。' : 'GoFurry brings navigation, game intelligence, open platform surfaces, and DeepFurry into a clearer, stable, and searchable frontend.' }}
      </p>

      <div class="mt-10 grid gap-6 md:grid-cols-2">
        <section v-for="item in sections" :key="item.title" class="border-t border-orange-200 pt-6">
          <h2 class="text-2xl font-semibold text-slate-950">{{ item.title }}</h2>
          <p class="mt-3 text-sm leading-7 text-slate-600">{{ item.desc }}</p>
        </section>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const { locale } = useI18n()
const isZh = computed(() => locale.value === 'zh-CN')

useSeoMeta({
  title: () => isZh.value ? '关于 GoFurry' : 'About GoFurry',
  description: () => isZh.value ? '了解 GoFurry 的导航、游戏情报和开放生态方向。' : 'Learn about GoFurry navigation, game intelligence, and open ecosystem direction.'
})

const sections = computed(() => isZh.value
  ? [
      { title: '导航与发现', desc: '把分散的网站入口、分组和基础状态整合成清晰的浏览路径。' },
      { title: '游戏情报', desc: '围绕兽游构建更新、检索、详情和作者视角。' },
      { title: '开放平台', desc: '承接后续开放能力、公共服务和外部扩展页面。' },
      { title: '长期维护', desc: '以公共性、开放性和可维护性作为主要方向。' }
    ]
  : [
      { title: 'Navigation', desc: 'Scattered site entry points, groups, and status signals are shaped into clearer paths.' },
      { title: 'Game Intel', desc: 'Furry game updates, search, details, and creator perspectives are organized for discovery.' },
      { title: 'Open Platform', desc: 'Future public services and ecosystem pages can extend from this surface.' },
      { title: 'Maintenance', desc: 'Public value, openness, and maintainability guide the project.' }
    ])
</script>
