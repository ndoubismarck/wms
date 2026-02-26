<script setup lang="ts">
interface IHelpCard {
  title: string
  description: string
  action: string
  icon: string
  iconBgClass: string
  to?: { name: string }
  href?: string
  isExternal?: boolean
}

const cards: IHelpCard[] = [
  {
    title: 'Chat',
    description: 'Connect with support instantly for real-time help and guided troubleshooting assistance.',
    action: 'Start Chat',
    icon: 'bx bx-message-dots',
    iconBgClass: 'bg-label-primary',
    href: 'https://fiverr.com/ndoubismarck',
    isExternal: true,
  },
  {
    title: 'Raise Ticket',
    description: 'Submit technical or operational issues and monitor progress until resolution.',
    action: 'Open Ticket',
    icon: 'bx bx-help-circle',
    iconBgClass: 'bg-label-primary',
    href: 'https://fiverr.com/ndoubismarck',
    isExternal: true,
  },
  {
    title: 'Documentation',
    description: 'Access detailed guides for inventory management, workflows, and system configuration.',
    action: 'View Documentation',
    icon: 'bx bx-book-open',
    iconBgClass: 'bg-label-primary',
    href: 'https://fiverr.com/ndoubismarck',
    isExternal: true,
  },
  {
    title: 'System Information',
    description: 'Review platform version, build details, uptime status, and recent release updates.',
    action: 'View System Info',
    icon: 'bx bx-info-circle',
    iconBgClass: 'bg-label-primary',
    href: 'https://fiverr.com/ndoubismarck',
    isExternal: true,
  },
  {
    title: 'Feature Requests',
    description: 'Propose enhancements and explore planned features and upcoming improvements.',
    action: 'View Roadmap',
    icon: 'bx bx-map',
    iconBgClass: 'bg-label-primary',
    href: 'https://fiverr.com/ndoubismarck',
    isExternal: true,
  },
  {
    title: 'FAQs',
    description: 'Find answers to common issues including stock mismatches and workflow conflicts.',
    action: 'Browse FAQs',
    icon: 'bx bx-help-circle',
    iconBgClass: 'bg-label-primary',
    href: 'https://fiverr.com/ndoubismarck',
    isExternal: true,
  },
  {
    title: 'Resources',
    description: 'Download templates, checklists, and reference materials for warehouse operations.',
    action: 'View Resources',
    icon: 'bx bx-folder-open',
    iconBgClass: 'bg-label-primary',
    href: 'https://fiverr.com/ndoubismarck',
    isExternal: true,
  },
  {
    title: 'System Status',
    description: 'Check current system uptime, service health, and scheduled maintenance notices.',
    action: 'View Status',
    icon: 'bx bx-pulse',
    iconBgClass: 'bg-label-primary',
    href: 'https://fiverr.com/ndoubismarck',
    isExternal: true,
  },
]
</script>

<template>
  <div class="card mb-4">
    <div class="card-body">
      <h5 class="card-title mb-2">Help, Support and Activation</h5>
      <p class="text-body-secondary mb-0">
        Manage product activation, raise tickets, chat with support, and access warehouse documentation.
      </p>
    </div>
  </div>

  <div class="row g-4">
    <div v-for="card in cards" :key="card.title" class="col-12 col-md-6 col-lg-3">
      <div class="card h-100 help-card">
        <div class="card-body d-flex flex-column">
          <div class="d-flex align-items-center mb-1">
            <div class="avatar">
              <div class="avatar-initial rounded-2" :class="card.iconBgClass">
                <i class="icon-base icon-lg" :class="card.icon"></i>
              </div>
            </div>
            <h6 class="mb-0 ms-3">{{ card.title }}</h6>
          </div>
          <hr />
          <p class="text-body-secondary flex-grow-1 mb-4">{{ card.description }}</p>

          <router-link v-if="card.to" :to="card.to" class="btn btn-sm btn-outline-primary align-self-start">
            {{ card.action }}
          </router-link>

          <a
            v-else-if="card.href"
            :href="card.href"
            :target="card.isExternal ? '_blank' : '_self'"
            :rel="card.isExternal ? 'noopener noreferrer' : undefined"
            class="btn btn-sm btn-outline-primary align-self-start">
            {{ card.action }}
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.help-card .avatar {
  width: 2.25rem;
  height: 2.25rem;
}

.help-card .avatar i {
  font-size: 1.1rem;
}
</style>
