export function stripHtml(value?: string | null) {
  return (value || '')
    .replace(/<[^>]*>/g, ' ')
    .replace(/\s+/g, ' ')
    .trim()
}

export function truncate(value?: string | null, length = 150) {
  const text = stripHtml(value)
  return text.length > length ? `${text.slice(0, length)}...` : text
}

export function safeJsonParse<T>(value: unknown, fallback: T): T {
  if (typeof value !== 'string') {
    return (value as T) ?? fallback
  }

  try {
    return JSON.parse(value) as T
  } catch {
    return fallback
  }
}

export function displayDate(value?: string | number | null) {
  if (!value) return ''

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return String(value)
  }

  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}
