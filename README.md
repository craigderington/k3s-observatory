# Observatory - Home Lab k3s Visualization

A beautiful, real-time 3D visualization dashboard for monitoring Kubernetes clusters.

## Project Structure

```
observatory/
├── backend/        # Go backend (API + WebSocket server)
├── frontend/       # React + Three.js frontend
├── k3s/            # Kubernetes manifests

```

## Quick Start

### Backend (Go)
```bash
cd backend
go run cmd/observatory/main.go
```

### Frontend (React)
```bash
cd frontend
npm install
npm run dev
```

## API Documentation

### REST Endpoints

#### GET /api/nodes
Fetch all nodes in the cluster.

**Response:** Array of Node objects with position, status, and resource info.

#### GET /api/pods
Fetch all pods across all namespaces.

**Response:** Array of Pod objects with container details, status, and 3D position.

### WebSocket

**Connect:** `ws://localhost:8000/ws`

**Client → Server:**
- `{"type": "ping"}` - Heartbeat (send every 30s to keep connection alive)

**Server → Client Events:**
- `pod_added` - New pod created
- `pod_modified` - Pod status/spec changed
- `pod_deleted` - Pod removed
- `node_added` - Node joined cluster
- `node_modified` - Node status changed
- `node_deleted` - Node removed

Each event contains a `data` field with the affected resource.

## Features

✅ Real-time 3D visualization
✅ Live WebSocket updates
✅ Toast notifications
✅ Namespace filtering
✅ Connection status indicator
✅ Interactive 3D controls

