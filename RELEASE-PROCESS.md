# Release Process for AI Terminal Helper

This document describes the complete release process using Pull Requests.

---

## 🔄 Release Workflow

```
1. Create release branch
2. Run release script
3. Create Pull Request
4. Review & approve PR
5. Merge to main
6. Create git tag
7. Push tag (triggers GitHub Actions)
8. GitHub Actions creates release automatically
9. Announce
```

---

## 📋 Step-by-Step Guide

### Step 1: Create Release Branch

```bash
# Make sure you're on main and up to date
git checkout main
git pull origin main

# Create release branch
git checkout -b release/v2.0.0
```

### Step 2: Prepare Release

```bash
# Update CHANGELOG.md with release date
# Review all changes since last release
vim CHANGELOG.md

# Commit changelog
git add CHANGELOG.md
git commit -m "docs: update changelog for v2.0.0"
```

### Step 3: Run Release Script

```bash
# This will:
# - Update VERSION file
# - Build for all platforms
# - Create checksums
# - Commit VERSION file
./scripts/release.sh 2.0.0
```

**Output:**
```
✅ Release v2.0.0 prepared successfully!

📝 Next steps:
  1. Push release branch
  2. Create Pull Request
  3. After PR merged, create tag
```

### Step 4: Push Release Branch

```bash
git push origin release/v2.0.0
```

### Step 5: Create Pull Request

1. Go to: https://github.com/amaslovskyi/ai-helper/compare
2. **Base:** `main`
3. **Compare:** `release/v2.0.0`
4. **Title:** `Release v2.0.0`
5. **Description:** Copy from `.github/RELEASE.md`

**PR Description Template:**
```markdown
# Release v2.0.0

## 🎉 What's New

### ✨ Command Validation (The Big Feature!)
Prevents AI hallucinations by validating suggested commands before showing them.

### 🚀 Performance
- 10x faster than bash version
- 5ms startup time
- Single 5.5MB binary

### 🔒 Security
- 18 dangerous patterns detected
- 100% local execution
- Safe for secrets

## 📋 Checklist

- [x] CHANGELOG.md updated
- [x] VERSION file updated
- [x] All tests pass
- [x] Documentation updated
- [x] Builds successfully on all platforms

## 🔗 Links

- Changelog: [CHANGELOG.md](./CHANGELOG.md)
- Release Notes: [.github/RELEASE.md](./.github/RELEASE.md)
- Documentation: [README.md](./README.md)
```

### Step 6: Review & Approve PR

**Review Checklist:**
- [ ] CHANGELOG.md has all changes
- [ ] VERSION file is correct
- [ ] Documentation is up to date
- [ ] All tests pass (CI/CD)
- [ ] Builds succeed for all platforms

**Approve and merge the PR** (use "Squash and merge" or "Create a merge commit")

### Step 7: Create Git Tag (After PR Merged)

```bash
# Switch to main and pull latest
git checkout main
git pull origin main

# Verify you're on the right commit
git log -1

# Create annotated tag
git tag -a v2.0.0 -m "Release v2.0.0

✨ Command Validation - Prevents AI hallucinations
🚀 10x Faster Performance
🔒 100% Local & Secure

See CHANGELOG.md for full details.
"

# Verify tag
git tag -l -n9 v2.0.0
```

### Step 8: Push Tag (Triggers GitHub Actions)

```bash
# Push tag to GitHub
git push origin v2.0.0
```

**This triggers GitHub Actions workflow which will:**
1. Build binaries for all platforms
2. Create checksums
3. Create GitHub release
4. Upload binaries and checksums

**Monitor the workflow:**
- Go to: https://github.com/amaslovskyi/ai-helper/actions
- Watch the "Release" workflow

### Step 9: Verify Release

1. Go to: https://github.com/amaslovskyi/ai-helper/releases
2. Verify release `v2.0.0` is published
3. Check all binaries are uploaded:
   - `ai-helper-darwin-amd64`
   - `ai-helper-darwin-arm64`
   - `ai-helper-linux-amd64`
   - `ai-helper-linux-arm64`
   - `checksums.txt`

### Step 10: Test Release Binaries

```bash
# Test macOS Intel binary
curl -L -o ai-helper https://github.com/amaslovskyi/ai-helper/releases/download/v2.0.0/ai-helper-darwin-amd64
chmod +x ai-helper
./ai-helper version
# Should output: AI Terminal Helper 2.0.0-go

# Test macOS Apple Silicon binary
curl -L -o ai-helper https://github.com/amaslovskyi/ai-helper/releases/download/v2.0.0/ai-helper-darwin-arm64
chmod +x ai-helper
./ai-helper version
```

### Step 11: Announce Release

**Social Media (Twitter, LinkedIn):**
```
🚀 AI Terminal Helper v2.0.0 is here!

✨ Prevents AI hallucinations with command validation
⚡ 10x faster than bash version
🔒 100% local, no cloud, no telemetry
📦 Single 5.5MB binary

Perfect for DevOps/SRE/MLOps!

#golang #devops #ai #terminal #opensource

https://github.com/amaslovskyi/ai-helper/releases/tag/v2.0.0
```

**Reddit (r/devops, r/golang, r/commandline):**
- See `RELEASE-SUMMARY.md` for full template

**Dev.to / Hashnode:**
- Write a blog post about the release
- Include examples and screenshots

---

## 🔄 Hotfix Process

If a critical bug is found after release:

### 1. Create Hotfix Branch
```bash
git checkout main
git pull origin main
git checkout -b hotfix/v2.0.1
```

### 2. Fix the Bug
```bash
# Make your fixes
git add .
git commit -m "fix: critical bug description"
```

### 3. Follow Same Release Process
```bash
./scripts/release.sh 2.0.1
git push origin hotfix/v2.0.1
# Create PR, merge, tag, push
```

---

## 📊 Release Checklist

### Pre-Release (1-2 days before)
- [ ] All tests pass
- [ ] Documentation updated
- [ ] CHANGELOG.md complete
- [ ] Test on macOS (Intel & Apple Silicon)
- [ ] Test on Linux (if possible)

### Release Day
- [ ] Create release branch
- [ ] Update CHANGELOG.md with date
- [ ] Run `./scripts/release.sh X.Y.Z`
- [ ] Push release branch
- [ ] Create Pull Request
- [ ] Review PR
- [ ] Merge PR to main
- [ ] Create git tag
- [ ] Push tag
- [ ] Verify GitHub Actions workflow
- [ ] Verify release is published
- [ ] Test release binaries
- [ ] Announce release

### Post-Release
- [ ] Monitor GitHub issues
- [ ] Respond to questions
- [ ] Update documentation if needed
- [ ] Plan next release

---

## 🎯 Version Numbering

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (2.0.0): Breaking changes
  - Example: Rewrite from bash to Go
  
- **MINOR** (2.1.0): New features, backwards-compatible
  - Example: Add kubectl validator
  
- **PATCH** (2.0.1): Bug fixes, backwards-compatible
  - Example: Fix cache corruption

---

## 🤖 GitHub Actions Workflow

The `.github/workflows/release.yml` workflow:

**Triggers:**
- On push of tags matching `v*` pattern

**Steps:**
1. Checkout code
2. Set up Go 1.21+
3. Build for all platforms (darwin/linux, amd64/arm64)
4. Create SHA256 checksums
5. Create GitHub release with:
   - Release notes from `.github/RELEASE.md`
   - All binaries
   - Checksums file

**No manual binary upload needed!** 🎉

---

## 🔧 Troubleshooting

### Issue: GitHub Actions Fails

**Check:**
1. Go version in workflow (should be 1.21+)
2. Build commands in Makefile
3. GitHub Actions logs

**Fix:**
```bash
# Test build locally
make build-all

# If it works locally, check workflow file
cat .github/workflows/release.yml
```

### Issue: Tag Already Exists

**Fix:**
```bash
# Delete local tag
git tag -d v2.0.0

# Delete remote tag
git push origin :refs/tags/v2.0.0

# Recreate tag
git tag -a v2.0.0 -m "Release v2.0.0"
git push origin v2.0.0
```

### Issue: Release Not Created

**Check:**
1. GitHub Actions workflow status
2. GitHub token permissions
3. Tag format (must be `vX.Y.Z`)

---

## 📝 Summary

**Quick Release (After PR Merged):**
```bash
git checkout main
git pull origin main
git tag -a v2.0.0 -m "Release v2.0.0"
git push origin v2.0.0
# GitHub Actions handles the rest!
```

**That's it!** 🚀

---

Built with ❤️ in Go. No more hallucinations! 🎉

