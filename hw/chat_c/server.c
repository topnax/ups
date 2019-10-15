#include <stdio.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>
#include <netinet/in.h>
#include <stdlib.h>
// kvuli iotctl
#include <sys/ioctl.h>

#define MSG_SIZE 80


int main (void){

	int server_socket;
	int client_socket, fd;
	int return_value;
	char cbuf;
	int len_addr;
	int a2read;
	int i;
	struct sockaddr_in my_addr, peer_addr;
	fd_set client_socks, tests;

	server_socket = socket(AF_INET, SOCK_STREAM, 0);

	memset(&my_addr, 0, sizeof(struct sockaddr_in));

	my_addr.sin_family = AF_INET;
	my_addr.sin_port = htons(10000);
	//_addr.sin_addr.s_addr = htonl(INADDR_ANY);
	my_addr.sin_addr.s_addr = inet_addr("127.0.0.1");

	return_value = bind(server_socket, (struct sockaddr *) &my_addr, \
		sizeof(struct sockaddr_in));

	if (return_value == 0) 
		printf("Bind - OK\n");
	else {
		printf("Bind - ERR\n");
		return -1;
	}

	return_value = listen(server_socket, 5);
	if (return_value == 0){
		printf("Listen - OK\n");
	} else {
		printf("Listen - ER\n");
	}

	// vyprazdnime sadu deskriptoru a vlozime server socket
	FD_ZERO( &client_socks );
	FD_SET( server_socket, &client_socks );

	 while (1) {
        tests = client_socks;
        select(FD_SETSIZE, &tests, NULL, NULL, NULL);
                    
        /* If there is activity, find which descriptor it's on using FD_ISSET */
        for (fd = 0; fd < FD_SETSIZE; fd++) {
           printf("hello");
           if (FD_ISSET(fd, &tests)) {
              
              if (fd == server_socket) { /* Accept a new connection request */
                 client_socket = accept(server_socket, NULL, NULL);
                 /*printf("client_sockfd: %d\n",client_sockfd);*/
                
                                
                 //if (num_clients < MAX_CLIENTS) {
                    FD_SET(client_socket, &client_socks);
                   // fd_array[num_clients]=client_sockfd;
                    /*Client ID*/
                    //printf("Client %d joined\n",num_clients++);
                    //fflush(stdout);
                    
                    //sprintf(msg,"M%2d",client_sockfd);
                    /*write 2 byte clientID */
					printf("A client has been accepted");
                    //send(client_sockfd,msg,strlen(msg),0);
                 //}
                // else {
                  //  sprintf(msg, "XSorry, too many clients.  Try again later.\n");
                   // write(client_sockfd, msg, strlen(msg));
                   // close(client_sockfd);
                 }
              }
			  /**
              else if (fd == 0)  {  
                 fgets(kb_msg, MSG_SIZE + 1, stdin);
                 //printf("%s\n",kb_msg);
                 if (strcmp(kb_msg, "quit\n")==0) {
                    sprintf(msg, "XServer is shutting down.\n");
                    for (i = 0; i < num_clients ; i++) {
                       write(fd_array[i], msg, strlen(msg));
                       close(fd_array[i]);
                    }
                    close(server_sockfd);
                    exit(0);
                 }
                 else {
                    //printf("server - send\n");
                    sprintf(msg, "M%s", kb_msg);
                    for (i = 0; i < num_clients ; i++)
                       write(fd_array[i], msg, strlen(msg));
                
				 }
              }
			  */
              else if(fd) {  /*Process Client specific activity*/
                 //printf("server - read\n");
                 //read data from open socket
				char msg[MSG_SIZE + 1];     

                 int result = read(fd, msg, MSG_SIZE);
                 
                 if(result==-1) perror("read()");
                 else if(result>0){
                    /*read 2 bytes client id*/
                    //sprintf(kb_msg,"M%2d",fd);
                    msg[result]='\0';
                    
                    /*concatinate the client id with the client's message*/
                    // strcat(kb_msg,msg+1);                                        
                    
                    /*print to other clients*/
                   // for(i=0;i<num_clients;i++){
                     //  if (fd_array[i] != fd)  /*dont write msg to same client*/
                       //   write(fd_array[i],kb_msg,strlen(kb_msg));
                    //}
                    /*print to server*/
                    printf("%s",msg);
                    
                     /*Exit Client*/
                    if(msg[0] == 'X'){
                       // exitClient(fd,&readfds, fd_array,&num_clients);
                    }   
                 }                                   
              }                  
              else {  /* A client is leaving */
                 // exitClient(fd,&readfds, fd_array,&num_clients);
              }//if
           }//if
        //for
     }//while

	return 0;
	
}
