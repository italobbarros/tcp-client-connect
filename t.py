import socket

def make_http_request():
    # Cria um socket TCP
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    try:
        # Conecta ao servidor
        server_address = ('httpbin.org', 80)
        sock.connect(server_address)

        # Envia a requisição HTTP
        message = 'GET /get HTTP/1.1\r\nHost: httpbin.org\r\nAccept: application/json\r\n\r\n'
        sock.sendall(message.encode('utf-8'))

        # Recebe a resposta
        data = sock.recv(4096)

        print(data.decode())

    finally:
        sock.close()

make_http_request()
