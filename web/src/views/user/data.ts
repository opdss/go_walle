import { BasicColumn } from '/@/components/Table';
import { FormSchema } from '/@/components/Table';
import {Status, StatusOptions, StatusTag} from "/@/enums/fieldEnum";
import {formatToDateTime} from "/@/utils/dateUtil";
export const columns: BasicColumn[] = [
  {
    title: '用户id',
    dataIndex: 'id',
  },
  {
    title: '用户名',
    dataIndex: 'username',
  },
  {
    title: '邮箱',
    dataIndex: 'email',
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

export const searchFormSchema: FormSchema[] = [
  {
    field: 'key_word',
    label: '名称',
    component: 'Input',
    colProps: { span: 8 },
  },
];


export const formIdSchema: FormSchema = {
  field: 'id',
  label: '用户id',
  required: true,
  show: false,
  component: 'InputNumber',
  componentProps: {
    disabled: true,
  },
};

export const formSchema: FormSchema[] = [
  {
    field: 'username',
    label: '用户名',
    required: true,
    component: 'Input',
  },
  {
    field: 'email',
    label: '邮箱',
    required: true,
    component: 'Input',
  },
  {
    field: 'password',
    label: '密码',
    required: true,
    component: 'InputPassword',
  },
  {
    field: 'status',
    label: '状态',
    component: 'RadioButtonGroup',
    defaultValue: Status.Enabled,
    required: true,
    componentProps: {
      options: StatusOptions,
    },
  },
];
