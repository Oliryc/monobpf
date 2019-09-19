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
sleep 5
argv=("$@")
val=${argv[0]}
echo "Launching $val nodes"
for ((i=1; i<=$val; i++))
do
rosrun rospy_tutorials talker.py 1>/dev/null &
TALK_PID[$i]=$!
rosrun rospy_tutorials listener.py 1>/dev/null &
#roslaunch rospy_tutorials talker_listener.launch 1>/dev/null &
LISTEN_PID[$i]=$!
done
export LAUNCH_PID
echo "Talker STARTED WITH PID ${TALK_PID[@]}"
echo "Listener STARTED WITH PID ${LISTEN_PID[@]}"
sleep 5
rostopic bw /chatter &
export APP_PID=$!
echo "PROCESS STARTED WITH PID $APP_PID"
sleep 100
kill -2 $APP_PID
for ((i=1; i<=$val; i++))
do
kill -2 ${TALK_PID[$i]}
kill -2 ${LISTEN_PID[$i]}
done
kill -2 ${TALK_PID[@]}
kill -2 ${LISTEN_PID[@]}
kill -2 $CORE_PID
sleep 5
kill -9 $APP_PID
for ((i=1; i<=$val; i++))
do
kill -9 $TALK_PID[$i]
kill -9 $LISTEN_PID[$i]
done
kill -9 $CORE_PID
kill -9 ${TALK_PID[@]}
kill -9 ${LISTEN_PID[@]}

rosclean purge -y
