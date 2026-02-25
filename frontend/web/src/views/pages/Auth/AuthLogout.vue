<script setup lang="ts">
import Preloader from '@/views/shared/components/Preloader.vue'
import { useApp } from '@/app/app.ts'
import { useRouter } from 'vue-router'
import { onMounted } from 'vue'

const app = useApp()
const router = useRouter()

onMounted(async () => {
  await app.services.auth.logout()
  await app.helpers.async.sleep(3000)
  await router.replace({ name: 'auth:login' })
})
</script>

<template>
  <div class="card px-sm-6 px-0">
    <div class="card-body position-relative">
      <preloader :overlay="true" />
    </div>
  </div>
</template>
<style scoped>
.card .card-body {
  min-height: 450px;
}
</style>
