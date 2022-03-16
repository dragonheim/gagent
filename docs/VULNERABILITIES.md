### [Source Code Scan](#source)

IGNORED: We are not using the protocol buffers provided by gogo/protobuf
```
2022-03-16T16:27:43.221-0700    INFO    Need to update DB
2022-03-16T16:27:43.221-0700    INFO    Downloading DB...
26.43 MiB / 26.43 MiB [---------------------------------------------------------------------------------------------------------------------------------------------------------------] 100.00% 5.71 MiB p/s 5s
2022-03-16T16:27:49.073-0700    INFO    Number of language-specific files: 1
2022-03-16T16:27:49.073-0700    INFO    Detecting gomod vulnerabilities...

go.sum (gomod)
==============
Total: 1 (UNKNOWN: 0, LOW: 0, MEDIUM: 0, HIGH: 1, CRITICAL: 0)

+--------------------------+------------------+----------+-------------------+---------------+--------------------------------------+
|         LIBRARY          | VULNERABILITY ID | SEVERITY | INSTALLED VERSION | FIXED VERSION |                TITLE                 |
+--------------------------+------------------+----------+-------------------+---------------+--------------------------------------+
| github.com/gogo/protobuf | CVE-2021-3121    | HIGH     | 1.1.1             | 1.3.2         | gogo/protobuf:                       |
|                          |                  |          |                   |               | plugin/unmarshal/unmarshal.go        |
|                          |                  |          |                   |               | lacks certain index validation       |
|                          |                  |          |                   |               | -->avd.aquasec.com/nvd/cve-2021-3121 |
+--------------------------+------------------+----------+-------------------+---------------+--------------------------------------+
```
---
### [Image Scan](#image)

NONE