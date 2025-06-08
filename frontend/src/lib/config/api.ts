
import axios from "axios";
import type { AxiosInstance, AxiosResponse,AxiosError,
  InternalAxiosRequestConfig, } from "axios";

const BASE_URL = "http://localhost:8080/"; 

const api: AxiosInstance = axios.create({
  baseURL: BASE_URL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});



export default api;

