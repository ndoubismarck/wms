<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, shallowRef, watch } from 'vue'
import SimpleBar from 'simplebar'
import type { RouteLocationAsPathGeneric, RouteLocationAsRelativeGeneric } from 'vue-router'

export interface ITableColumn {
  name: string
  class?: string
  value?: {
    text?: string
    html?: string
    link?: ITableValueLink
    image?: ITableValueImage
    dropdown?: ITableValueDropdown
  }
}

export interface ITableValueLink {
  href: string | RouteLocationAsRelativeGeneric | RouteLocationAsPathGeneric
  text: string
  rel?: string
  class?: string
  target?: '_blank'
}

export interface ITableValueImage {
  src: string
  alt?: string
  width?: number
  height?: number
  class?: string
}

export interface ITableValueDropdown {
  button: {
    text: string
    class?: string
  }
  class?: string
  items: {
    text?: string
    link?: ITableValueLink
    image?: ITableValueImage
    class?: string
  }[]
}

const props = defineProps<{
  rows: any[]
  columns: (row: any) => ITableColumn[]
  checkboxes: boolean
}>()

const emits = defineEmits(['checked', 'scrollend'])

const table = ref<HTMLElement | null>(null)
const checked = ref<string[]>([])
const tableColumns = ref<string[] | null>(null)
const tableSimpleBar = ref<SimpleBar | null>(null)
const scrollEndTrigger = ref<HTMLElement>()
const selectAllCheckbox = ref<HTMLInputElement | null>(null)

const observer = shallowRef<IntersectionObserver>()

watch(table, () => {
  if (table.value) {
    if (!tableSimpleBar.value) {
      tableSimpleBar.value = new SimpleBar(table.value)
    }
  }
})

watch(
  checked,
  () => {
    emits('checked', checked.value)
    if (checked.value.length == 0 && selectAllCheckbox.value) {
      selectAllCheckbox.value.checked = false
    }
  },
  {
    deep: true,
  },
)

const loadColumn = (row: any): ITableColumn[] => {
  const tableColumn = props.columns(row)
  if (!tableColumns.value) {
    tableColumns.value = []
    tableColumn.forEach((val: ITableColumn) => {
      tableColumns.value?.push(val.name)
    })
  }
  return tableColumn
}

const handleSelectAllCheckboxChange = (evt: Event) => {
  const target = evt.target as HTMLInputElement
  if (table.value) {
    const tableEl = table.value as HTMLTableElement
    const inputEls = tableEl.querySelectorAll('tr td input')
    inputEls.forEach((val: Element, key: number, parent: NodeListOf<Element>) => {
      const valEl = val as HTMLInputElement
      valEl.checked = target.checked
      valEl.dispatchEvent(new Event('change'))
    })
  }
}

const handleSelectSingleCheckboxChange = (evt: Event) => {
  const target = evt.target as HTMLInputElement
  const index = checked.value.findIndex((value) => `${value}` == `${target.value}`)
  if (target.checked) {
    if (index == -1) {
      checked.value.push(`${target.value}`)
    }
  } else {
    if (index >= 0) {
      checked.value.splice(index, 1)
    }
  }
}

onMounted(() => {
  observer.value = new IntersectionObserver(
    (entries) => {
      if (entries[0]?.isIntersecting) {
        emits('scrollend')
      }
    },
    {
      root: null,
      threshold: 0.1,
    },
  )
  if (scrollEndTrigger.value) {
    observer.value.observe(scrollEndTrigger.value)
  }
})

onBeforeUnmount(() => {
  if (observer.value && scrollEndTrigger.value) {
    observer.value.unobserve(scrollEndTrigger.value)
  }
})
</script>

<template>
  <div ref="table" class="table-responsive sticky-last-column text-nowrap">
    <table class="table">
      <thead>
        <tr>
          <th class="first-column">
            <div class="form-check">
              <input
                ref="selectAllCheckbox"
                @change="handleSelectAllCheckboxChange"
                type="checkbox"
                class="form-check-input" />
            </div>
          </th>
          <th v-for="column in tableColumns">
            {{ column }}
          </th>
          <th class="last-column"></th>
        </tr>
      </thead>
      <tbody class="table-border-bottom-0">
        <tr v-for="(row, idx) in rows" :key="idx">
          <td class="first-column">
            <div class="form-check">
              <input
                @change="handleSelectSingleCheckboxChange"
                type="checkbox"
                name="item"
                :value="row.id"
                class="form-check-input" />
            </div>
          </td>
          <td v-for="(column, idx) in loadColumn(row)" :key="idx">
            <template v-if="column.value?.text">
              {{ column.value.text }}
            </template>
            <template v-if="column.value?.link">
              <router-link :to="column.value.link.href" v-slot="{ href, navigate }" custom>
                <a
                  :href="href"
                  :rel="column.value.link.rel"
                  :class="column.value.link.class"
                  :target="column.value.link.target"
                  @click="navigate">
                  {{ column.value.link.text }}
                </a>
              </router-link>
            </template>
            <template v-else-if="column.value?.html">
              <div v-html="column.value.html" />
            </template>
            <template v-else-if="column.value?.image">
              <img
                :src="column.value.image.src"
                :alt="column.value.image.alt"
                :width="column.value.image.width"
                :height="column.value.image.height"
                :class="column.value.image.class" />
            </template>
            <template v-else-if="column.value?.dropdown">
              <div class="dropdown" :class="column.value.dropdown.class">
                <a
                  href="#"
                  class="btn dropdown-toggle"
                  data-bs-toggle="dropdown"
                  data-bs-boundary="viewport"
                  aria-expanded="false"
                  :class="column.value.dropdown.button.class">
                  {{ column.value.dropdown.button.text }}
                </a>
                <ul class="dropdown-menu">
                  <li v-for="(item, idx) in column.value.dropdown.items" :key="idx">
                    <template v-if="item.text">
                      <span class="dropdown-item">{{ item.text }}</span>
                    </template>
                    <template v-if="item.image">
                      <div class="dropdown-item" :class="item.class">
                        <img
                          :src="item.image.src"
                          :alt="item.image.alt"
                          :width="item.image.width"
                          :height="item.image.height"
                          :class="item.image.class" />
                      </div>
                    </template>
                    <template v-else-if="item.link">
                      <router-link :to="item.link.href" v-slot="{ href, navigate }" custom>
                        <a
                          :href="href"
                          :rel="item.link.rel"
                          :target="item.link.target"
                          @click="navigate"
                          class="dropdown-item"
                          :class="item.class">
                          {{ item.link.text }}
                        </a>
                      </router-link>
                    </template>
                  </li>
                </ul>
              </div>
            </template>
            <template v-else>&nbsp;</template>
          </td>
          <td class="last-column">
            <slot name="actions" :row="row"></slot>
          </td>
        </tr>
      </tbody>
    </table>
    <div ref="scrollEndTrigger"></div>
    <slot name="preloader" />
  </div>
</template>
<style scoped>
.form-check {
  display: block;
  padding-left: 1.8em;
  margin-bottom: 0;
}

.sticky-last-column th.first-column,
.sticky-last-column td.first-column {
  left: 0;
}

.sticky-last-column th.last-column,
.sticky-last-column td.last-column {
  right: 0;
}

.sticky-last-column th.first-column,
.sticky-last-column td.first-column,
.sticky-last-column th.last-column,
.sticky-last-column td.last-column {
  position: sticky;
  background-color: #ffffff;
}

.sticky-last-column th.first-column,
.sticky-last-column th.last-column {
  z-index: 10;
}

.table th .form-check .form-check-input {
  margin-top: 5px;
}

.table th .form-check .form-check-input,
.table tr .form-check .form-check-input {
  width: 18px;
  height: 18px;
}
</style>
