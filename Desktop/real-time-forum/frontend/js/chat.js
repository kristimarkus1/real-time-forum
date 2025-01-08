const chatContainer = document.getElementById("chat");

const socket = new WebSocket("ws://localhost:8080/ws");

// Log connection status
socket.onopen = () => {
    console.log("Connected to WebSocket server");
};

// Handle incoming messages
socket.onmessage = (event) => {
    const message = JSON.parse(event.data);
    displayMessage(message);
};

// Log errors
socket.onerror = (error) => {
    console.error("WebSocket error:", error);
};

// Log when the connection is closed
socket.onclose = () => {
    console.log("Disconnected from WebSocket server");
};

// Function to display messages
function displayMessage(message) {
    const messageElement = document.createElement("div");
    messageElement.textContent = `[${message.date}] ${message.sender}: ${message.text}`;
    chatContainer.appendChild(messageElement);
}

// Function to send a message
function sendMessage(text) {
    const message = {
        text,
        sender: "YourNickname", // Replace with logged-in user's nickname
        date: new Date().toLocaleString(),
    };
    socket.send(JSON.stringify(message));
}

// Example: Attach sendMessage to a button
const sendButton = document.createElement("button");
sendButton.textContent = "Send Message";
sendButton.onclick = () => {
    const messageInput = document.createElement("input");
    messageInput.placeholder = "Type a message...";
    chatContainer.appendChild(messageInput);

    sendMessage(messageInput.value);
    messageInput.value = ""; // Clear input after sending
};
chatContainer.appendChild(sendButton);
