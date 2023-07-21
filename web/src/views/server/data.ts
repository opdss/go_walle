import { BasicColumn } from '/@/components/Table';
import { FormSchema } from '/@/components/Table';
import {formatToDateTime} from "/@/utils/dateUtil";
export const columns: BasicColumn[] = [
  {
    title: '名称',
    dataIndex: 'name',
  },
  {
    title: '用户',
    dataIndex: 'user',
  },
  {
    title: 'IP',
    dataIndex: 'host',
  },
  {
    title: '端口',
    dataIndex: 'port',
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
  label: '主机id',
  required: true,
  show: false,
  component: 'InputNumber',
  componentProps: {
    disabled: true,
  },
};

export const formSchema: FormSchema[] = [
  {
    field: 'name',
    label: '名称',
    required: true,
    component: 'Input',
  },
  {
    field: 'user',
    label: '用户',
    defaultValue: '',
    required: true,
    component: 'Input',
  },
  {
    field: 'host',
    label: 'IP',
    required: true,
    component: 'Input',
  },
  {
    field: 'port',
    label: '端口',
    defaultValue: 22,
    required: true,
    component: 'InputNumber',
  },
  {
    label: '备注',
    field: 'description',
    component: 'InputTextArea',
  },
];
