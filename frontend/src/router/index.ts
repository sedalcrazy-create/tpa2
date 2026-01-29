import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory('/tpa/'),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { guest: true }
    },
    {
      path: '/',
      component: () => import('@/layouts/MainLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'dashboard',
          component: () => import('@/views/DashboardView.vue')
        },
        {
          path: 'claims',
          name: 'claims',
          component: () => import('@/views/ClaimsView.vue')
        },
        {
          path: 'claims/:id',
          name: 'claim-detail',
          component: () => import('@/views/ClaimDetailView.vue')
        },
        {
          path: 'packages',
          name: 'packages',
          component: () => import('@/views/PackagesView.vue')
        },
        {
          path: 'centers',
          name: 'centers',
          component: () => import('@/views/CentersView.vue')
        },
        {
          path: 'settlements',
          name: 'settlements',
          component: () => import('@/views/SettlementsView.vue')
        },
        {
          path: 'members',
          name: 'members',
          component: () => import('@/views/MembersView.vue')
        },
        {
          path: 'reports',
          name: 'reports',
          component: () => import('@/views/ReportsView.vue')
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('@/views/UsersView.vue'),
          meta: { roles: ['system_admin', 'insurer_admin'] }
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('@/views/SettingsView.vue'),
          meta: { roles: ['system_admin', 'insurer_admin'] }
        },
        // Commission Module Routes
        {
          path: 'commission',
          name: 'commission-cases',
          component: () => import('@/views/CommissionCasesView.vue')
        },
        {
          path: 'commission/:id',
          name: 'commission-case-detail',
          component: () => import('@/views/CommissionCaseDetailView.vue')
        },
        {
          path: 'case-types',
          name: 'case-types',
          component: () => import('@/views/CaseTypesView.vue'),
          meta: { roles: ['system_admin'] }
        },
        {
          path: 'verdict-templates',
          name: 'verdict-templates',
          component: () => import('@/views/VerdictTemplatesView.vue'),
          meta: { roles: ['system_admin'] }
        },
        {
          path: 'insured-persons',
          name: 'insured-persons',
          component: () => import('@/views/InsuredPersonsView.vue')
        },
        // Social Work Module Routes
        {
          path: 'social-work',
          name: 'social-work-list',
          component: () => import('@/views/SocialWorkListView.vue')
        },
        {
          path: 'social-work/create',
          name: 'social-work-create',
          component: () => import('@/views/CreateSocialWorkCaseView.vue')
        },
        {
          path: 'social-work/:id',
          name: 'social-work-detail',
          component: () => import('@/views/SocialWorkDetailView.vue')
        },
        // New TPA Module Routes
        {
          path: 'price-conditions',
          name: 'price-conditions',
          component: () => import('@/views/ItemPriceConditionsView.vue'),
          meta: { roles: ['system_admin', 'insurer_admin'] }
        },
        {
          path: 'prescriptions',
          name: 'prescriptions',
          component: () => import('@/views/PrescriptionsView.vue')
        },
        {
          path: 'insurance-rules',
          name: 'insurance-rules',
          component: () => import('@/views/InsuranceRulesView.vue'),
          meta: { roles: ['system_admin', 'insurer_admin'] }
        },
        {
          path: 'contracts',
          name: 'contracts',
          component: () => import('@/views/ContractsView.vue'),
          meta: { roles: ['system_admin', 'insurer_admin'] }
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/NotFoundView.vue')
    }
  ]
})

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'login', query: { redirect: to.fullPath } })
  } else if (to.meta.guest && authStore.isAuthenticated) {
    next({ name: 'dashboard' })
  } else if (to.meta.roles && !authStore.hasRole(to.meta.roles as string[])) {
    next({ name: 'dashboard' })
  } else {
    next()
  }
})

export default router
