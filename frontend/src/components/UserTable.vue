<template>
  <div>
    <div class="add--btn--wrapper">
      <el-button type="primary" :icon="Plus" @click="addModalOpen = true"
        >Add</el-button
      >
    </div>
    <el-table
      :data="users"
      class="user--table"
      v-loading="loading"
      element-loading-background="rgba(122, 122, 122, 0.8)"
    >
      <el-table-column label="ID" prop="userID" />
      <el-table-column label="First Name" prop="firstName" />
      <el-table-column label="Last Name" prop="lastName" />
      <el-table-column label="Age" prop="age" />
      <el-table-column label="Email" prop="email" />
      <el-table-column align="right">
        <template #header>
          <el-input
            v-model="search"
            size="small"
            placeholder="Type to search"
          />
        </template>
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.$index, scope.row)"
            >Edit</el-button
          >
          <el-button
            size="small"
            type="danger"
            @click="handleDelete(scope.$index, scope.row)"
            >Delete</el-button
          >
        </template>
      </el-table-column>
    </el-table>
    <div class="paginator--wrapper">
      <Paginator
        :total="totalPages"
        :page="page"
        :limit="limit"
        @pagination="onChangePagination($event)"
      />
    </div>
    <UserEditModal
      :open="editModalOpen"
      :selectedUser="toRaw(selectedUser)"
      @toggleModal="toggleModal($event)"
      @updateUser="onUpdateUser($event)"
    />

    <UserAddModal
      :open="addModalOpen"
      @toggleModal="toggleModal($event)"
      @createUser="onCreateUser($event)"
    />
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, toRaw, watch } from "vue";
import { onMounted, onUpdated } from "vue";
import { IUser } from "../types";
import { createUser, getUsers, deleteUser, updateUser } from "../api/users";
import { ElNotification } from "element-plus";
import UserEditModal from "./UserEditModal.vue";
import UserAddModal from "./UserAddModal.vue";
import Paginator from "./Paginator.vue";
import debounce from "lodash/debounce";
import dayjs from "dayjs";
import { Plus } from "@element-plus/icons-vue";

const loading = ref(true);
const editModalOpen = ref(false);
const addModalOpen = ref(false);
const users = ref([]);
const page = ref(1);
const search = ref("");
const selectedUser = ref<IUser>();
const limit = ref(5); // Default page size
const totalPages = ref(10);

onMounted(async () => {
  fetchUsers();
});

watch([page, search, limit], () => {
  fetchUsers();
});

const handleEdit = (index: number, row: IUser) => {
  editModalOpen.value = true;
  selectedUser.value = row;
};

const handleDelete = (index: number, row: IUser) => {
  deleteUser(row.id)
    .then(res => {
      users.value.splice(index, 1);
      fetchUsers();
      ElNotification({
        title: "Success",
        message: "User removed successfully",
        type: "success",
      });
    })
    .catch(error => {
      console.log(error);
      ElNotification({
        title: "Error",
        message: error.response?.data.error || error.message,
        type: "error",
      });
    });
};

const fetchUsers = debounce(
  () =>
    getUsers(page.value, limit.value, search.value)
      .then(res => {
        users.value = res.data.data;
        totalPages.value = res.data.total;
      })
      .catch(error => {
        console.log(error);
      })
      .finally(() => {
        loading.value = false;
      }),
  500,
);

const toggleModal = value => {
  value === "add"
    ? (addModalOpen.value = false)
    : (editModalOpen.value = false);
};

const onUpdateUser = user => {
  updateUser(user)
    .then(res => {
      editModalOpen.value = false;
      const index = users.value.findIndex(item => item.id === user.id);
      users.value[index] = { ...user };
      console.log(index);
      ElNotification({
        title: "Success",
        message: "User updated successfully",
        type: "success",
      });
    })
    .catch(error => {
      console.log(error);
      ElNotification({
        title: "Error",
        message: error.response.data.error || error.message,
        type: "error",
      });
    });
};

const onChangePagination = paginationInfo => {
  page.value = paginationInfo.page;
  limit.value = paginationInfo.limit;
};

const onCreateUser = async newUser => {
  createUser(newUser)
    .then(res => {
      users.value.push({ ...newUser, userID: res.data.userID });
      addModalOpen.value = false;
      ElNotification({
        title: "Success",
        message: "User created Successfully",
        type: "success",
      });
    })
    .catch(error => {
      console.log(error);
      ElNotification({
        title: "Error",
        message: error.response.data.error || error.message,
        type: "error",
      });
    });
};
</script>

<style scoped>
.paginator--wrapper {
  display: flex;
  justify-content: center;
  margin-top: 1rem;
}

.add--btn--wrapper {
  text-align: right;
  margin-bottom: 0.5rem;
}

.user--table {
  width: 100%;
  height: 350px;
  border: 1px solid #eff2f7;
}
</style>
