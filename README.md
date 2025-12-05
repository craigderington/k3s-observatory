# Observatory - Home Lab k3s Visualization

A beautiful, real-time 3D visualization dashboard for monitoring Kubernetes clusters.

## Project Structure

```
observatory/
├── backend/        # Go backend (API + WebSocket server)
├── frontend/       # React + Three.js frontend
├── k8s/            # Kubernetes manifests
└── CLAUDE.md       # Detailed project documentation
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

## Documentation

See [CLAUDE.md](./CLAUDE.md) for detailed project vision, architecture, and development guidelines.
