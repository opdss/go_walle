import type { AppRouteModule } from '/@/router/types';

import { LAYOUT } from '/@/router/constant';
//import { t } from '/@/hooks/web/useI18n';

const project: AppRouteModule = {
  path: '/member',
  name: 'Member',
  component: LAYOUT,
  redirect: '/member/index',
  meta: {
    hideChildrenInMenu: true,
    icon: 'ant-design:user-outlined',
    title: "成员管理",
  },
  children: [
    {
      path: 'index',
      name: 'memberIndex',
      component: () => import('/@/views/member/index.vue'),
      meta: {
        title: '成员管理',
        icon: 'ant-design:user-outlined',
        hideMenu: true,
      },
    },
  ],
};

export default project;
