# GitHub-User-Activity-CLI
A simple command-line interface (CLI) to fetch the recent activity of a GitHub user and display it in the terminal. 
This project is inspired by https://roadmap.sh/projects/github-user-activity to learn Go better.

## Prerequesties
`Go` is needed for this project (`go 1.26.2`)

## Installation
### Linux
1. Clone the repository:
```
git clone https://github.com/momin-chezgi/GitHub-User-Activity-CLI.git
```
2. Build the executable:
```
go build -o github-user-activity
```
3. Create the local binaries directory if it doesn't already exist:
```
mkdir -p ~/.local/bin
```

4. Copy the executable to it:
```
cp github-user-activity ~/.local/bin/
chmod +x ~/.local/bin/github-user-activity
```

5. Add ~/.local/bin to your PATH:
```
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

You can now run the program from anywhere:
```
github-user-activity
```


### Windows
1. Clone the repository:
```
git clone https://github.com/momin-chezgi/GitHub-User-Activity-CLI.git
cd GitHub-User-Activity-CLI
```

2. Build the executable:
```
go build -o github-user-activity.exe
```

3. Run the program:
```
.\github-user-activity.exe
```

To run the program from anywhere, add the directory containing github-user-activity.exe to your Windows PATH.
