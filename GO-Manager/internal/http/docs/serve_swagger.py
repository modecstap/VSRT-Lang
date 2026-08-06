import http.server
import socketserver
import os

PORT = 8081
BASE_DIR = os.path.dirname(__file__)

class Handler(http.server.SimpleHTTPRequestHandler):
    def translate_path(self, path):
        # Map /swagger/doc.yaml -> internal/http/docs/openapi.yaml
        if path == '/swagger/doc.yaml' or path == '/swagger/doc.yaml/':
            return os.path.join(BASE_DIR, 'openapi.yaml')
        # Serve swagger-ui files under /swagger/
        if path.startswith('/swagger/'):
            rel = path[len('/swagger/'):]
            if rel == '' or rel == '/':
                rel = 'index.html'
            return os.path.join(BASE_DIR, 'swagger-ui', rel)
        # Fallback to serving files from BASE_DIR
        return os.path.join(BASE_DIR, path.lstrip('/'))

    def log_message(self, format, *args):
        # quieter logs
        pass

with socketserver.TCPServer(("", PORT), Handler) as httpd:
    print(f"Serving at http://localhost:{PORT}")
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        httpd.shutdown()
