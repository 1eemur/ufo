# UFO - Universal File Operator

A modern, responsive web application for file uploads built with Go and vanilla JavaScript. Intended for quick file transfer in a trusted local network or lab environment. 

**N.B. This is not intended as a secure method of online or public file transfer and there are currently relatively limited input sanitization or safeguards against malicious uploads.**

## Features

- 🚀 **Modern UI**: Beautiful, responsive design with smooth animations
- 📁 **Drag & Drop**: Intuitive file upload with drag and drop support
- 📊 **Progress Tracking**: Real-time upload progress with visual feedback
- 📋 **File Management**: View uploaded files with metadata (size, date)
- ⚡ **Fast**: Built with Go for optimal performance
- 📱 **Responsive**: Works seamlessly on desktop and mobile devices

## Quick Start

### Prerequisites

- Go 1.21 or later
- Modern web browser

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/1eemur/ufo.git
   cd ufo
   ```

2. **Download dependencies**
   ```bash
   make deps
   ```

3. **Build and run**
   ```bash
   make run
   ```

4. **Open your browser**
   Navigate to `http://localhost:8080`

## Project Structure

```
ufo/
├── cmd/server/main.go              # Application entry point
├── internal/
│   ├── config/config.go            # Configuration management
│   ├── handlers/upload.go          # HTTP handlers
│   └── middleware/logging.go       # Middleware
├── web/
│   ├── static/
│   │   ├── css/style.css          # Stylesheets
│   │   └── js/upload.js           # Frontend JavaScript
│   └── templates/index.html        # HTML template
├── data/                           # Upload directory
├── Makefile                        # Build automation
└── README.md                       # Documentation
```

## Configuration

Configure the application using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DATA_DIR` | Upload directory | `./data` |
| `MAX_FILE_SIZE` | Maximum file size in bytes | `104857600` (100MB) |

### Examples

```bash
# Run on port 3000
PORT=3000 make run

# Custom upload directory
DATA_DIR=/uploads make run

# 50MB file size limit
MAX_FILE_SIZE=52428800 make run
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Main upload page |
| `POST` | `/upload` | Upload files |
| `GET` | `/files` | List uploaded files |
| `GET` | `/static/*` | Static assets |

### Available Commands

```bash
make help           # Show all available commands
make build          # Build the application
make run            # Build and run
make dev            # Development mode with auto-reload
make clean          # Clean build artifacts
make test           # Run tests
make deps           # Update dependencies
make setup          # Create directory structure
```

## To Do 

- **Password Protection**: Add functionality and configuration for login requirements.
- **File Downloads**: Implement functionality to download designated files from a separate folder on the server as well.
- **Error Handling**: Better user feedback and error handling for file size limits.