import { createRouter, createWebHistory } from 'vue-router';
import ShellLayout from '../components/ShellLayout.vue';
import DashboardView from '../views/DashboardView.vue';
import AssetsView from '../views/AssetsView.vue';
import IdentitiesView from '../views/IdentitiesView.vue';
import VerificationsView from '../views/VerificationsView.vue';
import DataView from '../views/DataView.vue';
import LoginView from '../views/LoginView.vue';
import RegisterView from '../views/RegisterView.vue';
import { hasPermission, isLoggedIn, loadProfile } from '../services/auth';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', name: 'login', component: LoginView },
    {
      path: '/',
      component: ShellLayout,
      redirect: '/dashboard',
      children: [
        { path: 'dashboard', name: 'dashboard', component: DashboardView },
        { path: 'assets', name: 'assets', component: AssetsView },
        { path: 'identities', name: 'identities', component: IdentitiesView },
        { path: 'verifications', name: 'verifications', component: VerificationsView },
        { path: 'data', name: 'data', component: DataView },
        { path: 'register', name: 'register', component: RegisterView, meta: { permission: 'user:create' } }
      ]
    }
  ]
});

router.beforeEach(async (to) => {
  if (to.name === 'login') {
    return isLoggedIn() ? '/dashboard' : true;
  }
  if (!isLoggedIn()) {
    return '/login';
  }
  try {
    await loadProfile();
  } catch {
    return '/login';
  }
  const permission = to.meta.permission as string | undefined;
  if (permission && !hasPermission(permission)) {
    return '/dashboard';
  }
  return true;
});

export default router;
