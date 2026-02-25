<script setup lang="ts">
import { toRef } from 'vue'
import { EAlertMessageType, type IAlertMessage } from '@/app/types'

const props = defineProps<{
  params: IAlertMessage
}>()

const params = toRef(props, 'params')
</script>

<template>
  <div
    class="alert"
    :class="{
      'alert-danger': params.type == EAlertMessageType.Error,
      'alert-success': params.type == EAlertMessageType.Success,
      'alert-warning': params.type == EAlertMessageType.Warning,
      'alert-info': params.type == EAlertMessageType.Info,
    }">
    <h4 v-if="params.title">
      {{ params.title }}
    </h4>
    <template v-if="params.body.text">
      <div>
        {{ params.body.text }}
      </div>
    </template>
    <template v-if="params.body.html">
      <div v-html="params.body.html"></div>
    </template>
  </div>
</template>
