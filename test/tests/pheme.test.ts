import {
  describe, it, expect, beforeAll, afterAll, afterEach,
} from '@jest/globals';
import request from 'supertest';

const authUrl = 'http://127.0.0.1:8000';
const phemeUrl = 'http://127.0.0.1:8001';
let cookie = '';
let userID = 0;

async function postPheme(): Promise<number> {
  const response = await request(phemeUrl)
    .post('/api/v1/pheme')
    .send({
      visibility: 0, category: 'main', text: 'Hello world!', userID,
    })
    .set('Cookie', cookie);

  expect(response.statusCode).toBe(200);
  expect(response.headers['content-type']).toContain('application/json');
  expect(response.body).toHaveProperty('id');

  return response.body.id;
}

beforeAll(async () => {
  // Create user
  await request(authUrl)
    .post('/api/v1/auth/register')
    .send({ name: 'test.pheme', email: 'test.pheme@test.com', password: 'test' });

  let response = await request(authUrl)
    .post('/api/v1/auth/login')
    .send({ email: 'test.pheme@test.com', password: 'test' });
  cookie = response.get('Set-Cookie')[0];

  response = await request(authUrl)
    .get('/api/v1/auth/user')
    .set('Cookie', cookie);
  userID = response.body.id;
});

afterAll(async () => {
  await request(authUrl)
    .delete('/api/v1/auth/user')
    .set('Cookie', cookie);
});

afterEach(async () => {
  const response = await request(phemeUrl)
    .get('/api/v1/pheme/mine')
    .set('Cookie', cookie);

  const phemes = response.body;
  await Promise.all(phemes.map(async (pheme: any) => {
    await request(phemeUrl)
      .delete(`/api/v1/pheme/${pheme.id}`)
      .set('Cookie', cookie);
  }));
});

describe('PostPheme endpoint', () => {
  it('new pheme missing category', async () => {
    const response = await request(phemeUrl)
      .post('/api/v1/pheme')
      .send({ visibility: 0, text: 'Hello world!', userID })
      .set('Cookie', cookie);

    expect(response.statusCode).toBe(400);
    expect(response.headers['content-type']).toContain('application/json');
  });

  it('new pheme missing text', async () => {
    const response = await request(phemeUrl)
      .post('/api/v1/pheme')
      .send({ visibility: 0, category: 'main', userID })
      .set('Cookie', cookie);

    expect(response.statusCode).toBe(400);
    expect(response.headers['content-type']).toContain('application/json');
  });

  it('new pheme missing userID', async () => {
    const response = await request(phemeUrl)
      .post('/api/v1/pheme')
      .send({ visibility: 0, category: 'main', text: 'Hello world!' })
      .set('Cookie', cookie);

    expect(response.statusCode).toBe(400);
    expect(response.headers['content-type']).toContain('application/json');
  });

  it('new pheme', async () => {
    const response = await request(phemeUrl)
      .post('/api/v1/pheme')
      .send({
        visibility: 0, category: 'main', text: 'Hello world!', userID,
      })
      .set('Cookie', cookie);
    expect(response.statusCode).toBe(200);
    expect(response.headers['content-type']).toContain('application/json');
    expect(response.body).toHaveProperty('id');
  });
});

describe('UpdatePheme endpoint', () => {
  let phemeID = -1;

  beforeEach(async () => {
    const response = await request(phemeUrl)
      .post('/api/v1/pheme')
      .send({
        visibility: 0, category: 'main', text: 'Hello world!', userID,
      })
      .set('Cookie', cookie);
    expect(response.statusCode).toBe(200);
    expect(response.headers['content-type']).toContain('application/json');

    phemeID = response.body.id;
  });

  it('update pheme without id', async () => {
    const response = await request(phemeUrl)
      .put('/api/v1/pheme')
      .send({
        visibility: 0, category: 'main', text: 'Hello world!', userID,
      })
      .set('Cookie', cookie);

    expect(response.statusCode).toBe(405);
  });

  it('update pheme without category', async () => {
    const response = await request(phemeUrl)
      .put(`/api/v1/pheme/${phemeID}`)
      .send({
        visibility: 0, text: 'Hello world!', userID,
      })
      .set('Cookie', cookie);

    expect(response.statusCode).toBe(400);
  });

  it('update pheme without text', async () => {
    const response = await request(phemeUrl)
      .put(`/api/v1/pheme/${phemeID}`)
      .send({
        visibility: 0, category: 'main', userID,
      })
      .set('Cookie', cookie);

    expect(response.statusCode).toBe(400);
  });

  it('update pheme without user ID', async () => {
    const response = await request(phemeUrl)
      .put(`/api/v1/pheme/${phemeID}`)
      .send({
        visibility: 0, category: 'main', text: 'Hello world!',
      })
      .set('Cookie', cookie);

    expect(response.statusCode).toBe(400);
  });

  it('update pheme', async () => {
    const response = await request(phemeUrl)
      .put(`/api/v1/pheme/${phemeID}`)
      .send({
        visibility: 0, category: 'main', text: 'Hello world!', userID,
      })
      .set('Cookie', cookie);

    expect(response.statusCode).toBe(200);
  });
});

describe('GetPheme endpoint', () => {
  const numPhemesToPost = 5;

  it('get wrong posted pheme', async () => {
    const lastPhemeID = await postPheme();

    const response = await request(phemeUrl)
      .get(`/api/v1/pheme/${lastPhemeID + 1}`)
      .set('Cookie', cookie);

    expect(response.statusCode).toBe(400);
    expect(response.headers['content-type']).toContain('application/json');
  });

  it('get posted pheme', async () => {
    const lastPhemeID = await postPheme();

    const response = await request(phemeUrl)
      .get(`/api/v1/pheme/${lastPhemeID}`)
      .set('Cookie', cookie);

    expect(response.statusCode).toBe(200);
    expect(response.headers['content-type']).toContain('application/json');
    expect(response.body).toHaveProperty('id');
  });

  it('get posted phemes', async () => {
    const postedPhemes: Number[] = [];

    await Promise.all(Array(numPhemesToPost).fill(0).map(async () => {
      const phemeID = await postPheme();
      postedPhemes.push(phemeID);
    }));

    await Promise.all(postedPhemes.map(async (id) => {
      const response = await request(phemeUrl)
        .get(`/api/v1/pheme/${id}`)
        .set('Cookie', cookie);

      expect(response.statusCode).toBe(200);
      expect(response.headers['content-type']).toContain('application/json');
      expect(response.body).toHaveProperty('id');
    }));
  });

  it('get posted phemes', async () => {
    const postedPhemes: Number[] = [];

    await Promise.all(Array(numPhemesToPost).fill(0).map(async () => {
      const phemeID = await postPheme();
      postedPhemes.push(phemeID);
    }));

    const response = await request(phemeUrl)
      .get('/api/v1/pheme/mine')
      .set('Cookie', cookie);

    expect(response.statusCode).toBe(200);
    expect(response.headers['content-type']).toContain('application/json');
    expect(Array.isArray(response.body)).toBe(true);
    expect(response.body.length).toBe(numPhemesToPost);
  });
});

describe('DeletePheme endpoint', () => {
  const numPhemesToPost = 5;

  it('delete wrong posted pheme', async () => {
    const lastPhemeID = await postPheme();

    const response = await request(phemeUrl)
      .delete(`/api/v1/pheme/${lastPhemeID + 10}`)
      .set('Cookie', cookie);

    expect(response.statusCode).toBe(400);
    expect(response.headers['content-type']).toContain('application/json');
  });

  it('delete posted pheme', async () => {
    const lastPhemeID = await postPheme();

    const response = await request(phemeUrl)
      .delete(`/api/v1/pheme/${lastPhemeID}`)
      .set('Cookie', cookie);

    expect(response.statusCode).toBe(200);
    expect(response.headers['content-type']).toContain('application/json');
    expect(response.body).toHaveProperty('id');
  });

  it('delete posted phemes', async () => {
    const postedPhemes: Number[] = [];

    await Promise.all(Array(numPhemesToPost).fill(0).map(async () => {
      const phemeID = await postPheme();
      postedPhemes.push(phemeID);
    }));

    await Promise.all(postedPhemes.map(async (id) => {
      const response = await request(phemeUrl)
        .delete(`/api/v1/pheme/${id}`)
        .set('Cookie', cookie);

      expect(response.statusCode).toBe(200);
      expect(response.headers['content-type']).toContain('application/json');
      expect(response.body).toHaveProperty('id');
    }));

    const response = await request(phemeUrl)
      .get('/api/v1/pheme/mine')
      .set('Cookie', cookie);

    expect(response.statusCode).toBe(200);
    expect(response.headers['content-type']).toContain('application/json');
    expect(Array.isArray(response.body)).toBe(true);
    expect(response.body.length).toBe(0);
  });
});
