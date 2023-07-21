<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" @click="handleCreate"> 新增</a-button>
      </template>
      <template #expandedRowRender="{ record }">
        <span> 备注：{{ record.description }} </span>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="recordActions(record)"
          />
        </template>
      </template>
    </BasicTable>
  </div>
</template>
<script lang="ts">
import {defineComponent} from 'vue';
import {ActionItem, BasicTable, TableAction, useTable} from '/@/components/Table';
import {columns, searchFormSchema} from './data';
import {useMessage} from '/@/hooks/web/useMessage';
import {deleteDeploy, getDeployListByPage, auditDeploy} from "/@/api/deploy"
import {DeployStatus} from "/@/enums/fieldEnum"
import {useGo} from "/@/hooks/web/usePage";

export default defineComponent({
  name: 'DeployManagement',
  components: {BasicTable, TableAction},
  setup() {
    const {createMessage} = useMessage();
    const go = useGo();
    const [registerTable, {reload}] = useTable({
      title: '部署管理',
      api: getDeployListByPage,
      columns,
      formConfig: {
        labelWidth: 120,
        schemas: searchFormSchema,
      },
      useSearchForm: true,
      showTableSetting: true,
      //bordered: true,
      showIndexColumn: false,
      expandRowByClick: true,
      actionColumn: {
        title: '操作',
        dataIndex: 'action',
        // slots: { customRender: 'action' },
        fixed: "right",
      },
    });

    function handleCreate() {
      go("/deploy/create")
    }


    function handleDelete(record: Recordable) {
      deleteDeploy(record.id).then(() => {
        createMessage.success("删除成功")
        reload()
      })
    }

    function handleAudit(record: Recordable, audit:boolean) {
      auditDeploy(record.id, audit).then(() => {
        reload()
      })
    }

    function recordActions(record: Recordable): ActionItem[] {
      const actions: { [key: string]: ActionItem } = {
        "release": {
          label: "上线",
          onClick: () => {
            go(`/deploy/release/${record.id}`)
          }
        },
        "rollback": {
          label: "回滚",
          color: 'error',
          onClick: () => {
            go(`/deploy/release/${record.id}`)
          }
        },
        "detail": {
          label: "查看",
          onClick: () => {
            go(`/deploy/release/${record.id}`)
          }
        },
        "cancel": {
          label: "取消发布",
          onClick: () => {
            go(`/deploy/release/${record.id}`)
          }
        },
        "audit": {
          label: "审核",
          popConfirm: {
            title: '是否确认发布该上线单？',
            placement: 'left',
            confirm: handleAudit.bind(null, record, true),
            cancelText: "拒绝",
            cancel:  handleAudit.bind(null, record, false),
          },
        },
        "del": {
          icon: 'ant-design:delete-outlined',
          color: 'error',
          popConfirm: {
            title: '是否确认删除',
            placement: 'left',
            confirm: handleDelete.bind(null, record),
          },
        }
      }
      let acts:string[] = []
      switch (record.status) {
        case DeployStatus.Waiting:
          acts = ['audit', 'cancel']
          break
        case DeployStatus.Audit:
          acts = ['release']
          break
        case DeployStatus.AuditReject:
          acts = ['detail']
          break
        case DeployStatus.Release:
          acts = ['cancel', 'detail']
          break
        case DeployStatus.ReleaseFail:
          acts = ['detail']
          break
        case DeployStatus.Finish:
          acts = ['detail']
          break
      }
      let res:ActionItem[]= []
      acts.forEach((v, _) => {
        res.push(actions[v])
      })
      return res
    }


    return {
      registerTable,
      handleCreate,
      handleDelete,
      recordActions,
    };
  },
});
</script>
