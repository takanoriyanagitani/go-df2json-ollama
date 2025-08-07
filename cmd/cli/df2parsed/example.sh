#!/bin/sh

export OLLAMA_MODEL=phi4-mini:3.8b          # less correct
export OLLAMA_MODEL=mistral:7b              #
export OLLAMA_MODEL=phi4:14b                #
export OLLAMA_MODEL=granite3.3:8b           #
export OLLAMA_MODEL=gemma3:12b              #
export OLLAMA_MODEL=mistral-nemo:12b        #
export OLLAMA_MODEL=magistral:24b           #
export OLLAMA_MODEL=nemotron-mini:4b        #
export OLLAMA_MODEL=phi4-reasoning:14b      #
export OLLAMA_MODEL=llama3-groq-tool-use:8b #
export OLLAMA_MODEL=mixtral:8x7b            #
export OLLAMA_MODEL=llama3.1:8b             #

export OLLAMA_MODEL=llama3.2:3b             #

df -h .

echo
./df2parsed |
	jq
