<script setup lang="ts">
import { onMounted, ref, toRef, watch } from 'vue'
import { ErrorMessage as VeeErrorMessage, Field as VeeField } from 'vee-validate'
import { useApp } from '@/app/app'
import Quill from 'quill'

const props = defineProps<{
  name: string
  rules?: string
  label: string
  value?: string
  disabled: boolean
  placeholder?: string
}>()

const app = useApp()

const model = ref<string>('')
const quill = ref<Quill | null>(null)
const editorRef = ref<HTMLDivElement | null>(null)

const name = toRef(props, 'name')
const label = toRef(props, 'label')
const rules = toRef(props, 'rules')
const value = toRef(props, 'value')
const placeholder = toRef(props, 'placeholder')

watch(model, () => {
  if (model.value.length == 0) {
    if (quill.value) {
      quill.value.root.innerHTML = ''
    }
  }
})

watch(editorRef, () => {
  if (editorRef.value) {
    quill.value = new Quill(editorRef.value, {
      theme: 'snow',
      modules: {
        toolbar: true,
      },
      placeholder: placeholder.value,
    })
    quill.value.on('text-change', () => {
      if (quill.value) {
        let innerHTML = quill.value.root.innerHTML
        //todo remove empty tags and any unwanted tags like <img>, <iframe>, <embed>, <script>..
        model.value = innerHTML
      }
    })
    if (value.value && quill.value) {
      quill.value.root.innerHTML = value.value
    }
  }
})

onMounted(async () => {
  if (value.value) {
    model.value = value.value
  }
})
</script>

<template>
  <div class="w-100">
    <vee-field v-slot="{ field, meta }" :name="name" :rules="rules" v-model="model">
      <textarea v-bind="field" class="d-none" />
      <div class="editor" :class="{ 'is-invalid ': !meta.valid && (meta.touched || meta.validated) }">
        <div ref="editorRef" />
      </div>
    </vee-field>
    <vee-error-message v-slot="{ message }" :name="name">
      <p class="invalid-feedback d-block p-0 m-0">
        {{ app.helpers.veeValidate.errorMessage(message) }}
      </p>
    </vee-error-message>
  </div>
</template>
<style>
.editor .ql-container {
  min-height: 200px;
}

.editor .ql-toolbar {
  font-size: 0.9375rem !important;
  font-family: var(--bs-body-font-family) !important;
  border-top-left-radius: var(--bs-border-radius) !important;
  border-top-right-radius: var(--bs-border-radius) !important;
}

.editor .ql-container {
  font-family: var(--bs-body-font-family) !important;
  font-size: 0.9375rem !important;
  border-bottom-left-radius: var(--bs-border-radius) !important;
  border-bottom-right-radius: var(--bs-border-radius) !important;
}

.editor.is-invalid .ql-toolbar,
.editor.is-invalid .ql-container {
  border-color: var(--bs-form-invalid-border-color) !important;
}
</style>
