### [Source Code Scan](#source)

IGNORED: We are not using the SSH features of golang.org/x/crypto
```
2021-08-30T07:10:13.085-0700    INFO    Detected OS: unknown
2021-08-30T07:10:13.085-0700    INFO    Number of PL dependency files: 1
2021-08-30T07:10:13.085-0700    INFO    Detecting gomod vulnerabilities...

go.sum
======
Total: 2 (UNKNOWN: 0, LOW: 0, MEDIUM: 0, HIGH: 2, CRITICAL: 0)

+---------------------+------------------+----------+-----------------------------------+------------------------------------+---------------------------------------+
|       LIBRARY       | VULNERABILITY ID | SEVERITY |         INSTALLED VERSION         |           FIXED VERSION            |                 TITLE                 |
+---------------------+------------------+----------+-----------------------------------+------------------------------------+---------------------------------------+
| golang.org/x/crypto | CVE-2020-29652   | HIGH     | 0.0.0-20190426145343-a29dc8fdc734 | v0.0.0-20201216223049-8b5274cf687f | golang: crypto/ssh: crafted           |
|                     |                  |          |                                   |                                    | authentication request can            |
|                     |                  |          |                                   |                                    | lead to nil pointer dereference       |
|                     |                  |          |                                   |                                    | -->avd.aquasec.com/nvd/cve-2020-29652 |
+                     +------------------+          +                                   +------------------------------------+---------------------------------------+
|                     | CVE-2020-9283    |          |                                   | v0.0.0-20200220183623-bac4c82f6975 | golang.org/x/crypto: Processing       |
|                     |                  |          |                                   |                                    | of crafted ssh-ed25519                |
|                     |                  |          |                                   |                                    | public keys allows for panic          |
|                     |                  |          |                                   |                                    | -->avd.aquasec.com/nvd/cve-2020-9283  |
+---------------------+------------------+----------+-----------------------------------+------------------------------------+---------------------------------------+
```
---
### [Image Scan](#image)

NONE