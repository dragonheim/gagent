### [Source Code Scan](#source)
```
2023-03-15T06:59:15.989-0700    INFO    Need to update DB
2023-03-15T06:59:15.989-0700    INFO    DB Repository: ghcr.io/aquasecurity/trivy-db
2023-03-15T06:59:15.989-0700    INFO    Downloading DB...
36.01 MiB / 36.01 MiB [-----------------------------------------------------------------------------------------------------------------------------------------------------------------] 100.00% 7.93 MiB p/s 4.7s
2023-03-15T06:59:21.791-0700    INFO    Vulnerability scanning is enabled
2023-03-15T06:59:21.791-0700    INFO    Secret scanning is enabled
2023-03-15T06:59:21.791-0700    INFO    If your scanning is slow, please try '--scanners vuln' to disable secret scanning
2023-03-15T06:59:21.791-0700    INFO    Please see also https://aquasecurity.github.io/trivy/v0.37/docs/secret/scanning/#recommendation for faster secret detection
2023-03-15T06:59:25.698-0700    INFO    Number of language-specific files: 1
2023-03-15T06:59:25.699-0700    INFO    Detecting gomod vulnerabilities...

Total: 0 (UNKNOWN: 0, LOW: 0, MEDIUM: 0, HIGH: 0, CRITICAL: 0)
```

---
### [Image Scan](#image)
```
2023-03-15T07:05:44.377-0700    INFO    Vulnerability scanning is enabled
2023-03-15T07:05:44.377-0700    INFO    Secret scanning is enabled
2023-03-15T07:05:44.377-0700    INFO    If your scanning is slow, please try '--scanners vuln' to disable secret scanning
2023-03-15T07:05:44.377-0700    INFO    Please see also https://aquasecurity.github.io/trivy/v0.37/docs/secret/scanning/#recommendation for faster secret detection
2023-03-15T07:05:44.731-0700    INFO    Detected OS: alpine
2023-03-15T07:05:44.731-0700    INFO    Detecting Alpine vulnerabilities...
2023-03-15T07:05:44.732-0700    INFO    Number of language-specific files: 1
2023-03-15T07:05:44.732-0700    INFO    Detecting gobinary vulnerabilities...

dragonheim/gagent:0.0.7 (alpine 3.17.2)

Total: 0 (UNKNOWN: 0, LOW: 0, MEDIUM: 0, HIGH: 0, CRITICAL: 0)
```