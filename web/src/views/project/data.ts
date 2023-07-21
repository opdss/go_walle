import { BasicColumn } from '/@/components/Table';
import { FormSchema } from '/@/components/Table';
import {getEnvironmentOptions} from "/@/api/environment";
import {getServerListByPage} from "/@/api/server";
// import {h} from "vue"
// import { CodeEditor } from '/@/components/CodeEditor';
import {formatToDateTime} from "/@/utils/dateUtil";
//import { MarkDown } from '/@/components/MarkDown';
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
    format: (val) => formatToDateTime(val)
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
  "基础配置": [
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
      componentProps: {
        api: getEnvironmentOptions,
        resultField: 'options',
        labelField: 'text',
        valueField: 'value',
      },
      required: true,
    },
    {
      label: '仓库类型',
      field: 'repo_type',
      component: 'RadioButtonGroup',
      required: true,
      defaultValue: 'git',
      componentProps:{
        options: [
          {
            label: 'Git',
            value: 'git',
          },
          {
            label: 'Svn',
            value: 'svn',
          },],
      },
      colProps: {
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
      component: 'RadioButtonGroup',
      required: true,
      defaultValue:"tag",
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
    },
    {
      label: '是否审核',
      field: 'task_audit',
      component: 'Switch',
      required: true,
      defaultValue:1,
      componentProps:{
        checkedChildren:"是",
        unCheckedChildren:"否",
        checkedValue: 1,
        unCheckedValue: 2,
      },
      colProps: {
      },
    },
    {
      field: 'description',
      component: 'InputTextArea',
      label: '项目简介',
      colProps: {
      },
    },
  ],
  '服务器配置' : [
    {
      label: '选择服务器',
      field: 'server_ids',
      component: 'ApiTransfer',
      required: true,
      componentProps: {
        api:getServerListByPage,
        params: {
          page_size: 1000,
        },
        showSearch: true,
        showSelectAll: true,
        // resultCallback: (res) => {
        //   if (res.items?.length > 0) {
        //     res.items.forEach((val) => {
        //       if (val.status == 2) {
        //         val.disabled = true
        //       }
        //     })
        //   }
        //   return res.items
        // },
        resultField: "items",
        labelField: "name",
        valueField: "id",
        listStyle: {
          width: '280px',
          height: '320px',
        },
      },
    },
    {
      label: '目标集群部署仓库',
      field: 'target_releases',
      component: 'Input',
      helpMessage: '部署服务器存放每个版本的代码的目录，target_root将会建立一个软链接指向到这里的某个版本目录',
      required: true,
    },
    {
      label: '目标集群部署路径',
      field: 'target_root',
      component: 'Input',
      helpMessage: '部署服务器目标路径，该值对应比如nginx的root，是一个软链接，配置时不应该存在，由发布系统自动创建',
      required: true,
    },
    {
      label: '目标集群部署仓库版本保留数',
      field: 'keep_version_num',
      component: 'InputNumber',
      defaultValue: 5,
      required: true,
    },
  ],
  "高级配置" : [
    {
      label: '排除文件',
      field: 'excludes',
      component: 'InputTextArea',
    },
    {
      label: '全局变量',
      field: 'task_vars',
      component: 'InputTextArea',
    },
    {
      label: '高级任务-Deploy前置任务',
      field: 'prev_deploy',
      component: 'InputTextArea',
    },
    {
      label: '高级任务-Deploy后置任务',
      field: 'post_deploy',
      component: 'InputTextArea',
    },
    {
      label: '高级任务-Release前置任',
      field: 'prev_release',
      component: 'InputTextArea',
    },
    {
      label: '高级任务-Release后置任务',
      field: 'post_release',
      component: 'InputTextArea',
    },
    // {
    //   label: '高级任务-Release后置任务',
    //   field: 'post_release',
    //   component: 'InputTextArea',
    //   required: true,
    //   render: ({ model, field }) => {
    //     return h(CodeEditor, {
    //       mode:"shell",
    //       value: model[field],
    //       onChange: (value: string) => {
    //         model[field] = value;
    //       },
    //     });
    //   },
    // },
  ]
}

