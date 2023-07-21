import { BasicPageParams, BasicFetchResult } from '/@/api/model/baseModel';
import {ListItem as serverListItem} from '/@/api/server/model'

export type ListReq = BasicPageParams & {
  name?: string;
  environment_id?: number;
};

export interface CreateReq {
  name: string;
  environment_id: number;
  repo_url: string;
  repo_mode: string;
  repo_type:string;
  task_audit: number;
  description: string;

  server_ids: number[];

  target_root: string;
  target_releases: string;
  keep_version_num: number;

  excludes : string;
  is_include: number;
  task_vars : string;
  prev_deploy : string;
  post_deploy : string;
  prev_release : string;
  post_release : string;

}

export interface UpdateReq extends CreateReq {
  id: number;
}

export interface Detail {
  id: number;
  name: string;
  environment_id: number;
  repo_url: string;
  repo_mode: string;
  repo_type:string;
  task_audit: number;
  description: string;

  servers:serverListItem[];

  target_root: string;
  target_releases: string;
  keep_version_num: number;

  excludes : string;
  is_include: number;
  task_vars : string;
  prev_deploy : string;
  post_deploy : string;
  prev_release : string;
  post_release : string;

  server_ids: number[],
  space: string;
  environment: number;
  created_at: string;
  [key: string] : any
}

export interface ListItem {
  id: number;
  name: string;
  space: string;
  environment: number;
  description: string;
  created_at: string;
  [key: string] : any
}

export type ListItemRes = BasicFetchResult<ListItem>;


export type TagItems = {
  name: string,
  hash: string,
}

export type BranchItems = TagItems

export type CommitItems = TagItems & {
  timestamp: string,
  message: string,
}
