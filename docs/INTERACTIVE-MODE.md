# Interactive Mode Guide

## 🎯 Overview

Interactive Mode gives you **full control** over when and how AI assistance is triggered. Instead of automatically analyzing every error, you can choose to show an interactive menu, require manual activation, or even disable AI completely.

This feature is inspired by modern terminal tools like Warp, but maintains our 100% local, privacy-first approach.

---

## 🚀 Activation Modes

AI Helper supports **4 activation modes**:

### 1. Auto Mode (Default)
**Behavior:** AI triggers automatically on command failures  
**Best For:** Fast-paced work, maximum automation  
**Example:**
```bash
kubectl get pods --sort=name
# ❌ Error
# 🤖 AI automatically suggests fix
```

**Set:**
```bash
ai-helper config-set mode auto
```

---

### 2. Interactive Mode ⭐ NEW!
**Behavior:** Shows menu on failures, you choose the action  
**Best For:** Controlled assistance, learning, avoiding interruptions  
**Example:**
```bash
kubectl get pods --sort=name
# ❌ Error

🤖 Command failed. What would you like to do?

  [1] Get AI suggestion - Let AI analyze and suggest a fix
  [2] Show manual - Display manual page for this command
  [3] Skip - Continue without fixing
  [4] Disable AI for session - Turn off AI until terminal restart

Your choice: _
```

**Set:**
```bash
ai-helper config-set mode interactive
```

---

### 3. Manual Mode
**Behavior:** AI only activates with explicit commands (`ask`, `analyze`)  
**Best For:** Expert users who want AI on-demand only  
**Example:**
```bash
kubectl get pods --sort=name
# ❌ Error
# (No AI suggestion - error is shown as normal)

# You explicitly call AI when needed:
ask how do I sort kubectl pods by creation time
```

**Set:**
```bash
ai-helper config-set mode manual
```

---

### 4. Disabled Mode
**Behavior:** AI is completely turned off  
**Best For:** Temporarily disabling AI, troubleshooting  
**Example:**
```bash
kubectl get pods --sort=name
# ❌ Error
# (AI is completely disabled - no output, no menu)
```

**Set:**
```bash
ai-helper config-set mode disabled
```

---

## ⚙️ Configuration Commands

### Show Current Configuration
```bash
ai-helper config-show
```

**Output:**
```
⚙️  Configuration:
  Activation Mode: interactive
  Auto Execute Safe: false
  Show Confidence: true
  Tool-Specific Modes:
    kubectl: interactive
    docker: auto
```

---

### Set Activation Mode
```bash
# Set global mode
ai-helper config-set mode interactive
ai-helper config-set mode auto
ai-helper config-set mode manual
ai-helper config-set mode disabled
```

---

### Set Tool-Specific Modes
You can set different modes for different tools:

```bash
# Interactive mode for kubectl (safety-critical)
ai-helper config-set tool-mode kubectl interactive

# Auto mode for docker (less risky)
ai-helper config-set tool-mode docker auto

# Manual mode for terraform (explicit control)
ai-helper config-set tool-mode terraform manual
```

**Example Configuration:**
```json
{
  "activation_mode": "auto",
  "tool_specific_modes": {
    "kubectl": "interactive",
    "terraform": "interactive",
    "docker": "auto"
  }
}
```

**Behavior:**
- `kubectl` errors → Show interactive menu
- `terraform` errors → Show interactive menu  
- `docker` errors → Auto-trigger AI
- All other tools → Auto-trigger AI (global default)

---

### Toggle Other Settings
```bash
# Show/hide confidence scores
ai-helper config-set confidence true
ai-helper config-set confidence false
```

---

### Reset Configuration
```bash
ai-helper config-reset
```
Resets all settings to defaults (asks for confirmation).

---

## 📋 Interactive Menu Options

When in Interactive Mode, you get these choices:

### [1] Get AI suggestion
- Triggers AI analysis
- Shows suggested fix with explanation
- Same behavior as Auto mode

### [2] Show manual
- Displays tip to use `man <command>`
- Helpful for learning command syntax
- No AI query needed

### [3] Skip
- Ignores the error
- Continues your workflow
- No AI query, no action

### [4] Disable AI for session
- Turns off AI until terminal restart
- Temporary disable (doesn't save to config)
- Useful when working on scripts with expected failures

---

## 🎓 Use Cases & Best Practices

### For Senior Engineers
**Recommended:** `manual` or `interactive` mode
- Less interruption during focused work
- AI on-demand when needed
- Full control over assistance

```bash
ai-helper config-set mode manual
# Or
ai-helper config-set mode interactive
```

---

### For Junior Engineers / Learning
**Recommended:** `auto` or `interactive` mode
- Automatic help when errors occur
- Learn correct syntax from AI suggestions
- Build command muscle memory

```bash
ai-helper config-set mode interactive  # Learning with control
# Or  
ai-helper config-set mode auto  # Maximum assistance
```

---

### For Production/Critical Systems
**Recommended:** `interactive` mode with tool-specific overrides
- Manual confirmation for dangerous operations
- Auto-assist for safe commands
- Prevent accidental destructive commands

```bash
# Global interactive mode
ai-helper config-set mode interactive

# Extra caution for critical tools
ai-helper config-set tool-mode kubectl interactive
ai-helper config-set tool-mode terraform interactive
ai-helper config-set tool-mode ansible interactive
```

---

### For Scripting / Automation
**Recommended:** `disabled` or `manual` mode
- No AI interruption in scripts
- Predictable behavior
- No unexpected prompts

```bash
ai-helper config-set mode disabled
```

---

## 🔧 Technical Details

### Configuration File
Location: `~/.ai/config.json`

**Example:**
```json
{
  "activation_mode": "interactive",
  "auto_execute_safe": false,
  "show_confidence": true,
  "preferred_model": "",
  "tool_specific_modes": {
    "kubectl": "interactive",
    "terraform": "manual"
  }
}
```

### Session Overrides
When you select "Disable AI for session" from the menu, it sets an in-memory flag:
- **Not saved to disk** - resets on terminal restart
- **Temporary disable** - doesn't affect config file
- **Quick toggle** - no file editing needed

---

## 🆚 Comparison with Warp Terminal

| Feature              | AI Helper      | Warp Terminal       |
| -------------------- | -------------- | ------------------- |
| **Interactive Menu** | ✅ Yes          | ❌ No (auto only)    |
| **Per-Tool Modes**   | ✅ Yes          | ❌ No                |
| **Session Disable**  | ✅ Yes          | Limited             |
| **Privacy**          | ✅ 100% local   | ❌ Cloud-based       |
| **Offline**          | ✅ Yes          | ❌ Requires internet |
| **Configuration**    | ✅ Fine-grained | ⚠️ Limited options   |

**Our Advantage:** More control, better privacy, works offline.

---

## 📚 Examples

### Example 1: Interactive Workflow
```bash
# Set mode
ai-helper config-set mode interactive

# Try a bad command
terraform plan --apply

# Menu appears:
🤖 Command failed. What would you like to do?
[1] Get AI suggestion
[2] Show manual
[3] Skip
[4] Disable AI for session

Your choice: 1

# AI analyzes:
✓ terraform plan -out=plan.tfplan && terraform apply plan.tfplan
Root: The --apply flag doesn't exist in terraform plan
Tip: Use separate plan and apply commands for safety
Confidence: ✅ High (95%)
```

---

### Example 2: Tool-Specific Configuration
```bash
# Setup
ai-helper config-set mode auto
ai-helper config-set tool-mode kubectl interactive
ai-helper config-set tool-mode terraform interactive

# kubectl errors show menu (high risk)
kubectl delete deployment prod-app
# → Interactive menu appears

# docker errors auto-fix (low risk)
docker ps --format invalid
# → AI suggestion appears immediately
```

---

### Example 3: Learning Mode
```bash
# Set interactive mode for learning
ai-helper config-set mode interactive

# Make mistakes and learn
git push --force main
# Menu: [1] AI [2] sudo [3] manual [4] skip [5] disable
# Choose [1]
# AI explains: "Never force push to main! Use: git push origin feature-branch"

# Eventually switch to auto when confident
ai-helper config-set mode auto
```

---

## 🐛 Troubleshooting

### AI Not Triggering
**Check mode:**
```bash
ai-helper config-show
```
Make sure mode is not `disabled` or `manual`.

---

### Menu Not Appearing
**Verify interactive mode:**
```bash
ai-helper config-set mode interactive
```

---

### Reset Not Working
**Manual reset:**
```bash
rm ~/.ai/config.json
ai-helper config-show  # Recreates with defaults
```

---

## 🎯 Quick Reference

```bash
# View config
ai-helper config-show

# Set modes
ai-helper config-set mode auto|interactive|manual|disabled

# Tool-specific
ai-helper config-set tool-mode <tool> <mode>

# Other settings
ai-helper config-set confidence true|false

# Reset
ai-helper config-reset
```

---

## 🚀 What's Next?

Future enhancements planned:
- **Auto-execute safe commands** - Run read-only commands automatically
- **Custom menu options** - Add your own actions
- **Workflow detection** - Multi-step command sequences
- **Learning mode** - Track which suggestions you accept/reject

---

**Version:** v2.3.0  
**Author:** Andrii Maslovskyi  
**License:** MIT

