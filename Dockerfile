diff --git a/Dockerfile b/Dockerfile
index 123456..789012 100644
--- a/Dockerfile
+++ b/Dockerfile
@@ -1,5 +1,7 @@
 FROM python:3.9-slim

 # Set environment variables
 ENV AXENTX_DB_USER=axentx
 ENV AXENTX_DB_PASSWORD=secret
+ENV AXENTX_DB_HOST=localhost
+ENV AXENTX_DB_PORT=5432

 # Copy requirements file
 COPY requirements.txt .

@@ -9,3 +11,6 @@
 RUN pip install -r requirements.txt

 # Expose port
 EXPOSE 8000

+ # Health check endpoint
+ HEALTHCHECK --interval=1m --timeout=3s --retries=3 --start-period=30s \
+   curl -f http://localhost:8000/health || exit 1