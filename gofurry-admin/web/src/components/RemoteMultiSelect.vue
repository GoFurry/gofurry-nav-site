<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { listJSON } from '../api'
import type { OptionItem } from '../types'

const props = defineProps<{ endpoint: string; modelValue: string[] }>()
const emit = defineEmits<{ 'update:modelValue': [value: string[]] }>()

const keyword = ref('')
const options = ref<OptionItem[]>([])

async function load() {
  const result = await listJSON<OptionItem>(props.endpoint, 1, 100, keyword.value)
  options.value = result.list
}

function toggle(id: string, checked: boolean) {
  const set = new Set(props.modelValue ?? [])
  if (checked) set.add(id)
  else set.delete(id)
  emit('update:modelValue', Array.from(set))
}

onMounted(load)
watch(keyword, load)
</script>

<template>
  <div class="space-y-2">
    <input
      v-model="keyword"
      class="w-full border border-[var(--line)] bg-black/20 px-3 py-2 text-sm outline-none focus:border-[var(--accent)]"
      placeholder="搜索可选项"
    />
    <div class="max-h-64 space-y-2 overflow-auto border border-[var(--line)] bg-black/10 p-3">
      <label v-for="item in options" :key="item.id" class="flex items-center gap-2 text-sm">
        <input
          type="checkbox"
          :checked="(modelValue ?? []).includes(item.id)"
          @change="toggle(item.id, ($event.target as HTMLInputElement).checked)"
        />
        <span>{{ item.label }}</span>
        <span class="text-[var(--text-muted)]">{{ item.extra }}</span>
      </label>
    </div>
  </div>
</template>
