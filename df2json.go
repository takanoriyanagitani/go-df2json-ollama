package df2json

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"unicode/utf8"

	oa "github.com/ollama/ollama/api"
)

var ErrInvalidDfString error = errors.New("invalid df string")

type BasicGenerateRequest struct {
	Model  string          `json:"model"`
	Prompt string          `json:"prompt"`
	Stream bool            `json:"stream"`
	Format json.RawMessage `json:"format"`
}

func (b BasicGenerateRequest) ToRequest() oa.GenerateRequest {
	return oa.GenerateRequest{
		Model:  b.Model,
		Prompt: b.Prompt,
		Stream: &b.Stream,
		Format: b.Format,
	}
}

const DfFormat string = `{
    "type": "object",
    "properties": {
        "filesystem":       {"type":"string"},
        "size":             {"type":"string"},
        "used":             {"type":"string"},
        "avail":            {"type":"string"},
        "capacity":         {"type":"string"},
        "iused":            {"type":"string"},
        "ifree":            {"type":"string"},
        "iused_percentage": {"type":"string"},
        "mounted_on":       {"type":"string"}
    },
    "required": ["filesystem", "size", "used", "avail", "capacity", "iused", "ifree", "iused_percentage", "mounted_on"]
}`

const PromptSuffixDefault string = " . Respond using JSON"

type RawDfString string

func (r RawDfString) ToPrompt(promptSuffix string) string {
	return string(r) + promptSuffix
}

func (r RawDfString) ToBasicRequestDefault(model string) BasicGenerateRequest {
	return BasicGenerateRequest{
		Model:  model,
		Prompt: r.ToPrompt(PromptSuffixDefault),
		Stream: false,
		Format: json.RawMessage(DfFormat),
	}
}

func (r RawDfString) ToRequestDefault(model string) oa.GenerateRequest {
	return r.ToBasicRequestDefault(model).ToRequest()
}

type Client struct{ *oa.Client }

func (c Client) Generate(
	ctx context.Context,
	req *oa.GenerateRequest,
) (res oa.GenerateResponse, e error) {
	e = c.Client.Generate(
		ctx,
		req,
		func(r oa.GenerateResponse) error {
			res = r
			return nil
		},
	)
	return
}

func (c Client) BasicGenerate(
	ctx context.Context,
	req BasicGenerateRequest,
) (oa.GenerateResponse, error) {
	greq := req.ToRequest()
	return c.Generate(ctx, &greq)
}

func (c Client) ParseDfString(
	ctx context.Context,
	df RawDfString,
	model string,
) (oa.GenerateResponse, error) {
	basicRequest := df.ToBasicRequestDefault(model)
	return c.BasicGenerate(ctx, basicRequest)
}

type RawDfSource func(context.Context) (RawDfString, error)

func RawDfSourceExec(ctx context.Context) (RawDfString, error) {
	out, e := exec.CommandContext(ctx, "df", "-h", ".").Output()
	if nil != e {
		return "", e
	}

	var s string = string(out)
	var valid bool = utf8.ValidString(s)
	if !valid {
		return "", fmt.Errorf("%w: len=%v", ErrInvalidDfString, len(out))
	}

	return RawDfString(s), nil
}

var RawDfSourceDefault RawDfSource = RawDfSourceExec

func (c Client) GetParsedDfDefault(
	ctx context.Context,
	model string,
) (oa.GenerateResponse, error) {
	rawDf, err := RawDfSourceDefault(ctx)
	if err != nil {
		return oa.GenerateResponse{}, fmt.Errorf("failed to retrieve disk usage data: %w", err)
	}

	return c.ParseDfString(ctx, rawDf, model)
}

func ResponseToJsonString(r oa.GenerateResponse) string {
	return r.Response
}
