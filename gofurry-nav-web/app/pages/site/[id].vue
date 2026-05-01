<script setup lang="ts">
const route = useRoute()
await navigateTo({
  path: `/sites/${route.params.id}`,
  query: route.query
}, { redirectCode: 301 })
</script>
