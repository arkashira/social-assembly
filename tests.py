# social-assembly/artifact.py
from typing import Dict, Optional, List
import json
from datetime import datetime

class SocialArtifact:
    def __init__(self, content: str, author: str, metadata: Optional[Dict] = None):
        """
        Initialize a SocialArtifact object.

        Args:
            content (str): The content of the social artifact.
            author (str): The author of the social artifact.
            metadata (Optional[Dict], optional): Additional metadata for the social artifact. Defaults to None.
        """
        self.content = content
        self.author = author
        self.metadata = metadata or {}
        self.timestamp = datetime.utcnow().isoformat()

    def to_dict(self) -> Dict:
        """
        Convert the SocialArtifact object to a dictionary.

        Returns:
            Dict: A dictionary representation of the SocialArtifact object.
        """
        return {
            "content": self.content,
            "author": self.author,
            "metadata": self.metadata,
            "timestamp": self.timestamp
        }

    @classmethod
    def from_dict(cls, data: Dict) -> "SocialArtifact":
        """
        Create a SocialArtifact object from a dictionary.

        Args:
            data (Dict): A dictionary containing the social artifact data.

        Returns:
            SocialArtifact: A SocialArtifact object created from the dictionary.
        """
        return cls(
            content=data["content"],
            author=data["author"],
            metadata=data.get("metadata", {})
        )

    def validate(self) -> bool:
        """
        Basic validation for social artifacts.

        Returns:
            bool: True if the social artifact is valid, False otherwise.
        """
        return bool(self.content and self.author)

# Example usage for API compatibility
if __name__ == "__main__":
    artifact = SocialArtifact("Hello, World!", "John Doe")
    print(artifact.to_dict())
    print(SocialArtifact.from_dict(artifact.to_dict()).validate())