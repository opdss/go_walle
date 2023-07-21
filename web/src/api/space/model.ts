import { BasicPageParams, BasicFetchResult } from '/@/api/model/baseModel';
import {ListItem as UserListItem} from '/@/api/user/model'

export type ListReq = BasicPageParams & {
  name?: string;
};

export interface CreateReq {
  user_id: number;
  name: string;
  status: number;
}

export interface UpdateReq extends CreateReq {
  id: number;
}

export interface ListItem {
  id: number;
  user_id: number;
  name: string;
  status: number;
  created_at: number;
  user: UserListItem;
  [key: string] : any;
}

export type ListItemRes = BasicFetchResult<ListItem>;
