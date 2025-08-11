
import axios from "axios";
import type { AxiosInstance, AxiosResponse,AxiosError,
  InternalAxiosRequestConfig, } from "axios";

const BASE_URL = import.meta.env.VITE_BASE_URL;

const api: AxiosInstance = axios.create({
  baseURL: BASE_URL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});



export default api;

