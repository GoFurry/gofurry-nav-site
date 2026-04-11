<script setup lang="ts">
const props = defineProps<{ modelValue: string[] }>()
const emit = defineEmits<{ 'update:modelValue': [value: string[]] }>()

function update(index: number, value: string) {
  const next = [...(props.modelValue ?? [])]
  next[index] = value
  emit('update:modelValue', next)
}

function addRow() {
  emit('update:modelValue', [...(props.modelValue ?? []), ''])
}

function removeRow(index: number) {
  const next = [...(props.modelValue ?? [])]
  next.splice(index, 1)
  emit('update:modelValue', next.length ? next : [''])
}
</script>

<template>
  <div class="space-y-2">
    <div v-for="(item, index) in modelValue" :key="index" class="flex gap-2">
      <input
        class="w-full border border-[var(--line)] bg-black/20 px-3 py-2 text-sm outline-none focus:border-[var(--accent)]"
        :value="item"
        @input="update(index, ($event.target as HTMLInputElement).value)"
      />
      <button type="button" class="border border-[var(--line)] px-3 text-sm text-[var(--text-muted)]" @click="removeRow(index)">删</button>
    </div>
    <button type="button" class="border border-[var(--line-strong)] px-3 py-2 text-sm" @click="addRow">添加一项</button>
  </div>
</template>
