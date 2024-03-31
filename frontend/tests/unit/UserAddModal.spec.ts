import { mount } from "@vue/test-utils";
import UserAddModal from "@/components/UserAddModal.vue";
jest.mock("@/api/axios", () => require("@/__mocks__/api/axios"));

describe("UserAddModal.vue", () => {
  it("opens the modal when the `open` prop is true", async () => {
    const wrapper = mount(UserAddModal, { props: { open: true } });
    expect(wrapper.find(".el-dialog__wrapper").isVisible()).toBe(true);
  });

  it("closes the modal on clicking the cancel button", async () => {
    const wrapper = mount(UserAddModal, { props: { open: true } });
    await wrapper.find('[aria-label="Close"]').trigger("click");
    expect(wrapper.emitted("update:open")).toBeTruthy();
  });

  it("submits valid form data correctly", async () => {
    const wrapper = mount(UserAddModal, { props: { open: true } });
    await wrapper.find('input[name="firstName"]').setValue("Jane");
    await wrapper.find('input[name="lastName"]').setValue("Doe");
    await wrapper.find('input[name="age"]').setValue(25);
    await wrapper.find('input[name="email"]').setValue("jane.doe@example.com");
    await wrapper.find("form").trigger("submit.prevent");

    expect(axios.post).toHaveBeenCalledWith("/user", {
      firstName: "Jane",
      lastName: "Doe",
      age: 25,
      email: "jane.doe@example.com",
    });
    expect(wrapper.emitted("createUser")).toBeTruthy();
  });

  it("shows validation errors for empty required fields", async () => {
    const wrapper = mount(UserAddModal, { props: { open: true } });
    await wrapper.find("form").trigger("submit.prevent");
    await wrapper.vm.$nextTick();
    expect(wrapper.find(".el-form-item__error").exists()).toBe(true);
  });
});
