import { mount, flushPromises } from "@vue/test-utils";
import UserTable from "@/components/UserTable.vue";
import axios from "@/api/axios";
import ElementPlus, { ElNotification } from "element-plus";

jest.mock("@/api/axios");
jest.mock("element-plus", () => ({
  ElNotification: jest.fn(),
}));

jest.mock("@/components/UserEditModal.vue", () => ({
  name: "UserEditModal",
  template: "<div />",
}));
jest.mock("@/components/UserAddModal.vue", () => ({
  name: "UserAddModal",
  template: "<div />",
}));
jest.mock("@/components/Paginator.vue", () => ({
  name: "Paginator",
  template: "<div />",
}));

const usersMock = [
  {
    userID: "1",
    firstName: "John",
    lastName: "Doe",
    age: 30,
    email: "john@example.com",
  },
];

beforeEach(() => {
  (axios.get as jest.Mock).mockResolvedValueOnce({
    data: { data: usersMock, total: 50 },
  });
});

describe("UserTable.vue", () => {
  it("loads users and displays them in the table on mount", async () => {
    const wrapper = mount(UserTable);
    await flushPromises();

    expect(axios.get).toHaveBeenCalledTimes(1);
    expect(wrapper.findAll(".el-table__row").length).toEqual(usersMock.length);
  });

  it("opens UserAddModal on add button click", async () => {
    const wrapper = mount(UserTable);
    await wrapper.find(".add--btn--wrapper .el-button").trigger("click");
    expect(wrapper.findComponent({ name: "UserAddModal" }).isVisible()).toBe(
      true,
    );
  });

  it("calls deleteUser API and removes a user on delete button click", async () => {
    (axios.delete as jest.Mock).mockResolvedValueOnce({});
    const wrapper = mount(UserTable);
    await flushPromises();

    await wrapper.findAll(".el-button--danger")[0].trigger("click");
    await flushPromises();

    expect(axios.delete).toHaveBeenCalledTimes(1);
    expect(ElNotification).toHaveBeenCalledWith(
      expect.objectContaining({ type: "success" }),
    );
  });

  it("handles search input and filters users", async () => {
    const wrapper = mount(UserTable);
    const input = wrapper.find(".user--table header input");
    await input.setValue("John");
    await flushPromises();

    expect(axios.get).toHaveBeenCalledWith(expect.stringContaining("John"));
  });

  it("changes pagination and fetches users accordingly", async () => {
    const wrapper = mount(UserTable);
    await wrapper
      .findComponent({ name: "Paginator" })
      .vm.$emit("pagination", { page: 2, limit: 10 });
    await flushPromises();

    expect(axios.get).toHaveBeenCalledWith(expect.stringContaining("page=2"));
    expect(axios.get).toHaveBeenCalledWith(
      expect.stringContaining("pageSize=10"),
    );
  });
});
