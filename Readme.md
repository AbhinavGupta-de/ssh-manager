# 🔑 SSH Manager

**One command to manage all your SSH keys across any device.**

Never manually configure SSH keys again. Perfect for developers who juggle multiple accounts (work, personal, freelance) across different machines.

## ✨ Features

- 🚀 **One-line installation** - Clone and run, that's it
- 🔄 **Easy switching** - Switch between SSH profiles instantly
- 🌍 **Cross-platform** - Works on macOS, Linux, and Windows
- ⚡ **Zero dependencies** - Single binary, no runtime needed
- 🔧 **Auto-configuration** - Manages SSH config, Git config, and SSH agent
- 📦 **Portable** - Copy the binary anywhere, it just works

## 🎯 Quick Start

### Installation

**macOS / Linux:**
```bash
git clone https://github.com/abhinavgupta-de/ssh-manager.git
cd ssh-manager
chmod +x install.sh
./install.sh
```

**Windows (PowerShell):**
```powershell
git clone https://github.com/abhinavgupta-de/ssh-manager.git
cd ssh-manager
PowerShell -ExecutionPolicy Bypass -File install.ps1
```

That's it! The script will:
- ✅ Build the binary for your platform
- ✅ Install it globally
- ✅ Add it to your PATH
- ✅ Make it executable

### First Use

Create your first profile:
```bash
sshm new
```

Follow the prompts:
```
Profile name: work
Email: you@company.com
Git username: johndoe
Git host: github.com
```

The tool will:
1. Generate an SSH key pair
2. Display your public key
3. Save your profile

**Copy the displayed public key to your GitHub/GitLab account!**

Create more profiles:
```bash
sshm new  # Create 'personal' profile
```

## 📖 Usage

### Create a new profile
```bash
sshm new
```

### List all profiles
```bash
sshm list
```

### Switch between profiles
```bash
sshm switch work
sshm switch personal
```

When you switch, the tool automatically:
- Updates your SSH config
- Changes Git global username and email
- Adds the key to SSH agent

### Check current profile
```bash
sshm current
```

### Delete a profile
```bash
sshm delete old-work
```

## 🎬 Example Workflow

```bash
# Monday morning - time for work
sshm switch work
git clone git@github.com:company/project.git

# Evening - personal project time
sshm switch personal
git clone git@github.com:me/side-project.git

# Check which profile is active
sshm current
```

## 🔧 What It Does Behind the Scenes

When you run `sshm switch work`:

1. **Updates `~/.ssh/config`**:
   ```
   Host github.com
       IdentityFile ~/.ssh/id_rsa_work
       IdentitiesOnly yes
   ```

2. **Updates Git config**:
   ```bash
   git config --global user.name "Your Work Name"
   git config --global user.email "work@company.com"
   ```

3. **Loads key to SSH agent**:
   ```bash
   ssh-add ~/.ssh/id_rsa_work
   ```

## 💡 Why SSH Manager?

**Problem**: You have multiple GitHub accounts (work, personal, freelance). Every time you:
- Clone a repo, you get permission errors
- Make a commit, it has the wrong email
- Setup a new laptop, you spend 30 minutes configuring SSH

**Solution**: One command switches everything. Setup new machines in seconds.

## 🏗️ Building from Source

```bash
# Build for your current platform
go build -o sshm main.go

# Cross-compile for all platforms
make build-all

# Or manually:
GOOS=darwin GOARCH=amd64 go build -o sshm-mac main.go
GOOS=linux GOARCH=amd64 go build -o sshm-linux main.go
GOOS=windows GOARCH=amd64 go build -o sshm.exe main.go
```

## 📁 Project Structure

```
ssh-manager/
├── main.go          # Main application code
├── install.sh       # Unix/Linux/macOS installer
├── install.ps1      # Windows PowerShell installer
├── Makefile         # Build automation
└── README.md        # This file
```

## 🔐 Security

- SSH keys are stored in standard `~/.ssh/` directory
- Profile metadata stored in `~/.ssh-profiles` in your home directory
- No keys or sensitive data transmitted anywhere
- All operations are local to your machine

## 🤝 Contributing

Contributions welcome! Please feel free to submit a Pull Request.

## 📝 License

MIT License - feel free to use this in your own projects!

## 🐛 Troubleshooting

**Command not found after installation:**
- Restart your terminal
- Or run: `source ~/.zshrc` (or `~/.bashrc`)

**Permission denied when connecting:**
- Make sure you copied the public key to your Git provider
- Run `sshm current` to verify the correct profile is active
- Check `ssh -T git@github.com` to test connection

**Keys not switching:**
- Run `ssh-add -D` to clear SSH agent
- Then run `sshm switch <profile>` again

## 🌟 Star History

If this tool saves you time, give it a star! ⭐

---

Made with ❤️ for developers tired of SSH configuration headaches