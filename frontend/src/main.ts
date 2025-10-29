import 'vue3-toastify/dist/index.css'
import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import i18n from './plugins/i18n' // Import i18n
import toastify from './plugins/toastify'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(i18n) // Use i18n
app.use(toastify)

app.mount('#app')
