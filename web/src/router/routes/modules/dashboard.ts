import type { AppRouteModule } from '/@/router/types';

import { LAYOUT } from '/@/router/constant';
import { t } from '/@/hooks/web/useI18n';

const dashboard: AppRouteModule = {
  path: '/dashboard',
  name: 'Dashboard',
  component: LAYOUT,
  redirect: '/dashboard/index',
  meta: {
    orderNo: 0,
    icon: 'ion:grid-outline',
    hideChildrenInMenu: true,
    title: t('routes.basic.dashboard'),
  },
  children: [
    {
      path: 'index',
      name: 'DashboardIndex',
      component: () => import('/@/views/dashboard/index.vue'),
      meta: {
        // affix: true,
        title: t('routes.basic.dashboard'),
      },
    },
  ],
};

export default dashboard;
