import {
  ListReq,
  CreateReq,
  ListItemRes, ListItem
} from './model';
import { defHttp } from '/@/utils/http/axios';

enum Api {
  Deploy = '/deploy',
  DeployId = '/deploy/{id}',
  DeployStart = '/deploy/{id}/release',
  DeployAudit = '/deploy/{id}/audit',
  DeployConsoleWs = 'ws://localhost:8989/api/deploy/{id}/console',
}

export const getDeployListByPage = (params?: ListReq) =>
  defHttp.get<ListItemRes>({ url: Api.Deploy, params });

export const createDeploy = (params: CreateReq) =>
  defHttp.post({url: Api.Deploy, params: params});

export const deleteDeploy = (id: number) =>
  defHttp.delete({ url: Api.DeployId.replace('{id}', id.toString()) });

export const detailDeploy = (id: number, notAlertErrMsg: boolean|undefined) =>
  defHttp.get<ListItem>({url: Api.DeployId.replace('{id}', id.toString())}, notAlertErrMsg ? {errorMessageMode: "none"} : {});

export const startDeploy = (id: number, notAlertErrMsg: boolean|undefined) =>
  defHttp.get<ListItem>({url: Api.DeployStart.replace('{id}', id.toString())}, notAlertErrMsg ? {errorMessageMode: "none"} : {});

export const auditDeploy = (id: number, audit:boolean) =>
  defHttp.post<ListItem>({url: Api.DeployAudit.replace('{id}', id.toString()), params:{audit:audit}}, );

export const getDeployConsoleWs = (id: number) =>
  Api.DeployConsoleWs.replace('{id}', id.toString())
