package tpl

import (
	"fmt"
	htmlTpl "html/template"
	"io"
	"io/fs"
	"os"
	txtTpl "text/template"
)

type Client struct {
	fs fs.FS
}

func NewClient(opts ...ClientOpt) *Client {
	c := &Client{
		fs: &osFS{},
	}

	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}

	return c
}

func (c *Client) ParseHTMLTemplateFromFile(path string) (*htmlTpl.Template, error) {
	s, err := c.readFileContents(path)
	if err != nil {
		return nil, err
	}

	tpl, err := htmlTpl.New("").Parse(s)
	if err != nil {
		return nil, fmt.Errorf("cannot parse as html template: file %q: %w", path, err)
	}

	return tpl, nil
}

func (c *Client) ParseTxtTemplateFromFile(path string) (*txtTpl.Template, error) {
	s, err := c.readFileContents(path)
	if err != nil {
		return nil, err
	}

	tpl, err := txtTpl.New("").Parse(s)
	if err != nil {
		return nil, fmt.Errorf("cannot parse as text template: file %q: %w", path, err)
	}

	return tpl, nil
}

func (c *Client) readFileContents(path string) (string, error) {
	f, err := c.fs.Open(path)
	if err != nil {
		return "", fmt.Errorf("cannot open: file %q: %w", path, err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("cannot read: file %q: %w", path, err)
	}

	return string(b), nil
}

type ClientOpt func(c *Client)

func WithFS(fs fs.FS) ClientOpt {
	return func(c *Client) {
		if c != nil && fs != nil {
			c.fs = fs
		}
	}
}

type osFS struct{}

func (o *osFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}
