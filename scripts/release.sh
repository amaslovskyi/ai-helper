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

# Verify we're on main/master branch
CURRENT_BRANCH=$(git branch --show-current)
if [[ "$CURRENT_BRANCH" != "main" && "$CURRENT_BRANCH" != "master" ]]; then
    echo -e "${YELLOW}⚠️  Warning: Not on main/master branch (current: ${CURRENT_BRANCH})${RESET}"
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
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

# Create git tag
echo -e "${BLUE}🏷️  Creating git tag v${VERSION}...${RESET}"
git add VERSION
git commit -m "Release v${VERSION}" || true
git tag -a "v${VERSION}" -m "Release v${VERSION}

See CHANGELOG.md for details.
"
echo -e "${GREEN}✅ Git tag created${RESET}"

# Summary
echo ""
echo -e "${GREEN}${BOLD}✅ Release v${VERSION} prepared successfully!${RESET}"
echo ""
echo -e "${BLUE}📝 Next steps:${RESET}"
echo "  1. Review the changes:"
echo "     git show v${VERSION}"
echo ""
echo "  2. Push the tag:"
echo "     git push origin v${VERSION}"
echo ""
echo "  3. Create GitHub release:"
echo "     - Go to: https://github.com/yourusername/ai-helper/releases/new"
echo "     - Tag: v${VERSION}"
echo "     - Title: AI Terminal Helper v${VERSION}"
echo "     - Description: Copy from .github/RELEASE.md"
echo "     - Upload binaries from bin/"
echo "     - Upload checksums.txt"
echo ""
echo "  4. Announce the release:"
echo "     - Update README.md badges"
echo "     - Post on social media"
echo "     - Update documentation"
echo ""
echo -e "${YELLOW}💡 Tip: Use 'git push --tags' to push all tags${RESET}"

