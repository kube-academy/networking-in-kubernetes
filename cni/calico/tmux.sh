#!/bin/bash
set -- $(stty size) # $1 = rows $2 = columns

SESSION="Calico"
SESSIONEXISTS=$(tmux list-sessions | grep $SESSION)

if [ "$SESSIONEXISTS" = "" ]
then
    tmux -2 new-session -d -s $SESSION -x "$2" -y "$(($1 - 1))"

    tmux split-window -v
    tmux select-pane -t 0
    tmux send-keys 'watch -c -t --differences kubectl get nodes -o wide' C-m

    tmux split-window -v
    tmux send-keys 'watch -c -t --differences kubectl get pod -o wide' C-m

    tmux split-window -v
    tmux select-pane -t 2
    tmux send-keys 'watch -c -t --differences docker exec -it calico-worker ip -c r' C-m

    tmux split-window -h
    tmux select-pane -t 3
    tmux send-keys 'watch -c -t --differences docker exec -it calico-worker2 ip -c r' C-m

    tmux split-window -h
    tmux select-pane -t 4
    tmux send-keys 'watch -c -t --differences docker exec -it calico-worker3 ip -c r' C-m
fi

tmux resize-pane -t 0 -y 5
tmux resize-pane -t 1 -y 7
tmux resize-pane -t 2 -y 10 -x 60
tmux resize-pane -t 3 -y 10 -x 60
tmux resize-pane -t 4 -y 10 -x 60

tmux select-pane -t 5

if [ "$SESSIONEXISTS" = "" ]
then
    tmux attach-session -t $SESSION:0
fi
