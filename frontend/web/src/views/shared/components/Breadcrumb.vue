<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { onMounted, ref } from 'vue'

interface IBreadcrumb {
  to: { name: string }
  text: string
  active: boolean
}

const route = useRoute()
const router = useRouter()

const breadcrumb = ref<IBreadcrumb[]>([])

onMounted(() => {
  if (route.meta.breadcrumb) {
    const routes = route.meta.breadcrumb as string[]
    routes.forEach((val: string, index: number) => {
      try {
        const targetRoute = router.resolve({ name: val })
        const targetRouteTitle = targetRoute.meta.title
        if (targetRouteTitle) {
          breadcrumb.value.push({
            to: { name: val },
            text: targetRouteTitle as string,
            active: index == routes.length - 1,
          })
        }
      } catch (ex) {}
    })
  }
})
</script>

<template>
  <nav class="mb-4">
    <ol class="breadcrumb">
      <li v-for="item in breadcrumb" class="breadcrumb-item" :class="{ active: item.active }">
        <template v-if="item.active">
          {{ item.text }}
        </template>
        <template v-else>
          <router-link :to="item.to">{{ item.text }}</router-link>
        </template>
      </li>
    </ol>
  </nav>
  <slot />
</template>

<style scoped></style>
