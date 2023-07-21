import type { AppRouteModule } from '/@/router/types';

import { LAYOUT } from '/@/router/constant';
import { t } from '/@/hooks/web/useI18n';

const deploy: AppRouteModule = {
  path: '/deploy',
  name: 'Deploy',
  component: LAYOUT,
  redirect: '/deploy/index',
  meta: {
    hideChildrenInMenu: true,
    icon: 'ant-design:code-outlined',
    title: t('routes.basic.deploy'),
  },
  children: [
    {
      path: 'index',
      name: 'DeployIndex',
      component: () => import('/@/views/deploy/index.vue'),
      meta: {
        title: t('routes.basic.deploy'),
        icon: 'ant-design:code-outlined',
        hideMenu: true,
      },
    },
    {
      path: 'create',
      name: 'DeployCreate',
      component: () => import('/@/views/deploy/create.vue'),
      meta: {
        title: "新建上线单",
        icon: 'ant-design:code-outlined',
        hideMenu: true,
      },
    },
    {
      path: 'create/:project_id',
      name: 'DeployCreateForm',
      component: () => import('/@/views/deploy/create_form.vue'),
      meta: {
        title: "新建上线单",
        icon: 'ant-design:code-outlined',
        hideMenu: true,
      },
    },
    {
      path: 'release/:id',
      name: 'DeployRelease',
      component: () => import('/@/views/deploy/release.vue'),
      meta: {
        title: "控制台",
        icon: 'ant-design:code-outlined',
        hideMenu: true,
      },
    },
    {
      path: 'process/:task_id',
      name: 'DeployProcessForm',
      component: () => import('/@/views/deploy/process.vue'),
      meta: {
        title: "发布中",
        icon: 'ant-design:code-outlined',
        hideMenu: true,
      },
    },
  ],
};

export default deploy;
