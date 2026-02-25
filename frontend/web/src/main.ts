import {createApp} from 'vue'
import {createPinia} from 'pinia'

import App from './App.vue'
import router from './router'

import 'simplebar'
import 'bootstrap'

import 'quill/dist/quill.core.css'
import 'quill/dist/quill.snow.css'
import 'quill/dist/quill.bubble.css'

import '@uppy/core/css/style.min.css'
import '@uppy/dashboard/css/style.min.css'
import '@uppy/webcam/css/style.min.css'
import '@uppy/image-editor/css/style.min.css'

import "vue3-select-component/styles";

import 'simplebar/dist/simplebar.css'
import 'bootstrap/dist/css/bootstrap.css'

import './assets/css/app/app.css'
import './assets/css/app/app-icons.css'
import './assets/css/app/app-auth.css'
import './assets/css/app/app-custom.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
