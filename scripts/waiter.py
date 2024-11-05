import json
from http.server import BaseHTTPRequestHandler, HTTPServer

class SimpleHTTPRequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        # Выводим информацию о запросе
        print(f"Получен GET-запрос: {self.path}")
        print(f"Запрос от: {self.client_address}")
        print(f"Заголовки запроса: {json.dumps(dict(self.headers), indent=4)}")
        
        # Отправляем ответ
        self.send_response(200)
        self.send_header('Content-type', 'text/plain')
        self.end_headers()
        self.wfile.write(b"success!\n")

    def do_POST(self):
        # Выводим информацию о POST-запросе
        print(f"Получен POST-запрос: {self.path}")
        print(f"Запрос от: {self.client_address}")
        print(f"Заголовки запроса: {json.dumps(dict(self.headers), indent=4)}")
        
        # Читаем тело запроса
        content_length = int(self.headers['Content-Length'])
        post_data = self.rfile.read(content_length)
        print(f"Тело запроса: {post_data.decode()}")

        # Отправляем ответ
        self.send_response(200)
        self.send_header('Content-type', 'text/plain')
        self.end_headers()
        self.wfile.write(b"POST-success!\n")

def run(server_class=HTTPServer, handler_class=SimpleHTTPRequestHandler, port=4318):
    server_address = ('', port)
    httpd = server_class(server_address, handler_class)
    print(f"Запуск сервера на порту {port}...")
    httpd.serve_forever()

if __name__ == '__main__':
    run(port=4318)
