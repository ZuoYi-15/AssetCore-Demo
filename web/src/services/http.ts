import axios from 'axios';
import { ElMessage } from 'element-plus';
import type { ApiResponse } from '../types/api';

export const http = axios.create({
  baseURL: '',
  timeout: 15000
});

http.interceptors.response.use(
  (response) => {
    const body = response.data as ApiResponse<unknown>;
    if (body && typeof body.code === 'number' && body.code !== 0) {
      ElMessage.error(body.message || '请求失败');
      return Promise.reject(new Error(body.message));
    }
    return response;
  },
  (error) => {
    ElMessage.error(error?.response?.data?.message || error.message || '网络请求失败');
    return Promise.reject(error);
  }
);

export function unwrap<T>(body: ApiResponse<T>): T {
  return body.data;
}
