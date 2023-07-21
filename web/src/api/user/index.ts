import { defHttp } from '/@/utils/http/axios';
import { LoginParams, LoginResultModel, GetUserInfoModel, CreateReq, ListItemRes, ListReq, UpdateReq } from './model';
import {GetOptionItemsModel} from "/@/api/model/baseModel"

import { ErrorMessageMode } from '/#/axios';

enum Api {
  Login = '/login',
  Logout = '/logout',
  RefreshToken = '/refresh_token',
  GetUserInfo = '/user_info',
  User = '/user',
  UserId = '/user/{id}',
  UserOptions = '/user/options',
}

/**
 * @description: user login api
 */
export function loginApi(params: LoginParams, mode: ErrorMessageMode = 'modal') {
  return defHttp.post<LoginResultModel>(
    {
      url: Api.Login,
      params,
    },
    {
      errorMessageMode: mode,
    },
  );
}

/**
 * @description: getUserInfo
 */
export function getUserInfo() {
  return defHttp.get<GetUserInfoModel>({ url: Api.GetUserInfo }, { errorMessageMode: 'none' });
}

export function refreshToken(refresh_token: string) {
  return defHttp.post<LoginResultModel>({ url: Api.RefreshToken, params: {refresh_token:refresh_token} });
}

export function doLogout() {
  return defHttp.post({ url: Api.Logout });
}

export const getUserListByPage = (params?: ListReq) =>
  defHttp.get<ListItemRes>({ url: Api.User, params });

export const createUser = (params: CreateReq) =>
  defHttp.post({url: Api.User, params: params});


export const updateUser = (params: UpdateReq) =>
  defHttp.put({ url: Api.User, params: params });

export const deleteUser = (id: number) =>
  defHttp.delete({ url: Api.UserId.replace("{id}", id.toString()) });

export const getUserOptions = (params?: ListReq) =>
  defHttp.get<GetOptionItemsModel>({ url: Api.UserOptions, params });
