# ouR/Place

This is an attempt to make a simplified version of r/Place from reddit using websockets and Go.

## Structure
### `internal`
Contains the `Client` struct that imlpements the `Subscriber` interface for the `SMessage` Type. Handles connections with the users.

Contains the `Canvas` struct. Represents the canvas of r/Place. Also implements the Subscriber State. On a new connection the entire canvas is sent to the user. Currently kept only in memory.

Contains the `Hub` struct which handles registered clients and notifies the attached clients about changes in the canvas state.

### `frontend`
Contains a react app that visualizes the thing. It is just a canvas that with mouse event handlers attached to it.

The connection to the websocket server exists in the `WSContext`. it is initiated in the `Canvas` component.

### TODO:
- [x] Setup authentication (preferably OAuth)
- [x] Define and implement the Publisher interface for the Hub
- [ ] Add handlers for different types of incoming messages for the Client struct.
- [ ] Move all message related code to a separate package.
- [ ] Add consistency across reloads using database/file
- [ ] Consider using protobufs instead of the hideous scheme in use right now.
- [ ] Try using WebGL instead of Canvas.
- [ ] Performance testing
- [ ] Deploy :,)
