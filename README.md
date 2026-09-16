# GitHub-User-Activity-CLI
A simple command-line interface (CLI) to fetch the recent activity of a GitHub user and display it in the terminal. 
This project is inspired by https://roadmap.sh/projects/github-user-activity to learn Go better.

## Installation
### Linux
1. Build the executable:
```
go build -o github-user-activity
```
Or download it from here

2. Create the local binaries directory if it doesn't already exist:
```
mkdir -p ~/.local/bin
```

3. Copy the executable to it:
```
cp github-user-activity ~/.local/bin/
chmod +x ~/.local/bin/github-user-activity
```

4. Add ~/.local/bin to your PATH:
```
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

You can now run the program from anywhere:
```
github-user-activity
```
