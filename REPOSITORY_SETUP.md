# Custom Repository Setup for go-pwr

go-pwr now supports using custom script repositories while keeping RocketPowerInc's scriptbin as the default.

## Quick Start

### Command Line Options

1. **View current repository:**

   ```bash
   go-pwr -show-repo
   ```

2. **Set a custom repository:**

   ```bash
   go-pwr -set-repo https://github.com/yourusername/your-scripts.git
   ```

3. **Reset to default repository:**
   ```bash
   go-pwr -reset-repo
   ```

### Through the UI

1. Start go-pwr normally
2. Navigate to the "Options" tab
3. Select "Repository Settings" from the left panel
4. Choose your action from the right panel:
   - **Set Custom Repository**: Instructions for setting up a custom repo
   - **Reset to Default**: Restore RocketPowerInc's scriptbin
   - **Current Repository**: View your current repository URL

## Repository Requirements

Your custom repository should:

1. **Be a Git repository** ending with `.git`
2. **Be publicly accessible** or you should have appropriate credentials configured
3. **Contain executable scripts** in a directory structure
4. **Use supported URL schemes**: https, http, git, or ssh

### Supported URL Formats

- `https://github.com/username/repo.git`
- `https://gitlab.com/username/repo.git`
- `git@github.com:username/repo.git`
- Any other valid Git repository URL

## Example Custom Repository Structure

```
your-scripts/
├── automation/
│   ├── backup.sh
│   └── deploy.py
├── utilities/
│   ├── file-organizer.ps1
│   └── system-info.sh
└── development/
    ├── setup-env.sh
    └── build-project.py
```

## Configuration Storage

Repository settings are stored in:

- **Linux/macOS**: `~/.config/go-pwr/config.json`
- **Windows**: `%USERPROFILE%\.config\go-pwr\config.json`

Example config file:

```json
{
  "theme": "Ocean Breeze",
  "repo_url": "https://github.com/yourusername/your-scripts.git"
}
```

## Notes

- **Default Behavior**: If no custom repository is set, go-pwr uses RocketPowerInc's scriptbin
- **Repository Validation**: URLs are validated before saving
- **Fresh Clone**: The repository is freshly cloned each time go-pwr starts
- **Multiple Repositories**: Different custom repositories are stored in separate directories

## Troubleshooting

1. **Invalid URL Error**: Ensure your repository URL ends with `.git` and uses a supported scheme
2. **Clone Failed**: Check that the repository is accessible and you have proper credentials
3. **No Scripts Found**: Verify your repository contains executable scripts in the expected structure

## Security Note

Only use repositories you trust, as scripts will be executed on your system.
