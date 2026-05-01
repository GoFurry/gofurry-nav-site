<template>
  <footer class="border-t border-orange-200/70 bg-white">
    <div class="mx-auto flex w-full max-w-7xl flex-col gap-3 px-4 py-8 text-sm text-slate-500 sm:px-6 md:flex-row md:items-center md:justify-between">
      <div>GoFurry © {{ year }}</div>
      <div class="flex flex-wrap gap-4">
        <NuxtLinkLocale to="/about" class="hover:text-orange-700">{{ locale === 'zh-CN' ? '关于' : 'About' }}</NuxtLinkLocale>
        <NuxtLinkLocale to="/updates" class="hover:text-orange-700">{{ locale === 'zh-CN' ? '更新日志' : 'Updates' }}</NuxtLinkLocale>
        <a href="https://github.com/GoFurry/gofurry-nav-site" target="_blank" rel="noopener noreferrer" class="hover:text-orange-700">GitHub</a>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
const year = new Date().getFullYear()
const { locale } = useI18n()
</script>
