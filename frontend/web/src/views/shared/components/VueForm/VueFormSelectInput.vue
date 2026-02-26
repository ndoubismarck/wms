<script setup lang="ts">
import {computed, onBeforeUnmount, onMounted, ref, toRef, watch} from 'vue'
import {ErrorMessage as VeeErrorMessage, Field as VeeField, type RuleExpression} from 'vee-validate'
import {useApp} from '@/app/app.ts'

export interface ISelectInputOption {
  label: string
  value: any
  selected?: boolean
}

const app = useApp()

const props = defineProps<{
  name: string
  label: string
  value?: number | string
  options: ISelectInputOption[]
  disabled?: boolean
  validation?: RuleExpression<any>
  placeholder: string
}>()

const emits = defineEmits<{
  (e: 'change', option?: ISelectInputOption): void
  (e: 'search', query: string): void
}>()

const value = toRef(props, 'value')
const options = toRef(props, 'options')
const validation = toRef(props, 'validation')

const query = ref('')
const model = ref<any>()
const isOpen = ref(false)
const rootEl = ref<HTMLElement | null>(null)
const inputEl = ref<HTMLInputElement | null>(null)
const emittedEmpty = ref(false)

const selectedOption = computed(() => {
  if (model.value === null || model.value === undefined) {
    return null
  }
  return options.value.find((o) => o.value === model.value) ?? null
})

const visibleOptions = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) {
    return options.value
  }
  const selectedLabel = (selectedOption.value?.label ?? '').trim().toLowerCase()
  if (selectedOption.value && q === selectedLabel) {
    return options.value
  }
  return options.value.filter((o) => (o.label ?? '').toLowerCase().includes(q))
})

const isOptionSelected = (opt: ISelectInputOption) => {
  return model.value !== null && model.value !== undefined && opt.value === model.value
}

watch(
  options,
  () => {
    loadOptions()
  },
  {deep: true},
)

watch(
  selectedOption,
  (opt) => {
    if (!opt) {
      if (model.value !== null && model.value !== undefined) {
        model.value = undefined
        query.value = ''
        emitEmptyOption()
      }
      return
    }
    emittedEmpty.value = false
    query.value = opt.label ?? ''
  },
  {immediate: true},
)

watch(
  value,
  (next) => {
    if (next === null || next === undefined) {
      if (model.value !== null && model.value !== undefined) {
        clearSelection()
      }
      return
    }
    emittedEmpty.value = false
    model.value = next
  },
  {immediate: true},
)

const close = () => {
  isOpen.value = false
}

const loadOptions = () => {
  if (options.value.length == 0) {
    if (model.value !== null && model.value !== undefined) {
      model.value = undefined
      query.value = ''
      emitEmptyOption()
    }
    return
  }

  const selected = options.value.find((item) => item.selected)
  if (selected) {
    model.value = selected.value
    query.value = selected.label ?? ''
    emittedEmpty.value = false
    return
  }

  if (model.value !== null && model.value !== undefined) {
    const index = options.value.findIndex((item) => item.value == model.value)
    if (index < 0) {
      model.value = undefined
      query.value = ''
      emitEmptyOption()
    }
  }
}

const clearSelection = () => {
  const hadSelection = model.value !== null && model.value !== undefined
  const hadQuery = query.value.trim().length > 0
  model.value = undefined
  query.value = ''
  if (hadSelection || hadQuery) {
    emitEmptyOption()
  }
}

const emitEmptyOption = () => {
  if (emittedEmpty.value) {
    return
  }
  emittedEmpty.value = true
  emits('change', undefined)
}

const emitSelectedOption = (opt: ISelectInputOption) => {
  emittedEmpty.value = false
  emits('change', opt)
}

const handleClose = (evt: Event) => {
  evt.preventDefault()
  close()
}

const handleFocus = (evt: Event, touchedCallback: (isTouched: boolean) => void) => {
  evt.preventDefault()
  isOpen.value = true
  touchedCallback(true)
}

const handleInputClick = (touchedCallback: (isTouched: boolean) => void) => {
  isOpen.value = true
  touchedCallback(true)
}

const handleOnBlur = (evt: Event, touchedCallback: (isTouched: boolean) => void) => {
  evt.preventDefault()
  touchedCallback(true)
}

const handleSelect = (evt: Event, opt: ISelectInputOption, touchedCallback: (isTouched: boolean) => void) => {
  evt.preventDefault()
  touchedCallback(true)
  query.value = opt.label ?? ''
  model.value = opt.value
  emitSelectedOption(opt)
  close()
}

const handleOnInput = (evt: Event, touchedCallback: (isTouched: boolean) => void) => {
  evt.preventDefault()
  touchedCallback(true)
  isOpen.value = true
  const target = evt.target as HTMLInputElement
  if (!target) {
    return
  }

  const searchValue = target.value.trim()
  if (!searchValue) {
    clearSelection()
    return
  }

  if (model.value !== null && model.value !== undefined) {
    model.value = undefined
    emitEmptyOption()
  }

  if (searchValue.length >= 3) {
    emits('search', searchValue)
  }
}

const handleDocMouseDown = (e: MouseEvent) => {
  const target = e.target as Node | null
  if (!target) {
    return
  }
  if (!rootEl.value?.contains(target)) {
    close()
  }
}

const handleFocusFirstVisible = (evt: Event) => {
  evt.preventDefault()
  if (!isOpen.value) {
    isOpen.value = true
  }
  const first = rootEl.value?.querySelector<HTMLButtonElement>('.dropdown-item')
  first?.focus()
}

onMounted(() => {
  document.addEventListener('mousedown', handleDocMouseDown)
})

onBeforeUnmount(() => {
  document.removeEventListener('mousedown', handleDocMouseDown)
})

</script>

<template>
  <div ref="rootEl" class="w-100">
    <vee-field v-slot="{ field, meta, setTouched }" :name="name" :rules="validation"
               v-model="model">
      <div class="position-relative d-flex w-100">
        <slot name="addon-start"/>
        <div class="w-100">
          <div class="form-control-fi w-100">
            <input
              ref="inputEl"
              type="text"
              class="form-control"
              :disabled="disabled"
              :placeholder="placeholder"
              v-model="query"
              :class="{ 'is-invalid': !disabled && options.length >0 &&  meta.touched && meta.validated && !meta.valid }"
              @blur="handleOnBlur($event, setTouched)"
              @focus="handleFocus($event, setTouched)"
              @click="handleInputClick(setTouched)"
              @input="handleOnInput($event, setTouched)"
              @keydown.down.prevent="handleFocusFirstVisible"
              @keydown.esc.prevent="handleClose($event)"
            />
            <label>{{ label }}</label>
            <input v-bind="field" type="hidden"/>
          </div>
          <ul class="dropdown-menu w-100" :class="{ show: isOpen }"
              style="max-height: 260px; overflow: auto">
            <li v-if="visibleOptions.length === 0">
              <span class="dropdown-item-text text-muted">No results</span>
            </li>

            <li v-for="opt in visibleOptions" :key="opt.value">
              <button type="button" class="dropdown-item" :class="{ active: isOptionSelected(opt) }"
                      @mousedown.prevent="handleSelect($event, opt, setTouched)">
                {{ opt.label }}
              </button>
            </li>
          </ul>
        </div>
        <slot name="addon-end"/>
      </div>
      <vee-error-message v-slot="{ message }" :name="name">
      <span v-if="!disabled && options.length > 0" class="invalid-feedback d-block p-0 m-0">
        {{ app.helpers.veeValidate.errorMessage(message) }}
      </span>
      </vee-error-message>
    </vee-field>
  </div>
</template>
