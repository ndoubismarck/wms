<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, toRef, watch } from 'vue'
import { ErrorMessage as VeeErrorMessage, Field as VeeField, type RuleExpression } from 'vee-validate'
import { useApp } from '@/app/app.ts'
import type { ISelectInputOption } from '@/views/shared/components/VueForm/VueFormSelectInput.vue'

const app = useApp()

const props = defineProps<{
  name: string
  label: string
  value?: any[]
  options: ISelectInputOption[]
  disabled?: boolean
  validation?: RuleExpression<any>
  placeholder: string
}>()

const emits = defineEmits<{
  (e: 'change', options: ISelectInputOption[]): void
  (e: 'search', query: string): void
}>()

const value = toRef(props, 'value')
const options = toRef(props, 'options')
const disabled = toRef(props, 'disabled', false)
const validation = toRef(props, 'validation')

const model = ref<any[]>([])
const query = ref('')
const isOpen = ref(false)
const rootEl = ref<HTMLElement | null>(null)
const inputEl = ref<HTMLInputElement | null>(null)

const sameValue = (left: any, right: any): boolean => {
  return `${left}` === `${right}`
}

const normalizeModel = (values: any[]): any[] => {
  const result: any[] = []
  const index: Record<string, boolean> = {}

  values.forEach((raw) => {
    const key = `${raw}`.trim()
    if (!key || index[key]) {
      return
    }
    const option = options.value.find((item) => sameValue(item.value, raw))
    if (!option) {
      return
    }
    index[key] = true
    result.push(option.value)
  })

  return result
}

const selectedOptions = computed<ISelectInputOption[]>(() => {
  return options.value.filter((opt) => model.value.some((val) => sameValue(val, opt.value)))
})

const hasValue = computed<boolean>(() => {
  return selectedOptions.value.length > 0 || query.value.trim().length > 0
})

const isFloating = computed<boolean>(() => {
  return isOpen.value && hasValue.value
})

const visibleOptions = computed<ISelectInputOption[]>(() => {
  const search = query.value.trim().toLowerCase()
  if (!search) {
    return options.value
  }
  return options.value.filter((opt) => (opt.label ?? '').toLowerCase().includes(search))
})

const isOptionSelected = (option: ISelectInputOption): boolean => {
  return model.value.some((val) => sameValue(val, option.value))
}

const emitSelected = () => {
  emits('change', selectedOptions.value)
}

const close = () => {
  isOpen.value = false
}

const loadOptions = () => {
  if (options.value.length == 0) {
    if (model.value.length > 0) {
      model.value = []
      emitSelected()
    }
    return
  }

  if (model.value.length == 0) {
    const selected = options.value.filter((item) => item.selected).map((item) => item.value)
    if (selected.length > 0) {
      model.value = normalizeModel(selected)
      emitSelected()
      return
    }
  }

  const normalized = normalizeModel(model.value)
  if (normalized.length != model.value.length) {
    model.value = normalized
    emitSelected()
  }
}

watch(
  options,
  () => {
    loadOptions()
  },
  { deep: true },
)

watch(
  value,
  (next) => {
    if (!Array.isArray(next)) {
      if (model.value.length > 0) {
        model.value = []
        emitSelected()
      }
      return
    }

    model.value = normalizeModel(next)
  },
  { immediate: true },
)

const handleInput = (evt: Event, touchedCallback: (isTouched: boolean) => void) => {
  evt.preventDefault()
  touchedCallback(true)
  isOpen.value = true

  const target = evt.target as HTMLInputElement
  if (!target) {
    return
  }

  const searchValue = target.value.trim()
  if (searchValue.length >= 3) {
    emits('search', searchValue)
  }
}

const handleFocus = (evt: Event, touchedCallback: (isTouched: boolean) => void) => {
  evt.preventDefault()
  touchedCallback(true)
  isOpen.value = true
}

const handleControlClick = (touchedCallback: (isTouched: boolean) => void) => {
  touchedCallback(true)
  if (disabled.value) {
    return
  }
  isOpen.value = true
  inputEl.value?.focus()
}

const handleToggleOption = (evt: Event, option: ISelectInputOption, touchedCallback: (isTouched: boolean) => void) => {
  evt.preventDefault()
  touchedCallback(true)

  if (isOptionSelected(option)) {
    model.value = model.value.filter((val) => !sameValue(val, option.value))
  } else {
    model.value = [...model.value, option.value]
    model.value = normalizeModel(model.value)
  }

  emitSelected()
}

const handleRemoveTag = (evt: Event, option: ISelectInputOption, touchedCallback: (isTouched: boolean) => void) => {
  evt.preventDefault()
  touchedCallback(true)
  model.value = model.value.filter((val) => !sameValue(val, option.value))
  emitSelected()
  inputEl.value?.focus()
}

const handleDocMouseDown = (evt: MouseEvent) => {
  const target = evt.target as Node | null
  if (!target) {
    return
  }
  if (!rootEl.value?.contains(target)) {
    close()
  }
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
    <vee-field v-slot="{ field, meta, setTouched }" :name="name" :rules="validation" v-model="model">
      <div class="position-relative">
        <div class="form-control-fi w-100">
          <div
            class="form-control multi-select-control"
            :class="{
              'is-invalid': meta.touched && meta.validated && !meta.valid,
              'bg-light': disabled,
              'has-value': hasValue,
              'is-floating': isFloating,
            }"
            @click="handleControlClick(setTouched)">
            <div class="d-flex flex-wrap align-items-center gap-1">
              <span
                v-for="option in selectedOptions"
                :key="`${name}-${option.value}`"
                class="badge bg-primary-subtle text-primary d-inline-flex align-items-center multi-select-tag">
                <span class="me-1">{{ option.label }}</span>
                <button
                  type="button"
                  class="btn-close multi-select-tag-close"
                  aria-label="Remove"
                  :disabled="disabled"
                  @click.stop="handleRemoveTag($event, option, setTouched)"></button>
              </span>

              <input
                ref="inputEl"
                type="text"
                class="multi-select-search"
                :disabled="disabled"
                :placeholder="selectedOptions.length == 0 ? placeholder : ''"
                v-model="query"
                @focus="handleFocus($event, setTouched)"
                @input="handleInput($event, setTouched)"
                @keydown.esc.prevent="close" />
            </div>
          </div>
          <label>{{ label }}</label>
          <input v-bind="field" type="hidden" />
        </div>

        <ul class="dropdown-menu w-100" :class="{ show: isOpen }" style="max-height: 260px; overflow: auto">
          <li v-if="visibleOptions.length == 0">
            <span class="dropdown-item-text text-muted">No results</span>
          </li>
          <li v-for="option in visibleOptions" :key="`${name}-option-${option.value}`">
            <button
              type="button"
              class="dropdown-item d-flex justify-content-between align-items-center"
              :class="{ active: isOptionSelected(option) }"
              @mousedown.prevent="handleToggleOption($event, option, setTouched)">
              <span>{{ option.label }}</span>
              <i v-if="isOptionSelected(option)" class="bx bx-check"></i>
            </button>
          </li>
        </ul>
      </div>
    </vee-field>

    <vee-error-message v-slot="{ message }" :name="name">
      <span class="invalid-feedback d-block p-0 m-0">
        {{ app.helpers.veeValidate.errorMessage(message) }}
      </span>
    </vee-error-message>
  </div>
</template>

<style scoped>
.multi-select-control {
  min-height: 42px;
  height: auto;
}

.multi-select-control.is-floating + label {
  top: 0;
  transform: translateY(-50%);
  font-size: 0.75rem;
  color: var(--bs-body-color, #212529);
}

.multi-select-control.has-value:not(.is-floating) + label {
  opacity: 0;
}

.multi-select-tag {
  padding: 0.25rem 0.5rem;
}

.multi-select-tag-close {
  width: 0.65rem;
  height: 0.65rem;
  font-size: 0.55rem;
}

.multi-select-search {
  flex: 1 1 140px;
  min-width: 120px;
  border: 0;
  outline: none;
  background: transparent;
  padding: 0.1rem 0.15rem;
}
</style>
