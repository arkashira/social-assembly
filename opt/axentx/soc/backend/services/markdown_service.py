from markdown import markdown
from flask import jsonify
from html import escape

class MarkdownService:
    @staticmethod
    def render_markdown(content):
        """Renders markdown content to HTML with proper escaping and styling.

        Args:
            content (str): The markdown content to be rendered

        Returns:
            tuple: (response dict, status code)
        """
        if not content:
            return jsonify({"error": "Content cannot be empty"}), 400

        try:
            # Basic HTML escaping for security
            safe_content = escape(content)

            # Convert markdown to HTML
            html_content = markdown(safe_content)

            # Wrap in a styled div for better presentation
            styled_html = f'<div class="markdown-content">{html_content}</div>'

            return jsonify({"html": styled_html}), 200

        except Exception as e:
            return jsonify({"error": f"Markdown rendering failed: {str(e)}"}), 500

    @staticmethod
    def render_header(text, level):
        """Renders a markdown header with specified level.

        Args:
            text (str): Header text
            level (int): Header level (1-6)

        Returns:
            tuple: (response dict, status code)
        """
        if not text or level < 1 or level > 6:
            return jsonify({"error": "Invalid header parameters"}), 400

        try:
            header = f"{'#' * level} {text}"
            html = markdown(header)
            return jsonify({"html": f'<div class="markdown-header">{html}</div>'}), 200
        except Exception as e:
            return jsonify({"error": str(e)}), 500

    @staticmethod
    def render_list(items):
        """Renders a markdown list.

        Args:
            items (list): List of strings to render as markdown list

        Returns:
            tuple: (response dict, status code)
        """
        if not items or not isinstance(items, list):
            return jsonify({"error": "Invalid list parameters"}), 400

        try:
            list_content = "\n".join([f"- {item}" for item in items])
            html = markdown(list_content)
            return jsonify({"html": f'<div class="markdown-list">{html}</div>'}), 200
        except Exception as e:
            return jsonify({"error": str(e)}), 500

    @staticmethod
    def render_code_block(code, language="text"):
        """Renders a markdown code block with syntax highlighting.

        Args:
            code (str): Code content
            language (str): Programming language for syntax highlighting

        Returns:
            tuple: (response dict, status code)
        """
        if not code:
            return jsonify({"error": "Code content cannot be empty"}), 400

        try:
            code_block = f"