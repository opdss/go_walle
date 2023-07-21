import {Tag} from 'ant-design-vue';
import {h} from 'vue';

export enum Status {
  Enabled = 1, //启用
  Disabled = 2, //禁用
}


export const StatusOptions = [
  {label: '正常', value: Status.Enabled},
  {label: '禁用', value: Status.Disabled},
];
export const StatusTag = (status: number) => {
  const enable = status == Status.Enabled;
  const color = enable ? 'green' : 'red';
  const text = enable ? '正常' : '禁用';
  return h(Tag, {color: color}, () => text);
};


export enum DeployStatus {
  Waiting = 1, //新建提交，等待审核
  Audit = 2, //审核通过
  AuditReject = 3, //审核拒绝
  Release = 4, //上线发布中
  ReleaseFail = 5, //上线失败
  Finish = 6, //上线完成
}

export const DeployStatusShowMsg = (deployStatus:number): [string, string ] => {
  let msg = ""
  let color = "info"
  switch (deployStatus) {
    case DeployStatus.Audit:
      msg = "审核通过，准备发布上线"
      break
    case DeployStatus.Waiting :
      msg = "等待审核中"
      break
    case DeployStatus.AuditReject:
      msg = "审核失败，不允许发布上线"
      color="error"
      break
    case DeployStatus.Release:
      msg = "发布上线中，请等待"
      color="warning"
      break
    case DeployStatus.ReleaseFail:
      msg ="发布失败"
      color="error"
      break
    case DeployStatus.Finish:
      msg="发布成功"
      color="success"
      break
    default:
  }
  return [msg, color]
}
