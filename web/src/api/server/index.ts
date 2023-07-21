import {
  ListReq,
  CreateReq,
  UpdateReq,
  ListItemRes, SetLogin,
} from './model';
import { defHttp } from '/@/utils/http/axios';

enum Api {
  Server = '/server',
  ServerId = '/server/{id}',
  ServerSshWs = 'ws://localhost:8989/api/server/{id}/terminal',
  ServerCheck =  "/server/{id}/check",
  ServerSetLogin =  "/server/setting",
}

export const getServerListByPage = (params?: ListReq) =>
  defHttp.get<ListItemRes>({ url: Api.Server, params });

export const createServer = (params: CreateReq) =>
  defHttp.post({url: Api.Server, params: params});


export const updateServer = (params: UpdateReq) =>
  defHttp.put({ url: Api.Server, params: params });

export const deleteServer = (id: number) =>
  defHttp.delete({ url: Api.ServerId.replace('{id}', id.toString()) });

export const checkServer = (id: number) =>
  defHttp.post({ url: Api.ServerCheck.replace('{id}', id.toString()) }, {errorMessageMode:"none"});

export const setLoginServer = (params:SetLogin) =>
  defHttp.post({ url: Api.ServerSetLogin, params: params}, {errorMessageMode:"none"});

export const getServerSshWs = (id: number) =>
  Api.ServerSshWs.replace('{id}', id.toString())
