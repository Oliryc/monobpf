#!/bin/bash
source /opt/ros/melodic/setup.bash
export TURTLEBOT3_MODEL=waffle
export EDITOR=/usr/bin/vim
export PCP_PMDAS_DIR=/var/lib/pcp/pmdas/
export PCP_LOG_DIR=/var/log/pcp/pmcd/
export PCP_DIR=/
export PCP_SHARE_DIR=/usr/share/pcp/
go run *.go
