from fastapi import APIRouter, HTTPException, status
from typing import List
from src.models.user import User, UserProfileCreate, UserProfileUpdate

router = APIRouter(prefix="/profiles", tags=["profiles"])

# In-memory storage for demonstration purposes
users_db = {}

@router.post("/", response_model=User)
def create_profile(profile: UserProfileCreate):
    """Create a new user profile"""
    # Check if user already exists
    if profile.username in users_db:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="User with this username already exists"
        )
    
    # Create new user
    user_id = f"user_{len(users_db) + 1}"
    new_user = User(
        id=user_id,
        username=profile.username,
        email=profile.email,
        display_name=profile.display_name,
        bio=profile.bio,
        avatar_url=profile.avatar_url,
        created_at=datetime.now(),
        updated_at=datetime.now()
    )
    
    users_db[profile.username] = new_user
    return new_user

@router.get("/{username}", response_model=User)
def get_profile(username: str):
    """Get a user profile by username"""
    if username not in users_db:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User not found"
        )
    return users_db[username]

@router.put("/{username}", response_model=User)
def update_profile(username: str, profile_update: UserProfileUpdate):
    """Update a user profile"""
    if username not in users_db:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User not found"
        )
    
    # Update user
    user = users_db[username]
    if profile_update.display_name is not None:
        user.display_name = profile_update.display_name
    if profile_update.bio is not None:
        user.bio = profile_update.bio
    if profile_update.avatar_url is not None:
        user.avatar_url = profile_update.avatar_url
    
    user.updated_at = datetime.now()
    users_db[username] = user
