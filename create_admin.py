import argparse
import os
import sys
from pathlib import Path

# Ensure we can import app modules
sys.path.append(str(Path(__file__).parent.parent / "src"))

try:
    from app.models.user import User
    from app.db.session import SessionLocal
    from app.core.security import get_password_hash
except ImportError as e:
    print(f"ERROR: Failed to import app modules: {e}", file=sys.stderr)
    sys.exit(1)

def create_admin(email: str, username: str, password: str) -> bool:
    db = SessionLocal()
    try:
        if db.query(User).filter(User.email == email).first():
            print(f"Admin user with email {email} already exists.")
            return False
        hashed_password = get_password_hash(password)
        admin = User(
            email=email,
            username=username,
            hashed_password=hashed_password,
            is_active=True,
            is_superuser=True,
        )
        db.add(admin)
        db.commit()
        print(f"Admin user created: {username} ({email})")
        return True
    except Exception as e:
        print(f"ERROR: Failed to create admin: {e}", file=sys.stderr)
        db.rollback()
        return False
    finally:
        db.close()

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Create admin user.")
    parser.add_argument("--email", required=True, help="Admin email")
    parser.add_argument("--username", required=True, help="Admin username")
    parser.add_argument("--password", required=True, help="Admin password")
    args = parser.parse_args()

    success = create_admin(args.email, args.username, args.password)
    sys.exit(0 if success else 1)