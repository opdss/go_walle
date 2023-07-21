import { BasicPageParams, BasicFetchResult } from '/@/api/model/baseModel';
import {ListItem as EnvListItem} from '/@/api/environment/model'
import {ListItem as ProjectListItem} from '/@/api/project/model'
import {ListItem as UserListItem} from '/@/api/user/model'
import {ListItem as ServerListItem} from '/@/api/server/model'

export type ListReq = BasicPageParams & {
  name?: string;
};

export interface CreateReq {
  name: string;
  project_id: number;
  branch?: string;
  tag?: string;
  commit?: string;
  server_ids?: string;
}

export interface ListItem {
  id: number;
  name: string;
  project: ProjectListItem;
  environment: EnvListItem;
  user: UserListItem;
  servers: ServerListItem[];
  commit_id: string;
  status: number;
  description: string;
  created_at: string;
  [key: string]: any;
}

export interface RecordItem {
  id : number;
  type: number;
  status: number;
  server_id: number;
  command: string;
  output: string;
  created_at: string;
  [key: string]:any;
}

export type ListItemRes = BasicFetchResult<ListItem>;
