<template>
  <div v-if="visibleSites.length" class="flex items-center gap-2 overflow-x-auto overflow-y-hidden">
    <a
        v-for="item in visibleSites"
        :key="item.id"
        :href="toExternalUrl(item.url)"
        target="_blank"
        rel="noopener noreferrer"
        class="group flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-white/15 bg-black/20 shadow-[0_10px_24px_rgba(15,23,42,0.18)] backdrop-blur-sm transition hover:-translate-y-0.5 hover:bg-white/15"
        :title="item.name"
        @click="handleVisit(item)"
    >
      <img
          v-if="!failedIcons[item.id]"
          :src="`https://favicon.im/${toExternalUrl(item.url)}?larger=true`"
          :alt="item.name"
          class="h-6 w-6 rounded-md object-cover"
          loading="lazy"
          @error="markIconFailed(item.id)"
      />
      <div
          v-else
          class="flex h-6 w-6 items-center justify-center rounded-md bg-slate-900/85 text-xs font-semibold text-white"
      >
        {{ item.name.slice(0, 1).toUpperCase() }}
      </div>
    </a>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { recordRecentSite, toExternalUrl } from '@/utils/recentSites'

export interface SiteStripItem {
  id: string
  name: string
  url: string
}

const props = withDefaults(defineProps<{
  sites: SiteStripItem[]
  maxItems?: number
}>(), {
  maxItems: 8,
})

const visibleSites = computed(() => props.sites.slice(0, props.maxItems))
const failedIcons = ref<Record<string, boolean>>({})

function markIconFailed(id: string) {
  failedIcons.value = {
    ...failedIcons.value,
    [id]: true,
  }
}

function handleVisit(site: SiteStripItem) {
  recordRecentSite(site)
}
</script>
