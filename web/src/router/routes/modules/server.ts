import type { AppRouteModule } from '/@/router/types';

import { LAYOUT } from '/@/router/constant';
import { t } from '/@/hooks/web/useI18n';

const server: AppRouteModule = {
  path: '/server',
  name: 'Server',
  component: LAYOUT,
  redirect: '/server/index',
  meta: {
    hideChildrenInMenu: true,
    icon: 'ant-design:hdd-outlined',
    title: t('routes.basic.server'),
    orderNo: 100000,
  },
  children: [
    {
      path: 'index',
      name: 'ServerIndex',
      component: () => import('/@/views/server/index.vue'),
      meta: {
        title: t('routes.basic.server'),
        icon: 'ant-design:hdd-outlined',
        hideMenu: true,
      },
    },
    {
      path: 'terminal/:id/:server',
      name: 'ServerTerminal',
      component: () => import('/@/views/server/terminal.vue'),
      meta: {
        title: "终端",
        icon: 'ant-design:hdd-outlined',
        hideMenu: true,
      },
    },
  ],
};

export default server;
