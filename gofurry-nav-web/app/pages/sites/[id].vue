<template>
  <div class="page-shell bg-orange-50">
    <section class="mx-auto w-full max-w-6xl px-4 py-10 sm:px-6">
      <ErrorState v-if="error" :title="isZh ? '加载站点详情失败' : 'Failed to load site detail'" :description="String(error.message || error)" />
      <LoadingState v-else-if="pending" :label="isZh ? '正在加载站点详情...' : 'Loading site detail...'" />

      <article v-else-if="site" class="space-y-8">
        <header class="rounded-lg border border-orange-200 bg-white p-6">
          <div class="flex flex-col gap-5 md:flex-row md:items-center">
            <img :src="site.icon || assets.siteDefaultLogo" :alt="site.name" class="h-20 w-20 rounded-lg object-cover" />
            <div class="min-w-0 flex-1">
              <div class="flex flex-wrap items-center gap-2">
                <h1 class="text-4xl font-semibold text-slate-950">{{ site.name }}</h1>
                <span v-if="site.nsfw === '1'" class="rounded bg-red-100 px-2 py-1 text-xs font-semibold text-red-700">NSFW</span>
                <span v-if="site.welfare === '1'" class="rounded bg-emerald-100 px-2 py-1 text-xs font-semibold text-emerald-700">{{ isZh ? '公益' : 'Welfare' }}</span>
              </div>
              <p v-if="domain" class="mt-2 text-sm text-slate-500">{{ domain }}</p>
              <p class="mt-4 max-w-3xl text-sm leading-7 text-slate-700">{{ site.info }}</p>
            </div>
          </div>
        </header>

        <div class="grid gap-6 lg:grid-cols-3">
          <section class="rounded-lg border border-orange-200 bg-white p-5">
            <h2 class="text-lg font-semibold text-slate-950">HTTP</h2>
            <dl class="mt-4 space-y-3 text-sm">
              <div class="flex justify-between gap-4"><dt class="text-slate-500">Status</dt><dd class="font-medium">{{ httpRecord?.statusCode || '-' }}</dd></div>
              <div class="flex justify-between gap-4"><dt class="text-slate-500">Server</dt><dd class="truncate font-medium">{{ httpRecord?.server || '-' }}</dd></div>
              <div class="flex justify-between gap-4"><dt class="text-slate-500">TLS</dt><dd class="font-medium">{{ httpRecord?.tlsVersion || '-' }}</dd></div>
            </dl>
          </section>
          <section class="rounded-lg border border-orange-200 bg-white p-5">
            <h2 class="text-lg font-semibold text-slate-950">Ping</h2>
            <dl class="mt-4 space-y-3 text-sm">
              <div class="flex justify-between gap-4"><dt class="text-slate-500">20</dt><dd class="font-medium">{{ pingRecord?.twenty?.avgDelay || '-' }}</dd></div>
              <div class="flex justify-between gap-4"><dt class="text-slate-500">60</dt><dd class="font-medium">{{ pingRecord?.sixty?.avgDelay || '-' }}</dd></div>
              <div class="flex justify-between gap-4"><dt class="text-slate-500">100</dt><dd class="font-medium">{{ pingRecord?.hundred?.avgDelay || '-' }}</dd></div>
            </dl>
          </section>
          <section class="rounded-lg border border-orange-200 bg-white p-5">
            <h2 class="text-lg font-semibold text-slate-950">DNS</h2>
            <dl class="mt-4 space-y-3 text-sm">
              <div class="flex justify-between gap-4"><dt class="text-slate-500">A</dt><dd class="font-medium">{{ dnsCount('a') }}</dd></div>
              <div class="flex justify-between gap-4"><dt class="text-slate-500">AAAA</dt><dd class="font-medium">{{ dnsCount('AAAA') }}</dd></div>
              <div class="flex justify-between gap-4"><dt class="text-slate-500">CNAME</dt><dd class="font-medium">{{ dnsCount('CNAME') }}</dd></div>
            </dl>
          </section>
        </div>
      </article>

      <EmptyState v-else :title="isZh ? '未找到站点' : 'Site not found'" />
    </section>
  </div>
</template>

<script setup lang="ts">
import type { DnsItem, DnsRecord } from '~/types/nav'
import { getSiteDetail, getSiteDnsRecord, getSiteHttpRecord, getSitePingRecord } from '~/services/site'
import { safeJsonParse, truncate } from '~/utils/format'

const route = useRoute()
const { locale } = useI18n()
const isZh = computed(() => locale.value === 'zh-CN')
const lang = computed(() => locale.value === 'zh-CN' ? 'zh' : 'en')
const assets = useSiteAssets()
const domain = computed(() => String(route.query.domain || ''))
const siteId = computed(() => String(route.params.id || ''))

const { data, pending, error } = await useAsyncData(
  () => `site-detail-${siteId.value}-${lang.value}-${domain.value}`,
  async () => {
    const site = await getSiteDetail(siteId.value, lang.value)
    const [httpRecord, dnsRecord, pingRecord] = domain.value
      ? await Promise.all([
          getSiteHttpRecord(domain.value).catch(() => null),
          getSiteDnsRecord(domain.value).catch(() => null),
          getSitePingRecord(domain.value).catch(() => null)
        ])
      : [null, null, null]

    return { site, httpRecord, dnsRecord, pingRecord }
  },
  { watch: [lang] }
)

const site = computed(() => data.value?.site || null)
const httpRecord = computed(() => data.value?.httpRecord || null)
const dnsRecord = computed(() => data.value?.dnsRecord || null)
const pingRecord = computed(() => data.value?.pingRecord || null)

function dnsCount(key: keyof DnsRecord) {
  const values = safeJsonParse<DnsItem[]>(dnsRecord.value?.[key], [])
  return values.length
}

useSeoMeta({
  title: () => site.value?.name ? `${site.value.name} - GoFurry` : 'GoFurry Site Detail',
  description: () => truncate(site.value?.info, 150) || 'GoFurry site detail',
  ogTitle: () => site.value?.name || 'GoFurry',
  ogDescription: () => truncate(site.value?.info, 150) || 'GoFurry site detail',
  ogImage: () => site.value?.icon || assets.siteDefaultLogo,
  twitterCard: 'summary_large_image'
})
</script>
