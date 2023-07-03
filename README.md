# Trivy Scanner

Trivy Scanner is a Go script that reads a manifest file containing a list of Docker images, pulls the images locally (if not already present), and performs a Trivy scan on each image. The Trivy scan reports are then stored in individual files.

## Presentation

The Trivy Scanner script performs the following steps:

1. Creates a "scan-reports" folder to store the scan reports.
2. Reads the manifest.yaml file containing the list of Docker images.
3. Parses the YAML content and extracts the list of saved images.
4. Iterates over each image and performs the following actions:
    - Checks if the Docker image already exists locally.
    - If the image does not exist locally, pulls the image.
    - Performs a Trivy scan on the image.
    - Stores the Trivy scan report in a separate file.

## How to Use

1. Ensure you have Docker and Trivy installed on your system.
2. Make sure to provide a correct `manifest.yaml`, place it at the root of the project.
3. Open a terminal and navigate to the root directory of the project.
4. Run the following command to execute the script:

   ```bash
   make run-scanner
   ```
# trivy-scanner
