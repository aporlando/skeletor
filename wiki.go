package main
import (
       "net/http"
       "log"

       "github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)	     // broadcast channel

// configure the upgrader
var upgrader = websocket.Upgrader{}

// define our message object
type Message struct {
     Email    string `json:"email"`
     Username string `json:"username"`
     Message  string `json:"message"`
}

func handleConnections(writer http.ResponseWriter, request *http.Request) {
     ws, err := upgrader.Upgrade(writer, request, nil)
     if err != nil {
     	log.Fatal(err)
     }

     // Make sure we close the connection when the function returns
     defer ws.Close()

     // Register our new client
     clients[ws] = true

     for {
     	 var msg Message
	 // Read in a new message as JSON and map it to a Message object
	 err := ws.ReadJSON(&msg)
	 if err != nil {
	    log.Printf("error: %v", err)
	    delete(clients, ws)
	    break
	 }

	 // Send the newly received message to the broadcast channel
	 broadcast <- msg
     }
}

func handleMessages() {
     for {
     	 // Grab the next message from the broadcast channel
	 msg := <-broadcast
	 // Send it out to every client that is currently connected
	 for client := range clients {
	     err := client.WriteJSON(msg)
	     if err != nil {
	     	log.Printf("error: %v", err)
		client.Close()
		delete(clients, client)
	     }
	 }
     }
}

func main() {
     // Create a simple file server
     fs := http.FileServer(http.Dir("../public"))
     http.Handle("/", fs)
     http.HandleFunc("/ws", handleConnections)

     // Start listening for incoming chat messages
     go handleMessages()

     // Start the server on localhost port 5000 and log any errors
     log.Println("http server started on :5000")
     err := http.ListenAndServe(":5000", nil)
     if err != nil {
     	log.Fatal("ListenAndServe: ", err)
     }
}
