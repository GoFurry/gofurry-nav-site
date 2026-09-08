import assert from 'node:assert/strict'
import { createServer } from 'node:http'
import { spawn } from 'node:child_process'
import { once } from 'node:events'
import { existsSync } from 'node:fs'
import { mkdtemp } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { setTimeout as delay } from 'node:timers/promises'
import { launchPerfBrowser, rootDir } from './perf/shared.mjs'
import { mockOverview, mockGamePanel } from './fixtures/insights-overview.mjs'

const sourcePaths = {
  nav: '/api/v2/nav/insights/overview',
  game: '/api/v2/game/insights/overview',
  panel: '/api/v2/game/panel/main',
}

// Invoked by insights:smoke -- --overview-fixtures. This deliberately tests only the
// Overview against local fixtures; it does not stand in for the full live-data smoke.
export async function runOverviewSmoke() {
  assert(existsSync(join(rootDir, '.output/server/index.mjs')), 'Build Nav Web before running the production fixture smoke')
  let failure = ''
  let siteHero = false
  const requests = []
  let upstreamUrl = ''
  const upstream = createServer((request, response) => {
    const path = new URL(request.url, 'http://localhost').pathname
    requests.push(path)
    if (path.startsWith('/media/')) {
      // Small local artwork for layout verification only; never shipped as product assets.
      const site = path.includes('site') || path.includes('default')
      response.writeHead(200, { 'Content-Type': 'image/svg+xml' })
      response.end(site
        ? '<svg xmlns="http://www.w3.org/2000/svg" width="72" height="72"><rect width="72" height="72" rx="14" fill="#82654f"/><path d="M16 54V18l20 18 20-18v36" fill="none" stroke="#f4e9dd" stroke-width="5"/></svg>'
        : '<svg xmlns="http://www.w3.org/2000/svg" width="460" height="215" viewBox="0 0 460 215"><rect width="460" height="215" fill="#394b44"/><circle cx="350" cy="62" r="28" fill="#c8b398"/><path d="M0 215V155L110 58l120 157M190 215l125-123 145 110v13" fill="#778978"/><text x="28" y="180" fill="#fff" font-family="sans-serif" font-size="18">GAME · FIXTURE</text></svg>')
      return
    }
    const source = Object.keys(sourcePaths).find(key => sourcePaths[key] === path)
    if (!source || failure === source || failure === 'all') {
      response.writeHead(503, { 'Content-Type': 'application/json' })
      response.end(JSON.stringify({ code: 0, message: 'Local fixture unavailable' }))
      return
    }
    const data = source === 'panel' ? mockGamePanel(`${upstreamUrl}/media`) : mockOverview(source === 'nav' ? 'site' : 'game', `${upstreamUrl}/media`)
    if (siteHero && source === 'nav') data.recent_changes[0].occurred_at = '2026-09-02T12:00:00Z'
    response.writeHead(200, { 'Content-Type': 'application/json' })
    response.end(JSON.stringify({ code: 1, data }))
  })
  await listen(upstream)
  upstreamUrl = `http://127.0.0.1:${upstream.address().port}`
  const reservation = createServer()
  await listen(reservation)
  const port = reservation.address().port
  await new Promise(resolve => reservation.close(resolve))
  const base = `http://127.0.0.1:${port}`
  const environment = {
    ...process.env, NITRO_HOST: '127.0.0.1', NITRO_PORT: String(port),
    NUXT_PUBLIC_NAV_API_BASE: '/api/v1', NUXT_PUBLIC_NAV_V2_API_BASE: '/api/v2',
    NUXT_PUBLIC_GAME_API_BASE: '/api/v1', NUXT_PUBLIC_GAME_V2_API_BASE: '/api/v2',
    NUXT_PUBLIC_SITE_LOGO_PREFIX_URL: `${upstreamUrl}/media/`,
    NUXT_PUBLIC_SITE_DEFAULT_LOGO: `${upstreamUrl}/media/default.svg`,
  }
  for (const name of ['NAV_API', 'NAV_V2_API', 'GAME_API', 'GAME_V2_API']) {
    environment[`NUXT_${name}_INTERNAL_BASE`] = `${upstreamUrl}/api/${name.includes('V2') ? 'v2' : 'v1'}`
  }
  const preview = spawn(process.execPath, ['.output/server/index.mjs'], { cwd: rootDir, env: environment, windowsHide: true, stdio: ['ignore', 'pipe', 'pipe'] })
  let logs = ''
  const record = data => { logs = (logs + data.toString()).slice(-8000) }
  preview.stdout.on('data', record)
  preview.stderr.on('data', record)
  let browser
  try {
    await ready(`${base}/web/background/gofurry-pattern.svg`, preview)
    for (const route of ['/insights', '/en/insights']) {
      requests.length = 0
      const response = await fetch(base + route)
      const html = await response.text()
      assert.equal(response.status, 200, `${route} HTTP status`)
      assert.match(html, /<h1>[^<]+<\/h1>/, 'visible localized H1 must SSR')
      assert.deepEqual(stats(html), ['238', '213', '47'], 'independent stats must SSR')
      for (const section of ['header', 'activity', 'sites', 'games', 'explore']) assert(html.includes(`data-overview-${section}`), `missing ${section}`)
      const prefix = route.startsWith('/en/') ? '/en' : ''
      for (const path of ['/site/41', '/site/42', '/games/82', '/games/83', '/games/91', '/games/92', '/games/93']) assert(html.includes(`href="${prefix}${path}"`), `missing localized ${path}`)
      assert(html.includes(`${upstreamUrl}/media/site.svg`) && html.includes(`${upstreamUrl}/media/default.svg`), 'Site refs/default must SSR')
      assert(html.includes('data-media-state="fallback"'), 'missing Game header must SSR a neutral fallback')
      assert.match(html, /<time datetime="2026-09-01T10:00:00.000Z"/, 'earlier snapshot time must SSR')
      assert.deepEqual(requests.filter(path => Object.values(sourcePaths).includes(path)).sort(), Object.values(sourcePaths).sort(), 'Overview must request exactly three independent sources')
      assert(!requests.some(path => /\/(?:sites|games)\/\d+/.test(path)), 'Overview made a per-entity lookup')
      console.log(`[overview] SSR ${route}: sections, stats, media, localized links and three requests PASS`)
    }
    for (const [source, expected] of [['nav', ['—', '213', '—']], ['game', ['238', '—', '—']], ['panel', ['238', '213', '47']], ['all', ['—', '—', '—']]]) {
      failure = source
      const response = await fetch(`${base}/insights`)
      const html = await response.text()
      assert.equal(response.status, 200, `${source} failure should stay local`)
      assert.deepEqual(stats(html), expected, `${source} failure fabricated a statistic`)
      assert.equal(html.includes('data-pulse="players"'), source !== 'panel' && source !== 'all', 'Panel failure independence')
      assert.equal(html.includes('data-metric="tls13"'), source !== 'nav' && source !== 'all', 'Site failure independence')
      assert.equal(html.includes('data-change-link'), source !== 'all', 'activity failure independence')
      console.log(`[overview] independent ${source} failure PASS`)
    }
    failure = ''
    browser = await launchPerfBrowser()
    const context = await browser.newContext({ locale: 'zh-CN' })
    const page = await context.newPage()
    const pageErrors = []
    page.on('pageerror', error => pageErrors.push(error.message))
    // Keep fixture runs local even if the shell loads optional third-party resources.
    await page.route('**/*', route => {
      const url = new URL(route.request().url())
      return ['127.0.0.1', 'localhost'].includes(url.hostname) || ['data:', 'blob:'].includes(url.protocol) ? route.continue() : route.abort()
    })
    const screenshots = await mkdtemp(join(tmpdir(), 'gofurry-overview-b2-'))
    for (const width of [1440, 1024, 390]) {
      await page.setViewportSize({ width, height: 1000 })
      for (const prefix of ['', '/en']) {
        requests.length = 0
        await page.goto(`${base}${prefix}/insights`, { waitUntil: 'networkidle' })
        assert.deepEqual(requests.filter(path => Object.values(sourcePaths).includes(path)).sort(), Object.values(sourcePaths).sort(), 'hydration repeated Overview requests')
        assert(!requests.some(path => /\/(?:sites|games)\/\d+/.test(path)), 'browser added a per-entity request')
        assert(await page.locator('main h1').isVisible(), 'H1 is not visible')
        assert.equal(await page.locator('[data-overview-activity] [data-change-link]').count(), 4, 'fixture activity did not render')
        assert.equal(await page.locator('[data-pulse="players"] strong').textContent(), 'Active game fixture')
        assert((await page.locator('[data-pulse="players"]').textContent()).includes('2,800'), 'available Panel peak was dropped')
        assert((await page.locator('[data-pulse="discount"]').textContent()).includes('5.99'), 'minor-unit USD price was not preserved')
        const progress = await page.locator('progress').evaluateAll(elements => elements.map(el => ({ value: el.value, max: el.max, label: el.getAttribute('aria-labelledby') })))
        assert(progress.every(item => item.value >= 0 && item.value <= 1 && item.max === 1 && item.label), 'progress contract')
        await assertLayout(page, width)
        await revealImages(page)
        await page.screenshot({ path: join(screenshots, `${prefix ? 'en' : 'zh'}-${width}.png`), fullPage: true })
        await page.locator('.insights-primary-nav a').last().focus()
        assert(await page.locator('.insights-primary-nav a').last().evaluate(el => getComputedStyle(el).outlineStyle !== 'none'), 'navigation keyboard focus missing')
        console.log(`[overview] ${prefix || 'zh'} ${width}px layout, media and keyboard focus PASS`)
      }
    }
    siteHero = true
    await page.goto(`${base}/insights`, { waitUntil: 'networkidle' })
    assert.equal(await page.locator('.insight-activity-item--hero').getAttribute('data-domain'), 'site')
    assert((await page.locator('.insight-activity-item--hero img').boundingBox()).width <= 72, 'Site icon was stretched into a cover')
    await page.screenshot({ path: join(screenshots, 'site-hero-mobile.png'), fullPage: true })
    await page.route('**/media/**', route => route.abort())
    await page.reload({ waitUntil: 'networkidle' })
    await revealImages(page)
    assert.equal(await page.locator('main .insight-entity-media img').count(), 0, 'failed Site/default/Game images did not reach fallback')
    assert(await page.locator('main [role="img"][aria-label]').count() > 0, 'fallback lost its accessible entity name')
    await assertLayout(page, 390)
    await page.screenshot({ path: join(screenshots, 'media-fallback-mobile.png'), fullPage: true })
    await page.getByRole('button', { name: '切换明暗主题图标', exact: true }).click()
    assert(await page.locator('html').evaluate(el => el.classList.contains('dark')), 'existing theme control did not enable dark mode')
    // The existing public background is fixed to the viewport; inspect dark mode
    // in viewports so a full-page capture does not extend beyond that surface.
    await page.screenshot({ path: join(screenshots, 'dark-fallback-mobile.png') })
    await page.locator('[data-overview-games]').scrollIntoViewIfNeeded()
    await page.screenshot({ path: join(screenshots, 'dark-games-mobile.png') })
    assert.deepEqual(pageErrors, [], 'Overview threw a browser exception')
    console.log('[overview] Site hero and image error fallback PASS')
    console.log(`[overview] screenshots (temporary, outside Git): ${screenshots}`)
  } catch (error) {
    console.error(logs)
    throw error
  } finally {
    await browser?.close()
    if (preview.exitCode === null) { preview.kill(); await once(preview, 'exit') }
    await new Promise(resolve => upstream.close(resolve))
  }
}

function stats(html) {
  const dl = html.match(/<dl class="overview-stats">([\s\S]*?)<\/dl>/)?.[1] || ''
  return [...dl.matchAll(/<dd>(.*?)<\/dd>/g)].map(match => match[1])
}
async function listen(server) {
  server.listen(0, '127.0.0.1')
  await once(server, 'listening')
}
async function ready(url, child) {
  for (let attempt = 0; attempt < 100; attempt++) {
    assert(child.exitCode === null, 'production server exited during startup')
    try { if ((await fetch(url)).ok) return } catch {}
    await delay(200)
  }
  throw new Error('production server did not become ready')
}
async function assertLayout(page, width) {
  const layout = await page.evaluate(() => {
    const site = document.querySelector('[data-overview-sites]').getBoundingClientRect()
    const game = document.querySelector('[data-overview-games]').getBoundingClientRect()
    return { overflow: document.documentElement.scrollWidth > window.innerWidth, site: { top: site.top, width: site.width }, game: { top: game.top, width: game.width } }
  })
  assert(!layout.overflow, `${width}px page has horizontal overflow`)
  if (width >= 1200) assert(layout.game.width > layout.site.width * 1.3 && Math.abs(layout.site.top - layout.game.top) < 2, 'desktop lost asymmetric split')
  if (width === 390) assert(layout.game.top > layout.site.top, 'mobile did not stack domains')
}
async function revealImages(page) {
  await page.evaluate(async () => {
    for (let top = 0; top < document.body.scrollHeight; top += 650) {
      window.scrollTo(0, top)
      await new Promise(resolve => setTimeout(resolve, 60))
    }
    window.scrollTo(0, 0)
  })
  await page.waitForTimeout(200)
}
