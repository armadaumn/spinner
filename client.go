package spinner

import (
  "github.com/armadanet/spinner/spinresp"
  "github.com/google/uuid"
  "github.com/armadanet/comms"
  "log"
)

// Single client (Captain) connection
type Client interface {
  Run()
  // Enter client into messenger system
  Register()
  // Quit the client
  Quit()
}

// client is abstration for a captain connection
type client struct {
  // pointing to spinner sync controller
  handler       *Handler
  // record the socket connection to captain
  socket        *comms.Socket
  // pointer for its instance in client messenger (in handler)
  self          *comms.Instance
  // Queue for submitted task
  spinup        chan interface{}
  // [taskId, requester id]
  responses     map[uuid.UUID]*uuid.UUID
  quit          chan struct{}
}

// Create new Client interface of client struct
func NewClient(h *Handler, socket *comms.Socket) Client {
  return &client{
    handler: h,
    socket: socket,
    spinup: make(chan interface{}),
    responses: make(map[uuid.UUID]*uuid.UUID),
    quit: make(chan struct{}),
    self: nil,
  }
}

// Get messages from the client
func (c *client) Run() {
  defer func(){
    c.handler.Unregister <- c.self
    (*c.socket).Close()
  }()
  read := (*c.socket).Reader()
  write := (*c.socket).Writer()
  for {
    select {
    // task execution response from captain (client)
    case response, ok := <- read:
      if !ok {return}
      resp, ok := response.(*spinresp.Response)
      if !ok {return}
      if resp.Id == nil {break}
      if identifier, ok := c.responses[*resp.Id]; ok {
        if resp.Code <= 0 {delete(c.responses, *resp.Id)}
        // new routine to send back to requester
        go func() {
          // identifier here is the requester id
          // send the response back to requester
          // return res back result channel -> only return when res is sent to writer of the requester
          if !c.handler.Requester.SendMessage(identifier, resp) {
            log.Printf("Failed: %+v\n", resp)
          }
        }()
      }
    // one task ready to send to captain (client)
    case data, ok := <- c.spinup:
      if !ok {break}
      task, ok := data.(*Task)
      if !ok {break}
      if task.Config.Id == nil {
        identifier := uuid.New()
        task.Config.Id = &identifier
      }
      c.responses[*task.Config.Id] = task.From
      write <- task.Config
    }
  }
}

// Register client with messenger and accept read/writes.
func (c *client) Register() {
  var resp spinresp.Response
  // socket reader and writer start buffering data
  // resp is the type of data reader use
  (*c.socket).Start(resp)
  // just create a client instance [captain id, captain spinup queue]
  c.self = c.handler.Requester.MakeInstance(c.spinup)
  // (BUG!!!) c.self = c.handler.clients.MakeInstance(c.spinup)
  c.handler.Register <- c.self
  go c.Run()
  log.Println("Client registered")

}

// Close the client connection
func (c *client) Quit() {close(c.quit)}
