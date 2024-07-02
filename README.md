# IP Address and Greeting Service (Go)

This Go program provides a basic HTTP server that does the following:

1. **Retrieves the client's IP address.**
2. **Returns a JSON response** containing:
   * The client's IP address
   * Get the current client's IP address
   * A personalized greeting with the client name and temperature

### How to Use

1. **Prerequisites:**
   * Go installed on your system ([https://golang.org/](https://golang.org/))

2. **Run the Server:**
   * Open your terminal.
   * Navigate to the directory where this code is saved.
   * Execute:  `go run main.go`

3. **Access the Service:**
   * Open your web browser or use a tool like `curl`.
   * Visit: `http://localhost:8000`
   * You'll see a JSON response similar to this:

   ```json
   {"ip": ${client's ip},"location":${client's location},"greeting":"Hello, ${client's name}, The temperature is 1${client's temperature} degrees Celsius in ${client's location}"}

### How to Build and Run with Docker

* **Build:** Open a terminal in that directory and run:

```sh
docker build -t ip-greeter-image .

# This creates a Docker image named ip-greeter-image.
```

- **Run**

```sh
docker run -p 8000:8000 ip-greeter-image

# This starts a container from the image, publishing port 8000 inside the container to port 8000 on your local machine.
```
