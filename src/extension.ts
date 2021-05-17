// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from "vscode";
import * as os from "os";
import * as fs from "fs";
import { LanguageClientOptions, LanguageClient, ServerOptions, TransportKind } from "vscode-languageclient/node";

let client: LanguageClient;
let previewHtml: string = "No preview available";
let previewPanel: vscode.WebviewPanel | null = null;

export function activate(context: vscode.ExtensionContext) {

	// Select correct language server binary for this platform.
	let platform = `${os.platform()}-${os.arch()}`;
	let binPath = context.asAbsolutePath(`out/bin/tadl-${platform}`);
	if (!fs.existsSync(binPath)) {
		vscode.window.showErrorMessage(`tadl-support has no binary for platform "${platform}" and will not work. Contact the developer to fix this.`);
		return;
	}

	let serverOptions: ServerOptions = {
		command: binPath,
		transport: TransportKind.stdio,
	};
	let clientOptions: LanguageClientOptions = {
		documentSelector: [{scheme: "file", language: "tadl"}],
	};
	client = new LanguageClient(
		"tadl-language-server",
		"Tadl Language Server",
		serverOptions,
		clientOptions
	);
	client.start();

	// Setup client listeners, once it is ready.
	client.onReady().then(() => {

		// Store HTML preview we got from the server.
		client.onNotification("custom/preview", (html) => {
			previewHtml = html;
			// If we have a panel, set previous HTML
			if (previewPanel !== null) {
				previewPanel.webview.html = previewHtml;
			}
		});
		
	});

	// Command for opening the preview panel.
	context.subscriptions.push(vscode.commands.registerCommand("tadl.previewWorkspace", () => {
		if (previewPanel === null) {
			// Create a new panel if it was not open
			previewPanel = vscode.window.createWebviewPanel(
				"tadl.previewPanel",
				"Tadl Preview",
				vscode.ViewColumn.Beside,
				{}
			);
			previewPanel.onDidDispose(() => {
				previewPanel = null;
			});
			previewPanel.webview.html = previewHtml;
		} else {
			// Bring old panel to front
			previewPanel.reveal();
		}
	}));
}

// Shut down language server and close preview panels when extension is deactivated
export function deactivate(): Thenable<void> | undefined {
	if (previewPanel !== null) {
		previewPanel.dispose();
	}
	if (!client) {
		return undefined;
	} else {
		return client.stop();
	}
}
