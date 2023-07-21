export enum RoleEnum {
  // super admin
  SUPER = 'super',

  // tester
  OWNER = 'owner',
  MASTER = 'master',
  DEVELOPER = 'developer',
}

export const RoleSUPER = [RoleEnum.SUPER]
export const RoleOWNER = [RoleEnum.OWNER, ...RoleSUPER]
export const RoleMaster = [RoleEnum.MASTER, ...RoleOWNER]
export const RoleDEVELOPER = [RoleEnum.DEVELOPER, ...RoleMaster]

export const RoleOptions = [
  {
    text:"超级管理员",
    value:RoleEnum.SUPER
  },
  {
    text:"空间所有者",
    value:RoleEnum.OWNER
  },
  {
    text:"项目管理员",
    value:RoleEnum.MASTER
  },
  {
    text:"开发者",
    value:RoleEnum.DEVELOPER
  }
]
