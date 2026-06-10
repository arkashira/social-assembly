// Define the Forum model
const mongoose = require('mongoose');

const forumSchema = new mongoose.Schema({
  title: String,
  description: String,
  posts: [{ type: mongoose.Schema.Types.ObjectId, ref: 'Post' }]
});

const Forum = mongoose.model('Forum', forumSchema);

module.exports = Forum;