#!/bin/bash

# Ensure the 'sale-description-generator' session exists
if ! tmux has-session -t sale-description-generator 2>/dev/null; then
  # Create a new session named 'sale-description-generator', detached
  cd /workspaces/seraphine
  tmux new-session -d -s sale-description-generator
  
  # Split the window horizontally (-h)
  # The left pane will remain a terminal
  # The right pane will run 'gh dash'
  tmux split-window -h -t sale-description-generator "gh dash"
fi
