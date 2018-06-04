#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/socket.h>
#include <arpa/inet.h>

#define DEFAULT_PORT 12345

int main (int argc, char *argv[]) {
    int sock;
    struct sockaddr_in echo_server_addr;
    unsigned int echo_string_len;
    char echo_buffer[32];
    int total_recv_bytes;
    int err;
    int recv_bytes;
    int send_len;
    char *echo_string;

    echo_string = argv[1];

    sock = socket(PF_INET, SOCK_STREAM, IPPROTO_TCP);
    if (sock < 0) {
        perror("socket() failed");
        exit(1);
    }

    memset(&echo_server_addr, 0, sizeof(echo_server_addr));
    echo_server_addr.sin_family = AF_INET;
    echo_server_addr.sin_addr.s_addr = inet_addr("127.0.0.1");
    echo_server_addr.sin_port = htons(DEFAULT_PORT);

    err = connect(sock, (struct sockaddr *) &echo_server_addr, sizeof(echo_server_addr));
    if (err < 0) {
        perror("connect() failed");
        exit(1);
    }

    send_len = echo_string_len = strlen(echo_string);
    printf("send_len = %d\n", send_len);
    if (send_len != echo_string_len) {
        perror("send() failed");
        exit(1);
    }

    send(sock, echo_string, echo_string_len, 0);

    total_recv_bytes = 0;
    while (total_recv_bytes < echo_string_len) {
        recv_bytes = recv(sock, echo_buffer, 32-1, 0);
        if (recv_bytes <= 0) {
            perror("recv() failed");
            exit(1);
        }
        printf("recv_bytes = %d\n", recv_bytes);
        total_recv_bytes += recv_bytes;
        echo_buffer[recv_bytes] = '\0';
        printf("%s", echo_buffer);
    }

    close(sock);
    exit(0);
}