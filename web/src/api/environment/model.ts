import { BasicPageParams, BasicFetchResult } from '/@/api/model/baseModel';
import {ListItem as SpaceListItem} from '/@/api/space/model'

export type ListReq = BasicPageParams & {
  name?: string;
};

export interface CreateReq {
  name: string;
  status: number;
  description: string;
  color: string;
}

export interface UpdateReq extends CreateReq {
  id: number;
}

export interface ListItem {
  id: number;
  name: string;
  status: number;
  description: string;
  color: string;
  created_at: number;
  space: SpaceListItem;
}

export type ListItemRes = BasicFetchResult<ListItem>;
