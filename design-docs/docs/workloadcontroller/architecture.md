---
sidebar_position: 1
---

# Architecture

```mermaid
sequenceDiagram
	participant frontend
	participant apiserver

	frontend->>apiserver: POST /app/:appID/deployments/new
	apiserver--)frontend: { "ok": true }
```
