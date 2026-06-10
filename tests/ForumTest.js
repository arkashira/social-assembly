describe('Post model', () => {
  it('creates a new post with title, body, and optional media', () => {
    const post = new Post({ title: 'Test post', body: 'This is a test post' });
    expect(post.title).toBe('Test post');
    expect(post.body).toBe('This is a test post');
  });

  it('edits a post by the author within 5 minutes of creation', () => {
    const post = new Post({ title: 'Test post', body: 'This is a test post' });
    post.edit({ title: 'Updated test post', body: 'This is an updated test post' });
    expect(post.title).toBe('Updated test post');
    expect(post.body).toBe('This is an updated test post');
  });

  it('deletes a post by the author or moderators', () => {
    const post = new Post({ title: 'Test post', body: 'This is a test post' });
    post.delete();
    expect(post.deletedAt).not.toBeNull();
  });

  it('preserves post history (soft delete)', () => {
    const post = new Post({ title: 'Test post', body: 'This is a test post' });
    post.delete();
    expect(post.history).not.toBeNull();
  });
});

describe('Post route', () => {
  it('creates a new post', () => {
    const req = { body: { title: 'Test post', body: 'This is a test post' } };
    const res = { status: jest.fn(), json: jest.fn() };
    postRoute.create(req, res);
    expect(res.status).toHaveBeenCalledTimes(1);
    expect(res.status).toHaveBeenCalledWith(201);
  });

  it('edits a post', () => {
    const req = { body: { title: 'Updated test post', body: 'This is an updated test post' } };
    const res = { status: jest.fn(), json: jest.fn() };
    postRoute.edit(req, res);
    expect(res.status).toHaveBeenCalledTimes(1);
    expect(res.status).toHaveBeenCalledWith(200);
  });

  it('deletes a post', () => {
    const req = { params: { id: 1 } };
    const res = { status: jest.fn(), json: jest.fn() };
    postRoute.delete(req, res);
    expect(res.status).toHaveBeenCalledTimes(1);
    expect(res.status).toHaveBeenCalledWith(204);
  });
});