<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { onMounted, ref, toRef, watch } from 'vue'

interface IBreadcrumb {
  text: string
  route: { name: string }
  active: boolean
}

interface IBreadcrumbHomeItem {
  text: string
  route: { name: string }
}

const route = useRoute()
const router = useRouter()

const props = defineProps<{
  homeItem: IBreadcrumbHomeItem
}>()

const breadcrumb = ref<IBreadcrumb[]>([])
const breadcrumbHomeItem = toRef(props, 'homeItem')

const loadBreadcrumb = () => {
  breadcrumb.value = []
  if (route.meta.breadcrumb) {
    const routes = route.meta.breadcrumb as string[]
    routes.forEach((val: string, index: number) => {
      try {
        const targetRoute = router.resolve({ name: val })
        const targetRouteTitle = targetRoute.meta.title
        if (targetRouteTitle) {
          breadcrumb.value.push({
            text: targetRouteTitle as string,
            route: { name: val },
            active: index == routes.length - 1,
          })
        }
      } catch (ex) {}
    })
  }
}

watch(route, () => {
  loadBreadcrumb()
})

onMounted(() => {
  loadBreadcrumb()
})
</script>

<template>
  <nav v-if="breadcrumb.length > 0" class="mb-4">
    <ol class="breadcrumb">
      <li class="breadcrumb-item">
        <router-link :to="breadcrumbHomeItem.route">{{ breadcrumbHomeItem.text }}</router-link>
      </li>
      <li v-for="item in breadcrumb" class="breadcrumb-item" :class="{ active: item.active }">
        <template v-if="item.active">
          {{ item.text }}
        </template>
        <template v-else>
          <router-link :to="item.route">{{ item.text }}</router-link>
        </template>
      </li>
    </ol>
  </nav>
</template>

<style scoped></style>
