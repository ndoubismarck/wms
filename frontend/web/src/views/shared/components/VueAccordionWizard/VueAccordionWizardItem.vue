<script setup lang="ts">
import {computed, getCurrentInstance, inject, onBeforeUnmount, onMounted, ref} from 'vue'

const props = defineProps<{
  title: string
}>()

const accordion = inject<any>('accordion')
if (!accordion) {
  throw new Error('VueAccordionItem must be inside VueAccordion')
}

const uid = getCurrentInstance()!.uid
const index = ref(0)

const isOpen = computed(() => accordion.currentIndex.value === index.value)
const transitionName = computed(() =>
  accordion.transitionDirection?.value === 'backward' ? 'wizard-step-backward' : 'wizard-step-forward',
)

onMounted(() => {
  index.value = accordion.register(uid)
})

onBeforeUnmount(() => {
  if (typeof accordion.unregister === 'function') {
    accordion.unregister(uid)
  }
  index.value = 0
})
</script>

<template>
  <div class="accordion-item">
    <div class="accordion-header">
      <div
        class="accordion-button cursor-default"
        :class="{ collapsed: !isOpen }"
        :aria-expanded="isOpen">
        <span class="badge bg-light text-dark me-3 rounded">{{ index + 1 }}</span>
        <span>{{ title }}</span>
      </div>
    </div>
    <div class="accordion-collapse">
      <transition :name="transitionName">
        <div v-show="isOpen" class="accordion-body wizard-step-panel py-3">
          <slot/>
        </div>
      </transition>
    </div>
  </div>
</template>

<style scoped>
.accordion-button {
  font-size: 1rem;
}

.accordion .accordion-button::after {
  display: none;
}

.wizard-step-panel {
  overflow: hidden;
}

.wizard-step-forward-enter-active,
.wizard-step-forward-leave-active,
.wizard-step-backward-enter-active,
.wizard-step-backward-leave-active {
  transition: opacity 0.24s ease, transform 0.24s ease, max-height 0.24s ease;
  will-change: opacity, transform, max-height;
}

.wizard-step-forward-enter-from,
.wizard-step-backward-leave-to {
  opacity: 0;
  transform: translateX(14px);
  max-height: 0;
}

.wizard-step-forward-leave-to,
.wizard-step-backward-enter-from {
  opacity: 0;
  transform: translateX(-14px);
  max-height: 0;
}

.wizard-step-forward-enter-to,
.wizard-step-forward-leave-from,
.wizard-step-backward-enter-to,
.wizard-step-backward-leave-from {
  max-height: 2000px;
}

@media (prefers-reduced-motion: reduce) {
  .wizard-step-forward-enter-active,
  .wizard-step-forward-leave-active,
  .wizard-step-backward-enter-active,
  .wizard-step-backward-leave-active {
    transition: none;
  }
}
</style>
