<script setup lang="ts">
import { onMounted, shallowRef } from 'vue'

const menu = shallowRef<AppMenu>()

class AppMenu {
  private readonly el: HTMLElement

  constructor() {
    this.el = document.querySelector('.menu-inner') as HTMLElement
    if (this.el) {
      this.el.addEventListener('click', this.handleClick.bind(this))
    }
    document.querySelectorAll('aside a').forEach((item) => {
      item.addEventListener('click', () => {
        document.body.classList.remove('layout-menu-expanded')
      })
    })
  }

  handleClick(e: Event) {
    const target = e.target as HTMLElement
    const toggle = target.closest('.menu-toggle')
    if (!toggle) return

    const menuItem = toggle.parentElement
    if (!menuItem) return

    const submenu = menuItem.querySelector('.menu-sub') as HTMLElement
    if (!submenu) return

    const siblings = menuItem.parentElement?.querySelectorAll(':scope > .menu-item.open') || []
    siblings.forEach((sib) => {
      if (sib !== menuItem) {
        this.closeMenu(sib as HTMLElement)
      }
    })
    menuItem.classList.contains('open') ? this.closeMenu(menuItem) : this.openMenu(menuItem)
  }

  openMenu(menuItem: HTMLElement) {
    const submenu = menuItem.querySelector('.menu-sub') as HTMLElement
    menuItem.classList.add('open')
    submenu.style.height = '0px'
    const fullHeight = submenu.scrollHeight
    submenu.style.transition = 'height 0.3s ease'
    requestAnimationFrame(() => {
      submenu.style.height = fullHeight + 'px'
    })
    submenu.addEventListener('transitionend', () => (submenu.style.height = 'auto'), { once: true })
  }

  closeMenu(menuItem: HTMLElement) {
    const submenu = menuItem.querySelector('.menu-sub') as HTMLElement
    const height = submenu.scrollHeight
    submenu.style.height = height + 'px'
    requestAnimationFrame(() => {
      submenu.style.height = '0px'
    })
    submenu.addEventListener(
      'transitionend',
      () => {
        menuItem.classList.remove('open')
        submenu.style.height = ''
      },
      { once: true },
    )
  }

  handleToggleMobileMenu(evt: Event) {
    evt.preventDefault()
    document.body.classList.toggle('layout-menu-expanded')
  }
}

onMounted(() => {
  menu.value = new AppMenu()
})
</script>

<template>
  <div class="layout-wrapper layout-content-navbar layout-menu-fixed">
    <div class="layout-container">
      <aside class="layout-menu menu-vertical menu bg-menu-theme">
        <div class="app-brand">
          <router-link :to="{ name: 'admin:home' }" class="app-brand-link">
            <span class="app-brand-logo">
              <img src="/assets/img/app/app-icon-dark.svg" alt="" />
            </span>
            <span class="app-brand-text menu-text fw-bold">WMS</span>
          </router-link>
          <!--
          <a @click="menu?.handleToggleMobileMenu"
             class="layout-menu-toggle menu-link text-center d-flex align-items-center p-1">
            <i class="bx bx-x d-block d-xl-none icon-lg"></i>
          </a>
          -->
        </div>
        <div class="menu-divider mt-0"></div>
        <div class="menu-inner-shadow"></div>
        <slot name="menu.sidebar" />
      </aside>
      <div class="layout-page">
        <nav
          class="layout-navbar container-xxl navbar-detached navbar navbar-expand-xl align-items-center bg-navbar-theme">
          <div class="layout-menu-toggle navbar-nav align-items-xl-center me-4 me-xl-0 d-xl-none">
            <a @click="menu?.handleToggleMobileMenu" class="nav-item nav-link px-0 me-xl-6">
              <i class="icon-base bx bx-menu icon-xl"></i>
            </a>
          </div>

          <div class="navbar-nav-right d-flex align-items-center justify-content-end" id="navbar-collapse">
            <div class="navbar-nav align-items-center me-auto">
              <div class="nav-item d-flex align-items-center">
                <span class="w-px-22 h-px-22 cursor-pointer">
                  <i class="icon-base bx bx-search icon-md"></i>
                </span>
                <button
                  type="button"
                  class="border-0 bg-transparent shadow-none ps-1 ps-sm-2 d-md-block d-none text-light">
                  Search [CTRL + K]
                </button>
              </div>
            </div>
            <ul class="navbar-nav flex-row align-items-center ms-md-auto">
              <li class="nav-item lh-1 me-4"></li>
              <li class="nav-item navbar-dropdown dropdown-user dropdown">
                <slot name="menu.user" />
              </li>
            </ul>
          </div>
        </nav>
        <div class="content-wrapper">
          <div class="container-xxl flex-grow-1 container-p-y">
            <slot name="page.breadcrumb" />
            <slot name="page.content" />
          </div>
          <footer class="content-footer footer bg-footer-theme">
            <div class="container-xxl">
              <div
                class="footer-container d-flex align-items-center justify-content-between py-6 flex-md-row flex-column">
                <slot name="page.footer" />
              </div>
            </div>
          </footer>
          <div class="content-backdrop fade"></div>
        </div>
      </div>
    </div>
    <div @click="menu?.handleToggleMobileMenu" class="layout-overlay layout-menu-toggle"></div>
  </div>
</template>

<style>
.menu-sub {
  overflow: hidden;
  height: auto;
  transition: height 0.3s;
}
</style>
