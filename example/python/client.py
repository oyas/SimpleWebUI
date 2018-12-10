#!/usr/bin/env python3
# coding: utf-8

import subprocess
import socket
import time
import signal
import json

SocketPath = "io.socket"
SimpleWebUI_cmd = "../../SimpleWebUI"


def recv():
    client = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
    client.connect(SocketPath)
    response = client.recv(4096)
    client.close()
    return json.loads(response)


def main():
    cmd = SimpleWebUI_cmd + ' -i mark.md -t ../../static/index.html'
    proc = subprocess.Popen(cmd.split(), stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    print('Please open this URL in your browser')
    print('')
    print('    http://localhost:8080/')
    print('')

    time.sleep(1)  # waiting to boot

    data = recv()

    proc.send_signal(signal.SIGINT)
    #proc.send_signal(signal.SIGTERM)
    #proc.kill()
    proc.wait()

    data = {d['name']: d['value'] for d in data}
    print('Your name is ' + data['name'] + '.')


if __name__ == '__main__':
    main()
