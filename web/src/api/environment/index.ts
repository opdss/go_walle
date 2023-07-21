import {
  ListReq,
  CreateReq,
  UpdateReq,
  ListItemRes
} from './model';
import { defHttp } from '/@/utils/http/axios';
import {GetOptionItemsModel} from "../model/baseModel";

enum Api {
  Environment = '/environment',
  EnvironmentId = '/environment/{id}',
  EnvironmentOptions = '/environment/options',
}

export const getEnvironmentListByPage = (params?: ListReq) =>
  defHttp.get<ListItemRes>({ url: Api.Environment, params });

export const createEnvironment = (params: CreateReq) =>
  defHttp.post({url: Api.Environment, params: params});


export const updateEnvironment = (params: UpdateReq) =>
  defHttp.put({ url: Api.Environment, params: params });

export const deleteEnvironment = (id: number) =>
  defHttp.delete({ url: Api.EnvironmentId.replace('{id}', id.toString()) });

export const getEnvironmentOptions = (params?: ListReq) =>
  defHttp.get<GetOptionItemsModel>({ url: Api.EnvironmentOptions, params });
