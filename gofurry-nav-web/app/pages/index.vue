<template>
  <div class="bg-[#f5ead8]">
    <section class="relative flex min-h-[calc(100svh-72px)] items-center overflow-hidden">
      <div class="absolute inset-0 content-grid"></div>
      <div class="absolute inset-0 bg-[linear-gradient(90deg,rgba(15,23,42,0.88),rgba(15,23,42,0.56)_46%,rgba(15,23,42,0.16))]"></div>
      <div class="relative mx-auto grid w-full max-w-7xl gap-10 px-4 py-16 sm:px-6 lg:grid-cols-[minmax(0,0.95fr)_minmax(360px,0.7fr)] lg:items-center">
        <div class="max-w-3xl text-white">
          <p class="text-sm font-semibold uppercase tracking-[0.24em] text-orange-200">GoFurry</p>
          <h1 class="mt-5 text-5xl font-semibold leading-tight tracking-normal sm:text-6xl">
            {{ isZh ? '更清晰的兽圈导航与游戏发现入口' : 'A clearer entry point for furry navigation and games' }}
          </h1>
          <p class="mt-6 max-w-2xl text-base leading-8 text-slate-100">
            {{ isZh ? '把分散的网站、项目、游戏情报和长期内容整理成一个可被搜索引擎读取、也适合日常浏览的前台体验。' : 'GoFurry organizes scattered sites, projects, game updates, and long-form work into an SEO-friendly public experience.' }}
          </p>
          <div class="mt-8 flex flex-wrap gap-3">
            <NuxtLinkLocale to="/nav" class="rounded-md bg-orange-400 px-5 py-3 text-sm font-semibold text-slate-950 transition hover:bg-orange-300">
              {{ isZh ? '进入导航' : 'Open Navigation' }}
            </NuxtLinkLocale>
            <NuxtLinkLocale to="/games" class="rounded-md border border-white/30 px-5 py-3 text-sm font-semibold text-white transition hover:bg-white/10">
              {{ isZh ? '浏览游戏' : 'Browse Games' }}
            </NuxtLinkLocale>
          </div>
        </div>

        <div class="grid gap-4 text-white">
          <div v-for="item in features" :key="item.title" class="border-l border-orange-200/50 pl-5">
            <h2 class="text-xl font-semibold">{{ item.title }}</h2>
            <p class="mt-2 text-sm leading-7 text-slate-200">{{ item.desc }}</p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const { locale } = useI18n()
const isZh = computed(() => locale.value === 'zh-CN')

useSeoMeta({
  title: () => isZh.value ? 'GoFurry - 兽圈导航与游戏发现' : 'GoFurry - Furry Navigation and Game Discovery',
  description: () => isZh.value
    ? 'GoFurry 提供兽圈站点导航、游戏情报、更新记录和开放生态入口。'
    : 'GoFurry provides furry site navigation, game intelligence, updates, and open ecosystem entry points.',
  ogTitle: 'GoFurry',
  ogDescription: 'GoFurry navigation and game discovery site.',
  twitterCard: 'summary_large_image'
})

const features = computed(() => isZh.value
  ? [
      { title: '导航', desc: '站点分组、基础状态和入口信息集中展示，适合搜索和快速访问。' },
      { title: '游戏', desc: '围绕 Steam 兽游提供更新、详情、作者和检索入口。' },
      { title: 'SEO', desc: '公开页面迁移到 Nuxt SSR，让搜索引擎直接获取完整 HTML。' }
    ]
  : [
      { title: 'Navigation', desc: 'Grouped sites, status signals, and entry points are organized for search and quick access.' },
      { title: 'Games', desc: 'Steam furry game updates, details, creators, and search are available from one surface.' },
      { title: 'SEO', desc: 'Public pages are migrated to Nuxt SSR so crawlers can read complete HTML.' }
    ])
</script>
