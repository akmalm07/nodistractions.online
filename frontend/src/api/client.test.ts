import { afterEach, describe, expect, it, vi } from 'vitest';

import { api } from './client';

afterEach(() => vi.unstubAllGlobals());

describe('API client', () => {
  it('uses cookie credentials for protected requests', async () => {
    const fetch = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ id: 'u1', email: 'person@example.test', name: 'Person' }), {
        status: 200,
      }),
    );
    vi.stubGlobal('fetch', fetch);

    await expect(api.me()).resolves.toMatchObject({ id: 'u1' });
    expect(fetch.mock.calls[0][1]).toMatchObject({ credentials: 'include' });
  });

  it('surfaces API error messages', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(
        new Response(JSON.stringify({ error: 'Authentication required.' }), { status: 401 }),
      ),
    );

    await expect(api.me()).rejects.toThrow('Authentication required.');
  });
});
