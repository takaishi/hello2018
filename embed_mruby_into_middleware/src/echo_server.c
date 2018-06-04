#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include "mruby.h"
#include "mruby/compile.h"
#include "mruby/string.h"

#define DEFAULT_PORT 12345

typedef struct em_t {
    int reply_count;
    char *reply_prefix;
    mrb_state *mrb;
} em;

char *em_mrb_value_to_str(em *core, mrb_value value) {
    char *str;
    enum mrb_vtype type = mrb_type(value);

    if (mrb_undef_p(value) || mrb_nil_p(value)) {
        printf("undef or nil");
        asprintf(&str, "(nil)");
        return str;
    }

    switch (type) {
        case MRB_TT_FIXNUM: {
            asprintf(&str, "%s (integer) %lld\n", core->reply_prefix, mrb_fixnum(value));
            break;

        }
        case MRB_TT_STRING: {
            asprintf(&str, "%s (string) %s\n", core->reply_prefix, mrb_str_to_cstr(core->mrb, value));
            break;
        }
    }

    return str;
}


char *em_mrb_eval(em *core, char *str) {
    char *res;

    printf("str = %s\n", str);

    mrb_value value;
    mrbc_context *cxt;

    cxt = mrbc_context_new(core->mrb);
    value = mrb_load_string_cxt(core->mrb, str, cxt);
    res = em_mrb_value_to_str(core, value);

    return res;
}

int main (int argc, char *argv[]) {
    int server_sock;
    int client_sock;
    struct sockaddr_in echo_server_addr;
    struct sockaddr_in echo_client_addr;
    unsigned int client_len;
    int err;

    em *core;
    core = malloc(sizeof(em));
    core->reply_count = 0;
    core->reply_prefix = "em>";
    core->mrb = mrb_open();

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
        echo_buffer[recv_size] = '\0';
        printf("recv_size = %d\n", recv_size);
        printf("echo_buffer = %s\n", echo_buffer);

        while (recv_size > 0) {
            char *res;
            res = em_mrb_eval(core, echo_buffer);
            printf("res = %s\n", res);
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
