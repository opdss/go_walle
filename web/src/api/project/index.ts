import {
  ListReq,
  CreateReq,
  UpdateReq,
  ListItemRes,
  BranchItems,
  TagItems,
  CommitItems, Detail,
} from './model';
import { defHttp } from '/@/utils/http/axios';
import {GetOptionItemsModel} from "../model/baseModel";

enum Api {
  Project = '/project',
  ProjectId = '/project/{id}',
  ProjectOptions = '/project/options',
  ProjectBranches = '/project/{id}/branches',
  ProjectTags = '/project/{id}/tags',
  ProjectCommits = '/project/{id}/commits',
}

export const getProjectListByPage = (params?: ListReq) =>
  defHttp.get<ListItemRes>({ url: Api.Project, params });

export const createProject = (params: CreateReq) =>
  defHttp.post({url: Api.Project, params: params});


export const updateProject = (params: UpdateReq) =>
  defHttp.put({ url: Api.Project, params: params });

export const deleteProject = (id: number) =>
  defHttp.delete({ url: Api.ProjectId.replace('{id}', id.toString()) });

export const detailProject = (id: number) =>
  defHttp.get<Detail>({ url: Api.ProjectId.replace('{id}', id.toString()) });


export const getProjectOptions = (params?: ListReq) =>
  defHttp.get<GetOptionItemsModel>({ url: Api.ProjectOptions, params });

export const getProjectBranches = (id: number) =>
  defHttp.get<BranchItems[]>({ url: Api.ProjectBranches.replace('{id}', id.toString()) });

export const getProjectTags = (id: number) =>
  defHttp.get<TagItems[]>({ url: Api.ProjectTags.replace('{id}', id.toString()) });

export const getProjectCommits = (id: number, branch?:string) =>
  defHttp.get<CommitItems[]>({ url: Api.ProjectCommits.replace('{id}', id.toString()), params: {"branch": branch} });
