// +build ignore

#include <linux/bpf.h>
#include <linux/in.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

struct {
    __uint(type, BPF_MAP_TYPE_DEVMAP);
    __type(key, int);
    __type(value, int);
    __uint(max_entries, 1);
} intf_map SEC(".maps");

SEC("xdp_redirect")
int xdp_redirect_func(struct xdp_md *ctx) {
    __u32 index = 0;
    __u32 *val = bpf_map_lookup_elem(&intf_map, &index);
    const char fmt[] = "Looked up value %d from index %d";
    bpf_trace_printk(fmt, sizeof(fmt), val, index);
    return bpf_redirect_map(&intf_map, index, 0);
}

char _license[] SEC("license") = "GPL";
