import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  compatibilityDate: '2026-05-01',
  modules: ['@pinia/nuxt', '@nuxtjs/i18n'],
  css: ['~/assets/css/main.css'],
  vite: {
    plugins: [tailwindcss()]
  },
  app: {
    head: {
      htmlAttrs: {
        lang: 'zh-CN'
      },
      title: 'GoFurry',
      meta: [
        { name: 'description', content: 'GoFurry navigation and game discovery site.' },
        { property: 'og:site_name', content: 'GoFurry' },
        { name: 'theme-color', content: '#f97316' }
      ],
      link: [
        { rel: 'icon', href: '/favicon.ico' }
      ]
    }
  },
  runtimeConfig: {
    navApiInternalBase: process.env.NAV_API_INTERNAL_BASE || process.env.NUXT_NAV_API_INTERNAL_BASE || 'http://192.168.153.1:9999/api',
    gameApiInternalBase: process.env.GAME_API_INTERNAL_BASE || process.env.NUXT_GAME_API_INTERNAL_BASE || 'http://192.168.153.1:9998/api',
    public: {
      siteUrl: process.env.NUXT_PUBLIC_SITE_URL || 'http://localhost:3000',
      navApiBase: process.env.NUXT_PUBLIC_NAV_API_BASE || '/api',
      gameApiBase: process.env.NUXT_PUBLIC_GAME_API_BASE || '/api',
      siteLogoPrefixUrl: process.env.NUXT_PUBLIC_SITE_LOGO_PREFIX_URL || 'https://qcdn.go-furry.com/nav/static/SiteLogos/',
      siteDefaultLogo: process.env.NUXT_PUBLIC_SITE_DEFAULT_LOGO || 'https://qcdn.go-furry.com/nav/static/SiteLogos/defaultLogo.svg',
      gameSiteLogoPrefixUrl: process.env.NUXT_PUBLIC_GAME_SITE_LOGO_PREFIX_URL || 'https://qcdn.go-furry.com/game/icons/',
      gamePrefixUrl: process.env.NUXT_PUBLIC_GAME_PREFIX_URL || 'https://qcdn.go-furry.com/game/',
      steamAppPrefixUrl: process.env.NUXT_PUBLIC_STEAM_APP_PREFIX_URL || 'https://store.steampowered.com/app/',
      steamCoverPrefixUrl: process.env.NUXT_PUBLIC_STEAM_COVER_PREFIX_URL || 'https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/'
    }
  },
  routeRules: {
    '/': { prerender: true },
    '/about': { prerender: true },
    '/nav': { ssr: true },
    '/sites': { ssr: true },
    '/sites/**': { ssr: true },
    '/site/**': { redirect: '/sites/**' },
    '/games': { ssr: true },
    '/games/**': { ssr: true },
    '/updates': { ssr: true },
    '/games/news/more': { ssr: true },
    '/games/search': { ssr: false },
    '/games/prize/**': { ssr: false },
    '/user/**': { ssr: false },
    '/settings/**': { ssr: false },
    '/panel': { ssr: false }
  },
  i18n: {
    defaultLocale: 'zh-CN',
    strategy: 'prefix_except_default',
    langDir: 'locales',
    compilation: {
      strictMessage: false,
      escapeHtml: false
    },
    locales: [
      {
        code: 'zh-CN',
        name: '简体中文',
        language: 'zh-CN',
        file: 'zh-CN.json'
      },
      {
        code: 'en',
        name: 'English',
        language: 'en-US',
        file: 'en.json'
      }
    ]
  }
})
