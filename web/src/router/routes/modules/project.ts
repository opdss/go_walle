import type { AppRouteModule } from '/@/router/types';

import { LAYOUT } from '/@/router/constant';
import { t } from '/@/hooks/web/useI18n';

const environment: AppRouteModule = {
  path: '/project',
  name: 'Project',
  component: LAYOUT,
  redirect: '/project/index',
  meta: {
    hideChildrenInMenu: true,
    icon: 'ant-design:database-outlined',
    title: t('routes.basic.project'),
    orderNo: 100000,
  },
  children: [
    {
      path: 'index',
      name: 'ProjectIndex',
      component: () => import('/@/views/project/index.vue'),
      meta: {
        title: t('routes.basic.project'),
        icon: 'ant-design:database-outlined',
        hideMenu: true,
      },
    },
    {
      path: 'create',
      name: 'ProjectCreate',
      component: () => import('/@/views/project/form.vue'),
      meta: {
        title: t('routes.basic.project'),
        icon: 'ant-design:database-outlined',
        hideMenu: true,
      },
    },
    {
      path: 'update/:id',
      name: 'ProjectUpdate',
      component: () => import('/@/views/project/form.vue'),
      meta: {
        title: t('routes.basic.project'),
        icon: 'ant-design:database-outlined',
        hideMenu: true,
      },
    },
    {
      path: 'detail/:id',
      name: 'ProjectDetail',
      component: () => import('/@/views/project/form.vue'),
      meta: {
        title: t('routes.basic.project'),
        icon: 'ant-design:database-outlined',
        hideMenu: true,
      },
    },
  ],
};

export default environment;
