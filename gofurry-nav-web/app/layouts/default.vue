<template>
  <div class="flex min-h-screen flex-col bg-slate-50 text-slate-950">
    <AppHeader />
    <main class="min-w-0 flex-1">
      <slot />
    </main>
    <AppFooter />
  </div>
</template>
