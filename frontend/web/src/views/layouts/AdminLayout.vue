<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppLayout from '@/views/shared/components/Layouts/AppLayout.vue'
import { useUserStore } from '@/stores/user_store.ts'
import Breadcrumb from '@/views/shared/components/Breadcrumb.vue'

const route = useRoute()
const userStore = useUserStore()

const user = computed(() => userStore.get().value)
const routeName = computed(() => route.name as string)
</script>

<template>
  <app-layout>
    <template #menu.user>
      <div class="nav-link dropdown-toggle hide-arrow cursor-pointer p-0" data-bs-toggle="dropdown">
        <div class="avatar avatar-online">
          <img src="/assets/img/user/avatar-rounded.png" alt="" class="w-px-40 h-auto rounded-circle" />
        </div>
      </div>
      <ul class="dropdown-menu dropdown-menu-end">
        <li>
          <a class="dropdown-item" href="#">
            <div class="d-flex">
              <div class="flex-shrink-0 me-3">
                <div class="avatar avatar-online">
                  <img src="/assets/img/user/avatar-rounded.png" alt="" class="w-px-40 h-auto rounded-circle" />
                </div>
              </div>
              <div class="flex-grow-1">
                <h6 class="mb-0">{{ user?.fullName }}</h6>
                <small class="text-body-secondary">{{ user?.emailAddress }}</small>
              </div>
            </div>
          </a>
        </li>
        <li>
          <div class="dropdown-divider my-1"></div>
        </li>
        <li>
          <router-link :to="{ name: 'admin:account' }" class="dropdown-item">
            <i class="icon-base bx bx-user icon-md me-3"></i>
            <span>Account</span>
          </router-link>
        </li>
        <li>
          <router-link :to="{ name: 'admin:settings' }" class="dropdown-item">
            <i class="icon-base bx bx-slider-alt icon-md me-3"></i>
            <span>Preferences</span>
          </router-link>
        </li>
        <li>
          <div class="dropdown-divider my-1"></div>
        </li>
        <li>
          <router-link :to="{ name: 'auth:logout' }" class="dropdown-item">
            <i class="icon-base bx bx-power icon-md me-3"></i>
            <span>Log Out</span>
          </router-link>
        </li>
      </ul>
    </template>
    <template #menu.sidebar>
      <ul class="menu-inner py-1">
        <li class="menu-item" :class="{ active: routeName.startsWith('admin:home') }">
          <router-link :to="{ name: 'admin:home' }" class="menu-link">
            <i class="menu-icon tf-icons bx bx-bar-chart"></i>
            <div class="text-truncate">Dashboard</div>
          </router-link>
        </li>
        <li
          class="menu-item"
          :class="{
            'active open':
              routeName.startsWith('admin:staff') ||
              routeName.startsWith('admin:roles') ||
              routeName.startsWith('admin:tasks'),
          }">
          <div class="menu-link menu-toggle cursor-pointer">
            <i class="menu-icon tf-icons bx bx-user"></i>
            <div class="text-truncate">Admin</div>
          </div>
          <ul class="menu-sub">
            <li class="menu-item" :class="{ active: routeName.startsWith('admin:staff') }">
              <router-link :to="{ name: 'admin:staff' }" class="menu-link">
                <div class="text-truncate">Staff</div>
              </router-link>
            </li>
            <li class="menu-item" :class="{ active: routeName.startsWith('admin:roles') }">
              <router-link :to="{ name: 'admin:roles' }" class="menu-link">
                <div class="text-truncate">Roles</div>
              </router-link>
            </li>
            <li class="menu-item" :class="{ active: routeName.startsWith('admin:tasks') }">
              <router-link :to="{ name: 'admin:tasks' }" class="menu-link">
                <div class="text-truncate">Tasks</div>
              </router-link>
            </li>
          </ul>
        </li>
        <li class="menu-item" :class="{ 'active open': routeName.startsWith('admin:inventory') }">
          <div class="menu-link menu-toggle cursor-pointer">
            <i class="menu-icon tf-icons bx bx-briefcase"></i>
            <div class="text-truncate">Inventory</div>
          </div>
          <ul class="menu-sub">
            <li class="menu-item" :class="{ active: routeName == 'admin:inventory' }">
              <router-link :to="{ name: 'admin:inventory' }" class="menu-link">
                <div class="text-truncate">Summery</div>
              </router-link>
            </li>
            <li class="menu-item" :class="{ active: routeName.startsWith('admin:inventory:movements') }">
              <router-link :to="{ name: 'admin:inventory:movements' }" class="menu-link">
                <div class="text-truncate">Movements</div>
              </router-link>
            </li>
            <li class="menu-item" :class="{ active: routeName.startsWith('admin:inventory:products') }">
              <router-link :to="{ name: 'admin:inventory:products' }" class="menu-link">
                <div class="text-truncate">Products</div>
              </router-link>
            </li>
          </ul>
        </li>

        <li
          class="menu-item"
          :class="{
            'active open': routeName.startsWith('admin:operations') || routeName.startsWith('admin:suppliers'),
          }">
          <div class="menu-link menu-toggle cursor-pointer">
            <i class="menu-icon tf-icons bx bx-cart"></i>
            <div class="text-truncate">Operations</div>
          </div>
          <ul class="menu-sub">
            <li class="menu-item" :class="{ active: routeName.startsWith('admin:operations:orders') }">
              <router-link :to="{ name: 'admin:operations:orders' }" class="menu-link">
                <div class="text-truncate">Orders</div>
              </router-link>
            </li>
            <li class="menu-item" :class="{ active: routeName.startsWith('admin:operations:shipments') }">
              <router-link :to="{ name: 'admin:operations:shipments' }" class="menu-link">
                <div class="text-truncate">Shipments</div>
              </router-link>
            </li>
            <li class="menu-item" :class="{ active: routeName.startsWith('admin:suppliers:customers') }">
              <router-link :to="{ name: 'admin:suppliers:customers' }" class="menu-link">
                <div class="text-truncate">Customers</div>
              </router-link>
            </li>
            <li class="menu-item" :class="{ active: routeName.startsWith('admin:suppliers:list') }">
              <router-link :to="{ name: 'admin:suppliers:list' }" class="menu-link">
                <div class="text-truncate">Suppliers</div>
              </router-link>
            </li>
          </ul>
        </li>
        <li class="menu-item" :class="{ active: routeName.startsWith('admin:warehouse') }">
          <router-link :to="{ name: 'admin:warehouse' }" class="menu-link">
            <i class="menu-icon tf-icons bx bx-building-house"></i>
            <div class="text-truncate">Locations</div>
          </router-link>
        </li>
        <li
          class="menu-item position-absolute bottom-0 mb-0 py-2 bg-white"
          :class="{ active: route.name === 'admin:help' }">
          <router-link :to="{ name: 'admin:help' }" class="menu-link">
            <i class="menu-icon tf-icons bx bx-help-circle"></i>
            <div class="text-truncate">Help</div>
          </router-link>
        </li>
        <li class="menu-item" :class="{ active: routeName.startsWith('admin:settings') }">
          <router-link :to="{ name: 'admin:settings' }" class="menu-link">
            <i class="menu-icon tf-icons bx bx-cog"></i>
            <div class="text-truncate">Settings</div>
          </router-link>
        </li>
      </ul>
    </template>
    <template #page.content>
      <slot />
    </template>
    <template #page.footer>
      <div class="mb-2 mb-md-0">
        © {{ new Date().getFullYear() }}
        <a href="https://fiverr.com/inbox/ndoubismarck" target="_blank" class="footer-link">
          Warehouse Management System
        </a>
      </div>
      <div class="d-none d-lg-inline-block">
        <a href="https://fiverr.com/inbox/ndoubismarck" target="_blank" class="footer-link me-4">Documentation</a>
        <a href="https://fiverr.com/inbox/ndoubismarck" target="_blank" class="footer-link me-4">Support</a>
      </div>
    </template>
    <template #page.breadcrumb>
      <breadcrumb :home-item="{ text: 'Admin', route: { name: 'admin:home' } }" />
    </template>
  </app-layout>
</template>
