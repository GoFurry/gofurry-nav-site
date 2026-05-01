<template>
  <header class="sticky top-0 z-30 border-b border-orange-200/70 bg-white/90 backdrop-blur">
    <div class="mx-auto flex h-[72px] w-full max-w-7xl items-center justify-between px-4 sm:px-6">
      <NuxtLinkLocale to="/" class="flex items-center gap-3">
        <img src="~/assets/svgs/logo-mini.svg" alt="GoFurry" class="h-10 w-10" />
        <div>
          <div class="text-lg font-semibold tracking-normal text-slate-950">GoFurry</div>
          <div class="text-xs text-slate-500">Navigation & Games</div>
        </div>
      </NuxtLinkLocale>

      <nav class="hidden items-center gap-1 text-sm font-medium text-slate-600 md:flex">
        <NuxtLinkLocale v-for="item in navItems" :key="item.to" :to="item.to" class="rounded-md px-3 py-2 transition hover:bg-orange-50 hover:text-orange-700">
          {{ item.label }}
        </NuxtLinkLocale>
      </nav>

      <div class="flex items-center gap-2">
        <button class="rounded-md border border-orange-200 px-3 py-2 text-sm text-orange-700 transition hover:bg-orange-50" type="button" @click="toggleLocale">
          {{ locale === 'zh-CN' ? 'EN' : '中' }}
        </button>
        <NuxtLinkLocale to="/games/search" class="rounded-md bg-orange-500 px-3 py-2 text-sm font-semibold text-white transition hover:bg-orange-600">
          {{ locale === 'zh-CN' ? '搜索' : 'Search' }}
        </NuxtLinkLocale>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
const { locale, setLocale } = useI18n()

const navItems = computed(() => {
  const isZh = locale.value === 'zh-CN'
  return [
    { to: '/nav', label: isZh ? '导航' : 'Nav' },
    { to: '/games', label: isZh ? '游戏' : 'Games' },
    { to: '/updates', label: isZh ? '更新' : 'Updates' },
    { to: '/about', label: isZh ? '关于' : 'About' }
  ]
})

function toggleLocale() {
  setLocale(locale.value === 'zh-CN' ? 'en' : 'zh-CN')
}
</script>
