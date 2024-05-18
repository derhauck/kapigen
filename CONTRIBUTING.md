# Contributing to Kapigen

Thank you for considering contributing to Kapigen! We welcome contributions from everyone. Below are the guidelines for contributing to this project.

## How to Contribute

1. **Fork the Repository**: Create a fork of the repository to your own GitLab account.

2. **Clone the Repository**:
   ```shell
   git clone https://gitlab.com/your-username/kapigen.git
   cd kapigen
   ```

3. **Create a Feature Branch**:
    - Branch names should be descriptive and include the issue number if applicable.
    - Use the following format for branch names: `feature/short-description-issue-number`.

   ```shell
   git checkout -b feature/my-new-feature-123
   ```

4. **Work on the Issue**:
    - Make your changes, ensuring that your code follows the project's style and guidelines.
    - Write tests for your code if applicable.
    - Ensure all tests pass before committing your changes.

5. **Commit Your Changes**:
    - Write clear and concise commit messages.
    - Reference the issue number in your commit message if applicable.

   ```shell
   git add .
   git commit -m "Add new feature to address issue #123"
   ```

6. **Push Your Branch**:
   ```shell
   git push origin feature/my-new-feature-123
   ```

7. **Create a Merge Request (MR)**:
    - Go to the original repository on GitLab.
    - You will see a prompt to create a merge request from your recently pushed branch.
    - Provide a clear title and description for your MR. Mention the issue number in the description if applicable.
    - Ensure the MR targets the correct branch (`main`).

## Review Process

- **Review**: Once you create an MR, a maintainer will review your code. They may request changes or provide feedback.
- **CI/CD**: Ensure that your MR passes all continuous integration checks.
- **Approval**: A maintainer needs to approve your MR before it can be merged.
- **Merge**: Once approved, a maintainer will merge your MR into the main branch.

## Coding Guidelines

- **Style**: Follow the coding style conventions used in the project.
- **Documentation**: Update documentation if your changes affect the public API or important workflows.
- **Tests**: Write tests to cover new functionality and ensure existing tests pass. (project coverage shouldn't be lower after your changes)

## Getting Help

If you have any questions or need help with your contribution, feel free to ask for help by opening an issue or by reaching out to the maintainers.

## Code of Conduct

Please note that this project adheres to a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## License

By contributing to Kapigen, you agree that your contributions will be licensed under the [GNU GPLv3](LICENSE).

Thank you for contributing to Kapigen!