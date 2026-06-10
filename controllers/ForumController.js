// Implement CRUD operations for forums
const Forum = require('../models/Forum');

exports.createForum = async (req, res) => {
  try {
    const forum = new Forum(req.body);
    await forum.save();
    res.status(201).json(forum);
  } catch (err) {
    res.status(500).json({ message: err.message });
  }
};

exports.getForums = async (req, res) => {
  try {
    const forums = await Forum.find().exec();
    res.json(forums);
  } catch (err) {
    res.status(500).json({ message: err.message });
  }
};