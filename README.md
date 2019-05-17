A Docker image with Kubernetes manifests for investigation and troubleshooting your cluster.

[![Build Status](https://travis-ci.org/digitalocean/doks-debug.svg?branch=master)](https://travis-ci.org/digitalocean/doks-debug)

# Purpose

The DOKS team provides this image for use as-is and for transparency as the image used when a request to "deploy a debug pod" is made to our customers, which may occur when deeper investigation is needed with direct access to a cluster.

# Usage

The included DaemonSet manifest will:

 1. Ensure a pod with our Docker image is running indefinitely on every node.
 2. Use `hostPID`, `hostIPC`, and `hostNetwork`.
 3. Mount the entire host filesystem to `/host` in the containers.
 4. Mount `/var/run/docker.sock` from the host.

In order to make use of these workloads, you can exec into a pod of choice by name:

```bash
kubectl exec -it my-pod-name bash
```

If you know the specific node name that you're interested in, you can exec into the debug pod on that node with:

```bash
NODE_NAME="my-node-name"
POD_NAME=$(kubectl get pods --field-selector spec.nodeName=${NODE_NAME} -ojsonpath='{.items[0].metadata.name}')
kubectl exec -it ${POD_NAME} bash
```

Once you're in, you have access to the set of tools listed in the `Dockerfile`. This includes:

 - [`vim`](https://github.com/vim/vim) - is a greatly improved version of the good old UNIX editor Vi. 
 - [`screen`](https://www.gnu.org/software/screen/) - is a full-screen window manager that multiplexes a physical terminal between several processes, typically interactive shells.
 - [`curl`](https://github.com/curl/curl) - is a command-line tool for transferring data specified with URL syntax.
 - [`jq`](https://github.com/stedolan/jq) - is a lightweight and flexible command-line JSON processor.
 - [`dnsutils`](https://packages.debian.org/stretch/dnsutils) - includes various client programs related to DNS that are derived from the BIND source tree, specifically [`dig`](https://linux.die.net/man/1/dig), [`nslookup`](https://linux.die.net/man/1/nslookup), and [`nsupdate`](https://linux.die.net/man/8/nsupdate).
 - [`iputils-ping`](https://packages.debian.org/stretch/iputils-ping) - includes the [`ping`](https://linux.die.net/man/8/ping) tool that sends ICMP `ECHO_REQUEST` packets to a host in order to test if the host is reachable via the network.
 - [`tcpdump`](https://www.tcpdump.org/) - a powerful command-line packet analyzer; and libpcap, a portable C/C++ library for network traffic capture.
 - [`traceroute`](https://linux.die.net/man/8/traceroute) - tracks the route packets taken from an IP network on their way to a given host.
 - [`net-tools`](https://packages.debian.org/stretch/net-tools) - includes the important tools for controlling the network subsystem of the Linux kernel, specifically [`arp`](http://man7.org/linux/man-pages/man8/arp.8.html), [`ifconfig`](https://linux.die.net/man/8/ifconfig), and [`netstat`](https://linux.die.net/man/8/netstat).
 - [`netcat`](https://linux.die.net/man/1/nc) - is a multi-tool for interacting with TCP and UDP; it can open TCP connections, send UDP packets, listen on arbitrary TCP and UDP ports, do port scanning, and deal with both IPv4 and IPv6.
 - [`iproute2`](https://wiki.linuxfoundation.org/networking/iproute2) - is a collection of utilities for controlling TCP / IP networking and traffic control in Linux.
 - [`strace`](https://github.com/strace/strace) - is a diagnostic, debugging and instructional userspace utility with a traditional command-line interface for Linux. It is used to monitor and tamper with interactions between processes and the Linux kernel, which include system calls, signal deliveries, and changes of process state.
 - [`docker`](https://docs.docker.com/engine/reference/commandline/cli/) - is the CLI tool used for interacting with Docker containers on the system.

 # Contributing

 At DigitalOcean we value and love our community! If you have any issues or would like to contribute, feel free to open an issue or PR and cc any of the maintainers.