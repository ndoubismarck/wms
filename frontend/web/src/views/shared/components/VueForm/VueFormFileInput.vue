<script setup lang="ts">
import { useApp } from '@/app/app'
import { inject, onBeforeUnmount, onMounted, type Ref, ref, shallowRef, toRef, watch } from 'vue'
import { ErrorMessage as VeeErrorMessage, Field as VeeField } from 'vee-validate'
import Uppy, { type Meta, type PluginTarget, type UppyFile } from '@uppy/core'
import UppyXHR from '@uppy/xhr-upload'
import UppyWebcam from '@uppy/webcam'
import UppyDashboard from '@uppy/dashboard'
import UppyImageEditor from '@uppy/image-editor'
//import UppyScreenCapture from '@uppy/screen-capture'

const props = defineProps<{
  name: string
  rules?: string
  label: string
  value?: string
  disabled?: boolean
  maxFileSize?: number
  maxNumberOfFiles?: number
  allowedFileTypes: string[]
}>()

const app = useApp()

const reset = inject<Ref<number>>('formReset')!

const label = toRef(props, 'label')
const disabled = toRef(props, 'disabled')
const maxFileSize = toRef(props, 'maxFileSize', 2e7)
const maxNumberOfFiles = toRef(props, 'maxNumberOfFiles', 10)
const allowedFileTypes = toRef(props, 'allowedFileTypes', [])

const model = ref<string[] | undefined>()
const uppy = shallowRef<Uppy | undefined>()
const uppyDashboardElementRef = ref<PluginTarget<any, any> | undefined>()

watch(reset, () => {
  const files = uppy.value?.getFiles()
  files?.forEach((file: any) => {
    uppy.value?.removeFile(file.id)
  })
})

const addFile = (path: string) => {
  if (!model.value) {
    model.value = []
  }
  model.value.push(path)
}

const removeFile = (path: string) => {
  if (model.value) {
    model.value = model.value.filter((item: string) => {
      return item !== path
    })
  }
}

onMounted(() => {
  uppy.value = new Uppy({
    id: 'Uppy',
    autoProceed: true,
    restrictions: {
      maxFileSize: maxFileSize.value,
      maxNumberOfFiles: maxNumberOfFiles.value,
      allowedFileTypes: allowedFileTypes.value,
    },
  })
  uppy.value.use(UppyDashboard, {
    note: label.value,
    target: uppyDashboardElementRef.value,
    inline: true,
    disabled: disabled.value,
    proudlyDisplayPoweredByUppy: false,
    showRemoveButtonAfterComplete: true,
  })

  uppy.value.use(UppyWebcam, {
    modes: ['picture'],
    mirror: false,
  })

  uppy.value.use(UppyImageEditor)

  //uppy.value.use(ScreenCapture)

  uppy.value.use(UppyXHR, {
    endpoint: app.api.private.url('/files/upload'),
    fieldName: 'file',
    onBeforeRequest: async (xhr: XMLHttpRequest) => {
      const headers = await app.api.private.getHeaders(true)
      for (const key in headers) {
        if (Object.hasOwnProperty.call(headers, key)) {
          xhr.setRequestHeader(key, headers[key] as string)
        }
      }
    },
    onAfterResponse: async (xhr: XMLHttpRequest, _: number) => {
      const response = JSON.parse(xhr.responseText)
      if (response.code === 'success') {
        if (response.data) {
          const fileInfo = response.data['file']
          if (fileInfo && fileInfo['path']) {
            addFile(fileInfo['path'])
          }
        }
      }
    },
  })

  uppy.value.on('file-removed', (file: UppyFile<Meta, Record<string, never>>) => {
    if (file.response && file.response.body && file.response.body.data) {
      const fileInfo = file.response.body.data['file']
      if (fileInfo && fileInfo['path']) {
        removeFile(fileInfo['path'])
      }
    }
  })
})

onBeforeUnmount(() => {
  uppy.value?.destroy()
})
</script>

<template>
  <div class="w-100">
    <vee-field ref="fields" v-slot="{ field, meta }" v-model="model" :name="name">
      <div ref="uppyDashboardElementRef" class="file-upload" :class="{ 'is-invalid': meta.touched && !meta.valid }" />
    </vee-field>
    <vee-field v-for="(field, idx) in model" :key="idx" :name="`${name}[${idx}]`" :value="field" type="text" />
    <vee-error-message v-slot="{ message }" :name="name">
      <p class="invalid-feedback d-block p-0 m-0">
        {{ app.helpers.veeValidate.errorMessage(message) }}
      </p>
    </vee-error-message>
  </div>
</template>
