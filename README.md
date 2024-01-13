# SecretDetector

## Description

Recursively searches for secrets from given start location.
Project leverages the default secret patterns available on [Gitlab](https://gitlab.com/gitlab-org/security-products/analyzers/secrets/-/blob/master/gitleaks.toml).
A sequential scan is performed initially to establish a baseline and to validate the concurrent implementation.

## Usage

**run.sh** can be used to launch the project. The script will:
- Build the **secretdetector** image
- Launch the container
- Copy the report file, **report.json**, from container upon execution completion.
- Clean up.

_NOTE:_
_Project is configured to start the scan at ./thirdParty directory located at the project root.
You can modify this location by updating the **SCAN_START_DIRECTORY** environment variable within the
Dockerfile._

### Sample Output
```
2024/01/13 14:57:33.581141 [+]retrieving secret patterns from ./data/secretpatterns.toml
2024/01/13 14:57:33.622815 [/]secret patterns loaded: 121
2024/01/13 14:57:33.622829 [+] Starting Sequential Secret Detection.
2024/01/13 14:57:34.013144 [/]start time: 1705132653
2024/01/13 14:57:34.013187 [/]end time: 1705132654 [/]total files processed: 63;  in 1s time
# of Potential Secrets Found: 24
2024/01/13 14:57:34.013201 [+] Starting Concurrent Secret Detection.
2024/01/13 14:57:34.066074 [/]start time: 1705132654
2024/01/13 14:57:34.066134 [/]end time: 1705132654 [/]total files processed: 63;  in 0s time
# of Potential Secrets Found: 24
```