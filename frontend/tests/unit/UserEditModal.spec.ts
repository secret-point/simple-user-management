import { mount } from "@vue/test-utils";
import UserEditModal from "@/components/UserEditModal.vue";
import flushPromises from "flush-promises";
import axios from "@/api/axios";

describe("UserEditModal.vue", () => {
  const selectedUser = {
    userID: "1",
    firstName: "John",
    lastName: "Doe",
    age: 30,
    email: "john@example.com",
  };

  it("renders and reacts to the `open` prop correctly", async () => {
    const wrapper = mount(UserEditModal, {
      props: { open: false, selectedUser },
    });
    expect(wrapper.vm.modalOpen).toBe(false);
    await wrapper.setProps({ open: true });
    expect(wrapper.vm.modalOpen).toBe(true);
  });

  it("closes the modal and emits event on close button click", async () => {
    const wrapper = mount(UserEditModal, {
      props: { open: true, selectedUser },
    });
    await wrapper.find(".el-button--default").trigger("click");
    expect(wrapper.emitted()).toHaveProperty("toggleModal");
  });

  it('validates form inputs and emits "updateUser" with correct data on form submission', async () => {
    const wrapper = mount(UserEditModal, {
      props: { open: true, selectedUser },
    });
    await wrapper.find(".el-form").trigger("submit.prevent");
    await flushPromises();
    expect(axios.put).toHaveBeenCalled();
    expect(wrapper.emitted()).toHaveProperty("updateUser");

    const updateUserEventPayload = wrapper.emitted("updateUser")[0][0];
    expect(updateUserEventPayload).toMatchObject(selectedUser);
  });

  it("displays validation errors for empty required fields", async () => {
    const wrapper = mount(UserEditModal, {
      props: { open: true, selectedUser: {} },
    });
    await wrapper.find(".el-form").trigger("submit.prevent");
    await flushPromises();
    expect(wrapper.findAll(".el-form-item__error").length).toBeGreaterThan(0);
    expect(axios.put).not.toHaveBeenCalled();
  });
});
