<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useApp } from '@/app/app.ts'
import { useSetupStore } from '@/stores/setup_store.ts'
import { useStatsStore } from '@/stores/stats_store.ts'
import { DateTime } from '@/app/core/date_time.ts'
import DashboardStatCard from '@/views/pages/Admin/Dashboard/DashboardStatCard.vue'
import Preloader from '@/views/shared/components/Preloader.vue'

defineOptions({
  name: 'AdminInventorySummaryPage',
})

const app = useApp()
const setupStore = useSetupStore()
const statsStore = useStatsStore()

const setup = computed(() => setupStore.get().value)
const stats = computed(() => statsStore.get().value)

const isLoading = ref<boolean>(false)
const isRefreshing = ref<boolean>(false)
const lastUpdatedAt = ref<string>('')

const emptyTrend = {
  change: 0,
  direction: 'stable',
} as const

const emptyMetric = {
  total: 0,
  week_trend: emptyTrend,
  month_trend: emptyTrend,
  annual_trend: emptyTrend,
}

const inventoryStats = computed(() => {
  const inventory = stats.value?.inventory
  return {
    products: inventory?.products ?? emptyMetric,
    productVariants: inventory?.productVariants ?? emptyMetric,
    stock: inventory?.stock ?? emptyMetric,
    bins: inventory?.bins ?? emptyMetric,
  }
})

const summaryMetrics = computed(() => {
  const productsTotal = inventoryStats.value.products.total
  const variantsTotal = inventoryStats.value.productVariants.total
  const stockTotal = inventoryStats.value.stock.total
  const binsTotal = inventoryStats.value.bins.total

  return {
    avgStockPerProduct: productsTotal > 0 ? (stockTotal / productsTotal).toFixed(2) : '0.00',
    avgStockPerVariant: variantsTotal > 0 ? (stockTotal / variantsTotal).toFixed(2) : '0.00',
    variantsPerProduct: productsTotal > 0 ? (variantsTotal / productsTotal).toFixed(2) : '0.00',
    stockPerBin: binsTotal > 0 ? (stockTotal / binsTotal).toFixed(2) : '0.00',
  }
})

const getStats = async (withLoader: boolean = false): Promise<void> => {
  if (isRefreshing.value) {
    return
  }
  if (!setup.value?.location?.id) {
    return
  }

  if (withLoader) {
    isLoading.value = true
  }
  isRefreshing.value = true

  try {
    const result = await app.services.stats.get(setup.value.location.id)
    if (result.success) {
      statsStore.set(result.stats)
      lastUpdatedAt.value = new DateTime().toDateTimeString()
    }
  } catch (ex) {
    app.services.logger.error(ex)
  } finally {
    isRefreshing.value = false
    if (withLoader) {
      isLoading.value = false
    }
  }
}

onMounted(async () => {
  const needsInitialLoader = !stats.value
  await getStats(needsInitialLoader)
})
</script>

<template>
  <div class="inventory-summary-page">
    <div class="card mb-6">
      <div class="card-body d-flex flex-column flex-md-row justify-content-between align-items-md-center">
        <div>
          <h5 class="card-title mb-2">Inventory Summary</h5>
          <p class="text-body-secondary mb-0">
            View aggregate stock levels, SKU spread, and bin utilization for
            <strong>{{ setup?.location?.name || 'current location' }}</strong>.
          </p>
        </div>
        <div class="d-flex align-items-center mt-4 mt-md-0">
          <small v-if="lastUpdatedAt" class="text-body-secondary me-3">Last updated {{ lastUpdatedAt }}</small>
          <button class="btn btn-outline-primary" :disabled="isRefreshing" @click="getStats(false)">
            <i class="bx bx-refresh me-2" />
            Refresh
          </button>
        </div>
      </div>
    </div>

    <div class="card mb-6" v-if="isLoading">
      <div class="card-body position-relative h-px-300">
        <preloader :overlay="true" />
      </div>
    </div>

    <template v-else>
      <div class="row">
        <div class="col-12 col-md-6 col-xl-3 mb-6">
          <dashboard-stat-card
            title="Products"
            :value="inventoryStats.products.total"
            icon-class="bx bx-package text-primary"
            icon-bg-class="bg-label-primary"
            :week-trend="inventoryStats.products.week_trend"
            :month-trend="inventoryStats.products.month_trend"
            :annual-trend="inventoryStats.products.annual_trend" />
        </div>
        <div class="col-12 col-md-6 col-xl-3 mb-6">
          <dashboard-stat-card
            title="Product Variants"
            :value="inventoryStats.productVariants.total"
            icon-class="bx bx-cube text-info"
            icon-bg-class="bg-label-info"
            :week-trend="inventoryStats.productVariants.week_trend"
            :month-trend="inventoryStats.productVariants.month_trend"
            :annual-trend="inventoryStats.productVariants.annual_trend" />
        </div>
        <div class="col-12 col-md-6 col-xl-3 mb-6">
          <dashboard-stat-card
            title="Stock Units"
            :value="inventoryStats.stock.total"
            icon-class="bx bx-box-alt text-success"
            icon-bg-class="bg-label-success"
            :week-trend="inventoryStats.stock.week_trend"
            :month-trend="inventoryStats.stock.month_trend"
            :annual-trend="inventoryStats.stock.annual_trend" />
        </div>
        <div class="col-12 col-md-6 col-xl-3 mb-6">
          <dashboard-stat-card
            title="Storage Bins"
            :value="inventoryStats.bins.total"
            icon-class="bx bx-box text-dark"
            icon-bg-class="bg-label-secondary"
            :week-trend="inventoryStats.bins.week_trend"
            :month-trend="inventoryStats.bins.month_trend"
            :annual-trend="inventoryStats.bins.annual_trend" />
        </div>
      </div>

      <div class="row">
        <div class="col-12 col-lg-6 mb-6">
          <div class="card h-100">
            <div class="card-header">
              <h5 class="mb-0">Summary Metrics</h5>
            </div>
            <div class="card-body">
              <div class="d-flex justify-content-between py-2 border-bottom">
                <span class="text-body-secondary">Average Stock per Product</span>
                <strong>{{ summaryMetrics.avgStockPerProduct }}</strong>
              </div>
              <div class="d-flex justify-content-between py-2 border-bottom">
                <span class="text-body-secondary">Average Stock per Variant</span>
                <strong>{{ summaryMetrics.avgStockPerVariant }}</strong>
              </div>
              <div class="d-flex justify-content-between py-2 border-bottom">
                <span class="text-body-secondary">Variants per Product</span>
                <strong>{{ summaryMetrics.variantsPerProduct }}</strong>
              </div>
              <div class="d-flex justify-content-between py-2">
                <span class="text-body-secondary">Stock per Bin</span>
                <strong>{{ summaryMetrics.stockPerBin }}</strong>
              </div>
            </div>
          </div>
        </div>

        <div class="col-12 col-lg-6 mb-6">
          <div class="card h-100">
            <div class="card-header">
              <h5 class="mb-0">Quick Actions</h5>
            </div>
            <div class="card-body">
              <p class="text-body-secondary mb-4">Use these shortcuts to act on inventory summary insights.</p>
              <div class="d-flex flex-wrap gap-2">
                <router-link :to="{ name: 'admin:inventory:products' }" class="btn btn-outline-primary">
                  <i class="bx bx-package me-2" />
                  View Products
                </router-link>
                <router-link :to="{ name: 'admin:inventory:products:add' }" class="btn btn-primary">
                  <i class="bx bx-plus me-2" />
                  Add Product
                </router-link>
                <router-link :to="{ name: 'admin:inventory:movements' }" class="btn btn-outline-secondary">
                  <i class="bx bx-transfer-alt me-2" />
                  View Movements
                </router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped></style>
