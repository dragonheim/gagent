### [Source Code Scan](#source)

IGNORED: We are not using the SSH features of golang.org/x/crypto
```
2021-11-13T10:25:13.188-0800    INFO    Need to update DB
2021-11-13T10:25:13.188-0800    INFO    Downloading DB...
24.70 MiB / 24.70 MiB [----------------------------------------------------------------------------------------------------------------------------------------------------------------------------] 100.00% 6.04 MiB p/s 4s
2021-11-13T10:25:18.570-0800    INFO    Detected OS: unknown
2021-11-13T10:25:18.570-0800    INFO    Number of PL dependency files: 2
2021-11-13T10:25:18.570-0800    INFO    Detecting gobinary vulnerabilities...
2021-11-13T10:25:18.571-0800    INFO    Detecting gomod vulnerabilities...

bin/gagent
==========
Total: 0 (UNKNOWN: 0, LOW: 0, MEDIUM: 0, HIGH: 0, CRITICAL: 0)


go.sum
======
Total: 1 (UNKNOWN: 0, LOW: 0, MEDIUM: 0, HIGH: 1, CRITICAL: 0)

+---------------------+------------------+----------+-----------------------------------+------------------------------------+---------------------------------------+
|       LIBRARY       | VULNERABILITY ID | SEVERITY |         INSTALLED VERSION         |           FIXED VERSION            |                 TITLE                 |
+---------------------+------------------+----------+-----------------------------------+------------------------------------+---------------------------------------+
| golang.org/x/crypto | CVE-2020-29652   | HIGH     | 0.0.0-20200622213623-75b288015ac9 | v0.0.0-20201216223049-8b5274cf687f | golang: crypto/ssh: crafted           |
|                     |                  |          |                                   |                                    | authentication request can            |
|                     |                  |          |                                   |                                    | lead to nil pointer dereference       |
|                     |                  |          |                                   |                                    | -->avd.aquasec.com/nvd/cve-2020-29652 |
+---------------------+------------------+----------+-----------------------------------+------------------------------------+---------------------------------------+
```
---
### [Image Scan](#image)

NONE