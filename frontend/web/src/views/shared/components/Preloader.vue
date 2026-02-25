<script setup lang="ts">
import { toRef } from 'vue'
import { EPreloaderSize } from '@/app/types'

const props = defineProps<{
  size?: EPreloaderSize
  overlay?: boolean
}>()

const size = toRef(props, 'size', EPreloaderSize.SM)
const overlay = toRef(props, 'overlay', false)
</script>

<template>
  <template v-if="overlay">
    <div
      class="preloader position-absolute top-0 start-0 d-flex align-items-center justify-content-center text-dark overlay">
      <div class="loader" :class="size" />
    </div>
  </template>
  <template v-else>
    <div
      class="preloader position-absolute top-0 start-0 w-100 h-100 d-flex align-items-center justify-content-center bg-white text-dark">
      <div class="loader" :class="size" />
    </div>
  </template>
</template>
<style scoped>
.overlay {
  z-index: 99;
  cursor: wait;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  background: rgba(255, 255, 255, 0.65);
}

.loader {
  width: 40px;
  aspect-ratio: 1;
  border-radius: 50%;
  border: 3px solid var(--bs-primary);
  animation:
    l20-1 0.8s infinite linear alternate,
    l20-2 1.6s infinite linear;
}

.loader.md {
  width: 50px;
  border-width: 4px;
}

.loader.lg {
  width: 60px;
  border-width: 5px;
}

@keyframes l20-1 {
  0% {
    clip-path: polygon(50% 50%, 0 0, 50% 0%, 50% 0%, 50% 0%, 50% 0%, 50% 0%);
  }
  12.5% {
    clip-path: polygon(50% 50%, 0 0, 50% 0%, 100% 0%, 100% 0%, 100% 0%, 100% 0%);
  }
  25% {
    clip-path: polygon(50% 50%, 0 0, 50% 0%, 100% 0%, 100% 100%, 100% 100%, 100% 100%);
  }
  50% {
    clip-path: polygon(50% 50%, 0 0, 50% 0%, 100% 0%, 100% 100%, 50% 100%, 0% 100%);
  }
  62.5% {
    clip-path: polygon(50% 50%, 100% 0, 100% 0%, 100% 0%, 100% 100%, 50% 100%, 0% 100%);
  }
  75% {
    clip-path: polygon(50% 50%, 100% 100%, 100% 100%, 100% 100%, 100% 100%, 50% 100%, 0% 100%);
  }
  100% {
    clip-path: polygon(50% 50%, 50% 100%, 50% 100%, 50% 100%, 50% 100%, 50% 100%, 0% 100%);
  }
}

@keyframes l20-2 {
  0% {
    transform: scaleY(1) rotate(0deg);
  }
  49.99% {
    transform: scaleY(1) rotate(135deg);
  }
  50% {
    transform: scaleY(-1) rotate(0deg);
  }
  100% {
    transform: scaleY(-1) rotate(-135deg);
  }
}
</style>
