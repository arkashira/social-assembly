# /opt/axentx/social-assembly/test_auth.py

import unittest
import os
import sqlite3
from auth import Auth

class TestAuth(unittest.TestCase):
    def setUp(self):
        self.db_path = 'test_social_assembly.db'
        self.auth = Auth(self.db_path)
        self.username = 'testuser'
        self.password = 'testpass'

    def tearDown(self):
        if os.path.exists(self.db_path):
            os.remove(self.db_path)

    def test_register_success(self):
        success, error = self.auth.register(self.username, self.password)
        self.assertTrue(success)
        self.assertIsNone(error)

    def test_register_failure_empty_fields(self):
        success, error = self.auth.register('', '')
        self.assertFalse(success)
        self.assertEqual(error, "Username and password are required")

    def test_register_failure_username_exists(self):
        self.auth.register(self.username, self.password)
        success, error = self.auth.register(self.username, 'anotherpass')
        self.assertFalse(success)
        self.assertEqual(error, "Username already exists")

    def test_login_success(self):
        self.auth.register(self.username, self.password)
        self.assertTrue(self.auth.login(self.username, self.password))
