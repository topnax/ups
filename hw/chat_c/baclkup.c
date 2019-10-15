/*
Dejan Tolj
License: GPLv2
Date: Jan 2005
*/

#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netdb.h>

#define MSG_SIZE 80
#define MAX_CLIENTS 95
#define MYPORT 7400
#define CLIENT_NAME_LEN 18

void exitClient(int fd, fd_set *readfds, char fd_array[], int *num_clients, char client_names[MAX_CLIENTS][CLIENT_NAME_LEN + 1])
{
    int i;
    char left_label[] = "has left the chat...\n";
    char left_message[CLIENT_NAME_LEN + 3 + sizeof(left_label)];

    /*concatinate the client id with the client's message*/
    sprintf(left_message, "%s %s", client_names[fd], left_label);
    printf("mess: %s", left_message);
    /*print to other clients*/
    for (i = 0; i < (*num_clients); i++)
    {
        if (fd_array[i] != fd) /*dont write msg to same client*/
        {
            printf("written to %d...\n", i);
            write(fd_array[i], left_message, strlen(left_message));
        }
    }

    printf("Exitting client %d...\n", fd);
    close(fd);
    FD_CLR(fd, readfds); //clear the leaving client from the set
    for (i = 0; i < (*num_clients) - 1; i++)
        if (fd_array[i] == fd)
            break;
    for (; i < (*num_clients) - 1; i++)
    {
        (fd_array[i]) = (fd_array[i + 1]);

        strcpy(client_names[i], client_names[i + 1]);
    }

    (*num_clients)--;
}

void server(int port)
{

    int read_length, i;
    int sockfd;
    struct hostent *hostinfo;

    char hostname[MSG_SIZE] = "127.0.0.1";
    struct sockaddr_in address;

    int response_length;

    int num_clients = 0;
    int server_sockfd, client_sockfd;
    struct sockaddr_in server_address;
    int addresslen = sizeof(struct sockaddr_in);
    int fd;
    char fd_array[MAX_CLIENTS];
    char client_names[MAX_CLIENTS][CLIENT_NAME_LEN + 1];
    fd_set readfds, testfds;
    char msg[MSG_SIZE + 1];
    char kb_msg[MSG_SIZE + 10];

    printf("starting server...\n");
    // printf("%s", hostname);
    printf("Using port: %d\n", port);

    fflush(stdout);

    /* Create and name a socket for the server */
    server_sockfd = socket(AF_INET, SOCK_STREAM, 0);
    server_address.sin_family = AF_INET;
    server_address.sin_addr.s_addr = htonl(INADDR_ANY);
    server_address.sin_port = htons(port);
    bind(server_sockfd, (struct sockaddr *)&server_address, addresslen);

    /* Create a connection queue and initialize a file descriptor set */
    listen(server_sockfd, 1);
    FD_ZERO(&readfds);
    FD_SET(server_sockfd, &readfds);
    FD_SET(0, &readfds); /* Add keyboard to file descriptor set */

    /*  Now wait for clients and requests */
    while (1)
    {
        testfds = readfds;
        select(FD_SETSIZE, &testfds, NULL, NULL, NULL);

        /* If there is activity, find which descriptor it's on using FD_ISSET */
        for (fd = 0; fd < FD_SETSIZE; fd++)
        {
            if (FD_ISSET(fd, &testfds))
            {

                if (fd == server_sockfd)
                { /* Accept a new connection request */
                    client_sockfd = accept(server_sockfd, NULL, NULL);
                    /*printf("client_sockfd: %d\n",client_sockfd);*/

                    if (num_clients < MAX_CLIENTS)
                    {
                        FD_SET(client_sockfd, &readfds);
                        fd_array[num_clients] = client_sockfd;
                        /*Client ID*/
                        printf("Client %d joined\n", num_clients++);
                        fflush(stdout);

                        //sprintf(msg, "M%2d", client_sockfd);
                        /*write 2 byte clientID */
                        //send(client_sockfd, msg, strlen(msg), 0);
                    }
                    else
                    {
                        sprintf(msg, "XSorry, too many clients.  Try again later.\n");
                        write(client_sockfd, msg, strlen(msg));
                        close(client_sockfd);
                    }
                }
                else if (fd == 0)
                { /* Process keyboard activity */
                    fgets(kb_msg, MSG_SIZE + 1, stdin);
                    //printf("%s\n",kb_msg);
                    if (strcmp(kb_msg, "quit\n") == 0)
                    {
                        sprintf(msg, "XServer is shutting down.\n");
                        for (i = 0; i < num_clients; i++)
                        {
                            write(fd_array[i], msg, strlen(msg));
                            close(fd_array[i]);
                        }
                        close(server_sockfd);
                        exit(0);
                    }
                    else
                    {
                        //printf("server - send\n");
                        sprintf(msg, "[BCAST] - %s", kb_msg);
                        for (i = 0; i < num_clients; i++)
                            write(fd_array[i], msg, strlen(msg));
                    }
                }
                else if (fd)
                { 
                    // message from client
                    read_length = read(fd, msg, MSG_SIZE);

                    if (read_length == -1)
                    {
                        printf("A client %s [%d] is leaving. RL == -1...\n", client_names[fd], fd);
                        perror("read()");
                    }
                    else if (read_length > 0)
                    {
                        sprintf(kb_msg, "M%1d", fd);
                        msg[read_length] = '\0';

                        if (msg[0] == 'n')
                        {
                            msg[read_length - 2] = '\0';
                            strcpy(client_names[fd], msg + 1);
                            sprintf(kb_msg, "%s has joined the chat...\n", client_names[fd]);
                            printf("%s", kb_msg);

                            /*print to other clients*/
                            for (i = 0; i < num_clients; i++)
                            {
                                if (fd_array[i] != fd) /*dont write msg to same client*/
                                    write(fd_array[i], kb_msg, strlen(kb_msg));
                            }
                        }
                        else if (msg[0] == 'm')
                        {
                            sprintf(kb_msg, "[%s]: %s", client_names[fd], msg + 1);
                            printf("%s", kb_msg);

                            for (i = 0; i < num_clients; i++)
                            {
                                if (fd_array[i] != fd)
                                {
                                    write(fd_array[i], kb_msg, strlen(kb_msg));
                                }
                            }
                        }

                        /*Exit Client*/
                        if (msg[0] == 'x')
                        {
                            printf("A client %s [%d] is leaving. X passed...\n", client_names[fd], fd);
                            exitClient(fd, &readfds, fd_array, &num_clients, client_names);
                        }
                    }
                }
                else
                {
                    printf("A client %s [%d] is leaving. RL == -1...\n", client_names[fd], fd);
                    exitClient(fd, &readfds, fd_array, &num_clients, client_names);
                }
            }
        }
    }

    exit(1);
}

void client(char hostname[MSG_SIZE], int port, char name[])
{
    /*Client variables=======================*/
    int sockfd;
    int result;
    struct hostent *hostinfo;
    struct sockaddr_in address;
    char alias[MSG_SIZE];
    int clientid;
    int fd;
    fd_set readfds, testfds, clientfds;
    char msg[MSG_SIZE + 1];
    char kb_msg[MSG_SIZE + 10];

    printf(".-> # Chat Client v0.0.1 # <-.\n");
    printf("Enter 'quit' to stop the client...\n");
    printf("Connecting to: %s\n", hostname);
    printf("Using port: %d\n", port);
    printf("Logging as: %s\n\n\n", name);

    fflush(stdout);

    sockfd = socket(AF_INET, SOCK_STREAM, 0);

    hostinfo = gethostbyname(hostname);
    address.sin_addr = *(struct in_addr *)*hostinfo->h_addr_list;
    address.sin_family = AF_INET;
    address.sin_port = htons(port);

    // Connect the socket to the server's socket
    if (connect(sockfd, (struct sockaddr *)&address, sizeof(address)) < 0)
    {
        perror("error connecting to the server's socket");
        exit(1);
    }

    fflush(stdout);

    FD_ZERO(&clientfds);
    FD_SET(sockfd, &clientfds);
    FD_SET(0, &clientfds); //stdin

    // "login" with the specified name
    sprintf(msg, "%s%s \n\0", "n", name);
    write(sockfd, msg, strlen(msg));

    while (1)
    {
        testfds = clientfds;
        select(FD_SETSIZE, &testfds, NULL, NULL, NULL);

        for (fd = 0; fd < FD_SETSIZE; fd++)
        {
            if (FD_ISSET(fd, &testfds))
            {
                if (fd == sockfd)
                {
                    // READ FROM SERVER FD
                    result = read(sockfd, msg, MSG_SIZE);
                    msg[result] = '\0';
                    printf("%s", msg);

                    if (msg[0] == 'X')
                    {
                        close(sockfd);
                        exit(0);
                    }
                }
                else if (fd == 0)
                {
                    // STD IN FD
                    fgets(kb_msg, MSG_SIZE + 1, stdin);
                    if (strcmp(kb_msg, "quit\n") == 0)
                    {
                        sprintf(msg, "xClient is shutting down.\n");
                        write(sockfd, msg, strlen(msg));
                        close(sockfd); //close the current client
                        exit(0);       //end program
                    }
                    else
                    {
                        sprintf(msg, "m%s", kb_msg);
                        write(sockfd, msg, strlen(msg));
                    }
                }
            }
        }
    }

    exit(-1);
}

int main(int argc, char *argv[])
{

    if (argc > 1)
    {
        if (strcmp("s", argv[1]) == 0)
        {
            if (argc == 3)
            {
                server(atoi(argv[2]));
            }
            else
            {
                printf("Server requires a port as a second parameter...\n");
                exit(1);
            }
        }
        else if (strcmp("c", argv[1]) == 0)
        {
            if (argc == 5)
            {
                client(argv[2], atoi(argv[3]), argv[4]);
            }
            else
            {
                printf("Client requires an IP address, a port and a name...\n");
                exit(1);
            }
        }
        else
        {
            printf("Unknown arguemnts...\n");
            exit(1);
        }
    }
    else
    {
        printf("Specify whether you want to run a server or a client...\n");
        exit(1);
    }
    return 0;
} //main
