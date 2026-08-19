export interface User {
  id: number;
  username: string;
  full_name?: string;
  role: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
  errors?: unknown;
}

export interface MeResponse {
  id: number;
  username: string;
  role: string;
}