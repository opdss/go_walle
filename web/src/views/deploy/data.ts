import { BasicColumn } from '/@/components/Table';
import { FormSchema } from '/@/components/Table';
import {getProjectBranches, getProjectCommits, getProjectTags} from "/@/api/project"
import {DeployStatus} from "/@/enums/fieldEnum"
import {h} from "vue";
import {Tag} from "ant-design-vue";
import {formatToDateTime} from "/@/utils/dateUtil";

export const columns: BasicColumn[] = [
  {
    title: '上线单',
    dataIndex: 'name',
  },
  {
    title: '提交用户',
    dataIndex: '_user',
    customRender: ({ record }) => record.user.username,
  },
  {
    title: '上线版本',
    dataIndex: '_version',
    customRender: ({ record }) => {
      return record.tag || `${record.branch}@${record.commit_id.substring(0,8)}`
    },
  },
  {
    title: '发布环境',
    dataIndex: '_env',
    customRender: ({ record }) => record.environment.name,
  },
  {
    title: '项目',
    dataIndex: '_project',
    customRender: ({ record }) => record.project.name,
  },
  {
    title: '状态',
    dataIndex: 'status',
    customRender: ({ record }) => {
      let text = ""
      let color = "blue"
      switch (record.status) {
        case DeployStatus.Waiting:
            text='等待审核'
          break
        case DeployStatus.Audit:
          text='审核通过'
          break
        case DeployStatus.AuditReject:
          text='审核驳回'
          color='yellow'
          break
        case DeployStatus.Release:
          text='发布中...'
          break
        case DeployStatus.ReleaseFail:
          text='上线失败'
          color='red'
          break
        case DeployStatus.Finish:
          text='发布完成'
          color='green'
          break
      }
      return h(Tag, {color: color}, () => text);
    },
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


export const formSchema: FormSchema[] = [
  {
    field: 'name',
    label: '上线单标题',
    required: true,
    component: 'Input',
  },
  {
    field: 'tag',
    label: '选取Tag',
    component: 'ApiSelect',
    componentProps: ({formModel}) => {
      if (!formModel["project_id"]) {
        return {}
      }
      return {
        api: getProjectTags,
        params: formModel["project_id"],
        resultField: '',
        immediate: false,
        alwaysLoad: true,
        labelField: 'name',
        valueField: 'name',
      }
    },
    required: true,
  },
  {
    field: 'branch',
    label: '选取分支',
    component: 'ApiSelect',
    componentProps: ({formModel}) => {
      if (!formModel["project_id"]) {
        return {}
      }
      return {
        api: getProjectBranches,
        params: formModel["project_id"],
        immediate: false,
        alwaysLoad: true,
        resultField: '',
        labelField: 'name',
        valueField: 'name',
        showSearch: true,
        onChange: () => {
          formModel.commit_id = ""
        },
      }
    },
    required: true,
  },
  {
    label: '选取版本',
    field: 'commit_id',
    component: 'ApiSelect',
    componentProps: ({formModel}) => {
      if (!formModel["project_id"] || !formModel['branch']) {
        return {}
      }
      return {
        api: (p) => getProjectCommits(p[0], p[1]),
        params: [formModel["project_id"], formModel['branch']],
        immediate: false,
        //alwaysLoad: true,
        resultField: '',
        labelField: 'name',
        valueField: 'hash',
        showSearch: true,
      }
    },
    required: true,
  },
  {
    field: 'description',
    label: '上线说明',
    component: 'InputTextArea',
  },
  {
    label: '选取服务器',
    field: 'server_ids',
    helpMessage:'默认是项目绑定的所有服务器，也可选择部分发布',
    component: 'CheckboxGroup',
    required: true,
    componentProps: {
      options:[],
    }
  },
];
