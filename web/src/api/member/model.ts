import {BasicFetchResult, BasicPageParams} from '/@/api/model/baseModel';
import {RoleEnum} from "/@/enums/roleEnum";

export type ListReq = BasicPageParams & {
};

export interface StoreReq {
  user_id: number;
  role: RoleEnum;
}


export interface ListItem {
  user_id: number;
  username: string;
  email: string;
  role: string;
  status: number;
  created_at: string;
}

export type ListItemRes = BasicFetchResult<ListItem>;
