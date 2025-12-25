#!/bin/bash
# Release script for AI Terminal Helper
# Usage: ./scripts/release.sh [version]

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
RESET='\033[0m'

# Get version
VERSION=${1:-$(cat VERSION)}
if [[ -z "$VERSION" ]]; then
    echo -e "${RED}❌ Error: No version specified${RESET}"
    echo "Usage: ./scripts/release.sh [version]"
    exit 1
fi

echo -e "${BLUE}🚀 Preparing release v${VERSION}${RESET}"
echo ""

# Verify we're on a release branch
CURRENT_BRANCH=$(git branch --show-current)
if [[ "$CURRENT_BRANCH" == "main" || "$CURRENT_BRANCH" == "master" ]]; then
    echo -e "${RED}❌ Error: Cannot run release script on main/master branch${RESET}"
    echo "Please create a release branch first:"
    echo "  git checkout -b release/v${VERSION}"
    exit 1
fi

# Check for uncommitted changes
if [[ -n $(git status -s) ]]; then
    echo -e "${RED}❌ Error: Uncommitted changes detected${RESET}"
    echo "Please commit or stash your changes first."
    exit 1
fi

# Update VERSION file
echo "$VERSION" > VERSION
echo -e "${GREEN}✅ Updated VERSION file${RESET}"

# Build for all platforms
echo -e "${BLUE}🔨 Building for all platforms...${RESET}"
make build-all

# Verify builds
if [[ ! -f "bin/ai-helper-darwin-amd64" ]] || \
   [[ ! -f "bin/ai-helper-darwin-arm64" ]] || \
   [[ ! -f "bin/ai-helper-linux-amd64" ]] || \
   [[ ! -f "bin/ai-helper-linux-arm64" ]]; then
    echo -e "${RED}❌ Error: Build failed${RESET}"
    exit 1
fi
echo -e "${GREEN}✅ All platforms built successfully${RESET}"

# Create checksums
echo -e "${BLUE}🔐 Creating checksums...${RESET}"
cd bin
sha256sum ai-helper-* > checksums.txt
cd ..
echo -e "${GREEN}✅ Checksums created${RESET}"

# Show build info
echo ""
echo -e "${BLUE}📊 Build Information:${RESET}"
ls -lh bin/ai-helper-*
echo ""
cat bin/checksums.txt
echo ""

# Commit VERSION file
echo -e "${BLUE}💾 Committing VERSION file...${RESET}"
git add VERSION
git commit -m "chore: bump version to ${VERSION}" || true
echo -e "${GREEN}✅ Version committed${RESET}"

# Summary
echo ""
echo -e "${GREEN}${BOLD}✅ Release v${VERSION} prepared successfully!${RESET}"
echo ""
echo -e "${BLUE}📝 Next steps:${RESET}"
echo ""
echo "  1. Push release branch:"
echo "     git push origin ${CURRENT_BRANCH}"
echo ""
echo "  2. Create Pull Request:"
echo "     - Go to: https://github.com/amaslovskyi/ai-helper/compare"
echo "     - Base: main"
echo "     - Compare: ${CURRENT_BRANCH}"
echo "     - Title: Release v${VERSION}"
echo "     - Description: Copy from .github/RELEASE.md"
echo ""
echo "  3. After PR is merged to main:"
echo "     git checkout main"
echo "     git pull origin main"
echo "     git tag -a v${VERSION} -m \"Release v${VERSION}\""
echo "     git push origin v${VERSION}"
echo ""
echo "  4. GitHub Actions will automatically:"
echo "     - Build binaries for all platforms"
echo "     - Create GitHub release"
echo "     - Upload binaries and checksums"
echo ""
echo "  5. Announce the release:"
echo "     - Post on social media"
echo "     - Update documentation"
echo ""
echo -e "${YELLOW}💡 Tip: Review RELEASE-SUMMARY.md for announcement templates${RESET}"

