// Test: inviteModel.test.js
describe('Invite Model', () => {
  beforeEach(async () => {
    await Invite.deleteMany({});
  });

  test('should create invite with valid data', async () => {
    const inviteData = {
      token: 'abc123def456',
      communityId: 'community123',
      createdBy: 'user123'
    };
    
    const invite = new Invite(inviteData);
    const savedInvite = await invite.save();
    
    expect(savedInvite.token).toBe(inviteData.token);
    expect(savedInvite.communityId).toBe(inviteData.communityId);
    expect(savedInvite.createdBy).toBe(inviteData.createdBy);
    expect(savedInvite.createdAt).toBeDefined();
  });

  test('should validate required fields', async () => {
    const invite = new Invite({});
    let error;
    
    try {
      await invite.save();
    } catch (err) {
      error = err;
    }
    
    expect(error).toBeDefined();
    expect(error.errors.token).toBeDefined();
    expect(error.errors.communityId).toBeDefined();
    expect(error.errors.createdBy).toBeDefined();
  });
});

// Test: inviteRoutes.test.js
describe('Invite Routes', () => {
  let server;
  let validToken;

  beforeAll(async () => {
    server = app.listen(0);
    // Setup valid JWT token
    validToken = jwt.sign({ userId: 'testUser' }, process.env.JWT_SECRET);
  });

  afterAll(async () => {
    await server.close();
  });

  describe('POST /invite/generate', () => {
    test('should return 401 for unauthorized request', async () => {
      const response = await request(app)
        .post('/invite/generate')
        .send({ communityId: 'testCommunity' });
      
      expect(response.status).toBe(401);
      expect(response.body.error).toBe('Unauthorized');
    });

    test('should return 400 for missing communityId', async () => {
      const response = await request(app)
        .post('/invite/generate')
        .set('Authorization', `Bearer ${validToken}`);
      
      expect(response.status).toBe(400);
      expect(response.body.error).toBe('Community ID is required');
    });

    test('should generate invite with valid request', async () => {
      const response = await request(app)
        .post('/invite/generate')
        .set('Authorization', `Bearer ${validToken}`)
        .send({ communityId: 'testCommunity' });
      
      expect(response.status).toBe(201);
      expect(response.body.token).toHaveLength(greaterThan(31));
      expect(response.body.token).toHaveLength(lessThan(65));
      expect(response.body.communityId).toBe('testCommunity');
    });
  });
});