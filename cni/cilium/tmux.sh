#!/bin/bash
set -- $(stty size) # $1 = rows $2 = columns

SESSION="Cilium"
SESSIONEXISTS=$(tmux list-sessions | grep $SESSION)

if [ "$SESSIONEXISTS" = "" ]
then
    tmux -2 new-session -d -s $SESSION -x "$2" -y "$(($1 - 1))"

    tmux split-window -v
    tmux select-pane -t 0
    tmux send-keys 'watch -c -t --differences kubectl get pods' C-m

    tmux split-window -h
    tmux select-pane -t 1
    tmux send-keys 'watch -c -t --differences "echo -n cilium-worker ipables count: ; docker exec -it cilium-worker iptables -L | wc -l" ' C-m
fi

tmux resize-pane -t 1 -x 35 -y 5
tmux select-pane -t 2

if [ "$SESSIONEXISTS" = "" ]
then
    tmux attach-session -t $SESSION:0
fi
