# Abstract

The cluster consists of nodes, each capable of taking on different roles:

- cluster manager (or manager node or simply manager) and
- executor node (or executor role or simply executor).

The manager is elected using the Raft consensus protocol.

A cluster exists only if at least one node is present. The first online node becomes the de facto manager as it is the
only node forming the cluster. As more nodes join the cluster, the Raft consensus ensures that only one manager remains
in control.

The cluster state is maintained between nodes using a single management channel on port 8972. When a node intends to
join the cluster, it must announce its arrival on this channel, i.e. this is Node Discovery. Since there might be
multiple clusters operating on the same network, the node will sends its identification details, ensuring that only the
appropriate cluster manager can process the request.

Once the cluster verifies the authenticity and the uniqueness of the node, the node must authenticate and be authorized
to join.

The manager accepts or rejects the request to join. When the node is accepted, then it can contribute to the cluster.
All nodes are updated with the cluster state.

# Details

Here are the steps of the protocol.

## The first node of a cluster

- Node broadcasts on 8972 with a cluster ID
- No answer after 10s the broadcast timeout
- Node sets up itself as the manager of the cluster ID

## Cluster ID exists

- Node broadcasts on 8972 with the cluster ID
- Cluster manager registers the node in the cluster state 
- manager sends back an ACK response
