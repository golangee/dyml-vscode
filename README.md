# tadl-support

A VS code extension to add support for the [Tadl-specification](https://github.com/golangee/tadl/blob/main/specification.adoc) (name pending). It will feature syntax highlighting, checking for semantic errors and previewing the current workspace's documentation.

Currently WIP.

## Development
You need [VS Code](https://code.visualstudio.com/) to try this extension. Open this project in it and select `Run > Start Debugging` or press `F5`. The go compiler needs to be installed for development.

To package the application into a `.vsix` file run `vsce package`. You might need to `npm install -g vsce` first. This package can then be installed in vscode by opening the overflow menu in the extension tab, and selecting `Install from VSIX`.