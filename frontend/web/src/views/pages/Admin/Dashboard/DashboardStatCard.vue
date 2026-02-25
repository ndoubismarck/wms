<script setup lang="ts">
import { toRef } from 'vue'

interface Trend {
  change: number
  direction: string | 'up' | 'down' | 'stable'
}

interface IDashboardStatsCardParams {
  title: string
  value: number
  iconClass: string
  iconBgClass: string
  weekTrend: Trend
  monthTrend: Trend
  annualTrend: Trend
}

const props = defineProps<IDashboardStatsCardParams>()

const title = toRef(props, 'title')
const value = toRef(props, 'value')
const iconClass = toRef(props, 'iconClass')
const iconBgClass = toRef(props, 'iconBgClass')
const weekTrend = toRef(props, 'weekTrend')
const monthTrend = toRef(props, 'monthTrend')
const annualTrend = toRef(props, 'annualTrend')
</script>

<template>
  <div class="card h-100">
    <div class="card-body">
      <div class="card-title d-flex align-items-start justify-content-between mb-4">
        <div class="avatar">
          <div class="avatar-initial rounded-2" :class="iconBgClass">
            <i class="icon-base icon-lg" :class="iconClass"></i>
          </div>
        </div>
        <div class="dropdown">
          <button class="btn p-0" type="button" data-bs-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
            <i class="icon-base bx bx-dots-vertical-rounded text-body-secondary"></i>
          </button>
          <div class="dropdown-menu dropdown-menu-end">
            <a class="dropdown-item" href="javascript:void(0);">View More</a>
            <a class="dropdown-item" href="javascript:void(0);">Delete</a>
          </div>
        </div>
      </div>
      <p class="mb-1">{{ title }}</p>
      <h4 class="card-title mb-3">{{ value }}</h4>
      <small
        v-if="weekTrend.change != 0"
        class="fw-medium"
        :class="{
          'text-success': weekTrend.direction === 'up',
          'text-danger': weekTrend.direction === 'down',
          'text-muted': weekTrend.direction === 'stable',
        }">
        <i
          class="icon-base bx"
          :class="{
            'bx-up-arrow-alt': weekTrend.direction === 'up',
            'bx-down-arrow-alt': weekTrend.direction === 'down',
            'bx-time': weekTrend.direction === 'stable',
          }"></i>
        {{ weekTrend.change.toFixed(2) }}%
      </small>
    </div>
  </div>
</template>
