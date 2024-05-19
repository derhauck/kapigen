
# Versioning

Kapigen automates the versioning process to ensure consistency and efficiency. Here’s how it works:

### Feature Branch Versioning

- **Automatic Tagging**: When you create a feature branch, Kapigen fetches the latest Git tag and appends your branch name to it. For example, if the latest tag is `v1.2.3` and your branch is `feature-xyz`, the version will be `v1.2.3-feature-xyz`.
- **Building**: This tagged version is used for all builds and tests in the feature branch pipeline.

### Merging into Main

- **Version Bump Based on Labels**: When you merge a feature branch into the main branch, Kapigen checks the labels on your merge request to determine how to bump the version:
    - `version::major` - Increases the major version (e.g., `v1.2.3` to `v2.0.0`).
    - `version::minor` - Increases the minor version (e.g., `v1.2.3` to `v1.3.0`).
    - `version::patch` - Increases the patch version (e.g., `v1.2.3` to `v1.2.4`).

- **New Tag Creation**: Kapigen creates a new tag with the updated version on the main branch.

### Release Process

- **Triggering Release Jobs**: When a new tag is created, all release jobs are triggered automatically. This includes building Docker images and pushing them to the registry with the new version tag.

### Optional Manual Versioning

If you prefer to manage versioning manually, you can disable the automated versioning job by setting `versioning: false` in your configuration.

## Example Workflow

1. **Merge Request**:
    - The current version is appended with the feature branch name.
    - For example, merging `feature-xyz` into `main` with the current version `v1.2.3` will create `v1.2.3-feature-xyz`.

2. **Main (optional)**:
    - Will tag with a new version based on the labels in the merge request.

3. **Tag**:
    - Creates new artifacts for the release.

By following these steps, Kapigen ensures a streamlined and consistent versioning and release process for your projects.
