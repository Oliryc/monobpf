#!/bin/bash
source /opt/ros/melodic/setup.bash
export TURTLEBOT3_MODEL=waffle
export EDITOR=/usr/bin/vim
export PCP_PMDAS_DIR=/var/lib/pcp/pmdas/
export PCP_LOG_DIR=/var/log/pcp/pmcd/
export PCP_DIR=/
export PCP_SHARE_DIR=/usr/share/pcp/
roscore 1>/dev/null &
export CORE_PID=$!
echo "ROSCORE STARTED WITH PID $CORE_PID"
roslaunch rospy_tutorials talker_listener.launch 1>/dev/null &
export LAUNCH_PID=$!
echo "ROSCORE STARTED WITH PID $LAUNCH_PID"
sleep 5
go run *.go &
export APP_PID=$!
echo "PROCESS STARTED WITH PID $APP_PID"
sleep 30
kill -2 $APP_PID
kill -2 $LAUNCH_PID
kill -2 $CORE_PID
sleep 5
kill -9 $APP_PID
kill -9 $LAUNCH_PID
kill -9 $CORE_PID
