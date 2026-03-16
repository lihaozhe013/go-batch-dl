import http.server
import socketserver
import threading
import time
import subprocess
import shutil
from pathlib import Path

# Global Configuration
PORT = 8080
HOST = "localhost"
BASE_URL = f"http://{HOST}:{PORT}"

# Calculate absolute paths
# Assuming tester.py is located at {PROJECT_ROOT}/tester/tester.py
CURRENT_DIR = Path(__file__).resolve().parent
PROJECT_ROOT = CURRENT_DIR.parent

TEST_SERVER_ROOT = PROJECT_ROOT / "test_server_root"
# We define DOWNLOAD_DIR name here, but use the absolute path for operations
DOWNLOAD_DIR_NAME = "test_downloads"
ABS_DOWNLOAD_DIR = PROJECT_ROOT / DOWNLOAD_DIR_NAME

# Number of test files to create
NUM_FILES = 10

def setup_test_files():
    """Create HTML and TXT files for testing."""
    # Clean up old server directory
    if TEST_SERVER_ROOT.exists():
        shutil.rmtree(TEST_SERVER_ROOT)
    TEST_SERVER_ROOT.mkdir(parents=True, exist_ok=True)

    links = []
    print(f"[Tester] Creating {NUM_FILES} test files in {TEST_SERVER_ROOT}...")
    for i in range(NUM_FILES):
        filename = f"file_{i}.txt"
        file_path = TEST_SERVER_ROOT / filename
        file_path.write_text(f"This is content for file {i}", encoding="utf-8")
        
        # Generate corresponding link in index.html
        links.append(f'<a href="{filename}">Download {filename}</a>')

    # Create index.html
    index_path = TEST_SERVER_ROOT / "index.html"
    index_content = (
        "<html><body>"
        "<h1>Test Files</h1>"
        + "<br>".join(links) +
        "</body></html>"
    )
    index_path.write_text(index_content, encoding="utf-8")
    print("[Tester] index.html created.")

def start_server():
    """Start background HTTP server."""
    print(f"[Tester] Starting HTTP server on port {PORT} serving {TEST_SERVER_ROOT}...")
    
    # Change to server root directory so SimpleHTTPRequestHandler serves files correctly
    # Note: This changes the CWD for the thread/process, but we use absolute paths elsewhere
    # Using os.chdir is acceptable here as SimpleHTTPRequestHandler defaults to CWD
    import os
    os.chdir(TEST_SERVER_ROOT)
    
    # Use socketserver to allow address reuse, avoiding errors on restart
    socketserver.TCPServer.allow_reuse_address = True
    with socketserver.TCPServer(("", PORT), http.server.SimpleHTTPRequestHandler) as httpd:
        httpd.serve_forever()

def run_go_program():
    """Run the Go downloader."""
    # Ensure download directory is clean
    if ABS_DOWNLOAD_DIR.exists():
        shutil.rmtree(ABS_DOWNLOAD_DIR)

    # Construct go run command
    # Using str() to convert Path objects to strings for subprocess
    cmd = [
        "go", "run", "cmd/gobatchdl/main.go",
        "-url", BASE_URL,
        "-ext", ".txt",
        "-dir", str(ABS_DOWNLOAD_DIR), # Use absolute path for safety
        "-workers", "4"
    ]
    
    print("-" * 50)
    print(f"[Tester] Running command: {' '.join(cmd)}")
    
    try:
        # Run go command in the project root directory
        # subprocess.run temporarily changes CWD to PROJECT_ROOT
        result = subprocess.run(
            cmd, 
            cwd=PROJECT_ROOT, 
            capture_output=True, 
            text=True
        )
        
        # Print output
        if result.stdout:
            print("[STDOUT]:\n" + result.stdout)
        if result.stderr:
            print("[STDERR]:\n" + result.stderr)
            
        print("-" * 50)
        
        # Verify download results
        if ABS_DOWNLOAD_DIR.exists():
            files = list(ABS_DOWNLOAD_DIR.iterdir())
            downloaded = len(files)
            print(f"[Tester] Result: Downloaded {downloaded}/{NUM_FILES} files.")
            if downloaded == NUM_FILES:
                print("[Tester] SUCCESS! All files downloaded.")
            else:
                print(f"[Tester] FAILED: File count mismatch. Found {downloaded}, expected {NUM_FILES}.")
        else:
            print("[Tester] FAILED: Download directory not found.")
            
    except Exception as e:
        print(f"[Tester] Error running Go program: {e}")

def main():
    try:
        # 1. Prepare files
        setup_test_files()

        # 2. Start server thread
        # Set daemon=True so server thread exits automatically when main program ends
        server_thread = threading.Thread(target=start_server, daemon=True)
        server_thread.start()
        
        # Give server a moment to start
        time.sleep(1)

        # 3. Run Go program
        run_go_program()
        
    except KeyboardInterrupt:
        print("\n[Tester] Interrupted.")
    except Exception as e:
        print(f"\n[Tester] Error: {e}")
    finally:
        pass
    
    print("[Tester] Done.")

if __name__ == "__main__":
    main()
