package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type SSHProfile struct {
	Name     string
	Email    string
	KeyPath  string
	Username string
	Host     string
}

const (
	configFile = ".ssh-profiles"
	sshDir     = ".ssh"
	version    = "1.0.0"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "new":
		createNewProfile()
	case "list":
		listProfiles()
	case "switch":
		if len(os.Args) < 3 {
			fmt.Println("Usage: sshm switch <profile-name>")
			return
		}
		switchProfile(os.Args[2])
	case "current":
		showCurrentProfile()
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: sshm delete <profile-name>")
			return
		}
		deleteProfile(os.Args[2])
	case "version", "-v", "--version":
		fmt.Printf("SSH Manager v%s\n", version)
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("🔑 SSH Key Manager v" + version)
	fmt.Println("\nUsage:")
	fmt.Println("  sshm new              - Create a new SSH profile")
	fmt.Println("  sshm list             - List all profiles")
	fmt.Println("  sshm switch <name>    - Switch to a profile")
	fmt.Println("  sshm current          - Show current active profile")
	fmt.Println("  sshm delete <name>    - Delete a profile")
	fmt.Println("  sshm version          - Show version")
	fmt.Println("\nExamples:")
	fmt.Println("  sshm new")
	fmt.Println("  sshm switch work")
	fmt.Println("  sshm switch personal")
}

func getHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}
	return home
}

func getSSHDir() string {
	return filepath.Join(getHomeDir(), sshDir)
}

func getConfigPath() string {
	return filepath.Join(getHomeDir(), configFile)
}

func createNewProfile() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Profile name (e.g., 'work', 'personal'): ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	if name == "" {
		fmt.Println("❌ Profile name cannot be empty")
		return
	}

	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Git username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Git host (github.com, gitlab.com, etc.): ")
	host, _ := reader.ReadString('\n')
	host = strings.TrimSpace(host)

	// Create SSH directory if it doesn't exist
	sshDir := getSSHDir()
	os.MkdirAll(sshDir, 0700)

	// Generate SSH key
	keyPath := filepath.Join(sshDir, fmt.Sprintf("id_rsa_%s", name))
	
	fmt.Printf("\n🔨 Generating SSH key for %s...\n", name)
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-C", email, "-f", keyPath, "-N", "")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Error generating SSH key:", err)
		return
	}

	// Save profile
	profile := SSHProfile{
		Name:     name,
		Email:    email,
		KeyPath:  keyPath,
		Username: username,
		Host:     host,
	}

	saveProfile(profile)
	
	fmt.Printf("\n✅ Profile '%s' created successfully!\n", name)
	fmt.Printf("📋 Public key location: %s.pub\n", keyPath)
	fmt.Println("\n📎 Copy this public key to your Git provider:")
	
	// Display public key
	pubKeyContent, err := os.ReadFile(keyPath + ".pub")
	if err != nil {
		fmt.Println("❌ Error reading public key:", err)
		return
	}
	
	fmt.Println("─────────────────────────────────────────")
	fmt.Print(string(pubKeyContent))
	fmt.Println("─────────────────────────────────────────")
	
	fmt.Printf("\n💡 Run 'sshm switch %s' to activate this profile\n", name)
}

func saveProfile(profile SSHProfile) {
	configPath := getConfigPath()
	
	// Read existing profiles
	profiles := loadProfiles()
	
	// Add or update profile
	found := false
	for i, p := range profiles {
		if p.Name == profile.Name {
			profiles[i] = profile
			found = true
			break
		}
	}
	
	if !found {
		profiles = append(profiles, profile)
	}
	
	// Write profiles back
	file, err := os.Create(configPath)
	if err != nil {
		fmt.Println("❌ Error saving profile:", err)
		return
	}
	defer file.Close()
	
	for _, p := range profiles {
		fmt.Fprintf(file, "%s|%s|%s|%s|%s\n", p.Name, p.Email, p.KeyPath, p.Username, p.Host)
	}
}

func loadProfiles() []SSHProfile {
	configPath := getConfigPath()
	var profiles []SSHProfile
	
	file, err := os.Open(configPath)
	if err != nil {
		return profiles
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) == 5 {
			profile := SSHProfile{
				Name:     parts[0],
				Email:    parts[1],
				KeyPath:  parts[2],
				Username: parts[3],
				Host:     parts[4],
			}
			profiles = append(profiles, profile)
		}
	}
	
	return profiles
}

func listProfiles() {
	profiles := loadProfiles()
	
	if len(profiles) == 0 {
		fmt.Println("No profiles found. Create one with 'sshm new'")
		return
	}
	
	fmt.Println("📋 Available SSH Profiles:")
	fmt.Println("─────────────────────────")
	
	current := getCurrentProfile()
	
	for _, profile := range profiles {
		active := ""
		if profile.Name == current {
			active = " ⭐ (active)"
		}
		fmt.Printf("📝 %s%s\n", profile.Name, active)
		fmt.Printf("   Email: %s\n", profile.Email)
		fmt.Printf("   Username: %s\n", profile.Username)
		fmt.Printf("   Host: %s\n", profile.Host)
		fmt.Printf("   Key: %s\n\n", profile.KeyPath)
	}
}

func switchProfile(profileName string) {
	profiles := loadProfiles()
	
	var selectedProfile *SSHProfile
	for _, profile := range profiles {
		if profile.Name == profileName {
			selectedProfile = &profile
			break
		}
	}
	
	if selectedProfile == nil {
		fmt.Printf("❌ Profile '%s' not found. Use 'sshm list' to see available profiles.\n", profileName)
		return
	}
	
	// Update SSH config
	updateSSHConfig(*selectedProfile)
	
	// Update Git config
	updateGitConfig(*selectedProfile)
	
	// Add key to SSH agent
	addKeyToAgent(selectedProfile.KeyPath)
	
	// Save current profile
	saveCurrentProfile(profileName)
	
	fmt.Printf("✅ Switched to profile: %s (%s)\n", selectedProfile.Name, selectedProfile.Email)
	fmt.Printf("🔧 Git configured as: %s <%s>\n", selectedProfile.Username, selectedProfile.Email)
}

func updateSSHConfig(profile SSHProfile) {
	sshDir := getSSHDir()
	configPath := filepath.Join(sshDir, "config")

	// Read existing config if it exists
	content, err := os.ReadFile(configPath)
	if err != nil {
		content = []byte("")
	}

	configStr := string(content)

	// Remove ALL existing Host entries for the target host (e.g., github.com)
	// This ensures no duplicate/conflicting entries remain
	lines := strings.Split(configStr, "\n")
	var cleanedLines []string
	skipUntilNextHost := false

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Check if this is a Host entry for our target host
		if strings.HasPrefix(trimmed, "Host ") {
			hostValue := strings.TrimSpace(strings.TrimPrefix(trimmed, "Host "))

			// If it matches our target host, skip this entire section
			if hostValue == profile.Host {
				skipUntilNextHost = true
				continue
			} else {
				// Different host, stop skipping
				skipUntilNextHost = false
				cleanedLines = append(cleanedLines, line)
			}
		} else if skipUntilNextHost {
			// Skip lines that are part of the host section we're removing
			// Stop skipping when we hit a comment that's not indented (new section)
			if strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
				skipUntilNextHost = false
				cleanedLines = append(cleanedLines, line)
			}
			// Otherwise continue skipping
		} else {
			cleanedLines = append(cleanedLines, line)
		}
	}

	// Clean up multiple consecutive newlines
	configStr = strings.Join(cleanedLines, "\n")
	for strings.Contains(configStr, "\n\n\n") {
		configStr = strings.ReplaceAll(configStr, "\n\n\n", "\n\n")
	}

	// Trim trailing newlines
	configStr = strings.TrimRight(configStr, "\n")

	// Add new SSH Manager section at the top for priority
	newSection := fmt.Sprintf(`# SSH Manager - Profile: %s
Host %s
    HostName %s
    User git
    IdentityFile %s
    IdentitiesOnly yes
    AddKeysToAgent yes`, profile.Name, profile.Host, profile.Host, profile.KeyPath)

	// Prepend the new section at the beginning for highest priority
	if len(configStr) > 0 {
		configStr = newSection + "\n\n" + configStr + "\n"
	} else {
		configStr = newSection + "\n"
	}

	// Write the updated config
	err = os.WriteFile(configPath, []byte(configStr), 0600)
	if err != nil {
		fmt.Println("❌ Error updating SSH config:", err)
	}
}

func removeSSHConfigEntry(profileName string) {
	sshDir := getSSHDir()
	configPath := filepath.Join(sshDir, "config")

	// Read existing config if it exists
	content, err := os.ReadFile(configPath)
	if err != nil {
		return // No config file exists, nothing to clean
	}

	existingContent := strings.Split(string(content), "\n")
	var cleanedContent []string
	skipUntilNextHost := false

	for _, line := range existingContent {
		trimmedLine := strings.TrimSpace(line)

		// Check if this is the SSH Manager managed section for the profile we're deleting
		if strings.HasPrefix(trimmedLine, "# SSH Manager - Profile: "+profileName) {
			skipUntilNextHost = true
			continue
		}

		// If we're skipping and find a new Host or comment that's not SSH Manager
		if skipUntilNextHost {
			if strings.HasPrefix(trimmedLine, "Host ") && !strings.Contains(line, "# SSH Manager") {
				skipUntilNextHost = false
				cleanedContent = append(cleanedContent, line)
			} else if strings.HasPrefix(trimmedLine, "#") && !strings.Contains(line, "SSH Manager") {
				skipUntilNextHost = false
				cleanedContent = append(cleanedContent, line)
			}
			// Skip lines that are part of the SSH Manager managed section we're deleting
			continue
		}

		cleanedContent = append(cleanedContent, line)
	}

	// Remove trailing empty lines
	for len(cleanedContent) > 0 && strings.TrimSpace(cleanedContent[len(cleanedContent)-1]) == "" {
		cleanedContent = cleanedContent[:len(cleanedContent)-1]
	}

	// Write the cleaned config back
	finalContent := strings.Join(cleanedContent, "\n")
	if len(cleanedContent) > 0 {
		finalContent += "\n"
	}
	os.WriteFile(configPath, []byte(finalContent), 0600)
}

func updateGitConfig(profile SSHProfile) {
	exec.Command("git", "config", "--global", "user.name", profile.Username).Run()
	exec.Command("git", "config", "--global", "user.email", profile.Email).Run()
}

func addKeyToAgent(keyPath string) {
	if runtime.GOOS == "windows" {
		exec.Command("ssh-add", keyPath).Run()
	} else {
		cmd := exec.Command("ssh-add", keyPath)
		cmd.Run()
	}
}

func saveCurrentProfile(profileName string) {
	currentPath := filepath.Join(getHomeDir(), ".ssh-current")
	os.WriteFile(currentPath, []byte(profileName), 0600)
}

func getCurrentProfile() string {
	currentPath := filepath.Join(getHomeDir(), ".ssh-current")
	content, err := os.ReadFile(currentPath)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(content))
}

func showCurrentProfile() {
	current := getCurrentProfile()
	if current == "" {
		fmt.Println("No active profile set.")
		fmt.Println("💡 Use 'sshm list' to see profiles and 'sshm switch <name>' to activate one.")
		return
	}
	
	profiles := loadProfiles()
	for _, profile := range profiles {
		if profile.Name == current {
			fmt.Printf("🔥 Current profile: %s\n", profile.Name)
			fmt.Printf("   Email: %s\n", profile.Email)
			fmt.Printf("   Username: %s\n", profile.Username)
			fmt.Printf("   Host: %s\n", profile.Host)
			return
		}
	}
	
	fmt.Printf("⚠️  Current profile '%s' not found in saved profiles.\n", current)
}

func deleteProfile(profileName string) {
	profiles := loadProfiles()
	
	var updatedProfiles []SSHProfile
	var deletedProfile *SSHProfile
	
	for _, profile := range profiles {
		if profile.Name == profileName {
			deletedProfile = &profile
		} else {
			updatedProfiles = append(updatedProfiles, profile)
		}
	}
	
	if deletedProfile == nil {
		fmt.Printf("❌ Profile '%s' not found.\n", profileName)
		return
	}
	
	os.Remove(deletedProfile.KeyPath)
	os.Remove(deletedProfile.KeyPath + ".pub")

	// Clean up SSH config entry for deleted profile
	removeSSHConfigEntry(deletedProfile.Name)

	configPath := getConfigPath()
	file, err := os.Create(configPath)
	if err != nil {
		fmt.Println("❌ Error updating profiles:", err)
		return
	}
	defer file.Close()
	
	for _, p := range updatedProfiles {
		fmt.Fprintf(file, "%s|%s|%s|%s|%s\n", p.Name, p.Email, p.KeyPath, p.Username, p.Host)
	}
	
	fmt.Printf("✅ Profile '%s' deleted successfully.\n", profileName)
}