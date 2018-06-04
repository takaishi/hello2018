#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/socket.h>
#include <arpa/inet.h>

#define DEFAULT_PORT 12345

int main (int argc, char *argv[]) {
    int server_sock;
    int client_sock;
    struct sockaddr_in echo_server_addr;
    struct sockaddr_in echo_client_addr;
    unsigned int client_len;
    int err;

    server_sock = socket(PF_INET, SOCK_STREAM, IPPROTO_TCP);
    if (server_sock < 0) {
        perror("socket() failed");
        exit(1);
    }

    memset(&echo_server_addr, 0, sizeof(echo_server_addr));
    echo_server_addr.sin_family = AF_INET;
    echo_server_addr.sin_addr.s_addr = htonl(INADDR_ANY);
    echo_server_addr.sin_port = htons(DEFAULT_PORT);

    err = bind(server_sock, (struct sockaddr *) &echo_server_addr, sizeof(echo_server_addr));
    if (err < 0) {
        perror("bind() failed");
        exit(1);
    }

    err = listen(server_sock, 5);
    if (err < 0) {
        perror("listen() failed");
        exit(1);
    }

    for (;;) {
        char echo_buffer[32];
        int recv_size;

        client_len = sizeof(echo_client_addr);
        client_sock = accept(server_sock, (struct sockaddr *) &echo_client_addr, &client_len);
        if (client_sock < 0) {
            perror("accept() failed");
            exit(1);
        }

        recv_size = recv(client_sock, echo_buffer, 32, 0);
        if (recv_size < 0) {
            perror("recv() failed");
            exit(1);
        }
        printf("recv_size = %d\n", recv_size);
        printf("echo_buffer = %s\n", echo_buffer);

        while (recv_size > 0) {
            err = send(client_sock, echo_buffer, recv_size, 0);
            if (err < 0) {
                perror("send() failed");
                exit(1);
            }

            recv_size = recv(client_sock, echo_buffer, 32, 0);
            if (recv_size < 0) {
                perror("recv() failed");
                exit(1);
            }
        }
        close(client_sock);
    }
}