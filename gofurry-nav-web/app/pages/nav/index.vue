<template>
  <div class="page-shell content-grid">
    <section class="mx-auto w-full max-w-7xl px-4 py-10 sm:px-6">
      <div class="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
        <div>
          <p class="text-sm font-semibold uppercase tracking-[0.24em] text-orange-700">Navigation</p>
          <h1 class="mt-3 text-4xl font-semibold text-slate-950">{{ isZh ? '站点导航' : 'Site Navigation' }}</h1>
          <p class="mt-3 max-w-2xl text-sm leading-7 text-slate-700">
            {{ isZh ? '按分组浏览已收录站点，SSR 首屏输出便于搜索引擎抓取。' : 'Browse indexed sites by group with SSR output for search engines.' }}
          </p>
        </div>
        <input v-model="keyword" class="w-full rounded-md border border-orange-200 bg-white/90 px-4 py-3 text-sm outline-none transition focus:border-orange-500 md:w-80" :placeholder="isZh ? '搜索站点' : 'Search sites'" />
      </div>

      <ErrorState v-if="error" class="mt-8" :title="isZh ? '加载导航失败' : 'Failed to load navigation'" :description="String(error.message || error)" />
      <LoadingState v-else-if="pending" :label="isZh ? '正在加载导航...' : 'Loading navigation...'" />

      <div v-else class="mt-10 space-y-10">
        <section v-for="group in visibleGroups" :key="group.id">
          <div class="mb-4 flex items-center justify-between border-b border-orange-200 pb-3">
            <h2 class="text-2xl font-semibold text-slate-950">{{ group.name }}</h2>
            <span class="text-sm text-slate-500">{{ sitesByGroup(group).length }}</span>
          </div>
          <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
            <NuxtLinkLocale
              v-for="site in sitesByGroup(group)"
              :key="site.id"
              :to="{ path: `/sites/${site.id}`, query: { domain: site.domain } }"
              class="group min-h-36 rounded-lg border border-orange-200/70 bg-white/85 p-5 transition hover:-translate-y-0.5 hover:border-orange-400 hover:shadow-lg"
            >
              <div class="flex items-start gap-4">
                <img :src="site.icon || assets.siteDefaultLogo" :alt="site.name" class="h-11 w-11 rounded-md object-cover" loading="lazy" />
                <div class="min-w-0">
                  <h3 class="truncate text-lg font-semibold text-slate-950 group-hover:text-orange-700">{{ site.name }}</h3>
                  <p class="mt-1 truncate text-xs text-slate-500">{{ site.domain }}</p>
                </div>
              </div>
              <p class="mt-4 line-clamp-2 text-sm leading-6 text-slate-600">{{ site.info }}</p>
            </NuxtLinkLocale>
          </div>
        </section>
        <EmptyState v-if="visibleSiteCount === 0" :title="isZh ? '暂无匹配站点' : 'No matching sites'" />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import type { Group, Site } from '~/types/nav'
import { getGroups, getPing, getSites } from '~/services/nav'

const { locale } = useI18n()
const isZh = computed(() => locale.value === 'zh-CN')
const keyword = ref('')
const assets = useSiteAssets()

useSeoMeta({
  title: () => isZh.value ? 'GoFurry 站点导航' : 'GoFurry Site Navigation',
  description: () => isZh.value ? '按分组浏览 GoFurry 收录的兽圈站点导航。' : 'Browse furry sites indexed by GoFurry by group.'
})

const lang = computed(() => locale.value === 'zh-CN' ? 'zh' : 'en')
const { data, pending, error } = await useAsyncData(
  () => `nav-page-${lang.value}`,
  async () => {
    const [sites, groups, ping] = await Promise.all([
      getSites(lang.value),
      getGroups(lang.value),
      getPing().catch(() => ({}))
    ])

    return { sites, groups, ping }
  },
  { watch: [lang] }
)

const sites = computed<Site[]>(() => data.value?.sites || [])
const groups = computed<Group[]>(() => data.value?.groups || [])
const normalizedKeyword = computed(() => keyword.value.trim().toLowerCase())

function sitesByGroup(group: Group) {
  const ids = new Set(group.sites || [])
  return sites.value.filter((site) => {
    const inGroup = ids.has(site.id)
    if (!inGroup) return false
    if (!normalizedKeyword.value) return true
    return `${site.name} ${site.domain} ${site.info}`.toLowerCase().includes(normalizedKeyword.value)
  })
}

const visibleGroups = computed(() => groups.value.filter((group) => sitesByGroup(group).length > 0))
const visibleSiteCount = computed(() => visibleGroups.value.reduce((count, group) => count + sitesByGroup(group).length, 0))
</script>
