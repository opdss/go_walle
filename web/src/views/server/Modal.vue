<template>
  <BasicModal v-bind="$attrs" @register="registerModal" :title="getTitle" @ok="handleSubmit">
    <BasicForm @register="registerForm" />
  </BasicModal>
</template>
<script lang="ts">
import { defineComponent, ref, computed, unref } from 'vue';
import { BasicModal, useModalInner } from '/@/components/Modal';
import { BasicForm, useForm } from '/@/components/Form/index';
import { formSchema , formIdSchema} from './data';
import { createServer, updateServer} from '/@/api/server';
import { useMessage } from '/@/hooks/web/useMessage';
export default defineComponent({
  name: 'Modal',
  components: { BasicModal, BasicForm },
  emits: ['success', 'register'],
  setup(_, { emit }) {
    const isUpdate = ref(true);
    const hasAppendIdForm = ref(false);
    const { createMessage} = useMessage()

    const [
      registerForm,
      {
        resetFields,
        removeSchemaByField,
        appendSchemaByField,
        setFieldsValue,
        validate,
      },
    ] = useForm({
      labelWidth: 100,
      baseColProps: { span: 24 },
      schemas: formSchema,
      showActionButtonGroup: false,
    });

    const [registerModal, { setModalProps , closeModal}] = useModalInner(async (data) => {
      await resetFields();
      setModalProps({ confirmLoading: false });
      isUpdate.value = !!data?.isUpdate;

      if (unref(isUpdate)) {
        if (!unref(hasAppendIdForm)) {
          await appendSchemaByField(formIdSchema, undefined, true);
          hasAppendIdForm.value = true;
        }
        await setFieldsValue({
          ...data.record,
        });
      } else {
        if (unref(hasAppendIdForm)) {
          await removeSchemaByField('id');
          hasAppendIdForm.value = false;
        }
      }
    });

    const getTitle = computed(() => (!unref(isUpdate) ? '新增' : '编辑'));

    async function handleSubmit() {
      try {
        const values = await validate();
        setModalProps({ confirmLoading: true });
        if (unref(isUpdate)) {
          updateServer(values)
            .then(() => {
              createMessage.success("更新成功")
              emit('success');
            })
        } else {
          createServer(values)
            .then(() => {
              createMessage.success("新建成功")
              emit('success');
            })
        }
        closeModal()
      } finally {
        setModalProps({ confirmLoading: false });
      }
    }

    return { registerModal, registerForm, getTitle, handleSubmit };
  },
});
</script>
