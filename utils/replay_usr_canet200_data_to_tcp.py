#!/usr/bin/env python3

import base64
import json
import socket
import sys
import time

TCP_HOST = "0.0.0.0"
TCP_PORT = 20002


def main():
    if (len(sys.argv) < 2):
        print('Must provide in file!!')
        sys.exit(1)

    in_file = sys.argv[1]

    packages = []

    with open(in_file, "r", encoding='utf-8') as f:
        packages = json.load(f)

    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)

    s.bind((TCP_HOST, TCP_PORT))
    s.listen()

    try:
        while True:
            conn, _ = s.accept()

            try:
                while True:
                    for package in packages:
                        data_b64 = package["data"]
                        data_b64_bytes = data_b64.encode("utf-8")
                        data_bytes = base64.b64decode(data_b64_bytes)
                        conn.send(data_bytes)
                        time.sleep(0.1)
            except Exception:
                pass
    except KeyboardInterrupt:
        pass

    print("Closing socket")
    s.shutdown(1)
    s.close()


if __name__ == '__main__':
    main()
