import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import type { Ref } from 'vue'

export const useLangStore = defineStore('lang', () => {
    const { locale, setLocale } = useI18n()
    const lang: Ref<'zh' | 'en'> = ref(locale.value === 'en' ? 'en' : 'zh')

    if (import.meta.client) {
        const savedLang = localStorage.getItem('lang') as 'zh' | 'en' | null
        if (savedLang === 'zh' || savedLang === 'en') {
            lang.value = savedLang
            setLocale(savedLang)
        }
    }

    function setLang(newLang: 'zh' | 'en') {
        lang.value = newLang
        setLocale(newLang)
        if (import.meta.client) {
            localStorage.setItem('lang', newLang)
        }
    }

    watch(lang, (newLang) => {
        setLocale(newLang)
        if (import.meta.client) {
            localStorage.setItem('lang', newLang)
        }
    })

    return { lang, setLang }
})
