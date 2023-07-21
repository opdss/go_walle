import { defHttp } from '/@/utils/http/axios';
import { ListReq, StoreReq, ListItemRes} from './model';


enum Api {
  Member = '/member',
}

export const getMemberListByPage = (params?: ListReq) =>
  defHttp.get<ListItemRes>({ url: Api.Member, params });

export const storeMember = (params: StoreReq) =>
  defHttp.post({url: Api.Member, params: params});

export const deleteMember = (id: number) =>
  defHttp.delete({ url: Api.Member.replace("{id}", id.toString()) });
