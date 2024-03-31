import type { IUser } from "@/types";
import {type FormRules} from "element-plus/es";
import { reactive } from "vue";
import * as Yup from "yup"

export const checkAge = (rule: any, value: any, callback: any) => {
  if (!value) {
    return callback(new Error('Please input the age'))
  }
  setTimeout(() => {
    if (!Number.isInteger(value)) {
      callback(new Error('Please input digits'))
    } else {
      if (value < 18) {
        callback(new Error('Age must be greater than 18'))
      } else {
        callback()
      }
    }
  })
}


export const formRules = reactive<FormRules<IUser>>({
  userID: [{ required: true, message: "User ID is required" }],
  firstName: [
    { required: true, message: "Please input first name", trigger: "blur" },
  ],
  lastName: [
    { required: true, message: "Please input last name", trigger: "blur" },
  ],
  age: [
    { required: true, message: "age is required" },
    { type: "number", message: "age must be a number" },
    { validator: checkAge },
  ],
  email: [
    {
      required: true,
      message: "Please input email address",
      trigger: "blur",
    },
    {
      type: "email",
      message: "Please input correct email address",
      trigger: ["blur", "change"],
    },
  ],
});

export const formValidationSchema = Yup.object().shape({
  userID: Yup.string()
    .required('User ID is required'),
  firstName: Yup.string()
      .required('First Name is required'),
  lastName: Yup.string()
      .required('Last name is required'),
  email: Yup.string()
      .required('Email is required')
      .email('Email is invalid'),
  age: Yup.string()
      .min(19, 'Age must be greater than 18')
      .required('Age is required'),
});
