package main

import (
	"bufio"
	"dyml-support"
	"dyml-support/protocol"
	"encoding/json"
	"log"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	server := dyml.NewServer()

	// Continuously read and respond to requests
	for {
		request, err := readRequest(reader)
		if err != nil {
			log.Println("Error while reading request:", err)
			continue
		}

		methodNameRaw, ok := request["method"]
		if !ok {
			log.Println("Got request with no method name!")
			continue
		}
		var methodName string
		if err := json.Unmarshal(methodNameRaw, &methodName); err != nil {
			log.Println("Error while unmarshalling method name:", err)
			continue
		}

		log.Printf("Got request with method '%s'\n", methodName)

		// Find id from request.
		// Notifications do not have an ID, so requestId might be 0.
		var requestId float64
		if requestIdRaw, ok := request["id"]; ok {
			if err := json.Unmarshal(requestIdRaw, &requestId); err != nil {
				log.Println("Error while unmarshalling id:", err)
			}
		}

		// Call the correct method on the server.
		switch methodName {
		case "initialize":
			var params protocol.InitializeParams
			if err := json.Unmarshal(request["params"], &params); err != nil {
				log.Println(err)
				continue
			}
			sendResponse(server.Initialize(&params), requestId)
		case "initialized":
			server.Initialized()
		case "$/cancelRequest":
			// Cancelling a request only makes sense for a multithreaded server.
		case "textDocument/hover":
			var params protocol.HoverParams
			if err := json.Unmarshal(request["params"], &params); err != nil {
				log.Println(err)
				continue
			}
			sendResponse(server.Hover(&params), requestId)
		case "textDocument/didSave":
			var params protocol.DidSaveTextDocumentParams
			if err := json.Unmarshal(request["params"], &params); err != nil {
				log.Println(err)
				continue
			}
			server.DidSaveTextDocument(&params)
		case "textDocument/didOpen":
			var params protocol.DidOpenTextDocumentParams
			if err := json.Unmarshal(request["params"], &params); err != nil {
				log.Println(err)
				continue
			}
			server.DidOpenTextDocument(&params)
		case "textDocument/didClose":
			var params protocol.DidCloseTextDocumentParams
			if err := json.Unmarshal(request["params"], &params); err != nil {
				log.Println(err)
				continue
			}
			server.DidCloseTextDocument(&params)
		case "textDocument/didChange":
			var params protocol.DidChangeTextDocumentParams
			if err := json.Unmarshal(request["params"], &params); err != nil {
				log.Println(err)
				continue
			}
			server.DidChangeTextDocument(&params)
		case "textDocument/semanticTokens/full":
			var params protocol.SemanticTokensParams
			if err := json.Unmarshal(request["params"], &params); err != nil {
				log.Println(err)
				continue
			}
			sendResponse(server.FullSemanticTokens(&params), requestId)
		case "custom/encodeXML":
			log.Println(string(request["params"]))
			var params []protocol.DocumentURI
			if err := json.Unmarshal(request["params"], &params); err != nil {
				log.Println(err)
				continue
			}
			sendResponse(server.EncodeXML(params[0]), requestId)
		default:
			log.Printf("Unknown method '%s'\n", methodName)
		}
	}
}

// Send response or log error message.
func sendResponse(response interface{}, requestId float64) {
	if err := dyml.SendResponse(response, requestId); err != nil {
		log.Println(err)
	}
}

// Read and parse a LSP request from stdin.
func readRequest(reader *bufio.Reader) (map[string]json.RawMessage, error) {
	// The request has a header and a JSON body.
	// The json.Decoder is smart enough read JSON without knowing the amount of expected bytes
	// from the header. So we just skip the header here, by reading until we encounter '\r\n'.
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if line == "\r\n" {
			// End of headers
			break
		}
	}

	decoder := json.NewDecoder(reader)
	var request map[string]json.RawMessage

	err := decoder.Decode(&request)
	if err != nil {
		return nil, err
	}

	return request, nil
}
