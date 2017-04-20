package lingo

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// this file provides IO support and type safety for brown clusters.
// The creation of brownclusters is not done here.
// Right now lingo does not generate clusters - use PercyLiang's excellent tool for that

// Cluster represents a brown cluster
type Cluster int

// ReadCluster reads PercyLiang's cluster file format and returns a map of strings to Cluster
func ReadCluster(r io.Reader) map[string]Cluster {
	scanner := bufio.NewScanner(r)
	clusters := make(map[string]Cluster)

	for scanner.Scan() {
		line := scanner.Text()

		splits := strings.Split(line, "\t")
		var word string
		var cluster, freq int

		word = splits[1]

		var i64 int64
		var err error
		if i64, err = strconv.ParseInt(splits[0], 2, 64); err != nil {
			panic(err)
		}
		cluster = int(i64)

		if freq, err = strconv.Atoi(splits[2]); err != nil {
			panic(err)
		}

		// if clusterer has only seen a word a few times, then the cluster is not reliable
		if freq >= 3 {
			clusters[word] = Cluster(cluster)
		} else {
			clusters[word] = Cluster(0)
		}
	}

	// expand clusters with recasing
	for word, clust := range clusters {
		lowered := strings.ToLower(word)
		if _, ok := clusters[lowered]; !ok {
			clusters[lowered] = clust
		}

		titled := strings.ToTitle(word)
		if _, ok := clusters[titled]; !ok {
			clusters[titled] = clust
		}

		uppered := strings.ToUpper(word)
		if _, ok := clusters[uppered]; !ok {
			clusters[uppered] = clust
		}
	}

	return clusters
}
