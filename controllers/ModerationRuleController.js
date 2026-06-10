// moderationRule.model.test.js
const ModerationRule = require('../models/ModerationRule');
const mongoose = require('mongoose');

describe('ModerationRule Model', () => {
  afterAll(() => mongoose.disconnect());

  test('creates a rule with version 1 and timestamps', async () => {
    const rule = await ModerationRule.create({
      title: 'No Spam',
      description: 'Spam posts are prohibited',
      severity: 'high',
    });
    expect(rule.version).toBe(1);
    expect(rule.createdAt).toBeInstanceOf(Date);
    expect(rule.updatedAt).toBeInstanceOf(Date);
    expect(rule.archived).toBeFalsy();
  });

  test('rejects creation with missing title', async () => {
    await expect(
      ModerationRule.create({ severity: 'low' })
    ).rejects.toThrow(/title.*required/);
  });

  test('rejects creation with invalid severity', async () => {
    await expect(
      ModerationRule.create({ title: 'Test', severity: 'critical' })
    ).rejects.toThrow(/severity.*enum/);
  });

  test('editing creates a new version, preserves old version', async () => {
    const original = await ModerationRule.create({
      title: 'No Hate',
      severity: 'medium',
    });
    const edited = await ModerationRule.edit(original.id, {
      title: 'No Hate Speech',
    });
    expect(edited.version).toBe(2);
    expect(edited.updatedAt).toBeInstanceOf(Date);
    // fetch historic version
    const history = await ModerationRule.find({ _id: original.id }).sort('version');
    expect(history).toHaveLength(2);
    expect(history[0].title).toBe('No Hate');
    expect(history[1].title).toBe('No Hate Speech');
  });

  test('archiving sets archived flag without changing version', async () => {
    const rule = await ModerationRule.create({
      title: 'No Off‑topic',
      severity: 'low',
    });
    const archived = await ModerationRule.archive(rule.id);
    expect(archived.archived).toBeTruthy();
    expect(archived.version).toBe(rule.version);
  });

  test('timestamps are stored in UTC and not future dated', async () => {
    const now = new Date();
    const rule = await ModerationRule.create({
      title: 'UTC Test',
      severity: 'low',
    });
    expect(rule.createdAt.getTime()).toBeLessThanOrEqual(now.getTime());
    expect(rule.updatedAt.getTime()).toBeLessThanOrEqual(now.getTime());
    expect(rule.createdAt.getTimezoneOffset()).toBe(0);
  });
});