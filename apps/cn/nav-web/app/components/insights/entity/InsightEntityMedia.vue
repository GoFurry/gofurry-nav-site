<template>
  <span class="insight-entity-media" :class="`insight-entity-media--${domain}`" :data-media-state="src ? 'image' : 'fallback'">
    <img v-if="src" ref="mediaImage" :key="src" :src="src" :alt="entity.name" :loading="eager ? 'eager' : 'lazy'" decoding="async" @error="onError" />
    <span v-else class="insight-entity-media__fallback" role="img" :aria-label="entity.name">
      <span aria-hidden="true">{{ initials }}</span>
    </span>
  </span>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import type { InsightDomain, InsightEntityRef } from '@/types/insights'

const props = defineProps<{ domain: InsightDomain, entity: InsightEntityRef, eager?: boolean }>()
const { siteLogoPrefixUrl, siteDefaultLogo } = useSiteAssets()
const failure = ref(0)
const mediaImage = ref<HTMLImageElement | null>(null)
const asset = computed(() => props.entity.visual?.kind === (props.domain === 'site' ? 'site_icon' : 'game_header')
  ? props.entity.visual.asset?.trim() || '' : '')
const defaultLogo = computed(() => siteAsset(String(siteDefaultLogo || '')))
const src = computed(() => {
  if (props.domain === 'game') return failure.value ? '' : asset.value
  if (failure.value > 1) return ''
  return failure.value || !asset.value ? defaultLogo.value : siteAsset(asset.value)
})
const initials = computed(() => Array.from(props.entity.name.trim()).slice(0, 2).join('').toUpperCase() || '—')
watch(() => [props.entity.id, props.domain, asset.value], () => { failure.value = 0 })
onMounted(() => {
  // An SSR image can fail before hydration attaches the error listener.
  if (mediaImage.value?.complete && mediaImage.value.naturalWidth === 0) onError()
})

function siteAsset(value: string) {
  if (!value || /^(?:https?:)?\/\//i.test(value) || value.startsWith('data:')) return value
  return `${String(siteLogoPrefixUrl).replace(/\/+$/, '')}/${value.replace(/^\/+/, '')}`
}
function onError() {
  failure.value = props.domain === 'site' && src.value !== defaultLogo.value ? 1 : 2
}
</script>
