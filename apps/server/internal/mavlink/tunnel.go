package mavlink

// TunnelPayload mirrors the custom MAVLink TUNNEL message (payload_type
// 0x8100) already validated in the existing Python prototype
// (mavlink_scraper.py) and the QGC CustomInstrumentWidget C++ fork.
// Struct layout: <ffffB (little-endian): cpu_pct, ram_pct, disk_pct,
// cpu_temp, node_mask.
type TunnelPayload struct {
	CPUPct   float32
	RAMPct   float32
	DiskPct  float32
	CPUTempC float32
	NodeMask uint8
}

var WatchedNodeNames = []string{
	"obstacle_avoidance",
	"path_planner",
	"camera_driver",
	"lidar_proc",
}

// DecodeTunnelPayload unpacks the 17-byte custom TUNNEL payload.
// TODO: implement with encoding/binary, matching struct.unpack_from("<ffffB", ...)
// from the existing mavlink_scraper.py implementation.
func DecodeTunnelPayload(raw []byte) (TunnelPayload, error) {
	panic("not implemented — port decode_tunnel() from mavlink_scraper.py")
}
