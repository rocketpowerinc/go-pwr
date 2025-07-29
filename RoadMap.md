# Features
- Make it so anyone can easily change the repo it clones to explore


# FYI
- Can't have the script run in the preview pane because gum menus wouldn't work
- Syntax Highling in preview pane broke boarders

# Bugs to fix

- Keyboard works good on Windows,Linux but Mac I have to use Iterm2 terminal emulator

# Test

- Test App on a Ubuntu Server

# Build

- Tag Templet in Scriptbin repo for bash and pwsh
  - works on Server/Desktop or both
  - works on Windows/Mac/Linux or all

# Script Bin Structure

- Bash
  - Mac/Linux
  - Mac
  - Linux
  - Server
  - Docker
- Pwsh
  - Universal
  - Windows
  - Mac
  - Linux

# My Scripts

- Alot of them will use gum -choose

# AI Continuous Help

- Ask it to optimise
- Ask it to look for vulnerabilty
- Ask it for error checking and dubugging
- Ask it to explain the main.go line by line
- Ask it for cool new features, but that will also work on a ubuntu server

# To Do
- Reupload go-pwr to github once all commits are complete so it does not look so sloppy

# AI Ideas

Search & Filtering Features
- **Fuzzy Search**: Press `/` to search script names/paths in real-time
- **Tag/Category Filtering**: Filter scripts by tags like #server, #backup, #docker
- **Recent Scripts**: Show recently run scripts at the top
- **Favorites**: Star frequently used scripts with `f` key

Enhanced Script Management
- **Script History**: Track which scripts you've run and when
- **Execution Time**: Show how long scripts took to run
- **Script Dependencies**: Show if a script requires certain tools (gum, docker, etc.)
- **Dry Run Mode**: Preview what a script would do without running it

Terminal Integration
- **Multiple Repo Support**: Switch between different script repos (work, personal, etc.)
- **Script Arguments**: Pass arguments to scripts before running
- **Environment Variables**: Set env vars for script execution
- **Working Directory**: Choose which directory to run scripts from

Visual Enhancements
- **Script Previews**: Show script description/comments in preview
- **Execution Status**: Visual indicator for running/completed/failed scripts
- **Progress Bars**: For long-running scripts
- **Color Themes**: Different color schemes

Server-Friendly Features
- **SSH Mode**: Run scripts on remote servers
- **Batch Execution**: Queue multiple scripts to run in sequence
- **Log Viewer**: Built-in log viewer for script outputs
- **Health Checks**: Quick system status scripts