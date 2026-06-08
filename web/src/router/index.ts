import { createRouter, createWebHistory } from 'vue-router';
import ShellLayout from '../components/ShellLayout.vue';
import DashboardView from '../views/DashboardView.vue';
import AssetsView from '../views/AssetsView.vue';
import IdentitiesView from '../views/IdentitiesView.vue';
import VerificationsView from '../views/VerificationsView.vue';
import DataView from '../views/DataView.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: ShellLayout,
      redirect: '/dashboard',
      children: [
        { path: 'dashboard', name: 'dashboard', component: DashboardView },
        { path: 'assets', name: 'assets', component: AssetsView },
        { path: 'identities', name: 'identities', component: IdentitiesView },
        { path: 'verifications', name: 'verifications', component: VerificationsView },
        { path: 'data', name: 'data', component: DataView }
      ]
    }
  ]
});

export default router;
