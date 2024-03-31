import axios from "axios";

// Create a mock instance
const mockAxios = axios.create();

// Mock methods
mockAxios.get = jest.fn(() =>
  Promise.resolve({
    data: {
      data: [
        {
          userID: "1",
          firstName: "Jane",
          lastName: "Doe",
          age: 28,
          email: "jane.doe@example.com",
        },
        {
          userID: "2",
          firstName: "John",
          lastName: "Smith",
          age: 30,
          email: "john.smith@example.com",
        },
      ],
    },
  }),
);
mockAxios.post = jest.fn(() =>
  Promise.resolve({ data: { message: "User updated successfully" } }),
);
mockAxios.put = jest.fn((url, data) =>
  Promise.resolve({ data: { ...data, id: url.split("/").pop() } }),
);
mockAxios.delete = jest.fn(() =>
  Promise.resolve({ data: { message: "User removed successfully" } }),
);

export default mockAxios;
