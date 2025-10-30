import 'vue3-toastify/dist/index.css'
import Vue3Toastify, { type ToastContainerOptions } from 'vue3-toastify'
import 'vue3-toastify/dist/index.css'
import { toast } from 'vue3-toastify'
import type { App } from 'vue'

export default {
    install: (app: App) => {
        app.use(Vue3Toastify, {
            autoClose: 3000,
            position: toast.POSITION.TOP_CENTER,
            clearOnUrlChange: false,
            // You can add other global options here
        } as ToastContainerOptions)
    },
}
