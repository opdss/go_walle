import type { AppRouteModule } from '/@/router/types';

import { LAYOUT } from '/@/router/constant';
import { t } from '/@/hooks/web/useI18n';

const project: AppRouteModule = {
  path: '/user',
  name: 'User',
  component: LAYOUT,
  redirect: '/user/index',
  meta: {
    hideChildrenInMenu: true,
    icon: 'ant-design:user-outlined',
    title: t('routes.basic.account'),
  },
  children: [
    {
      path: 'index',
      name: 'userIndex',
      component: () => import('/@/views/user/index.vue'),
      meta: {
        title: t('routes.basic.account'),
        icon: 'ant-design:user-outlined',
        hideMenu: true,
      },
    },
  ],
};

export default project;
