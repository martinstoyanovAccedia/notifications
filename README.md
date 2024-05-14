# notifications
Take home notification-sending system.

### Usage

- **Authentication**: First, authenticate the user by sending a POST request to `localhost:8000/login` with JSON containing the username and password.

  Example Request:
    ```json
    {
        "username": "test",
        "password": "test"
    }
    ```

- **Sending Email Notifications**: Use the `/email` endpoint to send email notifications. The system will dial the email client and wait for messages for 5 minutes.

- **Adding Messages**: Send JSON-formatted messages to the `/messages` endpoint to add messages to both email and potentially other channels.

  Example Request:
    ```json
    [
        {"title": "Notification 1", "message": "This is message 1"},
        {"title": "Notification 2", "message": "This is message 2"}
    ]
    ```
