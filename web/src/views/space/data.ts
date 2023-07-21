import { BasicColumn } from '/@/components/Table';
import { FormSchema } from '/@/components/Table';
import { StatusTag, StatusOptions, Status } from '/@/enums/fieldEnum';
import {getUserOptions} from "/@/api/user";
import {formatToDateTime} from "/@/utils/dateUtil";

export const columns: BasicColumn[] = [
  {
    title: '空间名称',
    dataIndex: 'name',
  },
  {
    title: '所属负责人',
    dataIndex: 'user',
    customRender: ({ record }) => record.user.username,
  },
  {
    title: '状态',
    dataIndex: 'status',
    customRender: ({ record }) => StatusTag(record.status),
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    format: (val) => formatToDateTime(val)
  },
];

export const formIdSchema: FormSchema = {
  field: 'id',
  label: '空间ID',
  show: false,
  required: true,
  component: 'InputNumber',
};

export const formSchema: FormSchema[] = [
  {
    field: 'name',
    label: '空间名称',
    required: true,
    component: 'Input',
  },
  {
    field: 'status',
    label: '状态',
    component: 'RadioButtonGroup',
    defaultValue: Status.Enabled,
    helpMessage:"选择禁用后，该空间所有操作将不被允许",
    componentProps: {
      options: StatusOptions,
    },
  },
  {
    label: '所属负责人',
    field: 'user_id',
    component: 'ApiSelect',
    componentProps: {
      api: getUserOptions,
      resultField: 'options',
      labelField: 'text',
      valueField: 'value',
    },
    required: true,
  },
];
