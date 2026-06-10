// Define API endpoints for forum interactions
const express = require('express');
const router = express.Router();
const forumController = require('../controllers/ForumController');

router.post('/forums', forumController.createForum);
router.get('/forums', forumController.getForums);

module.exports = router;