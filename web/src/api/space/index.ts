import {
  ListReq,
  CreateReq,
  UpdateReq,
  ListItemRes
} from './model';
import { defHttp } from '/@/utils/http/axios';

enum Api {
  Space = '/space',
  SpaceId = '/space/{id}',
}

export const getSpaceListByPage = (params?: ListReq) =>
  defHttp.get<ListItemRes>({ url: Api.Space, params });

export const createSpace = (params: CreateReq) =>
  defHttp.post({url: Api.Space, params: params});


export const updateSpace = (params: UpdateReq) =>
  defHttp.put({ url: Api.Space, params: params });

export const deleteSpace = (id: number) =>
  defHttp.delete({ url: Api.SpaceId.replace('{id}', id.toString()) });

