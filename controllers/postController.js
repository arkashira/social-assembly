// Test post creation with valid data
test('createPost should create a post with valid data', async () => {
  const postData = {
    title: 'Test Post',
    body: 'This is a test post.',
    media: 'test.jpg'
  };
  const response = await request(app)
    .post('/api/posts')
    .send(postData)
    .expect(201);

  expect(response.body.title).toBe(postData.title);
  expect(response.body.body).toBe(postData.body);
  expect(response.body.media).toBe(postData.media);
});

// Test post creation with missing title
test('createPost should return 400 if title is missing', async () => {
  const postData = {
    body: 'This is a test post.',
    media: 'test.jpg'
  };
  await request(app)
    .post('/api/posts')
    .send(postData)
    .expect(400);
});

// Test post creation with missing body
test('createPost should return 400 if body is missing', async () => {
  const postData = {
    title: 'Test Post',
    media: 'test.jpg'
  };
  await request(app)
    .post('/api/posts')
    .send(postData)
    .expect(400);
});

// Test post creation with empty title
test('createPost should return 400 if title is empty', async () => {
  const postData = {
    title: '',
    body: 'This is a test post.',
    media: 'test.jpg'
  };
  await request(app)
    .post('/api/posts')
    .send(postData)
    .expect(400);
});

// Test post creation with empty body
test('createPost should return 400 if body is empty', async () => {
  const postData = {
    title: 'Test Post',
    body: '',
    media: 'test.jpg'
  };
  await request(app)
    .post('/api/posts')
    .send(postData)
    .expect(400);
});

// Test post creation with valid data and no media
test('createPost should create a post with valid data and no media', async () => {
  const postData = {
    title: 'Test Post',
    body: 'This is a test post.'
  };
  const response = await request(app)
    .post('/api/posts')
    .send(postData)
    .expect(201);

  expect(response.body.title).toBe(postData.title);
  expect(response.body.body).toBe(postData.body);
  expect(response.body.media).toBeNull();
});