export type ApiResult<T> = {
  code: number
  message: string
  data: T
}

export type PageResult<T> = {
  total: number
  list: T[]
}

export type AuthState = {
  initialized: boolean
  authenticated: boolean
}

export type OptionItem = {
  id: string
  label: string
  extra?: string
}

export type SelectOption = {
  label: string
  value: string
}

export type KeyValue = {
  key: string
  value: string
}

export type FieldType =
  | 'text'
  | 'textarea'
  | 'number'
  | 'bool'
  | 'string-array'
  | 'kv-array'
  | 'select'
  | 'remote-select'
  | 'remote-multi'
  | 'datetime'

export type ResourceField = {
  key: string
  label: string
  type: FieldType
  placeholder?: string
  options?: SelectOption[]
  optionEndpoint?: string
}

export type ResourceConfig = {
  key: string
  title: string
  section: 'nav' | 'game'
  listEndpoint: string
  detailEndpoint: string
  columns: Array<{ key: string; label: string }>
  fields: ResourceField[]
  defaults: Record<string, unknown>
  bulkReplace?: {
    endpoint: string
    ownerField: ResourceField
    targetField: ResourceField
  }
}
