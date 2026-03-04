# opnlab

> A lightweight, modular home lab dashboard and control center for FreeBSD homelabs.

## Overview

opnlab is a modern dashboard for managing FreeBSD-based homelabs. It provides real-time monitoring and control of jails, system resources, backups, and custom services — all in a fast, modular interface.

## Goals

- **Lightweight** — Single binary, minimal dependencies
- **Fast** — Quick load times, efficient polling
- **Modular** — Easy to add new service providers
- **Dashboard + Control** — View status AND take action

## Providers

| Provider | Description | Actions |
|----------|-------------|---------|
| `jail` | FreeBSD jail management | start, stop, restart, console |
| `system` | CPU, RAM, disk, ZFS | - |
| `network` | Interfaces, bandwidth | - |
| `backup` | Backup jobs, ZFS snapshots | restore, snapshot, run |
| `health` | Service health checks | - |
| `custom` | User-defined scripts | (custom) |

## Architecture

```
┌──────────────────────────────────────────────────┐
│                   Frontend                      │
│              (React / HTMX / Vue)               │
└────────────────────┬────────────────────────────┘
                     │ REST + WebSocket
┌────────────────────▼────────────────────────────┐
│                     API                         │
│              (Go + Gin/Echo)                    │
└────────────────────┬────────────────────────────┘
                     │
    ┌────────────────┼────────────────┐
    │                │                │
┌───▼────┐    ┌─────▼────┐    ┌─────▼────┐
│ Jail   │    │ System   │    │ Backup   │
│Provider│    │ Provider │    │ Provider │
└────────┘    └──────────┘    └──────────┘
```

## Tech Stack

- **Backend:** Go
- **Frontend:** TBD (React, Svelte, or HTMX)
- **Database:** SQLite (embedded) or time-series for metrics
- **Deployment:** Single binary + static files in a FreeBSD jail

## Getting Started

Work in progress.

## License

MIT
