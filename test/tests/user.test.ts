import {
  describe, it, expect, beforeAll, afterAll, afterEach,
} from '@jest/globals';
import request from 'supertest';

const authUrl = 'http://127.0.0.1:8000';
const phemeUrl = 'http://127.0.0.1:8001';

class User {
  id: number = 0;

  userName: string = '';

  cookie: string = '';

  constructor(id: number, userName: string, cookie: string) {
    this.id = id;
    this.userName = userName;
    this.cookie = cookie;
  }
}

async function createUser(userName: string): Promise<User> {
  // Create user
  await request(authUrl)
    .post('/api/v1/auth/register')
    .send({ name: `${userName}`, email: `${userName}@user.com`, password: 'test' });

  let response = await request(authUrl)
    .post('/api/v1/auth/login')
    .send({ email: `${userName}@user.com`, password: 'test' });
  const cookie = response.get('Set-Cookie')[0];

  response = await request(authUrl)
    .get('/api/v1/auth/user')
    .set('Cookie', cookie);

  return new User(response.body.id, response.body.userName, cookie);
}

async function deleteUser(user: User) {
  await request(authUrl)
    .delete('/api/v1/auth/user')
    .set('Cookie', user.cookie);
}

async function deleteFriend(userA: User, userB: User) {
  await request(phemeUrl)
    .delete(`/api/v1/user/friend/${userB.id}`)
    .set('Cookie', userA.cookie);
}

async function deleteFollower(userA: User, userB: User) {
  await request(phemeUrl)
    .delete(`/api/v1/user/follower/${userB.id}`)
    .set('Cookie', userA.cookie);
}

async function postPheme(user: User): Promise<number> {
  const response = await request(phemeUrl)
    .post('/api/v1/pheme')
    .send({
      visibility: 0, category: 'main', text: 'Hello world!', id: user.id,
    })
    .set('Cookie', user.cookie);

  expect(response.statusCode).toBe(200);
  expect(response.headers['content-type']).toContain('application/json');
  expect(response.body).toHaveProperty('id');

  return response.body.id;
}

let testUser: User = new User(0, '', '');

beforeAll(async () => {
  testUser = await createUser('test.user');
});

afterAll(async () => {
  await deleteUser(testUser);
});

afterEach(async () => {
  const response = await request(phemeUrl)
    .get('/api/v1/pheme/mine')
    .set('Cookie', testUser.cookie);

  const phemes = response.body;
  await Promise.all(phemes.map(async (pheme: any) => {
    await request(phemeUrl)
      .delete(`/api/v1/pheme/${pheme.id}`)
      .set('Cookie', testUser.cookie);
  }));
});

describe('GetUser endpoint', () => {
  it('Missing username', async () => {
    const response = await request(phemeUrl)
      .get('/api/v1/user/wrong');

    expect(response.statusCode).toBe(200);
    expect(response.headers['content-type']).toContain('application/json');
    expect(response.body).toHaveLength(0);
  });

  it('Search single user', async () => {
    const response = await request(phemeUrl)
      .get(`/api/v1/user/${testUser.userName}`);

    expect(response.statusCode).toBe(200);
    expect(response.headers['content-type']).toContain('application/json');
    expect(response.body).toHaveLength(1);
    expect(response.body[0]).toHaveProperty('id');
    expect(response.body[0]).toHaveProperty('userName');
  });

  it('Search multiples users', async () => {
    const name = 'aname';
    const numUsers = 3;
    const users: Array<User> = [];

    await Promise.all(Array.from(Array(numUsers).keys()).map(async (index) => {
      const user = await createUser(`${name}${index}`);
      users.push(user);
    }));

    const response = await request(phemeUrl)
      .get(`/api/v1/user/${name}`);

    expect(response.statusCode).toBe(200);
    expect(response.headers['content-type']).toContain('application/json');
    expect(response.body).toHaveLength(numUsers);

    response.body.forEach((user: Object) => {
      expect(user).toHaveProperty('id');
      expect(user).toHaveProperty('userName');
    });

    await Promise.all(users.map(async (user) => {
      await deleteUser(user);
    }));
  });
});

describe('Friend endpoint', () => {
  it('Wrong user', async () => {
    const response = await request(phemeUrl)
      .put(`/api/v1/user/friend/${testUser.id + 1}`)
      .set('Cookie', testUser.cookie);

    expect(response.statusCode).toBe(400);
    expect(response.headers['content-type']).toContain('application/json');
  });

  it('Add same user', async () => {
    const response = await request(phemeUrl)
      .put(`/api/v1/user/friend/${testUser.id}`)
      .set('Cookie', testUser.cookie);

    expect(response.statusCode).toBe(400);
    expect(response.headers['content-type']).toContain('application/json');
  });

  it('Remove same user', async () => {
    const response = await request(phemeUrl)
      .delete(`/api/v1/user/friend/${testUser.id}`)
      .set('Cookie', testUser.cookie);

    expect(response.statusCode).toBe(400);
    expect(response.headers['content-type']).toContain('application/json');
  });

  it('Add Valid friend', async () => {
    const friend = await createUser('friend');

    const response = await request(phemeUrl)
      .put(`/api/v1/user/friend/${friend.id}`)
      .set('Cookie', testUser.cookie);

    expect(response.statusCode).toBe(200);
    expect(response.headers['content-type']).toContain('application/json');

    await deleteFriend(testUser, friend);

    await deleteUser(friend);
  });
});

describe('Follower endpoint', () => {
  it('Wrong user', async () => {
    const response = await request(phemeUrl)
      .put(`/api/v1/user/follower/${testUser.id + 1}`)
      .set('Cookie', testUser.cookie);

    expect(response.statusCode).toBe(400);
    expect(response.headers['content-type']).toContain('application/json');
  });

  it('Add same user', async () => {
    const response = await request(phemeUrl)
      .put(`/api/v1/user/follower/${testUser.id}`)
      .set('Cookie', testUser.cookie);

    expect(response.statusCode).toBe(400);
    expect(response.headers['content-type']).toContain('application/json');
  });

  it('Remove same user', async () => {
    const response = await request(phemeUrl)
      .delete(`/api/v1/user/follower/${testUser.id}`)
      .set('Cookie', testUser.cookie);

    expect(response.statusCode).toBe(400);
    expect(response.headers['content-type']).toContain('application/json');
  });

  it('Add Valid follower', async () => {
    const friend = await createUser('follower');

    const response = await request(phemeUrl)
      .put(`/api/v1/user/follower/${friend.id}`)
      .set('Cookie', testUser.cookie);

    expect(response.statusCode).toBe(200);
    expect(response.headers['content-type']).toContain('application/json');

    await deleteFollower(testUser, friend);

    await deleteUser(friend);
  });
});
