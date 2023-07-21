import type { AppRouteModule } from '/@/router/types';

import { LAYOUT } from '/@/router/constant';
import { t } from '/@/hooks/web/useI18n';

const environment: AppRouteModule = {
  path: '/environment',
  name: 'Environment',
  component: LAYOUT,
  redirect: '/environment/index',
  meta: {
    hideChildrenInMenu: true,
    icon: 'ant-design:global-outlined',
    title: t('routes.basic.environment'),
    orderNo: 100000,
  },
  children: [
    {
      path: 'index',
      name: 'EnvironmentIndex',
      component: () => import('/@/views/environment/index.vue'),
      meta: {
        title: t('routes.basic.environment'),
        icon: 'ant-design:global-outlined',
        hideMenu: true,
      },
    },
  ],
};

export default environment;
