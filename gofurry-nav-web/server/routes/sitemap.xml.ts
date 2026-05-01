type ApiResult<T> = {
  code: number
  data: T
}

type SiteRecord = {
  id: string
  domain: string
}

type GameRecord = {
  game_id: string
}

function escapeXml(value: string) {
  return value
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&apos;')
}

function urlEntry(loc: string) {
  return `<url><loc>${escapeXml(loc)}</loc></url>`
}

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig(event)
  const siteUrl = String(config.public.siteUrl).replace(/\/$/, '')
  const urls = new Set([
    '/',
    '/nav',
    '/games',
    '/updates',
    '/about',
    '/games/news/more',
    '/games/creator'
  ])

  const [sites, games] = await Promise.all([
    $fetch<ApiResult<SiteRecord[]>>('/nav/page/site/list', {
      baseURL: config.navApiInternalBase,
      query: { lang: 'zh' }
    }).then((res) => res.code === 1 ? res.data : []).catch(() => []),
    $fetch<ApiResult<GameRecord[]>>('/game/info/list', {
      baseURL: config.gameApiInternalBase
    }).then((res) => res.code === 1 ? res.data : []).catch(() => [])
  ])

  for (const site of sites) {
    urls.add(`/sites/${site.id}`)
  }

  for (const game of games) {
    urls.add(`/games/${game.game_id}`)
  }

  setHeader(event, 'content-type', 'application/xml; charset=utf-8')
  return [
    '<?xml version="1.0" encoding="UTF-8"?>',
    '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">',
    ...Array.from(urls).map((path) => urlEntry(`${siteUrl}${path}`)),
    '</urlset>'
  ].join('')
})
