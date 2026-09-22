import { FormEvent, useEffect, useState } from 'react';

import { api, type User } from './api/client';

export default function App() {
  const [user, setUser] = useState<User | null>(null);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    api.me().then(setUser).catch(() => undefined);
  }, []);

  async function submit(event: FormEvent) {
    event.preventDefault();
    setError('');

    try {
      const result = await api.login(email, password);
      setUser(result.user);
    } catch (error) {
      setError(error instanceof Error ? error.message : 'Request failed');
    }
  }

  if (!user) {
    return (
      <main className="page">
        <section className="card">
          <h1>Sign in</h1>
          <form onSubmit={submit}>
            <label>
              Email
              <input
                required
                type="email"
                value={email}
                onChange={(event) => setEmail(event.target.value)}
              />
            </label>
            <label>
              Password
              <input
                required
                minLength={12}
                type="password"
                value={password}
                onChange={(event) => setPassword(event.target.value)}
              />
            </label>
            {error && <p role="alert">{error}</p>}
            <button>Sign in</button>
          </form>
        </section>
      </main>
    );
  }

  return (
    <main className="page">
      <nav>
        <strong>My SaaS</strong>
        <button onClick={() => api.logout().then(() => setUser(null))}>Logout</button>
      </nav>
      <section className="card">
        <h1>Dashboard</h1>
        <p>Welcome, {user.email}</p>
      </section>
    </main>
  );
}
