<template>
  <section class="overview-games" aria-labelledby="overview-games-title" data-overview-games>
    <div class="overview-section-heading">
      <div><p class="overview-kicker">{{ $t('insights.editorial.gamesKicker') }}</p><h2 id="overview-games-title">{{ $t('insights.editorial.gamesTitle') }}</h2></div>
      <NuxtLink :to="localePath('/insights/games')" class="overview-text-link">{{ $t('insights.editorial.explore') }} <span aria-hidden="true">↗</span></NuxtLink>
    </div>
    <p v-if="!panel" class="overview-unavailable">{{ $t('insights.emptyStates.unavailable') }}</p>
    <div v-else class="overview-pulse">
      <template v-for="slot in slots" :key="slot.role">
        <NuxtLink v-if="slot.game" :to="localePath(`/games/${encodeURIComponent(slot.game.id)}`)" class="overview-pulse__item" :class="{ 'overview-pulse__item--lead': slot.role === 'players' }" :data-pulse="slot.role">
          <InsightEntityMedia domain="game" :entity="entity(slot.game)" />
          <span class="overview-pulse__body">
            <span class="overview-kicker">{{ $t(`insights.editorial.pulse.${slot.role}`) }}</span>
            <strong>{{ slot.game.name }}</strong>
            <template v-if="slot.role === 'players'">
              <span class="overview-pulse__count">{{ number(slot.game.online_count.count) }} <small>{{ $t('insights.editorial.playersUnit') }}</small></span>
              <span class="overview-pulse__detail">{{ $t('insights.editorial.latestObservation') }}</span>
              <span v-if="typeof slot.game.online_count.peak_count === 'number' && Number.isFinite(slot.game.online_count.peak_count)" class="overview-pulse__detail">{{ $t('game.panel.onlinePeak') }} {{ number(slot.game.online_count.peak_count) }}</span>
            </template>
            <span v-else-if="slot.role === 'discount' && discount" class="overview-pulse__detail">
              <b>{{ formatMinorAmount(discount.final_amount, discount.currency, locale) }}</b>
              · {{ $t('insights.editorial.discount', { percent: discount.discount_percent }) }} · {{ $t('insights.editorial.usRegion') }} / {{ discount.currency }}
            </span>
            <span v-else-if="slot.role === 'latest' && slot.game.release_date" class="overview-pulse__detail">{{ slot.game.release_date }}</span>
          </span>
        </NuxtLink>
        <div v-else class="overview-pulse__empty" :data-pulse="slot.role"><span class="overview-kicker">{{ $t(`insights.editorial.pulse.${slot.role}`) }}</span><p>{{ $t('insights.emptyStates.unavailable') }}</p></div>
      </template>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { GameV2ListItem, GameV2PanelRecord } from '@/types/game'
import type { InsightEntityRef } from '@/types/insights'
import { pulseDiscountPrice, selectOverviewPulse } from '@/utils/insightOverview'
import { formatMinorAmount } from '@/utils/insightPrices'
import InsightEntityMedia from '../entity/InsightEntityMedia.vue'

const props = defineProps<{ panel: GameV2PanelRecord | null }>()
const { locale } = useI18n()
const localePath = useLocalePath()
const pulse = computed(() => selectOverviewPulse(props.panel))
const slots = computed(() => (['players', 'discount', 'latest'] as const).map(role => ({ role, game: pulse.value[role] })))
const discount = computed(() => pulse.value.discount ? pulseDiscountPrice(pulse.value.discount) : null)
const number = (value: number) => new Intl.NumberFormat(locale.value === 'en' ? 'en-US' : 'zh-CN').format(value)
function entity(game: GameV2ListItem): InsightEntityRef {
  return { id: Number(game.id), name: game.name, visual: { kind: 'game_header', asset: game.header_url || game.capsule_url || null } }
}
</script>
