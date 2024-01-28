#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <time.h>

#define PORT 3032

int main() {
    int server_fd, new_socket;
    struct sockaddr_in address;
    int opt = 1;
    int addrlen = sizeof(address);
    char buffer[1024] = { 0 };

    FILE* logfile;
    logfile = fopen("logs/http.log", "a");
    if (logfile == NULL) {
        perror("open log file");
        exit(EXIT_FAILURE);
    }

    // Creating socket file descriptor
    if ((server_fd = socket(AF_INET, SOCK_STREAM, 0)) == 0) {
        perror("socket failed");
        exit(EXIT_FAILURE);
    }

    // Forcefully attaching socket to the port
    if (setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt))) {
        perror("setsockopt");
        exit(EXIT_FAILURE);
    }

    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY;
    address.sin_port = htons(PORT);

    // Forcefully attaching socket to the port
    if (bind(server_fd, (struct sockaddr*)&address, sizeof(address)) < 0) {
        perror("bind failed");
        exit(EXIT_FAILURE);
    }

    if (listen(server_fd, 3) < 0) {
        perror("listen");
        exit(EXIT_FAILURE);
    }

    while (1) {
        time_t now = time(NULL);
        char* timestamp = ctime(&now);
        timestamp[strlen(timestamp) - 1] = '\0'; // Remove newline at the end of timestamp

        printf("Waiting for connections...\n");

        if ((new_socket = accept(server_fd, (struct sockaddr*)&address, (socklen_t*)&addrlen)) < 0) {
            perror("accept");
            exit(EXIT_FAILURE);
        }

        // Read the HTTP request
        char request[3000];
        read(new_socket, request, 3000);

        // Log the HTTP request
        fprintf(logfile, "[%s] Received request: %s\n", timestamp, request);
        fflush(logfile);

        // Check if this is a GET request to /verify
        if (strncmp(request, "GET /verify ", 12) == 0) {

            char time[500];
            snprintf(time, sizeof(time), "%s", timestamp);
            char jsonValue[1024];
            snprintf(jsonValue, sizeof(jsonValue), "{\"isValid\":\"True\", \"time\":\"%s\"}", time);
            int jsonLength = strlen(jsonValue);

            char response[20484];
            sprintf(response, "HTTP/1.1 200 OK\nContent-Type: application/json\nContent-Length: %d\n\n%s", jsonLength, jsonValue);

            // Send HTTP response
            send(new_socket, response, strlen(response), 0);
        }
        else {
            // Send 404 Not Found for other requests
            char* notFound = "HTTP/1.1 404 Not Found\nContent-Type: text/plain\nContent-Length: 9\n\nNot Found";
            send(new_socket, notFound, strlen(notFound), 0);
        }

        // Close the socket for this connection
        close(new_socket);
    }
    fclose(logfile);
    // Close the main listening socket (this will not be reached in the current loop)
    close(server_fd);

    return 0;
}
