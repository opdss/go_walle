import type { AppRouteModule } from '/@/router/types';

import { LAYOUT } from '/@/router/constant';
import { t } from '/@/hooks/web/useI18n';

const environment: AppRouteModule = {
  path: '/space',
  name: 'Space',
  component: LAYOUT,
  redirect: '/space/index',
  meta: {
    hideChildrenInMenu: true,
    icon: 'ant-design:appstore-add-outlined',
    title: t('routes.basic.space'),
    orderNo: 100000,
  },
  children: [
    {
      path: 'index',
      name: 'SpaceIndex',
      component: () => import('/@/views/space/index.vue'),
      meta: {
        title: t('routes.basic.space'),
        icon: 'ant-design:appstore-add-outlined',
        hideMenu: true,
      },
    },
  ],
};

export default environment;
