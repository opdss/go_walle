import { BasicColumn } from '/@/components/Table';
import { FormSchema } from '/@/components/Table';
import {getEnvironmentOptions} from "/@/api/environment";
import {getServerListByPage} from "/@/api/server";
export const columns: BasicColumn[] = [
  {
    title: '名称',
    dataIndex: 'name',
  },

  {
    title: '环境',
    dataIndex: 'environment',
    customRender: ({ record }) => record.environment.name,
  },
  {
    title: '空间',
    dataIndex: 'space',
    customRender: ({ record }) => record.space.name,
  },
  {
    title: '备注',
    dataIndex: 'description',
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    width: 180,
  },
];

export const searchFormSchema: FormSchema[] = [
  {
    field: 'name',
    label: '名称',
    component: 'Input',
    colProps: { span: 8 },
  },
];


export const formIdSchema: FormSchema = {
  field: 'id',
  label: 'ID',
  required: true,
  show: false,
  component: 'InputNumber',
  componentProps: {
    disabled: true,
  },
};

export const formSchemas: {[key:string]: FormSchema[]} = {
  "base": [

  ]
}

export const formSchema: FormSchema[] = [
  {
    field: 'divider1',
    component: 'Divider',
    label: '基础字段',
    componentProps: {
      dashed: true,
    },
    colProps: {
      span: 24,
    },
  },
  {
    field: 'name',
    component: 'Input',
    label: '名称',
    required: true,
    colProps: {
    },
  },
  {
    label: '环境',
    field: 'environment_id',
    component: 'ApiSelect',
    defaultValue: 1,
    componentProps: {
      api: getEnvironmentOptions,
      resultField: 'options',
      labelField: 'text',
      valueField: 'value',
    },
    required: true,
    colProps: {
      offset: 2,
    },
  },
  {
    label: '仓库地址',
    field: 'repo_url',
    component: 'Input',
    required: true,
    colProps: {
    },
  },
  {
    label: '上线方式',
    field: 'repo_mode',
    component: 'Select',
    required: true,
    componentProps:{
      options: [
        {
          label: 'Tag',
          value: 'tag',
        },
        {
          label: 'Branch',
          value: 'branch',
        },],
    },
    colProps: {
      offset: 2,
    },
  },


  {
    field: 'divider2',
    component: 'Divider',
    label: '目标机器',
    componentProps: {
      dashed: true,
      type: "vertical",
    },
    colProps: {
      span: 24,
    },
  },
  {
    label: '选择服务器',
    field: 'server_ids',
    component: 'ApiTransfer',
    colProps: {
      span: 16,
      offset: 4
    },
    required: true,
    componentProps: {
      api:getServerListByPage,
      params: {
        page_size: 1000,
      },
      showSearch: true,
      showSelectAll: true,
      resultField: "items",
      labelField: "name",
      valueField: "id",
    },
  },
  {
    label: '目标集群部署路径',
    field: 'target_root',
    component: 'Input',
    required: true,
    colProps: {
      span:6
    },
  },
  {
    label: '目标集群部署仓库',
    field: 'target_releases',
    component: 'Input',
    required: true,
    colProps: {
      span:6,
      offset: 2,
    },
  },
  {
    label: '目标集群部署仓库版本保留数',
    field: 'keep_version_num',
    component: 'InputNumber',
    required: true,
    colProps: {
      span:6,
      offset: 2,
    },
  },

  {
    label: '任务配置',
    field: 'divider3',
    component: 'Divider',
    componentProps: {
      dashed: true,
    },
    colProps: {
      span: 24,
    },
  },
  {
    label: '排除文件',
    field: 'excludes',
    component: 'InputTextArea',
    colProps: {
      span: 10
    },
  },
  {
    label: '全局变量',
    field: 'task_vars',
    component: 'InputTextArea',
    colProps: {
      span: 10,
      offset: 2,
    },
  },
  {
    label: '高级任务-Deploy前置任务',
    field: 'prev_deploy',
    component: 'InputTextArea',
    colProps: {
      span: 10
    },
  },
  {
    label: '高级任务-Deploy后置任务',
    field: 'post_deploy',
    component: 'InputTextArea',
    colProps: {
      span: 10,
      offset: 2,
    },
  },
  {
    label: '高级任务-Release前置任',
    field: 'prev_release',
    component: 'InputTextArea',
    colProps: {
      span: 10
    },
  },
  {
    label: '高级任务-Release后置任务',
    field: 'post_release',
    component: 'InputTextArea',
    colProps: {
      span: 10,
      offset: 2,
    },
  },
];
