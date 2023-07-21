import { BasicColumn } from '/@/components/Table';
import { FormSchema } from '/@/components/Table';
import {Badge} from 'ant-design-vue'
import { h } from 'vue';
import { StatusTag, StatusOptions, Status } from '/@/enums/fieldEnum';
import {formatToDateTime} from "/@/utils/dateUtil";
export const columns: BasicColumn[] = [
  {
    title: '名称',
    dataIndex: 'name',
    customRender: ({ record }) => {
      return h(Badge, {color:record.color, text:record.name})
    },
  },
  {
    title: '状态',
    dataIndex: 'status',
    customRender: ({ record }) => StatusTag(record.status),
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

export const formSchema: FormSchema[] = [
  {
    field: 'name',
    label: '名称',
    required: true,
    component: 'Input',
  },
  {
    field: 'status',
    label: '状态',
    component: 'RadioButtonGroup',
    defaultValue: Status.Enabled,
    componentProps: {
      options: StatusOptions,
    },
  },
  {
    label: '备注',
    field: 'description',
    component: 'InputTextArea',
  },
];
