<script setup lang="ts">
import mobile from 'is-mobile'
import { onBeforeMount, onBeforeUnmount, onMounted, type Ref, ref, toRef, useSlots, watch } from 'vue'
import { Offcanvas } from 'bootstrap'
import SimpleBar from 'simplebar'

const slots = useSlots()

export interface IOffcanvas {
  show: () => void
  hide: () => void
  element: HTMLElement | Ref<HTMLElement | null> | null
  scrollTop: () => void
  scrollBottom: () => void
}

const emit = defineEmits(['hide', 'show'])

const props = defineProps<{
  show?: boolean
  size?: number
  title?: string
  parent?: IOffcanvas
  padding?: boolean
  disabled?: boolean
  overflow?: boolean
}>()

const parent = toRef(props, 'parent')
const overflow = toRef(props, 'overflow', true)

const size = ref<number>(props.size ? props.size : 0)
const start = ref<number>(0)
const zIndex = ref<number>(1099)
const element = ref<HTMLElement | null>(null)
const showView = ref<boolean>(props.show)
const isMobile = ref<boolean>(false)
const isShowing = ref<boolean>(false)
const offcanvas = ref<Offcanvas | null>(null)
const offcanvasBody = ref<HTMLElement | null>(null)
const offcanvasBodySimpleBar = ref<SimpleBar | null>(null)

const hide = () => {
  if (offcanvas.value) {
    offcanvas.value.hide()
  }
  showView.value = false
}

const show = () => {
  if (offcanvas.value) {
    offcanvas.value.show()
  }
}

const handleHide = () => {
  hide()
}

const scrollTop = () => {
  if (offcanvasBodySimpleBar.value) {
    const scrollElement = offcanvasBodySimpleBar.value.getScrollElement()
    if (scrollElement) {
      scrollElement.scrollTo({ top: 0, behavior: 'smooth' })
    }
  }
}

const scrollBottom = () => {
  if (offcanvasBodySimpleBar.value) {
    const scrollElement = offcanvasBodySimpleBar.value.getScrollElement()
    if (scrollElement) {
      scrollElement.scrollTo({ top: scrollElement.scrollHeight, behavior: 'smooth' })
    }
  }
}

const loadParent = () => {
  if (parent.value) {
    const body = document.body
    const parentElement = parent.value.element as HTMLElement | null
    if (parentElement) {
      size.value = parentElement.offsetWidth - 80
      const parentZIndex = parseInt(parentElement.style.zIndex)
      if (parentZIndex > 0) {
        zIndex.value = parentZIndex + 1
      }
      start.value = body.offsetWidth - parentElement.offsetWidth
    }
  }
}

watch(parent, () => {
  loadParent()
})

watch(offcanvasBody, () => {
  if (offcanvasBody.value) {
    if (!offcanvasBodySimpleBar.value) {
      offcanvasBodySimpleBar.value = new SimpleBar(offcanvasBody.value)
    }
  }
})

onMounted(() => {
  if (element.value) {
    offcanvas.value = new Offcanvas(element.value)
    element.value.addEventListener('shown.bs.offcanvas', () => {
      isShowing.value = true
      emit('show')
    })
    element.value.addEventListener('hidden.bs.offcanvas', () => {
      isShowing.value = false
      emit('hide')
    })
  }
  if (props.show) {
    show()
  }
  loadParent()
})

onBeforeMount(() => {
  isMobile.value = mobile()
  if (!size.value || size.value <= 0) {
    size.value = 800
  }
})

onBeforeUnmount(() => {
  if (offcanvas.value) {
    offcanvas.value.dispose()
  }
  if (offcanvasBodySimpleBar.value) {
    offcanvasBodySimpleBar.value.unMount()
  }
})

defineExpose<IOffcanvas>({
  show,
  hide,
  element,
  scrollTop,
  scrollBottom,
})
</script>

<template>
  <teleport to="body">
    <div
      ref="element"
      class="offcanvas offcanvas-end"
      :style="{ zIndex: zIndex + 1, minWidth: isMobile ? '100%' : `${size}px` }"
      tabindex="-1"
      data-bs-backdrop="false"
      data-bs-keyboard="false"
      aria-labelledby="offcanvas-label">
      <div class="offcanvas-header" :class="{ 'border-bottom': title }">
        <h5 class="offcanvas-title pe-3">
          {{ title }}
        </h5>
        <button
          :disabled="disabled"
          @click="handleHide"
          type="button"
          class="btn-close text-reset"
          aria-label="Close" />
      </div>
      <div ref="offcanvasBody" class="offcanvas-body" :class="{ 'p-3': !padding, 'overflow-y-hidden': !overflow }">
        <div>
          <slot name="body" />
        </div>
        <div :class="{ overlay: slots.overaly }">
          <slot name="overlay" />
        </div>
      </div>
    </div>
    <div v-if="isShowing" class="offcanvas-backdrop fade show" :style="{ zIndex: zIndex, left: `${start}px` }" />
  </teleport>
</template>
<style scoped>
.offcanvas .offcanvas-header {
  height: 65px;
  padding-top: 0;
  padding-bottom: 0;
}

.offcanvas .offcanvas-body {
  position: relative;
  background: var(--bs-body-bg);
  max-height: calc(100vh - 65px);
}

.offcanvas .offcanvas-body .overlay {
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  position: absolute;
  height: calc(100vh - 65px);
}

.offcanvas .btn-close {
  font-size: 20px;
}
</style>
