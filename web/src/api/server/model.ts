import { BasicPageParams, BasicFetchResult } from '/@/api/model/baseModel';

export type ListReq = BasicPageParams & {
  name?: string;
  host?: string;
};

export interface CreateReq {
  name: string;
  user: string;
  host: string;
  port: number;
  description: string;
}

export interface UpdateReq extends CreateReq {
  id: number;
}

export interface ListItem {
  id: number;
  user: string;
  name: string;
  host: string;
  port: number;
  description: string;
  created_at: string;
}

export interface SetLogin {
  id: number;
  password: string;
}

export type ListItemRes = BasicFetchResult<ListItem>;
