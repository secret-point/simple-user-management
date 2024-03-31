<template>
  <div>
    <el-dialog
      v-model="modalOpen"
      title="Edit User"
      width="500"
      @close="closeModal()"
      align-center
    >
      <el-form
        ref="ruleFormRef"
        style="max-width: 600px"
        :model="ruleForm"
        :rules="rules"
        label-width="auto"
        class="demo-ruleForm"
        size="large"
        status-icon
      >
        <el-form-item label="User ID" prop="userID">
          <el-input
            v-model.number="ruleForm.userID"
            type="text"
            autocomplete="off"
            placeholder="30"
          />
        </el-form-item>
        <el-form-item label="First Name" prop="firstName">
          <el-input v-model="ruleForm.firstName" placeholder="John" />
        </el-form-item>
        <el-form-item label="Last Name" prop="lastName">
          <el-input v-model="ruleForm.lastName" placeholder="Doe" />
        </el-form-item>
        <el-form-item label="age" prop="age">
          <el-input
            v-model.number="ruleForm.age"
            type="text"
            autocomplete="off"
            placeholder="30"
          />
        </el-form-item>
        <el-form-item prop="email" label="email">
          <el-input v-model="ruleForm.email" placeholder="john@example.com" />
        </el-form-item>
        <div class="dialog-footer">
          <el-form-item>
            <el-button
              type="primary"
              @click="submitForm(ruleFormRef)"
              :disabled="ruleForm.o"
            >
              Add
            </el-button>
            <el-button @click="closeModal()">Candel</el-button>
          </el-form-item>
        </div>
      </el-form>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, reactive, toRaw, onMounted, watchEffect } from "vue";
import { FormInstance, ElNotification } from "element-plus";
import FormRules from "element-plus/es";
import { createUser, getUsers, deleteUser, updateUser } from "../api/users";
import { IUser } from "../types";
import { formRules as rules } from "../contsants";

const emit = defineEmits(["toggleModal", "createUser"]);
const props = defineProps(["open", "selectedUser"]);
const ruleFormRef = ref<FormInstance>();
const ruleForm = reactive<IUser>({
  ...props.selectedUser,
});

watchEffect(() => {
  if (props.selectedUser) {
    Object.assign(ruleForm, props.selectedUser);
  } else {
    Object.assign(ruleForm, {
      firstName: "",
      lastName: "",
      email: "",
    });
  }
});

const submitForm = async (formEl: FormInstance) => {
  console.log(formEl);
  if (!formEl) return;
  await formEl.validate((valid, fields) => {
    console.log(valid);
    if (valid) {
      emit("createUser", toRaw(ruleForm));
    }
  });
};

const modalOpen = computed({
  get() {
    return props.open;
  },
  set(value) {},
});

const closeModal = () => {
  emit("toggleModal", "edit");
};
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: end;
}
</style>
