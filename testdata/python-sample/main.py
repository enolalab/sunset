"""Main application module."""

from models.user import User, UserRole
from utils.helper import validate_email


# Default admin email
DEFAULT_ADMIN = "admin@example.com"


def create_user(name: str, email: str, role: UserRole = "user") -> User:
    """Create a new user with validation.

    Args:
        name: The user's display name.
        email: The user's email address.
        role: The user's role, defaults to 'user'.

    Returns:
        A new User instance.

    Raises:
        ValueError: If the email is invalid.
    """
    if not validate_email(email):
        raise ValueError(f"Invalid email: {email}")
    return User(name=name, email=email, role=role)


class UserService:
    """Service for managing users."""

    def __init__(self, db_url: str):
        """Initialize with database connection."""
        self.db_url = db_url
        self._users: list[User] = []

    def get_user(self, user_id: str) -> User | None:
        """Find a user by their ID."""
        for user in self._users:
            if user.id == user_id:
                return user
        return None

    @staticmethod
    def hash_password(password: str) -> str:
        """Hash a password for storage."""
        return f"hashed_{password}"
