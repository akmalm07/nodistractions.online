export type User = {
  id: string;
  email: string;
  name: string;
};

const baseURL = import.meta.env.VITE_API_URL ?? 'http://localhost:3000/api/v1';

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(baseURL + path, {
    ...init,
    credentials: 'include',
    headers: {
      'content-type': 'application/json',
      ...init?.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new Error(error.error ?? 'Request failed');
  }

  return response.status === 204 ? (undefined as T) : response.json();
}

export const api = {
  register: (body: { email: string; password: string; name: string }) =>
    request<{ user: User }>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(body),
    }),
  login: (email: string, password: string) =>
    request<{ user: User }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),
  logout: () => request<void>('/auth/logout', { method: 'POST' }),
  me: () => request<User>('/users/me'),
  deleteAccount: () => request<void>('/users/me', { method: 'DELETE' }),
};
