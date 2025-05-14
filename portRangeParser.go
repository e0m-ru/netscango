package main

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

// ParsePortRanges parses a string of ports separated by '-' or ',' and returns a sorted slice of unique ports.
func ParsePortRanges(portStr string) ([]int, error) {
	seen := make(map[int]bool)
	var ports []int

	ranges := strings.Split(portStr, ",")

	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if strings.Contains(r, "-") {
			start, end, err := parseRange(r)
			if err != nil {
				return nil, err
			}
			for port := start; port <= end; port++ {
				if err := addPort(port, seen, &ports); err != nil {
					return nil, err
				}
			}
		} else {
			port, err := strconv.Atoi(r)
			if err != nil {
				return nil, errors.New("invalid port: " + r)
			}
			if err := addPort(port, seen, &ports); err != nil {
				return nil, err
			}
		}
	}

	sort.Ints(ports)
	return ports, nil
}

func parseRange(r string) (int, int, error) {
	parts := strings.Split(r, "-")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid range format: " + r)
	}

	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, errors.New("invalid start of range: " + parts[0])
	}

	end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, errors.New("invalid end of range: " + parts[1])
	}

	if start > end {
		start, end = end, start
	}

	return start, end, nil
}

func addPort(port int, seen map[int]bool, ports *[]int) error {
	if port < 1 || port > MAX_PORT {
		return errors.New("port out of range: " + strconv.Itoa(port))
	}
	if !seen[port] {
		*ports = append(*ports, port)
		seen[port] = true
	}
	return nil
}
