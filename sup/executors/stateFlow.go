func updateCluster() {
	for {
		if cluster.targetPaasVersion == cluster.currentPaasVersion {
			continue
		}
		if cluster.curState != PaasReady {
			continue
		}
		if cluster.Nodes[].curState != InfraReady {
			continue
		}
		if cluster.Nodes[].InfraVersion != cluster.targetInfraVersion {
			continue
		}
		cluster.curInfraVersion = cluster.targetInfraVersion
		cluster.curState = UpdatingPaas
	}
}

func resizeCluster() {
	if cluster.curState != PaasReady {
		return
	}
	// all current nodes should be in steady state
	if cluster.Nodes[].curState != InfraReady {
		return
	}
	// a node was either removed or added
	if cluster.Nodes[] != getNodes(InfraReady)  {
		cluster.curState = Resizing
	}
}

func settingUpPaasCluster() {
	if cluster.curState != BootStrapping {
		return
	}
	if len(getNodes(InfraReady)) < minSize {
		return
	}
	cluster.curState = settingUpPaas
}

func updatingInfraNode() {
	for {
		if cluster.curState != PaasReady {
			continue
		}
		for node in getNodes(InfraReady) {
			if node.InfraVersion == cluster.InfraVersion {
				continue
			}
			node.curState = UpdatingInfra
		}
	}
}

func addingInfraNode() {
	if node.curState != Onboarding  {
		return
	}

	if cluster.curState != BootStrapping && cluster.curState != PaasReady {
		return
	}

	node.curState = SettingUpInfra
}

func settingUpInfra() {
	for {
		result = execute("infra.yaml", getNodes(SettingUpInfra))
	    for node in result.Nodes {
			node.curState =   
		}
	}	
}

func updatingInfra() {
	for {
		result = execute("update_infra.yaml", cluster.targetInfraVersion, getNodes(UpdatingInfra))
	    for node in result.Nodes {
			node.InfraVersion = cluster.targetInfraVersion
			node.curState = InfraReady
		}
	}	
}

func updatingPaas() {
	for {
		if cluster.curState != UpdatingPaas {
			continue
		}
		result = execute("update_paas.yaml")
		if result.Success {
			cluster.curState = PaasReady
		}
	}
}

func settingUpPaas() {
	for {
		if cluster.curState != SettingUpPaas {
			continue
		}
		nodes = getNodes(InfraNode)
		result = execute("rke_start.yaml", nodes)
		if result.Success {
			cluster.curState = PaasReady
		}
	}
}


func resizingPaas() {
	for {
		if cluster.curState != Resizing {
			continue
		}
		newNodes = getNodes(InfraReady) - cluster.Nodes[]
		if newNodes {
			result = execute("rke_add_node.yaml", newNodes)
			if result.Success {
				cluster.Nodes[] += newNodes
			} else {
				continue
			}
		}

		removedNodes = cluster.Nodes[] - getNodes(Deleting)
		if removedNodes {
			result = execute("rke_remove_node.yaml", newNodes)
			if result.Success {
				cluster.Nodes[] -= removedNodes
			} else {
				continue
			}
		}
		cluster.curState = PaasReady
	}
}

func deletingInfraNode(node) {
	// can't take node out as it will lower the cluster size below minSize
	if node.curState == InfraReady && len(getNodes(InfraReady))-1 < minSize {
		return
	}

	node.curState = Deleting
	node.targetState = Deleted
}

func resizingInfra() {
	for {
		toBeDeletedNodes = getNodes(Deleting)
		for node in getNodes(Deleting) {
			if cluster.Nodes[].Contains(node) {
				continue
			}
			node.curState = Deleted
		}
	}
}

func updateInfraVersion(V) {
	cluster.InfraVersion = V
}

func updatePaasVersion(V) {
	cluster.PaasVersion = V
}