<script setup lang="ts">
import {computed, onBeforeUnmount, onMounted, ref, toRef, watch} from 'vue'
import VueSelect from 'vue3-select-component'
import {ErrorMessage as VeeErrorMessage, Field as VeeField, type RuleExpression} from 'vee-validate'
import {useApp} from '@/app/app'

export interface ISelectInputOption {
  label: string
  value: any
  selected?: boolean
}

const props = defineProps<{
  name: string
  label?: string
  value?: number | string
  isMulti?: boolean
  options: ISelectInputOption[]
  isLoading?: boolean
  disabled?: boolean
  validation?: RuleExpression<any>
  placeholder?: string
}>()

const app = useApp()
const emits = defineEmits(['change', 'search'])

const model = ref<any>()
const fields = ref<(typeof VeeField)[]>()

const value = toRef(props, 'value')
const options = toRef(props, 'options')
const selectOptions = computed(() => options.value as any[])

const loadOptions = () => {
  options.value.forEach((item) => {
    if (item.selected) {
      model.value = item.value
    }
  })

  const index = options.value.findIndex((item) => item.value == model.value)
  if (index < 0) {
    model.value = undefined
  }
}

const emitSelectedOption = () => {
  const selectedOption = options.value.find((item) => item && item.value == model.value)
  emits('change', selectedOption)
  if (!selectedOption) {
    model.value = null
  }
}

const emitSearch = (value: string) => {
  emits('search', value)
}

watch(
  options,
  () => {
    if (options.value.length == 0) {
      model.value = undefined
      return
    }
    loadOptions()
  },
  {deep: true},
)

watch(model, () => {
  emitSelectedOption()
})

onMounted(() => {
  loadOptions()
  if (value.value != null) {
    model.value = value.value
  }
})

onBeforeUnmount(() => {
  fields.value = []
})
</script>

<template>
  <div class="w-100">
    <div class="form-floating d-flex">
      <vee-field ref="fields" v-slot="{ field }" v-model="model" :name="name" :rules="validation">
        <slot name="addon-start"/>
        <vue-select
          v-bind="field"
          v-model="model"
          :is-multi="isMulti"
          :options="selectOptions"
          :is-loading="isLoading"
          :is-disabled="disabled"
          @search="emitSearch"
          :classes="{control:'form-control'}"
          :placeholder="placeholder ? placeholder : 'Select an option'"/>
        <slot name="addon-end"/>
      </vee-field>
    </div>
    <vee-error-message v-slot="{ message }" :name="name">
      <p class="invalid-feedback d-block p-0 m-0">
        {{ app.helpers.veeValidate.errorMessage(message) }}
      </p>
    </vee-error-message>
  </div>
</template>
<style>
:root {
  --vs-width: 100%;
  --vs-min-height: 58px;
  --vs-padding: 10px 8px;
  --vs-border: 1px solid #e4e4e7;
  --vs-border-radius: 4px;
  --vs-font-size: 16px;
  --vs-font-weight: 400;
  --vs-font-family: inherit;
  --vs-text-color: #18181b;
  --vs-line-height: 1.5;
  --vs-placeholder-color: #52525b;
  --vs-background-color: #fff;
  --vs-disabled-background-color: #f4f4f5;
  --vs-outline-width: 1px;
  --vs-outline-color: v-bind(--bs-primary);

  --vs-menu-offset-top: 8px;
  --vs-menu-height: 200px;
  --vs-menu-border: var(--vs-border);
  --vs-menu-background-color: var(--vs-background-color);
  --vs-menu-box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.05);
  --vs-menu-z-index: 2;

  --vs-option-width: 100%;
  --vs-option-padding: 8px 12px;
  --vs-option-cursor: pointer;
  --vs-option-font-size: var(--vs-font-size);
  --vs-option-font-weight: var(--vs-font-weight);
  --vs-option-text-align: -webkit-auto;
  --vs-option-text-color: var(--vs-text-color);
  --vs-option-hover-text-color: var(--vs-text-color);
  --vs-option-focused-text-color: var(--vs-text-color);
  --vs-option-selected-text-color: var(--vs-text-color);
  --vs-option-disabled-text-color: #52525b;
  --vs-option-background-color: var(--vs-menu-background);
  --vs-option-hover-background-color: #dbeafe;
  --vs-option-focused-background-color: #dbeafe;
  --vs-option-selected-background-color: #93c5fd;
  --vs-option-disabled-background-color: #f4f4f5;
  --vs-option-opacity-menu-open: 0.4;

  --vs-multi-value-margin: 2px;
  --vs-multi-value-border: 0px;
  --vs-multi-value-border-radius: 2px;
  --vs-multi-value-background-color: #f4f4f5;

  --vs-multi-value-label-padding: 4px 4px 4px 8px;
  --vs-multi-value-label-font-size: 12px;
  --vs-multi-value-label-font-weight: 400;
  --vs-multi-value-label-line-height: 1;
  --vs-multi-value-label-text-color: #3f3f46;

  --vs-multi-value-delete-padding: 0 3px;
  --vs-multi-value-delete-hover-background-color: #FF6467;
  --vs-multi-value-xmark-size: 16px;
  --vs-multi-value-xmark-cursor: pointer;
  --vs-multi-value-xmark-color: var(--vs-multi-value-label-text-color);
  --vs-multi-value-xmark-hover-color: #fff;

  --vs-indicators-gap: 0px;
  --vs-indicator-icon-size: 20px;
  --vs-indicator-icon-color: var(--vs-text-color);
  --vs-indicator-icon-cursor: pointer;
  --vs-indicator-dropdown-icon-transition: transform 0.2s ease-out;

  --vs-spinner-color: var(--vs-text-color);
  --vs-spinner-size: 16px;
}


</style>
