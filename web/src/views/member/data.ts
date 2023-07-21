import { BasicColumn } from '/@/components/Table';
import {StatusTag} from "/@/enums/fieldEnum";
import {formatToDateTime} from "/@/utils/dateUtil";
export const columns: BasicColumn[] = [
  {
    title: '用户id',
    dataIndex: 'user_id',
  },
  {
    title: '用户名',
    dataIndex: 'username',
  },
  {
    title: '角色',
    dataIndex: 'role',
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
