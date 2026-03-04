# opnlab Design

## Architecture
- Backend: Go
- Frontend: TBD
- Modular provider pattern

## Providers
- jail: FreeBSD jail management
- system: CPU, RAM, disk, ZFS  
- network: interfaces, bandwidth
- backup: ZFS snapshots, backup jobs
- health: service health checks
- custom: user scripts

## API
- REST for actions
- WebSocket for real-time updates
