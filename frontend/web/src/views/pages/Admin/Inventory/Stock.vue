<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useProductsStore } from '@/stores/products_store.ts'
import Breadcrumb from '@/views/shared/components/Breadcrumb.vue'
import VueTable, { type ITableColumn } from '@/views/shared/components/VueTable.vue'
import { useApp } from '@/app/app.ts'
import Preloader from '@/views/shared/components/Preloader.vue'
import type { ProductModel } from '@/app/models/product_model.ts'
import { useRouter } from 'vue-router'

const app = useApp()
const router = useRouter()
const productsStore = useProductsStore()

const products = computed(() => productsStore.get().value)

const page = ref<number>(1)
const isLoading = ref<boolean>(false)
const isFetching = ref<boolean>(false)
const checkedProductIds = ref<string[]>([])

const getProducts = async () => {
  try {
    const result = await app.services.products.getMany({
      page: page.value,
      limit: 25,
    })
    if (result.success) {
      if (result.pagination.hasNext) {
        page.value = result.pagination.next
      }
      productsStore.set(result.products)
    }
  } catch (ex) {
    app.services.logger.error(ex)
  }
}

const renderStatus = (value?: string): string => {
  if (value == 'available') {
    return `<span class="badge text-success bg-success-subtle text-capitalize">${value}</span>`
  }
  return `<span class="badge text-danger bg-danger-subtle text-capitalize">${value}</span>`
}

const getTableColumns = (row: ProductModel): ITableColumn[] => {
  return [
    {
      name: 'Name',
      value: {
        link: {
          text: row.name,
          href: {
            name: 'admin:inventory:products:view',
            params: {
              id: row.id,
            },
          },
          class: 'text-body',
        },
      },
    },
    {
      name: 'SKU',
      value: {
        link: {
          text: row.sku,
          href: {
            name: 'admin:inventory:products:view',
            params: { id: row.id },
          },
          class: 'text-body',
        },
      },
    },
    {
      name: 'Location',
      value: {
        dropdown: {
          button: {
            text: row.bin.name,
            class: 'fs-inherit text-body border-0',
          },
          items: [
            {
              link: {
                text: row.bin.name,
                href: {
                  name: 'admin:warehouse:bins:view',
                  params: { id: row.bin.id },
                },
              },
            },
            {
              link: {
                text: row.bin.shelfLevel.name,
                href: {
                  name: 'admin:warehouse:shelves:levels:view',
                  params: { id: row.bin.shelfLevel.id },
                },
              },
            },
            {
              link: {
                text: row.bin.shelfLevel.shelf.name,
                href: {
                  name: 'admin:warehouse:shelves:view',
                  params: { id: row.bin.shelfLevel.shelf.id },
                },
              },
            },
            {
              link: {
                text: row.bin.shelfLevel.shelf.bay.name,
                href: {
                  name: 'admin:warehouse:bays:view',
                  params: { id: row.bin.shelfLevel.shelf.bay.id },
                },
              },
            },
            {
              link: {
                text: row.bin.shelfLevel.shelf.bay.aisle.name,
                href: {
                  name: 'admin:warehouse:bins:view',
                  params: { id: row.bin.shelfLevel.shelf.bay.aisle.id },
                },
              },
            },
            {
              image: {
                src: app.api.private.url(`/products/${row.id}/barcode`),
                width: 250,
              },
              class: 'border-top pt-4',
            },
          ],
        },
      },
    },
    {
      name: 'Status',
      value: {
        html: renderStatus(row.status),
      },
    },
    {
      name: 'Date Updated',
      value: {
        text: row.updatedAt.toDateTimeString(),
      },
    },
    {
      name: 'Date Created',
      value: {
        text: row.updatedAt.toDateTimeString(),
      },
    },
  ]
}

const handleTabledChecked = (ids: string[]) => {
  checkedProductIds.value = ids
}

const handleTableScrollEnd = async () => {
  if (!isFetching.value) {
    isFetching.value = true
    await getProducts()
    await app.helpers.async.sleep(3000)
    isFetching.value = false
  }
}

onMounted(async () => {
  isLoading.value = products.value.length == 0
  await getProducts()
  await app.helpers.async.sleep(3000)
  isLoading.value = false
})
</script>

<template>
  <breadcrumb>
    <div class="card">
      <template v-if="isLoading">
        <div class="card-body position-relative h-px-300">
          <preloader :overlay="true" />
        </div>
      </template>
      <template v-else>
        <div class="card-header d-flex justify-content-start">
          <button
            style="height: 40px !important"
            @click="router.push({ name: 'admin:inventory:products:add' })"
            class="btn btn-outline-primary d-flex justify-content-center">
            <i class="bx bx-plus" />
            <span class="d-none d-sm-inline ms-2">Add Product</span>
          </button>
          <div class="dropdown ms-4">
            <button
              style="height: 40px !important"
              class="btn btn-outline-primary d-flex justify-content-center"
              data-bs-toggle="dropdown"
              data-bs-boundary="viewport"
              aria-expanded="false"
              :disabled="checkedProductIds.length == 0">
              <span class="d-none d-sm-inline me-2">With Selected</span>
              <i class="bx bx-chevron-down" />
            </button>
            <ul class="dropdown-menu">
              <li>
                <a href="#" class="dropdown-item text-primary">
                  <div class="d-flex justify-content-start"><i class="bx bx-download me-2"></i><span>Export</span></div>
                </a>
              </li>
              <li>
                <a href="#" class="dropdown-item text-danger">
                  <div class="d-flex justify-content-start"><i class="bx bx-trash me-2"></i><span>Delete</span></div>
                </a>
              </li>
            </ul>
          </div>
        </div>

        <div class="card-body p-0" style="min-height: 300px">
          <vue-table
            :rows="products"
            :columns="getTableColumns"
            @checked="handleTabledChecked"
            @scrollend="handleTableScrollEnd"
            :checkboxes="true">
            <template #actions="{ row }">
              <router-link
                :to="{ name: 'admin:inventory:products:view', params: { id: row.id } }"
                class="btn btn-lg text-primary border-0">
                <i class="bx bx-dots-vertical"></i>
              </router-link>
            </template>
            <template #preloader>
              <div v-if="isFetching" class="position-relative h-auto py-6 my-6">
                <preloader :overlay="true" />
              </div>
            </template>
          </vue-table>
        </div>
      </template>
    </div>
  </breadcrumb>
  <router-view />
</template>
