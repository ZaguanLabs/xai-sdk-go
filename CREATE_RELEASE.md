# Creating GitHub Release v0.1.1

The code has been committed and tagged. To complete the release, create a GitHub release manually:

## ✅ Completed Steps

1. ✅ Updated version to 0.1.1 in all files
2. ✅ Updated CHANGELOG.md with release notes
3. ✅ Committed all changes
4. ✅ Pushed to main branch
5. ✅ Created and pushed v0.1.1 tag

## 📋 Next Steps - Create GitHub Release

### Option 1: Using GitHub Web UI (Recommended)

1. Go to: https://github.com/ZaguanLabs/xai-sdk-go/releases/new

2. **Tag**: Select `v0.1.1` from the dropdown

3. **Release title**: `v0.1.1 - Bug Fix Release`

4. **Description**: Copy the contents from `RELEASE_NOTES_v0.1.1.md`

5. Click **"Publish release"**

### Option 2: Using GitHub CLI (if installed)

```bash
gh release create v0.1.1 \
  --title "v0.1.1 - Bug Fix Release" \
  --notes-file RELEASE_NOTES_v0.1.1.md
```

### Option 3: Using API

```bash
curl -X POST \
  -H "Authorization: token YOUR_GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/ZaguanLabs/xai-sdk-go/releases \
  -d @- <<EOF
{
  "tag_name": "v0.1.1",
  "name": "v0.1.1 - Bug Fix Release",
  "body": "$(cat RELEASE_NOTES_v0.1.1.md | jq -Rs .)",
  "draft": false,
  "prerelease": false
}
EOF
```

## 📦 What's Included in v0.1.1

### Bug Fixes
- Fixed models API proto definitions (package name, RPC methods, field numbers)
- Fixed metadata handling to preserve gRPC headers
- Removed manual content-type interceptor

### Improvements
- Enhanced models example with debug logging
- Updated models client API with type-specific methods

### Verified Working
- ✅ Lists 8 language models (grok-2, grok-3, grok-4 variants)
- ✅ Retrieves detailed model information
- ✅ Proper authentication and metadata handling

## 🔗 Links

- **Repository**: https://github.com/ZaguanLabs/xai-sdk-go
- **Tag**: https://github.com/ZaguanLabs/xai-sdk-go/releases/tag/v0.1.1
- **Comparison**: https://github.com/ZaguanLabs/xai-sdk-go/compare/v0.1.0...v0.1.1

## 📝 Release Notes

The full release notes are in `RELEASE_NOTES_v0.1.1.md`
