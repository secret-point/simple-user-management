import http from "./axios";
import { type IUser } from "../types";

export const getUsers = async (page: number, pageSize: number, search: string) => {
  return await http.get(`/users?page=${page}&pageSize=${pageSize}&search=${search}`);
};

export const deleteUser = async (id: string) => {
  return await http.delete(`/user/${id}`);
};

export const createUser = async (userInfo: Omit<IUser, "userID">) => {
  return await http.post("/user", userInfo);
};

export const updateUser = async (userInfo: IUser) => {
  return await http.put(`/user/${userInfo.id}`, userInfo);
};
