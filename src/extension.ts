// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from "vscode";
import * as os from "os";
import * as fs from "fs";
import { LanguageClientOptions, LanguageClient, ServerOptions, TransportKind } from "vscode-languageclient/node";

let client: LanguageClient;

export function activate(context: vscode.ExtensionContext) {

	// Select correct language server binary for this platform.
	let platform = `${os.platform()}-${os.arch()}`;
	let binPath = context.asAbsolutePath(`out/bin/dyml-${platform}`);
	if (!fs.existsSync(binPath)) {
		vscode.window.showErrorMessage(`dyml-support has no binary for platform "${platform}" and will not work. Contact the developer to fix this.`);
		return;
	}

	let serverOptions: ServerOptions = {
		command: binPath,
		transport: TransportKind.stdio,
	};
	let clientOptions: LanguageClientOptions = {
		documentSelector: [{scheme: "file", language: "dyml"}],
	};
	client = new LanguageClient(
		"dyml-language-server",
		"DYML Language Server",
		serverOptions,
		clientOptions
	);
	client.start();

	// Request an XML preview from the language server and show that result in a new editor.
	context.subscriptions.push(vscode.commands.registerCommand("dyml.encodeXML", () => {
		let doc = "file://" + vscode.window.activeTextEditor?.document.uri.fsPath;
		if (doc) {
			client.sendRequest("custom/encodeXML", doc).then((resp) => {
				vscode.workspace.openTextDocument({
					content: String(resp),
					language: "xml"
				}).then((document) => {
					vscode.window.showTextDocument(document);
				});
			});
		}
	}));
}

// Shut down language server and close preview panels when extension is deactivated
export function deactivate(): Thenable<void> | undefined {
	if (!client) {
		return undefined;
	} else {
		return client.stop();
	}
}
