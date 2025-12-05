# 🔭 Observatory - Real-Time 3D Kubernetes Visualization

> Transform your k3s/Kubernetes cluster monitoring from boring metrics into an engaging, animated 3D experience.

**Observatory** is a beautiful, real-time 3D visualization dashboard that lets you see your Kubernetes cluster's heartbeat. Watch pods spin up, scale, and disappear in stunning 3D space. No more endless `kubectl get pods` commands - just open your browser and observe your infrastructure come to life.

![k3s-observatory](https://raw.githubusercontent.com/craigderington/k3s-observatory/refs/heads/master/screenshots/screenshot-2025-12-05_14-26-12.png)


![Status](https://img.shields.io/badge/status-alpha-orange)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react)
![License](https://img.shields.io/badge/license-MIT-blue)

## ✨ Features

- **🌌 Immersive 3D Visualization** - Nodes as spheres, pods orbiting in space with Three.js/React Three Fiber
- **⚡ Real-Time Updates** - WebSocket-powered live updates with zero page refreshes
- **🎨 Color-Coded Status** - Instant visual feedback (🟢 Running, 🔵 Pending, 🔴 Failed, ⚠️ Warning)
- **📢 Smart Toast Notifications** - Non-intrusive alerts for pod/node lifecycle events
- **🏷️ Namespace Filtering** - Focus on what matters with dropdown namespace selection
- **🔄 Dynamic Redistribution** - Pods smoothly reposition when scaled up/down
- **💚 Connection Health** - Live/Offline indicator with automatic reconnection
- **🖱️ Interactive Controls** - Orbit, zoom, and pan through your cluster
- **🎯 Zero Config** - Works out of the box with your existing kubeconfig

## 📸 Screenshots

> *Coming soon! Add screenshots/GIFs of your Observatory in action*

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** - Backend runtime
- **Node.js 18+** - Frontend tooling
- **k3s/Kubernetes cluster** - The cluster you want to visualize
- **kubectl configured** - With access to your cluster

### Installation

**1. Clone the repository**
```bash
git clone https://github.com/craigderington/k3s-observatory.git
cd k3s-observatory
```

**2. Set up the backend**
```bash
cd backend
go mod download

# Point to your kubeconfig
export KUBECONFIG=/path/to/your/k3s.yaml

# Run the backend
go run cmd/observatory/main.go
```

The backend will start on **http://localhost:8000** and begin watching your cluster.

**3. Set up the frontend** (in a new terminal)
```bash
cd frontend
npm install
npm run dev
```

The frontend will start on **http://localhost:3000**.

**4. Open your browser**

Navigate to **http://localhost:3000** and watch your cluster come to life! 🎉

## 🎯 Usage

### Basic Operations

- **🖱️ Navigate**: Left-click and drag to orbit the camera
- **🔍 Zoom**: Scroll wheel to zoom in/out
- **🏷️ Filter**: Use the namespace dropdown to focus on specific namespaces
- **👆 Inspect**: Hover over pods/nodes to see their names
- **📊 Monitor**: Watch the Live status indicator and pod/node counts in the header

### Watching Real-Time Changes

Deploy something to your cluster and watch it appear:

```bash
# Scale a deployment
kubectl scale deployment my-app --replicas=5

# Delete a pod
kubectl delete pod my-pod-abc123

# Deploy a new app
kubectl apply -f deployment.yaml
```

You'll see pods appear, transition through states (Pending → Running), and redistribute smoothly around their nodes!

### Test Deployment

Try the included test deployment:

```bash
kubectl apply -f k8s/test-deployment.yaml
```

This creates 3 nginx pods in the `observatory-test` namespace. Watch them appear in real-time!

## 🐳 Docker Deployment

### Using Docker Compose (Recommended)

The easiest way to run Observatory is with Docker Compose:

**1. Set your kubeconfig location**
```bash
export KUBECONFIG=/path/to/your/k3s.yaml
```

**2. Start Observatory**
```bash
docker-compose up -d
```

**3. Access the UI**

Open **http://localhost:3000** in your browser!

**4. View logs**
```bash
# Both services
docker-compose logs -f

# Just backend
docker-compose logs -f backend

# Just frontend
docker-compose logs -f frontend
```

**5. Stop Observatory**
```bash
docker-compose down
```

### Building Individual Containers

**Backend:**
```bash
cd backend
docker build -t observatory-backend .
docker run -d \
  -p 8000:8000 \
  -v $KUBECONFIG:/root/.kube/config:ro \
  --name observatory-backend \
  observatory-backend
```

**Frontend:**
```bash
cd frontend
docker build -t observatory-frontend .
docker run -d \
  -p 3000:80 \
  --name observatory-frontend \
  observatory-frontend
```

### Container Features

- ✅ **Multi-stage builds** - Small, optimized images
- ✅ **Health checks** - Docker monitors service health
- ✅ **Auto-restart** - Containers restart on failure
- ✅ **Nginx proxy** - Frontend proxies API/WebSocket to backend
- ✅ **Volume mounts** - Kubeconfig mounted read-only

### Environment Variables

**Backend:**
- `PORT` - Backend port (default: 8000)
- `KUBECONFIG` - Path to kubeconfig (mounted as volume)

**Frontend:**
- No environment variables needed (configured via Nginx)

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────┐
│  Browser (React + Three.js)                     │
│  ├─ 3D Scene Renderer (React Three Fiber)       │
│  ├─ WebSocket Client (Auto-reconnect)           │
│  └─ Control Panel UI (React + CSS)              │
└─────────────────┬───────────────────────────────┘
                  │ WebSocket + REST
┌─────────────────▼───────────────────────────────┐
│  Observatory Backend (Go)                       │
│  ├─ WebSocket Hub (gorilla/websocket)           │
│  ├─ Kubernetes Watchers (client-go)             │
│  ├─ REST API Endpoints                          │
│  └─ Event Broadcaster                           │
└─────────────────┬───────────────────────────────┘
                  │ Kubernetes API
┌─────────────────▼───────────────────────────────┐
│  k3s/Kubernetes Cluster                         │
│  ├─ Nodes (compute resources)                   │
│  ├─ Pods (running workloads)                    │
│  └─ Watch API (real-time events)                │
└─────────────────────────────────────────────────┘
```

### How It Works

1. **Backend connects** to your k3s/Kubernetes cluster using your kubeconfig
2. **Watchers monitor** nodes and pods for any changes (add/modify/delete)
3. **WebSocket broadcasts** events to all connected frontend clients in real-time
4. **Frontend receives** events and updates the 3D visualization smoothly
5. **Dynamic positioning** recalculates pod orbits around nodes when the cluster changes

### Technology Stack

**Backend:**
- Go 1.21+
- [gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket server
- [client-go](https://github.com/kubernetes/client-go) - Kubernetes API client

**Frontend:**
- React 18
- [Three.js](https://threejs.org/) - 3D rendering engine
- [React Three Fiber](https://docs.pmnd.rs/react-three-fiber) - React renderer for Three.js
- [@react-three/drei](https://github.com/pmndrs/drei) - Useful helpers for R3F
- [Vite](https://vitejs.dev/) - Build tool and dev server
- TypeScript - Type safety

## 📡 API Documentation

### REST Endpoints

#### `GET /api/health`
Health check endpoint.

**Response:**
```json
{
  "status": "ok"
}
```

#### `GET /api/nodes`
Fetch all nodes in the cluster with their 3D positions.

**Response:**
```json
[
  {
    "id": "node-uid-123",
    "name": "node1",
    "status": "Ready",
    "cpu": {
      "used": 0,
      "total": 4.0
    },
    "memory": {
      "used": 0,
      "total": 8.0
    },
    "pods": [],
    "labels": {
      "kubernetes.io/hostname": "node1"
    },
    "position": {
      "x": 10.0,
      "y": 0.0,
      "z": 0.0
    }
  }
]
```

#### `GET /api/pods`
Fetch all pods across all namespaces with their 3D positions.

**Response:**
```json
[
  {
    "id": "pod-uid-456",
    "name": "nginx-deployment-abc123",
    "namespace": "default",
    "status": "Running",
    "nodeName": "node1",
    "containers": [
      {
        "name": "nginx",
        "status": "Running",
        "restarts": 0
      }
    ],
    "createdAt": "2025-01-15T10:30:00Z",
    "position": {
      "x": 12.5,
      "y": 0.0,
      "z": 1.2
    }
  }
]
```

### WebSocket

#### Connection

Connect to `ws://localhost:8000/ws` to receive real-time cluster events.

#### Client → Server Messages

**Heartbeat (Ping):**
```json
{
  "type": "ping"
}
```

Send every 30 seconds to keep the connection alive. The backend has a 60-second read timeout.

#### Server → Client Events

All events follow this structure:

```json
{
  "type": "event_type",
  "data": {
    "pod": { /* Pod object */ },
    "node": { /* Node object */ }
  }
}
```

**Event Types:**

- `pod_added` - New pod created
- `pod_modified` - Pod status/spec changed (e.g., Pending → Running)
- `pod_deleted` - Pod removed
- `node_added` - Node joined cluster
- `node_modified` - Node status changed (e.g., resource usage, conditions)
- `node_deleted` - Node removed from cluster

**Example:**
```json
{
  "type": "pod_added",
  "data": {
    "pod": {
      "id": "pod-uid-789",
      "name": "my-app-xyz",
      "namespace": "production",
      "status": "Pending",
      "nodeName": "",
      "position": { "x": 0, "y": 0, "z": 0 }
    }
  }
}
```

## 🛠️ Development

### Project Structure

```
k3s-observatory/
├── backend/
│   ├── cmd/observatory/      # Main entry point
│   │   └── main.go
│   ├── internal/
│   │   ├── k8s/             # Kubernetes client, watchers, data models
│   │   │   ├── client.go
│   │   │   ├── nodes.go
│   │   │   ├── pods.go
│   │   │   └── watcher.go
│   │   └── websocket/       # WebSocket hub and client management
│   │       ├── hub.go
│   │       └── client.go
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── components/      # React components
│   │   │   ├── Scene.tsx            # Main 3D scene
│   │   │   ├── NodeSphere.tsx       # Node visualization
│   │   │   ├── PodSphere.tsx        # Pod visualization
│   │   │   └── ToastContainer.tsx   # Toast notifications
│   │   ├── hooks/           # Custom React hooks
│   │   │   └── useWebSocket.ts      # WebSocket connection logic
│   │   ├── services/        # API clients
│   │   │   └── api.ts
│   │   ├── types/           # TypeScript types
│   │   │   └── index.ts
│   │   ├── App.tsx          # Main app component
│   │   ├── App.css          # Styles
│   │   └── main.tsx         # Entry point
│   ├── package.json
│   └── vite.config.ts
├── k8s/
│   └── test-deployment.yaml # Test deployment for demos
├── CLAUDE.md                # Detailed project vision & roadmap
└── README.md
```

### Running in Development Mode

**Backend with hot reload:**
```bash
cd backend
go run cmd/observatory/main.go
```

**Frontend with hot reload:**
```bash
cd frontend
npm run dev
```

Both will automatically reload when you make changes!

### Building for Production

**Backend:**
```bash
cd backend
go build -o observatory ./cmd/observatory
KUBECONFIG=/path/to/k3s.yaml ./observatory
```

**Frontend:**
```bash
cd frontend
npm run build
# Output in dist/ folder
```

## 🐛 Troubleshooting

### Backend won't connect to cluster

**Issue:** `Failed to connect to Kubernetes cluster`

**Solution:** Make sure your `KUBECONFIG` environment variable points to a valid kubeconfig file:

```bash
export KUBECONFIG=/path/to/your/k3s.yaml
```

For k3s clusters, the kubeconfig is usually at `/etc/rancher/k3s/k3s.yaml` on the server.

### WebSocket keeps disconnecting

**Issue:** Connection status shows "Offline" repeatedly

**Solution:**
1. Check that the backend is running on port 8000
2. Ensure no firewall is blocking WebSocket connections
3. Check browser console for specific error messages
4. Backend logs will show connection/disconnection events

### Pods appear at the center and don't move

**Issue:** Pods spawn at `(0, 0, 0)` and stay there

**Solution:** This happens when pods haven't been assigned to a node yet (still in Pending state). Once Kubernetes schedules them to a node, they'll receive a `pod_modified` event and move to orbit their assigned node.

### Nodes disappearing from view

**Issue:** Nodes vanish when pods are added/removed

**Solution:** This was fixed in recent updates. Hard refresh your browser (`Ctrl+Shift+R`) to ensure you have the latest frontend code.

## 🗺️ Roadmap

See [CLAUDE.md](./CLAUDE.md) for the detailed project vision and feature roadmap.

**Phase 1: MVP** ✅ Complete
- [x] Basic 3D visualization of nodes and pods
- [x] Real-time WebSocket updates
- [x] Toast notifications
- [x] Namespace filtering
- [x] Dynamic pod redistribution

**Phase 2: Resource Visualization** (In Progress)
- [ ] CPU usage shown as pulsing intensity/glow
- [ ] Memory usage shown as size scaling
- [ ] Network traffic as particle streams between pods
- [ ] Resource graphs in detail panel

**Phase 3: Historical & Analysis**
- [ ] Store metrics history in database
- [ ] Timeline scrubber to replay past states
- [ ] Playback speed controls
- [ ] Incident bookmarking

**Phase 4: Polish & Production**
- [ ] Multiple visualization modes
- [ ] Custom labels and annotations
- [ ] Screenshot/screen recording
- [ ] Mobile-responsive view
- [ ] Dark/light themes

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Built with ❤️ for the homelab community
- Inspired by the desire to make Kubernetes monitoring beautiful and intuitive
- Special thanks to the teams behind Three.js, React Three Fiber, and client-go

## 📞 Support

- 🐛 **Issues**: [GitHub Issues](https://github.com/craigderington/k3s-observatory/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/craigderington/k3s-observatory/discussions)

---

**Made with 🔭 by [Craig Derington](https://github.com/craigderington)**

*Transform your cluster into a work of art.*
