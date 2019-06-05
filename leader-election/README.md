# Cluster IP leader election
This POC performs cluster IP set up using leader election based on etcd. Essentially:
1. Use etcd for consensus
2. Use alias net interfaces for cluster IP
3. Use systemd as an orchestrator
