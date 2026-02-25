<script setup lang="ts">
import { computed, getCurrentInstance, inject, onMounted, ref } from 'vue'

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

function openSelf() {
  accordion.currentIndex.value = index.value
}

onMounted(() => {
  index.value = accordion.register(uid)
})
</script>

<template>
  <div class="accordion-item">
    <h2 class="accordion-header">
      <button
        type="button"
        class="accordion-button"
        :class="{ collapsed: !isOpen }"
        @click="openSelf"
        :aria-expanded="isOpen">
        <span class="badge bg-light text-dark me-3 rounded">{{ index + 1 }}</span>
        <span>{{ title }}</span>
      </button>
    </h2>
    <div class="accordion-collapse collapse" :class="{ show: isOpen }">
      <div class="accordion-body">
        <slot />
      </div>
    </div>
  </div>
</template>

<style scoped>
.accordion-button {
  font-size: 1rem;
}
</style>
