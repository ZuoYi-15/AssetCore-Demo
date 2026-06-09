import { reactive } from 'vue';
import { http, unwrap } from './http';
import type { AuthUser, LoginResponse, RegisterPayload } from '../types/api';

const TOKEN_KEY = 'asset_core_token';
const USER_KEY = 'asset_core_user';

export const authState = reactive<{
  token: string;
  user: AuthUser | null;
}>({
  token: localStorage.getItem(TOKEN_KEY) || '',
  user: readUser()
});

function readUser(): AuthUser | null {
  const raw = localStorage.getItem(USER_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as AuthUser;
  } catch {
    localStorage.removeItem(USER_KEY);
    return null;
  }
}

function setSession(token: string, user: AuthUser) {
  authState.token = token;
  authState.user = user;
  localStorage.setItem(TOKEN_KEY, token);
  localStorage.setItem(USER_KEY, JSON.stringify(user));
}

export function clearSession() {
  authState.token = '';
  authState.user = null;
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USER_KEY);
}

export function isLoggedIn() {
  return Boolean(authState.token && authState.user);
}

export function hasPermission(permission: string) {
  return authState.user?.permissions.includes(permission) || false;
}

export function hasRole(role: string) {
  return authState.user?.roles.includes(role) || false;
}

export async function login(username: string, password: string) {
  const res = await http.post('/api/v1/auth/login', { username, password });
  const data = unwrap<LoginResponse>(res.data);
  setSession(data.token, data.user);
  return data;
}

export async function loadProfile() {
  const res = await http.get('/api/v1/auth/me');
  const user = unwrap<AuthUser>(res.data);
  authState.user = user;
  localStorage.setItem(USER_KEY, JSON.stringify(user));
  return user;
}

export async function registerUser(payload: RegisterPayload) {
  const res = await http.post('/api/v1/auth/register', payload);
  return unwrap<AuthUser>(res.data);
}
