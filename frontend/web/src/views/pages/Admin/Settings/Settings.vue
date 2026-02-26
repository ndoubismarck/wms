<script setup lang="ts">
interface ISettingsSection {
  title: string
  description: string
  icon: string
  iconBgClass: string
  buttonLabel: string
}

const sections: ISettingsSection[] = [
  {
    title: 'Organization',
    description: 'Manage company profile, default timezone, regional settings, and operating calendar.',
    icon: ' bx bx-home-alt-2',
    iconBgClass: 'bg-label-primary',
    buttonLabel: 'Manage Organization',
  },
  {
    title: 'Warehouse Setup',
    description: 'Configure warehouse structure, location mapping, handling zones, and flow defaults.',
    icon: 'bx bx-buildings',
    iconBgClass: 'bg-label-primary',
    buttonLabel: 'Configure Warehouse',
  },
  {
    title: 'Inventory Rules',
    description: 'Set stock thresholds, replenishment behavior, reservation rules, and movement policies.',
    icon: 'bx bx-package',
    iconBgClass: 'bg-label-primary',
    buttonLabel: 'Configure Rules',
  },
  {
    title: 'Integrations',
    description: 'Connect external systems such as ERP, shipping partners, and notification providers.',
    icon: 'bx bx-link-alt',
    iconBgClass: 'bg-label-primary',
    buttonLabel: 'Manage Integrations',
  },
  {
    title: 'Licensing',
    description: 'Review plan details, activation status, feature entitlements, and license allocations.',
    icon: 'bx bx-lock-open',
    iconBgClass: 'bg-label-primary',
    buttonLabel: 'View License',
  },
  {
    title: 'Backup & Data',
    description: 'Control backup schedules, retention rules, export routines, and restore readiness.',
    icon: 'bx bx-cloud',
    iconBgClass: 'bg-label-primary',
    buttonLabel: 'Manage Backups',
  },
  {
    title: 'Advanced',
    description: 'Adjust technical defaults, audit behavior, automation flags, and system-level controls.',
    icon: 'bx bx-slider-alt',
    iconBgClass: 'bg-label-primary',
    buttonLabel: 'Advanced Settings',
  },
  {
    title: 'Audit & Logs',
    description: 'Review system activity logs, stock changes, and operational history.',
    icon: 'bx bx-clipboard',
    iconBgClass: 'bg-label-primary',
    buttonLabel: 'View Logs',
  },
  {
    title: 'Data Management',
    description: 'Manage data imports, exports, cleanup routines, and archival policies.',
    icon: 'bx bx-folder',
    iconBgClass: 'bg-label-primary',
    buttonLabel: 'Manage Data',
  },
]
</script>

<template>
  <div class="row g-4">
    <div v-for="section in sections" :key="section.title" class="col-12 col-md-6 col-xl-4">
      <div class="card h-100 settings-section-card">
        <div class="card-body d-flex flex-column">
          <div class="d-flex align-items-center mb-1">
            <div class="avatar">
              <div class="avatar-initial rounded-2" :class="section.iconBgClass">
                <i class="icon-base icon-lg" :class="section.icon"></i>
              </div>
            </div>
            <h6 class="mb-0 ms-3">{{ section.title }}</h6>
          </div>
          <hr />
          <p class="text-body-secondary flex-grow-1 mb-4">{{ section.description }}</p>
          <button type="button" class="btn btn-sm btn-outline-primary align-self-start">
            {{ section.buttonLabel }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-section-card .avatar {
  width: 2.25rem;
  height: 2.25rem;
}

.settings-section-card .avatar i {
  font-size: 1.1rem;
}
</style>
