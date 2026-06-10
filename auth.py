# /opt/axentx/social-assembly/auth.py

import bcrypt
import sqlite3
from typing import Optional, Tuple

class Auth:
    def __init__(self, db_path: str = 'social_assembly.db'):
        self.db_path = db_path
        self._init_db()

    def _init_db(self):
        with sqlite3.connect(self.db_path) as conn:
            cursor = conn.cursor()
            cursor.execute('''
                CREATE TABLE IF NOT EXISTS users (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    username TEXT UNIQUE NOT NULL,
                    password_hash TEXT NOT NULL
                )
            ''')
            conn.commit()

    def register(self, username: str, password: str) -> Tuple[bool, Optional[str]]:
        if not username or not password:
            return False, "Username and password are required"

        password_hash = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt())
        try:
            with sqlite3.connect(self.db_path) as conn:
                cursor = conn.cursor()
                cursor.execute('INSERT INTO users (username, password_hash) VALUES (?, ?)', (username, password_hash))
                conn.commit()
            return True, None
        except sqlite3.IntegrityError:
            return False, "Username already exists"

    def login(self, username: str, password: str) -> bool:
        with sqlite3.connect(self.db_path) as conn:
            cursor = conn.cursor()
            cursor.execute('SELECT password_hash FROM users WHERE username = ?', (username,))
            result = cursor.fetchone()
            if result:
                stored_hash = result[0]
                return bcrypt.checkpw(password.encode('utf-8'), stored_hash)
            return False